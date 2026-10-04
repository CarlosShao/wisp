#Requires -Version 5.1
<#
.SYNOPSIS
    SLO gate driver (ticket 08, SPEC-10 §3, D32 16.3.2): runs wisp.exe through
    the SLO states, samples each, checks the settle row, and self-tests the
    gate with a forced leak. Emits one aggregate JSON report.

.DESCRIPTION
    Per state: builds/uses build\wisp.exe, runs `wisp slo -state <X>` which
    boots the real runtime skeleton and samples the seven D32 metrics over
    the process tree (private working set gate units per docs/SLO.md §7).
    The driver's JSON verdict (exit 0 pass / 1 fail) is aggregated here.

    Subset "smoke" (windows-latest PR gate, SPEC-11 §6): Sleeping + Warm
    (memory/handle subset, no audio) + settle + leak self-test.
    Subset "full" (self-hosted wisp-slo runner / local pre-merge): all six
    D32 states + settle + leak self-test. Script is identical; only the
    state list differs. No subset may be skipped (D22 mode-6 ban).

    Before ANY state is sampled the script runs a sampling-validity precheck
    (ticket 134 AC#4): if foreign toolchain/wisp processes or a busy machine
    are found it prints one loud "NO CONCLUSION (machine-contended)" line per
    reason and leaves no number behind. What that refusal COLOURS is decided
    by ticket 134 AC#6 (owner 2026-09-23, Q-36): a contended run exits 0 and
    writes a machine-readable no-conclusion record instead of reddening the
    job, because a validity refusal is not a performance regression. The
    precheck itself may only ever REFUSE more often (AC#4's direction, not
    reversible by this or any later ticket); the two D32 thresholds it
    protects (Sleeping CPU <=0.5%, RSS <=25MB) live in `wisp slo` and are not
    evaluated here. "No conclusion" is not "no consequence": scripts/slo-fresh
    ness.sh probe P3 ages on the newest uploaded slo-report artifact, so a
    machine that keeps refusing to sample goes red on its own nail.

    Output JSON schema (slo-report.json):
    {
      "generated_at": RFC3339,
      "subset": "smoke"|"full",
      "machine": <env: COMPUTERNAME>,
      "results": [ { "state": "...", "exit_code": N, "pass": bool,
                     "report": <wisp slo StateReport JSON> } ],
      "settle":  { "exit_code": N, "pass": bool, "report": <SettleReport> },
      "leak_fixture": { "flipped_to_fail": bool, "pass": <leak run pass> },
      "all_pass": bool
    }
    Exit 0 iff all_pass, or iff this run produced no conclusion at all
    (machine-contended, ticket 134 AC#6 - see above: no number, no red).
    Exit 1 whenever a number WAS produced and did not pass; the script also
    FAILS if the forced leak does NOT flip the gate (a leak that passes means
    the sampler is broken).

.USAGE
    scripts/slo-check.ps1 [-Subset smoke|full] [-SecondsPerState 4]
                          [-WispExe build\wisp.exe] [-OutDir build\slo]
#>
[CmdletBinding()]
param(
    [ValidateSet('smoke', 'full')]
    [string]$Subset = 'smoke',
    [double]$SecondsPerState = 4,
    [string]$WispExe = '',
    [string]$OutDir = ''
)

$ErrorActionPreference = 'Stop'
Set-StrictMode -Version 2.0
$ProgressPreference = 'SilentlyContinue'

$RepoRoot = Split-Path -Parent $PSScriptRoot
if (-not $WispExe) { $WispExe = Join-Path $RepoRoot 'build\wisp.exe' }
if (-not $OutDir) { $OutDir = Join-Path $RepoRoot 'build\slo' }

function Fail([string]$Message) {
    Write-Host "slo-check.ps1: FATAL: $Message"
    exit 1
}

function Get-OwnProcessTree {
    # Everything reachable from this script as a parent or a child: the runner
    # host (Runner.Worker / Runner.Listener), the powershell host Actions used
    # to start us, and our own wisp.exe children must never be mistaken for
    # contention. Hops are bounded because a corrupt ParentProcessId would
    # otherwise loop forever, and a truncated tree only ever makes the precheck
    # refuse MORE (the one direction ticket 134 AC#4 allows).
    param([int]$RootPid, [array]$All)
    $ids = @($RootPid)
    for ($hop = 0; $hop -lt 16; $hop++) {
        $grown = $false
        foreach ($p in $All) {
            $procId = [int]$p.ProcessId
            $parentId = 0
            if ($p.ParentProcessId) { $parentId = [int]$p.ParentProcessId }
            if ($parentId -gt 0) {
                if (($ids -contains $procId) -and -not ($ids -contains $parentId)) { $ids += $parentId; $grown = $true }
            }
            if (($ids -contains $parentId) -and -not ($ids -contains $procId)) { $ids += $procId; $grown = $true }
        }
        if (-not $grown) { break }
    }
    return @($ids | Select-Object -Unique)
}

if (-not (Test-Path $WispExe)) {
    # Not a skippable step: the gate builds its own binary when missing.
    Write-Host "slo-check.ps1: wisp.exe missing; building (scripts/build.ps1 -Env dev)"
    & (Join-Path $RepoRoot 'scripts\build.ps1') -Env dev
    if ($LASTEXITCODE -ne 0) { Fail 'build.ps1 failed' }
}
if (-not (Test-Path $WispExe)) { Fail "wisp.exe not found at $WispExe" }

New-Item -ItemType Directory -Force -Path $OutDir | Out-Null
$env:WISP_ENV = 'test'
$env:WISP_TEST_DATA_DIR = Join-Path $OutDir 'data'

# --- clear stale numbers ---------------------------------------------------
# The runner reuses E:\work\base\actions-runner\_work\wisp\wisp across runs, so
# build\slo can still hold the previous run's JSON when this run refuses to
# sample. "No numbers" has to mean no numbers on disk, not just no new ones.
$stale = @(Get-ChildItem -Path $OutDir -Filter '*.json' -ErrorAction SilentlyContinue)
if ($stale.Count -gt 0) {
    Write-Host ("slo-check.ps1: clearing {0} stale report file(s) from {1}" -f $stale.Count, $OutDir)
    $stale | Remove-Item -Force -ErrorAction SilentlyContinue
}

# --- sampling validity precheck (ticket 134 AC#4, shape C) -----------------
# Why this exists (A103 in docs/reports/pending-and-issues.md, measured, not
# theorized): slo-full runs on a self-hosted runner on the SAME 6C12T laptop the
# review fleet compiles on, so one push can start a six-state sample while three
# agents are mid `go test`. The resulting D32 reading can then be falsely red
# (we chase a regression that does not exist) or falsely green-by-noise. Both
# are worse than no reading, and both are reversible only if we say so out loud.
#
# This block therefore never weakens a check. It only ever adds a reason to
# refuse to produce numbers, so a contended machine yields the word
# "machine-contended" instead of a number. It is NOT a threshold: the D32 rows
# (Sleeping CPU <=0.5%, RSS <=25MB) are evaluated by `wisp slo` itself and are
# untouched by this ticket. There is deliberately no -Skip / -Force / env
# override - a switch is a skip, and D22 mode-6 bans that.
Write-Host 'slo-check.ps1: sampling validity precheck (ticket 134 AC#4)'
try {
    $allProcs = @(Get-CimInstance -ClassName Win32_Process -ErrorAction Stop)
} catch {
    Fail "sampling validity precheck cannot enumerate processes: $_"
}
if ($allProcs.Count -lt 2) {
    Fail ('sampling validity precheck enumerated {0} process(es); refusing to sample on an unreadable process table (ticket 134 AC#4)' -f $allProcs.Count)
}
$ownIds = @(Get-OwnProcessTree -RootPid $PID -All $allProcs)

# Foreign toolchain activity. Names only - no path guessing, so a leftover
# wisp.exe from an earlier crash counts as contention too (it holds CPU and RSS
# while we sample). The gate's own tree is excluded via $ownIds, and at this
# point the gate has started no wisp.exe yet.
$loadNames = @('go.exe', 'gofmt.exe', 'cgo.exe', 'compile.exe', 'asm.exe', 'link.exe',
               'gcc.exe', 'g++.exe', 'cc1.exe', 'cc1plus.exe', 'as.exe', 'ld.exe',
               'wisp.exe', 'wisp-cli.exe', 'staticcheck.exe')

# "Anything under the runner's work root" is somebody else's CI checkout at
# work. Derived from the variables Actions sets, so a local run (where they are
# unset) simply loses this extra reason and keeps the name probe.
$workPrefixes = @()
if ($env:GITHUB_WORKSPACE) {
    $ws = [IO.Path]::GetFullPath($env:GITHUB_WORKSPACE)
    for ($depth = 0; $depth -lt 2; $depth++) {
        $ws = Split-Path -Parent $ws
        if ($ws) { $workPrefixes += $ws }
    }
}
if ($env:RUNNER_TEMP) { $workPrefixes += [IO.Path]::GetFullPath($env:RUNNER_TEMP) }

$offenders = @()
foreach ($p in $allProcs) {
    $procId = [int]$p.ProcessId
    if ($ownIds -contains $procId) { continue }
    $name = ''
    if ($p.Name) { $name = $p.Name.ToLowerInvariant() }
    $exePath = ''
    if ($p.ExecutablePath) { $exePath = [string]$p.ExecutablePath }
    $why = $null
    if ($loadNames -contains $name) { $why = 'foreign toolchain/wisp process present' }
    if (-not $why -and $exePath) {
        foreach ($prefix in $workPrefixes) {
            if ($exePath.StartsWith($prefix, [StringComparison]::OrdinalIgnoreCase)) {
                $why = 'process running under the runner work root'
            }
        }
    }
    if ($why) {
        $started = 'unknown'
        if ($p.CreationDate) { $started = ([datetime]$p.CreationDate).ToString('yyyy-MM-dd HH:mm:ss') }
        $offenders += [pscustomobject]@{
            pid    = $procId
            name   = $name
            path   = $exePath
            reason = $why
            started = $started
        }
    }
}

# Machine-wide load, so contention from a process this list does not name
# (another agent's already-built binary, an indexer) is also refused. Max of two
# 1s samples: a single sample can be a blip, and a 2s window costs 2s of a run
# that already spends ~40s sampling.
#
# Honest limitation, recorded rather than hidden: if the WMI perf class is not
# available the probe reports 'unavailable' and does NOT fail. Failing closed on
# a missing perf counter would turn both subsets permanently red, which is the
# "gate that never produces a verdict" disease this ticket exists to prevent.
# The identity probe above still gates on such a machine, so the direction stays
# "refuse at least as often".
$cpuBusyPct = 50
$cpuMax = -1
for ($sample = 0; $sample -lt 2; $sample++) {
    try {
        $raw = Get-CimInstance -ClassName Win32_PerfFormattedData_PerfOS_Processor -Filter "Name='_Total'" -ErrorAction Stop | Select-Object -ExpandProperty PercentProcessorTime
        if ($null -ne $raw) { $value = [int]$raw; if ($value -gt $cpuMax) { $cpuMax = $value } }
    } catch { $cpuMax = -1; break }
    if ($sample -lt 1) { Start-Sleep -Seconds 1 }
}

$reasons = @()
foreach ($o in $offenders) {
    $reasons += ('{0}: {1} pid={2} started={3} path={4}' -f $o.reason, $o.name, $o.pid, $o.started, $o.path)
}
if ($cpuMax -ge $cpuBusyPct) {
    $reasons += ('machine-wide cpu utilisation {0}% over a 1s window (>= {1}%)' -f $cpuMax, $cpuBusyPct)
}

if ($reasons.Count -gt 0) {
    # ---- colouring rule (ticket 134 AC#6, owner 2026-09-23 Q-36 "都按推荐") ----
    # This branch used to exit 1. On a self-hosted runner that shares one
    # laptop with a compiling review fleet (A103) that meant the D32 merge gate
    # was red for "the machine was busy", and A115 measured what that costs:
    # with the badge already red for two days straight, a real breakage
    # (ticket 131's Linux compile error) sat unnoticed for 4h17m. Owner's
    # ruling: a refused sample is a NO CONCLUSION, not a fail.
    #
    # What is NOT changed here, deliberately and permanently:
    #   * WHEN we refuse - the three reasons above are AC#4's and may only ever
    #     grow (this ticket does not touch a single condition).
    #   * the numbers themselves - D32's Sleeping CPU <=0.5% and private RSS
    #     <=25MB are evaluated by `wisp slo` and are one-byte identical.
    #   * the "cannot look" paths - Fail() above still exits 1, because "the
    #     instrument is broken" is not the same sentence as "no verdict today".
    #   * any number produced by a run that DID sample: `all_pass=False` still
    #     exits 1 at the bottom of this script.
    # The step in .github/workflows/ci.yml stays unconditional (D22 mode-6:
    # no `if:`, no continue-on-error), so the colour is carried by this exit
    # code and by the record written below - nothing else.
    Write-Host ('slo-check.ps1: NO CONCLUSION (machine-contended) - subset={0} refused to sample, no numbers were produced' -f $Subset)
    foreach ($reason in $reasons) {
        Write-Host ("slo-check.ps1: machine-contended reason: {0}" -f $reason)
    }
    Write-Host ('slo-check.ps1: NO CONCLUSION (machine-contended): {0} reason(s), 0 state file(s) written, slo-report.json NOT written' -f $reasons.Count)
    Write-Host 'slo-check.ps1: NO CONCLUSION (machine-contended): this is a SAMPLING VALIDITY verdict, not a'
    Write-Host 'slo-check.ps1: NO CONCLUSION (machine-contended): performance result - D32 (CPU <=0.5%, private'
    Write-Host 'slo-check.ps1: NO CONCLUSION (machine-contended): RSS <=25MB) stays unverified for this run.'
    Write-Host 'slo-check.ps1: NO CONCLUSION (machine-contended): Rerun when the machine is quiet (ticket 134'
    Write-Host 'slo-check.ps1: NO CONCLUSION (machine-contended): AC#4, shape C: refuse loudly instead of'
    Write-Host 'slo-check.ps1: NO CONCLUSION (machine-contended): reporting a contaminated sample as truth).'
    Write-Host 'slo-check.ps1: NO CONCLUSION (machine-contended): AC#6 recolors the refusal, it does not lift it.'
    Write-Host 'slo-check.ps1: NO CONCLUSION (machine-contended): what keeps "no conclusion" from rotting in'
    Write-Host 'slo-check.ps1: NO CONCLUSION (machine-contended): silence is probe P3 of scripts/slo-freshness.sh,'
    Write-Host 'slo-check.ps1: NO CONCLUSION (machine-contended): which ages on the newest uploaded SLO REPORT'
    Write-Host 'slo-check.ps1: NO CONCLUSION (machine-contended): artifact - never on this job, never on this'
    Write-Host 'slo-check.ps1: NO CONCLUSION (machine-contended): exit code. Enough busy days in a row and the'
    Write-Host 'slo-check.ps1: NO CONCLUSION (machine-contended): nail goes red by itself (ticket 134 AC#6).'

    # The record half of the colouring rule. This file is NOT the report: it
    # carries no state rows, no numbers and no all_pass flag, and ci.yml's
    # Upload step names build/slo/slo-report.json only, so nothing here can be
    # mistaken for a sample by the nail or by a human. It exists so "we refused,
    # and here is why" survives on disk after the log window closes.
    $noConclusion = [pscustomobject]@{
        generated_at         = (Get-Date).ToUniversalTime().ToString('yyyy-MM-ddTHH:mm:ssZ')
        subset               = $Subset
        verdict              = 'no-conclusion'
        reason               = 'machine-contended'
        machine              = $env:COMPUTERNAME
        reason_count         = $reasons.Count
        reasons              = $reasons
        offenders            = @($offenders | ForEach-Object {
            [pscustomobject]@{ name = $_.name; pid = $_.pid; path = $_.path; why = $_.reason; started = $_.started }
        })
        state_files_written  = 0
        slo_report_written   = $false
        d32_evaluated        = $false
        meaning              = 'No number was produced, so no number may be read out of this run. Rerun on a quiet machine.'
        nail                 = 'scripts/slo-freshness.sh P3 ages on the newest slo-full-report artifact, not on this record'
    }
    $noConclusionPath = Join-Path $OutDir 'slo-no-conclusion.json'
    $noConclusion | ConvertTo-Json -Depth 5 | Set-Content -Path $noConclusionPath -Encoding UTF8
    Write-Host ('slo-check.ps1: NO CONCLUSION (machine-contended): record written to {0} (NOT a report; slo-report.json is not written on this path)' -f $noConclusionPath)

    if ($env:GITHUB_STEP_SUMMARY) {
        @(
            '## slo-check: NO CONCLUSION (machine-contended)',
            '',
            ('Refused to sample subset `{0}` ({1} reason(s)); no SLO numbers produced. D32 is not evaluated and not disputed by this run.' -f $Subset, $reasons.Count),
            '',
            ('| reason | process | pid | path |'),
            ('|---|---|---|---|')
        ) + @($offenders | ForEach-Object {
            ('| {0} | {1} | {2} | {3} |' -f $_.reason, $_.name, $_.pid, $_.path)
        }) + @(
            '',
            ('Ticket 134 AC#6: this is reported as "no conclusion", not as a failure. The freshness nail (scripts/slo-freshness.sh P3) goes red if no uploaded SLO report is younger than its window.')
        ) | Add-Content -Path $env:GITHUB_STEP_SUMMARY -Encoding UTF8
    }
    Write-Host 'slo-check.ps1: NO CONCLUSION (machine-contended): exit 0'
    exit 0
}
$cpuText = if ($cpuMax -ge 0) { "$cpuMax%" } else { 'unavailable (perf class unreadable)' }
Write-Host ("slo-check.ps1: precheck ok - no foreign toolchain/runner process, machine-wide cpu max {0}" -f $cpuText)

$states = @('Sleeping', 'Warm')
if ($Subset -eq 'full') {
    $states = @('Sleeping', 'Armed', 'Warm', 'Conversation', 'PanelOpen', 'WorkPeak')
}

$results = @()
foreach ($state in $states) {
    $outFile = Join-Path $OutDir ("state-{0}.json" -f $state)
    Write-Host ("slo-check.ps1: sampling state {0} for {1}s" -f $state, $SecondsPerState)
    & $WispExe slo -state $state -seconds $SecondsPerState -interval-ms 250 -out $outFile
    $code = $LASTEXITCODE
    $pass = $false
    $reportJson = $null
    if ($code -eq 0 -and (Test-Path $outFile)) {
        try {
            $reportJson = Get-Content $outFile -Raw | ConvertFrom-Json
            $pass = [bool]$reportJson.pass
        } catch {
            Write-Host "slo-check.ps1: WARNING: state $state produced unparseable JSON: $_"
        }
    }
    Write-Host ("slo-check.ps1: state {0} exit={1} pass={2}" -f $state, $code, $pass)
    $results += [pscustomobject]@{ state = $state; exit_code = $code; pass = $pass; report = $reportJson }
}

# --- settle row (D32: 10s back under cap + FreeOSMemory counter > 0) ------
Write-Host "slo-check.ps1: settle check (10s window)"
$settleFile = Join-Path $OutDir 'settle.json'
& $WispExe slo -settle -seconds 10 -out $settleFile
$settleCode = $LASTEXITCODE
$settlePass = $false
$settleReport = $null
if ($settleCode -eq 0 -and (Test-Path $settleFile)) {
    try {
        $settleReport = Get-Content $settleFile -Raw | ConvertFrom-Json
        $settlePass = [bool]$settleReport.pass
        # Independent double-check: the counter must be > 0 in the report.
        if ($settleReport.settle.free_os_memory_count -le 0) {
            Write-Host 'slo-check.ps1: settle report has FreeOSMemoryCount == 0; forcing fail'
            $settlePass = $false
        }
    } catch {
        Write-Host "slo-check.ps1: WARNING: settle produced unparseable JSON: $_"
    }
}
Write-Host ("slo-check.ps1: settle exit={0} pass={1}" -f $settleCode, $settlePass)

# --- leak fixture self-test (must flip the gate RED) ----------------------
Write-Host 'slo-check.ps1: leak fixture self-test (100MB, must FAIL)'
& $WispExe slo -state Sleeping -leak -seconds 2 -interval-ms 250 -out (Join-Path $OutDir 'leak.json')
$leakCode = $LASTEXITCODE
$leakFlipped = ($leakCode -eq 1)
Write-Host ("slo-check.ps1: leak exit={0} flipped_to_fail={1}" -f $leakCode, $leakFlipped)
if (-not $leakFlipped) {
    Fail 'forced 100MB leak did NOT flip the gate red - the sampler is broken (D22 mode-6: gate must detect)'
}

# Force the filter into an array with @(): under Set-StrictMode 2.0 (line 48)
# a zero-match Where-Object returns $null and a one-match returns a scalar,
# neither of which carries .Count - so the old expression threw
# PropertyNotFoundException on exactly the ALL-PASS path, before line 141 could
# write slo-report.json (ticket 66 AC#4). It stayed hidden while A14 made every
# state fail closed with exit=2, because then the filter matched >=2 items.
# The gate must be able to report the verdict it actually computed.
$failingStates = @($results | Where-Object { -not $_.pass })
$allPass = ($failingStates.Count -eq 0) -and $settlePass

$report = [pscustomobject]@{
    generated_at = (Get-Date).ToUniversalTime().ToString('yyyy-MM-ddTHH:mm:ssZ')
    subset       = $Subset
    machine      = $env:COMPUTERNAME
    results      = $results
    settle       = [pscustomobject]@{ exit_code = $settleCode; pass = $settlePass; report = $settleReport }
    leak_fixture  = [pscustomobject]@{ flipped_to_fail = $leakFlipped }
    all_pass     = $allPass
}
$reportPath = Join-Path $OutDir 'slo-report.json'
$report | ConvertTo-Json -Depth 8 | Set-Content -Path $reportPath -Encoding UTF8
Write-Host ("slo-check.ps1: report written to {0} (all_pass={1})" -f $reportPath, $allPass)

if (-not $allPass) { exit 1 }
exit 0
