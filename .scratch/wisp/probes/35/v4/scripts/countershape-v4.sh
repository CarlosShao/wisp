#!/bin/bash
# 35-v4 adversarial counter-shape driver (validation leg; NOT the leg r4 wrote).
#
# Differences vs .scratch/wisp/probes/35/r4/scripts/orchestrator-countershape.sh:
#   1. OUT=/d/tmp/wisp35v4 (that driver wrote into another leg's dir /d/tmp/wisp35r4).
#   2. Every mutant source is re-created HERE from the current worktree files (the r4
#      mutants were read from disk by that leg's script; I do not inherit its sources).
#   3. Both overlay sides go through `cygpath -m` (forward slashes) - Go rejects a bare
#      backslash inside a JSON string ("invalid escape sequence \w").
#   4. The run is the WHOLE cmd/wisp package, not a 12-name roster, so a mutant is only
#      called "green" if nothing in the package reddens; attribution is a set difference
#      against my own clean baseline roster taken in the same session.
#   5. 11 mutants instead of 4; two of them are the orchestrator's ma-hook / ma-detector-gone
#      pair, re-run by me.
# Repo files are never touched: all breaks are -overlay copies outside the repo.
set -u
REPO="/d/work/workspace/projects plans/Wisp"
OUT=/d/tmp/wisp35v4
LOGS="$REPO/.scratch/wisp/probes/35/v4/logs"
mkdir -p "$OUT/mut" "$LOGS"
PROD=cmd/wisp/panel_host_windows.go
FIX=cmd/wisp/panel_transport_35r2_test.go
REPOPROD="$(cygpath -m "$REPO/$PROD")"
REPOFIX="$(cygpath -m "$REPO/$FIX")"
cd "$REPO" || exit 1

# one-line replacement by line number: $1 src, $2 line, $3 new line, $4 dst
mkmut() { mkdir -p "$(dirname "$4")"; awk -v n="$2" -v new="$3" 'NR==n{print new; next}{print}' "$1" > "$4"; }
mkprod() { mkmut "$REPO/$PROD" "$2" "$3" "$OUT/mut/$1-$PROD"; }
mkfix()  { mkmut "$REPO/$FIX"  "$2" "$3" "$OUT/mut/$1-$FIX"; }

# ---- production-hook mutants (one line each, panel_host_windows.go) ----
# 680: if (inside || typeof window.%[1]s !== "function") { return native.call(cw, message); }
mkprod v-native-window   680 '    if (inside || typeof window.%[1]s !== "function") { return native.call(window, message); }'
mkprod v-guard-truthy    680 '    if (inside || !window.%[1]s) { return native.call(cw, message); }'
mkprod v-ma-hook         680 '    if (inside || typeof window.%[1]s !== "function") { return native(message); }'
# 675: if (!cw || cw.__wispForwardInstalled) { return; }
mkprod v-idem-guard-gone 675 '  if (!cw) { return; }'
# 682: try { return window.%[1]s(message); } finally { inside = false; }
mkprod v-finally-stuck   682 '    try { return window.%[1]s(message); } finally { inside = true; }'
mkprod v-return-gone     682 '    try { window.%[1]s(message); } finally { inside = false; }'
# ---- fixture mutants (one line each, panel_transport_35r2_test.go) ----
# 1288: if recv != cw {
mkfix v-detector-gone 1288 '\t\tif false {'
# 1290: panic(&jsPanic{"TypeError: ... - called with " +   (panic removed, counter kept)
mkfix v-no-panic 1290 '\t\t\tbf.capMessage = ("TypeError: Illegal invocation: chrome.webview.postMessage lost its receiver - called with " +'
# 1947 / 1977: the world selector of the two named cases
mkfix v-ma-world-swap 1947 '\tbf, spy, thrown := runShippedHookInOneWorld35r4(t, true)'
mkfix v-mb-world-swap 1977 '\tbf, spy, thrown := runShippedHookInOneWorld35r4(t, false)'
# 1954: M-A's red sentence. 35-v4 judgment 2 shape: relabel it with M-B's sentence-start so
# a reader who attributed reds by MESSAGE TEXT could not tell the two faces apart. Attribution
# must stay with the test name, which is why this mutant is paired with the ma-hook break.
mkfix v-label-ma 1954 '\t\tt.Fatalf("M-B RED (fixture blind to writability): %d native-exit invocation(s) reached the host without chrome.webview as their receiver. A real browser answers that with Illegal invocation and the page'"'"'s letter never leaves the page, so this hook must not be certified as a delivery. %s", '
# 2048: same relabel trick on the door-absent case, using :1980's sentence-start.
mkfix v-label-mb 2048 '\t\tt.Fatalf("M-B RED (wrong death): with no door in the page the forwarding hook must fall back to the native exit; it instead invoked a non-function and the page'"'"'s post threw: %s. %s",'
# 164: jsObject.set's refusal of a write to a non-writable property (= the leg r4's mb-set-gone)
mkfix v-set-gone 164 '\tif _, exists := o.props[name]; exists && false {'

