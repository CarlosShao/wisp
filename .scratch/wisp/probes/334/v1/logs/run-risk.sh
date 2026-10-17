#!/usr/bin/env bash
# 334-v1 颜色尺台件：只在【仓外导出台树】里跑 go test，⛔ 共享工作树。
# 用法：run-risk.sh <树根> <名册正则> <日志名> [--ci-shape|--dshape]
#   --ci-shape  ＝ GOTMPDIR 逐字照 scripts/portable-tests.sh:152-159 的组法
#                 （$LOCALAPPDATA\\wisp-test-gotmp\\portable-<日期>-<pid>，TMP/TEMP 一字不动）
#   --dshape    ＝ GOTMPDIR 指 D: 上本腿自建临时根（先 mkdir -p，没建会红成 creating work dir）
# 四数尺：RUN=^=== RUN（含子测试）／PASS_TOP=^--- PASS（只顶层）／FAIL_TOP=^--- FAIL／SKIP=^--- SKIP，
# 四数不许混用、不许相减；一律 -count=1 -timeout 30m。
set -u
REPO="D:/work/workspace/projects plans/Wisp"
TREE="${1:?树根}"
ROSTER="${2:?名册正则}"
LOGNAME="${3:?日志名}"
SHAPE="${4:---ci-shape}"
LOGDIR="$REPO/.scratch/wisp/probes/334/v1/logs"
export GOPATH='D:\work\base\gopath'
export PATH="$TREE/third_party/sherpa-onnx:$TREE/build:$PATH"
STAMP="$(date +%Y%m%d-%H%M%S)-$$"
case "$SHAPE" in
  --ci-shape)
    base="${LOCALAPPDATA:-C:\\Users\\swq\\AppData\\Local}\\wisp-test-gotmp"
    GOTMP="$base\\portable-334v1-$STAMP" ;;
  --dshape)
    GOTMP="D:\\W334V1Gotmp\\v1-$STAMP" ;;
  *) echo "unknown shape $SHAPE"; exit 2 ;;
esac
mkdir -p "$GOTMP" && mkdir_gotmp_rc=$? || mkdir_gotmp_rc=$?
export GOTMPDIR="$GOTMP"
echo "date=$(date '+%F %T %z') tree=$TREE roster=$ROSTER shape=$SHAPE"
echo "mkdir_gotmp_rc=$mkdir_gotmp_rc GOTMPDIR=$GOTMPDIR TMP=$TMP TEMP=$TEMP"
echo "blacklist_blob_on_disk=$(git -C "$REPO" hash-object "$TREE/internal/risk/blacklist.go")"
cd "$TREE" || { echo "cd_rc=NONZERO"; exit 1; }
go test ./internal/risk/ -run "$ROSTER" -count=1 -timeout 30m -v > "$LOGDIR/$LOGNAME" 2>&1
rc=$?
echo "go_test_rc=$rc"
RUN=$(grep -c '^=== RUN' "$LOGDIR/$LOGNAME" || true)
PASS_TOP=$(grep -c '^--- PASS' "$LOGDIR/$LOGNAME" || true)
FAIL_TOP=$(grep -c '^--- FAIL' "$LOGDIR/$LOGNAME" || true)
SKIP=$(grep -c '^--- SKIP' "$LOGDIR/$LOGNAME" || true)
echo "FOURCOUNT RUN=$RUN PASS_TOP=$PASS_TOP FAIL_TOP=$FAIL_TOP SKIP=$SKIP rc=$rc"
grep -E '^ *(### )?FAIL|^ *--- FAIL' "$LOGDIR/$LOGNAME" | sed 's/([0-9.]*s)//' > "$LOGDIR/$LOGNAME.redlist" || true
echo "redlist=$(wc -l < "$LOGDIR/$LOGNAME.redlist")"
tail -40 "$LOGDIR/$LOGNAME"
