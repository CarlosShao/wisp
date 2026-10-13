$ErrorActionPreference = 'Continue'
function SizeMB($p) {
    try {
        $s = (Get-ChildItem -LiteralPath $p -Recurse -Force -File -ErrorAction SilentlyContinue | Measure-Object -Property Length -Sum).Sum
        if ($null -eq $s) { return 0 } else { return [math]::Round($s / 1MB, 1) }
    } catch { return -1 }
}
function CountDirs($p) {
    (Get-ChildItem -LiteralPath $p -Directory -Force -ErrorAction SilentlyContinue | Measure-Object).Count
}

"=== A. 应用自己的数据目录（不是测试残留）==="
foreach ($t in @(
        'C:\Users\swq\AppData\Local\wisp-test-gotmp',
        'C:\Users\swq\AppData\Local\wisp',
        'C:\Users\swq\AppData\Roaming\wisp',
        'C:\Users\swq\AppData\Roaming\wisp-dev',
        'C:\Users\swq\AppData\Local\Temp\wisp-test-gotmp')) {
    if (Test-Path -LiteralPath $t) {
        "{0} | sizeMB={1} | child_dirs={2}" -f $t, (SizeMB $t), (CountDirs $t)
    } else { "$t | ABSENT" }
}

"`n=== B. wisp-test-gotmp 专用根：每一枚运行根的名字/时间/大小（前 100）==="
$g = 'C:\Users\swq\AppData\Local\wisp-test-gotmp'
if (Test-Path -LiteralPath $g) {
    Get-ChildItem -LiteralPath $g -Directory -Force -ErrorAction SilentlyContinue |
        Sort-Object LastWriteTime |
        Select-Object -First 100 |
        ForEach-Object { "{0} | mtime={1} | sizeMB={2}" -f $_.Name, $_.LastWriteTime.ToString('MM-dd_HH:mm'), (SizeMB $_.FullName) }
    "total_roots=" + (CountDirs $g)
}

"`n=== C. %TEMP% 顶层：名字里带 wisp / Test / go-build 的 ==="
Get-ChildItem -LiteralPath $env:TEMP -Force -ErrorAction SilentlyContinue |
    Where-Object { $_.Name -like '*wisp*' -or $_.Name -like 'Test*' -or $_.Name -like 'go-build*' } |
    Sort-Object LastWriteTime |
    ForEach-Object { "{0} | type={1} | mtime={2} | sizeMB={3}" -f $_.Name, $(if ($_.PSIsContainer) { 'dir' } else { 'file' }), $_.LastWriteTime.ToString('MM-dd_HH:mm'), (SizeMB $_.FullName) }

"`n=== D. 用户目录下名字带 wisp 的目录/文件（深度 2）==="
Get-ChildItem -LiteralPath 'C:\Users\swq' -Force -Recurse -Depth 2 -ErrorAction SilentlyContinue |
    Where-Object { $_.Name -like '*wisp*' -and $_.FullName -notlike '*AppData\Local\Temp*' } |
    ForEach-Object { "{0} | dir={1} | mtime={2}" -f $_.FullName, $_.PSIsContainer, $_.LastWriteTime.ToString('MM-dd_HH:mm') }

"`n=== E. 两枚应用数据目录里的文件清单（名字+大小，前 40 枚最大）==="
foreach ($d in @('C:\Users\swq\AppData\Roaming\wisp', 'C:\Users\swq\AppData\Roaming\wisp-dev')) {
    if (Test-Path -LiteralPath $d) {
        "--- $d ---"
        Get-ChildItem -LiteralPath $d -Recurse -Force -File -ErrorAction SilentlyContinue |
            Sort-Object Length -Descending | Select-Object -First 12 |
            ForEach-Object { "{0} | sizeMB={1} | mtime={2}" -f ($_.FullName.Substring($d.Length + 1)), [math]::Round($_.Length / 1MB, 2), $_.LastWriteTime.ToString('MM-dd_HH:mm') }
    }
}
"`nrc=0"
