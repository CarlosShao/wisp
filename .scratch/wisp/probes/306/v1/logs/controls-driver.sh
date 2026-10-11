#!/usr/bin/env bash
# 306-v1 driver 3: HEAD re-anchor + before-state shots + no-fixture mutation controls.
set -u
REPO="D:/work/workspace/projects plans/Wisp"
TREE=$(cat "$REPO/.scratch/wisp/probes/306/v1/logs/TREE-PATH.txt" | tr -d '\r\n ')
LOG="$REPO/.scratch/wisp/probes/306/v1/logs"
SUM="$LOG/ZZ-controls.txt"
INJ="$TREE/internal/audio/wavinjector.go"
T306="$TREE/internal/audio/wavinjector_extensible_float_306_test.go"
PRIST="$TREE/_306v1-pristine"
HIDE="$TREE/_306v1-hidden"
export GOFLAGS=-mod=mod
cd "$REPO" || exit 9
: > "$SUM"
note() { echo "$@" >> "$SUM"; }

restore() { cp "$PRIST/wavinjector.go" "$INJ"; cp "$PRIST/test.go" "$T306"; }
hidefix() { mv "$T306" "$HIDE/h1.txt"; }
showfix() { [ -f "$HIDE/h1.txt" ] && mv "$HIDE/h1.txt" "$T306"; }
m192() { perl -i -pe 's/if fmtTag == 0xFFFE \{/if fmtTag == 0xFFFD {/' "$INJ"; }
m218() { perl -i -pe 's/case fmtTag == 3 && bits == 32:/case fmtTag == 2 && bits == 32:/' "$INJ"; }
m196b() { perl -i -pe 's/data\[body\+24 : body\+26\]/data[body+26 : body+28]/' "$INJ"; }
m196l() { perl -i -pe 's/data\[body\+24 : body\+26\]/data[body+26 : body+26]/' "$INJ"; }

pkgrun() { # name
  local name="$1" out="$LOG/$1.txt"
  (cd "$TREE" && go test ./internal/audio/ -count=1 -v) > "$out" 2>&1
  local rc=$?
  local f p s okln
  f=$(grep -c '^--- FAIL' "$out"); p=$(grep -c '^--- PASS' "$out"); s=$(grep -c '^--- SKIP' "$out"); okln=$(grep -c '^panic:' "$out")
  note "PKG run $name rc=$rc topFAIL=$f topPASS=$p topSKIP=$s panics=$okln"
  grep -E '^(--- FAIL|--- SKIP|panic:|ok |FAIL)' "$out" | sed "s#^#  [$name] #" >> "$SUM"
}

note "=== driver3 start $(date) ==="
# H1: HEAD drift since the export tree was taken
note "H1_head_now=$(git rev-parse HEAD)"
note "H1_head_used_for_tree=e4740e35ff62d9dc8a9732a20cfedbbf4d501a17"
git log --oneline e4740e35ff62d9dc8a9732a20cfedbbf4d501a17..HEAD | sed 's/^/  [H1 newcommit] /' >> "$SUM"
git diff --name-only e4740e35ff62d9dc8a9732a20cfedbbf4d501a17 HEAD -- internal | sed 's/^/  [H1 internal_diff] /' >> "$SUM"
git diff --name-only e4740e35ff62d9dc8a9732a20cfedbbf4d501a17 HEAD -- internal | grep -c . | sed 's/^/  [H1 internal_diff_count] /' >> "$SUM"
git diff --quiet e4740e35ff62d9dc8a9732a20cfedbbf4d501a17 HEAD -- internal/audio; note "H1 rc_diffquiet_internal_audio_e4740e35_vs_now=$?"

# H2: was TestLiveWasapiSmoke already there at the pre-leg baseline?
git grep -c "func TestLiveWasapiSmoke" 5480434f -- internal/audio | sed 's/^/  [H2 preexisting_skip_owner] /' >> "$SUM"
git grep -l "TestLiveWasapiSmoke" 5480434f -- internal/audio | sed 's/^/  [H2 file] /' >> "$SUM"
git diff --name-only 5480434f HEAD -- internal/audio | sed 's/^/  [H2 interval_files] /' >> "$SUM"

# D1: before state (fixture moved away), pristine production, shot #2
restore; hidefix; pkgrun D1-before-nofixture-pristine-shot2; showfix

# D2..D5: fixture moved away, each mutation (the ticket's "rc=0 toothless" control, re-run at HEAD)
restore; hidefix; m192; pkgrun D2-nofixture-mut192; showfix
restore; hidefix; m218; pkgrun D3-nofixture-mut218; showfix
restore; hidefix; m196b; pkgrun D4-nofixture-mut196boundsshift; showfix
restore; hidefix; m196l; pkgrun D5-nofixture-mut196literal; showfix

# D6: fixture present + mut196boundsshift, FULL package (attribute the red roster for AC#2b)
restore; m196b; pkgrun D6-fixture-mut196boundsshift-fullpkg

# D7: fixture present + mut218, FULL package (attribute for AC#1's :218 side)
restore; m218; pkgrun D7-fixture-mut218-fullpkg

restore
cmp -s "$INJ" "$PRIST/wavinjector.go"; note "D8 rc_final_cmp_inj=$?"
cmp -s "$T306" "$PRIST/test.go"; note "D8 rc_final_cmp_test=$?"
ls "$HIDE" | sed 's/^/  [D8 hide_dir_leftover] /' >> "$SUM"

# D9: honest diff rc for the two HEAD shots taken by driver2
diff "$LOG/B1-roster.txt" "$LOG/B2-roster.txt" > "$LOG/D9-b1-b2-roster-diff.txt" 2>&1; note "D9 rc_diff_B1_B2_rosters=$? bytes=$(wc -c < "$LOG/D9-b1-b2-roster-diff.txt")"
grep -cE '^internal/audio/' "$LOG/g4b-gofmt-worktree-fullrepo.txt" | sed 's/^/  [D10 gofmt_fullrepo_internalaudio_hits] /' >> "$SUM"
grep -c . "$LOG/g4b-gofmt-worktree-fullrepo.txt" | sed 's/^/  [D10 gofmt_fullrepo_lines] /' >> "$SUM"
tail -3 "$LOG/g4b-gofmt-worktree-fullrepo.txt" | sed 's#^#  [D10 tail] #' >> "$SUM"
note "=== driver3 done $(date) ==="
