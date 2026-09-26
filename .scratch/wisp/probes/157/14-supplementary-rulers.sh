#!/usr/bin/env bash
# 票 157 实现程 · 补尺：第 4 笔的"够不着"那一半、第 5 笔的函数体边界、第 8 笔的闭包逐包改动、几枚锚定 grep
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9
E155=docs/evidence/s1/155-three-unjudged-cells-r1.md
A=5365cb22

echo "### Q1. 第 4 笔的地基：Compressor 到底够不够得着 l.current"
echo "--- Compressor 结构体原文（$A）---"
git show $A:internal/agent/compress.go | sed -n '/^type Compressor struct/,/^}/p'
echo "--- Compressor 的全部方法签名（有没有一枚接 *Loop / *observe.Root）---"
git show $A:internal/agent/compress.go | grep -n '^func (c \*Compressor)'
echo "--- compress.go 全文里 current / observe.Root 的命中枚数 ---"
printf 'compress.go 内 current 命中 = '; git show $A:internal/agent/compress.go | grep -c 'current'
printf 'compress.go 内 observe\\.Root 命中 = '; git show $A:internal/agent/compress.go | grep -c 'observe\.Root'
echo
echo "### Q2. 第 5 笔的地基：:369 与 :396 之间有没有函数边界"
echo "--- \$A 版 loop.go :340..:396 里的 ^func / ^} 行（相对行号）---"
git show $A:internal/agent/loop.go | sed -n '340,396p' | grep -n '^func \|^}' || echo "(无：两枚行号在同一枚函数体内)"
echo "--- :369 那一行的缩进层级（制表符数）vs :396 ---"
git show $A:internal/agent/loop.go | sed -n '369p;396p' | cat -A | cut -c1-60
echo "--- 从 :339 起，for 循环从哪一枚行号开始（res 在循环之前）---"
git show $A:internal/agent/loop.go | sed -n '339,396p' | grep -n '^	for {\|^		for {\|^	res := Result'
echo
echo "### Q3. 第 8 笔的补尺：闭包逐包在 b23c7f7..6de3d1c5 的改动枚数（本程那 16 枚 test-dep 全列）"
for p in internal/buildinfo internal/config internal/llm internal/llm/adaptertest internal/llm/anthropic internal/llm/golden internal/llm/openaichat internal/llm/openairesponses internal/memory internal/observe internal/plugin internal/proc internal/risk internal/secret internal/statemachine internal/winsec internal/tools; do
  printf '%-30s ' "$p"; git diff --name-only b23c7f7 6de3d1c5 -- $p | wc -l
done
echo
echo "### Q4. 锚定 grep：几枚'表与盘'的形状尺"
echo "\$ grep -c '闭包' $E155"; grep -c '闭包' $E155
echo "\$ grep -c '同帧持有者\\|持有者' $E155"; grep -c '持有者' $E155
echo "--- 155 证据件里 AdmitTask / j.taskID / journal 的全部命中（第 6/7 笔是否真漏报）---"
printf '155 内 AdmitTask 命中 = '; grep -c 'AdmitTask' $E155
printf '155 内 j.taskID 命中 = '; grep -c 'j\.taskID' $E155
printf '155 内 journal 命中 = '; grep -c 'journal' $E155
printf '155 内 newTaskJournal 命中 = '; grep -c 'newTaskJournal' $E155
echo "--- 155 §2 那张'副本'表的四枚名册（第 4 笔宣称的四枚）---"
sed -n '111,117p' $E155 | cut -c1-46
echo "--- 155 证据件里 '3 处读点' 与 '晚于压缩' 的原句行号 ---"
grep -n '3 处读点\|晚于压缩' $E155
echo "--- 155 §3.3 那两枚'还原自证'行（第 8 笔射程声明所在）---"
grep -n '6de3d1c5` 的 blob 号逐字相同\|快照取版＝blob 级' $E155
