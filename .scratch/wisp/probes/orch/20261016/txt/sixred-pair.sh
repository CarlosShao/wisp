#!/usr/bin/env bash
# 成对台件（同台面·两发）：六枚 cmd/wisp 红，改前 fd01e7962 ⇄ 改后 HEAD
# ★台面补齐定式（我自己踩过的三枚坑，逐条带上）：
#   ① 建树必带 frontend（否则 assets.go:25 报 no required module ⇒ setup failed 假红）
#   ② build/ third_party/ 两枚目录未入库 ⇒ git archive 取不到 ⇒ 必须手工拷 DLL，
#      否则 0xc0000135（加载期缺 DLL）与 "no native DLLs in ..\..\third_party\sherpa-onnx"（测试自述的缺件红）
#   ③ frontend/dist 只有 1 枚跟踪件 ⇒ 真窗用例在导出树里页面是空的 ⇒ title="" 假红
#      ⇒ 必须把在树里的 frontend/dist 整枚拷进去（两发拷同一份 ⇒ 台面同）
set -u
REPO="/d/work/workspace/projects plans/Wisp"
LOGDIR="$REPO/.scratch/wisp/probes/orch/20261016/logs"
RUN='TestAC4FocusReturnToPriorWindowGap33r5|TestAC14GoSideEvalPushReachesThePage|Test258BridgeMutationNoSrcKeepsOldBinding|Test258OccupiedCombinationNamesTheNewValue|Test258V1ProbeSummonEditRebindsLiveBall|TestAC246DevLegIgnoresTheTestTaskInjection'
PRE=fd01e7962

build_tree() {
  local ref="$1" out="$2"
  rm -rf "$out" 2>/dev/null
  mkdir -p "$out"
  ( cd "$REPO" && git archive "$ref" go.mod go.sum internal cmd tools scripts frontend ) | tar -x -C "$out"
  cp -r "$REPO/frontend/dist/." "$out/frontend/dist/" 2>/dev/null
  mkdir -p "$out/build" "$out/third_party/sherpa-onnx"
  cp -p "$REPO/build/"*.dll "$out/build/" 2>/dev/null
  cp -p "$REPO/third_party/sherpa-onnx/"*.dll "$out/third_party/sherpa-onnx/" 2>/dev/null
  cp -p "$REPO/build/"*.dll "$out/cmd/wisp/" 2>/dev/null
  local n; n=$(find "$out" -name '*.go' | wc -l)
  echo "tree=$out ref=$ref go_files=$n dist_files=$(find "$out/frontend/dist" -type f | wc -l) dll=$(ls "$out/third_party/sherpa-onnx"/*.dll 2>/dev/null | wc -l)"
}

run_six() {
  local tag="$1" tree="$2" log="$3"
  local got="$tree/.gotmp"; mkdir -p "$got"
  local ngot; ngot=$(cygpath -w "$got")
  {
    echo "=== $tag ==="
    date '+%Y-%m-%d %H:%M:%S %z'
    ( cd "$tree" && git rev-parse HEAD 2>/dev/null || true )
    echo "anchor=$4"
  } > "$log" 2>&1
  ( cd "$tree" && GOTMPDIR="$ngot" TMP="$ngot" TEMP="$ngot" go test ./cmd/wisp -count=1 -timeout 20m -run "$RUN" >> "$log" 2>&1; echo "rc=$?" >> "$log" )
  {
    echo "--- names ---"
    grep '^--- FAIL' "$log" || echo "(no FAIL lines)"
    echo "FAIL_COUNT=$(grep -c '^--- FAIL' "$log")"
    echo "PASS_COUNT=$(grep -c '^--- PASS' "$log")"
    grep -E '^(ok|FAIL|panic|.*setup failed|.*cannot find|.*0xc000)' "$log" | head -8
    date '+%Y-%m-%d %H:%M:%S %z'
  } >> "$log" 2>&1
}

T1=/tmp/wisp-pair-pre; T2=/tmp/wisp-pair-post
build_tree "$PRE" "$T1" >> "$LOGDIR/orch-sixred-pair-114x.log" 2>&1
build_tree HEAD "$T2" >> "$LOGDIR/orch-sixred-pair-114x.log" 2>&1
cat "$LOGDIR/orch-sixred-pair-114x.log"
run_six "PRE-WIRING (fd01e7962)" "$T1" "$LOGDIR/orch-sixred-pair-pre.txt" "$PRE"
run_six "POST (HEAD)" "$T2" "$LOGDIR/orch-sixred-pair-post.txt" "$(cd "$REPO" && git rev-parse --short HEAD)"
echo ALLDONE
