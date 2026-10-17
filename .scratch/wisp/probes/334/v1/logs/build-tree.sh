#!/usr/bin/env bash
# 334-v1 台面配方（派单通用块那一节）：仓外导出树 = 八件套 + 仓根 docs/ + 手工三件。
# 用法：build-tree.sh <锚> <D 卷树根>
# ⛔ rm -rf：临时件只建不删；每发一棵新树。树根名故意 >8 字符，让 NTFS 生成 8.3 别名。
set -u
REPO="D:/work/workspace/projects plans/Wisp"
ANCHOR="${1:?需要锚}"
ROOT="${2:?需要树根}"
EIGHT="go.mod go.sum internal cmd tools scripts frontend models"
export GOPATH='D:\work\base\gopath'
echo "date=$(date '+%F %T %z') anchor=$ANCHOR root=$ROOT"
git -C "$REPO" log -1 --format='anchor_full=%H %ad' "$ANCHOR"
mkdir -p "$ROOT"
git -C "$REPO" archive "$ANCHOR" $EIGHT | tar -x -C "$ROOT" || { echo "archive_eight_rc=NONZERO"; exit 1; }
git -C "$REPO" archive "$ANCHOR" docs | tar -x -C "$ROOT" || { echo "archive_docs_rc=NONZERO"; exit 1; }
mkdir -p "$ROOT/build" "$ROOT/third_party/sherpa-onnx" "$ROOT/frontend/dist"
cp "$REPO"/build/*.dll "$ROOT/build/" 2>/dev/null
cp "$REPO"/third_party/sherpa-onnx/*.dll "$ROOT/third_party/sherpa-onnx/" 2>/dev/null
cp -r "$REPO"/frontend/dist/. "$ROOT/frontend/dist/" 2>/dev/null
echo "go_files=$(find "$ROOT/internal" "$ROOT/cmd" "$ROOT/tools" -name '*.go' | wc -l)"
echo "docs_plan=$([ -f "$ROOT/docs/PLAN.md" ] && echo present || echo MISSING)"
echo "dlls_build=$(ls "$ROOT"/build/*.dll 2>/dev/null | wc -l) dlls_sherpa=$(ls "$ROOT"/third_party/sherpa-onnx/*.dll 2>/dev/null | wc -l) dist=$([ -d "$ROOT/frontend/dist" ] && echo present || echo MISSING)"
cd "$(dirname "$ROOT")" || exit 1
echo "alias83=$(/c/Windows/System32/cmd.exe //c 'dir /X' 2>/dev/null | grep -i "$(basename "$ROOT")" | head -1)"
echo "blob_blacklist=$(git -C "$REPO" rev-parse "$ANCHOR:internal/risk/blacklist.go")"
echo "build_tree_rc=0"
