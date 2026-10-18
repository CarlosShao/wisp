#!/bin/bash
# 339-r1 交件侧台件：从锚新建一棵干净导出树，再把改后的两枚 .go 打进去。
# ⛔ 删任何件（临时件只建不删）⇒ 交件侧一律新目录，不复用基线那棵。
# 用法：bash apply.sh <linux|win> <tree-name>
set -u
REPO="/d/work/workspace/projects plans/Wisp"
KIND="$1"
NAME="$2"
ANCHOR=e6fd0f5051812a23c72db444d12d76b4d27285a7
ROOT="/d/tmp/wisp/20261018/339r1"
TREE="$ROOT/$KIND/$NAME"
cd "$REPO" || exit 9
mkdir -p "$TREE"
git archive "$ANCHOR" go.mod go.sum internal cmd tools scripts frontend models docs | tar -xf - -C "$TREE" || exit 8
cp internal/risk/anchor_miss_reason_334_test.go "$TREE/internal/risk/anchor_miss_reason_334_test.go"
cp internal/risk/anchor_miss_reason_339_other_test.go "$TREE/internal/risk/anchor_miss_reason_339_other_test.go"
echo "kind=$KIND tree=$TREE go=$(find "$TREE" -name '*.go' | wc -l)"
echo "risk_files: $(ls "$TREE/internal/risk" | grep -c 'anchor_miss_reason')"
echo "334_sha=$(sha1sum "$TREE/internal/risk/anchor_miss_reason_334_test.go" | cut -d' ' -f1)"
echo "339_sha=$(sha1sum "$TREE/internal/risk/anchor_miss_reason_339_other_test.go" | cut -d' ' -f1)"
