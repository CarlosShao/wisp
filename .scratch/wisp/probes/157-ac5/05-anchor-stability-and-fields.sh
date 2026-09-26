#!/usr/bin/env bash
# 157 AC#5 recheck — probe 05: anchor-stability control, 档 attribution, duration-field scan, roster union.
# U-A  the same section-anchored ruler on TWO anchors (entry 4f6c14c vs current HEAD) => identical?
# U-B  forward/reverse row counts at every version between the landing commit and HEAD
# U-C  15.1 rows: dereference their pointed sections too
# U-D  criterion field "多久": duration tokens inside the closure cells (+ the yardstick table as control)
# U-E  roster union / overlap between 12.2 (5), 15.2 (6) and 15.3 (7)
# U-F  attribution of AC#3 / AC#4 (no reverse row) — does a non-implementer grade exist for them?
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9
F=docs/evidence/s1/155-three-unjudged-cells-r1.md
G=docs/evidence/s1/153-trace-lies-unguarded-r1.md
A=docs/evidence/s1/155-three-unjudged-cells-r1-accept-r1.md
ACC=docs/evidence/s1/157-record-level-cleanup-r1-accept-r1.md
T157=.scratch/wisp/issues/157-record-level-cleanup-ten-table-vs-disk-mismatches-left-by-tickets-155-and-153.md
ENTRY=4f6c14c
NOW=$(git rev-parse --short HEAD)

sec() { git cat-file blob "$1:$2" | awk -v p="$3" 'index($0,p)==1 {f=1; next} f && /^### /{f=0} f'; }
sub() { git cat-file blob "$1:$2" | awk -v s="$3" -v e="$4" 'index($0,s)==1 {f=1;next} f && index($0,e)==1 {f=0} f'; }

echo "### U-A  anchor stability: entry anchor $ENTRY vs current HEAD $NOW"
printf '   %s:%s = %s\n' "$ENTRY" "$F" "$(git rev-parse "$ENTRY:$F")"
printf '   %s:%s = %s\n' "$NOW"   "$F" "$(git rev-parse "$NOW:$F")"
printf '   entry is ancestor of HEAD  = %s\n' "$(git merge-base --is-ancestor "$ENTRY" HEAD && echo YES || echo NO)"
printf '   commits added since entry  = %s\n' "$(git rev-list --count "$ENTRY"..HEAD)"
echo "   the SAME section-anchored ruler, both anchors:"
for a in "$ENTRY" "$NOW"; do
  printf '   %-8s 12.1fwd=%s 12.2rev=%s 15.1fwd=%s 15.2attr=%s 15.3rev=%s\n' "$a" \
   "$(sec "$a" "$F" "### B.12 " | grep -c '^| \*\*AC#')" \
   "$(sec "$a" "$F" "### B.12 " | grep -cE '^\| [0-9]+ \|')" \
   "$(sub "$a" "$F" '**15.1' '**15.2' | grep -c '^| \*\*AC#')" \
   "$(sub "$a" "$F" '**15.2' '**15.3' | grep -cE '^\| [0-9] ')" \
   "$(sub "$a" "$F" '**15.3' '**15.4' | grep -cE '^\| [0-9]+ \|')"
done

echo "### U-B  row counts at each version from the landing commit up to HEAD"
for a in ef1c474f c97d706 b60c6b7 89a35e9e 67c0177c cc04355 1bd16d0 76f7cc9 35bcbff 942dae0 "$ENTRY" "$NOW"; do
  sha=$(git rev-parse --short "$a" 2>/dev/null) || { printf '   %-9s unresolvable\n' "$a"; continue; }
  printf '   %-9s %-8s 12.1=%s 12.2=%s 15.1=%s 15.3=%s 15.2=%s\n' "$a" "$sha" \
   "$(sec "$sha" "$F" "### B.12 " | grep -c '^| \*\*AC#')" \
   "$(sec "$sha" "$F" "### B.12 " | grep -cE '^\| [0-9]+ \|')" \
   "$(sub "$sha" "$F" '**15.1' '**15.2' | grep -c '^| \*\*AC#')" \
   "$(sub "$sha" "$F" '**15.3' '**15.4' | grep -cE '^\| [0-9]+ \|')" \
   "$(sub "$sha" "$F" '**15.2' '**15.3' | grep -cE '^\| [0-9] ')"
