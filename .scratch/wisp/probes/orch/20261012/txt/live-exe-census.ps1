$ErrorActionPreference = 'Stop'
$rows = Get-CimInstance Win32_Process |
    Where-Object { $_.Name -like '*.test.exe' -or $_.Name -eq 'wisp.exe' -or $_.Name -eq 'balldebug.exe' }
"family_total=$($rows.Count)"
foreach ($p in $rows) {
    "pid=$($p.ProcessId) name=$($p.Name) created=$($p.CreationDate.ToString('MM-dd_HH:mm:ss'))"
    "  cmdline=$($p.CommandLine)"
}
$all = Get-CimInstance Win32_Process
"all_processes=$($all.Count)"
"svchost=$(($all | Where-Object { $_.Name -eq 'svchost.exe' }).Count)"
