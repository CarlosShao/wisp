#!/bin/sh
# 353-r1 台面：建一枚仓外导出树（★⛔ 复用别腿的树、⛔ 碰共享工作树字节、⛔ 在共享树跑 go test）。
# 八件套＋仓根 docs/＋手工补三件（build/*.dll、third_party/sherpa-onnx/*.dll、frontend/dist＝都⛔ 跟踪）。
# 用法：sh 01-build-tree.sh <anchor> <outname>
#   本腿的仪器件⛔ 在锚里（未 commit 前）⇒ 由调用方把物理字节 cp 进树（★打在仓外物理字节上）。
# ⛔ 删任何东西（临时件只建不删）。取出后先数 .go 证非空（"0"⛔ 是读数）。
set -u
REPO='/d/work/workspace/projects plans/Wisp'
ROOT=/d/tmp/wisp/20261020/353r1
ANCHOR="$1"
OUT="$2"
cd "$REPO" || exit 9
PIECES='go.mod go.sum internal cmd tools scripts frontend models docs'
T="$ROOT/$OUT"
mkdir -p "$ROOT/logs" "$T"
git archive --format=tar "$ANCHOR" $PIECES | tar -x -C "$T"
mkdir -p "$T/build" "$T/third_party/sherpa-onnx" "$T/frontend/dist"
cp -f build/*.dll "$T/build/" 2>/dev/null || true
cp -f third_party/sherpa-onnx/*.dll "$T/third_party/sherpa-onnx/" 2>/dev/null || true
echo "ANCHOR=$ANCHOR FULL=$(git rev-parse "$ANCHOR")"
echo "TREE=$T"
echo "GO_FILES=$(find "$T" -name '*.go' | wc -l)"
echo "WINSEC_GO_FILES=$(find "$T/internal/winsec" -name '*.go' | wc -l)"
echo "--- 承重三枚 blob ＋ 三枚本腿会落的件（内容锚，⛔ 绝对行号）---"
git rev-parse "$ANCHOR":internal/winsec/winsec.go | sed 's/^/BLOB_WINSEC_GO=/'
git rev-parse "$ANCHOR":internal/winsec/pinnedleaf_342_windows.go | sed 's/^/BLOB_PINNEDLEAF=/'
git rev-parse "$ANCHOR":internal/winsec/winsec_windows.go | sed 's/^/BLOB_WINSEC_WINDOWS=/'
echo "ANCHOR_HAS_353_TEST=$(test -f "$T/internal/winsec/nametopup_353_windows_test.go" && echo yes || echo no)"
echo "--- pristine：仓外留本腿自己的拷贝（还原靠它，⛔ git checkout）---"
for f in winsec.go pinnedleaf_342_windows.go; do
  cp -f "$T/internal/winsec/$f" "$ROOT/pristine-$OUT-$f"
done
sha256sum "$ROOT/pristine-$OUT-winsec.go" "$ROOT/pristine-$OUT-pinnedleaf_342_windows.go" \
  "$T/internal/winsec/pinnedleaf_342_windows.go"
