#!/bin/bash
# 339-r2 RED-before：把票 334 仪器摘掉平台档（只摘 tag，⛔ 补世界）⇒ Linux 档必须当场红。
# 这就是"仪器只在一个平台有世界"的原病，可重跑。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
BASE="/d/tmp/wisp/20261105/339r2"
LOGD="$REPO/.scratch/wisp/probes/339/r2/logs"
RUN='TestAnchorMissReasonNamesCandidateAndMissingAnchor|TestAnchorMissKeepsNonOverridableClassAndNeedL2'
mkdir -p "$LOGD"

rm -rf "$BASE/linux-red-tag"; mkdir -p "$BASE/linux-red-tag"
git -C "$REPO" archive HEAD | tar -x -C "$BASE/linux-red-tag"
F="$BASE/linux-red-tag/internal/risk/anchor_miss_reason_334_test.go"
head -3 "$F"
# 只摘掉首行 build tag（+ 其后空行），别的字节一字不动。
perl -0pi -e 's/\A\/\/go:build windows\n\n//' "$F"
echo "AFTER_HEAD:"; head -3 "$F"
echo "TAG_HITS=$(grep -c '^//go:build windows' "$F" || true)"

MSYS_NO_PATHCONV=1 docker run --rm \
  -v "$BASE/linux-red-tag":/src \
  -v "/d/work/base/gopath/pkg/mod":/go/pkg/mod:ro \
  -w /src \
  -e GOFLAGS=-mod=mod -e GOPROXY=off -e CGO_ENABLED=0 -e HOME=/src -e GOMODCACHE=/go/pkg/mod \
  golang:1.27 bash -c "mkdir -p /tmp/gt; export TMPDIR=/tmp/gt GOTMPDIR=/tmp/gt; go test ./internal/risk/ -run '$RUN' -count=1 -timeout 30m -v > /tmp/gt/o 2>&1; echo GO_TEST_RC=\$?; cat /tmp/gt/o" \
  > "$LOGD/linux-red-tag-only.log" 2>&1
O="$LOGD/linux-red-tag-only.log"
echo "rc-line: $(grep -E '^GO_TEST_RC=' "$O")"
echo "counts: RUN=$(grep -c '^=== RUN' "$O" || true) PASS=$(grep -c '^--- PASS' "$O" || true) FAIL=$(grep -c '^--- FAIL' "$O" || true) SKIP=$(grep -c '^--- SKIP' "$O" || true)"
grep -nE 'fixture guard|class: got|Gate: got|unlocked by an override' "$O" || true
