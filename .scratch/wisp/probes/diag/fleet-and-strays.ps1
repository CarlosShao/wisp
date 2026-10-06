"=== 1. 机器负载（能不能继续派活） ==="
$os = Get-CimInstance Win32_OperatingSystem
$totMB = [math]::Round($os.TotalVisibleMemorySize/1KB)
$freeMB = [math]::Round($os.FreePhysicalMemory/1KB)
"内存: 总 " + $totMB + " MB / 空闲 " + $freeMB + " MB / 已用 " + [math]::Round((1-$freeMB/$totMB)*100) + "%"
$cpu = (Get-CimInstance Win32_Processor | Measure-Object -Property LoadPercentage -Average).Average
"CPU 负载: " + $cpu + "%"

"=== 2. 我造出来的挂死 powershell（命令行里没有真参数、含 'Wisp.CPU' 这类坏字样的） ==="
Get-CimInstance Win32_Process -Filter "Name='powershell.exe'" | ForEach-Object {
    $c = $_.CommandLine
    if ($null -eq $c) { $c = "(空)" }
    $bad = if ($c -match 'Wisp\.(CPU|ProcessName|WorkingSet64|Id|StartTime)') { "<== 坏的（$_ 被 bash 抢着展开了）" } else { "" }
    "PID=" + $_.ProcessId + " 起于 " + $_.CreationDate + " " + $bad
}

"=== 3. 挂死的测试相关进程 ==="
foreach ($n in @('mockllm.exe','wisp.exe','go.exe','go','balldebug.exe')) {
    $ps = @(Get-CimInstance Win32_Process -Filter "Name='$n'")
    "进程 " + $n + " 枚数=" + $ps.Count
    foreach ($p in $ps) {
        "   PID=" + $p.ProcessId + " 起于 " + $p.CreationDate
        "     CMD: " + $p.CommandLine
    }
}

"=== 4. 弹窗本体 ==="
$ow = @(Get-CimInstance Win32_Process -Filter "Name='OpenWith.exe'")
"OpenWith.exe 枚数=" + $ow.Count
foreach ($p in $ow) { "   PID=" + $p.ProcessId + " 起于 " + $p.CreationDate + " CMD: " + $p.CommandLine }

"=== 5. 当前有没有窗口标题里带 '打开' / 'open' 的（弹窗是否真的可见） ==="
Get-CimInstance Win32_Process -Filter "Name='OpenWith.exe'" | ForEach-Object {
    "   OpenWith PID=" + $_.ProcessId
}
