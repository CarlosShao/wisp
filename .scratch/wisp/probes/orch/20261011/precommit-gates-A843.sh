#!/bin/sh
# 编排者本批门禁（docs-only 批：台账 A843＋停车点 §4.0ch＋编排者复跑件；⛔ 碰 .go）
cd "$(git rev-parse --show-toplevel)" || exit 9
D=.scratch/wisp/probes/orch/20261011
LOG=$D/logs/orch-A843-gates.txt

sh scripts/d22scan.sh > $D/logs/d22scan-A843.txt 2>&1; rc1=$?
sh scripts/check-path-length-budget.sh --with-self-test > $D/logs/lenbudget-A843.txt 2>&1; rc2=$?

git diff --cached --name-only > $D/logs/staged-A843.txt
n_all=$(grep -c '' $D/logs/staged-A843.txt || true)
n_go=$(grep -c '\.go$' $D/logs/staged-A843.txt || true)
git diff --cached --numstat -- docs/reports/pending-and-issues.md > $D/logs/ledger-numstat-A843.txt
git diff --cached --numstat -- docs/reports/HANDOVER.md > $D/logs/handover-numstat-A843.txt
dirty=$(git status --porcelain -- internal cmd | grep -c '' || true)

{
  echo "d22scan_rc=$rc1"
  echo "lenbudget_rc=$rc2"
  echo "staged_total=$n_all staged_go=$n_go"
  echo "ledger_numstat=$(cat $D/logs/ledger-numstat-A843.txt)"
  echo "handover_numstat=$(cat $D/logs/handover-numstat-A843.txt)"
  echo "worktree_dirty_internal_cmd=$dirty"
  echo "handover_4_0c_headers=$(grep -c '^### 4\\.0c' docs/reports/HANDOVER.md)"
  echo "ledger_A84x_headers=$(grep -c '^## A84' docs/reports/pending-and-issues.md)"
} > $LOG
cat $LOG
