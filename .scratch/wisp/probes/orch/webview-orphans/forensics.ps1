# Read-only forensics for the orphan msedgewebview2.exe population.
# Does NOT kill, start, or modify anything. Only queries CIM.
$ErrorActionPreference = 'Stop'
$out = Join-Path $PSScriptRoot 'forensics-1.txt'

$L = @()
$L += 'ruler = Get-CimInstance Win32_Process (read-only, no kill / no start / no modify)'
$L += 'captured_local = ' + (Get-Date -Format 'yyyy-MM-dd HH:mm:ss zzz')
$L += ''

$all = Get-CimInstance Win32_Process
$wv = $all | Where-Object { $_.Name -eq 'msedgewebview2.exe' }

$byPid = @{}
foreach ($p in $all) { $byPid[[int]$p.ProcessId] = $p }

$L += '== A. totals =='
$L += 'msedgewebview2.exe count = ' + (@($wv).Count)
$sumMb = [math]::Round((($wv | Measure-Object -Property WorkingSetSize -Sum).Sum) / 1MB, 1)
$L += 'sum working set MB = ' + $sumMb
$privMb = [math]::Round((($wv | Measure-Object -Property PrivatePageCount -Sum).Sum) / 1MB, 1)
$L += 'sum private bytes MB = ' + $privMb
$L += ''

$L += '== B. group by --type= token =='
$groups = $wv | Group-Object {
    if ($_.CommandLine -match '--type=([a-z\-]+)') { $Matches[1] } else { 'BROWSER(root)' }
}
foreach ($g in ($groups | Sort-Object Count -Descending)) {
    $L += ('  {0,-16} count={1}' -f $g.Name, $g.Count)
}
$L += ''

$L += '== C. per-PID row: pid, ppid, parent alive?, parent name, type, age, WS MB =='
$rows = @()
foreach ($p in ($wv | Sort-Object CreationDate)) {
    $ppid = [int]$p.ParentProcessId
    $parent = $byPid[$ppid]
    $palive = if ($null -ne $parent) { 'YES' } else { 'NO ' }
    $pname = if ($null -ne $parent) { $parent.Name } else { '-' }
    $t = if ($p.CommandLine -match '--type=([a-z\-]+)') { $Matches[1] } else { 'BROWSER' }
    $age = ((Get-Date) - [datetime]$p.CreationDate)
    $ws = [math]::Round($p.WorkingSetSize / 1MB, 1)
    $rows += [pscustomobject]@{
        Pid = $p.ProcessId; Ppid = $ppid; PAlive = $palive; PName = $pname
        Type = $t; Started = ([datetime]$p.CreationDate).ToString('MM-dd HH:mm:ss')
        AgeMin = [math]::Round($age.TotalMinutes, 0); WSMb = $ws
    }
}
$rows | ForEach-Object {
    $L += ('  pid={0,-7} ppid={1,-7} parentAlive={2} parentName={3,-22} type={4,-16} started={5} ageMin={6,-7} wsMb={7}' -f `
        $_.Pid, $_.Ppid, $_.PAlive, $_.PName, $_.Type, $_.Started, $_.AgeMin, $_.WSMb)
}
$L += ''

$L += '== D. root(s): a webview process whose parent is NOT another webview process =='
foreach ($p in $wv) {
    $ppid = [int]$p.ParentProcessId
    $parent = $byPid[$ppid]
    $pName = if ($null -ne $parent) { $parent.Name } else { '<GONE>' }
    if ($pName -ne 'msedgewebview2.exe') {
        $L += ('  ROOT pid={0} type={1} ppid={2} parent={3}' -f $p.ProcessId, `
            $(if ($p.CommandLine -match '--type=([a-z\-]+)') { $Matches[1] } else { 'BROWSER' }), $ppid, $pName)
        $cl = $p.CommandLine
        if ($null -eq $cl) { $cl = '<null>' }
        if ($cl.Length -gt 900) { $cl = $cl.Substring(0, 900) + ' ...[truncated]' }
        $L += ('       cmd = ' + $cl)
    }
}
$L += ''

$L += '== E. does any wisp.exe / balldebug.exe still exist on the machine? =='
$hosts = $all | Where-Object { $_.Name -in @('wisp.exe', 'balldebug.exe') }
if (@($hosts).Count -eq 0) { $L += '  none (count = 0)' }
foreach ($h in $hosts) {
    $L += ('  pid={0} name={1} started={2}' -f $h.ProcessId, $h.Name, ([datetime]$h.CreationDate).ToString('MM-dd HH:mm:ss'))
}
$L += ''

$L += '== F. parent PIDs of the webview roots, looked up by name even if pid reuse suspected =='
$rootPpids = @()
foreach ($p in $wv) {
    $parent = $byPid[[int]$p.ParentProcessId]
    if ($null -eq $parent -or $parent.Name -ne 'msedgewebview2.exe') { $rootPpids += [int]$p.ParentProcessId }
}
$L += '  distinct root ppids = ' + (($rootPpids | Sort-Object -Unique) -join ', ')
foreach ($rp in ($rootPpids | Sort-Object -Unique)) {
    $hit = $byPid[$rp]
    if ($null -eq $hit) { $L += ('  ppid {0} -> NOT PRESENT NOW' -f $rp) }
    else {
        $c = $hit.CommandLine
        if ($null -eq $c) { $c = '<null>' }
        if ($c.Length -gt 500) { $c = $c.Substring(0, 500) + ' ...[truncated]' }
        $L += ('  ppid {0} -> name={1} started={2} cmd={3}' -f $rp, $hit.Name, ([datetime]$hit.CreationDate).ToString('MM-dd HH:mm:ss'), $c)
    }
}
$L += ''
$L += 'end_of_forensics'

$L | Out-File -FilePath $out -Encoding utf8
Write-Output ('wrote ' + $out)
Write-Output ('count=' + (@($wv).Count))
