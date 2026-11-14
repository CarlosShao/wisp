# Ticket 443 repro harness (Windows, no POSIX sh on this bench).
#
# Runs ONE targeted test per process so every iteration pays the cold start the
# CI runner pays (the 90MB test binary is the fixture CLI itself; ticket 52's
# fixture re-executes it as the wrapped program).
#
# Usage:
#   pwsh -File repro.ps1 -N 10 -Name baseline
#   pwsh -File repro.ps1 -N 10 -Name delay700 -DelayMs 700
#
# -DelayMs sets WISP52_443_FIXTURE_DELAY_MS, the fixture's own startup-delay dial
# (see command_52_main_test.go). It is an amplification dial, not an assertion:
# the committed default (unset) is 0.
param(
    [int]$N = 10,
    [string]$Name = "run",
    [int]$DelayMs = 0,
    [string]$Run = "^Test52TimeoutKillsTheProcessTree$"
)

$root = (Resolve-Path "$PSScriptRoot\..\..\..\..\..").Path
$logDir = Join-Path $PSScriptRoot "..\logs"
New-Item -ItemType Directory -Force -Path $logDir | Out-Null

if ($DelayMs -gt 0) { $env:WISP52_443_FIXTURE_DELAY_MS = "$DelayMs" }
else { Remove-Item Env:WISP52_443_FIXTURE_DELAY_MS -ErrorAction SilentlyContinue }

Push-Location $root
$hits = 0
$abstains = 0
$greens = 0
$other = 0
$sw = [Diagnostics.Stopwatch]::StartNew()
for ($i = 1; $i -le $N; $i++) {
    $log = Join-Path $logDir ("{0}-{1:d2}.log" -f $Name, $i)
    $out = & go test -count=1 -run $Run ./internal/plugin/manifestslot 2>&1 | Out-String
    Set-Content -Path $log -Value $out -NoNewline
    if ($out -match 'never reported a grandchild pid') { $hits++ }
    elseif ($out -match 'ABSTAIN\(443:') { $abstains++ }
    elseif ($out -match '(?m)^ok\s') { $greens++ }
    else { $other++ }
    Write-Output ("run {0}/{1}: {2}" -f $i, $N, ($(if ($out -match 'never reported a grandchild pid') { 'RED(never-reported)' } elseif ($out -match 'ABSTAIN\(443:') { 'ABSTAIN' } elseif ($out -match '(?m)^ok\s') { 'GREEN' } else { 'OTHER' })))
}
$sw.Stop()
Pop-Location
Write-Output ("SUMMARY name={0} N={1} delay_ms={2} red-never-reported={3} abstain={4} green={5} other={6} elapsed_ms={7}" -f $Name, $N, $DelayMs, $hits, $abstains, $greens, $other, $sw.ElapsedMilliseconds)
