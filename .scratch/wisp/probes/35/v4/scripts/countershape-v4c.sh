#!/bin/bash
# 35-v4 driver part 3: (i) the fixed "detector keeps counting but stops throwing" shape - pass 2's
# v-no-panic did not compile (buildrc=1) because line 1291 closes the &jsPanic{...} literal, so the
# one-line edit has to keep the composite literal and drop only the panic; (ii) the re-entry cap's
# two directions, which is ledger #127's "再入帽数得对不对" question.
set -u
REPO="/d/work/workspace/projects plans/Wisp"
OUT=/d/tmp/wisp35v4
LOGS="$REPO/.scratch/wisp/probes/35/v4/logs"
PROD=cmd/wisp/panel_host_windows.go
FIX=cmd/wisp/panel_transport_35r2_test.go
REPOPROD="$(cygpath -m "$REPO/$PROD")"; REPOFIX="$(cygpath -m "$REPO/$FIX")"
cd "$REPO" || exit 1
mkfix() { mkdir -p "$(dirname "$4")"; awk -v n="$2" -v new="$3" 'NR==n{print new; next}{print}' "$REPO/$FIX" > "$4"; }
mkprod() { mkdir -p "$(dirname "$4")"; awk -v n="$2" -v new="$3" 'NR==n{print new; next}{print}' "$REPO/$PROD" > "$4"; }
M_NOP="$(cygpath -m "$OUT/mut3/v-no-panic-fix-$FIX")"
M_HOOK="$(cygpath -m "$OUT/mut3/v-ma-hook-$PROD")"
M_CAP2="$(cygpath -m "$OUT/mut3/v-cap-loose-$FIX")"
M_CAP1="$(cygpath -m "$OUT/mut3/v-cap-tight-$FIX")"
mkfix v-no-panic-fix 1290 '\t\t\t_ = (&jsPanic{"TypeError: Illegal invocation: chrome.webview.postMessage lost its receiver - called with " +' "$M_NOP"
mkprod v-ma-hook 680 '    if (inside || typeof window.%[1]s !== "function") { return native(message); }' "$M_HOOK"
mkfix v-cap-loose 1184 'const postMessageHopCap = 2' "$M_CAP2"
mkfix v-cap-tight 1184 'const postMessageHopCap = 1' "$M_CAP1"
printf '{"Replace":{"%s":"%s","%s":"%s"}}\n' "$REPOPROD" "$M_HOOK" "$REPOFIX" "$M_NOP" > "$OUT/overlay3-v-ma-hook+no-panic.json"
printf '{"Replace":{"%s":"%s"}}\n' "$REPOFIX" "$M_CAP2" > "$OUT/overlay3-v-cap-loose.json"
printf '{"Replace":{"%s":"%s"}}\n' "$REPOFIX" "$M_CAP1" > "$OUT/overlay3-v-cap-tight.json"
{ echo "35-v4 driver3 mutant landing ($(date '+%F %T%z')) - each must be exactly one changed line"
  for spec in "v-no-panic-fix:$FIX:$M_NOP" "v-ma-hook:$PROD:$M_HOOK" "v-cap-loose:$FIX:$M_CAP2" "v-cap-tight:$FIX:$M_CAP1"; do
    n="${spec%%:*}"; rest="${spec#*:}"; o="$REPO/${rest%%:*}"; r="${rest##*:}"
    echo "mutant $n changed-line-count=$(diff "$o" "$r" | grep -c '^[<>]')"
    diff "$o" "$r" | grep '^[<>]' | cat -A | cut -c1-190
  done; echo "rc=0"; } > "$LOGS/mut-landing3.txt" 2>&1
RUN="$(grep -m1 -o "RUN='\^([^']*'" "$REPO/.scratch/wisp/probes/35/v4/scripts/countershape-v4.sh" >/dev/null 2>&1; echo skip)"
RUNPAT='^(TestPagePostMessageEnvelopeReachesDispatchRawViaTransport|TestInboundRosterGuardRefusesUnregisteredNameOnPageEdge|TestInboundSourceGuardRefusesForeignSourceOnPageEdge|TestInboundRequestIDGuardRefusesMissingIDOnPageEdge|TestPagePostMessageEnvelopeReachesNativeExitViaShapeA3|TestShapeA3SecondPagePostStillDelivers|TestShapeA3ForwardingHookIsIdempotentInOneDocument|TestLegacySubShapeOneHookDiesInAReentryLoop|TestUnforwardedPageEnvelopeDiesAtTheUnboundSlot|TestForwardingHookMustNotLoseTheNativeExitReceiver|TestForwardingHookIsSilentlyUnarmedByANonWritableNativeExit|TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsAbsent)$'
BASE="$LOGS/names-baseline.txt"
{ echo "35-v4 driver3 ($(date '+%F %T%z')): 3 shapes, same 12-name roster as countershape-v4.sh pass 2"
  for m in "v-ma-hook+no-panic:overlay3-v-ma-hook+no-panic.json" "v-cap-loose:overlay3-v-cap-loose.json" "v-cap-tight:overlay3-v-cap-tight.json"; do
    n="${m%%:*}"; ov="$OUT/${m##*:}"; L="$LOGS/run3-$n.txt"
    { echo "=== run3 $n ($(date '+%F %T%z')) overlay=$ov ==="; cat "$ov"; echo "--- build ---"
      go test -c -overlay "$ov" -o "$OUT/pkg3-$n.exe" ./cmd/wisp/ 2>&1; echo "buildrc=$?"
      ( cd "$REPO/cmd/wisp" && "$OUT/pkg3-$n.exe" -test.run "$RUNPAT" -test.v 2>&1 ); echo "runrc=$?"; } > "$L" 2>&1
    grep -aoE -- '--- (FAIL|PASS|SKIP): [A-Za-z0-9_]+' "$L" | sed 's/--- \([A-Z]*\): /\1 /' > "$LOGS/names3-$n.txt"
    echo "--- $n buildrc=$(grep -a '^buildrc' "$L" | tr -d '\r') runrc=$(grep -a '^runrc' "$L" | tr -d '\r') PASS=$(grep -c '^PASS ' "$LOGS/names3-$n.txt") FAIL=$(grep -c '^FAIL ' "$LOGS/names3-$n.txt")"
    echo "  FAIL names:"; grep '^FAIL ' "$LOGS/names3-$n.txt" | sed 's/^/    /'
    echo "  vanished vs baseline (a missing name = it did not run):"; comm -13 <(cut -d' ' -f2- "$BASE" | sort -u) <(cut -d' ' -f2- "$LOGS/names3-$n.txt" | sort -u) | sed 's/^/    /'
    echo "  red-reasons:"; grep -a -A3 -- '--- FAIL' "$L" | grep -aoE '(M-A|M-B|AC#6|COUNTER-SHAPE|calibration|READING WRONG|RangeError)[^"]{0,100}' | sort -u | sed 's/^/    /'
  done
  echo "rc=0"; } > "$LOGS/mut-summary3.txt" 2>&1
echo "driver3 done"
