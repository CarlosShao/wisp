#!/bin/bash
# 339-r2 GREEN-after：把工作树里那三个件覆进 HEAD 导出树，两平台跑同一枚 -run 名册。
# 只覆本腿写面（⛔ 从共享 index 反推，逐枚手写路径）。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
BASE="/d/tmp/wisp/20261105/339r2"
LOGD="$REPO/.scratch/wisp/probes/339/r2/logs"
RUN='TestAnchorMiss|TestAnchorFallbackTeeth'
FILES="internal/risk/anchor_miss_reason_334_test.go internal/risk/anchor_miss_reason_334_windows_test.go internal/risk/anchor_miss_reason_339_other_test.go"
mkdir -p "$LOGD"

for platform in linux win; do
  T="$BASE/$platform-after"
  rm -rf "$T"; mkdir -p "$T"
  git -C "$REPO" archive HEAD | tar -x -C "$T"
  for f in $FILES; do
    mkdir -p "$T/$(dirname "$f")"
    cp "$REPO/$f" "$T/$f"
    echo "overlay $platform $f sha=$(git -C "$REPO" hash-object "$T/$f")"
  done
done
echo "TREE_GO_FILES=$(find "$BASE/linux-after" -name '*.go' | wc -l)"

# --- Linux ---
MSYS_NO_PATHCONV=1 docker run --rm \
  -v "$BASE/linux-after":/src \
  -v "/d/work/base/gopath/pkg/mod":/go/pkg/mod:ro \
  -w /src \
  -e GOFLAGS=-mod=mod -e GOPROXY=off -e CGO_ENABLED=0 -e HOME=/src -e GOMODCACHE=/go/pkg/mod \
  golang:1.27 bash -c "mkdir -p /tmp/gt; export TMPDIR=/tmp/gt GOTMPDIR=/tmp/gt; go vet ./internal/risk/ > /tmp/gt/v 2>&1; echo GO_VET_RC=\$?; cat /tmp/gt/v; go test ./internal/risk/ -run '$RUN' -count=1 -timeout 30m -v > /tmp/gt/o 2>&1; echo GO_TEST_RC=\$?; cat /tmp/gt/o" \
  > "$LOGD/linux-green-miss-suite.log" 2>&1
O="$LOGD/linux-green-miss-suite.log"
echo "linux-after: $(grep -E '^GO_(VET|TEST)_RC=' "$O" | tr '\n' ' ') RUN=$(grep -c '^=== RUN' "$O" || true) PASS=$(grep -c '^--- PASS' "$O" || true) FAIL=$(grep -c '^--- FAIL' "$O" || true) SKIP=$(grep -c '^--- SKIP' "$O" || true)"
grep -E '^(--- (PASS|FAIL|SKIP)|ok|FAIL)' "$O" || true

# --- Windows ---
export TEMP='D:\tmp\wisp\20261105\339r2\gt' TMP="$TEMP" TMPDIR="$TEMP" GOTMPDIR="$TEMP"
export GOFLAGS=-mod=mod GOPROXY=off CGO_ENABLED=0 GOCACHE='D:\tmp\wisp\20261105\339r2\gocache'
( cd "$BASE/win-after" && go test ./internal/risk/ -run "$RUN" -count=1 -timeout 30m -v > "$LOGD/win-green-miss-suite.log" 2>&1; echo "WIN_RC=$?" )
O="$LOGD/win-green-miss-suite.log"
echo "win-after: $(grep -E '^WIN_RC=' "$O" | tr '\n' ' ') RUN=$(grep -c '^=== RUN' "$O" || true) PASS=$(grep -c '^--- PASS' "$O" || true) FAIL=$(grep -c '^--- FAIL' "$O" || true) SKIP=$(grep -c '^--- SKIP' "$O" || true)"
grep -E '^(--- (PASS|FAIL|SKIP)|ok|FAIL)' "$O" || true
