#!/bin/bash
# 339-v1 台件②：Linux 容器跑 ./internal/risk/（裁决腿自己建的树，⛔ 复用腿／编排者那几棵）。
# 用法：bash linux.sh <tree-dir> <label> [run-pattern]
# 配方＝公共块：八件套+docs 导出树挂 /src；宿主 module 缓存只读挂 /go/pkg/mod；
#   GOFLAGS=-mod=mod GOPROXY=off CGO_ENABLED=0 HOME=/src GOMODCACHE=/go/pkg/mod；
#   容器内先 mkdir -p /tmp/gt 再指 TMPDIR/GOTMPDIR；挂载路径写死 D:/ 正斜杠形＋
#   MSYS_NO_PATHCONV=1（防 Git Bash 改写＝台账 A1027 那枚 rc=125 故障件同形）。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
TREE_WIN="$1"          # 例 D:/tmp/wisp/20261018/339v1/post/tree
LABEL="$2"
RUNPAT="${3:-.}"
LOGDIR="$REPO/.scratch/wisp/probes/339/v1/logs"
CACHE_WIN='D:/work/base/gopath/pkg/mod'
mkdir -p "$LOGDIR"
OUT="$LOGDIR/$LABEL.log"
echo "tree=$TREE_WIN label=$LABEL run=$RUNPAT"
MSYS_NO_PATHCONV=1 docker run --rm \
  -v "$TREE_WIN":/src \
  -v "$CACHE_WIN":/go/pkg/mod:ro \
  -w /src \
  -e GOFLAGS=-mod=mod -e GOPROXY=off -e CGO_ENABLED=0 -e HOME=/src \
  -e GOMODCACHE=/go/pkg/mod \
  golang:1.27 bash -c "mkdir -p /tmp/gt; export TMPDIR=/tmp/gt GOTMPDIR=/tmp/gt; go test ./internal/risk/ -run '$RUNPAT' -count=1 -timeout 30m -v > /tmp/gt/out.txt 2>&1; echo GO_TEST_RC=\$?; cat /tmp/gt/out.txt" > "$OUT" 2>&1
echo "wrote=$OUT bytes=$(wc -c < "$OUT")"
echo "RUN=$(grep -c '^=== RUN' "$OUT" || true)"
echo "PASS_TOP=$(grep -c '^--- PASS' "$OUT" || true)"
echo "FAIL_TOP=$(grep -c '^--- FAIL' "$OUT" || true)"
echo "SKIP=$(grep -c '^--- SKIP' "$OUT" || true)"
echo "=== red roster (top-level) ==="
grep -E '^--- FAIL' "$OUT" || true
echo "=== skip roster ==="
grep -E '^--- SKIP' "$OUT" || true
echo "=== 台面故障名册 ==="
grep -nE 'walk|build failed|no test files|setup failed|creating work dir|0xc0000135' "$OUT" || echo none
