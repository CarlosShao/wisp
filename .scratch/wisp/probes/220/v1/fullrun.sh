#!/usr/bin/env bash
# 220-v1 AC#5 整包两发（裁决腿自己跑·仓外 full 树·显式带足 -timeout）
# 原始 -v 输出留在仓外（体量可达数十 MB，⛔ 入库），仓内只落抽取件：四数＋红名名册＋红句逐字。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
PROBE="$REPO/.scratch/wisp/probes/220/v1"
LOGS="$PROBE/logs"
RAW=/tmp/wisp220v1/raw
mkdir -p "$LOGS" "$RAW"
TREE=/tmp/wisp220v1/full
# 台面补齐（本腿自加，具名写进读数件）：docs 与 design 两枚目录从 HEAD 取进导出树，
# 否则 internal/panel 里读 docs/PLAN.md 与 design/assets/tokens.css 的用例＝台面缺件假红。
( cd "$REPO" && git archive HEAD docs design ) | tar -x -C "$TREE" 2>/dev/null
echo "docs_files=$(find "$TREE/docs" -type f 2>/dev/null | wc -l) design_files=$(find "$TREE/design" -type f 2>/dev/null | wc -l)" > "$TREE/TREE-NOTE.txt"
got="$TREE/.gotmp"; mkdir -p "$got"; ng=$(cygpath -w "$got")

four() { # rawlog -> four numbers
  echo "RUN=$(grep -c '^=== RUN' "$1") PASS_TOP=$(grep -c '^--- PASS' "$1") FAIL=$(grep -c '^--- FAIL' "$1") SKIP=$(grep -c '^--- SKIP' "$1")"
}

run_pkg() { # tag pkgargs
  local tag="$1"; shift
  local raw="$RAW/$tag.log"
  {
    echo "=== $tag $(date '+%F %T %z') anchor=$(cd "$REPO" && git rev-parse HEAD) tree=$TREE"
    echo "cmd=go test $* -count=1 -timeout 30m"
    ( cd "$TREE" && PATH="$PWD/third_party/sherpa-onnx:$PATH" GOTMPDIR="$ng" TMP="$ng" TEMP="$ng" \
        go test "$@" -count=1 -timeout 30m ) > "$raw" 2>&1
    echo "rc=$?"
    echo "--- four ---"
    four "$raw"
    echo "--- 台面故障哨兵（任一命中⇒那一发不算终态）---"
    grep -nE 'panic: test timed out|0xc0000135|no tests to run|setup failed|no required module provides|cannot find the file specified|no native DLLs' "$raw" | head -10 || echo "(零命中)"
    echo "--- fail names ---"
    grep '^--- FAIL' "$raw" | sed 's/ (.*//' | sort | tee "$RAW/$tag.failnames.txt"
    echo "FAIL_TOTAL=$(grep -c '^--- FAIL' "$raw")  SUBFAIL=$(grep -c '^    --- FAIL' "$raw")"
    echo "raw_bytes=$(wc -c < "$raw") raw_lines=$(wc -l < "$raw")"
    echo "=== end $(date '+%F %T %z') ==="
  } | tee "$LOGS/40-full-$tag.txt"
  # 六枚已知红的耗时与红句逐字（本机状态两形，⛔ 计入判据）
  local SIX='TestAC4FocusReturnToPriorWindowGap33r5|TestAC14GoSideEvalPushReachesThePage|Test258BridgeMutationNoSrcKeepsOldBinding|Test258OccupiedCombinationNamesTheNewValue|Test258V1ProbeSummonEditRebindsLiveBall|TestAC246DevLegIgnoresTheTestTaskInjection'
  grep -A4 -E "^--- FAIL: ($SIX) " "$raw" > "$LOGS/41-$tag-sixred-verbatim.txt" 2>/dev/null
  echo "sixred_extract_bytes=$(wc -c < "$LOGS/41-$tag-sixred-verbatim.txt")" >> "$LOGS/40-full-$tag.txt"
}

run_pkg cmdwisp ./cmd/wisp/ -v
run_pkg internal ./internal/...
echo "ALLDONE $(date '+%F %T %z')"
