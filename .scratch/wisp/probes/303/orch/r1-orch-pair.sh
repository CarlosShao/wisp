#!/usr/bin/env bash
# 编排者复跑：票 303 `AC#1` 定向两发（自己打一遍，⛔ 引腿的自报）
# 台件＝腿留下的仓外 clone：$HOME/wisp-303-bisect（clone 里没有 build/ 与 third_party/ 的 DLL，
# PATH 必须指回母仓，且必须 shell 形式 /d/…，D:/… 那形会 0xc0000135）
set -u
REPO="/d/work/workspace/projects plans/Wisp"
CLONE="$HOME/wisp-303-bisect"
OUT="$REPO/.scratch/wisp/probes/303/orch"
RUNAME='TestAC14AwaitedBindingReplyReachesThePage'
mkdir -p "$OUT"
export PATH="$REPO/third_party/sherpa-onnx:$REPO/build:$PATH"

echo "date_start=$(date --iso-8601=seconds)"
echo "clone_dir=$CLONE"
git -C "$CLONE" rev-parse --short HEAD
echo "== 复跑前工作树脏枚数（clone 内）=="
git -C "$CLONE" status --porcelain | wc -l

# (a) 归因那一笔：应当红，且红句逐字含判据句
git -C "$CLONE" checkout -q --detach fb2fb802 > "$OUT/r1-checkout-culprit.txt" 2>&1
echo "rc_checkout_a=$?"
echo "HEAD_a=$(git -C "$CLONE" rev-parse --short HEAD)"
go -C "$CLONE" test -count=1 -timeout 420s -v -run "$RUNAME" ./cmd/wisp/ > "$OUT/r1-pair-a-culprit.txt" 2>&1
echo "rc_pair_a=$?"
grep -c -e "--- FAIL: $RUNAME" "$OUT/r1-pair-a-culprit.txt"
grep -e 'no report "ac14r-0"' "$OUT/r1-pair-a-culprit.txt" | head -1

# (b) 那一笔的父发：应当绿，且逐字带页面自己的话
git -C "$CLONE" checkout -q --detach f718e9b6 > "$OUT/r1-checkout-parent.txt" 2>&1
echo "rc_checkout_b=$?"
echo "HEAD_b=$(git -C "$CLONE" rev-parse --short HEAD)"
go -C "$CLONE" test -count=1 -timeout 420s -v -run "$RUNAME" ./cmd/wisp/ > "$OUT/r1-pair-b-parent.txt" 2>&1
echo "rc_pair_b=$?"
grep -c -e "--- PASS: $RUNAME" "$OUT/r1-pair-b-parent.txt"
grep -e "page's own words" "$OUT/r1-pair-b-parent.txt" | head -1

# 收尾：clone 停在母仓当前分支等价位（下次程要用就现量，⛔ 拿旧状态当已知）
git -C "$CLONE" checkout -q --detach origin/dev > "$OUT/r1-checkout-restore.txt" 2>&1
echo "rc_restore=$?"
echo "HEAD_final=$(git -C "$CLONE" rev-parse --short HEAD)"
echo "date_end=$(date --iso-8601=seconds)"
