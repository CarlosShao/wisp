#!/bin/sh
# 台损更正：票 252 那一族用 shortName252(仓根) 作前置，仓根目录名必须产生 8.3 别名。
# 上一发我把导出树命名为 tree（4 字符合 8.3＝短名⇄长名同串）＝我自己造的台损。
# 这一发＝同锚重导一棵名字会出别名的树，供 AC#5 决定性读数用。
set -e
ANCHOR="$1"
REPO='/d/work/workspace/projects plans/Wisp'
ROOT='/d/tmp/wisp/20261019/227v2'
TREE="$ROOT/deliv-tree-ac5"
cd "$REPO"
mkdir -p "$TREE/build" "$TREE/third_party/sherpa-onnx"
git archive "$ANCHOR" go.mod go.sum internal cmd tools scripts frontend models docs | tar -x -C "$TREE"
for f in build/onnxruntime.dll build/sherpa-onnx-c-api.dll build/sherpa-onnx-cxx-api.dll; do
  cp -p "$REPO/$f" "$TREE/$f"
done
for f in "$REPO"/third_party/sherpa-onnx/*.dll; do
  cp -p "$f" "$TREE/third_party/sherpa-onnx/$(basename "$f")"
done
cp -rp "$REPO/frontend/dist" "$TREE/frontend/dist"
echo "ANCHOR=$ANCHOR"
echo "TREE=$TREE"
echo "go_files=$(find "$TREE" -name '*.go' | wc -l)"
echo "dll_build=$(ls "$TREE"/build/*.dll | wc -l) dll_third=$(ls "$TREE"/third_party/sherpa-onnx/*.dll | wc -l)"
echo "dist_files=$(find "$TREE/frontend/dist" -type f | wc -l)"
