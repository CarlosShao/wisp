$found = Get-Process | Where-Object { $_.ProcessName -like '*.test' }
$count = ($found | Measure-Object).Count
Write-Output ("dotnet-test-family-process-count=" + $count)
if ($count -gt 0) { $found | ForEach-Object { Write-Output ("pid=" + $_.Id + " name=" + $_.ProcessName) } }
