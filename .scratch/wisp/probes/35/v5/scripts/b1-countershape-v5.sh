#!/bin/bash
# 35-v5 Q1+Q2 counter-shape driver. NOTHING is written into the repo: every break is a
# -overlay copy living outside the tree (/d/tmp/wisp35v5), exactly like the write leg's
# driver, and the restore proof is a per-file git hash-object roster taken before and after
# over the fixture + production + all winlive files (not just the ones I edit).
#
# Runs:
#  R0 pristine        : overlay a byte-identical prod copy        -> must be 13 PASS/0 FAIL.
#                       Catches "the overlay was never applied" only in PAIR with R1.
#  R1 z-syntax        : overlay prod with a deliberate syntax err  -> buildrc MUST be 1.
#                       This is the proof the overlay KEY PATH matches (35-v4 pass 1 and
#                       111-v1 both died on a missing pair of these two controls).
#  R2 my-a1           : typeof half -> `=== "undefined"`   (NOT tried by 35-r5)
#  R3 my-a2           : typeof half -> `window.X == null`  (NOT tried by 35-r5)
#  R4 my-a3-silentdrop: fallback branch -> `return;` (drops the envelope, never throws)
#                       -> the NEW case must go red on its non-throw assertions.
#  R5 weak-shape      : fixture's new case reduced to "no throw = pass", SHIPPED prod -> green
#  R6 weak+truthy     : same weak fixture + the write leg's truthy mutant -> still red
#                       (=> the no-throw clause alone has teeth for the throw face)
#  R7 b-truthy-new    : re-run of the write leg's mutant -> the red roster must be exactly 1.
set -u
REPO="/d/work/workspace/projects plans/Wisp"
OUT=/d/tmp/wisp35v5
LOGS="$REPO/.scratch/wisp/probes/35/v5/logs"
cd "$REPO" || exit 1
rm -rf "$OUT"; mkdir -p "$OUT/mut" "$LOGS"
PROD=cmd/wisp/panel_host_windows.go
FIX=cmd/wisp/panel_transport_35r2_test.go
REPOPROD="$(cygpath -m "$REPO/$PROD")"
REPOFIX="$(cygpath -m "$REPO/$FIX")"
DLLDIR="$(cygpath -m "$REPO/third_party/sherpa-onnx")"
RUN='TestForwarding|TestPagePostMessage|TestShapeA3|TestLegacySubShape|TestUnforwardedPageEnvelope'

# ---------- restore roster: fixture + prod + every winlive file, individually ----------
# ruler = the build CONSTRAINT carries winlive, .go only (grep -rlE '//go:build.*winlive'),
# which is the same ruler that produced 111-r6's roster-winline.txt = 12 files
# (cmd/wisp 7 + internal/ball 5). A bare `git grep -l 'go:build.*winlive'` hits 73 files
# because it also counts prose in tickets/ledgers/probes - that is a word-frequency ruler,
# not a structural one, and it would have silently shrunk this roster to 7.
WINLINE=$(grep -rlE '//go:build.*winlive' --include='*.go' cmd internal | tr '\n' ' ')
roster() { for f in "$FIX" "$PROD" $WINLINE; do echo "$(git hash-object "$f")  $f"; done; }
roster > "$LOGS/b1-hash-before.txt"
echo "winlive files in the roster: $(echo $WINLINE | wc -w)  total roster lines: $(wc -l < "$LOGS/b1-hash-before.txt")"

