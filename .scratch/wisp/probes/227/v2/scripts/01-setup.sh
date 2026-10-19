#!/bin/sh
# 227-v2 台面起手：自建仓外导出树（不复用编排者 orch227-b 那棵）。
# 锚＝本腿自己 git log -1 现认的那枚 sha，由 argv[1] 传入。
set -e
ANCHOR="$1"
REPO='/d/work/workspace/projects plans/Wisp'
ROOT='/d/tmp/wisp/20261019/227v2'
TREE="$ROOT/tree"
mkdir -p "$TREE"
cd "$REPO"
echo "ANCHOR=$ANCHOR"
git rev-parse --verify "$ANCHOR^{commit}" >/dev/null
git archive "$ANCHOR" go.mod go.sum internal cmd tools scripts frontend models docs | tar -x -C "$TREE"
# 手工补三件（都不 tracked）：build/*.dll + third_party/sherpa-onnx/*.dll + frontend/dist
mkdir -p "$TREE/build" "$TREE/third_party/sherpa-onnx"
for f in build/onnxruntime.dll build/sherpa-onnx-c-api.dll build/sherpa-onnx-cxx-api.dll; do
  [ -f "$REPO/$f" ] && cp -p "$REPO/$f" "$TREE/$f" && echo "copied $f"
done
for f in "$REPO"/third_party/sherpa-onnx/*.dll; do
  b=$(basename "$f"); cp -p "$f" "$TREE/third_party/sherpa-onnx/$b"; echo "copied third_party/sherpa-onnx/$b"
done
cp -rp "$REPO/frontend/dist" "$TREE/frontend/dist" && echo "copied frontend/dist"
echo "go_files=$(find "$TREE" -name '*.go' | wc -l)"
echo "dll_build=$(ls "$TREE"/build/*.dll 2>/dev/null | wc -l)"
echo "dll_third=$(ls "$TREE"/third_party/sherpa-onnx/*.dll 2>/dev/null | wc -l)"
echo "dist_files=$(find "$TREE/frontend/dist" -type f | wc -l)"
echo "docs_plan=$([ -f "$TREE/docs/PLAN.md" ] && echo yes || echo no)"
# 短 GOTMPDIR 两枚根（C 卷＝%LOCALAPPDATA% 同卷；甲＝反斜杠拼法，乙＝正斜杠拼法＝编排者那发的形状）
mkdir -p '/c/Users/swq/AppData/Local/wisp-test-gotmp/227v2-a' '/c/Users/swq/AppData/Local/wisp-test-gotmp/227v2-slash'
echo "gotmp_a=$(ls -d '/c/Users/swq/AppData/Local/wisp-test-gotmp/227v2-a')"
echo SETUP_DONE
