#!/usr/bin/sh
# 票 164 验收 v1 · AC#2 台件 —— "改前那一发"的**真**未修码读数（验收程现量）
#
# 为什么是这一形，而不是往 probes/ 里塞一份 f206e9f 的完整副本：
#   ① 主树里再落一份全仓副本＝把同一棵树提交两遍（几千枚文件进历史），
#      共享工作树里那是别人的现场风险，也违背"临时件只建不删"的轻量前提；
#   ② 本程要的那发读数**不依赖跑测试**：AC#2 问的是"这一发落不落在工具上"，
#      在锚点上它由三件事实完全决定 —— 名字在不在 Go 源码里、生产注册表
#      由哪几枚 Entries 组成、`Lookup` 失败支的文案是什么。三件都能就地量；
#   ③ 但为了把"那枚改前用例到底吃不吃这棵树"钉死，本件**仍然**在仓外
#      （/tmp，绝不落在仓库目录内）解一份锚点树，把同一个文件在
#      改前树与改后树各跑一遍，比的是**读数逐字是否相同**。相同＝它不吃树
#      ＝它对 AC#2 这一问是恒真形状。这一步是本件的结论性凭据。
#
# 只读：本件不改主树任何一字节；产物落在 /tmp 与本目录的日志。
set -eu

REPO=$(git rev-parse --show-toplevel)
ANCHOR=f206e9f   # 实现程 step-0 锚（其证据件 §1 自报）
cd "$REPO"

echo "== [1] 锚点 $ANCHOR 上 task.output 这个字面量在 Go 源码里是否存在 =="
git grep -n "task\.output" "$ANCHOR" -- internal cmd || echo "  -> 0 命中（锚点上全仓 Go 无此名）"

echo "== [2] 锚点上 internal/tools 有无 task*.go =="
git ls-tree "$ANCHOR" internal/tools/ --name-only | grep -i "task" || echo "  -> 无 task*.go"

echo "== [3] 锚点上生产注册表由哪几枚 Entries 组成（cmd/wisp/run.go） =="
git show "$ANCHOR":cmd/wisp/run.go | grep -nE "tools\.Builtin[A-Za-z]*Entries"

echo "== [4] 锚点上 Lookup 失败支的文案本体 =="
git show "$ANCHOR":internal/tools/bridge.go | grep -n "未知工具"

echo "== [5] 同一枚改前用例：改前树 vs 改后树，读数是否逐字相同 =="
# 每次现量都开一枚全新的仓外目录（mktemp -d）：本程**不跑任何删除命令**，
# 旧的 /tmp 副本留着不管（临时件只建不删），它在下一次运行前不影响任何读数。
BASE=$(mktemp -d /tmp/wisp-anchor-164v1.XXXXXX)
git archive "$ANCHOR" | tar -x -C "$BASE"
cp internal/tools/task_output_ac2_before_test.go "$BASE/internal/tools/"
echo "-- 改前树（$ANCHOR）--"
(cd "$BASE" && go test -count=1 -v -run TestTaskOutputAC2BeforeLegIsUnreachable ./internal/tools/ \
  | grep -E "BEFORE verbatim|^--- (PASS|FAIL)")
echo "-- 改后树（HEAD）--"
go test -count=1 -v -run TestTaskOutputAC2BeforeLegIsUnreachable ./internal/tools/ \
  | grep -E "BEFORE verbatim|^--- (PASS|FAIL)"
echo "  -> 两行 Text 逐字相同 ⇒ 该用例的读数与被测树无关（它自建注册表、且注册表是手抄的 BuiltinFSEntries）"

echo "== [6] 生产调用者那一问（HEAD 现量） =="
echo "-- 名册写者：cmd/ 与 internal/ 里对 TaskRoster.Record 的调用（排测试） --"
grep -rn "\.Record(" --include=*.go internal cmd | grep -v "_test.go" || echo "  -> 0 枚写者"
echo "-- RunAsync 的非测试命中 --"
grep -rn "RunAsync" --include=*.go internal cmd | grep -v "_test.go"

echo "== [7] 锚点树基线名册 vs HEAD 名册（两向 comm） =="
(cd "$BASE" && go test -count=1 -v ./internal/tools/ ./internal/agent/ > /tmp/164v1-anchor.log 2>&1 || true)
grep -oE '^--- (PASS|FAIL|SKIP): [A-Za-z0-9_/]+' /tmp/164v1-anchor.log | sed 's/^--- [A-Z]*: //' \
  | grep -v "^TestTaskOutputAC2BeforeLegIsUnreachable$" | sort > /tmp/164v1-anchor-names.txt
go test -count=1 -v ./internal/tools/ ./internal/agent/ > /tmp/164v1-head.log 2>&1 || true
grep -oE '^--- (PASS|FAIL|SKIP): [A-Za-z0-9_/]+' /tmp/164v1-head.log | sed 's/^--- [A-Z]*: //' | sort > /tmp/164v1-head-names.txt
echo "锚点(净)=$(wc -l < /tmp/164v1-anchor-names.txt) HEAD=$(wc -l < /tmp/164v1-head-names.txt)"
echo "消失的一栏=$(comm -23 /tmp/164v1-anchor-names.txt /tmp/164v1-head-names.txt | wc -l)"
echo "新增的一栏=$(comm -13 /tmp/164v1-anchor-names.txt /tmp/164v1-head-names.txt | wc -l)"
comm -13 /tmp/164v1-anchor-names.txt /tmp/164v1-head-names.txt
