#!/bin/sh
# 353-v1 R32 = truncated 非零行的"侧"计数（量的；尺＝逐文件 grep -oE，按文件名前缀定侧；⛔ 折叠）
L=/d/tmp/wisp/20261020/353r1/logs
V=/d/tmp/wisp/20261020/353v1
: > "$V/trunc-rows.txt"
for f in "$L"/*.raw.txt; do
  [ -f "$f" ] || continue
  n=$(basename "$f")
  grep -oE 'form=[a-z-]+ attempts=400 cycles=[0-9]+ provedReparse=[0-9]+ formBLegRefusals=[0-9]+ truncated=[0-9]+ deleted=[0-9]+' "$f" \
    | while read -r r; do
        tr=$(echo "$r" | sed 's/.*truncated=\([0-9]*\).*/\1/')
        [ "$tr" -gt 0 ] 2>/dev/null && printf '%-32s %s\n' "$n" "$r" >> "$V/trunc-rows.txt"
      done
done
echo "== truncated 非零分子行全表（逐文件）=="
cat "$V/trunc-rows.txt"
echo "== 总行数 =="
grep -c . "$V/trunc-rows.txt"
echo "== 按侧计数（文件名含 B/b 者归基线侧、含 L/l 者归落地侧；recheck-summary 例外＝腿自己的汇总件、⛔ 是 raw 原件）=="
grep -cE '^[a-z0-9-]*(B|b2|B2)[a-z0-9.-]*\.raw' "$V/trunc-rows.txt"
grep -cE '^[a-z0-9-]*(L|l)[a-z0-9.-]*\.raw' "$V/trunc-rows.txt"
echo "== 逐文件枚数一览 =="
cut -f1 -d' ' "$V/trunc-rows.txt" | tr -s ' ' | cut -d' ' -f1 | sort | uniq -c | sort -rn