# 2037: bf.initScripts = bf.initScripts[1:]  -> keep the Bind stub instead of the hook
mkfix v-init-keep-bind 2037 '\tbf.initScripts = bf.initScripts[0:1]'

# replicate the leg's own fourth shape: delete the guard's typeof half entirely
mkprod v-guard-gone 680 '    if (inside) { return native.call(cw, message); }'

# ---- overlay files: replace side only, forward slashes ----
mk_overlay() {
  local name="$1"; shift
  { printf '{"Replace":{'; local first=1 pair o r
    for pair in "$@"; do
      o="${pair%%|*}"; r="${pair##*|}"
      [ "$first" -eq 1 ] || printf ','
      printf '"%s":"%s"' "$o" "$r"; first=0
    done
    printf '}}\n'; } > "$OUT/overlay-$name.json"
}
# ★FIX after pass 1: the pair form is "repoPath|mutantPath" and BOTH sides must be Windows
# paths with forward slashes. Pass 1 passed the mutant path alone, so mk_overlay emitted
# {"Replace":{mutant:mutant}} - a no-op that maps a file onto itself - and every pass-1 run
# was a clean build pretending to be a mutant (12 PASS on all thirteen, v-ma-hook included).
# Named as my own instrument's false green in logs/self-catch-noop-overlay.txt.
P() { echo "$REPOPROD|$(cygpath -m "$OUT/mut/$1-$PROD")"; }
F() { echo "$REPOFIX|$(cygpath -m "$OUT/mut/$1-$FIX")"; }
mk_overlay v-native-window    "$(P v-native-window)"
mk_overlay v-guard-truthy     "$(P v-guard-truthy)"
mk_overlay v-ma-hook          "$(P v-ma-hook)"
mk_overlay v-idem-guard-gone  "$(P v-idem-guard-gone)"
mk_overlay v-finally-stuck    "$(P v-finally-stuck)"
mk_overlay v-return-gone      "$(P v-return-gone)"
mk_overlay v-detector-gone    "$(F v-detector-gone)"
mk_overlay v-ma-hook+detector-gone "$(P v-ma-hook)" "$(F v-detector-gone)"
mk_overlay v-ma-hook+no-panic      "$(P v-ma-hook)" "$(F v-no-panic)"
mk_overlay v-ma-world-swap    "$(F v-ma-world-swap)"
mk_overlay v-mb-world-swap    "$(F v-mb-world-swap)"
mk_overlay v-init-keep-bind   "$(F v-init-keep-bind)"
mk_overlay v-guard-gone       "$(P v-guard-gone)"
mk_overlay v-set-gone         "$(F v-set-gone)"
mk_overlay v-guard-gone+label-mb "$(P v-guard-gone)" "$(F v-label-mb)"
mk_overlay v-ma-hook+label-ma    "$(P v-ma-hook)"  "$(F v-label-ma)"

