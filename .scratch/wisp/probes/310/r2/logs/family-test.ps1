$p = Get-CimInstance Win32_Process | Where-Object { $_.Name -like "*.test.exe" }
"FAMILY_TEST_EXE_COUNT=" + @($p).Count
foreach ($x in $p) {
    "PID=" + $x.ProcessId + " PPID=" + $x.ParentProcessId + " NAME=" + $x.Name
}
