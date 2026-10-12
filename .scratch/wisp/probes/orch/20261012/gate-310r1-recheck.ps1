$ts = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
Write-Output "ts=$ts"
Write-Output "--- positive control (roster 1: never read an empty output as 0) ---"
Write-Output ("svchost_count=" + (@(Get-Process -Name svchost -ErrorAction SilentlyContinue).Count))
Write-Output ("all_exe_count=" + (@(Get-Process -ErrorAction SilentlyContinue).Count))
Write-Output "--- four exact names, tasklist each one shot ---"
foreach ($n in @("wisp.exe", "balldebug.exe", "mockllm.exe", "wisp.test.exe")) {
  $c = @(tasklist //FI "IMAGENAME eq $n" //NH 2>$null | Where-Object { $_ -ne "" }).Count
  Write-Output ("exact " + $n + " count=" + $c)
}
Write-Output "--- family ruler (canonical verbatim form, Get-CimInstance) ---"
$family = @(Get-CimInstance Win32_Process -Filter "Name like '%.test.exe'" -ErrorAction SilentlyContinue)
Write-Output ("family_count=" + $family.Count)
Write-Output "--- temp rosters ---"
$t = @(Get-ChildItem -Path $env:TEMP -Directory -ErrorAction SilentlyContinue | Where-Object { $_.Name -like "Test*" })
Write-Output ("temp_test_star=" + $t.Count)
Write-Output ("TEMP_root=" + $env:TEMP)
$la = @(Get-ChildItem -Path (Join-Path $env:LOCALAPPDATA "wisp-test-gotmp") -Directory -ErrorAction SilentlyContinue)
Write-Output ("localappdata_wisp_test_gotmp_roots=" + $la.Count)
$dt = @(Get-ChildItem -Path "D:\tmp\wisp-test-gotmp" -Directory -ErrorAction SilentlyContinue)
Write-Output ("d_tmp_wisp_test_gotmp_roots=" + $dt.Count)