# ---------- overlay builder (both halves mandatory; BOTH must be cygpath -m) ----------
# POST-MORTEM of this leg's pass 1: I passed the replacement side as a raw MSYS path
# (/d/tmp/...), Go resolved it against the current drive as D:\d\tmp\... and ALL EIGHT
# runs came back buildrc=1 - including the pristine control, which is exactly how the
# pair-of-controls rule caught it (a no-op overlay would have looked like "nothing is red").
# The builder now refuses any side that is not a drive-letter absolute path.
mk_overlay() { local name="$1"; shift
  local body="" first=1 pair o r ow rw bad=0
  for pair in "$@"; do o="${pair%%|*}"; r="${pair##*|}"
    [ -f "$o" ] || { echo "OVERLAY KEY MISSING ON DISK: $o" >&2; bad=1; }
    [ -f "$r" ] || { echo "OVERLAY VALUE MISSING ON DISK: $r" >&2; bad=1; }
    ow="$(cygpath -m "$o")"; rw="$(cygpath -m "$r")"
    case "$rw" in [A-Za-z]:/*) ;; *) echo "OVERLAY VALUE NOT WINDOWS-ABS: $rw" >&2; bad=1;; esac
    [ "$first" -eq 1 ] || body="$body,"
    body="$body\"$ow\":\"$rw\""; first=0
  done
  [ "$bad" -eq 0 ] || { echo "overlay $name REFUSED" | tee -a "$LOGS/b6-overlay-selfcheck.txt" >&2; return 1; }
  printf '{"Replace":{%s}}\n' "$body" > "$OUT/overlay-$name.json"
  echo "overlay $name OK : $body" >> "$LOGS/b6-overlay-selfcheck.txt"; }

MB="$(basename "$PROD")"; FB="$(basename "$FIX")"
cp "$PROD" "$OUT/mut/pristine-$MB"
sed 's/var panelPostMessageForwardInit = fmt/var panelPostMessageForwardInit = fmt(((syntax error/' "$PROD" > "$OUT/mut/broken-$MB"
# my three shapes; the guard's literal is  if (inside || typeof window.%[1]s !== "function") { return native.call(cw, message); }
sed 's/typeof window\.%\[1\]s !== "function"/typeof window.%[1]s === "undefined"/' "$PROD" > "$OUT/mut/a1-typeofundef-$MB"
sed 's/typeof window\.%\[1\]s !== "function"/window.%[1]s == null/'                  "$PROD" > "$OUT/mut/a2-isnull-$MB"
sed 's/{ return native\.call(cw, message); }/{ return; }/'                           "$PROD" > "$OUT/mut/a3-silentdrop-$MB"
sed 's/typeof window\.%\[1\]s !== "function"/!window.%[1]s/'                         "$PROD" > "$OUT/mut/truthy-$MB"
# ---- Q3 behavioural yard: the ORIGINAL should-red mutant of each re-worded case ----
# r8  v-ma-hook    : native.call(cw, message) -> native(message). The mutant that owns
#                    TestForwardingHookMustNotLoseTheNativeExitReceiver, whose red sentence
#                    is AC#8(iii) site 3 (it was rewritten in place inside a t.Fatalf).
# r9  v-set-gone   : jsObject.set stops refusing writes to non-writable props (fixture :164).
#                    The mutant that owns TestForwardingHookIsSilentlyUnarmedByANonWritable
#                    NativeExit, the case whose 4 red sentences changed prefix M-B->M-B1.
# r10 v-typeofgone : the guard's typeof half deleted (35-v2's original M-B) -> the door-ABSENT
#                    case must stay the M-B2 face and the new one the M-B3 face.
sed 's/{ return native\.call(cw, message); }/{ return native(message); }/'  "$PROD" > "$OUT/mut/r8-mahook-$MB"
sed 's/if (inside || typeof window\.%\[1\]s !== "function")/if (inside)/'   "$PROD" > "$OUT/mut/r10-typeofgone-$MB"
sed 's/exists \&\& o\.unwritable\[name\] {/exists \&\& false {/' "$FIX" > "$OUT/mut/r9-setgone-$FB"
# weak-shape fixture: drop every assertion after the "did it throw" check.
# `if bf.capTripped {` occurs 3x in this file, so the range MUST be NR-guarded to the new
# case (found by name, not by a hardcoded number) - an unguarded range deleted 474 lines.
NEWCASE_START=$(grep -n 'func TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsNotCallable' "$FIX" | cut -d: -f1)
echo "new case starts at line $NEWCASE_START"
awk -v s="$NEWCASE_START" 'NR>=s && /if bf\.capTripped \{/{d=1} /t\.Logf\("M-B3 face has teeth/{d=0} !d' \
  "$FIX" > "$OUT/mut/weak-$FB"

for spec in "pristine:$OUT/mut/pristine-$MB" "a1:$OUT/mut/a1-typeofundef-$MB" "a2:$OUT/mut/a2-isnull-$MB" \
            "a3:$OUT/mut/a3-silentdrop-$MB" "truthy:$OUT/mut/truthy-$MB" \
            "r8-mahook:$OUT/mut/r8-mahook-$MB" "r10-typeofgone:$OUT/mut/r10-typeofgone-$MB"; do
  n="${spec%%:*}"; f="${spec#*:}"
  echo "mutant $n changed-line-count=$(diff "$PROD" "$f" | grep -c '^[<>]')"
  diff "$PROD" "$f" | grep '^[<>]' | sed 's/^/     /'
done > "$LOGS/b2-mutant-landing.txt" 2>&1
{ echo "r9-setgone FIXTURE mutant changed-line-count=$(diff "$FIX" "$OUT/mut/r9-setgone-$FB" | grep -c '^[<>]')  [want 2]"
  diff "$FIX" "$OUT/mut/r9-setgone-$FB" | grep '^[<>]' | sed 's/^/     /'
  echo "weak fixture changed-line-count=$(diff "$FIX" "$OUT/mut/weak-$FB" | grep -c '^[<>]') (all deletions: no line added)"
  echo "weak fixture added-line-count=$(diff "$FIX" "$OUT/mut/weak-$FB" | grep -c '^>')"
  echo "M-B3 tokens surviving in the weak fixture: $(grep -c 'M-B3' "$OUT/mut/weak-$FB") (before=$(grep -c 'M-B3' "$FIX"))"
  echo "guard line in each prod mutant:"; for f in "$OUT/mut"/*"$MB"; do echo "-- $f"; grep -n 'inside' "$f" | sed 's/^/   /'; done
  echo "rc=0"; } >> "$LOGS/b2-mutant-landing.txt" 2>&1

mk_overlay r0-pristine    "$REPO/$PROD|$OUT/mut/pristine-$MB"
mk_overlay r1-zsyntax     "$REPO/$PROD|$OUT/mut/broken-$MB"
mk_overlay r2-a1          "$REPO/$PROD|$OUT/mut/a1-typeofundef-$MB"
mk_overlay r3-a2          "$REPO/$PROD|$OUT/mut/a2-isnull-$MB"
mk_overlay r4-a3          "$REPO/$PROD|$OUT/mut/a3-silentdrop-$MB"
mk_overlay r5-weak-shape  "$REPO/$FIX|$OUT/mut/weak-$FB"
mk_overlay r6-weak-truthy "$REPO/$PROD|$OUT/mut/truthy-$MB" "$REPO/$FIX|$OUT/mut/weak-$FB"
mk_overlay r7-truthy-new  "$REPO/$PROD|$OUT/mut/truthy-$MB"
mk_overlay r8-mahook      "$REPO/$PROD|$OUT/mut/r8-mahook-$MB"
mk_overlay r9-setgone     "$REPO/$FIX|$OUT/mut/r9-setgone-$FB"
mk_overlay r10-typeofgone "$REPO/$PROD|$OUT/mut/r10-typeofgone-$MB"

: > "$LOGS/b3-run-summary.txt"
for m in r0-pristine r1-zsyntax r2-a1 r3-a2 r4-a3 r5-weak-shape r6-weak-truthy r7-truthy-new r8-mahook r9-setgone r10-typeofgone; do
  L="$LOGS/run-$m.txt"
  { echo "=== $m ($(date '+%F %T%z')) overlay=$(cat "$OUT/overlay-$m.json") ==="
    ( cd "$REPO/cmd/wisp" \
      && PATH="$DLLDIR:$PATH" go test -c -overlay "$OUT/overlay-$m.json" \
             -o "$OUT/panel-$m.test.exe" ./ ) > "$OUT/build-$m.log" 2>&1
    echo "buildrc=$?"; tail -12 "$OUT/build-$m.log"
    if [ -f "$OUT/panel-$m.test.exe" ]; then
      # This leg's pass 2 died here with "error while loading shared libraries:
      # sherpa-onnx-c-api.dll" / runrc=127 even with the dll dir prepended to PATH (the dir
      # path contains a space, and a space in a mangled PATH entry is enough). The known-good
      # harness for this package is: the three dlls sit BESIDE the test exe AND the CWD is the
      # package dir. Both are done here; PATH is kept as belt-and-braces.
      cp -f "$REPO"/third_party/sherpa-onnx/*.dll "$OUT/" 2>/dev/null
      echo "dlls beside exe: $(ls "$OUT"/*.dll 2>/dev/null | wc -l)"
      ( cd "$REPO/cmd/wisp" && PATH="$DLLDIR:$PATH" "$OUT/panel-$m.test.exe" -test.run "$RUN" -test.v 2>&1 )
      echo "runrc=$?"
    else echo "(no binary - build failed on purpose)"; echo "runrc=n/a"; fi
  } > "$L" 2>&1
  { echo "=== $m ==="
    echo "PASS=$(grep -c -- '--- PASS' "$L") FAIL=$(grep -c -- '--- FAIL' "$L") SKIP=$(grep -c -- '--- SKIP' "$L") $(grep -E '^(buildrc|runrc)' "$L" | tr '\n' ' ')"
    echo "load-failure sentinel (0xc0000135): $(grep -c 'exit status 0xc0000135' "$L")"
    echo "--- red roster (names only) ---"; grep -- '--- FAIL' "$L" | sed 's/^/   /' || echo "   (none)"
    echo "--- red sentences seen, verbatim prefix ---"
    grep -o 'M-B[0-9]* RED\|M-B[0-9]* READING WRONG\|M-B[0-9]* CALIBRATION RED\|M-A RED' "$L" | sort | uniq -c | sed 's/^/   /'
    echo "--- did the NEW case appear at all? ---"
    echo "   newcase-run-lines=$(grep -c 'TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsNotCallable' "$L")"
  } >> "$LOGS/b3-run-summary.txt" 2>&1
done

roster > "$LOGS/b4-hash-after.txt"
{ echo "=== restore proof: before vs after, per file, whole roster ==="
  diff "$LOGS/b1-hash-before.txt" "$LOGS/b4-hash-after.txt" && echo "IDENTICAL: $(wc -l < "$LOGS/b1-hash-before.txt") files, zero byte moved"
  echo "prod blob worktree=$(git hash-object "$PROD") HEAD=$(git rev-parse "HEAD:$PROD")"
  echo "git status on the two files: [$(git status --porcelain -- "$PROD" "$FIX")]"
  echo "rc=0"; } > "$LOGS/b5-restore.txt" 2>&1
cat "$LOGS/b5-restore.txt"; echo; cat "$LOGS/b3-run-summary.txt"
