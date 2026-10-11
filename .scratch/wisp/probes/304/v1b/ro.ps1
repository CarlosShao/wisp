param([string]$Tree,[switch]$Kill)
# 304-v1b own roster for mockllm.exe under a named tree path.
# Attribution = command line contains the named tree path (never count, never parent-PID existence).
if (-not $Tree) { Write-Output "usage: ro.ps1 -Tree <windows path> [-Kill]"; exit 2 }
$t = $Tree.ToLower()
$procs = @(Get-CimInstance Win32_Process -Filter "Name='mockllm.exe'" | Where-Object { $_.CommandLine -and $_.CommandLine.ToLower().Contains($t) })
foreach ($p in $procs) {
  $parent = Get-CimInstance Win32_Process -Filter ("ProcessId=" + $p.ParentProcessId) -ErrorAction SilentlyContinue
  $pn = if ($parent) { $parent.Name } else { "GONE" }
  Write-Output ("hit pid=" + $p.ProcessId + " ppid=" + $p.ParentProcessId + " parent=" + $pn + " cmd=" + $p.CommandLine)
}
Write-Output ("count_tree=" + $procs.Count)
if ($Kill) {
  foreach ($p in $procs) {
    if ($p.Name -ne 'mockllm.exe') { Write-Output ("skip pid=" + $p.ProcessId + " reason=name"); continue }
    $parent = Get-CimInstance Win32_Process -Filter ("ProcessId=" + $p.ParentProcessId) -ErrorAction SilentlyContinue
    if ($parent) { Write-Output ("skip pid=" + $p.ProcessId + " reason=parent-present(" + $parent.Name + ")"); continue }
    Stop-Process -Id $p.ProcessId -Force
    Write-Output ("killed pid=" + $p.ProcessId + " guards=name=mockllm.exe+tree-path+parent-gone")
  }
  $left = @(Get-CimInstance Win32_Process -Filter "Name='mockllm.exe'" | Where-Object { $_.CommandLine -and $_.CommandLine.ToLower().Contains($t) })
  Write-Output ("count_tree_after_kill=" + $left.Count)
}
