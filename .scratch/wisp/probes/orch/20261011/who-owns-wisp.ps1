Get-CimInstance Win32_Process -Filter "Name like 'wisp%' or Name like 'balldebug%' or Name like 'msedgewebview2%'" | ForEach-Object {
  $p = Get-CimInstance Win32_Process -Filter "ProcessId=$($_.ParentProcessId)" -ErrorAction SilentlyContinue
  Write-Output ("PID=" + $_.ProcessId + " NAME=" + $_.Name + " START=" + $_.CreationDate + " PARENT=" + $_.ParentProcessId + " PARENTNAME=" + $p.Name + " EXE=" + $_.ExecutablePath)
  Write-Output ("   CMD=" + $_.CommandLine)
}
Write-Output "--- counts by name ---"
Get-Process -Name wisp, balldebug, msedgewebview2 -ErrorAction SilentlyContinue | Group-Object Name | ForEach-Object { Write-Output ($_.Name + "=" + $_.Count) }
