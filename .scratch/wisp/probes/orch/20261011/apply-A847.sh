#!/bin/sh
# A847 落盘：两张票追加＋翻框、台账追加、停车点追加（全部 append-only；⛔ 改原句）
cd "$(git rev-parse --show-toplevel)" || exit 9
D=.scratch/wisp/probes/orch/20261011
T298=.scratch/wisp/issues/298-two-gofmt-dirty-test-files-at-head-and-one-crlf-fake-cell.md
T309=.scratch/wisp/issues/309-assessorfor-empty-taskid-limb-silently-drops-r4-detector.md
T308=.scratch/wisp/issues/308-spelling-invariance-test-panics-when-temp-has-no-volume.md
LEDGER=docs/reports/pending-and-issues.md
HAND=docs/reports/HANDOVER.md
OUT=$D/logs/orch-A847-apply.txt
LOG=$D/logs

# 先记 HEAD 层前缀行数（追加完必须逐字节等于原 HEAD 正文）
LP=$(git show HEAD:$LEDGER | wc -l); HP=$(git show HEAD:$HAND | wc -l)

cat $D/ticket-309-append-A847.md >> $T309
cat $D/ticket-308-append-A847.md >> $T308
cat $D/A847.txt >> $LEDGER
cat $D/HANDOVER-4.0cl.md >> $HAND

# 翻框：⛔ 按行号，一律按内容锚（AC 号就在框行开头）
for n in 4 5; do
  sed -i "0,/^[[:space:]]*- \[ \] \*\*\`AC#$n\`/ s/^[[:space:]]*- \[ \] \*\*\`AC#$n\`/- [x] **\`AC#$n\`/" $T309
done
for n in 1 2 3 4 5; do
  sed -i "0,/^[[:space:]]*- \[ \] \*\*\`AC#$n\`/ s/^[[:space:]]*- \[ \] \*\*\`AC#$n\`/- [x] **\`AC#$n\`/" $T308
done

git add -- $D $T309 $T308 $LEDGER $HAND 2>$LOG/orch-A847-add-err.txt

{ echo "date=$(date '+%Y-%m-%d %H:%M:%S %z') HEAD=$(git rev-parse --short HEAD)"
  echo "=== 框普查（尺＝含缩进那一把） ==="
  for t in $T309 $T308 $T298; do
    echo "  $(basename $t | cut -c1-3): checked=$(grep -cE '^[[:space:]]*- \[x\]' $t) unchecked=$(grep -cE '^[[:space:]]*- \[ \]' $t)"
  done
  echo "=== 前缀必须⛔ 被动过（＝只追加的证明） ==="
  git show HEAD:$LEDGER | cmp -s - <(head -$LP $LEDGER) && echo "  prefix_ledger=IDENTICAL ($LP 行)" || echo "  prefix_ledger=DIFF"
  git show HEAD:$HAND | cmp -s - <(head -$HP $HAND) && echo "  prefix_handover=IDENTICAL ($HP 行)" || echo "  prefix_handover=DIFF"
  echo "  ledger_lines=$(wc -l < $LEDGER) (HEAD=$LP)   handover_lines=$(wc -l < $HAND) (HEAD=$HP)"
  echo "=== 删除行必须＝翻框枚数（⛔ 别的删除） ==="
  echo "  numstat_309=$(git diff --cached --numstat -- $T309 | tr '\n' ' ')"
  echo "  numstat_308=$(git diff --cached --numstat -- $T308 | tr '\n' ' ')"
  echo "  numstat_ledger=$(git diff --cached --numstat -- $LEDGER | tr '\n' ' ')"
  echo "  numstat_handover=$(git diff --cached --numstat -- $HAND | tr '\n' ' ')"
  echo "=== 逐枚删除行原文（人眼过一遍：应全是框行的 - [ ] 形式） ==="
  git diff --cached -- $T309 $T308 | grep '^-' | grep -v '^---'
  echo "=== 暂存名册卫生（⛔ 关 quotePath，A846 那把假指控就是这么修的） ==="
  git -c core.quotePath=false diff --cached --name-only > $LOG/staged-A847.txt
  echo "  staged_total=$(wc -l < $LOG/staged-A847.txt) staged_go=$(grep -c '[.]go$' $LOG/staged-A847.txt) foreign=$(grep -vc '^docs/reports/\|^\.scratch/wisp/issues/30[89]\|^\.scratch/wisp/probes/orch/20261011/' $LOG/staged-A847.txt)"
  echo "  zero_byte_staged=$(git diff --cached --name-only -z | while IFS= read -r -d '' f; do [ -s "$f" ] || echo "$f"; done | grep -c '')"
  echo "  worktree_dirty_internal_cmd=$(git status --porcelain -- internal cmd | wc -l | tr -d ' ')"
} > $OUT 2>&1
cat $OUT
