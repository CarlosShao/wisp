#!/bin/bash
# 353-v1 R17 = A 族（ancestorunlink_348）逐发行全量提取（量的；尺＝grep -oE 打在 r1 入库底片）
L=/d/tmp/wisp/20261020/353r1/logs
echo "###### 1) swing-* 全部 form 行（ancestorunlink_348 那族的分子）"
for t in finalB finalL; do for s in 1 2 3 4; do
  f="$L/swing-$t-$s.raw.txt"; [ -f "$f" ] || continue
  echo "-- $t shot=$s --"
  grep -oE 'form=[a-z-]+ attempts=400 cycles=[0-9]+ provedReparse=[0-9]+ formBLegRefusals=[0-9]+ truncated=[0-9]+ deleted=[0-9]+ returnedNil=[0-9]+ flipErrs=[0-9]+' "$f" | sed 's/^/     /'
done; done
echo "###### 2) 整包 final-full-* 里同族分子行（★不同台面、⛔ 与定向混）"
for t in B L; do for s in 1 2; do
  f="$L/final-full-$t-$s.raw.txt"; [ -f "$f" ] || continue
  echo "-- full $t shot=$s --"
  grep -oE 'form=[a-z-]+ attempts=[0-9]+ cycles=[0-9]+ provedReparse=[0-9]+ formBLegRefusals=[0-9]+ truncated=[0-9]+ deleted=[0-9]+ returnedNil=[0-9]+ flipErrs=[0-9]+' "$f" | sed 's/^/     /'
done; done
echo "###### 3) truncated>0 的全部底片（尺＝grep -rl；正控见 R16 第4节）"
grep -rlE 'truncated=[1-9]' "$L" | sed 's/^/   /'
echo "###### 4) 逐枚 truncated 非零行原文（⛔ 截断）"
grep -rhoE 'form=[a-z-]+ attempts=400 cycles=[0-9]+ provedReparse=[0-9]+ formBLegRefusals=[0-9]+ truncated=[0-9]+ deleted=[0-9]+' "$L" | grep -E 'truncated=[1-9]' | sort | uniq -c | sed 's/^/   /'
