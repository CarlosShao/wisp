#!/bin/bash
# 353-v1 R25 = 跨批新增红硬尺（六批落地侧红名并集 − 六批基线侧红名并集）＋两味精确行段
cd '/d/work/workspace/projects plans/Wisp' || exit 9
L=/d/tmp/wisp/20261020/353r1/logs
V=/d/tmp/wisp/20261020/353v1
mkdir -p $V/x
: > $V/x/B.red; : > $V/x/L.red
for f in "$L"/final-full-B-1.raw.txt "$L"/final-full-B-2.raw.txt "$L"/mx-treeB2-1.raw.txt "$L"/mx-treeB2-2.raw.txt "$L"/fullpair-s1-treeB.raw.txt "$L"/pername-s1-treeB2.raw.txt "$L"/colpair-s1-treeB2.raw.txt; do
  [ -f "$f" ] || continue
  grep -E '^--- FAIL' "$f" | sed 's/^--- FAIL: //; s/ *([0-9.]*s)$//' | sed "s|^|$(basename $f)\t|" >> $V/x/B.detail
  grep -E '^--- FAIL' "$f" | sed 's/^--- FAIL: //; s/ *([0-9.]*s)$//' >> $V/x/B.red
done
for f in "$L"/final-full-L-1.raw.txt "$L"/final-full-L-2.raw.txt "$L"/mx-treeL-1.raw.txt "$L"/mx-treeL-2.raw.txt "$L"/fullpair-s2-treeL.raw.txt "$L"/pername-s2-treeL.raw.txt "$L"/colpair-s2-treeL.raw.txt; do
  [ -f "$f" ] || continue
  grep -E '^--- FAIL' "$f" | sed 's/^--- FAIL: //; s/ *([0-9.]*s)$//' | sed "s|^|$(basename $f)\t|" >> $V/x/L.detail
  grep -E '^--- FAIL' "$f" | sed 's/^--- FAIL: //; s/ *([0-9.]*s)$//' >> $V/x/L.red
done
sort -u $V/x/B.red > $V/x/B.u; sort -u $V/x/L.red > $V/x/L.u
echo "  基线批数=7 落地批数=7；红名并集 B=$(grep -c . $V/x/B.u) L=$(grep -c . $V/x/L.u)"
echo "  跨批新增红（L−B）枚数=$(comm -13 $V/x/B.u $V/x/L.u | grep -c . || true)"
comm -13 $V/x/B.u $V/x/L.u | sed 's/^/   NEWRED_XBATCH: /'
echo "  跨批只在基线侧（B−L）=$(comm -23 $V/x/B.u $V/x/L.u | grep -c . || true)"
comm -23 $V/x/B.u $V/x/L.u | sed 's/^/   onlyB_XBATCH: /'
echo "  两侧都见过的名（=交集）："; comm -12 $V/x/B.u $V/x/L.u | sed 's/^/   /'
echo "###### 两味精确行段（HEAD blob 现量）"
N=$V/blob/HEAD-nametopup_353_windows_test.go.txt
awk 'NR>=28 && NR<=41{printf "   %3d| %.130s\n", NR, $0}' "$N"
echo "-- GREEN SCOPE 语句跨行 --"
awk 'NR>=146 && NR<=165{printf "   %3d| %.110s\n", NR, $0}' "$N"
echo "###### 仪器件里 damage/entered/frozen/refused 四枚计数器的使用点（⛔ 合并⛔ 相减硬尺）"
grep -nE 'entered|frozen|damage|refused' "$N" | grep -E '\+\+|==|>|if ' | cut -c1-150 | sed 's/^/   /'
