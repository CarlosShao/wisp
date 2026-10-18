#!/bin/bash
# 339-r1 交件侧 Linux 尺：指定树跑整包（配方同 linux.sh，树由参数给）。
# 用法：bash linuxtree.sh <tree-path> <label> [run-regex]
set -u
REPO="/d/work/workspace/projects plans/Wisp"
TREE="$1"
LABEL="$2"
RUNPAT="${3:-.}"
LOGD="$REPO/.scratch/wisp/probes/339/r1/logs"
mkdir -p "$LOGD"
OUT="$LOGD/$LABEL.log"
MSYS_NO_PATHCONV=1 docker run --rm \
  -v "$TREE":/src \
  -v "$(cd "$REPO" && go env GOMODCACHE)":/go/pkg/mod:ro \
  -w /src \
  -e GOFLAGS=-mod=mod -e GOPROXY=off -e CGO_ENABLED=0 -e HOME=/src -e GOMODCACHE=/go/pkg/mod \
  golang:1.27 bash -c "mkdir -p /tmp/gt; export TMPDIR=/tmp/gt GOTMPDIR=/tmp/gt; go test ./internal/risk/ -run '$RUNPAT' -count=1 -timeout 30m -v > /tmp/gt/out.txt 2>&1; echo GO_TEST_RC=\$?; cat /tmp/gt/out.txt" > "$OUT" 2>&1
echo "tree=$TREE label=$LABEL bytes=$(wc -c < "$OUT")"
echo "RUN=$(grep -c '^=== RUN' "$OUT" || true)"
echo "PASS_TOP=$(grep -c '^--- PASS' "$OUT" || true)"
echo "FAIL_TOP=$(grep -c '^--- FAIL' "$OUT" || true)"
echo "SKIP=$(grep -c '^--- SKIP' "$OUT" || true)"
grep '^GO_TEST_RC=' "$OUT" || true
echo "=== red roster ==="
grep -E '^--- FAIL' "$OUT" || true
echo "=== skip roster ==="
grep -E '^--- SKIP' "$OUT" || true
echo "=== fault words ==="
grep -nE 'walk|build failed|no test files|setup failed|creating work dir|0xc0000135' "$OUT" || echo "none"
