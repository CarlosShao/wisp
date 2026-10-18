#!/bin/sh
# 票 338 落地腿 338-r1 · AC#3 定向成对 v2（两处修正，逐字具名）：
#  修①＝v1 的 -run 串以 '/' 开头 ⇒ Go 把 pattern 按 '/' 切段，首段为空＝匹配一切＝根本没过滤
#        （v1 读数 RUN=154/153 实为整包，⛔ 当定向读数用）。本版改成不带斜杠的 ^(A|B|…)$ 形。
#  修②＝v1 那对树差 3 枚（87521cb4⇄93d0f39c 之间夹了别腿提交：docs/reports 台账＋
#        internal/config/writeguard_handedit_227_test.go）⇒ 不是单差异对。
#        本版＝同一棵树（93d0f39c）只把那枚 panel 件换成 87521cb4 的字节＝真·单差异对。
set -u
ROOT=/d/tmp/wisp/20261018/338r1
NEST=$(cat "$ROOT/nest-v2.txt")
cd "/d/work/workspace/projects plans/Wisp" || exit 9
export GOPATH=/d/work/base/gopath
unset GOTMPDIR TMPDIR
mkdir -p "$NEST/gotmp-tgt2"
export GOTMPDIR="$NEST/gotmp-tgt2"
P=cmd/wisp/panel_transport_35r2_test.go

for v in ctrl-before ctrl-after; do
  T="$NEST/tree-$v"
  [ -d "$T" ] || { mkdir -p "$T"; git archive 93d0f39c go.mod go.sum internal cmd tools scripts frontend models | tar -x -C "$T"; git archive 93d0f39c docs | tar -x -C "$T"; mkdir -p "$T/build" "$T/third_party/sherpa-onnx" "$T/frontend/dist"; cp -r build/. "$T/build/"; cp -r third_party/sherpa-onnx/. "$T/third_party/sherpa-onnx/"; [ -d frontend/dist ] && cp -r frontend/dist/. "$T/frontend/dist/"; }
  echo "$v go=$(find "$T" -name '*.go' | wc -l)"
done
git show "87521cb4:$P" > "$NEST/tree-ctrl-before/$P"
git show "93d0f39c:$P" > "$NEST/tree-ctrl-after/$P"
echo '===== 单差异核验（应只有 panel 一枚，且 gotmp 除外） ====='
diff -rq "$NEST/tree-ctrl-before" "$NEST/tree-ctrl-after" 2>/dev/null | grep -v gotmp
echo "before_panel_sha=$(sha256sum < "$NEST/tree-ctrl-before/$P" | cut -c1-16)"
echo "after_panel_sha=$(sha256sum < "$NEST/tree-ctrl-after/$P" | cut -c1-16)"

RE='^(TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsAbsent|TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsNotCallable|TestForwardingHookIsSilentlyUnarmedByANonWritableNativeExit|TestForwardingHookMustNotLoseTheNativeExitReceiver|TestLegacySubShapeOneHookDiesInAReentryLoop|TestPagePostMessageEnvelopeReachesNativeExitViaShapeA3|TestShapeA3ForwardingHookIsIdempotentInOneDocument|TestShapeA3SecondPagePostStillDelivers|TestUnforwardedPageEnvelopeDiesAtTheUnboundSlot)$'
for v in ctrl-before ctrl-after; do
  T="$NEST/tree-$v"; out="$ROOT/log-wisp-tgt2-$v.txt"
  echo "===== targeted2 $v ./cmd/wisp/ ====="
  ( cd "$T" && PATH="$T/third_party/sherpa-onnx:$T/build:$PATH" go test -count=1 -timeout 30m -v -run "$RE" ./cmd/wisp/ > "$out" 2>&1 )
  echo "rc=$?"
  echo "RUN=$(grep -c '^=== RUN' "$out" || true)"
  echo "PASS_TOP=$(grep -c '^--- PASS' "$out" || true)"
  echo "FAIL_TOP=$(grep -c '^--- FAIL' "$out" || true)"
  echo "SKIP=$(grep -c '^--- SKIP' "$out" || true)"
  echo '--- 逐名 ---'
  grep -E '^--- (PASS|FAIL|SKIP)' "$out" || true
  echo '--- 台面故障尺 ---'
  grep -nE 'walk|build failed|no test files|setup failed|0xc0000135' "$out" || echo none
  echo "bytes=$(wc -c < "$out")"
done
echo '===== 定向两侧名册作差 ====='
grep -E '^--- (PASS|FAIL)' "$ROOT/log-wisp-tgt2-ctrl-before.txt" | sed 's/ (.*//' | sort > "$ROOT/t2b.txt"
grep -E '^--- (PASS|FAIL)' "$ROOT/log-wisp-tgt2-ctrl-after.txt"  | sed 's/ (.*//' | sort > "$ROOT/t2a.txt"
echo "only_before=$(comm -23 "$ROOT/t2b.txt" "$ROOT/t2a.txt" | grep -c . || true)"
echo "only_after=$(comm -13 "$ROOT/t2b.txt" "$ROOT/t2a.txt" | grep -c . || true)"
comm -3 "$ROOT/t2b.txt" "$ROOT/t2a.txt"
