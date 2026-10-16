#!/usr/bin/env bash
# 成对台件·改前那一发：在 fd01e796（220-r2 接线之前）的仓外导出树上跑那六枚红
# 编排者本人跑；⛔ 动任何断言；⛔ 在工作树里跑"改前那发"
set -u
M=/tmp/wisp-orch-pre220b-6C3v9y
REPO="/d/work/workspace/projects plans/Wisp"
LOG="$REPO/.scratch/wisp/probes/orch/20261016/logs/orch-sixred-pre220-dl.txt"
GOT="$M/.gotmp2"
mkdir -p "$GOT"
NMGOT=$(cygpath -w "$GOT")
NMBLD=$(cygpath -w "$M/build")
NMPART=$(cygpath -w "$M/third_party/sherpa-onnx")
RUN='TestAC4FocusReturnToPriorWindowGap33r5|TestAC14GoSideEvalPushReachesThePage|Test258BridgeMutationNoSrcKeepsOldBinding|Test258OccupiedCombinationNamesTheNewValue|Test258V1ProbeSummonEditRebindsLiveBall|TestAC246DevLegIgnoresTheTestTaskInjection'

{
  echo "=== provenance ==="
  date '+%Y-%m-%d %H:%M:%S %z'
  echo "tree=$M"
  echo "anchor=fd01e7962 (== 2d10bd0f^，接线之前)"
  echo "go version: $(go version)"
  echo "dll next to pkg:"; ls -la "$M/cmd/wisp/"*.dll
  echo "=== run ==="
} > "$LOG" 2>&1

cd "$M" || exit 9
GOTMPDIR="$NMGOT" TMP="$NMGOT" TEMP="$NMGOT" PATH="$PATH:$NMBLD:$NMPART" \
  go test ./cmd/wisp -count=1 -timeout 20m -run "$RUN" >> "$LOG" 2>&1
RC=$?
{
  echo "=== tail ==="
  echo "rc=$RC"
  date '+%Y-%m-%d %H:%M:%S %z'
  grep -c '^--- FAIL' "$LOG" || true
  grep '^--- FAIL' "$LOG" || true
  grep -c '^--- PASS' "$LOG" || true
} >> "$LOG" 2>&1
echo "DONE rc=$RC"
