#!/bin/bash
# 339-v1 台件①：仓外导出树（裁决腿自己建，⛔ 复用编排者／实现腿那几棵）。
# 用法：bash tree.sh <pre|post|winpost> <commit>
# 配方＝公共块八件套 + 仓根 docs/；落点一律 D:\tmp\wisp\20261018\339v1\...（⛔ 盘符根、⛔ 删）。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
KIND="$1"
ANCHOR="$2"
ROOT="/d/tmp/wisp/20261018/339v1"
TREE="$ROOT/$KIND/tree"
cd "$REPO" || exit 9
mkdir -p "$ROOT/$KIND"
git archive "$ANCHOR" go.mod go.sum internal cmd tools scripts frontend models docs > "$ROOT/$KIND/tree.tar" || exit 8
mkdir -p "$TREE" && tar -xf "$ROOT/$KIND/tree.tar" -C "$TREE" || exit 7
GOC=$(find "$TREE" -name '*.go' | wc -l)
echo "anchor=$ANCHOR kind=$KIND tree=$TREE go_files=$GOC risk_go=$(find "$TREE/internal/risk" -name '*.go' | wc -l)"
echo "TAGLINE_334=$(head -1 "$TREE/internal/risk/anchor_miss_reason_334_test.go")"
echo "TAGLINE_339=$(head -1 "$TREE/internal/risk/anchor_miss_reason_339_other_test.go" 2>/dev/null || echo ABSENT)"
echo "blob_blacklist=$(git -C "$REPO" rev-parse "$ANCHOR:internal/risk/blacklist.go")"
echo "blob_pathresolverother=$(git -C "$REPO" rev-parse "$ANCHOR:internal/risk/pathresolver_other.go")"
echo "blob_334=$(git -C "$REPO" rev-parse "$ANCHOR:internal/risk/anchor_miss_reason_334_test.go")"
