#!/usr/bin/env bash
# 票 157 续程 r2 · 本程当轮复算（只读尺；取数一律打在 git 对象上，不读脏工作树）
# 两件事：(1) AC#5 那枚"票面 5 枚框 ↔ 本程格"双向对账——盘上到底有没有；
#         (2) B.1 那张全档表三行的"档"列到底被读成了什么（两把尺并排）。
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9

F155=docs/evidence/s1/155-three-unjudged-cells-r1.md
F153=docs/evidence/s1/153-trace-lies-unguarded-r1.md
A153=docs/evidence/s1/153-trace-lies-unguarded-r1-accept-r1.md
ACC=docs/evidence/s1/157-record-level-cleanup-r1-accept-r1.md
T157=.scratch/wisp/issues/157-record-level-cleanup-ten-table-vs-disk-mismatches-left-by-tickets-155-and-153.md
T153=.scratch/wisp/issues/153-the-compression-trace-has-two-unguarded-ways-to-lie-swapping-the-ran-guard-for-need-keeps-all-four-nails-green-and-the-trace-carries-no-taskid.md
PROBE11=.scratch/wisp/probes/157/11-items-1-2-3-8-table-and-range.sh

echo "### W0. 时刻（date 现取）"
date '+%Y-%m-%d %H:%M:%S %z'
echo
echo "### W1. 名册尺（--no-walk 必带）：本票 r1 实现程名下枚数"
echo '$ R=$(git log --pretty=%h --grep="evidence(157" | grep -v ^89a35e9$)   # 去掉验收程那一枚（它的主题也以 evidence(157 开头）'
R=$(git log --pretty=%h --grep='evidence(157' | grep -v '^89a35e9$' | tr '\n' ' ')
echo "  锚定 grep 枚数＝$(git log --pretty=%h --grep='evidence(157' | grep -vc '^89a35e9$')"
echo '$ git log --no-walk --pretty="%h %ad" --date=format:%H:%M $R'
git log --no-walk --pretty='  %h %ad %s' --date=format:%H:%M $R | cut -c1-58
echo "  --no-walk 名册（枚枚 git show --name-only 并集）＝$(git log --no-walk --pretty=tformat: --name-only $R | sort -u | grep -c .) 枚路径"
echo '  不带 --no-walk 的同一条尺（坑①，本程亲复）＝'
git log --pretty=tformat: --name-only $R | sort -u | wc -l
echo
echo "### W2. AC#5 那格到底交没交：锚点版逐枚数（'对账' / '双向' / 小节标题）"
for a in ec81582 ef1c474 b60c6b7 HEAD; do
  echo "--- $a（$(git log -1 --pretty='%ad' --date=format:%H:%M $a)）---"
  echo "  行数=$(git cat-file blob $a:$F155 | wc -l)  对账=$(git cat-file blob $a:$F155 | grep -c '对账')  双向=$(git cat-file blob $a:$F155 | grep -c '双向')"
  echo "  ### B.1x 小节=$(git cat-file blob $a:$F155 | grep -cE '^### B\.[0-9]+ ')"
  git cat-file blob $a:$F155 | grep -nE '^### B\.1[1-4] ' | sed 's/^/  /'
done
echo
echo "### W3. 验收程那发 '0/0' 读的是哪一版（它 run.sh 的 V6 命令原文）"
echo '$ grep -n "对账" probes/157-accept/run.sh'
grep -n '对账' .scratch/wisp/probes/157-accept/run.sh | sed 's/^/  /'
echo "  本程按它那一版（ec81582）复跑同一把尺＝"
git cat-file blob ec81582:$F155 | sed -n '389,760p' | grep -c '对账\|双向'
echo "  同一把尺打在 B.12 落盘那一枚（ef1c474）＝"
git cat-file blob ef1c474:$F155 | sed -n '389,760p' | grep -c '对账\|双向'
echo "  B.12 落盘时刻与验收程锚点差＝"
git log -1 --pretty='%h %ad  %s' --date=format:'%H:%M:%S' ef1c474 | cut -c1-60
git log -1 --pretty='%h %ad  %s' --date=format:'%H:%M:%S' ec81582 | cut -c1-60
git log -1 --pretty='%h %ad  %s' --date=format:'%H:%M:%S' 89a35e9 | cut -c1-60
echo
echo "### W4. 验收件 §3.1 的档（本程只具名转记、不自立档）"
echo '$ grep -n "^| \*\*AC#" 157-…-accept-r1.md'
grep -nE '^\| \*\*AC#' $ACC | sed -E 's/(.{150}).*/\1…/' | sed 's/^/  /'
echo
echo "### W5. 157 票面：框数与勾数（锚定计数，现量）"
echo "  票面未勾=$(git cat-file blob HEAD:$T157 | grep -c '^- \[ \]')  已勾=$(git cat-file blob HEAD:$T157 | grep -c '^- \[x\]')  AC 行=$(git cat-file blob HEAD:$T157 | grep -cE '^- \[ \] \*\*AC#')"
echo "  票面最后一次改动=$(git log -1 --pretty='%h %ad' --date=format:%H:%M -- $T157)"
ls .scratch/wisp/issues/ | grep -c '157.*-done' | sed 's/^/  157 名下 -done 文件＝/'
echo
echo "### W6. B.1 三行：两把尺并排（它尺＝index 取第一个 '| **'；本尺＝按 '|' 切格取**最后一格非空**）"
echo '$ git cat-file blob b319bab:<153 验收件> 逐行；本尺＝awk -F"|" 从 NF 往回找第一枚非空格'
lastcell() { git cat-file blob b319bab:$A153 | awk -F'|' -v N="$1" 'NR==N{for(i=NF;i>=2;i--){c=$i; gsub(/^ +| +$/,"",c); if(c!=""){print c; exit}}}'; }
echo '  本尺＝最后一格非空格；它尺＝substr(index($0,"| **"))（r1 那一把）：'
for n in 831 832 833 834 835 836 842 843; do
  A=$(git cat-file blob b319bab:$A153 | awk -v N=$n 'NR==N{n=index($0,"| **"); if(n>0) printf "%s", substr($0,n+2,22)}')
  B=$(lastcell $n)
  C=$(git cat-file blob b319bab:$A153 | awk -F'|' -v N=$n 'NR==N{c=$3; gsub(/^ +| +$/,"",c); printf "%s", c}')
  printf '  :%s 它尺=%-24s 本尺最后一格=%-22s 前一程那一格=%s\n' "$n" "$A" "$B" "$C"
