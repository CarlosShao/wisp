#!/bin/sh
# 票 155 对抗验收程 · 开场反查（前提不成立就报回来）
# 尺一律打在 git 对象上（git show / git cat-file），不读脏工作树。
R="D:/work/workspace/projects plans/Wisp"
cd "$R" || exit 9
echo "### 本发时刻"; date "+%Y-%m-%d %H:%M %z"
echo "### 被验锚点 / 当前 HEAD（共享树在动）"
git rev-parse --short bb3a7a1; git rev-parse --short HEAD
git log -1 --format="%h %ad %s" --date=format:"%m-%d %H:%M" HEAD | cut -c1-120

echo
echo "=== A. 派单/证据件点名的每一枚锚号：cat-file -t + 一行标题 ==="
for a in bb3a7a1 829abe5 a52ded8 9136820 c10f0c4 5f6646a 4c43b68 33c6306 2a303c5 fe13c11 \
         ff550f3 5365cb22 86b0161 b23c7f7 6de3d1c5 777d6cc ac7fb00 b319bab; do
  t=$(git cat-file -t "$a" 2>&1); h=$(git log -1 --format="%h|%ad|%s" --date=format:"%m-%d %H:%M" "$a" 2>/dev/null | cut -c1-70)
  printf "%-10s %-8s %s\n" "$a" "$t" "$h"
done
echo "-- ff550f3^ 展开 --"; git rev-parse --short "ff550f3^"
echo "-- 6de3d1c5 是否已进主线 --"; git merge-base --is-ancestor 6de3d1c5 HEAD && echo YES || echo NO
echo "-- b23c7f7 是否 6de3d1c5 的祖先 --"; git merge-base --is-ancestor b23c7f7 6de3d1c5 && echo YES || echo NO

echo
echo "=== B. 派单前提 P1：行号天然带版本（它写 run.go:576） ==="
for a in 5365cb22 bb3a7a1 HEAD; do
  printf "%-9s " "$a"; git show "$a:cmd/wisp/run.go" | grep -n 'res := loop.Run(ctx, task)'
done
echo "-- 同一发在 HEAD 的 run.go 总行数（漂移幅度） --"
git show 5365cb22:cmd/wisp/run.go | wc -l; git show HEAD:cmd/wisp/run.go | wc -l

echo
echo "=== C. 派单前提 P2：6de3d1c5 存在，但它指的是哪一版（internal/agent 三枚被测文件） ==="
for f in internal/agent/compress.go internal/agent/loop.go internal/agent/compress_trace_test.go; do
  printf "%-38s b23c7f7=%s  6de3d1c5=%s  %s\n" "$f" \
    "$(git rev-parse --short=12 b23c7f7:$f)" "$(git rev-parse --short=12 6de3d1c5:$f)" \
    "$(if [ "$(git rev-parse b23c7f7:$f)" = "$(git rev-parse 6de3d1c5:$f)" ]; then echo SAME; else echo DIFFER; fi)"
done
echo "-- b23c7f7..6de3d1c5 在 internal/ cmd/ go.mod 上到底动了什么 --"
git diff --name-only b23c7f7 6de3d1c5 -- internal cmd go.mod go.sum; echo "rc=$?（空＝没动）"
echo "-- 同一区间动了哪些面（全量，前 20 行） --"
git diff --name-only b23c7f7 6de3d1c5 | head -20
echo "-- 枚数 --"; git diff --name-only b23c7f7 6de3d1c5 | wc -l

echo
echo "=== D. 派单前提 P4：155 §6.2 两条注入点名的路径盘上反查 ==="
for p in .scratch/wisp/dispatches/2026-09-26-122x-accept-156.md \
         .scratch/wisp/tasks/bca7d7c.output \
         .scratch/wisp/issues/155-three-unjudged-cells-r1-continued.md; do
  if [ -e "$p" ]; then echo "EXISTS  $p"; else echo "ABSENT  $p"; fi
done
echo "-- .scratch/wisp 下有没有 tasks/ 或任何 156 件 --"
ls .scratch/wisp/; find .scratch -name "*156*" -o -name "*bca7d7c*" | head; echo "find rc=$?"
echo "-- dispatches 最新 5 枚 --"; ls .scratch/wisp/dispatches/ | tail -5

echo
echo "=== E. 工作树脏在哪（不是被验版本，逐枚登记，一律不碰不还原） ==="
git status --porcelain | head -40
echo "-- 枚数 --"; git status --porcelain | wc -l
