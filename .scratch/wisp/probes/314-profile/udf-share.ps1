# 314 profile UDF-sharing probe: start the test binary and, WHILE it runs, sample the
# msedgewebview2 processes that name a wisp-33r* user-data folder, so the sample shows
# which test pid OWNS the browser a later run's window is bound to.
# One test process at a time: this script's own poller is not a test.
# usage: pwsh -File udf-share.ps1 -Bin <bin> -Tag A1share
param(
  [Parameter(Mandatory=$true)][string]$Bin,
  [Parameter(Mandatory=$true)][string]$Tag
)
$ErrorActionPreference = 'Stop'
$root = (Resolve-Path (Join-Path $PSScriptRoot '..\..\..\..')).Path
$logDir = Join-Path $PSScriptRoot 'logs'
$re = '^(TestAC14GoSideEvalPushReachesThePage|TestAC13BringUpSurvivesAReusedThreadQuit|TestPanelThreadIsSTAAndExitsCleanly|TestAC13BringUpRefusesAThreadWithAQueuedClose)$'
function Launch($suffix) {
  $out = Join-Path $logDir "$Tag-$suffix.out.txt"
  $err = Join-Path $logDir "$Tag-$suffix.err.txt"
  if (Test-Path $out) { Remove-Item $out -Force }
  if (Test-Path $err) { Remove-Item $err -Force }
  $sw = [Diagnostics.Stopwatch]::StartNew()
  $p = Start-Process -FilePath $Bin -ArgumentList @('-test.run', $re, '-test.v', '-test.timeout', '10m') `
    -WorkingDirectory $root -RedirectStandardOutput $out -RedirectStandardError $err -PassThru -NoNewWindow
  return @{ proc = $p; sw = $sw; out = $out; err = $err }
}
# run 1, unattended: leaves (or cleans) the profile the way the group under test does
$a = Launch 'warm'
$null = $a.proc.WaitForExit(300000)
$warmRc = $a.proc.ExitCode
$warmMs = [math]::Round($a.sw.Elapsed.TotalSeconds, 2)
# run 2, with the sampler watching from the first millisecond
$b = Launch 'sample'
$seen = @{}
$ticks = @()
$deadline = (Get-Date).AddSeconds(25)
while (-not $b.proc.HasExited -and (Get-Date) -lt $deadline) {
  $tick = @()
  foreach ($q in (Get-CimInstance Win32_Process -Filter "Name='msedgewebview2.exe'" -ErrorAction SilentlyContinue)) {
    $m = [regex]::Match($q.CommandLine, '--user-data-dir="?([^" ]+)"?')
    if (-not $m.Success) { continue }
    $udf = $m.Groups[1].Value
    if ($udf -notmatch 'wisp-33r[567]-') { continue }
    $tick += [pscustomobject]@{ browser = [int]$q.ProcessId; parent = [int]$q.ParentProcessId; udf = $udf }
    $key = "$($q.ProcessId)|$($q.ParentProcessId)|$udf"
    if (-not $seen.ContainsKey($key)) { $seen[$key] = (Get-Date).ToString('HH:mm:ss.fff') }
  }
  if ($tick.Count -gt 0) { $ticks += ,@($tick) }
  Start-Sleep -Milliseconds 300
}
$hung = -not $b.proc.WaitForExit(300000)
if ($hung) { Stop-Process -Id $b.proc.Id -Force -ErrorAction SilentlyContinue }
$b.sw.Stop()
$sampleRc = if ($hung) { 'KILLED-300s' } else { $b.proc.ExitCode }
$rows = @()
foreach ($k in ($seen.Keys | Sort-Object)) {
  $parts = $k.Split('|')
  $rows += [pscustomobject]@{ browserPid = $parts[0]; parentPid = $parts[1]; udf = $parts[2]; firstSeen = $seen[$k] }
}
# Coexistence: a top-level browser host is one whose parent is not another browser in
# the same tick. Two hosts on ONE user-data folder in ONE tick is the shared-folder
# condition the ticket asks about, observed without running two tests at once.
$warmPid = $a.proc.Id
$samplePid = $b.proc.Id
$coexist = @()
$maxHosts = 0
$hostsPerUdf = @{}
for ($i = 0; $i -lt $ticks.Count; $i++) {
  $t = $ticks[$i]
  $allBrowserPids = @($t | ForEach-Object { $_.browser })
  $hosts = @($t | Where-Object { $allBrowserPids -notcontains $_.parent })
  if ($hosts.Count -gt $maxHosts) { $maxHosts = $hosts.Count }
  foreach ($g in ($hosts | Group-Object udf)) {
    $n = $g.Count
    if (-not $hostsPerUdf.ContainsKey($g.Name) -or $hostsPerUdf[$g.Name] -lt $n) { $hostsPerUdf[$g.Name] = $n }
    $w = @($g.Group | Where-Object { $_.parent -eq $warmPid }).Count
    $s = @($g.Group | Where-Object { $_.parent -eq $samplePid }).Count
    if ($w -gt 0 -and $s -gt 0) { $coexist += "tick=$i udf=$($g.Name) warmHosts=$w sampleHosts=$s" }
  }
}
$summary = [ordered]@{
  tag = $Tag
  warmRc = $warmRc; warmElapsed_s = $warmMs
  sampleRc = $sampleRc; sampleElapsed_s = [math]::Round($b.sw.Elapsed.TotalSeconds, 2)
  warmTestPid = $a.proc.Id; sampleTestPid = $b.proc.Id
  browsersSeenDuringSample = $rows.Count
  sampleTicks = $ticks.Count
  maxHostsInOneTick = $maxHosts
  maxHostsOnOneUdf = (($hostsPerUdf.Values | Measure-Object -Maximum).Maximum)
  coexistenceSameUdf = if ($coexist.Count -gt 0) { $coexist -join ' ; ' } else { 'none' }
  browserDetail = ($rows | ForEach-Object { "browser=$($_.browserPid) parent=$($_.parentPid) udf=$($_.udf) first=$($_.firstSeen)" }) -join ' ;; '
}
$line = ($summary | ConvertTo-Json -Compress)
Add-Content -Path (Join-Path $PSScriptRoot 'results-udf-share.jsonl') -Value $line -Encoding utf8
$line
