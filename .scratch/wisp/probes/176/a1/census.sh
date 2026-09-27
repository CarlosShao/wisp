#!/bin/sh
# 票 176 AC#1 —— "起跑口不存在"的可复算读数生成器（只读，零产码副作用）。
# 派单：.scratch/wisp/dispatches/2026-09-27-235x-readonly-176-a1-*.md
# 程：176-a1｜锚点：fbcdf441764a77ccb917348b057fc5fb5907bebf（本程序不写死锚点，锚点由 -- 前 HEAD 现量）
#
# ⚠ logdir 取脚本自身目录（派单 §5 明令）。
logdir=$(cd "$(dirname "$0")" && pwd)
# 仓根＝脚本目录上溯五级：.scratch/wisp/probes/176/a1 -> repo root
repo=$(cd "$logdir/../../../../.." && pwd)
cd "$repo" || { echo "census: cannot reach repo root from $logdir"; exit 9; }

# readings.txt 是本程第一次跑坏的记录（仓根算错成 .scratch/wisp，grep 全数报"目录不存在"）
# ⚠ 零删除 ⇒ 不覆盖它，另开 readings2.txt 作有效读数。
out="$logdir/readings2.txt"
: > "$out"

say() { printf '%s\n' "$*" >> "$out"; }

say "# 176-a1 起跑口普查读数"
say "# generated=$(date -Is)  HEAD=$(git rev-parse HEAD)  branch=$(git rev-parse --abbrev-ref HEAD)"
say ""

# --- 尺 1：RunAsync 到底存不存在（定义行） -----------------------------------
say "## R1 RunAsync 定义行（存在性；不存在则本票票面措辞要改写）"
say "\$ grep -rn 'func (l \*Loop) RunAsync' --include='*.go' internal/"
grep -rn 'func (l \*Loop) RunAsync' --include='*.go' internal/ >> "$out" 2>&1
say "count=$(grep -rn 'func (l \*Loop) RunAsync' --include='*.go' internal/ | wc -l)"
say ""

# --- 尺 2：RunAsync 生产调用点（排测试） -------------------------------------
say "## R2 RunAsync 生产调用点（点号调用形，排 _test.go）"
say "\$ grep -rn '\.RunAsync(' --include='*.go' internal/ cmd/ | grep -v _test.go"
grep -rn '\.RunAsync(' --include='*.go' internal/ cmd/ | grep -v _test.go >> "$out" 2>&1
say "count=$(grep -rn '\.RunAsync(' --include='*.go' internal/ cmd/ | grep -v _test.go | wc -l)"
say ""

say "## R2b 同尺含测试（对照：命中全在测试里）"
say "\$ grep -rn '\.RunAsync(' --include='*.go' internal/ cmd/ | wc -l"
say "count=$(grep -rn '\.RunAsync(' --include='*.go' internal/ cmd/ | wc -l)"
grep -rn '\.RunAsync(' --include='*.go' internal/ cmd/ >> "$out" 2>&1
say ""

# --- 尺 3：TaskRoster.Record 生产写者（排测试） ------------------------------
say "## R3 TaskRoster.Record 生产写者（派单给的尺形原文）"
say "\$ grep -rn '\.Record(' --include='*.go' internal/ cmd/ | grep -v _test.go"
grep -rn '\.Record(' --include='*.go' internal/ cmd/ | grep -v _test.go >> "$out" 2>&1
say "count=$(grep -rn '\.Record(' --include='*.go' internal/ cmd/ | grep -v _test.go | wc -l)"
say ""

say "## R3b 同尺含测试（对照）"
say "\$ grep -rn '\.Record(' --include='*.go' internal/ cmd/ | wc -l"
say "count=$(grep -rn '\.Record(' --include='*.go' internal/ cmd/ | wc -l)"
say ""

say "## R3c Record 定义行（存在性）"
say "\$ grep -rn 'func (r \*TaskRoster) Record' --include='*.go' internal/"
grep -rn 'func (r \*TaskRoster) Record' --include='*.go' internal/ >> "$out" 2>&1
say ""

# --- 尺 4：裸 goroutine 与受管 spawn 的先例 ----------------------------------
say "## R4 裸 go func( 生产命中（ban 1 口径，排测试）"
say "\$ grep -rn 'go func(' --include='*.go' internal/ cmd/ | grep -v _test.go"
grep -rn 'go func(' --include='*.go' internal/ cmd/ | grep -v _test.go >> "$out" 2>&1
say "count=$(grep -rn 'go func(' --include='*.go' internal/ cmd/ | grep -v _test.go | wc -l)"
say ""

say "## R5 observe.Registry.Spawn 生产调用点（受管先例）"
say "\$ grep -rn '\.Spawn(' --include='*.go' internal/ cmd/ | grep -v _test.go"
grep -rn '\.Spawn(' --include='*.go' internal/ cmd/ | grep -v _test.go >> "$out" 2>&1
say "count=$(grep -rn '\.Spawn(' --include='*.go' internal/ cmd/ | grep -v _test.go | wc -l)"
say ""

# --- 尺 6：形状禁区 (i) 的现量（G3 pattern 自己跑一遍） ---------------------
say "## R6 禁区 (i) 自证：Loop 上收 taskID 的导出方法（= gate-clauses G3 的 pattern）"
say "\$ grep -rnE '^func \\(l \\*Loop\\) [A-Z][A-Za-z0-9]*\\(.*taskID' internal/agent"
grep -rnE '^func \(l \*Loop\) [A-Z][A-Za-z0-9]*\(.*taskID' internal/agent >> "$out" 2>&1
say "count=$(grep -rnE '^func \(l \*Loop\) [A-Z][A-Za-z0-9]*\(.*taskID' internal/agent | wc -l)"
say ""

# --- 尺 7：起跑口落地后 C25 豁免载具是否可达（票 175 惰性那一条） ------------
say "## R7 Record 的 ArtifactPath 写入点（豁免载具的唯一上游；排测试）"
say "\$ grep -rn 'ArtifactPath:' --include='*.go' internal/ cmd/ | grep -v _test.go"
grep -rn 'ArtifactPath:' --include='*.go' internal/ cmd/ | grep -v _test.go >> "$out" 2>&1
say "count=$(grep -rn 'ArtifactPath:' --include='*.go' internal/ cmd/ | grep -v _test.go | wc -l)"
say ""

say "# 读数完毕。判读见 docs/evidence/s1/176-background-start-port-census-a1.md"

printf '%s\n' "$out"
