# Read-only follow-up on the two leftover mockllm.exe found by forensics2 (R2).
# No kill. Only: age, parent chain, whether their TempDir still exists, listening sockets.
$ErrorActionPreference = 'Stop'
$out = Join-Path $PSScriptRoot 'forensics-3.txt'
$L = @()
$L += 'ruler = Get-CimInstance Win32_Process + Test-Path + Get-NetTCPConnection (read-only)'
$L += 'captured_local = ' + (Get-Date -Format 'yyyy-MM-dd HH:mm:ss zzz')
$L += ''

$all = @(Get-CimInstance Win32_Process)
$byPid = @{}
foreach ($p in $all) { $byPid[[int]$p.ProcessId] = $p }

$mine = @($all | Where-Object { $_.Name -eq 'mockllm.exe' })
$L += 'mockllm.exe count = ' + $mine.Count
$L += ''

foreach ($p in $mine) {
    $L += ('--- pid=' + $p.ProcessId + ' ---')
    $age = ((Get-Date) - [datetime]$p.CreationDate)
    $L += ('  started        = ' + ([datetime]$p.CreationDate).ToString('yyyy-MM-dd HH:mm:ss') + '   ageMinutes = ' + [math]::Round($age.TotalMinutes, 0))
    $L += ('  workingSetMB   = ' + [math]::Round($p.WorkingSetSize / 1MB, 1) + '   privateMB = ' + [math]::Round($p.PrivatePageCount / 1MB, 1))
    $L += ('  cpuSeconds     = ' + [math]::Round([int64]$p.UserModeTime / 10000000 + [int64]$p.KernelModeTime / 10000000, 1))
    $ppid = [int]$p.ParentProcessId
    $parent = $byPid[$ppid]
    if ($null -eq $parent) {
        $L += ('  parent pid=' + $ppid + ' -> GONE (parent no longer exists)')
    } else {
        $pc = $parent.CommandLine
        if ($null -eq $pc) { $pc = '<null>' }
        if ($pc.Length -gt 260) { $pc = $pc.Substring(0, 260) + ' ...' }
        $L += ('  parent pid=' + $ppid + ' -> ALIVE name=' + $parent.Name + ' cmd=' + $pc)
    }
    $L += ('  cmdFull        = ' + $p.CommandLine)
    if ($p.CommandLine -match '(C:\\Users\\swq\\AppData\\Local\\Temp\\[^"]*?\\00\\[0-9]+)') { }
    $m = [regex]::Match($p.CommandLine, 'C:\\Users\\swq\\AppData\\Local\\Temp\\[^"]+?\\00[0-9]')
    if ($m.Success) {
        $td = $m.Value
        $L += ('  tempdir        = ' + $td)
        $L += ('  tempdir exists = ' + (Test-Path -LiteralPath $td))
        $L += ('  testname       = ' + ((Split-Path $td -Leaf) -replace '[0-9]{6,}$', ''))
        $parentTemp = Split-Path $td -Parent
        $kids = @(Get-ChildItem -LiteralPath $parentTemp -ErrorAction SilentlyContinue | Where-Object { $_.Name -like 'Test*' })
        $L += ('  sibling Test* dirs under ' + $parentTemp + ' = ' + $kids.Count)
    }
    $cons = @(Get-NetTCPConnection -OwningProcess $p.ProcessId -ErrorAction SilentlyContinue)
    if ($cons.Count -eq 0) { $L += '  tcp sockets    = 0' }
    foreach ($c in $cons) {
        $L += ('  tcp socket     = ' + $c.LocalAddress + ':' + $c.LocalPort + ' state=' + $c.State)
    }
    $L += ''
}

$L += '== how many Temp dirs named Test* are sitting in the profile temp root ==
'
$tmpRoot = 'C:\Users\swq\AppData\Local\Temp'
$tds = @(Get-ChildItem -LiteralPath $tmpRoot -Directory -ErrorAction SilentlyContinue | Where-Object { $_.Name -like 'Test*' })
$L += ('  count = ' + $tds.Count)
$L += '  (no recursive size scan: this profile temp root carried 2.4 GB of snapshots, see the /tmp census)'
$L += '  oldest five by CreationDate:'
foreach ($d in ($tds | Sort-Object CreationDate | Select-Object -First 5)) {
    $L += ('    ' + ('' + $d.CreationDate) + '  ' + $d.Name)
}
$L += '  newest five by CreationDate:'
foreach ($d in ($tds | Sort-Object CreationDate -Descending | Select-Object -First 5)) {
    $L += ('    ' + ('' + $d.CreationDate) + '  ' + $d.Name)
}
$L += ''
$L += 'end_of_forensics_3'
$L | Out-File -FilePath $out -Encoding utf8
Write-Output ('wrote ' + $out)
