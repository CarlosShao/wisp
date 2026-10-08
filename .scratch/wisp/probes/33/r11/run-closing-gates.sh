#!/bin/sh
# 33-r11 closing gates: the name-level difference of the two full-package runs, the
# overlay control, the three mutations, and d22scan. Nothing here edits a tracked
# file: mutations reach the compiler only through -overlay.
cd "D:/work/workspace/projects plans/Wisp" || exit 9
export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"
L=.scratch/wisp/probes/33/r11/logs
R=.scratch/wisp/probes/33/r11/mutations
NAMES='TestLockedSuffixNeverTakesTheLockItself33r11|TestRosteredLockedMethodReturnsWhileCallerHoldsTheLock33r11|TestLockedNamingCensusBitesItsOwnFixture33r11|TestLockedNamingCensusRejectsTheOldName33r11|TestCompletedWithinReportsAParkedCall33r11|TestAC13ColdStartPageOverEndsOnEntryContentNotTheProbe'

# 1) per-name difference between the pristine baseline (RUN A) and the post-change
# run (RUN B). A run that never reached a terminal state (panic: test timed out, or
# zero --- FAIL / zero ok lines) is NOT a red list and says so.
{
  echo "=== 50: full-package name diff, RUN A (pristine HEAD via overlay) vs RUN B (post-change) ==="
  for f in 30-baseline-full-package 31-postchange-full-package; do
    echo "--- $f: rc line, terminal-state markers"
    grep -E "^(baseline_rc|postchange_rc|### )" "$L/$f.txt"
    echo "ran_to_terminal=$(grep -cE '^ok |^FAIL\s+github|^panic: ' "$L/$f.txt") panic=$(grep -c 'panic: test timed out' "$L/$f.txt")"
    grep -E "^--- FAIL" "$L/$f.txt" | sed -E 's/^--- FAIL: ([^ ]+).*/\1/' | sort > "$L/$f.names"
    echo "fail_names=$(wc -l < "$L/$f.names")"
    grep -E "^--- PASS|^--- SKIP" "$L/$f.txt" | wc -l
  done
  echo "--- names red in B and not red in A (this is the number that must be zero) ---"
  comm -13 "$L/30-baseline-full-package.names" "$L/31-postchange-full-package.names"
  echo "comm_rc=$?"
  echo "--- names red in A and not in B (a red this leg removed, if any) ---"
  comm -23 "$L/30-baseline-full-package.names" "$L/31-postchange-full-package.names"
  echo "comm2_rc=$?"
  echo "--- both lists, for the record ---"
  cat "$L/30-baseline-full-package.names" "$L/31-postchange-full-package.names"
  echo "cat_rc=$?"
} > "$L/50-package-name-diff.txt" 2>&1

# 2) overlay landing control, then the three mutations
{
  echo "### CONTROL - the declaration is renamed inside the overlay copy while the call site still asks for the old name, so a tree that really used the copy cannot compile. rc=1 IS the pass condition: it proves -overlay feeds this compiler, hence every GREEN below is a green of the tree the overlay describes."
  echo "### build only (go vet), no -run denominator"
  go vet -overlay=$R/overlay-control.json ./cmd/wisp/
  echo "control_rc=$?"
} > "$L/40-mutation-control-undefined-symbol.txt" 2>&1

for tag in m1 m2 m3; do
  {
    case $tag in
      m1) echo "### M1 - the pre-rename name put back on the self-locking round trip (declaration plus its call site, so it still compiles). Expect the census to complain twice: unrostered, and named as if the caller held a lock its body takes.";  ;;
      m2) echo "### M2 - a brand new Locked-suffixed method that takes no lock at all. Expect ONLY the roster complaint: that is the exemption list's own tooth, independent of the body-shape tooth."; ;;
      m3) echo "### M3 - the rostered honest method made to lock itself. Expect the census red AND the behavioural case red, because the call parks under the m.mu the caller already holds."; ;;
    esac
    echo "### -run denominator: 6 named cases (this leg's 5 plus 33-r10's window-free page-over case, which must stay green in every mutation)"
    go test -count=1 -v -timeout 240s -overlay=$R/overlay-$tag.json -run "$NAMES" ./cmd/wisp
    echo "mutation_${tag}_rc=$?"
  } > "$L/41-mut-$tag.txt" 2>&1
done

# 3) proof that no tracked file was touched by any of the above
{
  echo "=== 60: worktree cleanliness after all overlay runs (empty listing = nothing to revert) ==="
  git status --porcelain -- cmd/wisp internal tools
  echo "status_rc=$?"
  echo "=== hash-object parity: worktree file vs the committed blob, per file this leg owns ==="
  for f in cmd/wisp/panel_host_windows.go cmd/wisp/panel_pageover_33r10_windows_test.go cmd/wisp/panel_locked_naming_33r11_windows_test.go; do
    wt=$(git hash-object "$f"); hb=$(git rev-parse "HEAD:$f")
    echo "$f worktree=$wt head=$hb equal=$([ "$wt" = "$hb" ] && echo yes || echo NO)"
  done
  echo "hashloop_rc=$?"
} > "$L/60-worktree-clean-and-hash-parity.txt" 2>&1

# 4) the repo's own scanner, started the way this repo prescribes
{
  echo "=== 61: sh scripts/d22scan.sh ==="
  sh scripts/d22scan.sh
  echo "d22scan_rc=$?"
} > "$L/61-d22scan.txt" 2>&1

echo "closing script finished at $(date '+%Y-%m-%d %H:%M:%S %z')" > "$L/69-closing-script-done.txt"
