Get-CimInstance Win32_Process -Filter "Name='powershell.exe'" | ForEach-Object {
    "PID=" + $_.ProcessId + " created=" + $_.CreationDate + " ppid=" + $_.ParentProcessId
    "  CMD: " + $_.CommandLine
}
""
"=== OpenWith.exe parents ==="
Get-CimInstance Win32_Process -Filter "Name='OpenWith.exe'" | ForEach-Object {
    "PID=" + $_.ProcessId + " created=" + $_.CreationDate + " ppid=" + $_.ParentProcessId
    "  CMD: " + $_.CommandLine
}
""
"=== what are those parents ==="
$ppids = (Get-CimInstance Win32_Process -Filter "Name='OpenWith.exe'").ParentProcessId | Sort-Object -Unique
foreach ($p in $ppids) {
    Get-CimInstance Win32_Process -Filter "ProcessId=$p" | ForEach-Object {
        "parent PID=" + $_.ProcessId + " name=" + $_.Name + " created=" + $_.CreationDate
        "  CMD: " + $_.CommandLine
    }
}
