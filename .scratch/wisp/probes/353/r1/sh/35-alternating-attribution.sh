#!/bin/sh
# 353-r1 连带定因的★决定性尺＝两侧**交替**打同一发（⛔ 先打完一側再打另一侧＝那形会把"台面负载随时间变化"误当成字节差）。
# 只打 A 族那两枚 400 发 racer，各 3 轮交替，取 provedReparse／flipErrs／truncated／deleted。
# ★量的＝这三枚分子；★算的＝"随字节走"还是"随时段走"的判语（⛔ 跑计时）。
set -u
R=/d/tmp/wisp/20261020/353r1
ROOTW='D:\tmp\wisp\20261020\353r1'
L="$R/logs"
PAT='TestAC2WindowsPrivateFileThroughASwappedAncestorNeverDeletesAForeignFile342r1|TestAC2WindowsLeafIdentityReCheckRefusesAForeignNameAndNeverDeletesIt342r1'
for round in 1 2 3; do
  for TREE in finalB finalL finalL finalB; do
    OUT="$L/alt-$round-$TREE.raw.txt"
    TS=$(date +%Y%m%d-%H%M%S)-$$
    GOTMPW="${LOCALAPPDATA:-C:\\Users\\swq\\AppData\\Local}\\wisp-test-gotmp\\353r1alt-$TS"
    mkdir -p "$(cygpath -u "$GOTMPW")" || { echo "CANNOT_CREATE_GOTMP"; exit 7; }
    DLLP=$(cygpath -w "$ROOTW\\$TREE\\build"); SHP=$(cygpath -w "$ROOTW\\$TREE\\third_party\\sherpa-onnx")
    ( cd "$ROOTW\\$TREE" && PATH="$SHP;$DLLP;$PATH" GOTMPDIR="$GOTMPW" TMPDIR="$GOTMPW" \
        GOCACHE="$ROOTW\\gocache-win" go test -count=1 -timeout 40m -run "$PAT" ./internal/winsec/ ) > "$OUT" 2>&1
    RC=$?
    L1=$(grep -oE 'form=leaf-symlink attempts=400 cycles=[0-9]+ provedReparse=[0-9]+ formBLegRefusals=[0-9]+ truncated=[0-9]+ deleted=[0-9]+ returnedNil=[0-9]+ flipErrs=[0-9]+' "$OUT" | head -1)
    D1=$(grep -oE 'form=symlink-dir attempts=400 cycles=[0-9]+ provedReparse=[0-9]+ formBLegRefusals=[0-9]+ truncated=[0-9]+ deleted=[0-9]+ returnedNil=[0-9]+ flipErrs=[0-9]+' "$OUT" | head -1)
    echo "ALT r=$round tree=$TREE RC=$RC"
    echo "     leaf-symlink: ${L1:-<none>}"
    echo "     symlink-dir : ${D1:-<none>}"
    echo "ALT r=$round tree=$TREE ${L1:-none} | ${D1:-none}" >> "$L/alt-summary.txt"
  done
done
