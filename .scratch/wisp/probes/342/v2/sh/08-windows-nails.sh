#!/bin/sh
# 342-v2 必答(ii)：Windows 三形（junction ⇄ symlink-dir ⇄ leaf-symlink）× 基线树 treeB ⇄ 落地树 treeL，
# 各侧 3 发。本机 go1.27.1 windows/amd64；台面＝自己 git archive 的仓外导出树（单差异＝internal/winsec 落地字节）。
# ★GOTMPDIR/TMPDIR ⛔ 指到 bind mount／D 卷异常形；本机用默认 %LOCALAPPDATA%\Temp（同卷＝身份门灵）。
# ★仓外导出树跑 Windows .exe 取 DLL 需 cygpath -w ＋ ';'：本包无 cgo/DLL 依赖，仍按公共块带上（正控＝不带也不会死，带＝定式）。
set -u
ROOT='D:\tmp\wisp\20261019\342v2'
LOGS=/d/tmp/wisp/20261019/342v2/logs
mkdir -p "$LOGS"
WINSET='TestAC1WindowsPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1|TestAC2WindowsPrivateFileThroughASwappedAncestorNeverDeletesAForeignFile342r1|TestAC2WindowsLeafIdentityReCheckRefusesAForeignNameAndNeverDeletesIt342r1'
for TREE in treeB treeL; do
  for SHOT in 1 2 3; do
    LABEL="ii01-${TREE}-win-shot${SHOT}"
    DLLP=$(cygpath -w "$ROOT\\$TREE\\build"); SHP=$(cygpath -w "$ROOT\\$TREE\\third_party\\sherpa-onnx")
    ( cd "$ROOT\\$TREE" && PATH="$SHP;$DLLP;$PATH" GOCACHE='D:\tmp\wisp\20261019\342v2\gocache-win' \
        go test ./internal/winsec/ -v -count=1 -timeout 40m -run "$WINSET" ) > "$LOGS/$LABEL.raw.txt" 2>&1
    RC=$?
    RUN=$(grep -c '^=== RUN' "$LOGS/$LABEL.raw.txt" || true)
    PT=$(grep -c '^--- PASS' "$LOGS/$LABEL.raw.txt" || true)
    FT=$(grep -c '^--- FAIL' "$LOGS/$LABEL.raw.txt" || true)
    SK=$(grep -c '^--- SKIP' "$LOGS/$LABEL.raw.txt" || true)
    echo "$LABEL RC=$RC RUN=$RUN PASS_TOP=$PT FAIL_TOP=$FT SKIP=$SK" | tee -a "$LOGS/four-counts-windows.tsv"
    grep -E ': (AC#1|AC#2) \(ticket' "$LOGS/$LABEL.raw.txt" | sed -E 's/^ +[a-z0-9_]+\.go:[0-9]+: //' | head -12
    grep -E 'leg that answered|RED:|RED instrument|RESIDUAL|residual reading|leaf-leg split' "$LOGS/$LABEL.raw.txt" | sed -E 's/^ +[a-z0-9_]+\.go:[0-9]+: //' | head -40
    echo "-----"
  done
done
