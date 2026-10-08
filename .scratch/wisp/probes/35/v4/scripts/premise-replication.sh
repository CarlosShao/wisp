#!/bin/bash
# 35-v4 premise replication: A684 §3 / ticket 35:75 say 35-v2's two mutants stayed GREEN on the
# pre-r4 fixture. Nobody in this wave re-measured that - the leg and the orchestrator both took
# 35-v2's word. Here it is measured: overlay the CURRENT production file (mutated one line) onto
# the pre-r4 fixture blob (git show 2fc5f5c9^:cmd/wisp/panel_transport_35r2_test.go, 1824 lines)
# and run the roster. The three r4 cases do not exist in that blob, so their names vanishing from
# the roster is expected and is recorded, not read as a colour.
set -u
REPO="/d/work/workspace/projects plans/Wisp"
OUT=/d/tmp/wisp35v4
L="$REPO/.scratch/wisp/probes/35/v4/logs"
cd "$REPO" || exit 1
RP="$(cygpath -m "$REPO/cmd/wisp/panel_host_windows.go")"
RF="$(cygpath -m "$REPO/cmd/wisp/panel_transport_35r2_test.go")"
OLD="$OUT/prer4-fixture.go"
git show '2fc5f5c9^:cmd/wisp/panel_transport_35r2_test.go' > "$OLD" || exit 1
wc -l "$OLD"
OLDM="$(cygpath -m "$OLD")"
MH="$(cygpath -m "$OUT/mut/v-ma-hook-cmd/wisp/panel_host_windows.go")"
GG="$(cygpath -m "$OUT/mut/v-guard-gone-cmd/wisp/panel_host_windows.go")"
printf '{"Replace":{"%s":"%s"}}\n' "$RF" "$OLDM" > "$OUT/overlay4-prer4-clean.json"
printf '{"Replace":{"%s":"%s","%s":"%s"}}\n' "$RF" "$OLDM" "$RP" "$MH" > "$OUT/overlay4-prer4+ma-hook.json"
printf '{"Replace":{"%s":"%s","%s":"%s"}}\n' "$RF" "$OLDM" "$RP" "$GG" > "$OUT/overlay4-prer4+guard-gone.json"
RUNPAT='^(TestPagePostMessageEnvelopeReachesDispatchRawViaTransport|TestInboundRosterGuardRefusesUnregisteredNameOnPageEdge|TestInboundSourceGuardRefusesForeignSourceOnPageEdge|TestInboundRequestIDGuardRefusesMissingIDOnPageEdge|TestPagePostMessageEnvelopeReachesNativeExitViaShapeA3|TestShapeA3SecondPagePostStillDelivers|TestShapeA3ForwardingHookIsIdempotentInOneDocument|TestLegacySubShapeOneHookDiesInAReentryLoop|TestUnforwardedPageEnvelopeDiesAtTheUnboundSlot|TestForwardingHookMustNotLoseTheNativeExitReceiver|TestForwardingHookIsSilentlyUnarmedByANonWritableNativeExit|TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsAbsent)$'
{ echo "35-v4 premise replication on the pre-r4 fixture ($(date '+%F %T%z'))"
  echo "fixture blob used = $(git hash-object "$OLD") (= git show 2fc5f5c9^:... , lines $(wc -l < "$OLD"))"
  for m in prer4-clean prer4+ma-hook prer4+guard-gone; do
    echo "=== $m ==="
    go test -c -overlay "$OUT/overlay4-$m.json" -o "$OUT/pkg4-$m.exe" ./cmd/wisp/ 2>&1
    echo "buildrc=$?"
    ( cd "$REPO/cmd/wisp" && "$OUT/pkg4-$m.exe" -test.run "$RUNPAT" -test.v 2>&1 )
    echo "runrc=$?"
  done; } > "$L/prem-replication.txt" 2>&1
{ echo "35-v4 premise replication roster table ($(date '+%F %T%z'))"
  for m in prer4-clean prer4+ma-hook prer4+guard-gone; do
    f="$L/prem-replication.txt"
    seg=$(awk -v m="$m" '$0=="=== "m" ==="{p=1;next} /^=== /{p=0} p' "$f")
    echo "$m: PASS=$(printf '%s' "$seg" | grep -ac -- '--- PASS' || true) FAIL=$(printf '%s' "$seg" | grep -ac -- '--- FAIL' || true) buildrc=$(printf '%s' "$seg" | grep -a '^buildrc' | tr -d '\r') runrc=$(printf '%s' "$seg" | grep -a '^runrc' | tr -d '\r')"
    printf '%s' "$seg" | grep -aoE -- '--- FAIL: [A-Za-z0-9_]+' | sed 's/^/    /'
    echo "    present-in-roster=$(printf '%s' "$seg" | grep -acoE -- '--- (PASS|FAIL): [A-Za-z0-9_]+' || true) of 12 named; r4's three cases are absent from that blob by construction"
  done
  echo "rc=0"; } > "$L/prem-replication-table.txt" 2>&1
echo done
