#!/bin/sh
# 353-r1 定因最后一尺＝A 族那枚 342 `truncated` 到底是"每次落地的伤率变了"还是"落地的次数变了"：
#   尺＝逐发取 truncated 与 provedReparse 算比值（★量的＝分子；★算的＝比值）；＋两侧 settle-teardown _fault_ 计数。
#   ★⛔ 跑计时；★小样本⛔ 外推，只报形状。
set -u
L=/d/tmp/wisp/20261020/353r1/logs
{
  echo "== 逐发：leaf-symlink 那枚的 cycles ⇄ truncated ⇄ 比值 =="
  for t in finalB finalL; do
    for s in 1 2 3 4; do
      f="$L/swing-$t-$s.raw.txt"; [ -f "$f" ] || continue
      line=$(grep -oE 'form=leaf-symlink attempts=400 cycles=[0-9]+ provedReparse=[0-9]+ formBLegRefusals=[0-9]+ truncated=[0-9]+ deleted=[0-9]+' "$f" | head -1)
      cyc=$(echo "$line" | sed 's/.*provedReparse=\([0-9]*\).*/\1/')
      tr=$(echo "$line" | sed 's/.*truncated=\([0-9]*\).*/\1/')
      del=$(echo "$line" | sed 's/.*deleted=\([0-9]*\).*/\1/')
      echo "   $t shot=$s cycles=$cyc truncated=$tr deleted=$del rate=$(awk -v a="$tr" -v b="$cyc" 'BEGIN{printf "%.3f", (b>0? a/b : 0)}')"
    done
  done
  echo "== settle-teardown fault 计数（那枚仪器自己的 fixture 收尾，两侧都会响＝台振）=="
  for t in finalB finalL; do
    c=0
    for s in 1 2 3 4; do
      f="$L/swing-$t-$s.raw.txt"; [ -f "$f" ] || continue
      grep -q 'settle the world back to the real state' "$f" && c=$((c+1))
    done
    echo "   $t: settle_fault_shots=$c / 4"
  done
  echo "== 整包那两发里同一族的分子（★⛔ 与定向发混用＝不同台面）=="
  for t in finalB finalL; do for s in 1 2; do
    f="$L/final-full-$t-$s.raw.txt"
    [ -f "$f" ] || f="$L/final-full-${t/B/b}-${s}.raw.txt"
    :
  done; done
  grep -hoE 'form=(symlink-dir|leaf-symlink) attempts=400 cycles=[0-9]+ provedReparse=[0-9]+ formBLegRefusals=[0-9]+ truncated=[0-9]+ deleted=[0-9]+' "$L"/final-full-B-*.raw.txt "$L"/final-full-L-*.raw.txt 2>/dev/null | sort | uniq -c | sed 's/^/   /'
} 2>&1
