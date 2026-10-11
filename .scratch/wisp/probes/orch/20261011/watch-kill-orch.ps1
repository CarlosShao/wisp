param(
  [Parameter(Mandatory=$true)][int]$ParentPid,
  [string]$Label = 'shot'
)
# 编排者的形 (c) 尺：在同一个进程里"盯住那枚孩子 → 一看到就外力杀掉父进程 → 回头看孩子还在不在"。
# 为什么不用 bash 轮询：一发 bash↔powershell 往返要 0.5-1 s，而孩子只在用例那几秒里活着，慢一拍就读成"没起世界"。
$ErrorActionPreference = 'Stop'
$goldenPat = '*Wisp\internal\llm\testdata\golden*'
$found = $null
$deadline = (Get-Date).AddSeconds(60)
while ($null -eq $found -and (Get-Date) -lt $deadline) {
  $found = @(Get-CimInstance Win32_Process -Filter "Name='mockllm.exe'" |
             Where-Object { $_.ParentProcessId -eq $ParentPid -and $_.CommandLine -like $goldenPat }) | Select-Object -First 1
  if ($null -eq $found) { Start-Sleep -Milliseconds 150 }
}
if ($null -eq $found) {
  Write-Output ("$label world=NOT-CONSTRUCTED parent=$ParentPid child=none  ⇒ 这一发读数作废，⛔ 拿 0 当绿")
  exit 2
}
$childPid = [int]$found.ProcessId
Write-Output ("$label world=CONSTRUCTED parent=$ParentPid child=$childPid child_age_sec=" +
  ([math]::Round(((Get-Date) - $found.CreationDate).TotalSeconds, 2)))
Stop-Process -Id $ParentPid -Force
Write-Output ("$label external_kill=issued at=$(Get-Date -Format 'HH:mm:ss.fff')")
Start-Sleep -Milliseconds 1500
$still = @(Get-CimInstance Win32_Process -Filter "Name='mockllm.exe'" |
           Where-Object { $_.ProcessId -eq $childPid })
$anyGolden = @(Get-CimInstance Win32_Process -Filter "Name='mockllm.exe'" |
               Where-Object { $_.CommandLine -like $goldenPat })
Write-Output ("$label child_after_kill=" + $(if ($still.Count -gt 0) { 'STILL-ALIVE(=leak)' } else { 'GONE(=reaped)' }) +
  " repo_golden_now=" + $anyGolden.Count)
if ($still.Count -gt 0) { exit 1 }
exit 0
