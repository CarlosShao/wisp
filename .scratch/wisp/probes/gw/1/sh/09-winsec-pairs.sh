#!/bin/sh
# gw-1 味③ winsec 成对色驱动＝同一把尺（05-run-pair.sh）打两对锚：
#   对甲＝3a509e50（乙本体落地前那一笔＝gwW）⇄今天 HEAD（gwL）      ＝票 353 AC#3 那对
#   对乙＝0fb399c3（352 那一行的父＝gwB）⇄今天 HEAD（gwL）           ＝味②票 352 AC#2 的具名连带侧
# ★⛔ 整包 ./cmd/wisp（票 352 票面逐字归编排者那窗）；★⛔ ./... ；★⛔ 共享工作树跑 go test
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 9
echo "== driver START anchor =="
git rev-parse HEAD | sed 's|^|HEAD_FULL=|'
cd .scratch/wisp/probes/gw/1 || exit 9
echo
echo "### 对甲 host（gwW ⇄ gwL）###"
sh sh/05-run-pair.sh host p353 gwW gwL ./internal/winsec/
echo
echo "### 对乙 host（gwB ⇄ gwL）###"
sh sh/05-run-pair.sh host p352w gwB gwL ./internal/winsec/
echo
echo "### 对甲 container（gwW ⇄ gwL）###"
sh sh/05-run-pair.sh container c353 gwW gwL ./internal/winsec/
echo
echo "### 对乙 container（gwB ⇄ gwL）###"
sh sh/05-run-pair.sh container c352w gwB gwL ./internal/winsec/
echo
echo "== driver END anchor =="
git rev-parse HEAD | sed 's|^|HEAD_FULL_END=|'
