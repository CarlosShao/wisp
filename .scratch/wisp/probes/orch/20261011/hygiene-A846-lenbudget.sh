#!/bin/sh
# A846 车队卫生第五笔：长度门那枚 RED 归位（⛔ 删内容、⛔ 加长程豁免＝只改名一枚探针件，逐字节 cmp 自证）
cd "$(git rev-parse --show-toplevel)" || exit 9
D=.scratch/wisp/probes/orch/20261011
W=$D/writeface-1/logs
OLD=$W/head-296-hotkey-reload-source-skips-applyhotkeydefaults-keys-die-1s-after-boot.md
NEW=$W/head-ticket-296.md
LOG=$D/logs/orch-A846-hygiene-lenbudget.txt

{ echo "date=$(date '+%Y-%m-%d %H:%M:%S %z') HEAD=$(git rev-parse --short HEAD)"
  echo "RED 出处逐字（⛔ 抹，件留盘）＝$(grep 'NOT in the roster' $D/logs/lenbudget-A846.txt | head -1 | cut -c1-200)"
  echo "枚数尺＝$(grep 'denominator:' $D/logs/lenbudget-A846.txt | cut -c1-190)"
  echo "引用尺：谁引过这枚旧名？grep -rln 'head-296-hotkey-reload' 全仓（除它自己那枚 RED 读数件）＝$(grep -rln 'head-296-hotkey-reload' --include='*.md' --include='*.txt' --include='*.sh' . 2>/dev/null | grep -v "$OLD" | grep -c '')"
} > $LOG 2>&1

# 内容自证＝改名前后各对 HEAD 那枚 blob 比一次（⛔ 另建一份重复件进仓，49 KB 不该复制两遍）
git mv "$OLD" "$NEW" 2>> $LOG
{ echo "cmp 改名后件 ⇄ HEAD 的旧名 blob：$(git show HEAD:"$OLD" | cmp -s - "$NEW" && echo IDENTICAL_rc0 || echo DIFF) 字节=$(wc -c < "$NEW")"
  echo "cmp 改名后件 ⇄ HEAD 的票 296 全文（同一份内容的另一名）：$(git show HEAD:.scratch/wisp/issues/296-hotkey-reload-source-skips-applyhotkeydefaults-keys-die-1s-after-boot.md | cmp -s - "$NEW" && echo IDENTICAL_rc0 || echo DIFF)"
  echo "改名后相对长度=$(printf '%s' "$NEW" | wc -c | tr -d ' ')（帽＝121）；名长=$(basename "$NEW" | wc -c | tr -d ' ')（帽＝100，wc -c 多数到的换行 1）"
  echo "旧名件现在还在吗：$([ -e "$OLD" ] && echo YES_不该 || echo NO_已改名)"
} >> $LOG 2>&1

sh scripts/check-path-length-budget.sh --with-self-test > $D/logs/lenbudget-A846-after-rename.txt 2>&1
echo "lenbudget_after_rename_rc=$?" >> $LOG
tail -6 $D/logs/lenbudget-A846-after-rename.txt >> $LOG 2>&1
cat $LOG
