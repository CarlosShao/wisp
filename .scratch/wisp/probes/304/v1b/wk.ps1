param([string]$Tree,[int]$TimeoutSec=90)
# 304-v1b: single-process watcher+killer. Attribution = named tree path in command line
# (test binary = ac2t.exe under my tree; child = mockllm.exe carrying the tree's golden dir).
$t = $Tree.ToLower()
$sw = [Diagnostics.Stopwatch]::StartNew()
while ($sw.Elapsed.TotalSeconds -lt $TimeoutSec) {
  $p = @(Get-CimInstance Win32_Process -Filter "Name='ac2t.exe'" | Where-Object { $_.CommandLine -and $_.CommandLine.ToLower().Contains($t) })
  if ($p.Count -gt 0) {
    $mock = @(Get-CimInstance Win32_Process -Filter "Name='mockllm.exe'" | Where-Object { $_.CommandLine -and $_.CommandLine.ToLower().Contains($t) })
    if ($mock.Count -gt 0) {
      $kids = ($mock | ForEach-Object { $_.ProcessId }) -join ','
      Write-Output ("world=CONSTRUCTED child_pids=" + $kids + " t=" + [math]::Round($sw.Elapsed.TotalSeconds,2))
      foreach ($x in $p) {
        Write-Output ("kill_testbinary pid=" + $x.ProcessId + " parent_alive_at_kill=YES")
        Stop-Process -Id $x.ProcessId -Force
      }
      exit 0
    }
  }
  Start-Sleep -Milliseconds 40
}
Write-Output "world=NOT-CONSTRUCTED (timeout, no reading)"
exit 3
