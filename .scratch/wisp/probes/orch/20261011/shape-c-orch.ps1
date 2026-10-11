param(
  [Parameter(Mandatory=$true)][string]$Exe,
  [Parameter(Mandatory=$true)][string[]]$TestArgs,
  [Parameter(Mandatory=$true)][string]$Label,
  [string]$WorkDir = 'D:\work\workspace\projects plans\Wisp'
)
# 编排者的形 (c) 一发到底：起测试二进制 → 盯住它的 mockllm 孩子（命令行含本仓 golden 目录才算）→
# 一看到就外力 Stop-Process 那枚二进制（＝腿被杀／shell 超时／Ctrl+C 那一形）→ 回头看孩子还在不在。
# 起不来世界就 exit 2 并把这一发标成作废，⛔ 拿"后来没人"读成"守住了"。
$ErrorActionPreference = 'Stop'
$goldenPat = '*Wisp\internal\llm\testdata\golden*'
$out = "D:\tmp\wisp-orch-20261011\shape-$Label.out.txt"
$p = Start-Process -FilePath $Exe -ArgumentList $TestArgs -WorkingDirectory $WorkDir `
     -PassThru -NoNewWindow -RedirectStandardOutput $out -RedirectStandardError "$out.err"
Write-Output "$Label parent=$($p.Id) started"
$found = $null
$deadline = (Get-Date).AddSeconds(90)
while ($null -eq $found -and (Get-Date) -lt $deadline) {
  $found = @(Get-CimInstance Win32_Process -Filter "Name='mockllm.exe'" |
             Where-Object { $_.ParentProcessId -eq $p.Id -and $_.CommandLine -like $goldenPat }) | Select-Object -First 1
  if ($null -eq $found) { Start-Sleep -Milliseconds 120 }
}
if ($null -eq $found) {
  $alive = @(Get-Process -Id $p.Id -ErrorAction SilentlyContinue).Count
  Write-Output "$Label world=NOT-CONSTRUCTED child=none parent_alive_now=$alive  ⇒ 读数作废"
  if ($alive -gt 0) { Stop-Process -Id $p.Id -Force }
  exit 2
}
$childPid = [int]$found.ProcessId
Write-Output "$Label world=CONSTRUCTED child=$childPid (confirmed alive while parent ran)"
Stop-Process -Id $p.Id -Force
Write-Output "$Label external_kill_issued"
Start-Sleep -Milliseconds 2000
$still = @(Get-CimInstance Win32_Process -Filter "Name='mockllm.exe'" | Where-Object { $_.ProcessId -eq $childPid }).Count
$anyGolden = @(Get-CimInstance Win32_Process -Filter "Name='mockllm.exe'" | Where-Object { $_.CommandLine -like $goldenPat }).Count
Write-Output "$Label child_after_kill=$(if ($still -gt 0) {'STILL-ALIVE=LEAK'} else {'GONE=reaped'}) repo_golden_now=$anyGolden"
if ($still -gt 0) { exit 1 }
exit 0
