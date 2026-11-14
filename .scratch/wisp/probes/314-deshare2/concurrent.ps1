# 314-deshare2 control: two test processes started AT THE SAME MOMENT on the same
# -test.run selection, while a sampler reads every msedgewebview2 command line so
# the run can be asked the only question that matters for the shared-folder
# surface: "does one WebView2 user-data folder ever carry two top-level browser
# hosts in one tick?"
#
# Adapted from probes/314-deshare/concurrent.ps1 (ticket 314, landed three sites).
# The only changes: the trial tag defaults are this leg's, and the competitor count
# looks for this leg's own binary names.
#
# usage: pwsh -File concurrent.ps1 -Bin <test.exe> -Tag <tag> [-TestRun <regex>]
param(
  [Parameter(Mandatory=$true)][string]$Bin,
  [Parameter(Mandatory=$true)][string]$Tag,
  [string]$TestRun = '^Test253r6TheSilentDropBranchIsLoud$',
  [string]$TestRunA = '',
  [string]$TestRunB = '',
  [int]$SampleMs = 250,
  [int]$ExitBudgetSec = 300
)
$ErrorActionPreference = 'Stop'
if (-not $TestRunA) { $TestRunA = $TestRun }
if (-not $TestRunB) { $TestRunB = $TestRun }
$root = (Resolve-Path (Join-Path $PSScriptRoot '..\..\..\..')).Path
$logDir = Join-Path $PSScriptRoot 'logs'
$results = Join-Path $PSScriptRoot 'results-concurrent.jsonl'
New-Item -ItemType Directory -Force -Path $logDir | Out-Null

