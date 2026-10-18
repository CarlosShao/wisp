#!/bin/bash
# 编排者自跑：票 334 那两枚用例在 GOOS=linux 档红不红（⛔ 抄 322-r2 的读数）
set -u
REPO="/d/work/workspace/projects plans/Wisp"
WORK="/d/tmp/wisp/20261018/orch/linuxred334"
cd "$REPO" || exit 9
mkdir -p "$WORK"
rm -rf "$WORK/tree" "$WORK/tree.tar" 2>/dev/null
git archive HEAD go.mod go.sum internal cmd tools scripts frontend models docs > "$WORK/tree.tar"
mkdir -p "$WORK/tree" && tar -xf "$WORK/tree.tar" -C "$WORK/tree"
echo "go files in tree: $(find "$WORK/tree" -name '*.go' | wc -l)"
echo "334 test present: $(ls "$WORK/tree/internal/risk/" | grep -c anchor_miss_reason_334)"
# 宿主 module 缓存 → 容器内 GOPATH
CACHE="$(go env GOMODCACHE)"
echo "host modcache: $CACHE"
MSYS_NO_PATHCONV=1 docker run --rm \
  -v "$WORK/tree":/src \
  -v "$CACHE":/go/pkg/mod:ro \
  -w /src \
  -e GOFLAGS=-mod=mod -e GOPROXY=off -e CGO_ENABLED=0 -e HOME=/src \
  -e GOMODCACHE=/go/pkg/mod \
  golang:1.27 bash -c 'mkdir -p /tmp/gt; export TMPDIR=/tmp/gt GOTMPDIR=/tmp/gt; go test ./internal/risk/ -run "TestAnchorMissReasonNamesCandidateAndMissingAnchor|TestAnchorMissKeepsNonOverridableClassAndNeedL2" -count=1 -timeout 30m -v 2>&1 | tail -45; echo "PIPESTATUS_RC=${PIPESTATUS[0]}"'
