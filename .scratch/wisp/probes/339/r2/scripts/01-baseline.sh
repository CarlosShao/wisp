#!/bin/bash
# 339-r2 AC#0 起手现跑：两平台成对基线（锚＝HEAD，具名提交树，仓外导出）。
# 复用判据：⛔ 抄 339-r1/347 的读数，逐枚重打。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
BASE="/d/tmp/wisp/20261105/339r2"
LOGD="$REPO/.scratch/wisp/probes/339/r2/logs"
RUN='TestAnchorMissReasonNamesCandidateAndMissingAnchor|TestAnchorMissKeepsNonOverridableClassAndNeedL2|TestAnchorMiss|TestAnchorFallbackTeeth'
mkdir -p "$LOGD"

HEAD=$(git -C "$REPO" rev-parse HEAD)
echo "ANCHOR_HEAD=$HEAD"

# --- 仓外导出树：before（= HEAD 原样） ---
rm -rf "$BASE/linux-before"; mkdir -p "$BASE/linux-before"
git -C "$REPO" archive HEAD | tar -x -C "$BASE/linux-before"
echo "TREE_BEFORE_GO_FILES=$(find "$BASE/linux-before" -name '*.go' | wc -l)"

# --- Linux 档基线 ---
MSYS_NO_PATHCONV=1 docker run --rm \
  -v "$BASE/linux-before":/src \
  -v "/d/work/base/gopath/pkg/mod":/go/pkg/mod:ro \
  -w /src \
  -e GOFLAGS=-mod=mod -e GOPROXY=off -e CGO_ENABLED=0 -e HOME=/src -e GOMODCACHE=/go/pkg/mod \
  golang:1.27 bash -c "mkdir -p /tmp/gt; export TMPDIR=/tmp/gt GOTMPDIR=/tmp/gt; go test ./internal/risk/ -run '$RUN' -count=1 -timeout 30m -v > /tmp/gt/o 2>&1; echo GO_TEST_RC=\$?; cat /tmp/gt/o" \
  > "$LOGD/linux-base-miss-suite.log" 2>&1

O="$LOGD/linux-base-miss-suite.log"
echo "linux-base: $(grep -E '^GO_TEST_RC=' "$O") RUN=$(grep -c '^=== RUN' "$O" || true) PASS=$(grep -c '^--- PASS' "$O" || true) FAIL=$(grep -c '^--- FAIL' "$O" || true) SKIP=$(grep -c '^--- SKIP' "$O" || true)"
grep -E '^\s*--- (PASS|FAIL|SKIP)' "$O" || true

# --- Windows 档基线（本机原生，仓外树同锚） ---
rm -rf "$BASE/win-before"; mkdir -p "$BASE/win-before"
git -C "$REPO" archive HEAD | tar -x -C "$BASE/win-before"
export TEMP='D:\tmp\wisp\20261105\339r2\gt' TMP="$TEMP" TMPDIR="$TEMP" GOTMPDIR="$TEMP"
export GOFLAGS=-mod=mod GOPROXY=off CGO_ENABLED=0 GOCACHE='D:\tmp\wisp\20261105\339r2\gocache'
( cd "$BASE/win-before" && go test ./internal/risk/ -run "$RUN" -count=1 -timeout 30m -v > "$LOGD/win-base-miss-suite.log" 2>&1; echo "WIN_RC=$?" )
O="$LOGD/win-base-miss-suite.log"
echo "win-base: $(grep -E '^WIN_RC=' "$O" || true) $(tail -1 "$O") RUN=$(grep -c '^=== RUN' "$O" || true) PASS=$(grep -c '^--- PASS' "$O" || true) FAIL=$(grep -c '^--- FAIL' "$O" || true) SKIP=$(grep -c '^--- SKIP' "$O" || true)"
grep -E '^\s*--- (PASS|FAIL|SKIP)' "$O" || true