done
echo '  两把尺同/不同逐枚判（同＝它尺读数**以**最后一格开头；不同＝它尺吃到了邻格）：'
same=0; diff=0; difflist=""
for n in 831 832 833 834 835 836 842 843; do
  A=$(git cat-file blob b319bab:$A153 | awk -v N=$n 'NR==N{k=index($0,"| **"); if(k>0) printf "%s", substr($0,k+2)}')
  B=$(lastcell $n)
  A2=${A// /}; B2=${B// /}
  case "$A2" in
    "$B2"*) printf '  :%s 同\n' "$n"; same=$((same+1)) ;;
    *)      printf '  :%s 不同（它尺读到邻格：%s｜真档：%s）\n' "$n" "$(echo "$A2" | cut -c1-14)" "$B2"; diff=$((diff+1)); difflist="$difflist :$n" ;;
  esac
done
echo "  同=$same 不同=$diff 不同那几枚=$difflist"
echo '  盘上那一行最后一格逐字（三枚）：'
for n in 832 833 843; do printf '  :%s => %s\n' "$n" "$(lastcell $n)"; done
echo '  同一把尺打在工单点名的 :834/:835/:836（应仍是档）：'
for n in 834 835 836; do printf '  :%s => %s\n' "$n" "$(lastcell $n)"; done
echo
echo "### W7. 本件 B.1 表 HEAD 版那三行的字面（旧行不抹，只取原文）"
grep -nE '^\| \(b\) 问①|^\| \(c\) 问②|^\| AC#3 句②' $F155 | sed 's/^/  /'
echo
echo "### W8. 11-…sh:71 那发挂了 '| head -20' 的位置与不挂 head 的完整名册"
echo '$ sed -n "71p" probes/157/11-…sh'
sed -n '71p' $PROBE11 | sed 's/^/  /'
echo '$ git diff --name-only b23c7f7 6de3d1c5 | grep -v "^docs/" | wc -l'
git diff --name-only b23c7f7 6de3d1c5 | grep -v '^docs/' | wc -l
echo '$ 同上 | head -20 | wc -l（＝那发实际显示的枚数）'
git diff --name-only b23c7f7 6de3d1c5 | grep -v '^docs/' | head -20 | wc -l
echo '$ 被截掉的那几枚（diff -12 之后）：'
git diff --name-only b23c7f7 6de3d1c5 | grep -v '^docs/' | tail -n +21
echo '$ 完整名册里非 .scratch 的枚数＝'
git diff --name-only b23c7f7 6de3d1c5 | grep -vE '^docs/|^\.scratch/' | wc -l
git diff --name-only b23c7f7 6de3d1c5 | grep -vE '^docs/|^\.scratch/'
echo
echo "### W9. 只追加自证：本程进场基版（HEAD＝$(git rev-parse --short HEAD)）起 155 件前 385 行"
mkdir -p /d/tmp/157r2
git cat-file blob HEAD:$F155 | head -385 > /d/tmp/157r2/155-base-385.txt
cmp <(git cat-file blob HEAD:$F155 | head -385) <(git cat-file blob 424ee363:$F155 | head -385) && echo '  cmp 静默＝§0–§11 一字未改'
echo "  153 件前 589 行：$(cmp <(git cat-file blob HEAD:$F153 | head -589) <(git cat-file blob 424ee363:$F153 | head -589) >/dev/null 2>&1 && echo 静默 || echo 不同)"
echo "  本件 '## 附录 B（票 157 代记' 枚数=$(grep -c '^## 附录 B（票 157 代记' $F155)"
echo
echo "### W10. probes/157 枚数（本程新增前后）"
ls -1 .scratch/wisp/probes/157/ | wc -l
ls -1 .scratch/wisp/probes/157/ | sed 's/^/  /'
echo
echo "### W11. B.1 表八行逐枚对：写上去的档 vs 盘上最后一格（档字对不对／括号对不对分开数）"
git cat-file blob HEAD:$F155 | awk -F'|' '/^\| \(a\)|^\| \(b\)|^\| \(c\)|^\| \(d\)|^\| \(e\)|^\| \(f\)|^\| AC#3/{
  row=$2; grade=$3; src=$4; gsub(/^ +| +$/,"",row); gsub(/^ +| +$/,"",grade); gsub(/^ +| +$/,"",src);
  printf "%s\t%s\t%s\n", row, grade, src }' | while IFS=$'\t' read -r row grade src; do
  ln=$(echo "$src" | grep -oE ':[0-9]{3}' | head -1 | tr -d ':')
  real=$(git cat-file blob b319bab:$A153 | awk -F'|' -v N="$ln" 'NR==N{for(i=NF;i>=2;i--){c=$i; gsub(/^ +| +$/,"",c); if(c!=""){print c; exit}}}')
  wg=$(echo "$grade" | grep -oE '^\*\*[^*]+\*\*')
  rg=$(echo "$real" | grep -oE '^\*\*[^*]+\*\*')
  wp=$(echo "$grade" | grep -oE '（[^）]*）' | head -1)
  rp=$(echo "$real" | grep -oE '（[^）]*）' | head -1)
  if [ "$wg" = "$rg" ]; then gsame=档同; else gsame=档不同; fi
  if [ "$wp" = "$rp" ]; then pnote=括号同; else pnote=括号不同; fi
  printf '  %-14s 出处%-8s 本程写=%-40s 盘上真档=%-22s %s %s（本程括号=%s｜盘上括号=%s）\n' "$row" "$src" "$grade" "$real" "$gsame" "$pnote" "$wp" "$rp"
