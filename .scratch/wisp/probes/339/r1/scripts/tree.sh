#!/bin/bash
# 339-r1 AC#0 台件①：仓外导出树（Linux 容器腿 + 本机 Windows 腿各一棵）。
# 用法：bash tree.sh <linux|win> <commit>
# 八件套 + 仓根 docs/，手工三件（dll / frontend/dist）⛔ 跟踪、单独补。
# 落点一律 D:\tmp\wisp\20261018\339r1\...（⛔ 盘符根）。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
KIND="$1"
ANCHOR="$2"
ROOT="/d/tmp/wisp/20261018/339r1"
TREE="$ROOT/$KIND/tree"
cd "$REPO" || exit 9
mkdir -p "$ROOT/$KIND"
git archive "$ANCHOR" go.mod go.sum internal cmd tools scripts frontend models docs > "$ROOT/$KIND/tree.tar" || exit 8
mkdir -p "$TREE" && tar -xf "$ROOT/$KIND/tree.tar" -C "$TREE" || exit 7
# 证非空：先数 .go
GOC=$(find "$TREE" -name '*.go' | wc -l)
echo "anchor=$ANCHOR kind=$KIND tree=$TREE go_files=$GOC"
# 手工三件（都⛔ 跟踪，构建/加载需要）
mkdir -p "$TREE/build" "$TREE/third_party/sherpa-onnx" "$TREE/frontend"
cp -r /d/work/workspace/projects\ plans/Wisp/frontend/dist "$TREE/frontend/" 2>/dev/null || echo "note: frontend/dist copy skipped rc=$?"
ls "$REPO/build"/*.dll 2>/dev/null | while read -r f; do cp "$f" "$TREE/build/"; done
ls "$REPO/third_party/sherpa-onnx"/*.dll 2>/dev/null | while read -r f; do cp "$f" "$TREE/third_party/sherpa-onnx/"; done
echo "risk_go=$(find "$TREE/internal/risk" -name '*.go' | wc -l)"
echo "instr_present=$(ls "$TREE/internal/risk" | grep -c anchor_miss_reason_334)"
echo "PRISTINE_SHA334=$(git -C "$REPO" rev-parse "$ANCHOR:internal/risk/anchor_miss_reason_334_test.go")"
