#!/bin/sh
# 票 338 落地腿 338-r1 · AC#3 颜色 v2（嵌套台面＝父目录有 8.3 别名，修 v1 那枚 252 族台面红）
set -u
ROOT=/d/tmp/wisp/20261018/338r1
NEST=$(cat "$ROOT/nest-v2.txt")
TB="$NEST/tree-before"; TA="$NEST/tree-after"
export GOPATH=/d/work/base/gopath
unset GOTMPDIR TMPDIR
mkdir -p "$NEST/gotmp"
export GOTMPDIR="$NEST/gotmp"
echo "NEST=$NEST  GOTMPDIR=$GOTMPDIR"
echo "别名尺: $(cmd //c "dir /x D:\\tmp\\wisp\\20261018\\338r1" 2>&1 | grep -c -- 'TREES-~1')"

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
  echo "--- FAIL 逐名 ---"; grep '^--- FAIL' "$out" || true
  echo "--- SKIP 逐名 ---"; grep '^--- SKIP' "$out" || true
  echo "--- 台面故障尺 ---"; grep -nE 'walk|build failed|no test files|setup failed|0xc0000135' "$out" || echo none
  echo "bytes=$(wc -c < "$out")"
}
run_pkg before "$TB" internal/tools "$ROOT/log-tools-v2-before.txt"
run_pkg after  "$TA" internal/tools "$ROOT/log-tools-v2-after.txt"
echo '===== 两侧 FAIL 名册作差 ====='
grep '^--- FAIL' "$ROOT/log-tools-v2-before.txt" | sed 's/ (.*//' | sort > "$ROOT/f-v2-before.txt"
grep '^--- FAIL' "$ROOT/log-tools-v2-after.txt"  | sed 's/ (.*//' | sort > "$ROOT/f-v2-after.txt"
echo "only_before=$(comm -23 "$ROOT/f-v2-before.txt" "$ROOT/f-v2-after.txt" | grep -c . || true)"
echo "only_after=$(comm -13 "$ROOT/f-v2-before.txt" "$ROOT/f-v2-after.txt" | grep -c . || true)"
comm -3 "$ROOT/f-v2-before.txt" "$ROOT/f-v2-after.txt"