done
echo
echo "### W12. 只读别人的写面：验收件 §3.2 反向那六枚 + 它 V6 命令钉的锚点"
grep -nE '^\| [0-9] \|' $ACC | sed -E 's/(.{110}).*/\1…/' | sed 's/^/  /'
echo
echo
echo "### W13. 编排者给本程的三件裁定，逐件现量（只反查号不够，事也要量）"
echo '$ git cat-file -t bb3a7a1:.scratch/wisp/probes/155/00-anchor-reverse-check.txt'
git cat-file -t bb3a7a1:.scratch/wisp/probes/155/00-anchor-reverse-check.txt | sed 's/^/  /'
echo '  那枚凭据自己的 date 读数＝'
git cat-file blob bb3a7a1:.scratch/wisp/probes/155/00-anchor-reverse-check.txt | grep -A1 '^\$ date' | sed 's/^/    /'
echo '$ git show b319bab:<153 票面> | sed -n "25,29p"（第 2 笔：六件到底在不在票面）'
git cat-file blob b319bab:$T153 | sed -n '25,29p' | sed 's/^/  /'
echo '  票面里数得出的枚数：三问①②③=' "$(git cat-file blob b319bab:$T153 | sed -n '25,29p' | grep -o '[①②③]' | wc -l)" '枚 / ⚑ 弹头=' "$(git cat-file blob b319bab:$T153 | sed -n '25,29p' | grep -c '⚠')" '枚 / 该条款自身＝1 枚交付物 ⇒ 合计 6'
echo '  票面 :81 那一行（字母／"拆 6 件事"在不在票面）＝'
git cat-file blob b319bab:$T153 | sed -n '81p' | sed 's/^/    /'
echo '$ 台账里 A293② 当前状态（编排者是否已自裁；本程不写台账）'
grep -n 'A293' docs/reports/pending-and-issues.md | head -3 | sed -E 's/(.{96}).*/\1…/' | sed 's/^/  /'
echo '  台账 A293② 作废那一行的行号＝'
grep -n '本节②整笔作废' docs/reports/pending-and-issues.md | sed 's/\(:[0-9]*\).*/\1/' | sed 's/^/  /'
echo
echo "### W14. 验收件自己两把尺钉的不是同一版（§0 的行数读数 vs V6 钉的 blob）"
for a in ec81582 ef1c474 c97d706 b60c6b7; do
  printf '    %s（%s）=> 155 件 %s 行 / 对账 %s 枚\n' "$a" "$(git log -1 --pretty=%ad --date=format:%H:%M:%S $a)" "$(git cat-file blob $a:$F155 | wc -l)" "$(git cat-file blob $a:$F155 | grep -c '对账')"
done
echo '  它 V6 那把尺钉的是 ec81582（678 行、对账 0），它验收件 §0 那张漂移表报的"收口"行数是 760＋（＝ef1c474 及以后的形状）。'
echo "### END"
