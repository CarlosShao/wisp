#!/bin/sh
# 353-r1 摆动名册定因尺＝对那一枚 342/348 族的 racer 仪器（⛔ 本票的东西、400 发定时尝试＝非确定性）
#   在【finalB 基线】⇄【finalL 落地】两侧各连打 N 发，只数"哪一发响＋响在哪一支"。
# ★公共块＝撞红先安静复量、竞态读数⛔ 要求逐字相同（复现的是形状⇄量级）；★串行。
set -u
R=/d/tmp/wisp/20261020/353r1
ROOTW='D:\tmp\wisp\20261020\353r1'
LOGS="$R/logs"
N="${1:-4}"
PAT='TestAC2WindowsPrivateFileThroughASwappedAncestorNeverDeletesAForeignFile342r1|TestAC1WindowsLeafRecheckAnswersOnTheBareSiblingForm342r2|TestAC2WindowsLeafIdentityReCheckRefusesAForeignNameAndNeverDeletesIt342r1'
for TREE in finalB finalL; do
  for s in $(seq 1 "$N"); do
    OUT="$LOGS/swing-$TREE-$s.raw.txt"
    TS=$(date +%Y%m%d-%H%M%S)-$$
    GOTMPW="${LOCALAPPDATA:-C:\\Users\\swq\\AppData\\Local}\\wisp-test-gotmp\\353r1sw-$TS"
    mkdir -p "$(cygpath -u "$GOTMPW")" || { echo "CANNOT_CREATE_GOTMP"; exit 7; }
    DLLP=$(cygpath -w "$ROOTW\\$TREE\\build"); SHP=$(cygpath -w "$ROOTW\\$TREE\\third_party\\sherpa-onnx")
    ( cd "$ROOTW\\$TREE" && PATH="$SHP;$DLLP;$PATH" GOTMPDIR="$GOTMPW" TMPDIR="$GOTMPW" \
        GOCACHE="$ROOTW\\gocache-win" go test -v -count=1 -timeout 40m -run "$PAT" ./internal/winsec/ ) > "$OUT" 2>&1
    RC=$?
    echo "SWING $TREE shot=$s RC=$RC FAIL_TOP=$(grep -c '^--- FAIL' "$OUT" || true) deleted_GT0=$(grep -cE 'AC#2 RED: [0-9]+ of [0-9]+ attempts found the foreign witness GONE' "$OUT" || true) instrRed=$(grep -c 'RED instrument' "$OUT" || true)" | tee -a "$LOGS/swing-summary.txt"
    grep -E '^--- FAIL' "$OUT" | sed 's/^--- FAIL: /   RED /; s/ *([0-9.]*s)$//' | tee -a "$LOGS/swing-summary.txt"
    grep -oE 'form=symlink-dir attempts=400 cycles=[0-9]+ provedReparse=[0-9]+ formBLegRefusals=[0-9]+ truncated=[0-9]+ deleted=[0-9]+ returnedNil=[0-9]+ flipErrs=[0-9]+' "$OUT" | head -1 | sed 's/^/   NUM /' | tee -a "$LOGS/swing-summary.txt"
  done
done
echo "== 汇总（★只数形状，⛔ 把单发当定论）=="
for TREE in finalB finalL; do
  echo "  $TREE shots=$N redshots=$(grep -c "^SWING $TREE .*FAIL_TOP=[1-9]" "$LOGS/swing-summary.txt" || true) damageGt0=$(grep -c 'found the foreign witness GONE' "$LOGS"/swing-$TREE-*.raw.txt 2>/dev/null | grep -v ':0' | wc -l)"
done
