#!/bin/bash
# Orchestrator re-run of 35-r4's counter-shape driver (the leg died on the daily quota
# before any mutant produced a reading).
# ONLY CHANGE vs .scratch/../countershape.sh: the two repo-side paths go through
# `cygpath -m` instead of `cygpath -w`, because Go rejects a single backslash inside a
# JSON string ("parsing overlay JSON: invalid escape sequence `\w`" => buildrc=1, runrc=127).
# Mutant sources, the -run roster and the assertions are the leg's, byte-for-byte.
set -u
REPO="/d/work/workspace/projects plans/Wisp"
OUT=/d/tmp/wisp35r4
cd "$REPO" || exit 1
LOGS="$OUT/logs-orch"
mkdir -p "$LOGS"
PROD=cmd/wisp/panel_host_windows.go
FIX=cmd/wisp/panel_transport_35r2_test.go
REPOPROD="$(cygpath -m "$REPO/$PROD")"
REPOFIX="$(cygpath -m "$REPO/$FIX")"

mk_overlay() { # $1 = name, $2.. = "path|replacement"
  local name="$1"; shift
  { printf '{"Replace":{'; local first=1 pair o r
    for pair in "$@"; do
      o="${pair%%|*}"; r="${pair##*|}"
      [ "$first" -eq 1 ] || printf ','
      printf '"%s":"%s"' "$o" "$r"; first=0
    done
    printf '}}\n'; } > "$OUT/overlay-orch-$name.json"
}

M1="$(cygpath -m "$OUT/mut/ma-hook-panel_host_windows.go")"
M4="$(cygpath -m "$OUT/mut/mb-typeof-gone-panel_host_windows.go")"
M2="$(cygpath -m "$OUT/mut/ma-detector-gone-panel_transport_35r2_test.go")"
M3="$(cygpath -m "$OUT/mut/mb-set-gone-panel_transport_35r2_test.go")"

mk_overlay ma-hook          "$REPOPROD|$M1"
mk_overlay mb-typeof-gone   "$REPOPROD|$M4"
mk_overlay ma-detector-gone "$REPOPROD|$M1" "$REPOFIX|$M2"
mk_overlay mb-set-gone      "$REPOFIX|$M3"

# ---- re-confirm the four mutants are still exactly one changed line -------
{ echo "orchestrator mutant re-landing ($(date '+%F %T%z')): counts taken against the CURRENT working tree"
  for spec in "ma-hook:$PROD:$OUT/mut/ma-hook-panel_host_windows.go" \
              "mb-typeof-gone:$PROD:$OUT/mut/mb-typeof-gone-panel_host_windows.go" \
              "ma-detector-gone:$FIX:$OUT/mut/ma-detector-gone-panel_transport_35r2_test.go" \
              "mb-set-gone:$FIX:$OUT/mut/mb-set-gone-panel_transport_35r2_test.go"; do
    n="${spec%%:*}"; rest="${spec#*:}"; o="${rest%%:*}"; r="${rest##*:}"
    echo "mutant $n changed-line-count=$(diff "$o" "$r" | grep -c '^[<>]')"
    diff "$o" "$r" | grep '^[<>]'
  done; } > "$LOGS/orch-mut-landing.txt" 2>&1

BEFORE_PROD="$(git hash-object "$PROD")"
BEFORE_FIX="$(git hash-object "$FIX")"

RUN='TestForwarding|TestPagePostMessageEnvelopeReachesNativeExitViaShapeA3|TestShapeA3SecondPagePostStillDelivers|TestShapeA3ForwardingHookIsIdempotentInOneDocument|TestLegacySubShapeOneHookDiesInAReentryLoop|TestUnforwardedPageEnvelopeDiesAtTheUnboundSlot|TestPagePostMessageEnvelopeReachesDispatchRawViaTransport|TestInbound'
: > "$LOGS/orch-mut-summary.txt"
for m in ma-hook ma-detector-gone mb-set-gone mb-typeof-gone; do
  L="$LOGS/orch-mut-$m.txt"
  { echo "=== orchestrator mutant $m ($(date '+%F %T%z')) overlay=$OUT/overlay-orch-$m.json ==="
    cat "$OUT/overlay-orch-$m.json"
    echo "--- build ---"
    go test -c -overlay "$OUT/overlay-orch-$m.json" -o "$OUT/panel-orch-$m.exe" ./cmd/wisp/ 2>&1
    echo "buildrc=$?"
    echo "--- run (CWD=cmd/wisp) ---"
    ( cd "$REPO/cmd/wisp" && "$OUT/panel-orch-$m.exe" -test.run "$RUN" -test.v 2>&1 )
    echo "runrc=$?"
  } > "$L" 2>&1
  { echo "=== $m ==="
    grep -- '--- FAIL' "$L" || echo "(no FAIL line)"
    echo "counts: PASS=$(grep -c -- '--- PASS' "$L") FAIL=$(grep -c -- '--- FAIL' "$L") $(grep -E '^(buildrc|runrc)' "$L" | tr '\n' ' ')"
  } >> "$LOGS/orch-mut-summary.txt"
done

{ echo "orchestrator repo-untouched after overlays ($(date '+%F %T%z'))"
  echo "prod  before=$BEFORE_PROD after=$(git hash-object "$PROD")"
  echo "fix   before=$BEFORE_FIX after=$(git hash-object "$FIX")"
  echo "prod-untouched=$([ "$BEFORE_PROD" = "$(git hash-object "$PROD")" ] && echo YES || echo NO)"
  echo "fix-untouched=$([ "$BEFORE_FIX" = "$(git hash-object "$FIX")" ] && echo YES || echo NO)"
  echo "rc=0"; } > "$LOGS/orch-repo-untouched.txt" 2>&1
echo "driver done rc=0"
