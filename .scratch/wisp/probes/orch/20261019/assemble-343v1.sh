#!/bin/bash
set -u
D="/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/orch/20261019"
SIG=$(printf '\302\247')   # U+00A7 section sign（被禁的那枚）
cat "$D/343v1-head.txt" "$D/common-20261019e.txt" "$D/343v1-specific.txt" > "$D/343v1-dispatch.txt"
for f in 343v1-head.txt 343v1-specific.txt common-20261019e.txt 343v1-dispatch.txt; do
  echo "$f lines=$(wc -l < "$D/$f") bytes=$(wc -c < "$D/$f") U+00A7=$(grep -c "$SIG" "$D/$f" 2>/dev/null || true)"
done
echo "=== 波次字样核对（d 须 0、e/第四波须≥1）==="
echo "  第四波=$(grep -c '第四波' "$D/343v1-dispatch.txt")  第三波 d=$(grep -c '第三波 d' "$D/343v1-dispatch.txt")  common-c=$(grep -c 'common-20261018c' "$D/343v1-dispatch.txt")"
echo "=== 持有者表述（343-v1＝你 须≥1）==="
grep -c '342-v1' "$D/343v1-dispatch.txt"
echo "=== ★放腿前名册尺：锚定行首（改正形，⛔ --grep 正文命中、⛔ 硬编码分隔符）==="
cd "/d/work/workspace/projects plans/Wisp" || exit 9
git log --all --format='%H%x09%s' > "$D/allsubj.txt"
for p in "票 343 终裁腿 343-v1" "票 343 落地腿 343-r1" "票 342 终裁腿 342-v1"; do
  n=$(awk -F'\t' -v P="$p" 'index($2,P)==1' "$D/allsubj.txt" | wc -l)
  echo "  $p => $n 笔"
done
echo "=== 归档目录就位＋三枚写面静默 ==="
[ -d "$D" ] && echo "  orch/20261019 存在"
for x in "$D/343v1-head.txt" "$D/343v1-specific.txt" "$D/343v1-dispatch.txt" "$D/common-20261019e.txt"; do [ -e "$x" ] && echo "  OK $(basename $x)" || echo "  MISSING $(basename $x)"; done
echo "=== 写面被动（40 分）risk/winsec/config ==="
for f in internal/risk internal/winsec internal/config; do echo "  $f=$(find "$f" -type f -newermt '-40 minutes' 2>/dev/null | wc -l)"; done
