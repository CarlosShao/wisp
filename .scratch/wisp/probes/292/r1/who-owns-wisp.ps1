Get-CimInstance Win32_Process -Filter "Name='wisp.exe'" |
  Select-Object ProcessId, ParentProcessId, CreationDate, ExecutablePath, CommandLine |
  Format-List
Write-Output "--- parent chain ---"
$pp = (Get-CimInstance Win32_Process -Filter "Name='wisp.exe'" | Select-Object -First 1).ParentProcessId
Get-CimInstance Win32_Process -Filter "ProcessId=$pp" |
  Select-Object ProcessId, Name, CreationDate, CommandLine | Format-List
