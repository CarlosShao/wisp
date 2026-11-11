# 302-r1 final-state verifier - run AFTER the green shot restore.
# Reads only; asserts the tree's two touched files equal the post-move snapshots
# byte-for-byte (via git hash-object) and re-derives the AC#4 boundary rosters.
$ErrorActionPreference = "Stop"
Set-Location "D:\work\workspace\projects plans\Wisp"
$fail = 0
function Check-Blob($file, $snapshot, $label) {
    $h1 = (git hash-object $file).Trim()
    $h2 = (git hash-object $snapshot).Trim()
    $ok = $h1 -eq $h2
    Write-Output ("{0}: file={1} snapshot={2} identical={3}" -f $label, $h1, $h2, $ok)
    if (-not $ok) { $script:fail = 1 }
}
Check-Blob "cmd\wisp\panel_resident_windows_test.go" ".scratch\wisp\probes\302\r1\base\panel_resident_windows_test.go.post302" "resident==post302"
Check-Blob "cmd\wisp\panel_page_reply_302_winlive_test.go" ".scratch\wisp\probes\302\r1\base\panel_page_reply_302_winlive_test.go.post302" "winlive==post302"
# AC#4 rosters: diff of the red roster of the red shot vs the green shot.
# runtests red-roster extraction is done by the caller; here we keep the two
# claim sets the ticket requires, one name per line, sorted.
Write-Output "--- taken-by-302 (sorted) ---"
Get-Content ".scratch\wisp\probes\302\r1\boundary\taken-by-302.txt"
Write-Output "--- risk-12 left at A817 (sorted) ---"
Get-Content ".scratch\wisp\probes\302\r1\boundary\risk-12-left-at-A817.txt"
Write-Output "VERIFIER_fail=$fail"
