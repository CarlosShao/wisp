# fetch-red-window.ps1 -- ticket 447: pull the FAILING RUN's own words.
#
# For one run id and one test name, print every step-10 line that
#   (a) sits between "=== RUN   <name>" and "--- FAIL: <name>", or
#   (b) is that test's "=== RUN"/"--- FAIL" line itself,
# with the log's own stamp stripped, so the evidence file can quote the runner
# verbatim instead of paraphrasing it. Usage:
#   pwsh -File fetch-red-window.ps1 -Run 43583105791 -Test TestPanelHostRealWindowHopAndLifecycle
param(
    # [long], not [int]: a GitHub run id is past Int32 (43575529807 overflows)
    [Parameter(Mandatory = $true)][long]$Run,
    [Parameter(Mandatory = $true)][string]$Test
)
$ErrorActionPreference = 'Stop'
$gh = 'C:\Program Files\GitHub CLI\gh.exe'
$root = 'D:\work\workspace\projects plans\Wisp'
Set-Location $root
$jobs = (& $gh run view $Run --json jobs) | ConvertFrom-Json
$job = $jobs.jobs | Where-Object { $_.name -eq 'test-windows' } | Select-Object -First 1
$step = $job.steps | Where-Object { $_.number -eq 10 } | Select-Object -First 1
$lines = & $gh run view $Run --job $job.databaseId --log 2>&1
$on = $false
foreach ($l in $lines) {
    if ($l -isnot [string]) { continue }
    $c = ($l -split "`t", 3)[2]
    if (-not $c) { continue }
    $c = $c -replace '^\uFEFF', ''
    if ($c -match '^\d{4}-\d\d-\d\dT[\d:.]+Z\s') { $c = $c.Substring($Matches[0].Length) }
    if ($c -match "^=== RUN\s+$([regex]::Escape($Test))\b") { $on = $true }
    if ($c -match "^--- (FAIL|PASS|SKIP): $([regex]::Escape($Test))\b") {
        $c
        break
    }
    if ($on) { $c }
}
