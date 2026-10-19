#!/bin/sh
# 342-v2 容器驱动：在仓外导出树里跑指定的 -run，四数＋逐字 -v 输出落 log。
# 用法：sh 03-container-run.sh <treename> <label> <-run regex> [extra go test args...]
set -e
TREE="$1"; LABEL="$2"; shift 2
LOGDIR=/d/tmp/wisp/20261019/342v2/logs
mkdir -p "$LOGDIR"
RAW="$LOGDIR/$LABEL.raw.txt"
docker run --rm -i \
  -v "/d/tmp/wisp/20261019/342v2/$TREE:/src" \
  -v "/d/work/base/gopath/pkg/mod:/gm:ro" \
  -w /src \
  -e GOFLAGS=-mod=mod \
  -e GOPROXY=off \
  -e CGO_ENABLED=0 \
  -e GOMODCACHE=/gm \
  -e GOCACHE=/root/.cache/go-build \
  -e HOME=/root \
  golang:1.27 sh -c "ulimit -c 0; go test ./internal/winsec/ -v -count=1 -timeout 30m $*" > "$RAW" 2>&1
echo "RC=$?"
RUN=$(grep -c '^=== RUN' "$RAW" || true)
PT=$(grep -c '^--- PASS' "$RAW" || true)
FT=$(grep -c '^--- FAIL' "$RAW" || true)
SK=$(grep -c '^--- SKIP' "$RAW" || true)
echo "LABEL=$LABEL TREE=$TREE RUN=$RUN PASS_TOP=$PT FAIL_TOP=$FT SKIP=$SK"
echo "$LABEL TREE=$TREE RUN=$RUN PASS_TOP=$PT FAIL_TOP=$FT SKIP=$SK" >> "$LOGDIR/four-counts.tsv"
