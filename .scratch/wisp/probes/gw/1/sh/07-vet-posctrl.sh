#!/bin/sh
# gw-1 味① 第 5 把的**正控**＝"那枚 rc=0 是验过的读数"：
# 在**仓外**副本 cmd/wisp 里造一枚坏 printf（vet 的 printf analyzer 射程内的真错误形），
# 同一把尺在副本上必须 rc=1 且逐字打印诊断＝证明这把尺在该射程有牙。
# ★⛔ 碰共享工作树任何字节；★副本＝gwL（锚 a19d5bf2 的 git archive 导出树＝今天 HEAD 的字节）。
set -u
R=/d/tmp/wisp/20261020/gw1
LOGS="$R/logs"
T="$R/gwL"
MUT="$R/gwL-vetmut"
unset GOTMPDIR TMPDIR TEMP TMP
echo "== pristine 锚 ⇄ 副本 =="
git -C '/d/work/workspace/projects plans/Wisp' rev-parse HEAD | sed 's|^|HEAD_FULL=|'
[ -d "$MUT" ] || cp -a "$T" "$MUT"
echo "MUT_TREE=$MUT"
echo "GO_FILES=$(find "$MUT" -name '*.go' | wc -l)"
TARGET="$MUT/cmd/wisp/vet_posctrl_gw1.go.txt"
F="$MUT/cmd/wisp/vet_posctrl_gw1.go"
echo "TARGET=$TARGET"
printf 'package main\n\nimport "fmt"\n\n// vetPosCtrlGW1 exists only to prove go vet -tags winlive has teeth\n// in the printf analyzer range on this box.\nfunc vetPosCtrlGW1() {\n\tfmt.Printf("%s %s\\n", 1)\n}\n' > "$TARGET"
echo "PRE_TEST_RC_UNUSED=placeholder"
echo "PRE_LIST=$(ls "$MUT/cmd/wisp" | grep -c 'vet_posctrl_gw1')"
echo "SUFFIXED=.go.txt＝⛔ 编译面＝此时 vet 必须照 0（负档：证明那枚 0 ⛔ 是被新件搅出来的）"
TW=$(cygpath -w "$MUT")
DLLP=$(cygpath -w "$MUT/build"); SHP=$(cygpath -w "$MUT/third_party/sherpa-onnx")
( cd "$TW" && PATH="$SHP;$DLLP;$PATH" GOFLAGS= go vet -tags winlive ./cmd/wisp/ ) > "$LOGS/vet-posctrl-negative.txt" 2>&1
echo "VET_NEG_RC=$?"
echo "VET_NEG_BYTES=$(wc -c < "$LOGS/vet-posctrl-negative.txt")"
sed 's|^|  NEG |' "$LOGS/vet-posctrl-negative.txt"
echo
echo "== 落刀＝把 .go.txt 改成 .go（⛔ 跟踪面＝全在仓外副本）=="
mv "$TARGET" "$F"
echo "MOVED=$F  EXISTS=$(test -f "$F" && echo yes)"
( cd "$TW" && PATH="$SHP;$DLLP;$PATH" GOFLAGS= go vet -tags winlive ./cmd/wisp/ ) > "$LOGS/vet-posctrl-positive.txt" 2>&1
echo "VET_POS_RC=$?"
echo "VET_POS_BYTES=$(wc -c < "$LOGS/vet-posctrl-positive.txt")"
sed 's|^|  POS |' "$LOGS/vet-posctrl-positive.txt"
echo
echo "== 还原＝删掉副本里那枚文件（只动副本；⛔ git checkout／⛔ 动共享树）=="
mv "$F" "$TARGET"
( cd "$TW" && PATH="$SHP;$DLLP;$PATH" GOFLAGS= go vet -tags winlive ./cmd/wisp/ ) > "$LOGS/vet-posctrl-restored.txt" 2>&1
echo "VET_RESTORE_RC=$?"
echo "VET_RESTORE_BYTES=$(wc -c < "$LOGS/vet-posctrl-restored.txt")"
sed 's|^|  RES |' "$LOGS/vet-posctrl-restored.txt"
echo
echo "== 还原哈希对回＝副本 cmd/wisp 与 gwL cmd/wisp 逐枚 cmp（差异枚数须 0）=="
diff -rq "$T/cmd/wisp" "$MUT/cmd/wisp" > "$LOGS/vet-posctrl-restore-diff.txt" 2>&1
echo "RESTORE_DIFF_RC=$?"
echo "RESTORE_DIFF_LINES=$(grep -c . "$LOGS/vet-posctrl-restore-diff.txt" || true)"
grep . "$LOGS/vet-posctrl-restore-diff.txt" | sed 's|^|  DIFF |'
echo
echo "== 共享树复量（⛔ 任何字节改动；再打一次那把 vet 证其 0 复现）=="
cd '/d/work/workspace/projects plans/Wisp' || exit 9
GOFLAGS= go vet -tags winlive ./cmd/wisp/ > "$LOGS/vet-sharedtree-repeat.txt" 2>&1
echo "VET_SHARED_REPEAT_RC=$?"
echo "VET_SHARED_REPEAT_BYTES=$(wc -c < "$LOGS/vet-sharedtree-repeat.txt")"
git rev-parse HEAD | sed 's|^|HEAD_FULL_END=|'
