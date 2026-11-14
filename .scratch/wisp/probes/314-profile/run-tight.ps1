# 314 profile A/B tight-loop: one test process at a time, but the NEXT one starts
# immediately after the previous exits (no inspection gap), so the predecessor's
# msedgewebview2 children are still tearing down while the next boot names the UDF.
# usage: pwsh -File run-tight.ps1 -Bin <bin> -Tag A1t -Runs 3
param(
  [Parameter(Mandatory=$true)][string]$Bin,
  [Parameter(Mandatory=$true)][string]$Tag,
  [Parameter(Mandatory=$true)][int]$Runs
)
$ErrorActionPreference = 'Stop'
$root = (Resolve-Path (Join-Path $PSScriptRoot '..\..\..\..')).Path
$logDir = Join-Path $PSScriptRoot 'logs'
$re = '^(TestAC14GoSideEvalPushReachesThePage|TestAC13BringUpSurvivesAReusedThreadQuit|TestPanelThreadIsSTAAndExitsCleanly|TestAC13BringUpRefusesAThreadWithAQueuedClose)$'
$rows = @()
for ($i = 1; $i -le $Runs; $i++) {
  $out = Join-Path $logDir "$Tag-r$i.out.txt"
  $err = Join-Path $logDir "$Tag-r$i.err.txt"
  if (Test-Path $out) { Remove-Item $out -Force }
  if (Test-Path $err) { Remove-Item $err -Force }
  $competitors = @(Get-Process wisp.test -ErrorAction SilentlyContinue).Count
  $bwBefore = @(Get-Process msedgewebview2 -ErrorAction SilentlyContinue).Count
  $sw = [Diagnostics.Stopwatch]::StartNew()
  $p = Start-Process -FilePath $Bin -ArgumentList @('-test.run', $re, '-test.v', '-test.timeout', '10m') `
    -WorkingDirectory $root -RedirectStandardOutput $out -RedirectStandardError $err -PassThru -NoNewWindow
  $hung = -not $p.WaitForExit(300000)
  if ($hung) { Stop-Process -Id $p.Id -Force -ErrorAction SilentlyContinue }
  $sw.Stop()
  $rc = if ($hung) { 'KILLED-300s' } else { $p.ExitCode }
  $bwAfter = @(Get-Process msedgewebview2 -ErrorAction SilentlyContinue).Count
  $t = if (Test-Path $out) { Get-Content $out -Raw } else { '' }
  $e = if (Test-Path $err) { Get-Content $err -Raw } else { '' }
  $both = $t + $e
  $rows += [ordered]@{
    tag = $Tag; run = $i; rc = $rc
    rcHex = if ($rc -is [int]) { '0x{0:X8}' -f ([uint32]($rc -band 0xFFFFFFFF)) } else { $null }
    elapsed_s = [math]::Round($sw.Elapsed.TotalSeconds, 2)
    passLines = ([regex]::Matches($both, '(?m)^--- PASS:')).Count
    failLines = ([regex]::Matches($both, '(?m)^--- FAIL:')).Count
    skipLines = ([regex]::Matches($both, '(?m)^--- SKIP:')).Count
    panic = ([regex]::Matches($both, '(?m)^panic:')).Count
    bootFailed = ([regex]::Matches($both, 'webview2: .*boot failed')).Count
    embedFailed = ([regex]::Matches($both, 'WebView2 Embed failed')).Count
    createWindow = ([regex]::Matches($both, 'MARK299 after CreateWindowExW')).Count
    accessViolation = ([regex]::Matches($both, '0xc0000005|access violation')).Count
    competitorsAtStart = $competitors
    browsersBefore = $bwBefore; browsersRightAfterExit = $bwAfter
  }
}
foreach ($r in $rows) {
  $line = ($r | ConvertTo-Json -Compress)
  Add-Content -Path (Join-Path $PSScriptRoot 'results-tight.jsonl') -Value $line -Encoding utf8
  $line
}
