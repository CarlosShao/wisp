#!/bin/sh
# 352-r1 台面：自 `git archive` 建一枚仓外导出树（锚＝现认 commit），★⛔ 复用别腿的树、⛔ 碰共享树字节。
# 八件套＋仓根 docs/＋手工补三件（build/*.dll、third_party/sherpa-onnx/*.dll、frontend/dist＝都⛔ 跟踪）。
# ⛔ 删任何东西（临时件只建不删）。取出后先数 .go 证非空（"0"⛔ 是读数）。
# 用法：sh 01-build-tree.sh <anchor> <outname> <HEAD|WORKTREE>
#   WORKTREE＝把共享树当前的 internal/winsec/winsec.go 拷进去（＝落地那一笔的物理字节，入库前的自证用）
set -u
REPO='/d/work/workspace/projects plans/Wisp'
ROOT=/d/tmp/wisp/20261020/352r1
ANCHOR="$1"
OUT="$2"
WINSRC="${3:-HEAD}"
cd "$REPO" || exit 9
PIECES='go.mod go.sum internal cmd tools scripts frontend models docs'
T="$ROOT/$OUT"
mkdir -p "$T"
git archive --format=tar "$ANCHOR" $PIECES | tar -x -C "$T"
case "$WINSRC" in
  HEAD) : ;;
  WORKTREE) cp -f internal/winsec/winsec.go "$T/internal/winsec/winsec.go" ;;
  *) echo "BAD_WINSRC=$WINSRC"; exit 8 ;;
esac
mkdir -p "$T/build" "$T/third_party/sherpa-onnx" "$T/frontend/dist"
cp -f build/*.dll "$T/build/" 2>/dev/null || true
cp -f third_party/sherpa-onnx/*.dll "$T/third_party/sherpa-onnx/" 2>/dev/null || true
echo "ANCHOR=$ANCHOR FULL=$(git rev-parse "$ANCHOR")"
echo "TREE=$T"
echo "GO_FILES=$(find "$T" -name '*.go' | wc -l)"
echo "WINSEC_GO_FILES=$(find "$T/internal/winsec" -name '*.go' | wc -l)"
echo "--- wrapPath 那一枚格式串（内容锚，⛔ 绝对行号）---"
grep -n 'ErrNotSealable, path, err' "$T/internal/winsec/winsec.go"
echo "FMT_V=$(grep -c 'fmt.Errorf("%w: %s: %v", ErrNotSealable, path, err)' "$T/internal/winsec/winsec.go" || true)"
echo "FMT_W=$(grep -c 'fmt.Errorf("%w: %s: %w", ErrNotSealable, path, err)' "$T/internal/winsec/winsec.go" || true)"
echo "--- pristine：仓外留本腿自己的拷贝（还原靠它，⛔ git checkout）---"
cp -f "$T/internal/winsec/winsec.go" "$ROOT/pristine-winsec-$OUT.go"
sha256sum "$T/internal/winsec/winsec.go" "$ROOT/pristine-winsec-$OUT.go"
echo "BLOB_AT_ANCHOR=$(git rev-parse "$ANCHOR":internal/winsec/winsec.go)"
echo "WORKTREE_BLOB=$(git hash-object internal/winsec/winsec.go)"
