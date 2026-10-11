Get-CimInstance Win32_Process -Filter "Name like 'go%' or Name like 'git%'" | ForEach-Object {
  Write-Output ("PID=" + $_.ProcessId + " NAME=" + $_.Name + " PARENT=" + $_.ParentProcessId + " CMD=" + $_.CommandLine)
}
