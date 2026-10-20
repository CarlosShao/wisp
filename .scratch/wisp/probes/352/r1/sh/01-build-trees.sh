#!/bin/sh
# ★本件＝已被 01-build-tree.sh 取代的早期草稿（一次建两棵），⛔ 未_used_于任何读数；
#   按"临时件只建不删"留在盘上入库当形状记录。真用尺＝01-build-tree.sh（单棵＋WORKTREE 支）。
# 352-r1 台面：自 `git archive` 建两棵仓外导出树（基线＝起手现认锚 ⇄ 落地＝那一行那笔 commit），
# ★⛔ 复用别腿的树、★⛔ 碰共享工作树字节。八件套＋仓根 docs/＋手工补三件（都⛔ 跟踪）。
# ⛔ 删任何东西（临时件只建不删）。取出后先数 .go 证非空（"0"⛔ 是读数）。
# 用法：sh 01-build-trees.sh <anchor> <outname> <winsec_go_source>
#   winsec_go_source = HEAD (整笔 commit 的字节) 或 WORKTREE (共享树当前字节，落地那一笔之前用)
set -u
REPO='/d/work/workspace/projects plans/Wisp'
ROOT=/d/tmp/wisp/20261020/352r1
ANCHOR="$1"
OUT="$2"
WINSRC="$3"
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
echo "TREE=$T"
echo "GO_FILES=$(find "$T" -name '*.go' | wc -l)"
echo "WINSEC_GO_FILES=$(find "$T/internal/winsec" -name '*.go' | wc -l)"
echo "--- winsec.go 那一枚格式串（内容锚，⛔ 绝对行号）---"
grep -n 'fmt.Errorf("%w: %s: %', "$T/internal/winsec/winsec.go" | head -5
echo "WRAP_PATH_FMT=$(grep -c 'fmt.Errorf("%w: %s: %v", ErrNotSealable, path, err)' "$T/internal/winsec/winsec.go" || true)"
echo "WRAP_PATH_W=$(grep -c 'fmt.Errorf("%w: %s: %w", ErrNotSealable, path, err)' "$T/internal/winsec/winsec.go" || true)"
sha256sum "$T/internal/winsec/winsec.go"
# 仓外留一枚本腿自己的 pristine 拷贝（还原靠它，⛔ git checkout）
cp -f "$T/internal/winsec/winsec.go" "$ROOT/pristine-winsec-$OUT.go"
sha256sum "$ROOT/pristine-winsec-$OUT.go"
echo "SHARED_TREE_winsec_blob=$(git rev-parse "$ANCHOR":internal/winsec/winsec.go)"
echo "SHARED_TREE_winsec_worksha=$(git hash-object internal/winsec/winsec.go)"
echo "GONE_RC_IF_NONEMPTY_ABOVE=ok"
