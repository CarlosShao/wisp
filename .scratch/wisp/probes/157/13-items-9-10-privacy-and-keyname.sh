#!/usr/bin/env bash
# 票 157 实现程 · 第 9/10 笔（153 证据件缺的两笔账）的当轮复算尺
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9
E153=docs/evidence/s1/153-trace-lies-unguarded-r1.md
A153=docs/evidence/s1/153-trace-lies-unguarded-r1-accept-r1.md

echo "### P0. 尺的锚定：两枚文件此刻的行数与工作树脏否"
wc -l $E153 $A153
echo "\$ git status --porcelain -- $E153"
git status --porcelain -- $E153; echo "rc=$?（空＝与 HEAD 同版）"
echo
echo "### P1. 第 9 笔：'件内缺那笔隐私账'——153 证据件里那六个词的全部命中"
for w in 隐私 privacy keep_transcript 私有数据 Q-31 SealDir; do
  printf '%-16s 命中 = ' "$w"; grep -c -- "$w" $E153
done
echo "--- 对照：同一把尺打在 153 验收件上（证明这些词是验收程先量到的）---"
for w in 隐私 privacy keep_transcript Q-31 SealDir; do
  printf '%-16s 验收件命中 = ' "$w"; grep -c -- "$w" $A153
done
echo
echo "### P2. 第 9 笔的地基（本程自己现量，不抄验收件）"
echo "--- (i) 改前锚点 86b0161 上已经带 \"task\" 键的非测试行（验收件说三处）---"
echo "\$ git show 86b0161:internal/agent/loop.go | grep -n '\"task\"'"
git show 86b0161:internal/agent/loop.go | grep -n '"task"'
echo "--- (ii) 交付锚点 6de3d1c 上非测试侧带 \"task\" 的行数（验收件说 7）---"
git grep -c '"task"' 6de3d1c5 -- internal/agent ':!*_test.go'
echo "\$ git grep -n '\"task\"' 6de3d1c5 -- internal/agent ':!*_test.go' | wc -l"
git grep -n '"task"' 6de3d1c5 -- internal/agent ':!*_test.go' | wc -l
git grep -n '"task"' 6de3d1c5 -- internal/agent ':!*_test.go'
echo "--- (iii) compress.go 里那句注释（交付版行号现取）---"
git grep -n 'correlation id' 6de3d1c5 -- internal/agent/compress.go
git show 6de3d1c5:internal/agent/compress.go | sed -n '230,240p'
echo "--- (iv) keep_transcript 是不是硬编码 false、写 true 会不会报错 ---"
git grep -n 'keep_transcript\|KeepTranscript' 6de3d1c5 -- internal/config | head -20
git show 6de3d1c5:internal/config/validate.go | sed -n '78,90p'
echo "--- (v) winsec.SealDir 在 6de3d1c 非测试代码里的调用者枚数 ---"
echo "\$ git grep -n 'SealDir' 6de3d1c5 -- '*.go' ':!*_test.go'"
git grep -n 'SealDir' 6de3d1c5 -- '*.go' ':!*_test.go'
echo "非测试调用者（去掉定义与注释那一枚文件本身）："
git grep -l 'SealDir' 6de3d1c5 -- '*.go' ':!*_test.go'
echo
echo "### P3. 第 10 笔：日志侧 task vs DB 侧 task_id"
echo "--- (i) 规格那两枚出处 ---"
echo "\$ sed -n '74p' docs/specs/SPEC-02*.md"; sed -n '74p' docs/specs/SPEC-02*.md
echo "\$ sed -n '2697p' docs/PLAN.md"; sed -n '2697p' docs/PLAN.md
echo "--- (ii) DB 侧列名在交付码里的形状 ---"
git grep -n 'task_id' 6de3d1c5 -- internal | head -20
echo "--- (iii) 日志侧键名 \"task\" 的行（上面 P2(ii) 那 7 枚）：两面键名对照 ---"
echo "DB 列＝task_id / 日志键＝task"
echo "\$ git grep -c 'task_id' 6de3d1c5 -- internal"
git grep -c 'task_id' 6de3d1c5 -- internal
echo "--- (iv) 153 证据件里 task_id / 键名 有没有被写过 ---"
printf '153 证据件 task_id 命中 = '; grep -c 'task_id' $E153
printf '153 证据件 "task" 键那一枚键名讨论命中 = '; grep -c '键名' $E153
printf '153 验收件 task_id 命中 = '; grep -c 'task_id' $A153
echo
echo "### P4. 这两笔落在 153 票面/验收件的确切出处行号"
grep -n '本格档位：成立，附两笔待补' $A153
sed -n '416p' $A153
grep -n '件内缺那笔隐私账' $A153
