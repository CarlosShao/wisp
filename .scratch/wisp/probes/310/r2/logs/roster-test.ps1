$ErrorActionPreference = "Stop"
$all = Get-ChildItem -Path $env:TEMP -Directory
$test = $all | Where-Object { $_.Name -like "Test*" } | Sort-Object Name
"POSITIVE_CONTROL_all_top_level_dirs=" + $all.Count
"COUNT=" + $test.Count
foreach ($d in $test) {
    "NAME`t" + $d.Name + "`tCREATED`t" + $d.CreationTime.ToString("yyyy-MM-dd HH:mm:ss.fff")
}
