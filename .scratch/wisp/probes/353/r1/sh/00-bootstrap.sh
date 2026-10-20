#!/bin/sh
# 353-r1 起手：现认锚＋建两棵树（treeB＝基线侧，乙未落；treeL＝落地侧，随后同步本体字节进去）。
# ★⛔ 复用别腿的树、★⛔ 在共享工作树跑 go test／go build。
set -u
ROOT=/d/tmp/wisp/20261020/353r1
cd '/d/work/workspace/projects plans/Wisp' || exit 9
ANCHOR=$(git rev-parse HEAD)
echo "ANCHOR=$ANCHOR"
echo "ANCHOR_SUBJECT=$(git log -1 --format='%s')"
sh "$ROOT/sh/01-build-tree.sh" "$ANCHOR" treeB
sh "$ROOT/sh/01-build-tree.sh" "$ANCHOR" treeL
