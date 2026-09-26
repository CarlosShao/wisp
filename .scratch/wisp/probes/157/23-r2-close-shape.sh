#!/usr/bin/env bash
# 票 157 续程 r2 · 收口形状自证（只读尺；名册一律 --no-walk ＋ 锚定 grep，禁挂 head）
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9

F155=docs/evidence/s1/155-three-unjudged-cells-r1.md
F153=docs/evidence/s1/153-trace-lies-unguarded-r1.md
T157=.scratch/wisp/issues/157-record-level-cleanup-ten-table-vs-disk-mismatches-left-by-tickets-155-and-153.md
T155=.scratch/wisp/issues/155-ticket-153-leaves-three-cells-unjudged-because-my-attack-point-grid-and-the-ticket-ac-grid-are-not-the-same-coordinate-system-done.md
BASE=67c0177                      # 本程进场基版（派单存档那一枚）
T153=.scratch/wisp/issues/153-the-compression-trace-has-two-unguarded-ways-to-lie-swapping-the-ran-guard-for-need-keeps-all-four-nails-green-and-the-trace-carries-no-taskid.md
R2=$(git log --pretty=%h --grep='evidence(157 续程 r2' | tr '\n' ' ')
R2BAD=$(git log --pretty=%h --grep='157 续程 r2' | tr '\n' ' ')

echo "### K0. 时刻 ＋ 进场基版 ＋ 本程名下 commit（两把锚定尺并排）"
date '+%Y-%m-%d %H:%M:%S %z'
echo "  基版 $BASE ＝$(git cat-file -t $BASE)｜HEAD＝$(git rev-parse --short HEAD)"
echo "  ⚠ 尺坏一枚（本程自己踩的）：锚定 grep 的锚得锚到**本程自己的主题前缀**——"
echo "     --grep='157 续程 r2' ＝ $(echo $R2BAD | wc -w) 枚（把编排者 $BASE 那枚派单存档也算进来了）"
echo "     --grep='evidence(157 续程 r2' ＝ $(echo $R2 | wc -w) 枚（本程名下）"
git log --no-walk --pretty='  %h %ad %s' --date=format:%H:%M:%S $R2 | cut -c1-64
echo
echo "### K1. 名册（--no-walk ＋ 枚枚 --name-only 并集）"
echo '$ git log --no-walk --pretty=tformat: --name-only <枚枚> | sort -u'
git log --no-walk --pretty=tformat: --name-only $R2 | sort -u > /d/tmp/157r2/roster.txt
cat /d/tmp/157r2/roster.txt | sed 's/^/  /'
echo "  名册枚数＝$(grep -c . /d/tmp/157r2/roster.txt)"
echo '  越界（名册减掉声明的两族写面）＝'
grep -vE '^docs/evidence/s1/155-three-unjudged-cells-r1\.md$|^\.scratch/wisp/probes/157/' /d/tmp/157r2/roster.txt; echo "    grep rc=$?"
echo '  坑①复现：同一把尺不带 --no-walk ＝'
git log --pretty=tformat: --name-only $R2 | sort -u | wc -l
echo
echo "### K2. 禁改面逐支（19 支＝B.11 那 16 支 ＋ 本程另加 3 支；尺＝同一枚名册）"
for pat in '^docs/PLAN\.md' '^docs/specs/' '^internal/' '^cmd/' '^third_party/' 'thresholds\.go' 'golden' 'allowlist\.txt' '^scripts/slo-check\.ps1' '^tools/d22scan/' '^docs/reports/pending-and-issues\.md' '^docs/reports/HANDOVER\.md' '^docs/reports/injection-timeline\.md' '^\.scratch/wisp/issues/' '^docs/evidence/s1/152-' '^docs/evidence/s1/154-' '^docs/evidence/s1/156-' '^frontend/' '^design/'; do
  c=$(grep -cE "$pat" /d/tmp/157r2/roster.txt)
  printf '  %-46s %s 枚\n' "$pat" "$c"
