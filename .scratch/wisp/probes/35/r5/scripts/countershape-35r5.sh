#!/bin/bash
# 35-r5 AC#8(i) counter-shape self-check. Every run is `go test -c -overlay` from OUTSIDE
# the repo: no tracked file is written by this script (verified at the end, per-file blob
# hashes before vs after). Five runs:
#   Z-compile-probe : overlay a prod copy with a deliberate syntax error -> the BUILD must
#                     fail. This is the honest proof that the overlay's key path matched at
#                     all (a pristine control alone cannot tell "replaced by an identical
#                     file" from "never applied" - 35-v4 pass 1 died exactly there).
#   A-pristine      : overlay a byte-identical copy of prod -> must not change any colour.
#   B-truthy-new    : guard's typeof half -> truthiness test, current yard  -> AC#8(i)'s new
#                     case must be the ONLY red.
#   C-truthy-old    : same mutant, the pre-35-r5 fixture (blob at daf1f5a0) -> all green,
#                     i.e. this leg reproduces 35-v4's tautology reading (12 green) and
#                     proves the new case is the sole carrier of that face.
#   D-typeofgone-new: 35-v2's original M-B mutant (typeof half deleted), current yard -> the
#                     door-ABSENT case goes red and the new not-callable case stays green
#                     => the two faces are separate detectors, not one detector twice.
set -u
REPO="/d/work/workspace/projects plans/Wisp"
OUT="/d/tmp/wisp35r5"
LOGS="$REPO/.scratch/wisp/probes/35/r5/logs"
cd "$REPO" || exit 1
mkdir -p "$OUT/mut" "$LOGS"
PROD=cmd/wisp/panel_host_windows.go
FIX=cmd/wisp/panel_transport_35r2_test.go
REPOPROD="$(cygpath -m "$REPO/$PROD")"
REPOFIX="$(cygpath -m "$REPO/$FIX")"

mk_overlay() { # $1 = name, $2.. = "orig|replacement" pairs (BOTH halves required)
  local name="$1"; shift
  { printf '{"Replace":{'; local first=1 pair o r
    for pair in "$@"; do
      o="${pair%%|*}"; r="${pair##*|}"
      [ "$first" -eq 1 ] || printf ','
      printf '"%s":"%s"' "$o" "$r"; first=0
    done
    printf '}}\n'; } > "$OUT/overlay-$name.json"
}

MB="$(basename "$PROD")"   # panel_host_windows.go
FB="$(basename "$FIX")"     # panel_transport_35r2_test.go
cp "$PROD" "$OUT/mut/pristine-$MB"
sed 's/var panelPostMessageForwardInit = fmt/var panelPostMessageForwardInit = fmt(((syntax error/' "$PROD" > "$OUT/mut/broken-$MB"
sed 's/typeof window\.%\[1\]s !== "function"/!window.%[1]s/' "$PROD" > "$OUT/mut/truthy-$MB"
sed 's/if (inside || typeof window\.%\[1\]s !== "function")/if (inside)/' "$PROD" > "$OUT/mut/typeofgone-$MB"
git show "daf1f5a0:$FIX" > "$OUT/mut/oldfixture-$FB"

P_BROKEN="$(cygpath -m "$OUT/mut/broken-$MB")"
P_PRIST="$(cygpath -m "$OUT/mut/pristine-$MB")"
P_TRUTHY="$(cygpath -m "$OUT/mut/truthy-$MB")"
P_TYGONE="$(cygpath -m "$OUT/mut/typeofgone-$MB")"
F_OLD="$(cygpath -m "$OUT/mut/oldfixture-$FB")"

mk_overlay z-compile-probe "$REPOPROD|$P_BROKEN"
mk_overlay a-pristine      "$REPOPROD|$P_PRIST"
mk_overlay b-truthy-new    "$REPOPROD|$P_TRUTHY"
mk_overlay c-truthy-old    "$REPOPROD|$P_TRUTHY" "$REPOFIX|$F_OLD"
mk_overlay d-typeofgone-new "$REPOPROD|$P_TYGONE"

