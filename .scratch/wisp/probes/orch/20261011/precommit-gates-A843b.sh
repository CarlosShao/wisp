#!/bin/sh
# 编排者本批门禁第二把（改写正文后重跑；原始件＝logs/orch-A843-gates.txt ⛔ 覆盖⛔ 删，那枚坏尺读数留在原地）
cd "$(git rev-parse --show-toplevel)" || exit 9
D=.scratch/wisp/probes/orch/20261011
LOG=$D/logs/orch-A843-gates-b.txt

sh scripts/d22scan.sh > $D/logs/d22scan-A843-b.txt 2>&1; rc1=$?
sh scripts/check-path-length-budget.sh --with-self-test > $D/logs/lenbudget-A843-b.txt 2>&1; rc2=$?

PAT_HDR='^### 4\.0c[a-z]? 停车点'
PAT_LED='^## A84[0-9]'

git add -- docs/reports/pending-and-issues.md docs/reports/HANDOVER.md $D 2>/dev/null
git diff --cached --name-only > $D/logs/staged-A843-b.txt
n_all=$(grep -c '' $D/logs/staged-A843-b.txt || true)
n_go=$(grep -c '\.go$' $D/logs/staged-A843-b.txt || true)
git diff --cached --numstat -- docs/reports/pending-and-issues.md > $D/logs/ledger-numstat-A843-b.txt
git diff --cached --numstat -- docs/reports/HANDOVER.md > $D/logs/handover-numstat-A843-b.txt
dirty=$(git status --porcelain -- internal cmd | wc -l | tr -d ' ')

{
  echo "d22scan_rc=$rc1"
  echo "lenbudget_rc=$rc2"
  echo "staged_total=$n_all staged_go=$n_go"
  echo "ledger_numstat=$(cat $D/logs/ledger-numstat-A843-b.txt)"
  echo "handover_numstat=$(cat $D/logs/handover-numstat-A843-b.txt)"
  echo "worktree_dirty_internal_cmd=$dirty"
  echo "handover_4_0c_headers_b=$(grep -c "$PAT_HDR" docs/reports/HANDOVER.md)"
  echo "ledger_A84x_headers_b=$(grep -c "$PAT_LED" docs/reports/pending-and-issues.md)"
  echo "prefix_ledger_vs_head=$(git show HEAD:docs/reports/pending-and-issues.md | cmp -s - <(head -$(git show HEAD:docs/reports/pending-and-issues.md | wc -l) docs/reports/pending-and-issues.md) && echo IDENTICAL || echo DIFF)"
  echo "prefix_handover_vs_head=$(git show HEAD:docs/reports/HANDOVER.md | cmp -s - <(head -$(git show HEAD:docs/reports/HANDOVER.md | wc -l) docs/reports/HANDOVER.md) && echo IDENTICAL || echo DIFF)"
  echo "lines_ledger=$(grep -c '' docs/reports/pending-and-issues.md) lines_ledger_head=$(git show HEAD:docs/reports/pending-and-issues.md | wc -l)"
  echo "lines_handover=$(grep -c '' docs/reports/HANDOVER.md) lines_handover_head=$(git show HEAD:docs/reports/HANDOVER.md | wc -l)"
  echo "A843_file_lines=$(grep -c '' $D/A843.txt) sec40ch_file_lines=$(grep -c '' $D/HANDOVER-4.0ch.md)"
} > $LOG 2>&1
cat $LOG
