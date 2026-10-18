#!/bin/bash
# 339-r1 AC#0/AC#4 台件③：Windows 本机档跑 ./internal/risk/（仓外导出树，⛔ 共享工作树）。
# 用法：bash win.sh <run-regex> <label>
# 落点全在 D:\tmp\wisp\20261018\339r1\win\（⛔ 盘符根）；GOTMPDIR 先 mkdir -p。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
ROOT="/d/tmp/wisp/20261018/339r1"
RUNPAT="${1:-.}"
LABEL="${2:-win}"
TREE="$ROOT/win/tree"
LOGDIR="$REPO/.scratch/wisp/probes/339/r1/logs"
mkdir -p "$LOGDIR" "$ROOT/gt"
export TMPDIR="$ROOT/gt" TEMP="$ROOT/gt" TMP="$ROOT/gt" GOTMPDIR="$ROOT/gt"
export GOFLAGS=-mod=mod GOPROXY=off CGO_ENABLED=0
cd "$TREE" || exit 9
OUT="$LOGDIR/$LABEL.log"
echo "label=$LABEL tree=$TREE run=$RUNPAT pwd=$(pwd)"
go test ./internal/risk/ -run "$RUNPAT" -count=1 -timeout 30m -v > "$OUT" 2>&1
RC=$?
echo "GO_TEST_RC=$RC"
echo "bytes=$(wc -c < "$OUT")"
echo "RUN=$(grep -c '^=== RUN' "$OUT" || true)"
echo "PASS_TOP=$(grep -c '^--- PASS' "$OUT" || true)"
echo "FAIL_TOP=$(grep -c '^--- FAIL' "$OUT" || true)"
echo "SKIP=$(grep -c '^--- SKIP' "$OUT" || true)"
echo "=== red roster ==="
grep -E '^--- FAIL' "$OUT" || true
echo "=== skip roster ==="
grep -E '^--- SKIP' "$OUT" || true
echo "===台面故障名册==="
grep -nE 'walk|build failed|no test files|setup failed|creating work dir|0xc0000135' "$OUT" || echo "none"
