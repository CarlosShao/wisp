#!/bin/sh
# A846 读数 第三发：只重打 §7 那四把尺（b 发里我在双引号串里写了反引号包 `.go`＝那段真被当命令跑了，
# 两行标签被吞、尺的【值】⛔ 受影响＝5/5/334/22/5 在 b 发已照原样量出。这枚坑是我名下的第三次，具名记我。）
cd "$(git rev-parse --show-toplevel)" || exit 9
D=.scratch/wisp/probes/orch/20261011
LOG=$D/logs/orch-A846-readings-c.txt
DOTGO='[.]go$'
SCR='^\.scratch/'

NARROW=$(git ls-files --eol cmd internal tools | grep 'w/crlf')
WIDEALL=$(git ls-files --eol | grep 'w/crlf')
NON=$(printf '%s\n' "$WIDEALL" | awk '{print $NF}' | grep -v "$SCR")

{ echo "date=$(date '+%Y-%m-%d %H:%M:%S %z')  HEAD=$(git rev-parse --short HEAD)"
  echo "### 7-c) CRLF 假枚名册四把尺，射程＋单位名逐把写明（这条⛔ 带反引号）"
  echo "  尺1  射程=cmd internal tools      单位=件(全扩展名) = $(printf '%s\n' "$NARROW" | grep -c 'w/crlf')"
  echo "  尺2  射程=cmd internal tools      单位=.go 件       = $(printf '%s\n' "$NARROW" | awk '{print $NF}' | grep -c "$DOTGO")"
  echo "  尺3  射程=全仓 tracked            单位=件(⛔剥 .scratch) = $(printf '%s\n' "$WIDEALL" | grep -c 'w/crlf')"
  echo "  尺4  射程=全仓 tracked 剥 .scratch 单位=件           = $(printf '%s\n' "$NON" | grep -c '')"
  echo "  尺4b 射程=全仓 tracked 剥 .scratch 单位=.go 件        = $(printf '%s\n' "$NON" | grep -c "$DOTGO")"
  echo "  尺4b 件名 = $(printf '%s\n' "$NON" | grep "$DOTGO" | tr '\n' ' ')"
  echo "  尺4 全 22 枚件名 = $(printf '%s\n' "$NON" | tr '\n' ' ')"
} > $LOG 2>&1
cat $LOG
