# 450 AC#3 -- the ISOLATED arm: the cases the two 400-second deaths stopped inside,
# run WITHOUT the whole-package load, one process per arm.
#
# WHY THIS EXISTS: AC#3 asks whether Test145ArgsJSON* / TestAC4PriorFocus* die on
# their own too, or only inside the whole-package run. Each arm below is one
# carrier invocation (scripts/wisp-cli-staged-tests.sh, the only sanctioned way to
# run cmd/wisp on this host) with a full-name-anchored -test.run, so the reading
# is "this case, alone, in a fresh process".
#
# WHAT IT DOES NOT DO: it does not touch the resident ball (no kill/stop/start),
# does not run two arms at once, and does not weaken any flag - the arms differ
# from the whole-package run only in -test.run.
#
# THE THIRD ARM is the death REGION together (the two cases above plus the
# environment-gated focus ruler that is one of this machine's two known reds),
# which is the "panel region alone, no resident/voice load" reading.
#
# Readings land in ac3-readings.txt next to this file; per-arm logs keep their
# own stage dir.

$ErrorActionPreference = 'Stop'
$probe = $PSScriptRoot
$repo  = Split-Path -Parent (Split-Path -Parent (Split-Path -Parent (Split-Path -Parent $probe)))
if (-not (Test-Path (Join-Path $repo 'go.mod'))) { $repo = 'D:\work\workspace\projects plans\Wisp' }
$bash  = 'D:\work\soft\Git\bin\bash.exe'
if (-not (Test-Path $bash)) { Write-Host "ac3-isolated.ps1: no Git bash at $bash; nothing was run."; exit 2 }

$arms = @(
    [pscustomobject]@{ name = 'A1 Test145ArgsJSON alone';                run = '^Test145ArgsJSONComesFromTheBookedForensicsRow$';                                   stage = 'stage-ac3-a1' },
    [pscustomobject]@{ name = 'A2 TestAC4PriorFocus alone';              run = '^TestAC4PriorFocusSurvivesARefusedPanelSample$';                                    stage = 'stage-ac3-a2' },
    [pscustomobject]@{ name = 'A3 the death region together (3 cases)';  run = '^(Test145ArgsJSONComesFromTheBookedForensicsRow|TestAC4PriorFocusSurvivesARefusedPanelSample|TestAC4FocusReturnToPriorWindowGap33r5)$'; stage = 'stage-ac3-a3' }
)

$readings = Join-Path $probe 'ac3-readings.txt'
"# 450 AC#3 isolated arms -- $(Get-Date -Format o)" | Out-File -Encoding utf8 $readings

foreach ($arm in $arms) {
    $otherTest = @(Get-Process -Name 'wisp.test' -ErrorAction SilentlyContinue)
    if ($otherTest.Count -gt 0) {
        "ARM $($arm.name): REFUSED - a wisp.test.exe is already running (PID $(($otherTest | ForEach-Object Id) -join ',')); one run at a time. Nothing was run." |
            Out-File -Append -Encoding utf8 $readings
        continue
    }
    $env:WISP_CLI_STAGE_DIR = Join-Path $probe $arm.stage
    $out = Join-Path $probe ("ac3-" + ($arm.stage -replace 'stage-ac3-', '') + ".out.txt")
    $t0 = Get-Date
    Push-Location $repo
    try { & $bash 'scripts/wisp-cli-staged-tests.sh' "-test.run=$($arm.run)" *> $out; $rc = $LASTEXITCODE }
    finally { Pop-Location }
    $t1 = Get-Date
    $log = Join-Path $env:WISP_CLI_STAGE_DIR 'wisp-cli-tests.log'
    $lines = @(Get-Content $log)
    $ran     = @($lines | Where-Object { $_ -match '^=== RUN' }).Count
    $passed  = @($lines | Where-Object { $_ -match '^--- PASS' }).Count
    $failed  = @($lines | Where-Object { $_ -match '^--- FAIL' }).Count
    $skipped = @($lines | Where-Object { $_ -match '^--- SKIP' }).Count
    $last    = if ($lines.Count -gt 0) { $lines[-1] } else { '<empty>' }
    $verdict = @($lines | Where-Object { $_ -match '^(ok|FAIL)\s+github\.com/CarlosShao/wisp/cmd/wisp' } | Select-Object -Last 1)
    $perCase = @($lines | Where-Object { $_ -match '^--- (PASS|FAIL|SKIP): ' }) -join ' ; '
    $verdictText = if ($verdict) { $verdict[0] } else { '<none>' }
    @(
        "ARM $($arm.name)",
        "  pattern    = $($arm.run)",
        "  start/end  = $($t0.ToString('HH:mm:ss')) / $($t1.ToString('HH:mm:ss'))  ($([int]($t1-$t0).TotalSeconds) s incl. build)",
        "  exit_code  = $rc",
        "  RUN=$ran PASS=$passed FAIL=$failed SKIP=$skipped",
        "  verdict    = $verdictText",
        "  cases      = $perCase",
        "  last_line  = $last",
        "  log        = $log"
    ) | Out-File -Append -Encoding utf8 $readings
    Write-Host "ARM $($arm.name): exit=$rc RUN=$ran PASS=$passed FAIL=$failed SKIP=$skipped"
}
Write-Host "ac3-isolated.ps1: readings -> $readings"
