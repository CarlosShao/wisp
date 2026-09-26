#!/usr/bin/env bash
# 157 AC#5 recheck — probe 06: CORRECTED column ruler; supersedes the cell-level numbers in 02 (R-D/R-E), 03 (S-B/S-D) and 05 (U-D/U-E).
# Why this file exists (my own instrument defect, named):
#   gawk 5.3.2 regex constant /\\\|/  ==  (literal backslash) ALTERNATION (empty)  -> the empty branch matches EVERYWHERE,
#   so gsub(/\\\|/,"SH") inserted "SH" between every character and inflated every byte count.
#   Correct form: /\\[|]/  (backslash followed by a bracketed pipe). Caught by the positive control in 02 R-D0
#   (the yardstick rows printed as "rowSH SH1SH SH" — a shape no real cell has).
# Section counts that used grep -c on whole lines are NOT affected and are not re-run here.
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9
F=docs/evidence/s1/155-three-unjudged-cells-r1.md
G=docs/evidence/s1/153-trace-lies-unguarded-r1.md
A=docs/evidence/s1/155-three-unjudged-cells-r1-accept-r1.md
NOW=$(git rev-parse --short HEAD)
ANCH=4f6c14c

echo "### V-A  self-test of the fixed ruler (must be 0 spurious SH insertions)"
bad=$(printf '| a | b \\| c | d |\n' | awk '{gsub(/\\[|]/,"Q"); if ($0 ~ /Q c/) print "ok"}' )
printf '   fixed form keeps the row readable: %s\n' "${bad:-FAILED}"
printf '| x | \\| y |\n' | awk '{n=gsub(/\\[|]/,"Q"); m=gsub(/Q/,"Q"); printf "   gsub hits=%d  Q count=%d  (must be equal)\n", n, m}'
printf '| x | plain row without escapes |\n' | awk '{n=gsub(/\\[|]/,"Q"); printf "   gsub hits on an escape-free row=%d  (must be 0)\n", n}'

echo "### V-B  B.12 12.2 (5 rows) cell map, fixed ruler"
git cat-file blob "$ANCH:$F" | awk -F'|' 'index($0,"### B.12 ")==1{f=1;next} f&&/^### /{f=0} f&&/^\| [0-9]+ \|/{
  gsub(/\\[|]/,"Q"); printf "   row%s | item=%.54s | why=%.44s | closure-bytes=%d | closure=%.70s\n", $2, $3, $4, length($(NF-1)), $(NF-1)}'

echo "### V-C  B.15 15.3 (7 rows) cell map, fixed ruler"
git cat-file blob "$ANCH:$F" | awk -F'|' '/^\*\*15\.3/{f=1;next} f&&/^\*\*15\.4/{f=0} f&&/^\| [0-9]+ \|/{
  gsub(/\\[|]/,"Q"); printf "   row%s | item=%.54s | closure-bytes=%d | closure=%.70s\n", $2, $3, length($(NF-1)), $(NF-1)}'

echo "### V-D  field presence per row with the FIXED ruler (who / which-ruler-or-action / cost / clock-duration)"
run_fields() { # $1=label $2=rows
  printf '%s\n' "$2" | awk -F'|' -v lbl="$1" '{gsub(/\\[|]/,"Q"); cl=$(NF-1); n++;
    if (cl ~ /编排者|验收程|本程|任一枚|它/) W++;
    if (cl ~ /尺|一发|一次|按|重跑|并排|现读|名册|show|grep|sed|awk|抽|换/) R++;
    if (cl ~ /不需|需要|批准/) C++;
    if (cl ~ /分钟|小时|秒/) D++;
    if (length(cl)<10) S++ }
    END{printf "   %-26s rows=%-3s who=%-3s action-or-ruler=%-3s cost=%-3s clock=%-3s cell-under-10-bytes=%s\n", lbl, n, W+0, R+0, C+0, D+0, S+0}'
}
run_fields "B.12 12.2" "$(git cat-file blob "$ANCH:$F" | awk 'index($0,"### B.12 ")==1{f=1;next} f&&/^### /{f=0} f&&/^\| [0-9]+ \|/')"
run_fields "B.15 15.3" "$(git cat-file blob "$ANCH:$F" | awk '/^\*\*15\.3/{f=1;next} f&&/^\*\*15\.4/{f=0} f&&/^\| [0-9]+ \|/')"
run_fields "ctl 155-accept 7.2" "$(awk 'NR>=265 && NR<=300' "$A" | grep -E '^\| [0-9]+ \|')"
echo "   (the ctl line in probe 05 read who=2/7 with the BROKEN ruler; here is the same yardstick re-measured)"