# ---- mut-landing ruler: every mutant must be exactly ONE changed line ----
{ echo "35-v4 mutant landing ($(date '+%F %T%z')), counts vs the CURRENT worktree files"
  rc=0
  for spec in "v-native-window:$PROD:$(P v-native-window)" "v-guard-truthy:$PROD:$(P v-guard-truthy)" \
              "v-ma-hook:$PROD:$(P v-ma-hook)" "v-idem-guard-gone:$PROD:$(P v-idem-guard-gone)" \
              "v-finally-stuck:$PROD:$(P v-finally-stuck)" "v-return-gone:$PROD:$(P v-return-gone)" \
              "v-detector-gone:$FIX:$(F v-detector-gone)" "v-no-panic:$FIX:$(F v-no-panic)" \
              "v-ma-world-swap:$FIX:$(F v-ma-world-swap)" "v-mb-world-swap:$FIX:$(F v-mb-world-swap)" \
              "v-init-keep-bind:$FIX:$(F v-init-keep-bind)" "v-guard-gone:$PROD:$(P v-guard-gone)" \
              "v-set-gone:$FIX:$(F v-set-gone)" "v-label-mb:$FIX:$(F v-label-mb)" "v-label-ma:$FIX:$(F v-label-ma)"; do
    n="${spec%%:*}"; rest="${spec#*:}"; o="${rest%%:*}"; r="${rest##*:}"
    c=$(diff "$o" "$r" | grep -c '^[<>]')
    [ "$c" = 2 ] || rc=1
    echo "mutant $n changed-line-count=$c"
    diff "$o" "$r" | grep '^[<>]' | cat -A | cut -c1-190
  done
  echo "all-one-line=$([ $rc = 0 ] && echo YES || echo NO)"; echo "rc=$rc"; } > "$LOGS/mut-landing.txt" 2>&1

BEFORE_PROD="$(git hash-object "$PROD")"; BEFORE_FIX="$(git hash-object "$FIX")"

# ROSTER: every test in the repo that consumes the shipping wiring (installPanelTransport /
# panelPostMessageForwardInit as a VALUE).尺 = grep below, recorded in logs/roster-coverage.txt.
# The winlive file that also reads the hook string is excluded by its build tag (no real window
# was opened by this leg: no -tags winlive anywhere).
RUN='^(TestPagePostMessageEnvelopeReachesDispatchRawViaTransport|TestInboundRosterGuardRefusesUnregisteredNameOnPageEdge|TestInboundSourceGuardRefusesForeignSourceOnPageEdge|TestInboundRequestIDGuardRefusesMissingIDOnPageEdge|TestPagePostMessageEnvelopeReachesNativeExitViaShapeA3|TestShapeA3SecondPagePostStillDelivers|TestShapeA3ForwardingHookIsIdempotentInOneDocument|TestLegacySubShapeOneHookDiesInAReentryLoop|TestUnforwardedPageEnvelopeDiesAtTheUnboundSlot|TestForwardingHookMustNotLoseTheNativeExitReceiver|TestForwardingHookIsSilentlyUnarmedByANonWritableNativeExit|TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsAbsent)$'
{ echo "35-v4 roster coverage ruler ($(date '+%F %T%z')): which repo files read the shipping wiring at all"
  grep -rn 'installPanelTransport\|panelPostMessageForwardInit' --include=*.go cmd internal | grep -v '^\S*:[0-9]*:\s*//' | grep -v '_live_35v2_windows_test.go'
  echo "--- build tags of the live file (why it is not in any run today)"
  head -1 cmd/wisp/panel_transport_live_35v2_windows_test.go
  echo "rc=0"; } > "$LOGS/roster-coverage.txt" 2>&1

