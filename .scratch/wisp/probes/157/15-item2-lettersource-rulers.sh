#!/usr/bin/env bash
# 票 157 实现程 · 第 2 笔的当轮复算尺（结论先不写在这里——这一笔判〔不成立〕，读数才是正文）
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9
E155=docs/evidence/s1/155-three-unjudged-cells-r1.md
A153=docs/evidence/s1/153-trace-lies-unguarded-r1-accept-r1.md
T153=.scratch/wisp/issues/153-the-compression-trace-has-two-unguarded-ways-to-lie-swapping-the-ran-guard-for-need-keeps-all-four-nails-green-and-the-trace-carries-no-taskid.md
A155=docs/evidence/s1/155-three-unjudged-cells-r1-accept-r1.md

echo "### U1. 155 证据件里'6 件'那一句的全部出现处（验收件只点了 :71）"
grep -n '6 件' $E155
echo
echo "### U2. 153 票面 AC#2 那一条（:25-:29）里可数的编号件"
sed -n '25,29p' $T153
echo "--- 票面里 (a) 到 (f) 那套字母的全部出现处（锚定到票面全文，逐枚带行号）---"
grep -n '6 件事' $T153
grep -n '(a)(d)(e)(f)\|(b) 问①\|(c) 问②' $T153
echo
echo "### U3. '票面拆 6 件事'那一行是什么时候进票面的（逐枚锚点现查）"
for r in 5365cb22 ff550f3 b319bab bb3a7a1 HEAD; do
  printf '%-9s 命中 = ' "$r"
  git show $r:$T153 2>/dev/null | grep -c '6 件事'
done
echo "--- b319bab 那一枚是什么（时刻＋作者＋首句）---"
git log -1 --format='%h %ad %an %s' --date=format:'%m-%d %H:%M' b319bab | cut -c1-90
echo "--- 票面最后一次改动 ---"
git log -1 --format='%h %ad %s' --date=format:'%m-%d %H:%M' -- $T153 | cut -c1-80
echo
echo "### U4. 153 验收件 §12 :827 那句原话（归属是谁写的）"
sed -n '827p' $A153
echo
echo "### U5. 验收件自己那一笔的原句（判它说的是什么）"
echo "--- 155 验收件 §4 缺什么表第 2 行 + 最小闭合② ---"
sed -n '172p' $A155
sed -n '177p' $A155
echo
echo "### U6. 工单第 2 笔那一行（本程据以复算的那句宣称）"
sed -n '14p' .scratch/wisp/issues/157-record-level-cleanup-ten-table-vs-disk-mismatches-left-by-tickets-155-and-153.md
echo
echo "### U7. 153 票面 AC#2 条款可数件数的锚定尺（本程那把'件'尺）"
printf 'AC#2 那条的交付物＋三问＋两枚 ⚠ 逐枚行号：\n'
awk 'NR>=25 && NR<=29 {printf ":%s | 交付物=%s 问=%s ⚠=%s\n", NR, (substr($0,1,2)=="- " && index($0,"AC#2")>0)?"1":"0", gsub(/[①②③]/,"&",$0), gsub(/⚠/,"&",$0)}' $T153