echo "### V-E  forward table: per-row grade cell, fixed ruler (the criterion 2 field '明写未裁')"
echo "-- 12.1 (implementer self-proof, 5 rows):"
git cat-file blob "$ANCH:$F" | awk -F'|' 'index($0,"### B.12 ")==1{f=1;next} f&&/^### /{f=0} f&&/^\| \*\*AC#/{gsub(/\\[|]/,"Q"); printf "   %-9s section-col=%.30s | self-evidence-bytes=%d | grade=%.34s\n", $2, $3, length($4), $(NF-1)}'
echo "-- 15.1 (grades transcribed from the non-implementer, 5 rows):"
git cat-file blob "$ANCH:$F" | awk -F'|' '/^\*\*15\.1/{f=1;next} f&&/^\*\*15\.2/{f=0} f&&/^\| \*\*AC#/{gsub(/\\[|]/,"Q"); printf "   %-9s section-col=%.34s | grade-from-acceptor=%.26s | state-after-r2=%.40s\n", $2, $3, $4, $(NF-1)}'

echo "### V-F  the 15.2 attribution table item column, fixed ruler (6 rows)"
git cat-file blob "$ANCH:$F" | awk -F'|' '/^\*\*15\.2/{f=1;next} f&&/^\*\*15\.3/{f=0} f&&/^\| [0-9] /{gsub(/\\[|]/,"Q"); printf "   |%s| who-says=%.58s\n", $2, $3}'

echo "### V-G  does 15.1/15.2/15.3 carry the AC#3 and AC#4 discipline items into a reverse row? (union check)"
for n in 1 2 3 4 5; do
  printf '   AC#%s : 12.2=%s 15.2=%s 15.3=%s 15.1-fwd=%s 12.1-fwd=%s\n' "$n" \
   "$(git cat-file blob "$ANCH:$F" | awk 'index($0,"### B.12 ")==1{f=1;next} f&&/^### /{f=0} f&&/^\| [0-9]+ \|/' | grep -c "AC#$n")" \
   "$(git cat-file blob "$ANCH:$F" | awk '/^\*\*15\.2/{f=1;next} f&&/^\*\*15\.3/{f=0} f&&/^\| [0-9] /' | grep -c "AC#$n")" \
   "$(git cat-file blob "$ANCH:$F" | awk '/^\*\*15\.3/{f=1;next} f&&/^\*\*15\.4/{f=0} f&&/^\| [0-9]+ \|/' | grep -c "AC#$n")" \
   "$(git cat-file blob "$ANCH:$F" | awk '/^\*\*15\.1/{f=1;next} f&&/^\*\*15\.2/{f=0} f&&/^\| \*\*AC#/' | grep -c "AC#$n")" \
   "$(git cat-file blob "$ANCH:$F" | awk 'index($0,"### B.12 ")==1{f=1;next} f&&/^### /{f=0} f&&/^\| \*\*AC#/' | grep -c "AC#$n")"
done
echo "   discipline content by keyword (not by AC number) in the reverse union:"
for k in '翻勾' '不勾' 'done' '契约轴' '禁改' '零字节'; do
  printf '   %-7s 12.2=%s 15.2=%s 15.3=%s\n' "$k" \
   "$(git cat-file blob "$ANCH:$F" | awk 'index($0,"### B.12 ")==1{f=1;next} f&&/^### /{f=0} f&&/^\| [0-9]+ \|/' | grep -c "$k")" \
   "$(git cat-file blob "$ANCH:$F" | awk '/^\*\*15\.2/{f=1;next} f&&/^\*\*15\.3/{f=0} f&&/^\| [0-9] /' | grep -c "$k")" \
   "$(git cat-file blob "$ANCH:$F" | awk '/^\*\*15\.3/{f=1;next} f&&/^\*\*15\.4/{f=0} f&&/^\| [0-9]+ \|/' | grep -c "$k")"
done

echo "### V-H  anchor equality re-check (same fixed ruler on both anchors)"
for a in "$ANCH" "$NOW"; do
  printf '   %-8s 12.1 rows=%s 12.2 rows=%s 15.1 rows=%s 15.3 rows=%s | 12.2 未裁-cells=%s\n' "$a" \
   "$(git cat-file blob "$a:$F" | awk 'index($0,"### B.12 ")==1{f=1;next} f&&/^### /{f=0} f' | grep -c '^| \*\*AC#')" \
   "$(git cat-file blob "$a:$F" | awk 'index($0,"### B.12 ")==1{f=1;next} f&&/^### /{f=0} f' | grep -cE '^\| [0-9]+ \|')" \
   "$(git cat-file blob "$a:$F" | awk '/^\*\*15\.1/{f=1;next} f&&/^\*\*15\.2/{f=0} f' | grep -c '^| \*\*AC#')" \
   "$(git cat-file blob "$a:$F" | awk '/^\*\*15\.3/{f=1;next} f&&/^\*\*15\.4/{f=0} f' | grep -cE '^\| [0-9]+ \|')" \
   "$(git cat-file blob "$a:$F" | awk -F'|' 'index($0,"### B.12 ")==1{f=1;next} f&&/^### /{f=0} f&&/^\| \*\*AC#/{gsub(/\\[|]/,"Q"); if($(NF-1) ~ /未裁/) c++} END{print c+0}')"
done
