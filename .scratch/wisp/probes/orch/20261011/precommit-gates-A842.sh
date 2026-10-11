#!/bin/sh
# 编排者本批门禁（docs-only 批：只动 .md／.txt/.sh 证据件，⛔ 碰 .go）
# 两相：先扫门（写原始件），后按 index 出账（这枚件自己最后落，所以它的 rc 行⛔ 依赖自己）
cd "$(git rev-parse --show-toplevel)" || exit 9

D=.scratch/wisp/probes/orch/20261011
LOG=$D/logs/orch-A842-gates.txt

# 相 1：全仓形状门＋长度门（先跑，读数落原始件）
sh scripts/d22scan.sh > $D/logs/d22scan-A842.txt 2>&1; rc1=$?
sh scripts/check-path-length-budget.sh --with-self-test > $D/logs/lenbudget-A842.txt 2>&1; rc2=$?

# 相 2：按暂存面出账（调用前必须已经 git add 完）
git diff --cached --name-only > $D/logs/staged-A842.txt
n_all=$(grep -c '' $D/logs/staged-A842.txt || true)
n_go=$(grep -c '\.go$' $D/logs/staged-A842.txt || true)
git diff --cached --numstat -- docs/reports/pending-and-issues.md > $D/logs/ledger-numstat.txt
git diff --cached --numstat -- docs/reports/HANDOVER.md > $D/logs/handover-numstat.txt
git diff --cached --numstat -- .scratch/wisp/issues/ > $D/logs/ticket-numstat.txt
dirty=$(git status --porcelain | grep -c '' || true)

{
  echo "d22scan_rc=$rc1"
  echo "lenbudget_rc=$rc2"
  echo "staged_total=$n_all staged_go=$n_go"
  echo "ledger_numstat=$(cat $D/logs/ledger-numstat.txt)"
  echo "handover_numstat=$(cat $D/logs/handover-numstat.txt)"
  echo "ticket_numstat=$(cat $D/logs/ticket-numstat.txt)"
  echo "worktree_dirty_lines=$dirty"
} > $LOG
cat $LOG
