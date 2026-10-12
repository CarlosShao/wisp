$t = @(Get-ChildItem -Path $env:TEMP -Directory -ErrorAction SilentlyContinue | Where-Object { $_.Name -like "Test*" })
Write-Output ("temp_test_star=" + $t.Count + " positive_control_all_temp_dirs=" + @(Get-ChildItem -Path $env:TEMP -Directory -ErrorAction SilentlyContinue).Count)
$one = @($t | Where-Object { $_.Name -eq "TestRoundTripLoadMarshalLoad4222770164" })
Write-Output ("leg_named_leak_present=" + $one.Count + " created=" + $(if ($one.Count -gt 0) { $one[0].CreationTime.ToString("yyyy-MM-dd HH:mm:ss") } else { "NONE" }))
$probe = @($t | Where-Object { $_.Name -eq "TestProbeStepByStep3843095677" })
Write-Output ("unowned_probe_still_present=" + $probe.Count)
Write-Output "--- roots named by the leg's 5 entry readings ---"
$la = @(Get-ChildItem -Path (Join-Path $env:LOCALAPPDATA "wisp-test-gotmp") -Directory -ErrorAction SilentlyContinue)
$d = @(Get-ChildItem -Path "D:\tmp\wisp-test-gotmp" -Directory -ErrorAction SilentlyContinue)
Write-Output ("lappdata_roots=" + $la.Count + " d_roots=" + $d.Count)
Write-Output "--- residue dirs inside the leg's own roots (named by its report) ---"
foreach ($n in @("TestRoundTripLoadMarshalLoad3712892429")) {
  $hit = @($la | Where-Object { @(Get-ChildItem -Path $_.FullName -Directory -ErrorAction SilentlyContinue | Where-Object { $_.Name -eq $n }).Count -gt 0 })
  Write-Output ($n + " found_in_roots=" + $hit.Count)
}
