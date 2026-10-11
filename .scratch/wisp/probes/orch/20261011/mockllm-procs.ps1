param(
  [switch]$Quick,
  [string]$Names = '',
  [string]$CleanIds = ''
)
$ErrorActionPreference = 'Stop'
$goldenPat = '*Wisp\internal\llm\testdata\golden*'

function Get-Mine {
  @(Get-CimInstance Win32_Process -Filter "Name='mockllm.exe'" |
    Where-Object { [string]$_.CommandLine -like $goldenPat })
}

if ($Names -ne '') {
  $rows = @(Get-CimInstance Win32_Process -Filter "Name='$Names'")
  Write-Output ("name=" + $Names + " count=" + $rows.Count + " pids=[" + (($rows | ForEach-Object { $_.ProcessId }) -join ',') + "]")
  exit 0
}

if ($CleanIds -ne '') {
  $want = $CleanIds -split ',' | Where-Object { $_ -ne '' } | ForEach-Object { [int]$_ }
  $mine = Get-Mine
  foreach ($id in $want) {
    $r = @($mine | Where-Object { $_.ProcessId -eq $id })
    if ($r.Count -ne 1) { Write-Output ("skip pid=" + $id + " reason=not-in-repo-golden-roster"); continue }
    $ppid = $r[0].ParentProcessId
    $parentAlive = $false
    if ($ppid) { $parentAlive = [bool](Get-Process -Id $ppid -ErrorAction SilentlyContinue) }
    if ($parentAlive) { Write-Output ("skip pid=" + $id + " reason=parent-still-alive ppid=" + $ppid); continue }
    Stop-Process -Id $id -Force
    Write-Output ("killed pid=" + $id + " ppid=" + $ppid + " guards=name=mockllm.exe+repo-golden-path+parent-gone")
  }
  exit 0
}

$rows = @(Get-CimInstance Win32_Process -Filter "Name='mockllm.exe'")
$mine = @($rows | Where-Object { [string]$_.CommandLine -like $goldenPat })
if ($Quick) {
  Write-Output ("count_total=" + $rows.Count + " count_repo_golden=" + $mine.Count + " pids=[" + (($mine | ForEach-Object { $_.ProcessId }) -join ',') + "]")
  exit 0
}
foreach ($r in $rows) {
  $cl = [string]$r.CommandLine
  $m = if ($cl -like $goldenPat) { 'REPO_GOLDEN' } else { 'OTHER_GOLDEN' }
  $pp = 'NO_PARENT'
  if ($r.ParentProcessId) {
    if (Get-Process -Id $r.ParentProcessId -ErrorAction SilentlyContinue) { $pp = 'PARENT_ALIVE' } else { $pp = 'PARENT_GONE' }
  }
  $cpu = $null
  $p = Get-Process -Id $r.ProcessId -ErrorAction SilentlyContinue
  if ($p) { $cpu = [math]::Round($p.CPU, 2) }
  $port = ''
  if ($cl -match '-addr\s+127\.0\.0\.1:(\d+)') { $port = $Matches[1] }
  Write-Output ("pid=" + $r.ProcessId + " ppid=" + $r.ParentProcessId + " start=" + $r.CreationDate.ToString('yyyy-MM-dd HH:mm:ss') + " " + $m + " " + $pp + " cpu=" + $cpu + " addr_port=" + $port)
}
$tmp = @(Get-ChildItem -Path $env:TEMP -Directory -Filter 'Test*' -ErrorAction SilentlyContinue)
Write-Output ("count_total=" + $rows.Count + " count_repo_golden=" + $mine.Count + " temp_test_dirs=" + $tmp.Count)
Write-Output ("repo_golden_pids=[" + (($mine | ForEach-Object { $_.ProcessId }) -join ',') + "]")
