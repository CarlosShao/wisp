#!/bin/sh
# 票 338 落地腿 338-r1 · AC#3 颜色（cmd/wisp 枚＝lead 甲形条件③点名的成对口）
# 台面＝嵌套导出树（父目录有 8.3 别名）；改前⇄改后只差我这一枚＝单差异对
#   before＝87521cb4（panel 未归位／tools 已归位）  after＝93d0f39c（两枚都已归位）
# 四数尺：RUN=^=== RUN（含子测试）／PASS_TOP=^--- PASS（只顶层）／FAIL_TOP=^--- FAIL／SKIP=^--- SKIP
# ★跑 cmd/wisp 必带 PATH＝树内 third_party/sherpa-onnx＋build，否则 0xc0000135＋零枚 --- FAIL＝台面故障
set -u
ROOT=/d/tmp/wisp/20261018/338r1
NEST=$(cat "$ROOT/nest-v2.txt")
cd "/d/work/workspace/projects plans/Wisp" || exit 9
export GOPATH=/d/work/base/gopath
unset GOTMPDIR TMPDIR

build_tree() {
  rev=$1; T="$NEST/$2"
  mkdir -p "$T"
  git archive "$rev" go.mod go.sum internal cmd tools scripts frontend models | tar -x -C "$T"
  git archive "$rev" docs | tar -x -C "$T"
  mkdir -p "$T/build" "$T/third_party/sherpa-onnx" "$T/frontend/dist"
  cp -r build/. "$T/build/"; cp -r third_party/sherpa-onnx/. "$T/third_party/sherpa-onnx/"
  [ -d frontend/dist ] && cp -r frontend/dist/. "$T/frontend/dist/"
  echo "$2 rev=$rev go=$(find "$T" -name '*.go' | wc -l) dll=$(find "$T" -name '*.dll' | wc -l)"
}
build_tree 87521cb4 tree-wisp-before
build_tree 93d0f39c tree-wisp-after
echo '===== 两枚树应只差 panel 一枚 ====='
diff -rq "$NEST/tree-wisp-before" "$NEST/tree-wisp-after" 2>/dev/null | grep -v gotmp | head

mkdir -p "$NEST/gotmp-wisp"
export GOTMPDIR="$NEST/gotmp-wisp"

run_wisp() {
  tag=$1; tree=$2; out=$3
  echo "===== $tag ./cmd/wisp/ (whole package, -v) ====="
  ( cd "$tree" && PATH="$tree/third_party/sherpa-onnx:$tree/build:$PATH" \
      go test -count=1 -timeout 30m -v ./cmd/wisp/ > "$out" 2>&1 )
  echo "rc=$?"
  echo "RUN=$(grep -c '^=== RUN' "$out" || true)"
  echo "PASS_TOP=$(grep -c '^--- PASS' "$out" || true)"
  echo "FAIL_TOP=$(grep -c '^--- FAIL' "$out" || true)"
  echo "SKIP=$(grep -c '^--- SKIP' "$out" || true)"
  echo "testbinary_exit_lines=$(grep -cE '0xc0000135|exit status' "$out" || true)"
  echo '--- 台面故障尺 ---'
  grep -nE 'walk|build failed|no test files|setup failed|0xc0000135' "$out" || echo none
  echo '--- FAIL 逐名 ---'
  grep '^--- FAIL' "$out" || true
  echo '--- SKIP 逐名 ---'
  grep '^--- SKIP' "$out" || true
  echo "bytes=$(wc -c < "$out")"
}
run_wisp before "$NEST/tree-wisp-before" "$ROOT/log-wisp-before.txt"
run_wisp after  "$NEST/tree-wisp-after"  "$ROOT/log-wisp-after.txt"

echo '===== FAIL 名册逐名作差（必须 0 枚新增＝lead 条件③） ====='
grep '^--- FAIL' "$ROOT/log-wisp-before.txt" | sed 's/ (.*//' | sort > "$ROOT/fw-before.txt"
grep '^--- FAIL' "$ROOT/log-wisp-after.txt"  | sed 's/ (.*//' | sort > "$ROOT/fw-after.txt"
echo "only_before=$(comm -23 "$ROOT/fw-before.txt" "$ROOT/fw-after.txt" | grep -c . || true)"
echo "only_after=$(comm -13 "$ROOT/fw-before.txt" "$ROOT/fw-after.txt" | grep -c . || true)"
echo '--- 作差逐名 ---'; comm -3 "$ROOT/fw-before.txt" "$ROOT/fw-after.txt"
echo '--- SKIP 名册作差 ---'
grep '^--- SKIP' "$ROOT/log-wisp-before.txt" | sed 's/ (.*//' | sort > "$ROOT/sw-before.txt"
grep '^--- SKIP' "$ROOT/log-wisp-after.txt"  | sed 's/ (.*//' | sort > "$ROOT/sw-after.txt"
comm -3 "$ROOT/sw-before.txt" "$ROOT/sw-after.txt"
