#!/bin/bash
# 339-r2 收尾尺：从**已提交树**（git archive HEAD）重建两平台台件，跑同一 -run 名册。
# 目的＝证明读数属于提交字节，不是工作树里的未入库形状。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
BASE="/d/tmp/wisp/20261105/339r2"
LOGD="$REPO/.scratch/wisp/probes/339/r2/logs"
RUN='TestAnchorMiss|TestAnchorFallbackTeeth'
HEAD=$(git -C "$REPO" rev-parse HEAD)
echo "ANCHOR_HEAD=$HEAD"

for p in linux win; do rm -rf "$BASE/final-$p"; mkdir -p "$BASE/final-$p"; git -C "$REPO" archive HEAD | tar -x -C "$BASE/final-$p"; done
echo "go_files=$(find "$BASE/final-linux" -name '*.go' | wc -l)"

MSYS_NO_PATHCONV=1 docker run --rm \
  -v "$BASE/final-linux":/src -v "/d/work/base/gopath/pkg/mod":/go/pkg/mod:ro -w /src \
  -e GOFLAGS=-mod=mod -e GOPROXY=off -e CGO_ENABLED=0 -e HOME=/src -e GOMODCACHE=/go/pkg/mod \
  golang:1.27 bash -c "mkdir -p /tmp/gt; export TMPDIR=/tmp/gt GOTMPDIR=/tmp/gt; go test ./internal/risk/ -run '$RUN' -count=1 -timeout 30m -v > /tmp/gt/o 2>&1; echo GO_TEST_RC=\$?; cat /tmp/gt/o" \
  > "$LOGD/linux-final-miss-suite.log" 2>&1
O="$LOGD/linux-final-miss-suite.log"
echo "linux-final: $(grep -E '^GO_TEST_RC=' "$O") RUN=$(grep -c '^=== RUN' "$O" || true) PASS=$(grep -c '^--- PASS' "$O" || true) FAIL=$(grep -c '^--- FAIL' "$O" || true) SKIP=$(grep -c '^--- SKIP' "$O" || true)"
grep -E '^(--- (PASS|FAIL)|ok|FAIL)' "$O" || true

export TEMP='D:\tmp\wisp\20261105\339r2\gt' TMP="$TEMP" TMPDIR="$TEMP" GOTMPDIR="$TEMP"
export GOFLAGS=-mod=mod GOPROXY=off CGO_ENABLED=0 GOCACHE='D:\tmp\wisp\20261105\339r2\gocache'
( cd "$BASE/final-win" && go test ./internal/risk/ -run "$RUN" -count=1 -timeout 30m -v > "$LOGD/win-final-miss-suite.log" 2>&1; echo "WIN_RC=$?" )
O="$LOGD/win-final-miss-suite.log"
echo "win-final: $(grep -E '^WIN_RC=' "$O") RUN=$(grep -c '^=== RUN' "$O" || true) PASS=$(grep -c '^--- PASS' "$O" || true) FAIL=$(grep -c '^--- FAIL' "$O" || true) SKIP=$(grep -c '^--- SKIP' "$O" || true)"
grep -E '^(--- (PASS|FAIL)|ok|FAIL)' "$O" || true
