#!/bin/sh
# gw-1 量色腿 台面＝仓外导出树（★⛔ 在共享工作树跑 go test）。
# 用法：sh 02-build-tree.sh <anchor> <outname>
# 八件套＋仓根 docs/＋手工补三件（build/*.dll、third_party/sherpa-onnx/*.dll、frontend/dist＝都⛔ 跟踪）。
# ★取出后先数 .go 证非空（"0"⛔ 是读数）。⛔ 删任何东西（临时件只建不删）。
set -u
REPO='/d/work/workspace/projects plans/Wisp'
ROOT=/d/tmp/wisp/20261020/gw1
ANCHOR="$1"
OUT="$2"
cd "$REPO" || exit 9
PIECES='go.mod go.sum internal cmd tools scripts frontend models docs'
T="$ROOT/$OUT"
mkdir -p "$ROOT/logs" "$T"
git archive --format=tar "$ANCHOR" $PIECES | tar -x -C "$T"
RC=$?
mkdir -p "$T/build" "$T/third_party/sherpa-onnx" "$T/frontend/dist"
cp -f build/*.dll "$T/build/" 2>/dev/null || true
cp -f third_party/sherpa-onnx/*.dll "$T/third_party/sherpa-onnx/" 2>/dev/null || true
cp -a frontend/dist/. "$T/frontend/dist/" 2>/dev/null || true
echo "ARCHIVE_RC=$RC"
echo "ANCHOR_ARG=$ANCHOR FULL=$(git rev-parse "$ANCHOR")"
echo "TREE=$T"
echo "GO_FILES=$(find "$T" -name '*.go' | wc -l)"
echo "WINSEC_GO=$(find "$T/internal/winsec" -name '*.go' | wc -l)"
echo "SECRET_GO=$(find "$T/internal/secret" -name '*.go' | wc -l)"
echo "CONFIG_GO=$(find "$T/internal/config" -name '*.go' | wc -l)"
echo "MEMORY_GO=$(find "$T/internal/memory" -name '*.go' | wc -l)"
echo "AGENT_GO=$(find "$T/internal/agent" -maxdepth 1 -name '*.go' | wc -l)"
echo "PANEL_GO=$(find "$T/internal/panel" -name '*.go' | wc -l)"
echo "CMDWISP_GO=$(find "$T/cmd/wisp" -name '*.go' | wc -l)"
echo "DLL_BUILD=$(ls "$T/build" | grep -c . || true) DLL_SHERPA=$(ls "$T/third_party/sherpa-onnx" | grep -c . || true) DIST_FILES=$(find "$T/frontend/dist" -type f | wc -l)"
echo "--- 承重 blob（内容锚，⛔ 绝对行号；量＝git rev-parse <锚>:<path>）---"
for p in internal/winsec/winsec.go internal/winsec/pinnedleaf_342_windows.go internal/agent/approval/queue.go cmd/wisp/panel_control_windows.go; do
  git rev-parse "$ANCHOR:$p" 2>/dev/null | sed "s|^|BLOB $p=|"
done
echo "--- 用例件在位尺（具名连带三枚）---"
for f in internal/winsec/private_fail_test.go internal/winsec/openleaf_322_other_test.go internal/winsec/placement_symlink_113_other_test.go internal/winsec/nametopup_353_windows_test.go; do
  test -f "$T/$f" && echo "PRESENT $f" || echo "ABSENT  $f"
done
