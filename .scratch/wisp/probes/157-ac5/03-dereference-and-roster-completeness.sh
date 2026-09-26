#!/usr/bin/env bash
# 157 AC#5 recheck — probe 03: semantic-level dereference + reverse-roster completeness.
# S-A  every 凭据 cited inside B.12 / B.15 (probe file + section marker) — does the target exist at MY anchor?
# S-B  forward 档 column: is "未裁" literally present for each of the 5 boxes (12.1) and each grade in 15.1?
# S-C  cross-pointer: which forward rows point into the reverse table by number ("见 12.2 第 N 行")
# S-D  reverse roster printed in full (no truncation) for the overlap analysis B.12-12.2 vs B.15-15.3 vs 15.2
# S-E  ticket-box coverage: each of the 5 boxes + 3 discipline contents — named anywhere in the two reverse tables?
# S-F  criterion-text count (票面 "4 枚框" vs 现量 5) on both the ticket face and the two dispatches
# S-G  worktree-vs-anchor guard for every file this program reads (never adjudicate on a dirty tree)
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9
F=docs/evidence/s1/155-three-unjudged-cells-r1.md
G=docs/evidence/s1/153-trace-lies-unguarded-r1.md
A=docs/evidence/s1/155-three-unjudged-cells-r1-accept-r1.md
ACC=docs/evidence/s1/157-record-level-cleanup-r1-accept-r1.md
D1=.scratch/wisp/dispatches/2026-09-26-131x-impl-157-record-level.md
D2=.scratch/wisp/dispatches/2026-09-26-132x-accept-157-r1.md
D3=.scratch/wisp/dispatches/2026-09-26-134x-impl-157-r2-one-cell.md

echo "### S-A  cited 凭据 targets exist at the anchor (blob object type + cited marker present)"
for p in .scratch/wisp/probes/157/20-delivery-shape-and-forbidden-faces.txt \
         .scratch/wisp/probes/157/22-r2-ac5-cell-and-b1-columns.txt \
         .scratch/wisp/probes/157/23-r2-close-shape.txt \
         .scratch/wisp/probes/157/15-item2-lettersource-rulers.txt \
         .scratch/wisp/probes/157/21-final-self-check.txt ; do
  t=$(git cat-file -t "4f6c14c:$p" 2>&1)
  printf '  %-62s git cat-file -t => %s' "$p" "$t"
  if [ "$t" = "blob" ]; then
    printf '  lines=%s' "$(git cat-file blob "4f6c14c:$p" | wc -l)"
    for m in V1 V2 V3 V4 V5 V6 V7 W1 W2 W3 W5 W6 W8 W10 W11 W12 W13 K0 K1 K2 K3 K4 K5 K6 K7; do
      c=$(git cat-file blob "4f6c14c:$p" | grep -cE "^(###|##)? ?$m[ .．、:]")
      [ "$c" -gt 0 ] && printf ' %s=%s' "$m" "$c"
    done
  fi
  echo
done

echo "### S-B  forward grade column, per row (12.1 = implementer self-proof, 15.1 = grades from the acceptor)"
echo "-- 12.1 rows: does the last cell literally contain 未裁 ?"
git cat-file blob 4f6c14c:$F | awk -F'|' 'index($0,"### B.12 ")==1{f=1;next} f&&/^### /{f=0} f&&/^\| \*\*AC#/{gsub(/\\\|/,"SH"); cell=$(NF-1); printf "   %s -> 未裁?%s  col=%s\n", $2, (cell ~ /未裁/)?"YES":"NO ", substr(cell,1,26)}'
echo "-- 15.1 rows: does the grade column carry a NON-implementer 档 ?"
git cat-file blob 4f6c14c:$F | awk -F'|' '/^\*\*15\.1/{f=1;next} f&&/^\*\*15\.2/{f=0} f&&/^\| \*\*AC#/{gsub(/\\\|/,"SH"); printf "   %s -> col3(验收程的档)=%s\n", $2, substr($4,1,40)}'
echo "-- POSITIVE CONTROL for S-B: 155-accept 7.1 grade column (must read 成立 for its 4 boxes)"
awk -F'|' 'NR>=256 && NR<=264 && /^\| \*\*AC#/{gsub(/\\\|/,"SH"); printf "   ctl %s -> %s\n", $2, substr($(NF-1),1,30)}' "$A"

echo "### S-C  cross-pointer from the forward table into the reverse table"
printf '  rows in 12.1 whose grade cell names "见 12.2" = '
git cat-file blob 4f6c14c:$F | awk 'index($0,"### B.12 ")==1{f=1;next} f&&/^### /{f=0} f&&/^\| \*\*AC#/' | grep -c '见 12.2'
printf '  rows in 15.1 whose state cell names "15.3"    = '
git cat-file blob 4f6c14c:$F | awk '/^\*\*15\.1/{f=1;next} f&&/^\*\*15\.2/{f=0} f&&/^\| \*\*AC#/' | grep -c '15.3'

