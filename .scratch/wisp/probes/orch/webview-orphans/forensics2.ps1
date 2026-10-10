# Read-only follow-up: does ANY process on this machine belong to Wisp?
# Ruler 2 = full command-line scan for a webview hosted by wisp/balldebug.
# Ruler 3 = leftover test binaries (*test*.exe / *.test.exe).
$ErrorActionPreference = 'Stop'
$out = Join-Path $PSScriptRoot 'forensics-2.txt'
$L = @()
$L += 'ruler = Get-CimInstance Win32_Process, full CommandLine scan (read-only)'
$L += 'captured_local = ' + (Get-Date -Format 'yyyy-MM-dd HH:mm:ss zzz')
$L += ''

$all = @(Get-CimInstance Win32_Process)
$L += 'total processes on machine = ' + $all.Count
$L += ''

$L += '== R2. any process whose command line names Wisp as a WebView2 host? =='
$mine = $all | Where-Object {
    ($null -ne $_.CommandLine) -and (
        $_.CommandLine -match 'webview-exe-name=wisp' -or
        $_.CommandLine -match 'webview-exe-name=balldebug' -or
        $_.CommandLine -match 'projects plans' -or
        $_.CommandLine -match 'wisp\.exe' -or
        $_.CommandLine -match 'balldebug\.exe')
}
if (@($mine).Count -eq 0) { $L += '  HITS = 0' }
foreach ($p in $mine) {
    $c = $p.CommandLine
    if ($c.Length -gt 300) { $c = $c.Substring(0, 300) + ' ...' }
    $L += ('  HIT pid={0} name={1} cmd={2}' -f $p.ProcessId, $p.Name, $c)
}
$L += ''

$L += '== R2b. distinct --webview-exe-name= values across the whole machine =='
$names = @()
foreach ($p in $all) {
    if ($null -ne $p.CommandLine -and $p.CommandLine -match '--webview-exe-name=([^\s]+)') {
        $names += $Matches[1]
    }
}
$names | Group-Object | Sort-Object Count -Descending | ForEach-Object {
    $L += ('  {0,-28} count={1}' -f $_.Name, $_.Count)
}
$L += ''

$L += '== R3. leftover test binaries =='
$tb = $all | Where-Object { $_.Name -like '*test*.exe' -or $_.Name -like '*.test.exe' }
if (@($tb).Count -eq 0) { $L += '  none (count = 0)' }
foreach ($p in $tb) { $L += ('  pid={0} name={1}' -f $p.ProcessId, $p.Name) }
$L += ''

$L += '== R3b. go / go test / build.exe currently running? =='
$go = $all | Where-Object { $_.Name -in @('go.exe', 'go_gofmt.exe', 'wisp.exe', 'balldebug.exe') -or $_.Name -like 'go*' }
foreach ($p in $go) { $L += ('  pid={0} name={1}' -f $p.ProcessId, $p.Name) }
if (@($go).Count -eq 0) { $L += '  none (count = 0)' }
$L += ''
$L += 'end_of_forensics_2'
$L | Out-File -FilePath $out -Encoding utf8
Write-Output ('wrote ' + $out)
