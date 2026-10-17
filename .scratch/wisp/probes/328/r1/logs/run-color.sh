#!/usr/bin/env bash
# 328-r1 颜色尺：在仓外导出台树上跑一枚包，四数逐枚点名（^=== RUN 含子测试 ⇄ ^--- PASS 只顶层）。
# 用法：run-color.sh <树根> <包> <-run 样式|ALL> <日志名>
set -u
ROOT="${1:?}"
PKG="${2:?}"
PATTERN="${3:?}"
LOGNAME="${4:?}"
LOGDIR="D:/work/workspace/projects plans/Wisp/.scratch/wisp/probes/328/r1/logs"
export GOPATH='D:\work\base\gopath'
export GOTMPDIR='D:\tmp\328r1-gotmp'
mkdir -p "$GOTMPDIR"
cd "$ROOT" || { echo "cd_rc=NONZERO"; exit 1; }
LOG="$LOGDIR/$LOGNAME.log"
{ echo "date=$(date '+%F %T %z') root=$ROOT pkg=$PKG pattern=$PATTERN"; echo "gofiles=$(find internal cmd tools -name '*.go' | wc -l)"; } > "$LOG"
if [ "$PATTERN" = "ALL" ]; then
  go test -v "./$PKG" -count=1 -timeout 30m >> "$LOG" 2>&1
else
  go test -v "./$PKG" -run "$PATTERN" -count=1 -timeout 30m >> "$LOG" 2>&1
fi
echo "go_test_rc=$?" >> "$LOG"
# 四数与红名册先算进变量（⛔ 让 grep 一边读同一枚文件一边被追加），再一次性追加落盘。
RUN=$(grep -c '^=== RUN' "$LOG" || true)
PASS_TOP=$(grep -c '^--- PASS' "$LOG" || true)
FAIL_TOP=$(grep -c '^--- FAIL' "$LOG" || true)
SKIP=$(grep -c '^--- SKIP' "$LOG" || true)
RED=$(grep '^--- FAIL' "$LOG" || echo "NO_TOP_FAIL")
VERD=$(grep -E '^(ok|FAIL|# )' "$LOG" | head -40 || true)
{
  echo "RUN=$RUN"
  echo "PASS_TOP=$PASS_TOP"
  echo "FAIL_TOP=$FAIL_TOP"
  echo "SKIP=$SKIP"
  echo "--- red roster ---"
  echo "$RED"
  echo "--- package verdict lines ---"
  echo "$VERD"
} >> "$LOG"
tail -30 "$LOG"
