#!/bin/sh
# 编排者收件复跑（336-r1 决定性那一发＝交件锚确认三发串行）：⛔ 抄腿的读数
# ★腿的台件把 LOGDIR 写死在它自己的写面里 ⇒ 我复跑⛔ 往它那儿落日志（会污染它的收件面＋造"谁写的"歧义）
#   ⇒ 复制一份、只改 LOGDIR 到我自己的目录，其余逐字不动（⛔ 改判据逻辑）
set -u
REPO="/d/work/workspace/projects plans/Wisp"
W="/d/tmp/wisp/20261018/orch/336-verify"
ANCHOR=e5021339
mkdir -p "$W/tree" "$W/gotmp" "$W/logs"
cd "$REPO" || exit 9
sed 's#^LOGDIR=.*#LOGDIR="'"$W"'/logs"#' .scratch/wisp/probes/336/r1/logs/run-confirm-anchor.sh > "$W/run-orch.sh"
echo "=== LOGDIR line in my copy ==="; grep -n '^LOGDIR=' "$W/run-orch.sh"
{
echo "### 编排者复跑 336-r1 交件锚确认发 anchor=$ANCHOR"
git archive "$ANCHOR" go.mod go.sum internal cmd tools scripts frontend models docs | tar -x -C "$W/tree"
echo "go files = $(find "$W/tree" -name '*.go' | wc -l)"
mkdir -p "$W/tree/build" "$W/tree/third_party/sherpa-onnx"
cp -f build/*.dll "$W/tree/build/" 2>/dev/null || echo "WARN build dll"
cp -f third_party/sherpa-onnx/*.dll "$W/tree/third_party/sherpa-onnx/" 2>/dev/null || echo "WARN sherpa dll"
ANCHOR="$ANCHOR" TREE="$W/tree" GOTMP="$W/gotmp" sh "$W/run-orch.sh"
} > "$W/logs/orch-336-recheck.log" 2>&1
echo "LOG=$W/logs/orch-336-recheck.log"
tail -45 "$W/logs/orch-336-recheck.log"
