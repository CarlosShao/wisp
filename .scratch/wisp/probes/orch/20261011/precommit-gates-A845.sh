#!/bin/sh
# 编排者 A845 批门禁（纯账目＋两枚我自己名下的旧证据件入库；⛔ 碰 .go；⛔ 跑 go build/vet/test——298-v1 独占编译面）
cd "$(git rev-parse --show-toplevel)" || exit 9
D=.scratch/wisp/probes/orch/20261011
LOG=$D/logs/orch-A845-gates.txt
EV1=.scratch/wisp/probes/orch/296-ac3-runsheet.md
EV2=.scratch/wisp/probes/300/orch/logs/r3-orch-v4-attacks-rcfixed-20261010-204046.txt

sh scripts/d22scan.sh > $D/logs/d22scan-A845.txt 2>&1; rc1=$?
sh scripts/check-path-length-budget.sh --with-self-test > $D/logs/lenbudget-A845.txt 2>&1; rc2=$?

# ★⛔ 把在飞腿 writeface-1 的目录扫进我这批（它此刻正在写；它的件归它自己的 commit）
git add -- docs/reports/pending-and-issues.md docs/reports/HANDOVER.md $D ":(exclude)$D/writeface-1" $EV1 $EV2 2>/dev/null
git diff --cached --name-only > $D/logs/staged-A845.txt
# 只算"我这一批"：writeface-1 已被 staged 的条目＝在飞腿的账，⛔ 算进我这批的台账（本轮同族＝按名字认人）
grep "^$D/" $D/logs/staged-A845.txt | grep -v "^$D/writeface-1/" > /tmp/own-batch.txt
n_all=$(grep -c '' /tmp/own-batch.txt)
wf=$(git diff --cached --name-only | grep -c "^$D/writeface-1/")
n_go=$(grep -c '\.go$' /tmp/own-batch.txt)
n_out=$(grep -c '\.out$' /tmp/own-batch.txt)
foreign=$(grep -v '^\.scratch/wisp/probes/orch/\|^\.scratch/wisp/probes/300/orch/\|^docs/reports/' /tmp/own-batch.txt | grep -c '')
dirty=$(git status --porcelain -- internal cmd | wc -l | tr -d ' ')
{
  echo "d22scan_rc=$rc1"
  echo "lenbudget_rc=$rc2"
  echo "staged_total=$n_all staged_go=$n_go staged_dot_out=$n_out foreign_staged=$foreign"
  echo "inflight_writeface1_also_staged=$wf  # ⛔ 我这批的账；commit pathspec 用 :(exclude) 剥掉"
  echo "worktree_dirty_internal_cmd=$dirty"
  echo "numstat_ledger=$(git diff --cached --numstat -- docs/reports/pending-and-issues.md)"
  echo "numstat_handover=$(git diff --cached --numstat -- docs/reports/HANDOVER.md)"
  echo "prefix_ledger_vs_head=$(git show HEAD:docs/reports/pending-and-issues.md | cmp -s - <(head -$(git show HEAD:docs/reports/pending-and-issues.md | wc -l) docs/reports/pending-and-issues.md) && echo IDENTICAL || echo DIFF)"
  echo "prefix_handover_vs_head=$(git show HEAD:docs/reports/HANDOVER.md | cmp -s - <(head -$(git show HEAD:docs/reports/HANDOVER.md | wc -l) docs/reports/HANDOVER.md) && echo IDENTICAL || echo DIFF)"
  echo "lines_ledger=$(grep -c '' docs/reports/pending-and-issues.md) head_lines_ledger=$(git show HEAD:docs/reports/pending-and-issues.md | wc -l)"
  echo "lines_handover=$(grep -c '' docs/reports/HANDOVER.md) head_lines_handover=$(git show HEAD:docs/reports/HANDOVER.md | wc -l)"
  # ERE ruler (the BRE '?' bug is recorded in logs/orch-A843-gates-b.txt; keep the ruler form honest)
  echo "handover_4_0c_headers_ERE=$(grep -cE '^### 4\.0c[a-z]{1,2} 停车点' docs/reports/HANDOVER.md)"
  echo "handover_last_header=$(grep -E '^### 4\.0c[a-z]{1,2} 停车点' docs/reports/HANDOVER.md | tail -1 | cut -c1-24)"
  echo "ledger_A84x_headers=$(grep -cE '^## A84[0-9]' docs/reports/pending-and-issues.md)"
  echo "ledger_last=$(grep -E '^## A8[0-9]+｜' docs/reports/pending-and-issues.md | tail -1 | cut -c1-16)"
  echo "ev1_bytes=$(wc -c < $EV1) ev2_bytes=$(wc -c < $EV2)"
  echo "staged_outside_own_dirs=$(git diff --cached --name-only | grep -v '^\.scratch/wisp/probes/orch/\|^\.scratch/wisp/probes/300/orch/\|^docs/reports/' | grep -c '')"
} > $LOG 2>&1
cat $LOG
