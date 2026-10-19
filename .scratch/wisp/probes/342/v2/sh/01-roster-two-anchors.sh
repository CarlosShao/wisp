#!/bin/sh
# 独立名册尺（票 342 落地腿 342-r1）：两把锚并集，主题起头锚定行首（A1045 定式）。
# 不采信编排者数的 17 笔；自己另起一把形。整腿 diff --name-status 在 02 号件。
cd "$(git rev-parse --show-toplevel)" || exit 9
OUT=/d/tmp/wisp/20261019/342v2
mkdir -p "$OUT" || exit 9
git log --reverse --format='%H%x09%s' > "$OUT/subjects-all.tsv"
printf 'TOTAL_LINES=%s\n' "$(wc -l < "$OUT/subjects-all.tsv")"
# 单次导出 + 内存内并集（保持 log 顺序，⛔ sort -u）
awk -F'\t' '
  index($2,"票 342 落地腿 342-r1")==1 || index($2,"342-r1")==1 {print}
' "$OUT/subjects-all.tsv" > "$OUT/roster-union.tsv"
printf 'UNION_COUNT=%s\n' "$(wc -l < "$OUT/roster-union.tsv")"
awk -F'\t' 'index($2,"票 342 落地腿 342-r1")==1' "$OUT/subjects-all.tsv" | wc -l | sed 's/^/ANCHOR1_HEAD=/'
awk -F'\t' 'index($2,"342-r1")==1' "$OUT/subjects-all.tsv" | wc -l | sed 's/^/ANCHOR2_HEAD_BARE=/'
awk -F'\t' 'index($2,"342-r1")>0' "$OUT/subjects-all.tsv" | wc -l | sed 's/^/ANYWHERE_SUBSTRING=/'
printf '--- UNION FIRST(oldest)=%s LAST(newest)=%s ---\n' "$(head -1 "$OUT/roster-union.tsv" | cut -c1-12)" "$(tail -1 "$OUT/roster-union.tsv" | cut -c1-12)"
printf '--- UNION SUBJECTS ---\n'
cut -f2 "$OUT/roster-union.tsv" | cut -c1-90
