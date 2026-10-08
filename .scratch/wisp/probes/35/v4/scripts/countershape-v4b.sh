#!/bin/bash
# 35-v4 driver part 2: the four shapes that judge "恰红 1 枚" name-uniqueness, plus the two
# replications of leg r4's own mutants (mb-typeof-gone = v-guard-gone, mb-set-gone = v-set-gone).
# Mutant sources and overlay JSONs were already created by countershape-v4.sh (same 15 one-line
# mutants, ruler in logs/mut-landing.txt); this file only builds and runs them.
set -u
REPO="/d/work/workspace/projects plans/Wisp"
OUT=/d/tmp/wisp35v4
LOGS="$REPO/.scratch/wisp/probes/35/v4/logs"
cd "$REPO" || exit 1
RUN='^(TestPagePostMessageEnvelopeReachesDispatchRawViaTransport|TestInboundRosterGuardRefusesUnregisteredNameOnPageEdge|TestInboundSourceGuardRefusesForeignSourceOnPageEdge|TestInboundRequestIDGuardRefusesMissingIDOnPageEdge|TestPagePostMessageEnvelopeReachesNativeExitViaShapeA3|TestShapeA3SecondPagePostStillDelivers|TestShapeA3ForwardingHookIsIdempotentInOneDocument|TestLegacySubShapeOneHookDiesInAReentryLoop|TestUnforwardedPageEnvelopeDiesAtTheUnboundSlot|TestForwardingHookMustNotLoseTheNativeExitReceiver|TestForwardingHookIsSilentlyUnarmedByANonWritableNativeExit|TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsAbsent)$'

runpkg() {
  local n="$1" ov="$2" L="$LOGS/run2-$1.txt"
  { echo "=== run2 $n ($(date '+%F %T%z')) overlay=${ov:-none} roster=12 transport names ==="
    [ -n "$ov" ] && cat "$ov"
    echo "--- build ---"
    if [ -n "$ov" ]; then go test -c -overlay "$ov" -o "$OUT/pkg2-$n.exe" ./cmd/wisp/ 2>&1
    else go test -c -o "$OUT/pkg2-$n.exe" ./cmd/wisp/ 2>&1; fi
    echo "buildrc=$?"
    ( cd "$REPO/cmd/wisp" && "$OUT/pkg2-$n.exe" -test.run "$RUN" -test.v 2>&1 )
    echo "runrc=$?"
  } > "$L" 2>&1
  grep -aoE -- '--- (FAIL|PASS|SKIP): [A-Za-z0-9_]+' "$L" | sed 's/--- \([A-Z]*\): /\1 /' > "$LOGS/names2-$n.txt"
  echo "$n buildrc=$(grep -a '^buildrc' "$L" | tr -d '\r') PASS=$(grep -c '^PASS ' "$LOGS/names2-$n.txt") FAIL=$(grep -c '^FAIL ' "$LOGS/names2-$n.txt")"
}

{ echo "35-v4 driver2 ($(date '+%F %T%z')): name-uniqueness shapes. Red attribution is read by TEST NAME only; the two label mutants rewrite the red SENTENCE so both cases say the same thing."
  runpkg baseline ""
  BASE="$LOGS/names2-baseline.txt"
  for m in v-guard-gone v-set-gone v-guard-gone+label-mb v-ma-hook+label-ma; do
    echo "--- $m"
    runpkg "$m" "$OUT/overlay-$m.json"
    echo "  FAIL names:"; grep '^FAIL ' "$LOGS/names2-$m.txt" | sed 's/^/    /'
    echo "  newly-FAIL vs baseline:"; comm -13 <(grep '^FAIL ' "$BASE" | sort -u) <(grep '^FAIL ' "$LOGS/names2-$m.txt" | sort -u) | sed 's/^/    /'
    echo "  red-reasons:"; grep -a -A3 -- '--- FAIL' "$LOGS/run2-$m.txt" | grep -aoE '(M-A|M-B|AC#6|COUNTER-SHAPE|calibration|READING WRONG)[^"]{0,90}' | sort -u | sed 's/^/    /'
  done
  echo "rc=0"; } > "$LOGS/mut-summary2.txt" 2>&1
echo "driver2 done"
