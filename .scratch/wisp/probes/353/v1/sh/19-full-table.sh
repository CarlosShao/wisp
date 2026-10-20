#!/bin/bash
# 353-v1 R19 = 逐枚底片的 A 族分子全表（量的；尺＝grep -oE 逐文件，⛔ 折叠）
L=/d/tmp/wisp/20261020/353r1/logs
echo "###### 逐文件 × 逐 form：truncated/deleted/flipErrs/provedReparse"
for f in "$L"/alt-*.raw.txt "$L"/swing-*.raw.txt "$L"/recheck*-*.raw.txt "$L"/mx-*.raw.txt "$L"/final-full-*.raw.txt "$L"/fullpair-*.raw.txt "$L"/pername-*.raw.txt; do
  [ -f "$f" ] || continue
  n=$(basename "$f")
  grep -oE 'form=[a-z-]+ attempts=400 cycles=[0-9]+ provedReparse=[0-9]+ formBLegRefusals=[0-9]+ truncated=[0-9]+ deleted=[0-9]+ returnedNil=[0-9]+ flipErrs=[0-9]+' "$f" | while read -r r; do
    fr=$(echo "$r" | sed 's/ .*//')
    tr=$(echo "$r" | sed 's/.*truncated=\([0-9]*\).*/\1/')
    de=$(echo "$r" | sed 's/.*deleted=\([0-9]*\).*/\1/')
    fe=$(echo "$r" | sed 's/.*flipErrs=\([0-9]*\)/\1/')
    pr=$(echo "$r" | sed 's/.*provedReparse=\([0-9]*\).*/\1/')
    printf '%-30s %-16s provedReparse=%-4s truncated=%-2s deleted=%s flipErrs=%s\n' "$n" "$fr" "$pr" "$tr" "$de" "$fe"
  done
done
echo "###### 全底片 truncated 非零总表（含 form ⇄ 文件）"
grep -rEo 'form=[a-z-]+ attempts=400 cycles=[0-9]+ provedReparse=[0-9]+ formBLegRefusals=[0-9]+ truncated=[0-9]+ deleted=[0-9]+' "$L" | grep -E 'truncated=[1-9]' | sed 's|/d/tmp/wisp/20261020/353r1/logs/||' | sort | sed 's/^/  /'
echo "###### 正控（同一把尺命中的非零枚数应 >0；零枚＝尺坏）"
