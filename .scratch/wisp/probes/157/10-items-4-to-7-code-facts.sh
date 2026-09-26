#!/usr/bin/env bash
# 票 157 实现程 · 第 4/5/6/7 笔的当轮复算尺（全部打在 git 对象上，不读脏工作树）
# 锚点：5365cb22 = 改前那版（= ff550f3^，155 证据件 §0 已证同源）
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9
A=5365cb22

echo "### S1. 第 4 笔：l.current 的全部提及（宣称 3 处读点、只列 2 枚、其中一枚是写点）"
echo "\$ git grep -nE '\\.current\\b' $A -- internal/agent/loop.go"
git grep -nE '\.current\b' $A -- internal/agent/loop.go; echo "rc=$?"
echo "--- 逐枚看那一行的原文（判定它是注释／读／写）---"
for n in 189 348 353 523 987; do
  printf ':%s  ' "$n"
  git show $A:internal/agent/loop.go | sed -n "${n}p"
done
echo "--- 字段声明那一枚（第 4 枚提及？）---"
echo "\$ git grep -n 'current \\*observe.Root' $A -- internal/agent/loop.go"
git grep -n 'current \*observe.Root' $A -- internal/agent/loop.go; echo "rc=$?"
echo "--- setCurrent 的定义与全部调用点 ---"
git grep -n 'setCurrent' $A -- internal/agent/loop.go
echo
echo "### S2. 第 5 笔：Result.TaskID 的赋值时刻 vs 压缩那一发（宣称'晚于压缩'）"
echo "\$ git show $A:internal/agent/loop.go | grep -n 'res := Result{\|l.comp.Compress(\|func (l \*Loop) run(' "
git show $A:internal/agent/loop.go | grep -n 'res := Result{\|l\.comp\.Compress(\|func (l \*Loop) run('
echo "--- :369 与 :396 之间是否同一函数体、是否有 return/for 边界 ---"
git show $A:internal/agent/loop.go | sed -n '339,341p;365,372p;392,400p'
echo "--- 行号差 ---"
printf '396-369 = %s\n' "$((396-369))"
echo "\$ git grep -n 'TaskID' $A -- internal/agent/loop.go（全部 4 枚声明/赋值/读取点）"
git grep -n 'TaskID' $A -- internal/agent/loop.go
echo
echo "### S3. 第 6 笔：j.taskID 是同帧第二枚持有者（155 §2 表漏报）"
echo "\$ git grep -n 'newTaskJournal' $A -- internal/agent"
git grep -n 'newTaskJournal' $A -- internal/agent
echo "--- :356 那一行原文 ---"
git show $A:internal/agent/loop.go | sed -n '356p'
echo "\$ git grep -nE 'taskID' $A -- internal/agent/journal.go"
git grep -nE 'taskID' $A -- internal/agent/journal.go
echo "--- taskJournal 结构体原文 + newTaskJournal 函数体 ---"
git show $A:internal/agent/journal.go | sed -n '/^type taskJournal struct/,/^}/p'
git show $A:internal/agent/journal.go | sed -n '/^func newTaskJournal/,/^}/p'
echo "--- journal.go 是不是同一枚 package（同包 ⇒ 字段可直接访问）---"
git show $A:internal/agent/journal.go | sed -n '1,3p'
echo
echo "### S4. 第 7 笔：l.opt.AdmitTask(taskID) 那一刻 id 已出包（155 §2 表漏报）"
echo "\$ git grep -n 'AdmitTask' $A -- internal cmd"
git grep -n 'AdmitTask' $A -- internal cmd
echo "--- loop.go:362 那一行原文 ---"
git show $A:internal/agent/loop.go | sed -n '360,364p'
echo "--- Options.AdmitTask 的签名 ---"
git grep -n 'AdmitTask ' $A -- internal/agent/loop.go
echo
echo "### S5. 第 6/7 笔共同的地基：压缩那一发 :396 之前这几枚是否都已在同一作用域"
echo "\$ git show $A:internal/agent/loop.go | sed -n '339,400p' 的行号图（只取关键行）"
git show $A:internal/agent/loop.go | sed -n '339,400p' | grep -n 'root :=\|setCurrent\|j :=\|AdmitTask\|res := Result\|Compress(\|for {'
