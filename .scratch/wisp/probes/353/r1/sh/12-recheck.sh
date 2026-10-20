#!/bin/sh
# 353-r1 摆动复量尺：同一发定向 `-run` 在【treeB2（产码＝基线）】⇄【treeL（产码＝乙已落）】各连打 3 发。
# ★公共块＝撞红先安静复量 3 发再归因；★竞态／非确定性读数⛔ 要求逐字相同，复现的是量级与形状。
# 用法：sh 12-recheck.sh <TestNameRegex>
set -u
PAT="$1"
ROOT=/d/tmp/wisp/20261020/353r1
ROOTW='D:\tmp\wisp\20261020\353r1'
LOGS="$ROOT/logs"
mkdir -p "$LOGS"
for TREE in treeB2 treeL; do
  for shot in 1 2 3; do
    OUT="$LOGS/recheck-$TREE-$shot.raw.txt"
    TS=$(date +%Y%m%d-%H%M%S)-$$
    GOTMPW="${LOCALAPPDATA:-C:\\Users\\swq\\AppData\\Local}\\wisp-test-gotmp\\353r1-rc-$TS"
    mkdir -p "$(cygpath -u "$GOTMPW")" || { echo "CANNOT_CREATE_GOTMP=$GOTMPW"; exit 7; }
    DLLP=$(cygpath -w "$ROOTW\\$TREE\\build"); SHP=$(cygpath -w "$ROOTW\\$TREE\\third_party\\sherpa-onnx")
    ( cd "$ROOTW\\$TREE" && PATH="$SHP;$DLLP;$PATH" GOTMPDIR="$GOTMPW" TMPDIR="$GOTMPW" \
        GOCACHE="$ROOTW\\gocache-win" go test -count=1 -timeout 40m -run "$PAT" ./internal/winsec/ ) \
      > "$OUT" 2>&1
    RC=$?
    {
      echo "RECHECK tree=$TREE shot=$shot RC=$RC FAIL_TOP=$(grep -c '^--- FAIL' "$OUT" || true) RUN=$(grep -c '^=== RUN' "$OUT" || true)"
      grep -E 'truncated=|deleted=|provedReparse=|flipErrs=' "$OUT" | head -4
    } | tee -a "$LOGS/recheck-summary.txt"
  done
done
