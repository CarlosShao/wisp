#!/bin/sh
# 322-r2 台面配方：仓外导出树（八件套＋docs＋手工补三件 untracked 前置）。
# ⛔ 在共享工作树跑 go test/go build —— 一切读数都在这棵树里取。
set -eu
# 机主 2026-10-18 10:1x 新铁律（覆盖原台面配方）：盘符根零污染，唯一允许落点
# D:\tmp\wisp\20261018\322r2\*。任何腿⛔ 删任何东西 ⇒ 本脚本⛔ rm -rf，
# 每棵树都是 mktemp 唯一新名（旧树留在盘上，只建不删）。
ANCHOR="$1"
ROOT="${2:-/d/tmp/wisp/20261018/322r2}"   # 盘符根下零物件；一律在这棵树里
PREFIX="${3:-tree}"
mkdir -p "$ROOT"
TREE=$(mktemp -d "$ROOT/$PREFIX.XXXXXX")
REPO="/d/work/workspace/projects plans/Wisp"

cd "$REPO"

# 八件套 + docs（有些测试读 docs/PLAN.md；internal/panel/git_test.go 一族要它）
git archive "$ANCHOR" go.mod go.sum internal cmd tools scripts frontend models docs | tar -x -C "$TREE"

# 手工补三枚未跟踪前置（⛔ 跟踪）
mkdir -p "$TREE/build" "$TREE/third_party/sherpa-onnx"
cp "$REPO/build/onnxruntime.dll" "$TREE/build/"
cp "$REPO/build/sherpa-onnx-c-api.dll" "$TREE/build/"
cp "$REPO/build/sherpa-onnx-cxx-api.dll" "$TREE/build/"
cp "$REPO/third_party/sherpa-onnx/onnxruntime.dll" "$TREE/third_party/sherpa-onnx/"
cp "$REPO/third_party/sherpa-onnx/sherpa-onnx-c-api.dll" "$TREE/third_party/sherpa-onnx/"
cp "$REPO/third_party/sherpa-onnx/sherpa-onnx-cxx-api.dll" "$TREE/third_party/sherpa-onnx/"
mkdir -p "$TREE/frontend"
cp -r "$REPO/frontend/dist" "$TREE/frontend/"

# 取出后先数 .go 证非空（"0" ⛔ 是读数）
NGO=$(find "$TREE" -name '*.go' -type f | wc -l)
echo "ANCHOR=$ANCHOR"
echo "TREE=$TREE"
echo "GOCOUNT=$NGO"
echo "GOMOD=$(head -1 "$TREE/go.mod")"
[ "$NGO" -gt 0 ] || { echo "TREE EMPTY = 台面故障"; exit 90; }
cd "$TREE" && git rev-parse HEAD >/dev/null 2>&1 || true
echo "SHORTNAME:"; cmd //c "dir /x $ROOT" | grep -i "$(basename "$TREE")" || echo "(see dir /x output above)"
