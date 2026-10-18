#!/usr/bin/env bash
# 336-r1 颜色尺台件：只在【仓外导出台树】里跑 go test，⛔ 共享工作树。
# 用法：run-risk.sh <树根> <名册正则> <日志名> [--ci-shape|--dshape]
#   --ci-shape ＝ GOTMPDIR 逐字照 scripts/portable-tests.sh:152-159 的组法（$LOCALAPPDATA\wisp-test-gotmp\...）
#   --dshape   ＝ GOTMPDIR＝D:\tmp\wisp\20261018\336r1\gotmp\<日志名>（机主 10-18 新铁律：⛔ 盘符根目录）
# 四数尺：RUN=^=== RUN（含子测试）／PASS_TOP=^--- PASS（只顶层）／FAIL_TOP=^--- FAIL／SKIP=^--- SKIP，
# ⛔ 混用⛔ 相减；一律 -count=1 -timeout 30m。
set -u
REPO="D:/work/workspace/projects plans/Wisp"
TREE="${1:?树根}"
ROSTER="${2:?名册正则}"
LOGNAME="${3:?日志名}"
SHAPE="${4:---dshape}"
LOGDIR="$REPO/.scratch/wisp/probes/336/r1/logs"
GOTMPROOT="${GOTMPROOT:-/d/tmp/wisp/20261018/336r1/gotmp}"
export GOPATH='D:\work\base\gopath'
export PATH="$TREE/third_party/sherpa-onnx:$TREE/build:$PATH"
STAMP="$(date +%Y%m%d-%H%M%S)-$$"
case "$SHAPE" in
  --ci-shape)
    base="${LOCALAPPDATA:-C:\\Users\\swq\\AppData\\Local}\\wisp-test-gotmp"
    GOTMP="$base\\portable-336r1-$STAMP" ;;
  --dshape)
    GOTMP="$GOTMPROOT/$LOGNAME" ;;
  *) echo "unknown shape $SHAPE"; exit 2 ;;
esac
mkdir -p "$GOTMP" && mkdir_gotmp_rc=$? || mkdir_gotmp_rc=$?
export GOTMPDIR="$GOTMP"
echo "date=$(date '+%F %T %z') tree=$TREE roster=$ROSTER shape=$SHAPE"
echo "mkdir_gotmp_rc=$mkdir_gotmp_rc GOTMPDIR=$GOTMPDIR TMP=$TMP TEMP=$TEMP"
echo "blacklist_blob_on_disk=$(git -C "$REPO" hash-object "$TREE/internal/risk/blacklist.go")"
echo "teeth_blob_on_disk=$(git -C "$REPO" hash-object "$TREE/internal/risk/anchor_fallback_teeth_336_windows_test.go" 2>/dev/null || echo ABSENT)"
cd "$TREE" || { echo "cd_rc=NONZERO"; exit 1; }
go test ./internal/risk/ -run "$ROSTER" -count=1 -timeout 30m -v > "$LOGDIR/$LOGNAME" 2>&1
rc=$?
echo "go_test_rc=$rc"
echo "FOURCOUNT RUN=$(grep -c '^=== RUN' "$LOGDIR/$LOGNAME" || true) PASS_TOP=$(grep -c '^--- PASS' "$LOGDIR/$LOGNAME" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$LOGDIR/$LOGNAME" || true) SKIP=$(grep -c '^--- SKIP' "$LOGDIR/$LOGNAME" || true) rc=$rc"
grep -E '^ *(### )?FAIL|^ *--- FAIL' "$LOGDIR/$LOGNAME" | sed 's/([0-9.]*s)//' > "$LOGDIR/$LOGNAME.redlist" || true
echo "redlist=$(wc -l < "$LOGDIR/$LOGNAME.redlist")"
grep -E 'TEETH336' "$LOGDIR/$LOGNAME" | sed 's/^/  /' || true
tail -25 "$LOGDIR/$LOGNAME"
