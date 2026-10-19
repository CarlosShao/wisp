#!/bin/sh
# 342-v2 必答(v) 实测尺：同一枚分类探针件打进【基线树拷贝】⇄【落地树拷贝】，POSIX 容器档。
# ★探针件⛔ 入库＝只在 D:\tmp 的 mut-E-* 拷贝里；⛔ 动跟踪字节。
set -u
ROOT=/d/tmp/wisp/20261019/342v2
LOGS="$ROOT/logs"
PROBE="$ROOT/zzprobe/probe-err-unix_test.go.txt"
for TREE in treeB treeL; do
  MUT="mut-E-err-$TREE"
  [ -d "$ROOT/$MUT" ] || cp -a "$ROOT/$TREE" "$ROOT/$MUT"
  cp -f "$PROBE" "$ROOT/$MUT/internal/winsec/probe_342v2_err_test.go"
  LABEL="v02-errclass-${TREE}"
  docker run --rm \
    -v "$ROOT/$MUT:/src" -v /d/work/base/gopath/pkg/mod:/gm:ro -v "$ROOT/gocache:/root/.cache/go-build" -w /src \
    -e GOPROXY=off -e GOFLAGS=-mod=mod -e CGO_ENABLED=0 -e GOMODCACHE=/gm \
    golang:1.27 \
    go test ./internal/winsec/ -v -count=1 -timeout 20m -run 'TestP342V2ErrClassifyPOSIX' > "$LOGS/$LABEL.raw.txt" 2>&1
  echo "$LABEL RC=$?"
  grep -E 'C[0-9]+b? |--- (FAIL|PASS|SKIP)' "$LOGS/$LABEL.raw.txt" | sed -E 's/^ +[a-z0-9_]+\.go:[0-9]+: //' | head -20
  echo "-----"
done
