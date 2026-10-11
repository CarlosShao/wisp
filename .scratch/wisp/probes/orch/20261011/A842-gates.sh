#!/bin/sh
# 编排者本批门禁（docs-only 批：只动 .md／.txt 证据件，⛔ 碰 .go）
# 每格自己落一行 rc=N（0 字节的件＝那格没交）
cd "$(git rev-parse --show-toplevel)" || exit 9

LOG=.scratch/wisp/probes/orch/20261011/logs/orch-A842-gates.txt
: > "$LOG"

# 1) d22scan（全仓形状门）
sh scripts/d22scan.sh > .scratch/wisp/probes/orch/20261011/logs/d22scan-A842.txt 2>&1
echo "d22scan_rc=$?" >> "$LOG"

# 2) 长度门（票 309 新文件名要过预算）
sh scripts/check-path-length-budget.sh --with-self-test > .scratch/wisp/probes/orch/20261011/logs/lenbudget-A842.txt 2>&1
echo "lenbudget_rc=$?" >> "$LOG"

# 3) 本批⛔ 含 .go（所以 gofmt／build／vet 的分母与颜色不由我改动）
git diff --cached --name-only > .scratch/wisp/probes/orch/20261011/logs/staged-A842.txt 2>&1
n_go=$(grep -c '\.go$' .scratch/wisp/probes/orch/20261011/logs/staged-A842.txt)
echo "staged_go=$n_go" >> "$LOG"

# 4) 暂存枚数与越界面
n_all=$(wc -l < .scratch/wisp/probes/orch/20261011/logs/staged-A842.txt)
echo "staged_total=$n_all" >> "$LOG"

# 5) 台账／停车点 append-only 尺
echo "ledger_numstat=$(git diff --numstat -- docs/reports/pending-and-issues.md | awk '{print $1" 加 "$2" 删"}')" >> "$LOG"
echo "handover_numstat=$(git diff --numstat -- docs/reports/HANDOVER.md | awk '{print $1" 加 "$2" 删"}')" >> "$LOG"

cat "$LOG"
