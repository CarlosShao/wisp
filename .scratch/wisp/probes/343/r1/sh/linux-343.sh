#!/bin/bash
# 343-r1 AC#4 Linux 侧台件（配方逐字抄自 339-r1 的 linux.sh，仅换树/日志落点＝⛔ 复用别腿的树）。
# 用法：bash linux-343.sh <tree-path> <label>
set -u
REPO="/d/work/workspace/projects plans/Wisp"
TREE="$1"
LABEL="$2"
LOGDIR="$REPO/.scratch/wisp/probes/343/r1/logs"
mkdir -p "$LOGDIR"
OUT="$LOGDIR/$LABEL.log"
CACHE="$(cd "$REPO" && go env GOMODCACHE)"
echo "tree=$TREE label=$LABEL modcache=$CACHE"
MSYS_NO_PATHCONV=1 docker run --rm \
  -v "$TREE":/src \
  -v "$CACHE":/go/pkg/mod:ro \
  -w /src \
  -e GOFLAGS=-mod=mod -e GOPROXY=off -e CGO_ENABLED=0 -e HOME=/src -e GOMODCACHE=/go/pkg/mod \
  golang:1.27 bash -c "mkdir -p /tmp/gt; export TMPDIR=/tmp/gt GOTMPDIR=/tmp/gt; go test ./internal/risk/ -count=1 -timeout 30m -v > /tmp/gt/out.txt 2>&1; echo GO_TEST_RC=\$?; cat /tmp/gt/out.txt" > "$OUT" 2>&1
echo "wrote=$OUT bytes=$(wc -c < "$OUT")"
echo "=== four numbers (RUN / PASS_TOP / FAIL_TOP / SKIP) ==="
echo "RUN=$(grep -c '^=== RUN' "$OUT" || true)"
echo "PASS_TOP=$(grep -c '^--- PASS' "$OUT" || true)"
echo "FAIL_TOP=$(grep -c '^--- FAIL' "$OUT" || true)"
echo "SKIP=$(grep -c '^--- SKIP' "$OUT" || true)"
echo "=== red roster (top-level names) ==="
grep -E '^--- FAIL' "$OUT" || true
echo "=== skip roster ==="
grep -E '^--- SKIP' "$OUT" || true
echo "=== GO_TEST_RC ==="
grep -E '^GO_TEST_RC=' "$OUT" || true
echo "=== 台面故障名册 ==="
grep -nE 'walk|build failed|no test files|setup failed|creating work dir|0xc0000135' "$OUT" || echo "none"
