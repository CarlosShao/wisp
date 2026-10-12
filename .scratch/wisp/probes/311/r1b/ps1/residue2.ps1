# residue2.ps1 - 311-r1b: locked-leaf scan, DEEPENED.
# The first version (ps1/residue.ps1) probed three levels below the run root
# (root/testName/NNN) and reported unreadable_leaf_count=0 - a false negative,
# because the leaf Go's cleanup cannot delete is ONE level deeper
# (root/TestSealNarrows.../002/la). This version recurses to depth 6 and records the
# directory that actually throws, with the exception type. Positive control is in
# the same trip: the roots other legs left on D: (ticket 311's registered locked
# leaves) must come back non-zero, or this ruler is not proven to be able to see one.
$ErrorActionPreference = 'SilentlyContinue'
$script:leaves = @()
$script:scanned = 0

function Walk($path, $depth, $rootName) {
    if ($depth -le 0) { return }
    $kids = $null
    $err = $null
    try {
        $kids = @(Get-ChildItem -LiteralPath $path -Directory -ErrorAction Stop)
    } catch {
        $script:leaves += ($rootName + ' :: ' + $path.Replace('D:\tmp\wisp-test-gotmp\', '') + ' :: ' + $_.Exception.GetType().Name)
        return
    }
    $script:scanned++
    foreach ($k in $kids) { Walk $k.FullName ($depth - 1) $rootName }
}

foreach ($base in @('D:\tmp\wisp-test-gotmp', (Join-Path $env:LOCALAPPDATA 'wisp-test-gotmp'))) {
    if (-not (Test-Path -LiteralPath $base)) { Write-Output ("base MISSING " + $base); continue }
    Write-Output ("BASE " + $base)
    $roots = @(Get-ChildItem -LiteralPath $base -Directory -ErrorAction SilentlyContinue)
    Write-Output ("  top_root_dir_count=" + $roots.Count)
    foreach ($r in $roots) {
        $before = $script:leaves.Count
        Walk $r.FullName 6 $r.Name
        $mine = $script:leaves.Count - $before
        if ($mine -gt 0) {
            Write-Output ("  root=" + $r.Name + " created=" + $r.CreationTime.ToString('yyyy-MM-dd HH:mm:ss') + " unreadable_leaves=" + $mine)
        }
    }
}
Write-Output ('dirs_listed_ok=' + $script:scanned)
Write-Output ('unreadable_leaf_total=' + $script:leaves.Count)
foreach ($l in $script:leaves) { Write-Output ('LEAF ' + $l) }
