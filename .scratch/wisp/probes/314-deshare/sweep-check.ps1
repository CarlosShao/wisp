# 314 de-share sweep control: prove the startup sweep deletes ONLY old
# prefix-<pid>-<nanotime> siblings of the prefix a run is about to use.
#
# Four decoys are planted in $env:TEMP before the run and checked after it:
#   old-leaf      wisp-299r1-panel-profile-<n>-<n>   mtime -2h   EXPECT: gone
#   fresh-leaf    wisp-299r1-panel-profile-<n>-<n>   mtime now   EXPECT: kept
#   other-prefix  wisp-33r1-panel-profile-<n>-<n>    mtime -2h   EXPECT: kept
#   fixed-legacy  wisp-299r1-panel-profile           mtime -2h   EXPECT: kept
#   lookalike     wisp-299r1-panel-profile-abc-123   mtime -2h   EXPECT: kept
#
# The fresh leaf is the stand-in for "a concurrent run's folder": the whole point
# of the age gate is that a live run's leaf is seconds old and must survive.
#
# usage: pwsh -File sweep-check.ps1 -Bin <test.exe>
param(
  [Parameter(Mandatory=$true)][string]$Bin
)
$ErrorActionPreference = 'Stop'
$root = (Resolve-Path (Join-Path $PSScriptRoot '..\..\..\..')).Path
$logDir = Join-Path $PSScriptRoot 'logs'
New-Item -ItemType Directory -Force -Path $logDir | Out-Null
$tmp = $env:TEMP
$stamp = [DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds()
function Plant([string]$name, [int]$ageHours) {
  $p = Join-Path $tmp $name
  New-Item -ItemType Directory -Force -Path $p | Out-Null
  Set-Content -Path (Join-Path $p 'planted-by-314-deshare-sweep-check.txt') -Value $name -Encoding utf8
  (Get-Item $p).LastWriteTime = (Get-Date).AddHours(-$ageHours)
  return $p
}
$oldLeaf = Plant "wisp-299r1-panel-profile-4242-$stamp" 2
$freshLeaf = Plant "wisp-299r1-panel-profile-4343-$stamp" 0
$otherPrefix = Plant "wisp-33r1-panel-profile-4444-$stamp" 2
$fixedLegacy = Plant 'wisp-299r1-panel-profile' 2
$lookalike = Plant "wisp-299r1-panel-profile-abc-$stamp" 2
# A -2h mtime is what makes a folder "old"; the file inside is written before the
# timestamp is forced back so the directory mtime is the one under test.
$plants = [ordered]@{
  oldLeaf = $oldLeaf; freshLeaf = $freshLeaf; otherPrefix = $otherPrefix
  fixedLegacy = $fixedLegacy; lookalike = $lookalike
}
$before = [ordered]@{}
foreach ($k in $plants.Keys) { $before[$k] = Test-Path $plants[$k] }

$out = Join-Path $logDir 'sweep-check-299.out.txt'
$err = Join-Path $logDir 'sweep-check-299.err.txt'
$sw = [Diagnostics.Stopwatch]::StartNew()
$p = Start-Process -FilePath $Bin -ArgumentList @('-test.run', '^TestPanelPageFetchesItsOwnAssets299$', '-test.v', '-test.timeout', '5m') `
  -WorkingDirectory $root -RedirectStandardOutput $out -RedirectStandardError $err -PassThru -NoNewWindow
$null = $p.WaitForExit(300000)
$sw.Stop()

$after = [ordered]@{}
foreach ($k in $plants.Keys) { $after[$k] = Test-Path $plants[$k] }
$expect = [ordered]@{ oldLeaf = $false; freshLeaf = $true; otherPrefix = $true; fixedLegacy = $true; lookalike = $true }
$pass = $true
foreach ($k in $plants.Keys) { if ($after[$k] -ne $expect[$k]) { $pass = $false } }
$summary = [ordered]@{
  tag = 'sweep-check'
  bin = $Bin
  rc = $p.ExitCode
  elapsed_s = [math]::Round($sw.Elapsed.TotalSeconds, 2)
  before = ($before | ConvertTo-Json -Compress)
  after = ($after | ConvertTo-Json -Compress)
  expect = ($expect | ConvertTo-Json -Compress)
  pass = $pass
}
$line = ($summary | ConvertTo-Json -Compress)
Add-Content -Path (Join-Path $PSScriptRoot 'results-sweep.jsonl') -Value $line -Encoding utf8
$line