runpkg() { # $1 = name, $2 = overlay file or ""
  local n="$1" ov="$2" L="$LOGS/run-$1.txt"
  { echo "=== run $n ($(date '+%F %T%z')) overlay=${ov:-none} roster=12 transport names ==="
    if [ -n "$ov" ]; then cat "$ov"; fi
    echo "--- build ---"
    if [ -n "$ov" ]; then go test -c -overlay "$ov" -o "$OUT/pkg-$n.exe" ./cmd/wisp/ 2>&1
    else go test -c -o "$OUT/pkg-$n.exe" ./cmd/wisp/ 2>&1; fi
    echo "buildrc=$?"
    echo "--- run roster, CWD=cmd/wisp, no -tags winlive ---"
    ( cd "$REPO/cmd/wisp" && "$OUT/pkg-$n.exe" -test.run "$RUN" -test.v 2>&1 )
    echo "runrc=$?"
  } > "$L" 2>&1
  grep -aoE -- '--- (FAIL|PASS|SKIP): [A-Za-z0-9_]+' "$L" | sed 's/--- \([A-Z]*\): /\1 /' > "$LOGS/names-$n.txt"
  echo "$n buildrc=$(grep -a '^buildrc' "$L" | tr -d '\r') runrc=$(grep -a '^runrc' "$L" | tr -d '\r') PASS=$(grep -c '^PASS ' "$LOGS/names-$n.txt") FAIL=$(grep -c '^FAIL ' "$LOGS/names-$n.txt") SKIP=$(grep -c '^SKIP ' "$LOGS/names-$n.txt")"
}

{ echo "35-v4 run summary ($(date '+%F %T%z'))  attribution = FAIL-name set difference vs my own baseline roster (12 transport names, same 12 the leg r4 ran)"
  echo "targeted roster chosen because a FULL-package run measured >30 min on this box (150/346 tests in ~14 min); the roster is the complete set of tests that read the shipping wiring -尺 in logs/roster-coverage.txt"
  runpkg baseline ""
  BASE="$LOGS/names-baseline.txt"
  for m in v-ma-hook v-ma-hook+detector-gone v-ma-hook+no-panic v-native-window v-guard-truthy v-idem-guard-gone \
           v-finally-stuck v-return-gone v-detector-gone v-ma-world-swap v-mb-world-swap v-init-keep-bind \
           v-guard-gone v-set-gone v-guard-gone+label-mb v-ma-hook+label-ma; do
    echo "--- $m"
    runpkg "$m" "$OUT/overlay-$m.json"
    echo "  newly-FAIL:";       comm -13 <(grep '^FAIL ' "$BASE" | sort -u) <(grep '^FAIL ' "$LOGS/names-$m.txt" | sort -u) | sed 's/^/    /'
    echo "  no-longer-FAIL:";   comm -23 <(grep '^FAIL ' "$BASE" | sort -u) <(grep '^FAIL ' "$LOGS/names-$m.txt" | sort -u) | sed 's/^/    /'
    echo "  vanished-names(should be none; a missing name = that test did not run):"
    comm -13 <(cut -d' ' -f2- "$BASE" | sort -u) <(cut -d' ' -f2- "$LOGS/names-$m.txt" | sort -u) | sed 's/^/    /'
    echo "  red-reasons (verbatim first line of each failure):"
    grep -a -A3 -- '--- FAIL' "$LOGS/run-$m.txt" | grep -aoE '(M-A|M-B|AC#6|COUNTER-SHAPE|calibration|READING WRONG)[^"]{0,110}' | sed 's/^/    /' | sort -u
  done
  echo "rc=0"; } > "$LOGS/mut-summary.txt" 2>&1

{ echo "35-v4 repo untouched after all overlays ($(date '+%F %T%z'))"
  echo "prod before=$BEFORE_PROD after=$(git hash-object "$PROD")"
  echo "fix  before=$BEFORE_FIX  after=$(git hash-object "$FIX")"
  echo "prod-untouched=$([ "$BEFORE_PROD" = "$(git hash-object "$PROD")" ] && echo YES || echo NO)"
  echo "fix-untouched=$([ "$BEFORE_FIX" = "$(git hash-object "$FIX")" ] && echo YES || echo NO)"
  echo "head-blob-prod=$(git rev-parse HEAD:$PROD) head-blob-fix=$(git rev-parse HEAD:$FIX)"
  echo "git-status-cmd-wisp=[$(git status --porcelain -- cmd/wisp/ | tr '\n' ';')]"
  echo "rc=0"; } > "$LOGS/repo-untouched.txt" 2>&1
echo "driver done"
