#!/bin/sh
# 342-v2 必答(iv) 本机档：./internal/winsec/ 全包成对四数（Windows 出货平台），基线树⇄落地树各 2 发。
# ★带 -v；★仓外树跑 .exe 取 DLL 用 cygpath -w＋';'（本包无 cgo，仍按定式带）。
set -u
ROOTW='D:\tmp\wisp\20261019\342v2'
LOGS=/d/tmp/wisp/20261019/342v2/logs
for TREE in treeB treeL; do
  for SHOT in 1 2; do
    LABEL="iv02-${TREE}-winfullpkg-${SHOT}"
    DLLP=$(cygpath -w "$ROOTW\\$TREE\\build"); SHP=$(cygpath -w "$ROOTW\\$TREE\\third_party\\sherpa-onnx")
    ( cd "$ROOTW\\$TREE" && PATH="$SHP;$DLLP;$PATH" GOCACHE='D:\tmp\wisp\20261019\342v2\gocache-win' \
        go test ./internal/winsec/ -v -count=1 -timeout 40m ) > "$LOGS/$LABEL.raw.txt" 2>&1
    RC=$?
    RUN=$(grep -c '^=== RUN' "$LOGS/$LABEL.raw.txt" || true)
    PT=$(grep -c '^--- PASS' "$LOGS/$LABEL.raw.txt" || true)
    FT=$(grep -c '^--- FAIL' "$LOGS/$LABEL.raw.txt" || true)
    SK=$(grep -c '^--- SKIP' "$LOGS/$LABEL.raw.txt" || true)
    echo "$LABEL RC=$RC RUN=$RUN PASS_TOP=$PT FAIL_TOP=$FT SKIP=$SK" | tee -a "$LOGS/four-counts-windows-fullpkg.tsv"
    grep -E '^--- FAIL' "$LOGS/$LABEL.raw.txt" | sed 's/^--- FAIL: //' | sort > "$LOGS/$LABEL.rednames.txt"
    echo "  RED_NAMES=$(wc -l < "$LOGS/$LABEL.rednames.txt")"; cat "$LOGS/$LABEL.rednames.txt"
  done
done
