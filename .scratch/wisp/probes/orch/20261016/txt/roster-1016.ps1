Get-CimInstance Win32_Process |
  Where-Object { $_.Name -like '*.test.exe' -or $_.Name -eq 'go.exe' -or $_.Name -eq 'wisp.exe' -or $_.Name -like 'gofumpt*' } |
  ForEach-Object { '{0} {1} ppid={2}' -f $_.ProcessId, $_.Name, $_.ParentProcessId }
Write-Output '--- counts'
$all = Get-CimInstance Win32_Process | Where-Object { $_.Name -like '*.test.exe' -or $_.Name -eq 'go.exe' -or $_.Name -eq 'wisp.exe' }
Write-Output ('TOTAL=' + ($all | Measure-Object).Count)
$all | ForEach-Object { ($_.CommandLine) } | Where-Object { $_ } | ForEach-Object { 'CMD: ' + $_ }
