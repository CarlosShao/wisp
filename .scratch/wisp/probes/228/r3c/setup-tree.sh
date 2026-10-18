#!/bin/sh
# 228-r3c stage 1: build the OUT-OF-REPO export tree in the sanctioned landing
# path (D:\tmp\wisp\20261018\228r3c\tree == /d/tmp/wisp/20261018/228r3c/tree).
# Never run go test / go build in the shared worktree.
set -u
ANCHOR="$1"
ROOT=/d/tmp/wisp/20261018/228r3c
TREE=$ROOT/tree
REPO='/d/work/workspace/projects plans/Wisp'
LOG=$ROOT/logs-setup.txt
: >"$LOG"
say() { echo "$@" >>"$LOG"; echo "$@"; }

cd "$REPO" || { say "cd repo FAILED"; exit 9; }
say "anchor=$ANCHOR"
say "repo=$(git rev-parse HEAD)"

mkdir -p "$TREE" "$ROOT/gotmp" "$ROOT/tdst"
say "mkdir tree/gotmp/tdst rc=$?"

# 八件套 + 仓根 docs/
git archive "$ANCHOR" go.mod go.sum internal cmd tools scripts frontend models docs \
  | tar -x -C "$TREE"
say "git-archive|tar rc=$?"

# 手工三件（都非跟踪）
mkdir -p "$TREE/build" "$TREE/third_party/sherpa-onnx"
cp -f build/*.dll "$TREE/build/" 2>>"$LOG"; say "cp build/*.dll rc=$?"
cp -f third_party/sherpa-onnx/*.dll "$TREE/third_party/sherpa-onnx/" 2>>"$LOG"; say "cp sherpa dll rc=$?"
cp -rf frontend/dist "$TREE/frontend/" 2>>"$LOG"; say "cp frontend/dist rc=$?"

say "--- build dlls in tree ---"; ls -1 "$TREE/build" | grep -c '\.dll$' >>"$LOG"
say "--- sherpa dlls in tree ---"; ls -1 "$TREE/third_party/sherpa-onnx" | grep -c '\.dll$' >>"$LOG"
say "--- frontend/dist present ---"; ls -1 "$TREE/frontend/dist" | head -3 >>"$LOG"

# 取出后先数 .go 证非空
say "--- .go count in tree ---"
find "$TREE" -name '*.go' -type f | wc -l >>"$LOG"
say "--- docs/PLAN.md present ---"; test -f "$TREE/docs/PLAN.md" && say "docs/PLAN.md YES" >>"$LOG"
say "DONE setup"
