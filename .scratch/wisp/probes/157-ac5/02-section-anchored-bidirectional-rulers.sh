#!/usr/bin/env bash
# 157 AC#5 recheck — probe 02: section-anchored bidirectional rulers (replaces the line-window V6 ruler).
# R-A  section slicer by heading (version-proof)  /  R-B forward+reverse row counts, line-anchored
# R-C  dereference every forward pointer ("which section delivered that box") and verify the target exists
# R-D  per-row field presence in the reverse tables (who / which ruler / cost marker)
# R-E  empty-shell scan (待补 / 以后再 / 回头补 / closure cell under 10 bytes)
# R-F  shape comparison against the two tables the acceptor itself prescribed (155-accept 7.1 / 7.2)
# R-G  the verdict under re-judgement, verbatim  /  R-H appendix-wide 未裁 per version  /  R-I landing commit
# read-only, no code, no test runners
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9
F=docs/evidence/s1/155-three-unjudged-cells-r1.md
G=docs/evidence/s1/153-trace-lies-unguarded-r1.md
A=docs/evidence/s1/155-three-unjudged-cells-r1-accept-r1.md
ACC=docs/evidence/s1/157-record-level-cleanup-r1-accept-r1.md

sec() { # $1=anchor $2=path $3=heading prefix (exact, incl. trailing space)
  git cat-file blob "$1:$2" | awk -v p="$3" 'index($0,p)==1 {f=1; next} f && /^### /{f=0} f'
}
sub() { # $1=anchor $2=path $3=start bold marker $4=stop bold marker
  git cat-file blob "$1:$2" | awk -v s="$3" -v e="$4" 'index($0,s)==1 {f=1;next} f && index($0,e)==1 {f=0} f'
}

echo "### R-B  row counts per version, both directions, section-anchored"
for V in ec815820 89a35e9e 4f6c14c; do
  echo "=== VERSION $V (155 file: $(git cat-file blob "$V:$F" | wc -l) lines) ==="
  printf '     B.12 forward rows (^| **AC#)   = '; sec "$V" "$F" "### B.12 " | grep -c '^| \*\*AC#'
  printf '     B.12 reverse rows (^| N |)     = '; sec "$V" "$F" "### B.12 " | grep -cE '^\| [0-9]+ \|'
  for t in 对账 双向 未裁 最小闭合; do
    printf '       B.12 token %-9s = ' "$t"; sec "$V" "$F" "### B.12 " | grep -c "$t"
  done
  printf '     15.1 forward rows              = '; sub "$V" "$F" '**15.1' '**15.2' | grep -c '^| \*\*AC#'
  printf '     15.2 attribution rows          = '; sub "$V" "$F" '**15.2' '**15.3' | grep -cE '^\| [0-9] '
  printf '     15.3 reverse rows              = '; sub "$V" "$F" '**15.3' '**15.4' | grep -cE '^\| [0-9]+ \|'
  for t in 对账 双向 未裁 最小闭合; do
    printf '       B.15 token %-9s = ' "$t"; sec "$V" "$F" "### B.15 " | grep -c "$t"
  done
  echo
done

echo "### R-C  dereference B.12 forward pointers at the anchor (semantic level)"
git cat-file blob 4f6c14c:$F | awk 'index($0,"### B.12 ")==1{f=1;next} f&&/^### /{f=0} f&&/^\| \*\*AC#/{print}' | while IFS= read -r row; do
  box=$(printf '%s' "$row" | sed -E 's/^\| (\*\*AC#[0-9]\*\*).*/\1/')
  ptr=$(printf '%s' "$row" | awk -F'|' '{print $3}' | cut -c1-46)
  grade=$(printf '%s' "$row" | awk -F'|' '{print $(NF-1)}' | cut -c1-22)
  printf '  %-14s pointer=%-46s grade-col=%s\n' "$box" "$ptr" "$grade"
