#!/bin/sh
# 票 338 落地腿 338-r1 · AC#3 颜色（仓外导出树，改前⇄改后，同一配方；四数逐枚点名⛔ 混用⛔ 相减）
# 尺：RUN=^=== RUN（含子测试）／PASS_TOP=^--- PASS／FAIL_TOP=^--- FAIL／SKIP=^--- SKIP
set -u
ROOT=/d/tmp/wisp/20261018/338r1
TB=$(grep '^before=' "$ROOT/tree-paths.txt" | tail -1 | cut -d= -f2)
TA=$(grep '^after=' "$ROOT/tree-paths.txt" | tail -1 | cut -d= -f2)
export GOPATH=/d/work/base/gopath
unset GOTMPDIR TMPDIR
mkdir -p "$ROOT/gotmp"
export GOTMPDIR="$ROOT/gotmp"

run_pkg() {
  tag=$1; tree=$2; pkg=$3; out=$4
  echo "===== $tag $pkg ====="
  ( cd "$tree" && PATH="$tree/third_party/sherpa-onnx:$tree/build:$PATH" \
      go test -count=1 -timeout 30m -v "./$pkg/" > "$out" 2>&1 )
  echo "rc=$?"
  echo "RUN=$(grep -c '^=== RUN' "$out" || true)"
  echo "PASS_TOP=$(grep -c '^--- PASS' "$out" || true)"
  echo "FAIL_TOP=$(grep -c '^--- FAIL' "$out" || true)"
  echo "SKIP=$(grep -c '^--- SKIP' "$out" || true)"
  echo "--- FAIL 逐名 ---"
  grep '^--- FAIL' "$out" || true
  echo "--- SKIP 逐名 ---"
  grep '^--- SKIP' "$out" || true
  echo "--- 台面故障尺（walk|build failed|no test files|setup failed|0xc0000135） ---"
  grep -nE 'walk|build failed|no test files|setup failed|0xc0000135' "$out" || echo "none"
  echo "bytes=$(wc -c < "$out")"
}

run_pkg before "$TB" internal/tools "$ROOT/log-tools-before.txt"
run_pkg after  "$TA" internal/tools "$ROOT/log-tools-after.txt"
