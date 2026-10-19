#!/bin/sh
# 343-v1 自己的 8 枚牙复打：env 指我自己 git archive 的树（treeDel＝31d4a476＝HEAD 交件态）。
# 台件＝.scratch/wisp/probes/336/v1/sh/run-ac2-sweep.sh 盘上原件，一字⛔ 动。
set -u
unset GOTMPDIR TMPDIR TEMP TMP
REPO='/d/work/workspace/projects plans/Wisp'
V1LOG=/d/tmp/wisp/20261019/343v1/logs/teeth
mkdir -p "$V1LOG"
mkdir -p /d/tmp/wisp/20261019/343v1/gotmp
export ANCHOR=31d4a47626ec7997666b6e0b190a760cdc535da0
export TREE='D:/tmp/wisp/20261019/343v1/treeDel'
export LOGDIR="$V1LOG"
export GOTMPROOT=/d/tmp/wisp/20261019/343v1/gotmp
echo "SweepTreeStart"
sh "$REPO/.scratch/wisp/probes/336/v1/sh/run-ac2-sweep.sh"
echo "SWEEP_RC=$?"
echo "harness_sha256=$(sha256sum "$REPO/.scratch/wisp/probes/336/v1/sh/run-ac2-sweep.sh" | cut -d' ' -f1)"
echo "harness_blob=$(git -C "$REPO" rev-parse HEAD:.scratch/wisp/probes/336/v1/sh/run-ac2-sweep.sh)"