echo "### S-D  reverse roster printed in full (3 tables, NO truncation)"
echo "-- 12.2 item column (5 rows):"
git cat-file blob 4f6c14c:$F | awk -F'|' 'index($0,"### B.12 ")==1{f=1;next} f&&/^### /{f=0} f&&/^\| [0-9]+ \|/{gsub(/\\\|/,"SH"); printf "   |%s| %s\n", $2, substr($3,1,88)}'
echo "-- 15.2 attribution table (6 rows):"
git cat-file blob 4f6c14c:$F | awk -F'|' '/^\*\*15\.2/{f=1;next} f&&/^\*\*15\.3/{f=0} f&&/^\| [0-9] /{gsub(/\\\|/,"SH"); printf "   |%s| %s\n", $2, substr($3,1,80)}'
echo "-- 15.3 item column (7 rows):"
git cat-file blob 4f6c14c:$F | awk -F'|' '/^\*\*15\.3/{f=1;next} f&&/^\*\*15\.4/{f=0} f&&/^\| [0-9]+ \|/{gsub(/\\\|/,"SH"); printf "   |%s| %s\n", $2, substr($3,1,88)}'

echo "### S-E  coverage: is each ticket box named in the reverse union (12.2 + 15.3 + 15.2)?"
for n in 1 2 3 4 5; do
  printf '  AC#%s named in: 12.2=%s 15.3=%s 15.2=%s\n' "$n" \
   "$(git cat-file blob 4f6c14c:$F | awk 'index($0,"### B.12 ")==1{f=1;next} f&&/^### /{f=0} f&&/^\| [0-9]+ \|/' | grep -c "AC#$n")" \
   "$(git cat-file blob 4f6c14c:$F | awk '/^\*\*15\.3/{f=1;next} f&&/^\*\*15\.4/{f=0} f&&/^\| [0-9]+ \|/' | grep -c "AC#$n")" \
   "$(git cat-file blob 4f6c14c:$F | awk '/^\*\*15\.2/{f=1;next} f&&/^\*\*15\.3/{f=0} f&&/^\| [0-9] /' | grep -c "AC#$n")"
done
echo "  discipline contents (AC#3 不勾框 / -done, AC#4 契约轴/禁改面) anywhere in the appendix:"
printf '  翻勾|勾 anywhere in 12.2=%s 15.3=%s | -done in 12.1=%s 15.1=%s | 禁改面|零字节|契约轴 in 12.1=%s 15.1=%s\n' \
 "$(git cat-file blob 4f6c14c:$F | awk 'index($0,"### B.12 ")==1{f=1;next} f&&/^### /{f=0} f&&/^\| [0-9]+ \|/' | grep -cE '翻勾|不勾') " \
 "$(git cat-file blob 4f6c14c:$F | awk '/^\*\*15\.3/{f=1;next} f&&/^\*\*15\.4/{f=0} f&&/^\| [0-9]+ \|/' | grep -cE '翻勾|不勾')" \
 "$(git cat-file blob 4f6c14c:$F | awk 'index($0,"### B.12 ")==1{f=1;next} f&&/^### /{f=0} f&&/^\| \*\*AC#/' | grep -c 'done') " \
 "$(git cat-file blob 4f6c14c:$F | awk '/^\*\*15\.1/{f=1;next} f&&/^\*\*15\.2/{f=0} f&&/^\| \*\*AC#/' | grep -c 'done')" \
 "$(git cat-file blob 4f6c14c:$F | awk 'index($0,"### B.12 ")==1{f=1;next} f&&/^### /{f=0} f&&/^\| \*\*AC#/' | grep -cE '禁改|零字节|契约轴') " \
 "$(git cat-file blob 4f6c14c:$F | awk '/^\*\*15\.1/{f=1;next} f&&/^\*\*15\.2/{f=0} f&&/^\| \*\*AC#/' | grep -cE '禁改|零字节|契约轴')"

