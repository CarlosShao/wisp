#!/bin/bash
# 307 编排者独立复跑：两台面（default / aliased）× 两状态（post=母仓 HEAD / pre=仓外 clone 858ee7f4）
# 尺口径：go test ./internal/risk/ -count=1 -v，四数用 grep -cE（含子测试行），每发自己落一行 rc=
REPO="/d/work/workspace/projects plans/Wisp"
OUT="$REPO/.scratch/wisp/probes/307/orch/logs"
PRE="/d/work/workspace/bisect/307/pre-858ee7f4"
ALIAS_SHORT='C:\Users\swq\tmp\RISKAT~1\RUNNER~1'

run_one() {
  local name="$1" dir="$2" table="$3"
  local file="$OUT/orch-$name-$table.txt"
  : > "$file"
  if [ "$table" = "aliased" ]; then
    export TMPDIR="$ALIAS_SHORT" TEMP="$ALIAS_SHORT" TMP="$ALIAS_SHORT"
  else
    unset TMPDIR
    export TEMP='C:\Users\swq\AppData\Local\Temp'
    export TMP='C:\Users\swq\AppData\Local\Temp'
  fi
  {
    echo "### $name · table=$table · dir=$dir"
    date '+%F %T %z'
    echo "HEAD_now=$(git -C "$dir" rev-parse HEAD)"
    echo "dirty_in_dir=$(git -C "$dir" status --porcelain | wc -l) 行"
    echo "env TMP=[$TEMP] TMPDIR=[$TMPDIR]"
    echo "----- 命令（逐字）-----"
    echo "go test -C \"$dir\" ./internal/risk/ -count=1 -v"
    echo "----- 原始输出 -----"
  } >> "$file"
  go test -C "$dir" ./internal/risk/ -count=1 -v >> "$file" 2>&1
  local rc=$?
  {
    echo "rc_test=$rc"
    echo "four_counts:"
    echo "  RUN=$(grep -cE '^[[:space:]]*=== RUN' "$file")"
    echo "  PASS=$(grep -cE '^[[:space:]]*--- PASS' "$file")"
    echo "  FAIL=$(grep -cE '^[[:space:]]*--- FAIL' "$file")"
    echo "  SKIP=$(grep -cE '^[[:space:]]*--- SKIP' "$file")"
    echo "fail_roster:"
    grep -E '^[[:space:]]*--- FAIL' "$file" | sed 's/^[[:space:]]*//'
  } >> "$file"
  echo "== $name/$table rc=$rc RUN=$(grep -cE '^[[:space:]]*=== RUN' "$file") PASS=$(grep -cE '^[[:space:]]*--- PASS' "$file") FAIL=$(grep -cE '^[[:space:]]*--- FAIL' "$file") SKIP=$(grep -cE '^[[:space:]]*--- SKIP' "$file")"
}

echo "start $(date '+%F %T %z')"
run_one post-head "$REPO" default
run_one post-head "$REPO" aliased
run_one pre-clone "$PRE" default
run_one pre-clone "$PRE" aliased
echo "end $(date '+%F %T %z')"
