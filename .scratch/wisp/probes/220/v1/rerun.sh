#!/usr/bin/env bash
# 220-v1 AC#5 复跑（同一台面·同一棵树·⛔ 换锚⛔ 换机·派单要求"在同一台面上复跑一次"）
set -u
REPO="/d/work/workspace/projects plans/Wisp"
LOGS="$REPO/.scratch/wisp/probes/220/v1/logs"
TREE=/tmp/wisp220v1/full
RAW=/tmp/wisp220v1/raw
got="$TREE/.gotmp"; mkdir -p "$got"; ng=$(cygpath -w "$got")
tag="$1"; shift
raw="$RAW/$tag.log"
{
  echo "=== $tag $(date '+%F %T %z') tree=$TREE 台件＝同一枚（docs/design 已补，⛔ 重建）"
  echo "cmd=go test $* -count=1 -timeout 30m"
  ( cd "$TREE" && PATH="$PWD/third_party/sherpa-onnx:$PATH" GOTMPDIR="$ng" TMP="$ng" TEMP="$ng" \
      go test "$@" -count=1 -timeout 30m ) > "$raw" 2>&1
  echo "rc=$?"
  echo "RUN=$(grep -c '^=== RUN' "$raw") PASS_TOP=$(grep -c '^--- PASS' "$raw") FAIL=$(grep -c '^--- FAIL' "$raw") SKIP=$(grep -c '^--- SKIP' "$raw")"
  echo "包级行：$(grep -E '^(ok|FAIL)\s+github' "$raw" | tr '\n' ' ')"
  echo "--- 末三行（断点在哪一枚）---"
  grep -E '^(=== RUN|--- (PASS|FAIL)|    [a-z_0-9]+\.go:|ok|FAIL)' "$raw" | tail -4
  echo "--- 台面/中断哨兵 ---"
  grep -nE 'panic: test timed out|0xc0000135|no tests to run|setup failed|no required module provides|Exception|fatal error' "$raw" | head -5 || echo "(零命中)"
  echo "--- 六枚已知名在我这发的颜色 ---"
  for n in TestAC4FocusReturnToPriorWindowGap33r5 TestAC14GoSideEvalPushReachesThePage Test258BridgeMutationNoSrcKeepsOldBinding Test258OccupiedCombinationNamesTheNewValue Test258V1ProbeSummonEditRebindsLiveBall TestAC246DevLegIgnoresTheTestTaskInjection; do
    line=$(grep -E "^--- (PASS|FAIL): $n " "$raw" | head -1); echo "$n : ${line:-未跑到}"
  done
  echo "=== end $(date '+%F %T %z') ==="
} | tee "$LOGS/42-rerun-$tag.txt"