done
echo '  正控（同一把名册尺）：git show --name-only b23c7f7 上 ^internal/ 命中＝'
git show --pretty=tformat: --name-only b23c7f7 | grep -cE '^internal/'
echo '  另一把尺（对象不同，别互抄）：git ls-tree -r b23c7f7 -- internal/ ＝'
git ls-tree -r --name-only b23c7f7 -- internal/ | grep -cE '^internal/'
echo '  别程在飞的四族（152/154/156 的活、cmd/wisp WIP、design/frontend）在本程名册命中＝'
grep -cE '152-|154-|156-|^cmd/wisp/|^design/|^frontend/' /d/tmp/157r2/roster.txt
echo
echo "### K3. 只追加自证（旧行不抹）"
echo "  本件 diff 基版..HEAD 的增删＝"
git diff --numstat $BASE HEAD -- $F155 | sed 's/^/    /'
echo "  §0–§11（前 385 行）与 r1 基版 424ee363 比＝"
cmp <(git cat-file blob 424ee363:$F155 | head -385) <(git cat-file blob HEAD:$F155 | head -385) && echo '    静默（一字未改）'
echo "  B.1 那一发＝纯插入：基版第 443 行起 == 本版第 445 行起同长度＝"
BL=$(git cat-file blob $BASE:$F155 | wc -l); TAILN=$((BL - 442))
cmp <(git cat-file blob $BASE:$F155 | tail -n +443) <(git cat-file blob HEAD:$F155 | tail -n +445 | head -n $TAILN) && echo "    静默（尾部 $TAILN 行逐字节相同；其后全是 B.15 的新增）"
echo "  153 件与基版比＝"
cmp <(git cat-file blob $BASE:$F153) <(git cat-file blob HEAD:$F153) && echo '    静默（本程一字未碰）'
echo "  附录 B 标题枚数＝$(grep -c '^## 附录 B（票 157 代记，只追加）' $F155)／行数＝$(wc -l < $F155)"
echo "  新增节标题：$(grep -nE '^### B\.1[0-9] ' $F155 | tr '\n' ' ')"
echo
echo "### K4. AC#5 那一格现在的形状（正向 5 行／反向 5 行／本程那一节 7 行）"
echo "  本件 '对账'=$(grep -c '对账' $F155) '双向'=$(grep -c '双向' $F155)"
echo "  B.12 正向那表行数=$(awk '/^### B\.12 /,/^### B\.13 /' $F155 | grep -cE '^\| \*\*AC#')  B.12 反向行数=$(awk '/^\*\*12\.2/,/^\*\*15\.0|^### B\.13/' $F155 | grep -cE '^\| [0-9] \|')"
echo "  B.15 正向行数=$(awk '/^\*\*15\.1/,/^\*\*15\.2/' $F155 | grep -cE '^\| \*\*AC#')  B.15 反向行数=$(awk '/^\*\*15\.3/,/^\*\*15\.4/' $F155 | grep -cE '^\| [0-9] \|')"
echo
echo "### K5. 票面：三枚票本程一字未动（锚定到票路径的最后一次改动）"
for t in $T157 $T155 $T153; do printf '  %-70s 最后改动＝%s\n' "$(basename $t)" "$(git log -1 --pretty='%h %ad' --date=format:%H:%M -- $t)"; done
echo "  157 未勾=$(grep -c '^- \[ \]' $T157) 已勾=$(grep -c '^- \[x\]' $T157)"
echo "  issues 里含 157 的 -done 文件＝$(ls .scratch/wisp/issues/ | grep -c '157.*-done')"
echo
echo "### K6. probes/157 枚数与名册（本程新增 22／23 两组）"
ls -1 .scratch/wisp/probes/157/ | wc -l
ls -1 .scratch/wisp/probes/157/ | grep -E '^2[23]-' | sed 's/^/  /'
echo
echo "### K7. push／remote／工作树边界（本程只 commit 不 push）"
echo '  本机 remote（**早已存在，非本程所加**）＝'; git remote -v | sed 's/^/    /'
echo '  上游＝'; git rev-parse --abbrev-ref --symbolic-full-name '@{upstream}' 2>&1 | sed 's/^/    /'
echo '  上游..HEAD 领先枚数（含别程与编排者的，不是本程一家的）＝'; git rev-list --count '@{upstream}..HEAD' | sed 's/^/    /'
echo '  本程名下四枚是否已被推走（同一发尺）＝'
git log --pretty=%h --grep='evidence(157 续程 r2' '@{upstream}..HEAD' | wc -l | sed 's/^/    /'
echo '  未提交别程半件仍在原状（本程未 add／未还原）＝'
git status --porcelain -- docs/evidence/s1/152-subject-death-never-measured-r1-accept-r1.md .scratch/wisp/probes/152/my152.py | sed 's/^/    /'
echo "### END"
