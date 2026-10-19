#!/bin/sh
# 342-v2 必答(iv)：./internal/winsec/ 全包成对四数（POSIX 容器档），基线树⇄落地树各 2 发。
# ★必须带 -v（无 -v 时四数取不到＝尺失灵，342-r1 自抓过的那味）。
set -u
ROOT=/d/tmp/wisp/20261019/342v2
LOGS="$ROOT/logs"; mkdir -p "$LOGS"
for TREE in treeB treeL; do
  for SHOT in 1 2; do
    LABEL="iv01-${TREE}-fullpkg-${SHOT}"
    docker run --rm \
      -v "$ROOT/$TREE:/src" \
      -v /d/work/base/gopath/pkg/mod:/gm:ro \
      -v "$ROOT/logs:/logs" \
      -v "$ROOT/gocache:/root/.cache/go-build" \
      -w /src \
      -e GOPROXY=off -e GOFLAGS=-mod=mod -e CGO_ENABLED=0 -e GOMODCACHE=/gm \
      golang:1.27 \
      go test ./internal/winsec/ -v -count=1 -timeout 40m > "$LOGS/$LABEL.raw.txt" 2>&1
    RC=$?
    RUN=$(grep -c '^=== RUN' "$LOGS/$LABEL.raw.txt" || true)
    PT=$(grep -c '^--- PASS' "$LOGS/$LABEL.raw.txt" || true)
    FT=$(grep -c '^--- FAIL' "$LOGS/$LABEL.raw.txt" || true)
    SK=$(grep -c '^--- SKIP' "$LOGS/$LABEL.raw.txt" || true)
    echo "$LABEL RC=$RC RUN=$RUN PASS_TOP=$PT FAIL_TOP=$FT SKIP=$SK" | tee -a "$LOGS/four-counts-posix-fullpkg.tsv"
    grep -E '^--- FAIL' "$LOGS/$LABEL.raw.txt" | sed 's/^--- FAIL: //' | sort > "$LOGS/$LABEL.rednames.txt"
    echo "  RED_NAMES=$(wc -l < "$LOGS/$LABEL.rednames.txt")"
    grep -E 'build failed|^FAIL|no test files|setup failed|panic:' "$LOGS/$LABEL.raw.txt" | head -5
  done
done
