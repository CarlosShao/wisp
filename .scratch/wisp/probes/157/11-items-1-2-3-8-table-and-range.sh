#!/usr/bin/env bash
# 票 157 实现程 · 第 1/2/3 笔（表与名册）＋ 第 8 笔（射程声明）的当轮复算尺
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9
E155=docs/evidence/s1/155-three-unjudged-cells-r1.md
A153=docs/evidence/s1/153-trace-lies-unguarded-r1-accept-r1.md
T153=.scratch/wisp/issues/153-the-compression-trace-has-two-unguarded-ways-to-lie-swapping-the-ran-guard-for-need-keeps-all-four-nails-green-and-the-trace-carries-no-taskid.md

echo "### T1. 第 1 笔：155 §4 推进表的列形状（有没有'该枚 AC 全档'那一列）"
echo "--- §4 那节标题与表头逐行（带行号）---"
grep -n '^## 4\. 三块' $E155
sed -n '257,266p' $E155 | cat -n
echo "--- 锚定尺：'全档' 那一列在 155 证据件里的命中枚数 ---"
echo "\$ grep -c '全档' $E155"; grep -c '全档' $E155
echo "--- 155 §4 表里 (a)(e)(f) 三枚字母各出现几回（锚定到行首竖线）---"
for L in '(a)' '(e)' '(f)' '(b)' '(c)' '(d)'; do
  printf '155 §4 内 %-5s 命中 = ' "$L"
  sed -n '257,266p' $E155 | grep -o -F "$L" | wc -l
done
echo "--- 档位词在 §4 那十行里的枚数（成立／附条件／未裁）---"
for W in 成立 附条件 未裁; do
  printf '155 §4 内 %s = ' "$W"; sed -n '257,266p' $E155 | grep -c "$W"
done
echo
echo "### T2. 第 1 笔的地基：153 验收件 §12 那张 AC#2 六件表在盘上的真实行号与档位"
echo "\$ grep -n 'AC#2 痕带 taskID' $A153"; grep -n 'AC#2 痕带 taskID' $A153
echo "--- 逐行看 :829-:836 的行首字母与行尾档位列 ---"
git show HEAD:$A153 | sed -n '827,837p' | cut -c1-24
echo "--- 六件各自的档位（只取'| **' 之后那一段）---"
awk 'NR>=831 && NR<=836 {n=index($0,"| **"); if(n>0){printf ":%s  %s\n", NR, substr($0,n,40)}}' $A153
echo
echo "### T3. 第 2 笔：'票面拆出的 6 件'那一句在 155 证据件里的全部出现处"
echo "\$ grep -n '6 件' $E155"; grep -n '6 件' $E155
echo "--- 153 验收件 §12 里'票面把它拆成 6 件事'的原句（第 :827 行）---"
sed -n '827p' $A153
echo "--- 153 票面 AC#2 那一条的原文枚数盘点（锚定尺）---"
echo "\$ grep -c '^- \\[ \\]' $T153"; grep -c '^- \[ \]' $T153
echo "\$ grep -n 'AC#2' $T153"; grep -n 'AC#2' $T153
echo "--- AC#2 那一条跨几行、里面有几个编号件 ---"
awk 'NR>=26 && NR<=29 {printf ":%s %s\n", NR, substr($0,1,58)}' $T153
echo "--- 票面有没有 (a)/(b)/(c)/(d)/(e)/(f) 这套字母（锚定到票面全文）---"
for L in '(a)' '(b)' '(c)' '(d)' '(e)' '(f)'; do
  printf '153 票面 %-5s 命中 = ' "$L"; grep -c -F "$L" $T153
done
echo
echo "### T4. 第 3 笔：155 §4 的 AC#3 行 vs 153 验收件 §12.1 句① 的档"
echo "--- 155 §4 表末行（AC#3 那一行）原文 ---"
sed -n '263p' $E155
echo "--- 153 验收件 §12.1 AC#3 那两行 :842/:843 的档位段 ---"
awk 'NR>=840 && NR<=843 {n=index($0,"| **"); if(n>0) printf ":%s %s => %s\n", NR, substr($0,1,10), substr($0,n,26)}' $A153
echo "--- 句① 那一格前一程落在哪一节（验收件 :842 行首）---"
sed -n '842p' $A153 | cut -c1-120
echo "--- 155 证据件 §3.4/§4 里'句①'的全部出现处与是否带档字 ---"
grep -n '句①' $E155
echo
echo "### T5. 第 8 笔：6de3d1c5 等不等于交付那一版（本程自己那把尺）"
echo "--- 三枚被测文件在 b23c7f7 / 6de3d1c5 / bb3a7a1 三枚锚点上的 blob ---"
for f in internal/agent/compress.go internal/agent/loop.go internal/agent/compress_trace_test.go; do
  for r in b23c7f7 6de3d1c5 bb3a7a1; do
    printf '%-38s %-9s %s\n' "$f" "$r" "$(git rev-parse --short=12 $r:$f)"
  done
done
echo "--- b23c7f7..6de3d1c5 在 internal/agent 上的改动枚数 ---"
echo "\$ git diff --name-only b23c7f7 6de3d1c5 -- internal/agent | wc -l"
git diff --name-only b23c7f7 6de3d1c5 -- internal/agent | wc -l
echo "--- 同一区间 internal/** 上唯一动过的那枚 ---"
git diff --name-only b23c7f7 6de3d1c5 -- internal/
echo "--- 祖先关系 ---"
git merge-base --is-ancestor b23c7f7 6de3d1c5 && echo "ANCESTOR=YES" || echo "ANCESTOR=NO"
echo "--- 区间在 docs/ 之外还动了什么（只列非 docs 路径，供射程句写全）---"
git diff --name-only b23c7f7 6de3d1c5 | grep -v '^docs/' | head -20
echo
echo "### T6. 第 8 笔的第二半：internal/agent 的依赖闭包里到底有没有 internal/tools"
echo "--- 6de3d1c5 版 internal/agent 全部 go 文件里 import wisp/internal/* 的名字册（去重）---"
for f in $(git ls-tree -r --name-only 6de3d1c5 -- internal/agent); do
  git show 6de3d1c5:$f | grep -o 'github.com/CarlosShao/wisp/internal/[a-z_/]*'
done | sed 's|github.com/CarlosShao/wisp/||' | sort | uniq -c
echo "--- 6de3d1c5 上引用 internal/tools 的文件名册（全仓，非测试）---"
git grep -l 'wisp/internal/tools' 6de3d1c5 -- '*.go' | sed 's|^6de3d1c5:||' | sort
echo "--- 闭包里那几包各自在该区间的改动枚数 ---"
for p in internal/observe internal/llm internal/memory internal/plugin internal/risk internal/config internal/secret internal/winsec internal/agent internal/tools; do
  printf '%-18s ' "$p"; git diff --name-only b23c7f7 6de3d1c5 -- $p | wc -l
done
