#!/usr/bin/env bash
# 334-r1 台面配方（派单"台面配方"那一节）：仓外导出树 = 八件套 + 仓根 docs/ +
# 三件手工补（build/*.dll、third_party/sherpa-onnx/*.dll、frontend/dist，全部不跟踪）。
# 树根目录名故意超过 8 字符，好让 NTFS 生成 8.3 短名别名（internal/tools 那一族要读它）。
# ⛔ rm -rf：临时件只建不删（issues/README 规则 8），每发一棵新树。
# 用法：build-export-tree.sh <锚> <D 卷树根>
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
echo "go_files=$(find "$ROOT/internal" "$ROOT/cmd" "$ROOT/tools" -name '*.go' | wc -l) (取出后先数 .go 证非空)"
echo "top_level_dirs=$(ls "$ROOT" | tr '\n' ' ')"
echo "docs_plan=$([ -f "$ROOT/docs/PLAN.md" ] && echo present || echo MISSING)"
echo "dlls_build=$(ls "$ROOT"/build/*.dll 2>/dev/null | wc -l) dlls_sherpa=$(ls "$ROOT"/third_party/sherpa-onnx/*.dll 2>/dev/null | wc -l) dist=$([ -d "$ROOT/frontend/dist" ] && echo present || echo MISSING)"
# 8.3 短名别名（尺：在树根的【父目录】里跑 dir /X，看树根那一行的短名列；
# ⛔ cmd 走 PATH 会命中 npm 的同名 CLI，必须指真身 System32\cmd.exe）
cd "$(dirname "$ROOT")" || exit 1
echo "alias83=$(/c/Windows/System32/cmd.exe //c 'dir /X' 2>/dev/null | grep -i "$(basename "$ROOT")" | head -1)"
echo "build_export_tree_rc=0"
