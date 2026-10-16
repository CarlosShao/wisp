#!/usr/bin/env bash
# 220-v1 第三发（排除那枚把进程打断的用例，⛔ 动任何断言、⛔ 当票 220 的终态判据，只取名册）
# ＋ ./internal/... 补 models/ 后的第二发（消掉三枚台面缺件假红）
set -u
REPO="/d/work/workspace/projects plans/Wisp"
LOGS="$REPO/.scratch/wisp/probes/220/v1/logs"
TREE=/tmp/wisp220v1/full
RAW=/tmp/wisp220v1/raw
( cd "$REPO" && git archive HEAD models ) | tar -x -C "$TREE" 2>/dev/null
echo "models_files=$(find "$TREE/models" -type f 2>/dev/null | wc -l)" >> "$TREE/TREE-NOTE.txt"
got="$TREE/.gotmp"; mkdir -p "$got"; ng=$(cygpath -w "$got")

raw="$RAW/cmdwisp3.log"
{
  echo "=== cmdwisp3 $(date '+%F %T %z') tree=$TREE"
  echo "cmd=go test ./cmd/wisp/ -v -count=1 -timeout 30m -skip '^TestAC13BringUpSurvivesAReusedThreadQuit$'"
  echo "排除理由＝前两发（40-full-cmdwisp.txt／42-rerun-cmdwisp2.txt）在同一枚用例处进程硬终止、零 panic 句、"
  echo "  RUN=203 PASS_TOP=145 FAIL=3 两发同数 ⇒ 那一枚（票 33 的 poisoned-thread 用例）在本机拦腰斩断整包；"
  echo "  ⛔ 本发算票 220 AC#5 的终态判据，只用来取名册"
  ( cd "$TREE" && PATH="$PWD/third_party/sherpa-onnx:$PATH" GOTMPDIR="$ng" TMP="$ng" TEMP="$ng" \
      go test ./cmd/wisp/ -v -count=1 -timeout 30m -skip '^TestAC13BringUpSurvivesAReusedThreadQuit$' ) > "$raw" 2>&1
  echo "rc=$?"
  echo "RUN=$(grep -c '^=== RUN' "$raw") PASS_TOP=$(grep -c '^--- PASS' "$raw") FAIL=$(grep -c '^--- FAIL' "$raw") SKIP=$(grep -c '^--- SKIP' "$raw")"
  echo "raw_bytes=$(wc -c < "$raw") raw_lines=$(wc -l < "$raw")"
  echo "--- 台面/中断哨兵 ---"
  grep -nE 'panic: test timed out|0xc0000135|no tests to run|setup failed|no required module provides|fatal error' "$raw" | head -5
  echo "--- 末四行 ---"; grep -E '^(=== RUN|--- (PASS|FAIL)|ok|FAIL)' "$raw" | tail -4
  echo "--- 红名名册（顶层）---"; grep '^--- FAIL' "$raw" | sed 's/ (.*//'
  echo "--- 六枚已知名颜色 ---"
  for n in TestAC4FocusReturnToPriorWindowGap33r5 TestAC14GoSideEvalPushReachesThePage Test258BridgeMutationNoSrcKeepsOldBinding Test258OccupiedCombinationNamesTheNewValue Test258V1ProbeSummonEditRebindsLiveBall TestAC246DevLegIgnoresTheTestTaskInjection TestRunPacketLightsTheRosterRowOfItsLiveL1Window; do
    line=$(grep -E "^--- (PASS|FAIL): $n " "$raw" | head -1); echo "$n : ${line:-未跑到}"
  done
  echo "=== end $(date '+%F %T %z') ==="
} | tee "$LOGS/43-cmdwisp-skipAC13.txt"
grep -A5 -E '^--- FAIL: ' "$raw" | head -160 > "$LOGS/43-cmdwisp-failblocks.txt"
wc -c "$LOGS/43-cmdwisp-failblocks.txt" >> "$LOGS/43-cmdwisp-skipAC13.txt"

raw2="$RAW/internal2.log"
{
  echo "=== internal2 $(date '+%F %T %z') tree=$TREE（已补 models/）"
  echo "cmd=go test ./internal/... -count=1 -timeout 30m"
  ( cd "$TREE" && PATH="$PWD/third_party/sherpa-onnx:$PATH" GOTMPDIR="$ng" TMP="$ng" TEMP="$ng" \
      go test ./internal/... -count=1 -timeout 30m ) > "$raw2" 2>&1
  echo "rc=$?"
  echo "RUN=$(grep -c '^=== RUN' "$raw2") PASS_TOP=$(grep -c '^--- PASS' "$raw2") FAIL=$(grep -c '^--- FAIL' "$raw2") SKIP=$(grep -c '^--- SKIP' "$raw2")"
  echo "包级 FAIL=$(grep -cE '^FAIL\s+github' "$raw2") 包级 ok=$(grep -cE '^ok\s+github' "$raw2") 无测试文件=$(grep -c 'no test files' "$raw2")"
  echo "--- 红名名册 ---"; grep '^--- FAIL' "$raw2" | sed 's/ (.*//'
  echo "--- 哨兵 ---"; grep -nE 'panic|setup failed|no required module|0xc0000135' "$raw2" | head -5
  echo "=== end $(date '+%F %T %z') ==="
} | tee "$LOGS/44-internal-with-models.txt"
