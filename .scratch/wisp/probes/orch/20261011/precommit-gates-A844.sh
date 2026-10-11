#!/bin/sh
# 编排者 A844 批门禁（纯账目＋一枚翻框：台账 A844＋停车点 §4.0ci＋票 309 追加节＋编排者件；⛔ 碰 .go）
cd "$(git rev-parse --show-toplevel)" || exit 9
D=.scratch/wisp/probes/orch/20261011
LOG=$D/logs/orch-A844-gates.txt
T309=.scratch/wisp/issues/309-assessorfor-empty-taskid-limb-silently-drops-r4-detector.md

sh scripts/d22scan.sh > $D/logs/d22scan-A844.txt 2>&1; rc1=$?
sh scripts/check-path-length-budget.sh --with-self-test > $D/logs/lenbudget-A844.txt 2>&1; rc2=$?

git add -- docs/reports/pending-and-issues.md docs/reports/HANDOVER.md $T309 $D 2>/dev/null
git diff --cached --name-only > $D/logs/staged-A844.txt
n_all=$(grep -c '' $D/logs/staged-A844.txt)
n_go=$(grep -c '\.go$' $D/logs/staged-A844.txt)
n_out=$(grep -c '\.out$' $D/logs/staged-A844.txt)
foreign=$(git diff --cached --name-only | grep -v "^$D/\|^docs/reports/\|^$T309" | grep -c '')
dirty=$(git status --porcelain -- internal cmd | wc -l | tr -d ' ')
PAT_HDR='^### 4\.0c[a-z]\{0,1\}[a-z]\{0,1\} 停车点'

{
  echo "d22scan_rc=$rc1"
  echo "lenbudget_rc=$rc2"
  echo "staged_total=$n_all staged_go=$n_go staged_dot_out=$n_out foreign_staged=$foreign"
  echo "worktree_dirty_internal_cmd=$dirty"
  echo "numstat_ledger=$(git diff --cached --numstat -- docs/reports/pending-and-issues.md)"
  echo "numstat_handover=$(git diff --cached --numstat -- docs/reports/HANDOVER.md)"
  echo "numstat_ticket309=$(git diff --cached --numstat -- $T309)"
  echo "prefix_ledger_vs_head=$(git show HEAD:docs/reports/pending-and-issues.md | cmp -s - <(head -$(git show HEAD:docs/reports/pending-and-issues.md | wc -l) docs/reports/pending-and-issues.md) && echo IDENTICAL || echo DIFF)"
  echo "prefix_handover_vs_head=$(git show HEAD:docs/reports/HANDOVER.md | cmp -s - <(head -$(git show HEAD:docs/reports/HANDOVER.md | wc -l) docs/reports/HANDOVER.md) && echo IDENTICAL || echo DIFF)"
  echo "lines_ledger=$(grep -c '' docs/reports/pending-and-issues.md) head_lines_ledger=$(git show HEAD:docs/reports/pending-and-issues.md | wc -l)"
  echo "lines_handover=$(grep -c '' docs/reports/HANDOVER.md) head_lines_handover=$(git show HEAD:docs/reports/HANDOVER.md | wc -l)"
  echo "handover_4_0c_headers=$(grep -c "$PAT_HDR" docs/reports/HANDOVER.md)  # BRE-safe ruler (the '?' bug is recorded in logs/orch-A843-gates-b.txt)"
  echo "ledger_A84x_headers=$(grep -cE '^## A84[0-9]' docs/reports/pending-and-issues.md)"
  echo "ticket309_checked=$(grep -cE '^[[:space:]]*- \[x\]' $T309) ticket309_unchecked=$(grep -cE '^[[:space:]]*- \[ \]' $T309)"
  echo "ticket309_name_len=$(basename $T309 | wc -c | tr -d ' ')  # wc -c counts the newline too => subtract 1"
} > $LOG 2>&1
cat $LOG
