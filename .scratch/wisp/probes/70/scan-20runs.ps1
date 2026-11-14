$ErrorActionPreference = 'Continue'
$out = "D:\work\workspace\projects plans\Wisp\.scratch\wisp\probes\70\history"
New-Item -ItemType Directory -Force -Path $out | Out-Null
$ids = (Get-Content "D:\work\workspace\projects plans\Wisp\.scratch\wisp\probes\70\run-list-20.json" -Raw | ConvertFrom-Json)
$summary = @()
$n = 0
foreach ($r in $ids) {
    $n++
    $id = $r.databaseId
    $f = Join-Path $out "$id.log"
    if (-not (Test-Path $f)) {
        gh run view $id --log-failed 2>&1 | Set-Content -Path $f -Encoding utf8
    }
    $lines = Get-Content $f -ErrorAction SilentlyContinue
    $size = (Get-Item $f).Length
    $t221fail = ($lines | Select-String -Pattern '^--- FAIL: Test221SubagentCannotStopSiblingOrItself' -SimpleMatch:$false).Count
    $t221any = ($lines | Select-String -Pattern 'Test221SubagentCannotStopSiblingOrItself').Count
    $leg452 = ($lines | Select-String -Pattern 'task_cancel_221_legs_test\.go:45').Count
    $four = ($lines | Select-String -Pattern 'four numbers \(all from -v output\)').Count
    $fourline = ($lines | Select-String -Pattern 'portable-tests\.sh: four numbers \(all from -v output\)' | ForEach-Object { $_.Line }) -join ' || '
    $summary += [pscustomobject]@{ n=$n; id=$id; sha=$r.headSha.Substring(0,9); conclusion=$r.conclusion; bytes=$size; t221_fail=$t221fail; t221_any=$t221any; leg452=$leg452; fournum=$four; fourline=$fourline }
    "{0}/{1} {2} {3} bytes={4} t221fail={5} any={6} leg45x={7} four={8}" -f $n, $ids.Count, $id, $r.conclusion, $size, $t221fail, $t221any, $leg452, $four
}
$summary | ConvertTo-Json -Depth 4 | Set-Content -Path (Join-Path $out "summary.json") -Encoding utf8