done

echo "### U-C  dereference the 15.1 rows (semantic level)"
git cat-file blob "$NOW:$F" | awk -F'|' '/^\*\*15\.1/{f=1;next} f&&/^\*\*15\.2/{f=0} f&&/^\| \*\*AC#/{gsub(/\\\|/,"SH"); printf "   %s pointer-col=%s\n", $2, substr($3,1,54)}'
printf '   pointed sections exist: 155 B.1-B.8=%s | 153 B.9/B.10=%s | 155 B.11=%s | 155 B.12=%s | 155 B.15=%s\n' \
 "$(git cat-file blob "$NOW:$F" | grep -cE '^### B\.[1-8] ')" "$(git cat-file blob "$NOW:$G" | grep -cE '^### B\.(9|10) ')" \
 "$(git cat-file blob "$NOW:$F" | grep -cE '^### B\.11 ')" "$(git cat-file blob "$NOW:$F" | grep -cE '^### B\.12 ')" \
 "$(git cat-file blob "$NOW:$F" | grep -cE '^### B\.15 ')"

echo "### U-D  the field the recheck ruler asks for: 谁 / 用什么尺 / 多久"
for tgt in 12.2 15.3 ctl; do
  case "$tgt" in
    12.2) src=$(sec "$NOW" "$F" "### B.12 " | grep -E '^\| [0-9]+ \|'); lbl="B.12 12.2" ;;
    15.3) src=$(sub "$NOW" "$F" '**15.3' '**15.4' | grep -E '^\| [0-9]+ \|'); lbl="B.15 15.3" ;;
    ctl)  src=$(awk 'NR>=265 && NR<=300' "$A" | grep -E '^\| [0-9]+ \|'); lbl="155-accept 7.2 (yardstick)" ;;
  esac
  printf '%s\n' "$src" | awk -F'|' -v lbl="$lbl" '
    {gsub(/\\\|/,"SH"); cl=$(NF-1); n++;
     if (cl ~ /分钟|小时|秒/) D++;
     if (cl ~ /一发|一次|本轮|下一张|重跑|重走|现读|并排|按|换一把|抽/) R++;
     if (cl ~ /不需|需要|批准/) C++;
     if (cl ~ /编排者|验收程|本程|任一枚|它/) W++}
    END{printf "   %-28s rows=%-3s who=%-3s ruler-or-action=%-3s cost=%-3s clock-duration=%s\n", lbl, n, W+0, R+0, C+0, D+0}'
done
printf '   ticket face AC#5 asks literally for: 最小闭合集合 (hits in the ticket=%s); the word 分钟/小时 appears in its AC#5 line? %s\n' \
 "$(grep -c '最小闭合集合' "$T157")" "$(grep -aE '^\- \[ \] \*\*AC#5' "$T157" | grep -cE '分钟|小时')"

echo "### U-E  roster union / overlap"
git cat-file blob "$NOW:$F" | awk -F'|' '/^\*\*15\.2/{f=1;next} f&&/^\*\*15\.3/{f=0} f&&/^\| [0-9] /{printf "   15.2 key: %-34s -> who:%s\n", substr($2,1,34), substr($3,1,44)}'
echo "   are 12.2 rows 3 and 4 re-listed in 15.3 or 15.2?"
for k in '153 的 AC#2' '四枚已定的勾' '翻勾'; do
  printf '   token "%s": 12.2=%s 15.2=%s 15.3=%s\n' "$k" \
   "$(sec "$NOW" "$F" "### B.12 " | grep -E '^\| [0-9]+ \|' | grep -c "$k")" \
   "$(sub "$NOW" "$F" '**15.2' '**15.3' | grep -c "$k")" \
   "$(sub "$NOW" "$F" '**15.3' '**15.4' | grep -E '^\| [0-9]+ \|' | grep -c "$k")"
done

echo "### U-F  AC#3 / AC#4 grades exist on the non-implementer side?"
grep -aE '^\| \*\*AC#[1-5]\*\*' "$ACC" | cut -c1-100 | sed 's/^/   157-accept 3.1: /'
