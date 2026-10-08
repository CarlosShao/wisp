#!/bin/bash
# 35-r4 counter-shape self-check: every new detector gets (a) a mutant that must REDDEN it
# and (b) a "detector removed" run that must turn the mutant GREEN again.
# All via `go test -overlay`: nothing inside the repo is written.
set -u
REPO="/d/work/workspace/projects plans/Wisp"
OUT=/d/tmp/wisp35r4
cd "$REPO" || exit 1
mkdir -p "$OUT/mut" "$OUT/logs"
LOGS="$OUT/logs"
PROD=cmd/wisp/panel_host_windows.go
FIX=cmd/wisp/panel_transport_35r2_test.go
REPOPROD="$(cygpath -w "$REPO/$PROD")"
REPOFIX="$(cygpath -w "$REPO/$FIX")"

mk_overlay() { # $1 = name, $2.. = "winorig|winreplacement"
  local name="$1"; shift
  { printf '{"Replace":{'; local first=1 pair o r
    for pair in "$@"; do
      o="${pair%%|*}"; r="${pair##*|}"
      [ "$first" -eq 1 ] || printf ','
      printf '"%s":"%s"' "$o" "$r"; first=0
    done
    printf '}}\n'; } > "$OUT/overlay-$name.json"
}

# ---- mutant sources (each exactly one changed line) -----------------------
# ma-hook: shipped hook loses the receiver (A684 M-A).
sed 's/native\.call(cw, message)/native(message)/' "$PROD" > "$OUT/mut/ma-hook-panel_host_windows.go"
# mb-typeof-gone: the guard's typeof half deleted (A684 M-B, first half).
sed 's/if (inside || typeof window\.%\[1\]s !== "function")/if (inside)/' "$PROD" > "$OUT/mut/mb-typeof-gone-panel_host_windows.go"
# ma-detector-gone: the fixture's receiver requirement removed (used WITH ma-hook).
sed 's/if recv != cw {/if false {/' "$FIX" > "$OUT/mut/ma-detector-gone-panel_transport_35r2_test.go"
# mb-set-gone: the fixture's writability check removed.
sed 's/exists \&\& o\.unwritable\[name\] {/exists \&\& false {/' "$FIX" > "$OUT/mut/mb-set-gone-panel_transport_35r2_test.go"

M1="$(cygpath -m "$OUT/mut/ma-hook-panel_host_windows.go")"
M4="$(cygpath -m "$OUT/mut/mb-typeof-gone-panel_host_windows.go")"
M2="$(cygpath -m "$OUT/mut/ma-detector-gone-panel_transport_35r2_test.go")"
M3="$(cygpath -m "$OUT/mut/mb-set-gone-panel_transport_35r2_test.go")"

mk_overlay ma-hook          "$REPOPROD|$M1"
mk_overlay mb-typeof-gone   "$REPOPROD|$M4"
mk_overlay ma-detector-gone "$REPOPROD|$M1" "$REPOFIX|$M2"
mk_overlay mb-set-gone      "$REPOFIX|$M3"

# ---- landing check --------------------------------------------------------
{ echo "counter-shape mutant landing ($(date '+%F %T%z')): changed-line-count 2 = one line replaced"
  for spec in "ma-hook:$PROD:$OUT/mut/ma-hook-panel_host_windows.go" \
              "mb-typeof-gone:$PROD:$OUT/mut/mb-typeof-gone-panel_host_windows.go" \
              "ma-detector-gone:$FIX:$OUT/mut/ma-detector-gone-panel_transport_35r2_test.go" \
              "mb-set-gone:$FIX:$OUT/mut/mb-set-gone-panel_transport_35r2_test.go"; do
    n="${spec%%:*}"; rest="${spec#*:}"; o="${rest%%:*}"; r="${rest##*:}"
    echo "mutant $n changed-line-count=$(diff "$o" "$r" | grep -c '^[<>]')"
    diff "$o" "$r" | grep '^[<>]'
  done; } > "$LOGS/mut-landing.txt" 2>&1
echo "landing rc=$?" >> "$LOGS/mut-landing.txt"

# ---- snapshot the working tree before the overlay runs --------------------
BEFORE_PROD="$(git hash-object "$PROD")"
BEFORE_FIX="$(git hash-object "$FIX")"

# ---- run the four overlays ------------------------------------------------
RUN='TestForwarding|TestPagePostMessageEnvelopeReachesNativeExitViaShapeA3|TestShapeA3SecondPagePostStillDelivers|TestShapeA3ForwardingHookIsIdempotentInOneDocument|TestLegacySubShapeOneHookDiesInAReentryLoop|TestUnforwardedPageEnvelopeDiesAtTheUnboundSlot|TestPagePostMessageEnvelopeReachesDispatchRawViaTransport|TestInbound'
: > "$LOGS/mut-summary.txt"
for m in ma-hook ma-detector-gone mb-set-gone mb-typeof-gone; do
  L="$LOGS/mut-$m.txt"
  { echo "=== mutant $m ($(date '+%F %T%z')) overlay=$OUT/overlay-$m.json ==="
    cat "$OUT/overlay-$m.json"
    echo "--- build ---"
    go test -c -overlay "$OUT/overlay-$m.json" -o "$OUT/panel-$m.exe" ./cmd/wisp/ 2>&1
    echo "buildrc=$?"
    echo "--- run (CWD=cmd/wisp) ---"
    ( cd "$REPO/cmd/wisp" && "$OUT/panel-$m.exe" -test.run "$RUN" -test.v 2>&1 )
    echo "runrc=$?"
  } > "$L" 2>&1
  { echo "=== $m ==="
    grep -- '--- FAIL' "$L" || echo "(no FAIL line)"
    echo "counts: PASS=$(grep -c -- '--- PASS' "$L") FAIL=$(grep -c -- '--- FAIL' "$L") $(grep -E '^(buildrc|runrc)' "$L" | tr '\n' ' ')"
  } >> "$LOGS/mut-summary.txt"
done

# ---- repo untouched after all overlays ------------------------------------
{ echo "repo-untouched-after-overlays ($(date '+%F %T%z'))"
  echo "prod  before=$BEFORE_PROD after=$(git hash-object "$PROD") HEAD=$(git rev-parse HEAD:$PROD)"
  echo "fix   before=$BEFORE_FIX  after=$(git hash-object "$FIX")  (HEAD=$(git rev-parse HEAD:$FIX), edited by 35-r4 and committed separately)"
  echo "git status --porcelain on the two faces (expect only the r4 edit to FIX, or nothing once committed):"
  git status --porcelain -- "$PROD" "$FIX"
  echo "prod-untouched-by-overlays=$([ "$BEFORE_PROD" = "$(git hash-object "$PROD")" ] && echo YES || echo NO)"
  echo "fix-untouched-by-overlays=$([ "$BEFORE_FIX" = "$(git hash-object "$FIX")" ] && echo YES || echo NO)"
  echo "rc=0"; } > "$LOGS/repo-untouched-after-overlays.txt" 2>&1
