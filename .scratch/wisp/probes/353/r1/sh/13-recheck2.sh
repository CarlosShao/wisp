#!/bin/sh
# 353-r1 摆动复量尺（第二枚）＝定向 `-run` 在 treeB2⇄treeL 各连打 3 发。
# ★公共块＝撞红先安静复量 3 发再归因；竞态／非确定性读数⛔ 要求逐字相同（复现的是量级⇄形状）。
# 用法：sh 13-recheck2.sh <TestNameRegex>
set -u
PAT="$1"
ROOT=/d/tmp/wisp/20261020/353r1
ROOTW='D:\tmp\wisp\20261020\353r1'
LOGS="$ROOT/logs"
mkdir -p "$LOGS"
for TREE in treeB2 treeL; do
  for shot in 1 2 3; do
    OUT="$LOGS/recheck2-$TREE-$shot.raw.txt"
    TS=$(date +%Y%m%d-%H%M%S)-$$
    GOTMPW="${LOCALAPPDATA:-C:\\Users\\swq\\AppData\\Local}\\wisp-test-gotmp\\353r1-r2-$TS"
    mkdir -p "$(cygpath -u "$GOTMPW")" || { echo "CANNOT_CREATE_GOTMP=$GOTMPW"; exit 7; }
    DLLP=$(cygpath -w "$ROOTW\\$TREE\\build")
    SHP=$(cygpath -w "$ROOTW\\$TREE\\third_party\\sherpa-onnx")
    ( cd "$ROOTW\\$TREE" && PATH="$SHP;$DLLP;$PATH" GOTMPDIR="$GOTMPW" TMPDIR="$GOTMPW" \
        GOCACHE="$ROOTW\\gocache-win" go test -count=1 -timeout 40m -run "$PAT" ./internal/winsec/ ) \
      > "$OUT" 2>&1
    RC=$?
    {
      echo "RECHECK2 tree=$TREE shot=$shot RC=$RC RUN=$(grep -c '^=== RUN' "$OUT" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$OUT" || true) SKIP=$(grep -c '^--- SKIP' "$OUT" || true)"
      grep -E "leaf-leg split|AC#2 RED|no leaf leg" "$OUT" | head -3 | sed 's/^/    /'
    } | tee -a "$LOGS/recheck2-summary.txt"
  done
done