echo "### S-F  the '4 枚框' text vs the disk count (three sources side by side)"
printf '  ticket face 157 : grep -c "^- [ ] **AC#"   = %s\n' "$(grep -c '^- \[ \] \*\*AC#' .scratch/wisp/issues/157-record-level-cleanup-ten-table-vs-disk-mismatches-left-by-tickets-155-and-153.md)"
printf '  ticket face 157 : AC#5 line contains "4 枚框" = %s\n' "$(grep -c '4 枚框' .scratch/wisp/issues/157-record-level-cleanup-ten-table-vs-disk-mismatches-left-by-tickets-155-and-153.md)"
printf '  impl dispatch D1 :58 "5 枚框" hits            = %s\n' "$(grep -c '5 枚框' "$D1")"
printf '  r2 dispatch D3   "5 枚框" hits                = %s\n' "$(grep -c '5 枚框' "$D3")"
printf '  B.12 heading itself says "5 枚框"             = %s\n' "$(git cat-file blob 4f6c14c:$F | grep -c '^### B\.12 票面 5 枚框')"
echo "  ticket 155's face box count (the source of the '4 枚框' wording in ticket 157 AC#5):"
printf '  ticket 155: matches of ^155 in issues/ = %s file(s) (only one, so no truncation risk)\n' "$(ls .scratch/wisp/issues | grep -c '^155')"
T155=$(ls .scratch/wisp/issues | grep '^155' | head -1)
printf '  ticket 155 file = %s\n' "$T155"
printf '  ticket 155 anchored box count (any state) = %s ; unchecked = %s  <= the origin of the "四枚框" wording\n' \
 "$(grep -cE '^- \[[ x]\] \*\*AC#' ".scratch/wisp/issues/$T155")" \
 "$(grep -c '^- \[ \] \*\*AC#' ".scratch/wisp/issues/$T155")"

echo "### S-G  worktree vs anchor guard for every file this program reads (dirty tree is never an anchor)"
for p in "$F" "$G" "$A" "$ACC" "$D1" "$D2" "$D3" .scratch/wisp/issues/157-record-level-cleanup-ten-table-vs-disk-mismatches-left-by-tickets-155-and-153.md; do
  st=$(git status --porcelain -- "$p")
  if [ -z "$st" ]; then same=$(cmp -s <(git cat-file blob "4f6c14c:$p") "$p" && echo IDENTICAL || echo DIFFERS)
  else same="worktree dirty => not used as anchor"; fi
  printf '  %-70s status=[%s] %s\n' "$(basename "$p")" "$st" "$same"
done

echo "### S-H  AC#5 row of the acceptor verdict, verbatim in full (the text being re-judged)"
sed -n '170p' "$ACC"
echo "### S-I  the acceptor's prescribed closure, verbatim in full"
sed -n '178p' "$ACC"
echo "### S-J  does the delivered table satisfy that prescription line by line? (4 sub-asks)"
printf '  prescription says: 形状照 155-accept 7.1/7.2 两表 -> 12.1 rows=%s cols=%s | 12.2 rows=%s cols=%s\n' \
 "$(git cat-file blob 4f6c14c:$F | awk 'index($0,"### B.12 ")==1{f=1;next} f&&/^### /{f=0} f' | grep -c '^| \*\*AC#')" \
 "$(git cat-file blob 4f6c14c:$F | awk 'index($0,"### B.12 ")==1{f=1;next} f&&/^### /{f=0} f' | grep -m1 '^| 票面那枚框' | awk -F'|' '{print NF-2}')" \
 "$(git cat-file blob 4f6c14c:$F | awk 'index($0,"### B.12 ")==1{f=1;next} f&&/^### /{f=0} f' | grep -cE '^\| [0-9]+ \|')" \
 "$(git cat-file blob 4f6c14c:$F | awk 'index($0,"### B.12 ")==1{f=1;next} f&&/^### /{f=0} f' | grep -m1 '^| # ' | awk -F'|' '{print NF-2}')"
printf '  prescription says: 含"未裁"清单 -> 12.1 grade cells with 未裁 = %s / 5\n' \
 "$(git cat-file blob 4f6c14c:$F | awk -F'|' 'index($0,"### B.12 ")==1{f=1;next} f&&/^### /{f=0} f&&/^\| \*\*AC#/{gsub(/\\\|/,"SH"); if($(NF-1) ~ /未裁/) c++} END{print c+0}')"
printf '  prescription says: ＋最小闭合 -> 12.2 rows with a closure cell >=40 bytes = %s / 5\n' \
 "$(git cat-file blob 4f6c14c:$F | awk -F'|' 'index($0,"### B.12 ")==1{f=1;next} f&&/^### /{f=0} f&&/^\| [0-9]+ \|/{gsub(/\\\|/,"SH"); if(length($(NF-1))>=40) c++} END{print c+0}')"
printf '  prescription: "157 自己往任一节的末尾补" -> the table lives at the END of the appendix? B.12 line no vs appendix-B start line no\n'
git cat-file blob 4f6c14c:$F | grep -n '^## 附录 B\|^### B\.12 \|^### B\.15 ' | cut -c1-46 | sed 's/^/    /'