done
echo "  -- does each pointed target exist at the anchor? (line-anchored ruler)"
printf '     B.1-B.8 sections (AC#1)         = '; git cat-file blob 4f6c14c:$F | grep -cE '^### B\.[1-8] '
printf '     B.9/B.10 in 153 file (AC#2)     = '; git cat-file blob 4f6c14c:$G | grep -cE '^### B\.(9|10) '
printf '     153 file appendix-B heading     = '; git cat-file blob 4f6c14c:$G | grep -c '^## 附录 B'
printf '     B.11 section (AC#4)             = '; git cat-file blob 4f6c14c:$F | grep -cE '^### B\.11 '
printf '     B.12 section (AC#5)             = '; git cat-file blob 4f6c14c:$F | grep -cE '^### B\.12 '
printf '     B.15 section (AC#5 second half) = '; git cat-file blob 4f6c14c:$F | grep -cE '^### B\.15 '

echo "### R-D  per-row closure-field presence (awk -F'|' taking the LAST non-empty cell)"
echo "-- R-D0 POSITIVE CONTROL: same three greps over the 7 last-cells of 155-accept 7.2 --"
awk -F'|' 'NR>=265 && NR<=280 && /^\| [0-9]+ \|/{gsub(/\\\|/,"SH"); print $(NF-1)}' "$A" | while IFS= read -r cl; do
  printf '   ctl closure-bytes=%-5s who=%s ruler=%s cost=%s\n' "$(printf '%s' "$cl" | wc -c)" \
    "$(printf '%s' "$cl" | grep -cE '编排者|验收程|本程|任一枚|它')" \
    "$(printf '%s' "$cl" | grep -cE '尺|一发|按|重跑|并排|抽|名册|现读|show|grep|sed|换')" \
    "$(printf '%s' "$cl" | grep -cE '不需|需要|批准|一发')"
done
echo "-- B.12 12.2 rows:"
git cat-file blob 4f6c14c:$F | awk 'index($0,"### B.12 ")==1{f=1;next} f&&/^### /{f=0} f&&/^\| [0-9]+ \|/' | while IFS= read -r row; do
  n=$(printf '%s' "$row" | awk -F'|' '{print $2}' | tr -d ' ')
  cl=$(printf '%s' "$row" | sed 's/\\|/SH/g' | awk -F'|' '{print $(NF-1)}')
  printf '   row %-2s closure-bytes=%-5s who=%s ruler=%s cost=%s\n' "$n" "$(printf '%s' "$cl" | wc -c)" \
    "$(printf '%s' "$cl" | grep -cE '编排者|验收程|本程|任一枚')" \
    "$(printf '%s' "$cl" | grep -cE '尺|一发|按|重跑|并排|抽|名册|现读')" \
    "$(printf '%s' "$cl" | grep -cE '不需|需要新仪器|批准')"
done
echo "-- B.15 15.3 rows:"
git cat-file blob 4f6c14c:$F | awk '/^\*\*15\.3/{f=1;next} f&&/^\*\*15\.4/{f=0} f&&/^\| [0-9]+ \|/' | while IFS= read -r row; do
  n=$(printf '%s' "$row" | awk -F'|' '{print $2}' | tr -d ' ')
  cl=$(printf '%s' "$row" | sed 's/\\|/SH/g' | awk -F'|' '{print $(NF-1)}')
  printf '   row %-2s closure-bytes=%-5s who=%s ruler=%s cost=%s\n' "$n" "$(printf '%s' "$cl" | wc -c)" \
    "$(printf '%s' "$cl" | grep -cE '编排者|验收程|本程|任一枚|它')" \
    "$(printf '%s' "$cl" | grep -cE '尺|一发|按|重跑|换|并排|名册|现读|摘掉')" \
    "$(printf '%s' "$cl" | grep -cE '不需|需要|批准|一发')"
done

echo "### R-E  empty-shell scan"
for h in "### B.12 " "### B.15 "; do
  printf '  %s 待补/以后再/回头补/以后加固 hits = ' "$h"; sec 4f6c14c "$F" "$h" | grep -cE '待补|以后再|回头补|以后加固'
