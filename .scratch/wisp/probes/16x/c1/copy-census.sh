#!/usr/bin/env bash
# 16x-c1 同源拷贝普查（两把尺）— 只读取数，不落任何仓库写面。
# 锚点：0d1764731973a998c5d6da5839070151fb9b04a9（本程 step 0 `git rev-parse HEAD` 现量）
# 用法：bash .scratch/wisp/probes/16x/c1/copy-census.sh
set -u
ANCHOR=0d1764731973a998c5d6da5839070151fb9b04a9
# 普查语料＝"自称是拷贝的那一类文件"：定稿 + specs + 薄索引 + 派单规矩 + 生产码/CI。
# 排除：.scratch（工单池/探针＝历史读数，不是契约拷贝）、docs/evidence、docs/reports（裁决史）。
CORPUS_DOCS=(docs/PLAN.md docs/specs AGENTS.md .scratch/wisp/issues/README.md)

hr() { printf '\n===== %s =====\n' "$1"; }

hr "尺A-1  词形：自称'逐字/唯一来源/权威表/一一对应'的条款（定稿+specs+薄索引）"
grep -rnE "逐字|照抄|verbatim|唯一来源|权威表|一一对应|与 D3[0-9] 表一致|与 D43 转移表" \
  docs/PLAN.md docs/specs AGENTS.md 2>/dev/null
printf '尺A-1 命中行数='
grep -rcE "逐字|照抄|verbatim|唯一来源|权威表|一一对应|与 D3[0-9] 表一致|与 D43 转移表" \
  docs/PLAN.md docs/specs AGENTS.md 2>/dev/null | awk -F: '{s+=$NF} END {print s+0}'

hr "尺A-2  同一词形在 README/工单规矩里的命中"
grep -rnE "逐字|唯一来源|权威表" .scratch/wisp/issues/README.md 2>/dev/null
printf 'rc(0=有命中,1=无命中)='
grep -rqE "逐字|唯一来源|权威表" .scratch/wisp/issues/README.md 2>/dev/null; echo $?

hr "尺B-1  D34 shell.exec 行（PLAN.md:2563）的原文串在语料里的字面命中"
grep -rnE '`shell\.exec`.*(默认禁用|argv 向量)' docs/PLAN.md docs/specs AGENTS.md 2>/dev/null

hr "尺B-2  D34 task.* 行（PLAN.md:2561）的字面命中（含 specs）"
grep -rnE 'task\.list' docs/PLAN.md docs/specs AGENTS.md 2>/dev/null

hr "尺B-3  D43 计数断言（'20 态'／'40 条'／'42 条'）"
grep -rnE "20 态|40 条|42 条" docs/PLAN.md docs/specs AGENTS.md 2>/dev/null

hr "尺B-4  AwaitingApproval 在语料（非 .scratch／非 evidence／非 reports）里的命中枚数"
git grep -c "AwaitingApproval" "$ANCHOR" -- docs/PLAN.md 'docs/specs/*' AGENTS.md 2>/dev/null

hr "尺B-5  D46 那句'保持不变'"
grep -rn "保持不变" docs/PLAN.md docs/specs 2>/dev/null

hr "尺B-6  D33 那句'argv: string\[\]'"
grep -rn "argv: string\[\]" docs/PLAN.md docs/specs 2>/dev/null

hr "尺B-7  生产注册表里 fs.* 工具名枚数（尺1＝Name() 实现）"
git grep -nE 'func \([a-zA-Z0-9_ *]+\) Name\(\) string \{ return "' "$ANCHOR" \
  -- '*.go' ':!*_test.go' ':!.scratch' | wc -l

hr "尺B-8  生产注册表里 fs.* 工具名枚数（尺2＝字符串字面量）"
git grep -nE 'return "fs\.[a-z]+"' "$ANCHOR" -- 'internal/tools/' ':!*_test.go' | wc -l

hr "尺B-9  OpenScope 的生产调用枚数（尺1＝字面量，尺2＝文件名去重）"
git grep -n "OpenScope" "$ANCHOR" -- '*.go' ':!*_test.go' ':!.scratch' \
  | grep -v "provenance.go" | sed 's/^[^:]*://' | sed 's/:.*//' | sort -u
git grep -l "OpenScope" "$ANCHOR" -- '*.go' ':!*_test.go' ':!.scratch' \
  | grep -v "internal/risk/provenance.go"

hr "尺C   design/ 与 frontend/ 对 D34 单元格字面的命中（只枚文件名，不读内容）"
git grep -l "shell\.exec" "$ANCHOR" -- design frontend 2>/dev/null
