#!/bin/sh
# 票 338 落地腿 338-r1 · AC#3 红名册逐名归位尺（整包两侧 3＋2 枚红 ⇒ 同一台面定向复跑）
set -u
ROOT=/d/tmp/wisp/20261018/338r1
NEST=$(cat "$ROOT/nest-v2.txt")
export GOPATH=/d/work/base/gopath
unset GOTMPDIR TMPDIR
mkdir -p "$NEST/gotmp-attr"
export GOTMPDIR="$NEST/gotmp-attr"
RE='^(TestAC2SealNoticeLandsInTheRunLegLogFile|TestCleanCheckoutBuilds_AC11|TestAC14GoSideEvalPushReachesThePage)$'
for v in ctrl-before ctrl-after; do
  T="$NEST/tree-$v"; out="$ROOT/log-wisp-attr-$v.txt"
  echo "===== attr $v ====="
  ( cd "$T" && PATH="$T/third_party/sherpa-onnx:$T/build:$PATH" go test -count=1 -timeout 30m -v -run "$RE" ./cmd/wisp/ > "$out" 2>&1 )
  echo "rc=$?"
  echo "RUN=$(grep -c '^=== RUN' "$out" || true) PASS_TOP=$(grep -c '^--- PASS' "$out" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$out" || true) SKIP=$(grep -c '^--- SKIP' "$out" || true)"
  grep -E '^--- (PASS|FAIL|SKIP)' "$out" || true
  echo '--- 红因首行逐名 ---'
  grep -E '^\s+\S+_test\.go:[0-9]+:' "$out" | head -12
done
