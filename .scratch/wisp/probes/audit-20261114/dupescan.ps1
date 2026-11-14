# dupescan.ps1 — 数「同一个 commit 的 tree 里，同一票号的文件有多少枚」
#
# 用途：HEAD 双计数（同一张票的旧名 + 新名同时在树里）不是看磁盘，是看**提交树的 tree**。
#       2026-11-14 编排者结 441／274 时各犯过一次（见 docs/evidence/s1/audit-20261114-closures.md §2.1）。
#
# 用法（在仓库根跑）：
#   pwsh -NoProfile -File .scratch/wisp/probes/audit-20261114/dupescan.ps1
#   pwsh -NoProfile -File .scratch/wisp/probes/audit-20261114/dupescan.ps1 -Tickets 441,442,274,443,444 -Since 'b062d52a7'
#
# 判据：对 `$Since^..HEAD` 每一笔 commit，取 `git ls-tree -r --name-only <c> -- .scratch/wisp/issues`，
#       按 "/<票号>-" 数；count > 1 即双计数（打印该 commit 的短 hash、时间、哪几张票是几枚）。
#       -Since 的默认值 `b062d52a7` ＝ 票 441 立案那一笔，往前没有这五张票。

param(
    [int[]]$Tickets = @(441, 442, 274, 443, 444),
    [string]$Since = 'b062d52a7'
)

$ErrorActionPreference = 'Stop'
Push-Location (git rev-parse --show-toplevel)
try {
    $commits = git rev-list --reverse "$Since^..HEAD"
    $hot = 0
    foreach ($c in $commits) {
        $names = git ls-tree -r --name-only $c -- .scratch/wisp/issues
        $counts = @{}
        foreach ($t in $Tickets) {
            $n = ($names | Where-Object { $_ -match "/$t-" }).Count
            if ($n -gt 0) { $counts[$t] = $n }
        }
        $dup = $counts.GetEnumerator() | Where-Object { $_.Value -gt 1 } | ForEach-Object { "$($_.Key)=$($_.Value)" }
        if ($dup) {
            $hot++
            $stamp = git show -s --format='%h %ci' $c
            "$stamp  DOUBLE: $($dup -join ',')"
        }
    }
    if ($hot -eq 0) { "no commits with duplicate ticket paths under .scratch/wisp/issues (checked $($commits.Count) commits)" }

    # 现 HEAD 的对照：每票必须 =1
    ""
    "== HEAD ($(git rev-parse --short HEAD)) =="
    $headNames = git ls-tree -r --name-only HEAD -- .scratch/wisp/issues
    foreach ($t in $Tickets) {
        "{0} => {1}" -f $t, (($headNames | Where-Object { $_ -match "/$t-" }).Count)
    }
}
finally {
    Pop-Location
}
