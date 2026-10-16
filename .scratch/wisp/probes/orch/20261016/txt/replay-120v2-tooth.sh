#!/bin/sh
# 编排者逐字重放 120-v2 的 tooth.sh（三发去牙＋基线＋还原复确认）
# ★我只改了台面三行（腿的原树 /d/tmp/probe120v2 与它的 driver.log 一字不动，⛔ 覆盖别人的件）：
#   V、LOG、BAK 所在目录，加 docker 卷名。突变逻辑⛔ 动一行。
# ★我自己的树用 mktemp -d 新建（不删任何人的东西）。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
D=$(mktemp -d)
V="$D/head"
LOG="$D/driver.log"
mkdir -p "$V"
( cd "$REPO" && git archive HEAD go.mod go.sum internal cmd tools scripts frontend ) | tar -x -C "$V"
NGO=$(find "$V" -name '*.go' | wc -l)
echo "tree=$V go_files=$NGO" > "$LOG"
echo "winsec_other.go sha256=$(sha256sum "$V/internal/winsec/winsec_other.go" | cut -d' ' -f1)" >> "$LOG"

sed -e "s#^V=/d/tmp/probe120v2/head#V=$V#" \
    -e "s#^LOG=/d/tmp/probe120v2/v2-tooth-driver.log#LOG=$LOG#" \
    -e "s#BAK=\"\\\$F.bak-120v2-teeth-on\"#BAK=\"$D/winsec_other.go.bak-orch\"#" \
    -e "s#probe120v2build#orch120v2build#" \
    "$REPO/.scratch/wisp/probes/120/v2/tooth.sh" > "$D/tooth-orch.sh"
echo "--- diff vs leg script (只应见台面三行) ---" >> "$LOG"
diff "$REPO/.scratch/wisp/probes/120/v2/tooth.sh" "$D/tooth-orch.sh" >> "$LOG" 2>&1
cp "$LOG" /d/tmp/orch-120v2-replay-head.txt 2>/dev/null
echo "TREE=$D"
sh "$D/tooth-orch.sh"
echo "ORCHEXIT=$?" >> "$LOG"
echo '===== driver tail ====='
tail -45 "$LOG"
echo "LOGPATH=$LOG"
