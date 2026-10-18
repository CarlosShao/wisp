#!/bin/bash
# 编排者收件复跑（339-r1 尺③）：Linux 容器里跑 ./internal/risk/，验"两枚 334 用例转绿 + 包能编译"
set -u
REPO="/d/work/workspace/projects plans/Wisp"
W="/d/tmp/wisp/20261018/orch/339-verify"
mkdir -p "$W/tree" "$W/logs"
cd "$REPO" || exit 9
{
echo "### 编排者复跑 339-r1 决定性那一发（Linux 容器）"
for pair in "base:e6fd0f50" "final:HEAD"; do
  name="${pair%%:*}"; rev="${pair##*:}"
  rm -rf "$W/$name"; mkdir -p "$W/$name/tree" "$W/$name/gotmp"
  echo "=== [$name] rev=$(git rev-parse --short $rev) ==="
  git archive "$rev" go.mod go.sum internal cmd tools scripts frontend models docs | tar -x -C "$W/$name/tree"
  MSYS_NO_PATHCONV=1 docker run --rm \
    -v "$W/$name/tree":/src \
    -v "/d/work/base/gopath/pkg/mod":/go/pkg/mod:ro \
    -w /src \
    -e GOFLAGS=-mod=mod -e GOPROXY=off -e CGO_ENABLED=0 -e HOME=/src \
    -e GOMODCACHE=/go/pkg/mod -e GOCACHE=/src/.gocache \
    golang:1.27 bash -c "mkdir -p /src/gt /src/.gocache; export TMPDIR=/src/gt GOTMPDIR=/src/gt; echo '--- go vet (Linux build must compile) ---'; go vet ./internal/risk/ 2>&1 | tail -5; echo VET_RC=\$?; echo '--- full package four numbers ---'; go test ./internal/risk/ -count=1 -timeout 30m -v > /src/pkg.log 2>&1; echo TEST_RC=\$?; echo RUN=\$(grep -c '^=== RUN' /src/pkg.log) PASS_TOP=\$(grep -c '^--- PASS' /src/pkg.log) FAIL_TOP=\$(grep -c '^--- FAIL' /src/pkg.log) SKIP=\$(grep -c '^--- SKIP' /src/pkg.log); echo '--- FAIL names ---'; grep '^--- FAIL' /src/pkg.log | sed 's/([0-9.]*s)//'; echo '--- the two 334 tests specifically ---'; grep -E '^(    )?--- (PASS|FAIL): TestAnchorMiss' /src/pkg.log | cut -c1-95"
done
} > "$W/logs/orch-339-recheck.log" 2>&1
echo "LOG=$W/logs/orch-339-recheck.log"
cat "$W/logs/orch-339-recheck.log"
