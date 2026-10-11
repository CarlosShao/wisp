#!/bin/sh
# 编排者 A846 批门禁与暂存（账目＋五枚票面＋自家探针件；⛔ 碰 .go；⛔ 收在飞腿的件）
cd "$(git rev-parse --show-toplevel)" || exit 9
D=.scratch/wisp/probes/orch/20261011
LED=docs/reports/pending-and-issues.md
HAND=docs/reports/HANDOVER.md
LOG=$D/logs/orch-A846-gates.txt
T298=.scratch/wisp/issues/298-two-gofmt-dirty-test-files-at-head-and-one-crlf-fake-cell.md
T308=.scratch/wisp/issues/308-spelling-invariance-test-panics-when-temp-has-no-volume.md
T294=.scratch/wisp/issues/294-no-production-shape-test-for-the-mute-attach-hop.md
T295=.scratch/wisp/issues/295-ra-started-bare-bool-written-off-thread.md
T296=.scratch/wisp/issues/296-hotkey-reload-source-skips-applyhotkeydefaults-keys-die-1s-after-boot.md

sh scripts/d22scan.sh > $D/logs/d22scan-A846.txt 2>&1; rc1=$?
sh scripts/check-path-length-budget.sh --with-self-test > $D/logs/lenbudget-A846.txt 2>&1; rc2=$?

git add -- $LED $HAND $D $T298 $T308 $T294 $T295 $T296 2>/dev/null
git diff --cached --name-only > $D/logs/staged-A846.txt
n_all=$(grep -c '' $D/logs/staged-A846.txt)
n_go=$(grep -c '[.]go$' $D/logs/staged-A846.txt)
n_out=$(grep -c '[.]out$' $D/logs/staged-A846.txt)
foreign=$(git diff --cached --name-only | grep -v '^docs/reports/\|^\.scratch/wisp/issues/\|^\.scratch/wisp/probes/orch/20261011/' | grep -c '')
inflight=$(git diff --cached --name-only | grep -E '^internal/|^cmd/|^\.scratch/wisp/probes/(309|308)/r1/|^\.scratch/wisp/issues/309' | grep -c '')
zero=$(git diff --cached --name-only | xargs -I{} sh -c '[ -s "{}" ] || echo {}' 2>/dev/null | grep -c '')

{
  echo "d22scan_rc=$rc1"
  echo "lenbudget_rc=$rc2"
  echo "staged_total=$n_all staged_go=$n_go staged_dot_out=$n_out foreign_staged=$foreign inflight_staged=$inflight zero_byte_staged=$zero"
  echo "worktree_dirty_internal_cmd_after_stage=$(git status --porcelain -- internal cmd | wc -l | tr -d ' ')  # 剩下这些⛔ 我的"
  echo "numstat_ledger=$(git diff --cached --numstat -- $LED | tr '\n' ' ')"
  echo "numstat_handover=$(git diff --cached --numstat -- $HAND | tr '\n' ' ')"
  for t in "$T298" "$T308" "$T294" "$T295" "$T296"; do
    echo "numstat_ticket_$(basename "$t" | cut -c1-3)=$(git diff --cached --numstat -- "$t" | tr '\n' ' ')"
  done
  echo "prefix_ledger_vs_head=$(git show HEAD:$LED | cmp -s - <(head -"$(git show HEAD:$LED | wc -l)" "$LED") && echo IDENTICAL || echo DIFF)"
  echo "prefix_handover_vs_head=$(git show HEAD:$HAND | cmp -s - <(head -"$(git show HEAD:$HAND | wc -l)" "$HAND") && echo IDENTICAL || echo DIFF)"
  echo "lines_ledger=$(grep -c '' $LED) head_lines_ledger=$(git show HEAD:$LED | wc -l)"
  echo "lines_handover=$(grep -c '' $HAND) head_lines_handover=$(git show HEAD:$HAND | wc -l)"
  echo "handover_4_0c_headers=$(grep -cE '^### 4[.]0c[a-z]{1,2} 停车点' $HAND)  last=$(grep -E '^### 4[.]0c[a-z]{1,2} 停车点' $HAND | tail -1 | cut -c5-30)"
  echo "ledger_A84x_headers=$(grep -cE '^## A84[0-9]' $LED)"
  echo "box_298=$(grep -cE '^[[:space:]]*- \[x\]' $T298)/$(grep -cE '^[[:space:]]*- \[ \]' $T298) box_308=$(grep -cE '^[[:space:]]*- \[x\]' $T308)/$(grep -cE '^[[:space:]]*- \[ \]' $T308) box_294=$(grep -cE '^[[:space:]]*- \[x\]' $T294)/$(grep -cE '^[[:space:]]*- \[ \]' $T294) box_295=$(grep -cE '^[[:space:]]*- \[x\]' $T295)/$(grep -cE '^[[:space:]]*- \[ \]' $T295) box_296=$(grep -cE '^[[:space:]]*- \[x\]' $T296)/$(grep -cE '^[[:space:]]*- \[ \]' $T296)"
  echo "d22scan_tail=$(tail -3 $D/logs/d22scan-A846.txt | tr '\n' '|')"
  echo "lenbudget_tail=$(tail -3 $D/logs/lenbudget-A846.txt | tr '\n' '|')"
} > $LOG 2>&1
cat $LOG
