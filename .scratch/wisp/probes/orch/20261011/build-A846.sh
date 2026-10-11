#!/bin/sh
# A846 批落盘（纯账目＋票面追加＋框翻勾；⛔ 碰 .go；⛔ 动在飞腿的写面）
cd "$(git rev-parse --show-toplevel)" || exit 9
D=.scratch/wisp/probes/orch/20261011
LED=docs/reports/pending-and-issues.md
HAND=docs/reports/HANDOVER.md
T298=.scratch/wisp/issues/298-two-gofmt-dirty-test-files-at-head-and-one-crlf-fake-cell.md
T308=.scratch/wisp/issues/308-spelling-invariance-test-panics-when-temp-has-no-volume.md
T294=.scratch/wisp/issues/294-no-production-shape-test-for-the-mute-attach-hop.md
T295=$(ls .scratch/wisp/issues/295-*.md)
T296=$(ls .scratch/wisp/issues/296-*.md)
PRE=$D/logs/A846-prefix-check.txt

# ---- 0) 起手：确认我只动自家与账目面，在飞腿的写面⛔ 被我碰 ----
{ echo "date=$(date '+%Y-%m-%d %H:%M:%S %z') HEAD=$(git rev-parse --short HEAD)"
  echo "inflight_dirty_internal=$(git status --porcelain -- internal cmd | tr '\n' '|')"
  echo "issues_dirty_before=$(git status --porcelain -- .scratch/wisp/issues | tr '\n' '|')"
  echo "docs_reports_dirty_before=$(git status --porcelain -- $LED $HAND | grep -c '')"
} > $PRE 2>&1
[ "$(git status --porcelain -- $LED $HAND | grep -c '')" = "0" ] || { echo "ABORT: 账目两件起手非干净（别人正改着）" >> $PRE; cat $PRE; exit 8; }

# ---- 1) 台账＝HEAD blob + 空行 + A846.txt，先证前缀逐字节相同 ----
git show HEAD:$LED > $D/tmp-ledger-head.bin
HL=$(git show HEAD:$LED | wc -l | tr -d ' ')
head -$HL $LED > $D/tmp-ledger-now-headlines.bin
cmp $D/tmp-ledger-head.bin $D/tmp-ledger-now-headlines.bin && PFX_LED=IDENTICAL || PFX_LED=DIFF
{ echo "$LED prefix_vs_HEAD=$PFX_LED head_lines=$HL"
  echo "A846_new_lines=$(grep -c '' $D/A846.txt)"
} >> $PRE
[ "$PFX_LED" = "IDENTICAL" ] || { echo "ABORT: 台账前缀漂移" >> $PRE; cat $PRE; exit 7; }
{ cat $D/tmp-ledger-head.bin; printf '\n'; cat $D/A846.txt; } > $LED

# ---- 2) 停车点＝同一手法 ----
git show HEAD:$HAND > $D/tmp-hand-head.bin
HH=$(git show HEAD:$HAND | wc -l | tr -d ' ')
head -$HH $HAND > $D/tmp-hand-now-headlines.bin
cmp $D/tmp-hand-head.bin $D/tmp-hand-now-headlines.bin && PFX_HAND=IDENTICAL || PFX_HAND=DIFF
echo "$HAND prefix_vs_HEAD=$PFX_HAND head_lines=$HH" >> $PRE
[ "$PFX_HAND" = "IDENTICAL" ] || { echo "ABORT: 停车点前缀漂移" >> $PRE; cat $PRE; exit 6; }
{ cat $D/tmp-hand-head.bin; printf '\n'; cat $D/HANDOVER-4.0ck.md; } > $HAND

# ---- 3) 五枚票面＝现件 + 空行 + 追加节（⛔ 改原句；本节只追加） ----
for pair in "$T298:ticket-298-append.md" "$T308:ticket-308-append.md" "$T294:ticket-294-append.md" "$T295:ticket-295-append.md" "$T296:ticket-296-append.md"; do
  f=${pair%%:*}; a=${pair##*:}
  OL=$(grep -c '' "$f")
  { cat "$f"; printf '\n'; cat "$D/$a"; } > "$f.tmp" && mv "$f.tmp" "$f"
  echo "ticket $f 追加前=$OL 追加后=$(grep -c '' $f) 新增=$(( $(grep -c '' $f) - OL ))" >> $PRE
done

# ---- 4) 翻框普查（逐票）＋节头数（HANDOVER 只追加⛔ 吞旧节） ----
{ echo "box_census_298=x $(grep -cE '^[[:space:]]*- \[x\]' $T298) open $(grep -cE '^[[:space:]]*- \[ \]' $T298)"
  echo "box_census_308=x $(grep -cE '^[[:space:]]*- \[x\]' $T308) open $(grep -cE '^[[:space:]]*- \[ \]' $T308)"
  echo "box_census_294=x $(grep -cE '^[[:space:]]*- \[x\]' $T294) open $(grep -cE '^[[:space:]]*- \[ \]' $T294)"
  echo "box_census_295=x $(grep -cE '^[[:space:]]*- \[x\]' $T295) open $(grep -cE '^[[:space:]]*- \[ \]' $T295)"
  echo "box_census_296=x $(grep -cE '^[[:space:]]*- \[x\]' $T296) open $(grep -cE '^[[:space:]]*- \[ \]' $T296)"
  echo "handover_4_0c_headers=$(grep -cE '^### 4\.0c[a-z]{1,2} 停车点' $HAND)  last=$(grep -E '^### 4\.0c[a-z]{1,2} 停车点' $HAND | tail -1 | cut -c1-40)"
  echo "ledger_A84x_headers=$(grep -cE '^## A84[0-9]' $LED)  last=$(grep -E '^## A84[0-9]' $LED | tail -1 | cut -c1-30)"
  echo "zero_byte_files_in_D=$(find $D -type f -size 0 | grep -c '')  dot_out_in_D=$(find $D -type f -name '*.out' | grep -c '')"
} >> $PRE
rm -f $D/tmp-ledger-head.bin $D/tmp-ledger-now-headlines.bin $D/tmp-hand-head.bin $D/tmp-hand-now-headlines.bin 2>/dev/null
cat $PRE
