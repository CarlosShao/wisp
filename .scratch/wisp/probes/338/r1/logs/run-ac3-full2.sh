#!/bin/sh
# 票 338 落地腿 338-r1 · AC#3 整包成对 v2（★单差异台面＝tree-ctrl-before ⇄ tree-ctrl-after，
# 只差 panel 那一枚字节；v1 那对差 3 枚＝夹了别腿提交＝台面缺陷，读数作废、具名登记）
set -u
ROOT=/d/tmp/wisp/20261018/338r1
NEST=$(cat "$ROOT/nest-v2.txt")
export GOPATH=/d/work/base/gopath
unset GOTMPDIR TMPDIR
mkdir -p "$NEST/gotmp-full2"
export GOTMPDIR="$NEST/gotmp-full2"
for v in ctrl-before ctrl-after; do
  T="$NEST/tree-$v"; out="$ROOT/log-wisp2-$v.txt"
  echo "===== whole-package $v ./cmd/wisp/ ====="
  ( cd "$T" && PATH="$T/third_party/sherpa-onnx:$T/build:$PATH" go test -count=1 -timeout 30m -v ./cmd/wisp/ > "$out" 2>&1 )
  echo "rc=$?"
  echo "RUN=$(grep -c '^=== RUN' "$out" || true)"
  echo "PASS_TOP=$(grep -c '^--- PASS' "$out" || true)"
  echo "FAIL_TOP=$(grep -c '^--- FAIL' "$out" || true)"
  echo "SKIP=$(grep -c '^--- SKIP' "$out" || true)"
  echo '--- FAIL 逐名 ---'; grep '^--- FAIL' "$out" || true
  echo '--- SKIP 逐名 ---'; grep '^--- SKIP' "$out" || true
  echo '--- 台面故障尺 ---'; grep -nE 'walk|build failed|no test files|setup failed|0xc0000135' "$out" || echo none
  echo '--- 最后一枚顶层结果＋尾部 ---'; grep -E '^--- (PASS|FAIL|SKIP)' "$out" | tail -1; tail -2 "$out" | cut -c1-120
  echo "bytes=$(wc -c < "$out")"
done
echo '===== 名册逐名作差 ====='
grep -E '^--- (PASS|FAIL|SKIP)' "$ROOT/log-wisp2-ctrl-before.txt" | sed 's/ (.*//' | sort > "$ROOT/w2b.txt"
grep -E '^--- (PASS|FAIL|SKIP)' "$ROOT/log-wisp2-ctrl-after.txt"  | sed 's/ (.*//' | sort > "$ROOT/w2a.txt"
echo "only_before=$(comm -23 "$ROOT/w2b.txt" "$ROOT/w2a.txt" | grep -c . || true)"
echo "only_after=$(comm -13 "$ROOT/w2b.txt" "$ROOT/w2a.txt" | grep -c . || true)"
echo '--- 逐名 ---'; comm -3 "$ROOT/w2b.txt" "$ROOT/w2a.txt"
