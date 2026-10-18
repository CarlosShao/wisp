#!/bin/bash
# 339-r1 AC#0 台件②：Linux 容器跑 ./internal/risk/（⛔ 复用编排者那棵树）。
# 用法：bash linux.sh <commit> <label> [extra go test args...]
# 配方：八件套+docs 导出树挂 /src；宿主 module 缓存只读挂 /go/pkg/mod；
#       env=GOFLAGS=-mod=mod GOPROXY=off CGO_ENABLED=0 HOME=/src GOMODCACHE=/go/pkg/mod；
#       TMPDIR/GOTMPDIR 容器内先 mkdir -p；MSYS_NO_PATHCONV=1 防 /src 被改写。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
ROOT="/d/tmp/wisp/20261018/339r1"
ANCHOR="$1"
LABEL="$2"
RUNPAT="${3:-.}"
LOGDIR="$REPO/.scratch/wisp/probes/339/r1/logs"
TREE="$ROOT/linux/tree"
cd "$REPO" || exit 9
mkdir -p "$LOGDIR"
CACHE="$(go env GOMODCACHE)"
OUT="$LOGDIR/$LABEL.log"
echo "anchor=$ANCHOR label=$LABEL run=$RUNPAT modcache=$CACHE"
MSYS_NO_PATHCONV=1 docker run --rm \
  -v "$TREE":/src \
  -v "$CACHE":/go/pkg/mod:ro \
  -w /src \
  -e GOFLAGS=-mod=mod -e GOPROXY=off -e CGO_ENABLED=0 -e HOME=/src \
  -e GOMODCACHE=/go/pkg/mod \
  golang:1.27 bash -c "mkdir -p /tmp/gt; export TMPDIR=/tmp/gt GOTMPDIR=/tmp/gt; go test ./internal/risk/ -run '$RUNPAT' -count=1 -timeout 30m -v > /tmp/gt/out.txt 2>&1; echo GO_TEST_RC=\$?; cat /tmp/gt/out.txt" > "$OUT" 2>&1
echo "wrote=$OUT bytes=$(wc -c < "$OUT")"
echo "=== four numbers (RUN / PASS_TOP / FAIL_TOP / SKIP_TOP) ==="
echo "RUN=$(grep -c '^=== RUN' "$OUT" || true)"
echo "PASS_TOP=$(grep -c '^--- PASS' "$OUT" || true)"
echo "FAIL_TOP=$(grep -c '^--- FAIL' "$OUT" || true)"
echo "SKIP=$(grep -c '^--- SKIP' "$OUT" || true)"
echo "=== red roster (top-level names) ==="
grep -E '^--- FAIL' "$OUT" || true
echo "=== GO_TEST_RC ==="
grep -E '^GO_TEST_RC=' "$OUT" || true
echo "===台面故障名册（walk / build failed / no test files / setup failed / creating work dir）==="
grep -nE 'walk|build failed|no test files|setup failed|creating work dir|0xc0000135' "$OUT" || echo "none"