function Launch([string]$suffix, [string]$regex) {
  $out = Join-Path $logDir "$Tag-$suffix.out.txt"
  $err = Join-Path $logDir "$Tag-$suffix.err.txt"
  foreach ($f in @($out, $err)) { if (Test-Path $f) { Remove-Item $f -Force } }
  $sw = [Diagnostics.Stopwatch]::StartNew()
  $p = Start-Process -FilePath $Bin -ArgumentList @('-test.run', $regex, '-test.v', '-test.timeout', '10m') `
    -WorkingDirectory $root -RedirectStandardOutput $out -RedirectStandardError $err -PassThru -NoNewWindow
  return @{ proc = $p; sw = $sw; out = $out; err = $err }
}

$binName = [IO.Path]::GetFileNameWithoutExtension($Bin)
$competitorsAtStart = @(Get-Process -Name $binName -ErrorAction SilentlyContinue).Count
# Both starts are issued back to back; the skew between them is reported so the
# reader can see this is simultaneity, not a handoff.
$a = Launch 'A' $TestRunA
$b = Launch 'B' $TestRunB
$skewMs = [math]::Round((($b.sw.Elapsed - $a.sw.Elapsed).TotalMilliseconds), 1)
$testPidsOnBox = @(Get-Process -Name $binName -ErrorAction SilentlyContinue | ForEach-Object { $_.Id })

$seen = @{}
$ticks = @()
$deadline = (Get-Date).AddSeconds($ExitBudgetSec)
while ((-not $a.proc.HasExited -or -not $b.proc.HasExited) -and (Get-Date) -lt $deadline) {
  $tick = @()
  foreach ($q in (Get-CimInstance Win32_Process -Filter "Name='msedgewebview2.exe'" -ErrorAction SilentlyContinue)) {
    $m = [regex]::Match([string]$q.CommandLine, '--user-data-dir="?([^" ]+)"?')
    if (-not $m.Success) { continue }
    $udf = $m.Groups[1].Value
    if ($udf -notmatch 'wisp-' -or $udf -notmatch 'panel-profile') { continue }
    $tick += [pscustomobject]@{ browser = [int]$q.ProcessId; parent = [int]$q.ParentProcessId; udf = $udf }
    $key = "$($q.ProcessId)|$($q.ParentProcessId)|$udf"
    if (-not $seen.ContainsKey($key)) { $seen[$key] = (Get-Date).ToString('HH:mm:ss.fff') }
  }
  if ($tick.Count -gt 0) { $ticks += ,@($tick) }
  Start-Sleep -Milliseconds $SampleMs
}

$aHung = -not $a.proc.WaitForExit($ExitBudgetSec * 1000)
$bHung = -not $b.proc.WaitForExit($ExitBudgetSec * 1000)
if ($aHung) { Stop-Process -Id $a.proc.Id -Force -ErrorAction SilentlyContinue }
if ($bHung) { Stop-Process -Id $b.proc.Id -Force -ErrorAction SilentlyContinue }
$a.sw.Stop(); $b.sw.Stop()

function ReadCounts([string]$out, [string]$err, [bool]$hung) {
  $both = ''
  if (Test-Path $out) { $both += Get-Content $out -Raw }
  if (Test-Path $err) { $both += Get-Content $err -Raw }
  return [ordered]@{
    rc = if ($hung) { 'KILLED-EXITBUDGET' } else { $null }
    pass = ([regex]::Matches($both, '(?m)^--- PASS:')).Count
    fail = ([regex]::Matches($both, '(?m)^--- FAIL:')).Count
    skip = ([regex]::Matches($both, '(?m)^--- SKIP:')).Count
    panic = ([regex]::Matches($both, '(?m)^panic:')).Count
    bootFailed = ([regex]::Matches($both, 'webview2: .*boot failed')).Count
    embedFailed = ([regex]::Matches($both, 'WebView2 Embed failed')).Count
    accessViolation = ([regex]::Matches($both, '0xc0000005|access violation')).Count
  }
}
$ca = ReadCounts $a.out $a.err $aHung
$cb = ReadCounts $b.out $b.err $bHung
if (-not $aHung) { $ca.rc = $a.proc.ExitCode }
if (-not $bHung) { $cb.rc = $b.proc.ExitCode }

# A top-level browser host is one whose parent is not another browser in the same
# tick. Two hosts naming ONE folder in ONE tick is the shared-folder condition.
$coexist = @()
$maxHostsOnUdf = 0
$udfsByParentPid = @{}
$ticksWithHosts = 0
for ($i = 0; $i -lt $ticks.Count; $i++) {
  $t = $ticks[$i]
  $browserPids = @($t | ForEach-Object { $_.browser })
  $hosts = @($t | Where-Object { $browserPids -notcontains $_.parent })
  if ($hosts.Count -gt 0) { $ticksWithHosts++ }
  foreach ($h in $hosts) {
    if (-not $udfsByParentPid.ContainsKey($h.parent)) { $udfsByParentPid[$h.parent] = @{} }
    $udfsByParentPid[$h.parent][$h.udf] = $true
  }
  foreach ($g in ($hosts | Group-Object udf)) {
    if ($g.Count -gt $maxHostsOnUdf) { $maxHostsOnUdf = $g.Count }
    if ($g.Count -ge 2) {
      $parents = ($g.Group | ForEach-Object { $_.parent }) -join ','
      $coexist += "tick=$i udf=$($g.Name) topLevelHosts=$($g.Count) parents=$parents browsers=$(($g.Group | ForEach-Object { $_.browser }) -join ',')"
    }
  }
}
$udfListA = @(if ($udfsByParentPid.ContainsKey($a.proc.Id)) { $udfsByParentPid[$a.proc.Id].Keys } else { @() })
$udfListB = @(if ($udfsByParentPid.ContainsKey($b.proc.Id)) { $udfsByParentPid[$b.proc.Id].Keys } else { @() })
$udfPathsDistinct = ($udfListA.Count -gt 0 -and $udfListB.Count -gt 0 -and (@($udfListA | Where-Object { $udfListB -contains $_ }).Count -eq 0))

$summary = [ordered]@{
  tag = $Tag
  bin = $Bin
  testRun = $TestRun
  testRunA = $TestRunA
  testRunB = $TestRunB
  launchSkewMs = $skewMs
  a = [ordered]@{ pid = $a.proc.Id; rc = $ca.rc; elapsed_s = [math]::Round($a.sw.Elapsed.TotalSeconds, 2); pass = $ca.pass; fail = $ca.fail; skip = $ca.skip; panic = $ca.panic; bootFailed = $ca.bootFailed; embedFailed = $ca.embedFailed; accessViolation = $ca.accessViolation; udfs = ($udfListA -join ';') }
  b = [ordered]@{ pid = $b.proc.Id; rc = $cb.rc; elapsed_s = [math]::Round($b.sw.Elapsed.TotalSeconds, 2); pass = $cb.pass; fail = $cb.fail; skip = $cb.skip; panic = $cb.panic; bootFailed = $cb.bootFailed; embedFailed = $cb.embedFailed; accessViolation = $cb.accessViolation; udfs = ($udfListB -join ';') }
  competitorsAtStart = $competitorsAtStart
  testPidsOnBoxAtLaunch = ($testPidsOnBox -join ',')
  sampleTicks = $ticks.Count
  ticksWithHosts = $ticksWithHosts
  browsersSeen = $seen.Count
  maxTopLevelHostsOnOneUdf = $maxHostsOnUdf
  sameUdfCoexistenceTicks = $coexist.Count
  sameUdfCoexistence = if ($coexist.Count -gt 0) { $coexist -join ' ; ' } else { 'none' }
  udfPathsDistinctBetweenRuns = $udfPathsDistinct
  browserDetail = (($seen.Keys | Sort-Object) | ForEach-Object {
      $p = $_.Split('|'); "browser=$($p[0]) parent=$($p[1]) udf=$($p[2]) first=$($seen[$_])" }) -join ' ;; '
}
$line = ($summary | ConvertTo-Json -Compress -Depth 5)
Add-Content -Path $results -Value $line -Encoding utf8
$line
