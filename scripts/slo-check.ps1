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
    Exit 0 iff all_pass; the script FAILS if the forced leak does NOT flip
    the gate (a leak that passes means the sampler is broken).

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

$allPass = (($results | Where-Object { -not $_.pass }).Count -eq 0) -and $settlePass

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
