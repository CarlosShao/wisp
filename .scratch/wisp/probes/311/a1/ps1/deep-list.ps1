# Probe 311-a1: walk a root level by level (no -Recurse, which silently drops
# children of directories this token cannot enter) and name every directory.
# Usage: powershell -NoProfile -File deep-list.ps1 <root> [<root> ...]
$ErrorActionPreference = 'Continue'
Write-Output ('now=' + (Get-Date -Format 'yyyy-MM-dd HH:mm:ss'))
foreach ($r in $args) {
    Write-Output ('===== ROOT=' + $r)
    if (-not (Test-Path -LiteralPath $r)) { Write-Output 'root_exists=NO'; continue }
    $queue = New-Object System.Collections.Queue
    $queue.Enqueue($r)
    $count = 0
    $errcount = 0
    while ($queue.Count -gt 0) {
        $cur = $queue.Dequeue()
        $kids = @()
        try { $kids = @(Get-ChildItem -LiteralPath $cur -Directory -Force -ErrorAction Stop) }
        catch {
            $errcount = $errcount + 1
            Write-Output ('CANNOT_ENUMERATE=' + $cur + ' err=' + $_.Exception.GetType().Name + ': ' + $_.Exception.Message)
            continue
        }
        foreach ($k in $kids) {
            $count = $count + 1
            Write-Output ('DIR|creation=' + $k.CreationTime.ToString('yyyy-MM-dd HH:mm:ss.fff') + '|' + $k.FullName)
            $queue.Enqueue($k.FullName)
        }
    }
    Write-Output ('dir_count_excluding_root=' + $count)
    Write-Output ('enumerate_denied_count=' + $errcount)
}
Write-Output ('ps_errors=' + ($Error.Count))
