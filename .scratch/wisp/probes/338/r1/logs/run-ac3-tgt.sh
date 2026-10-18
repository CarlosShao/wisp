#!/bin/sh
# 票 338 落地腿 338-r1 · AC#3 定向成对（cmd/wisp＝我那枚文件里的 9 枚顶层测试）
set -u
ROOT=/d/tmp/wisp/20261018/338r1
NEST=$(cat "$ROOT/nest-v2.txt")
export GOPATH=/d/work/base/gopath
unset GOTMPDIR TMPDIR
mkdir -p "$NEST/gotmp-tgt"
export GOTMPDIR="$NEST/gotmp-tgt"
RE='/TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsAbsent|TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsNotCallable|TestForwardingHookIsSilentlyUnarmedByANonWritableNativeExit|TestForwardingHookMustNotLoseTheNativeExitReceiver|TestLegacySubShapeOneHookDiesInAReentryLoop|TestPagePostMessageEnvelopeReachesNativeExitViaShapeA3|TestShapeA3ForwardingHookIsIdempotentInOneDocument|TestShapeA3SecondPagePostStillDelivers|TestUnforwardedPageEnvelopeDiesAtTheUnboundSlot'
for tag in before after; do
  T="$NEST/tree-wisp-$tag"; out="$ROOT/log-wisp-tgt-$tag.txt"
  echo "===== targeted $tag ./cmd/wisp/ ====="
  ( cd "$T" && PATH="$T/third_party/sherpa-onnx:$T/build:$PATH" go test -count=1 -timeout 30m -v -run "$RE" ./cmd/wisp/ > "$out" 2>&1 )
  echo "rc=$?"
  echo "RUN=$(grep -c '^=== RUN' "$out" || true)"
  echo "PASS_TOP=$(grep -c '^--- PASS' "$out" || true)"
  echo "FAIL_TOP=$(grep -c '^--- FAIL' "$out" || true)"
  echo "SKIP=$(grep -c '^--- SKIP' "$out" || true)"
  grep -E '^(--- (PASS|FAIL|SKIP)|ok|FAIL|PASS)' "$out" | head -24
  echo "bytes=$(wc -c < "$out")"
done