# ---- mutant landing check: each mutant differs from the tracked file by ONE line -------
{ echo "=== mutant landing ($(date '+%F %T%z')): changed-line-count must be 2 (= one line replaced), oldfixture must differ only by 35-r5's edit ==="
  for spec in "z-compile-probe:$PROD:$OUT/mut/broken-$MB" "a-pristine:$PROD:$OUT/mut/pristine-$MB" \
              "b-truthy:$PROD:$OUT/mut/truthy-$MB" "d-typeofgone:$PROD:$OUT/mut/typeofgone-$MB"; do
    n="${spec%%:*}"; rest="${spec#*:}"; o="${rest%%:*}"; r="${rest##*:}"
    echo "mutant $n changed-line-count=$(diff "$o" "$r" | grep -c '^[<>]')"
    diff "$o" "$r" | grep '^[<>]'
  done
  echo "oldfixture changed-line-count=$(diff "$FIX" "$OUT/mut/oldfixture-$FB" | grep -c '^[<>]') (the AC#8(i) case this leg added, reversed out)"
  echo "guard line in each mutant:"
  for f in "$OUT/mut"/*panel_host_windows.go; do echo "-- $f"; grep -n 'inside ||' "$f" || echo "(no such line)"; done
  echo "rc=0"; } > "$LOGS/mutant-landing.txt" 2>&1

BEFORE_PROD="$(git hash-object "$PROD")"; BEFORE_FIX="$(git hash-object "$FIX")"

RUN='TestForwarding|TestPagePostMessage|TestShapeA3|TestLegacySubShape|TestUnforwardedPageEnvelope'
: > "$LOGS/countershape-summary.txt"
for m in z-compile-probe a-pristine b-truthy-new c-truthy-old d-typeofgone-new; do
  L="$LOGS/run-$m.txt"
  { echo "=== run $m ($(date '+%F %T%z')) overlay=$OUT/overlay-$m.json ==="
    cat "$OUT/overlay-$m.json"
    echo "--- CWD=$REPO/cmd/wisp, sherpa-onnx dlls on PATH (load-time requirement) ---"
    ( cd "$REPO/cmd/wisp" \
      && PATH="/d/work/workspace/projects plans/Wisp/third_party/sherpa-onnx:$PATH" \
         go test -c -overlay "$OUT/overlay-$m.json" -o "$OUT/panel-$m.test.exe" ./ ) \
      > "$OUT/build-$m.log" 2>&1
    echo "buildrc=$?"
    tail -25 "$OUT/build-$m.log"
    if [ -f "$OUT/panel-$m.test.exe" ]; then
      ( cd "$REPO/cmd/wisp" \
        && PATH="/d/work/workspace/projects plans/Wisp/third_party/sherpa-onnx:$PATH" \
           "$OUT/panel-$m.test.exe" -test.run "$RUN" -test.v 2>&1 )
      echo "runrc=$?"
    else
      echo "(no test binary produced - build failed on purpose for z-compile-probe)"
      echo "runrc=n/a"
    fi
  } > "$L" 2>&1
  { echo "=== $m ==="
    grep -- '--- FAIL' "$L" || echo "(no FAIL line)"
    echo "counts: PASS=$(grep -c -- '--- PASS' "$L") FAIL=$(grep -c -- '--- FAIL' "$L") SKIP=$(grep -c -- '--- SKIP' "$L") $(grep -E '^(buildrc|runrc)' "$L" | tr '\n' ' ')"
    echo "red-sentence prefixes seen: $(grep -o 'M-B[0-9]* RED\|M-B[0-9]* READING WRONG\|M-B[0-9]* CALIBRATION RED\|M-A RED' "$L" | sort | uniq -c | tr '\n' ';')"
  } >> "$LOGS/countershape-summary.txt"
done

{ echo "=== repo untouched by all overlays ($(date '+%F %T%z')) ==="
  echo "prod before=$BEFORE_PROD after=$(git hash-object "$PROD") HEAD=$(git rev-parse "HEAD:$PROD")"
  echo "fix  before=$BEFORE_FIX  after=$(git hash-object "$FIX")"
  echo "prod-untouched-by-overlays=$([ "$BEFORE_PROD" = "$(git hash-object "$PROD")" ] && echo YES || echo NO)"
  echo "fix-untouched-by-overlays=$([ "$BEFORE_FIX" = "$(git hash-object "$FIX")" ] && echo YES || echo NO)"
  echo "prod-vs-daf1f5a0 diff line count (expect 0, zero production changes this leg):"
  git diff --numstat daf1f5a0 -- "$PROD"; echo "numstat-rc=$?"
  git status --porcelain -- "$PROD"
  echo "rc=0"; } > "$LOGS/repo-untouched.txt" 2>&1
tail -6 "$LOGS/countershape-summary.txt"
