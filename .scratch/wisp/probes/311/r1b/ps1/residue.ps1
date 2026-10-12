# residue.ps1 - 311-r1b: inventory of what this leg left on disk (delete-forbidden:
# ticket 310's iron rule, restated verbatim on ticket 311's forbidden list). Everything
# here is NAME + CreationTime only; nothing is removed, chmod'd or takeown'd.
$ErrorActionPreference = 'SilentlyContinue'

function Show-Tree($label, $path, $depth) {
    if (-not (Test-Path -LiteralPath $path)) { Write-Output ("$label MISSING path=$path"); return }
    $item = Get-Item -LiteralPath $path
    $kind = if ($item.PSIsContainer) { 'dir' } else { 'file' }
    Write-Output ("$label path=" + $item.FullName + " kind=" + $kind + " created=" + $item.CreationTime.ToString('yyyy-MM-dd HH:mm:ss'))
    if ($depth -le 0) { return }
    $kids = @(Get-ChildItem -LiteralPath $path -Directory -ErrorAction SilentlyContinue)
    Write-Output ("$label child_dir_count=" + $kids.Count)
    foreach ($k in $kids) {
        $gk = @(Get-ChildItem -LiteralPath $k.FullName -Directory -ErrorAction SilentlyContinue)
        Write-Output ("  root=" + $k.Name + " created=" + $k.CreationTime.ToString('yyyy-MM-dd HH:mm:ss') + " subdirs=" + $gk.Count)
        if ($depth -ge 2) {
            foreach ($g in $gk) {
                Write-Output ("    sub=" + $g.Name + " created=" + $g.CreationTime.ToString('yyyy-MM-dd HH:mm:ss'))
            }
        }
    }
}

Show-Tree 'xvol_base' 'D:\tmp\wisp-test-gotmp\311r1b-xvol' 2
Show-Tree 'samevol' (Join-Path $env:LOCALAPPDATA '311r1b-samevol') 2
Show-Tree 'samevol_pre' (Join-Path $env:LOCALAPPDATA '311r1b-samevol-pre') 1
Show-Tree 'default_base' (Join-Path $env:LOCALAPPDATA 'wisp-test-gotmp') 0
Write-Output ('planted_notadir=' + (Test-Path 'D:\tmp\311r1b-notadir') + ' planted_notadir_file=' + (Test-Path 'D:\tmp\311r1b-notadir-file'))
Write-Output ('export_trees_top=' + ((Get-ChildItem -LiteralPath 'D:\tmp\311r1b' -Directory -ErrorAction SilentlyContinue).Count))

# The locked leaves: named so a later leg's "entries in root" arithmetic counts them.
# Detection = Get-ChildItem with -ErrorAction Stop inside try; the catch is the finding.
$leaves = @()
$scanned = 0
foreach ($r in @(Get-ChildItem -LiteralPath 'D:\tmp\wisp-test-gotmp\311r1b-xvol' -Directory -ErrorAction SilentlyContinue)) {
    foreach ($t in @(Get-ChildItem -LiteralPath $r.FullName -Directory -ErrorAction SilentlyContinue)) {
        foreach ($i in @(Get-ChildItem -LiteralPath $t.FullName -Directory -ErrorAction SilentlyContinue)) {
            $scanned++
            $ok = $true
            $err = ''
            try { $null = @(Get-ChildItem -LiteralPath $i.FullName -ErrorAction Stop) } catch { $ok = $false; $err = $_.Exception.GetType().Name }
            if (-not $ok) { $leaves += ($r.Name + '\' + $t.Name + '\' + $i.Name + ' err=' + $err) }
        }
    }
}
Write-Output ('leaf_dirs_scanned=' + $scanned)
Write-Output ('unreadable_leaf_count=' + $leaves.Count)
foreach ($l in $leaves) { Write-Output ('  unreadable_leaf=' + $l) }
