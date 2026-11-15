# 450 AC#2(b) -- the BALL-STOPPED arm of the whole-package bisect, as ONE command.
#
# WHY THIS FILE EXISTS: ticket 450's AC#2(b) is "the same whole-package command
# with the resident ball NOT running". That arm needs the ball to be gone, and
# the owner's own rule forbids this leg from killing or restarting it. So the
# arm is handed over as this script: the owner exits the ball from the tray, then
# (or the orchestrator, once the owner says it is down) runs
#
#     powershell -NoProfile -File ".scratch\wisp\probes\450\run-balloff.ps1"
#
# WHAT IT NEVER DOES: it does not stop, kill, suspend, restart or inspect-into
# the ball. It only LOOKS for a live wisp.exe at the repo root and REFUSES
# (exit 3) while one is alive. A refusal is the honest answer, not a failure of
# the measurement.
#
# WHAT IT PRINTS (the same readings AC#1/AC#2(a) record, so the three arms are
# comparable): start/end wall clock, the carrier's exit code, the four numbers
# read off the -v stream (=== RUN / --- PASS / --- FAIL / --- SKIP), the last
# line of output, and whether the process ended WITHOUT a verdict (no `--- FAIL`
# line while the four numbers fall short of a complete run = the "died mid-test"
# shape ticket 450 is about, as opposed to an assertion red).
#
# The run itself still goes through the carrier scripts/wisp-cli-staged-tests.sh
# (the only sanctioned way to run the whole package on this host), one run at a
# time, with its own stage dir so it cannot collide with a leg's staging area.

$ErrorActionPreference = 'Stop'
$repo = Split-Path -Parent (Split-Path -Parent (Split-Path -Parent (Split-Path -Parent (Split-Path -Parent $PSScriptRoot))))
# $PSScriptRoot = <repo>\.scratch\wisp\probes\450
if (-not (Test-Path (Join-Path $repo 'go.mod'))) { $repo = 'D:\work\workspace\projects plans\Wisp' }

$probe = Join-Path $repo '.scratch\wisp\probes\450'
$stamp = Get-Date -Format 'yyyyMMdd-HHmmss'
$out   = Join-Path $probe "run-balloff-$stamp.out.txt"
$env:WISP_CLI_STAGE_DIR = Join-Path $probe 'stage-balloff'

# --- precondition 1: no live ball -------------------------------------------------
$liveBall = @(Get-CimInstance Win32_Process -Filter "Name='wisp.exe'" -ErrorAction SilentlyContinue |
    Where-Object { $_.ExecutablePath -and $_.ExecutablePath.Equals((Join-Path $repo 'wisp.exe'), [StringComparison]::OrdinalIgnoreCase) })
if ($liveBall.Count -gt 0) {
    $ids = ($liveBall | ForEach-Object { "$($_.ProcessId) (started $($_.CreationDate))" }) -join ', '
    Write-Host "run-balloff.ps1: REFUSED - the resident ball is still running: PID $ids"
    Write-Host "run-balloff.ps1: AC#2(b) is exactly the arm where it is NOT running. Exit it from the tray"
    Write-Host "run-balloff.ps1: (or ask the owner to), re-check, then run this file again."
    Write-Host "run-balloff.ps1: this script never stops it; nothing was run."
    exit 3
}

# --- precondition 2: no other whole-package run in flight --------------------------
$otherTest = @(Get-Process -Name 'wisp.test' -ErrorAction SilentlyContinue)
if ($otherTest.Count -gt 0) {
    Write-Host "run-balloff.ps1: REFUSED - a wisp.test.exe is already running (PID $(($otherTest | ForEach-Object Id) -join ',')):"
    Write-Host "run-balloff.ps1: 'go test ./cmd/wisp' is one-run-at-a-time on this host, and a second run"
    Write-Host "run-balloff.ps1: invalidates both readings. Nothing was run."
    exit 3
}

# --- the carrier ------------------------------------------------------------------
$bash = 'D:\work\soft\Git\bin\bash.exe'
if (-not (Test-Path $bash)) {
    $cand = @(
        (Get-Command bash.exe -ErrorAction SilentlyContinue | Select-Object -First 1 -ExpandProperty Source),
        'C:\Program Files\Git\bin\bash.exe'
    ) | Where-Object { $_ -and (Test-Path $_) }
    if (-not $cand) { Write-Host "run-balloff.ps1: no msys/cygwin bash found; the WSL relay would run zero tests. Nothing was run."; exit 2 }
    $bash = $cand[0]
}
if ($bash -like '*System32*') {
    Write-Host "run-balloff.ps1: $bash is the WSL relay and runs zero tests. Nothing was run."
    exit 2
}

Write-Host "run-balloff.ps1: ball NOT running (precondition met). Carrier: scripts/wisp-cli-staged-tests.sh"
Write-Host "run-balloff.ps1: output -> $out"
$t0 = Get-Date
Push-Location $repo
try {
    & $bash 'scripts/wisp-cli-staged-tests.sh' *> $out
    $rc = $LASTEXITCODE
} finally {
    Pop-Location
}
$t1 = Get-Date

$log = Join-Path $env:WISP_CLI_STAGE_DIR 'wisp-cli-tests.log'
if (-not (Test-Path $log)) { $log = $out }
$lines   = @(Get-Content $log)
$ran     = @($lines | Where-Object { $_ -match '^=== RUN' }).Count
$passed  = @($lines | Where-Object { $_ -match '^--- PASS' }).Count
$failed  = @($lines | Where-Object { $_ -match '^--- FAIL' }).Count
$skipped = @($lines | Where-Object { $_ -match '^--- SKIP' }).Count
$last    = if ($lines.Count -gt 0) { $lines[-1] } else { '<empty>' }
$verdictLine = @($lines | Where-Object { $_ -match '^(ok|FAIL)\s+github\.com/CarlosShao/wisp/cmd/wisp' } | Select-Object -Last 1)

# The death shape, named rather than described: no FAIL line, yet the run did not
# reach a package verdict of its own making (the carrier prints the verdict line
# itself, so a missing one means the carrier never got to it, and a FAIL verdict
# with zero --- FAIL lines is the process-disappeared shape).
$deathShape = if ($failed -eq 0 -and $verdictLine -and $verdictLine[0] -match '^FAIL') { 'YES (no --- FAIL line, verdict FAIL = process disappeared mid-test)' } else { 'no' }

$report = [ordered]@{
    'arm'          = 'AC#2(b) ball NOT running'
    'start'        = $t0.ToString('o')
    'end'          = $t1.ToString('o')
    'duration_s'   = [int]($t1 - $t0).TotalSeconds
    'exit_code'    = $rc
    'RUN'          = $ran
    'PASS'         = $passed
    'FAIL'         = $failed
    'SKIP'         = $skipped
    'verdict_line' = if ($verdictLine) { $verdictLine[0] } else { '<none>' }
    'last_line'    = $last
    'died_midway'  = $deathShape
    'log'          = $log
}
Write-Host ''
Write-Host '=== 450 AC#2(b) reading ==='
foreach ($k in $report.Keys) { Write-Host ("{0,-13} {1}" -f $k, $report[$k]) }
$report.GetEnumerator() | ForEach-Object { "$($_.Key)=$($_.Value)" } |
    Out-File -Append -Encoding utf8 (Join-Path $probe 'balloff-readings.txt')
Write-Host ''
Write-Host "run-balloff.ps1: readings appended to $(Join-Path $probe 'balloff-readings.txt')"
exit $rc
