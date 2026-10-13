$ErrorActionPreference = 'Continue'
function SizeMB($p) {
    try {
        $s = (Get-ChildItem -LiteralPath $p -Recurse -Force -File -ErrorAction SilentlyContinue | Measure-Object -Property Length -Sum).Sum
        if ($null -eq $s) { return 0 } else { return [math]::Round($s / 1MB, 1) }
    } catch { return -1 }
}
"=== 候选（用户目录下 wisp-*/wisp*）逐枚：大小、子项数、最深层有无不可读件 ==="
$cands = @(
    'C:\Users\swq\wisp-301-r2-tmp',
    'C:\Users\swq\wisp-303-bisect',
    'C:\Users\swq\wisp-303-steps',
    'C:\Users\swq\wisp-303-v1',
    'C:\Users\swq\wisp-303-v1-blob',
    'C:\Users\swq\wisp-303-v1-mut-base',
    'C:\Users\swq\wisp-303-v1-mut-M1',
    'C:\Users\swq\wisp-303-v1-mut-M2',
    'C:\Users\swq\wisp-303-v1-mut-M3',
    'C:\Users\swq\wisp-303-v1-mut-M4',
    'C:\Users\swq\wisp35v7mut',
    'C:\Users\swq\wisp84-snap'
)
foreach ($c in $cands) {
    if (Test-Path -LiteralPath $c) {
        $files = Get-ChildItem -LiteralPath $c -Recurse -Force -File -ErrorAction SilentlyContinue
        $denied = @(Get-ChildItem -LiteralPath $c -Recurse -Force -Directory -ErrorAction SilentlyContinue | Where-Object { -not (Get-ChildItem -LiteralPath $_.FullName -Force -ErrorAction SilentlyContinue) -and (Test-Path -LiteralPath $_.FullName) })
        "{0} | sizeMB={1} | files={2} | subdirs={3}" -f $c, (SizeMB $c), $files.Count, (Get-ChildItem -LiteralPath $c -Recurse -Force -Directory -ErrorAction SilentlyContinue).Count
        $files | Sort-Object Length -Descending | Select-Object -First 3 | ForEach-Object { "     {0} | {1} KB" -f $_.Name, [math]::Round($_.Length / 1KB, 1) }
    } else { "$c | ABSENT" }
}
"`n=== 专用根 82 枚里：有内容的（>0.05MB）逐枚列 ==="
$g = 'C:\Users\swq\AppData\Local\wisp-test-gotmp'
Get-ChildItem -LiteralPath $g -Directory -Force -ErrorAction SilentlyContinue | ForEach-Object {
    $s = SizeMB $_.FullName
    if ($s -gt 0.05) { "{0} | sizeMB={1} | mtime={2}" -f $_.Name, $s, $_.LastWriteTime.ToString('MM-dd_HH:mm') }
} 
"`n=== 专用根里 0MB 的抽三枚，看里面到底有没有东西（含隐藏/不可读）==="
$zero = Get-ChildItem -LiteralPath $g -Directory -Force -ErrorAction SilentlyContinue | Where-Object { (SizeMB $_.FullName) -le 0.05 } | Select-Object -First 3
foreach ($z in $zero) {
    "--- $($z.Name) ---"
    Get-ChildItem -LiteralPath $z.FullName -Force -Recurse -ErrorAction SilentlyContinue | Select-Object -First 8 | ForEach-Object { "   {0} | dir={1} | KB={2}" -f $_.FullName.Substring($z.FullName.Length), $_.PSIsContainer, [math]::Round($_.Length/1KB,1) }
}
"`nrc=0"