done
echo "  rows whose last content cell is 1..9 bytes (empty-shell candidates):"
git cat-file blob 4f6c14c:$F | awk -F'|' 'index($0,"### B.12 ")==1{f=1;next} f&&/^### /{f=0} f&&/^\| /{gsub(/\\\|/,"SH"); c=$(NF-1); gsub(/[ \-|]/,"",c); if(length(c)>0 && length(c)<10) print "   B12 weak: " substr($0,1,70)}'
git cat-file blob 4f6c14c:$F | awk -F'|' '/^\*\*15\.3/{f=1;next} f&&/^\*\*15\.4/{f=0} f&&/^\| /{gsub(/\\\|/,"SH"); c=$(NF-1); gsub(/[ \-|]/,"",c); if(length(c)>0 && length(c)<10) print "   B15 weak: " substr($0,1,70)}'
echo "  POSITIVE CONTROL for the last-cell ruler (same checker on 155-accept 7.2, 7 rows):"
awk 'NR>=265 && NR<=280 && /^\| [0-9]+ \|/{gsub(/\\\|/,"SH"); c=$(NF-1); gsub(/^ +| +$/,"",c); printf "   ctl row %s bytes=%s starts=%s\n", $2, length(c), substr(c,1,28)}' FS='|' "$A"

echo "### R-F  shape yardstick the acceptor itself prescribed (155-accept 7.1/7.2) vs what landed"
printf '  155-accept 7.1 forward rows   = '; sed -n '256,264p' "$A" | grep -c '^| \*\*AC#'
printf '  155-accept 7.2 reverse rows   = '; sed -n '265,300p' "$A" | grep -cE '^\| [0-9]+ \|'
printf '  155-accept 7.1 header cols    = '; sed -n '258p' "$A" | awk -F'|' '{print NF-2}'
printf '  155-accept 7.2 header cols    = '; sed -n '267p' "$A" | awk -F'|' '{print NF-2}'
printf '  B.12 12.1 forward rows        = '; sec 4f6c14c "$F" "### B.12 " | grep -c '^| \*\*AC#'
printf '  B.12 12.2 reverse rows        = '; sec 4f6c14c "$F" "### B.12 " | grep -cE '^\| [0-9]+ \|'
printf '  B.12 12.1 header cols         = '; sec 4f6c14c "$F" "### B.12 " | grep -m1 '^| 票面那枚框' | awk -F'|' '{print NF-2}'
printf '  B.12 12.2 header cols         = '; sec 4f6c14c "$F" "### B.12 " | grep -m1 '^| # ' | awk -F'|' '{print NF-2}'
printf '  B.15 15.1 forward rows        = '; sub 4f6c14c "$F" '**15.1' '**15.2' | grep -c '^| \*\*AC#'
printf '  B.15 15.3 reverse rows        = '; sub 4f6c14c "$F" '**15.3' '**15.4' | grep -cE '^\| [0-9]+ \|'

echo "### R-G  the verdict under re-judgement, verbatim (157-accept 3.1 AC#5 row + 3.2 row 3)"
grep -n 'AC#5' "$ACC" | cut -c1-240

echo "### R-H  appendix-wide and per-section 未裁 counts per version (section ruler)"
for a in ec815820 89a35e9e 4f6c14c; do
  printf '  %-9s appendix-B 未裁=%-4s B.12 内=%-4s B.15 内=%s\n' "$a" \
   "$(git cat-file blob "$a:$F" | awk 'index($0,"## 附录 B")>=1{f=1} f' | grep -c '未裁')" \
   "$(sec "$a" "$F" "### B.12 " | grep -c '未裁')" \
   "$(sec "$a" "$F" "### B.15 " | grep -c '未裁')"
done

echo "### R-I  landing commit of each heading (pickaxe, path-scoped) + the roster of that commit"
for pat in '### B.12 ' '### B.15 '; do
  echo "  literal search string: '$pat'  (git log -S is literal by default, so do NOT backslash-escape the dot)"
  echo "    full pickaxe list (path-scoped, NOT truncated; first line = the commit that introduced it):"
  git log --reverse -S"$pat" --pretty='    %h %ad %s' --date=format:'%H:%M:%S' -- "$F" | cut -c1-96
  echo "    same string via --pickaxe-regex (regex ruler, must agree on the first line):"
  git log --reverse --pickaxe-regex -S"^$pat" --pretty='    %h %ad' --date=format:'%H:%M:%S' -- "$F" | cut -c1-40
done
git log -1 --no-walk --pretty='  ef1c474 author-date=%ad' --date=format:'%H:%M:%S' ef1c474f
git show --name-only --pretty=tformat: ef1c474f | sed 's/^/    ef1c474 touched: /'
git show --name-only --pretty=tformat: 1bd16d0 | sed 's/^/    1bd16d0 touched: /'
