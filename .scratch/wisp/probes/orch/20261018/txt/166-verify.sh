#!/bin/bash
# 编排者收件复跑（166-r1 决定性那一发）：四枚装配可达用例在 HEAD 必须绿；基线 c11163c8 上必须"没有这些用例"
set -u
REPO="/d/work/workspace/projects plans/Wisp"
W="/d/tmp/wisp/20261018/orch/166-verify"
mkdir -p "$W/head/tree" "$W/base/tree" "$W/gotmp" "$W/logs" "$W/tdst"
cd "$REPO" || exit 9
{
for pair in "head:HEAD" "base:c11163c8"; do
  name="${pair%%:*}"; rev="${pair##*:}"
  echo "=== [$name] export $rev ==="
  git archive "$rev" go.mod go.sum internal cmd tools scripts frontend models docs | tar -x -C "$W/$name/tree"
  cd "$W/$name/tree" || exit 9
  echo "go files = $(find . -name '*.go' | wc -l)"
  export GOTMPDIR="$W/gotmp" TMPDIR="$W/tdst" TEMP="$W/tdst" TMP="$W/tdst"
  echo "--- assembly tests ---"
  go test ./internal/session/ -run 'TestAC2Outlet' -count=1 -timeout 30m -v 2>&1 | grep -E '^(=== RUN|--- (PASS|FAIL)|ok  |FAIL|no test|testing:)' | head -10
  echo "RC=${PIPESTATUS[0]}"
  echo "--- whole package four numbers ---"
  go test ./internal/session/ -count=1 -timeout 30m 2>&1 | tail -3
  cd "$REPO" || exit 9
done
echo "=== wiring ruler (both revs) ==="
echo "HEAD NewManager callers = $(git grep -c 'session\.NewManager' HEAD -- '*.go' 2>/dev/null | wc -l)"
echo "base NewManager callers = $(git grep -c 'session\.NewManager' c11163c8 -- '*.go' 2>/dev/null | wc -l)"
} > "$W/logs/orch-166-recheck.log" 2>&1
echo "LOG=$W/logs/orch-166-recheck.log"
cat "$W/logs/orch-166-recheck.log"
