#!/bin/bash
# 353-v1 R16 = §C 那味独立复量（尺＝grep/awk 打在 r1 入库的仓外底片 alt-*/swing-*；量的）
L=/d/tmp/wisp/20261020/353r1/logs
echo "###### 0) 腿自己的 alt-summary 逐字"
cat "$L/alt-summary.txt"
echo "###### 1) 我从 alt-*.raw.txt 逐枚重取分子（provedReparse/flipErrs/truncated/deleted）"
for f in "$L"/alt-*-*.raw.txt; do
  n=$(basename "$f")
  line=$(grep -hoE 'provedReparse=[0-9]+[^ ]*|flipErrs=[0-9]+|truncated=[0-9]+|deleted=[0-9]+' "$f" | tr '\n' ' ')
  form=$(grep -hoE 'form=[a-z-]+' "$f" | sort -u | tr '\n' ',' )
  printf '%-26s forms[%s] %s\n' "$n" "$form" "${line:-(无分子)}"
done
echo "###### 2) 分子行原文逐字（每发第一枚含 flipErrs 的行）"
for f in "$L"/alt-*-*.raw.txt; do
  echo "-- $(basename $f)"
  grep -h 'flipErrs' "$f" | sed -e 's/^\(.\{280\}\).*/\1/' | sed -n '1,3p'
done
echo "###### 3) 另一批（sh/27 swing-*）：truncated 逐枚"
for f in "$L"/swing-*.raw.txt; do
  n=$(basename "$f")
  printf '%-24s %s\n' "$n" "$(grep -hoE 'truncated=[0-9]+' "$f" | sort | uniq -c | tr '\n' ' ')"
done
echo "-- swing-summary 逐字 --"
cat "$L/swing-summary.txt"
echo "###### 4) 正控：这些尺若报 0，先证明分子串存在于底片"
printf 'grep provedReparse 命中文件数=%s  总命中=%s\n' "$(grep -rl 'provedReparse' "$L" | wc -l)" "$(grep -rho 'provedReparse' "$L" | wc -l)"
printf 'grep flipErrs    命中文件数=%s  总命中=%s\n' "$(grep -rl 'flipErrs' "$L" | wc -l)" "$(grep -rho 'flipErrs' "$L" | wc -l)"
printf 'grep truncated   命中文件数=%s  总命中=%s\n' "$(grep -rl 'truncated' "$L" | wc -l)" "$(grep -rho 'truncated' "$L" | wc -l)"
printf 'grep "deleted="  命中文件数=%s  总命中=%s\n' "$(grep -rl 'deleted=' "$L" | wc -l)" "$(grep -rho 'deleted=' "$L" | wc -l)"
