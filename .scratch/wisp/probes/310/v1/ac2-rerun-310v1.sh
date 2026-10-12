#!/usr/bin/env bash
# 310-v1 AC#2 re-run, my own two shots (never copied from the leg's 20-paired-readings.md).
# (a) the step removed -> a death shape leaks ONE dir into %TEMP% top level
# (b) the step in place -> top level gains nothing, residue lands inside the dedicated root
# Rules honoured: out-of-repo export tree only (never mutate the shared worktree), test output
# goes to its OWN file (never into a reading file), rosters are name lists diffed with comm -13,
# nothing is deleted, the GOTMPDIR root stays on the SAME volume as %TEMP% (gate-roster-304.md S1).
set -u
REPO="D:/work/workspace/projects plans/Wisp"
D="$REPO/.scratch/wisp/probes/310/v1/logs"
PSR="$REPO/.scratch/wisp/probes/310/v1/roster-310v1.ps1"
MAIN="$D/20-ac2-run.txt"
: > "$MAIN"
say() { echo "$@" >> "$MAIN"; }

CDIR=$(cygpath -u "$LOCALAPPDATA")
say "localappdata_unix=$CDIR"
say "temp_unix=$(cygpath -u "$TEMP")"

T=$(mktemp -d /d/tmp/310v1-tree.XXXXXX)
say "export_tree=$T"
( cd "$REPO" && git archive HEAD go.mod go.sum internal cmd tools scripts ) | tar -x -C "$T"
say "archive_rc=$?"
GO=$(find "$T" -name '*.go' -type f | wc -l | tr -d ' ')
say "go_files_in_tree=$GO  (proof the tree is NOT empty before any zero reading)"
SIX=$( (cd "$T" && ls -d go.mod go.sum internal cmd tools scripts 2>/dev/null) | wc -l | tr -d ' ')
say "six_components_present=$SIX (must be 6)"
STEP=$(grep -c 'ticket 310 AC#1' "$T/scripts/portable-tests.sh" || true)
say "ticket310_step_lines_in_export_tree=$STEP (HEAD already carries the leg's step; my shots below drive go test directly, so GOTMPDIR is mine to set or unset - see the judgement file for what this does and does not prove)"

roster() { # $1 label $2 out-names-file
  powershell -NoProfile -ExecutionPolicy Bypass -File "$PSR" -Out "$2" -Round "$1" > "$2.meta" 2> "$2.err"
  say "roster $1 rc=$? names=$(wc -l < "$2") meta_count=$(grep -o 'count=[0-9]*' "$2.meta" | head -1) err_bytes=$(wc -c < "$2.err")"
}

CASE='^TestRoundTripLoadMarshalLoad$'
PKG='./internal/config'

# ---------- warm the build cache so the 2s cap lands INSIDE the loop ----------
( cd "$T" && env -u GOTMPDIR go test -count=1 -run "$CASE" "$PKG" ) > "$D/20-warm.txt" 2>&1
say "warm_rc=$? warm_tail=$(grep -c 'ok  ' "$D/20-warm.txt" || true)"

# ---------- (a) step REMOVED, death shape = -timeout panic ----------
roster a-before "$D/20-roster-a2-before.txt"
( cd "$T" && env -u GOTMPDIR go test -count=4000 -timeout 2s -run "$CASE" "$PKG" ) > "$D/20-shot-a2-test.txt" 2>&1
say "a2_shot_rc=$?"
say "a2_evidence_line=$(grep -m1 -o 'panic: test timed out after 2s' "$D/20-shot-a2-test.txt" || true)"
say "a2_run_lines=$(grep -c '^=== RUN' "$D/20-shot-a2-test.txt" || true) a2_pass_lines=$(grep -c -- '--- PASS' "$D/20-shot-a2-test.txt" || true)"
roster a-after "$D/20-roster-a2-after.txt"
comm -13 "$D/20-roster-a2-before.txt" "$D/20-roster-a2-after.txt" > "$D/20-a2-new.txt"
say "a2_top_level_new=$(wc -l < "$D/20-a2-new.txt" | tr -d ' ') names_in_diff=$(cut -c1-120 "$D/20-a2-new.txt" | tr '\n' ' ')"
say "empty above = a count, zero here means the death leaked nothing (it did not, see names)"

# ---------- (b) step IN PLACE (GOTMPDIR -> same-volume dedicated root), same death ----------
ROOT="$CDIR/wisp-test-gotmp/v1-$(date +%Y%m%d-%H%M%S)-$$"
mkdir -p "$ROOT"
say "b_root=$ROOT mkdir_rc=$?"
roster b-before "$D/20-roster-b-before.txt"
( cd "$T" && GOTMPDIR="$ROOT" go test -count=4000 -timeout 2s -run "$CASE" "$PKG" ) > "$D/20-shot-b-test.txt" 2>&1
say "b_shot_rc=$?"
say "b_evidence_line=$(grep -m1 -o 'panic: test timed out after 2s' "$D/20-shot-b-test.txt" || true)"
say "b_run_lines=$(grep -c '^=== RUN' "$D/20-shot-b-test.txt" || true)"
roster b-after "$D/20-roster-b-after.txt"
comm -13 "$D/20-roster-b-before.txt" "$D/20-roster-b-after.txt" > "$D/20-b-new.txt"
say "b_top_level_new=$(wc -l < "$D/20-b-new.txt" | tr -d ' ') empty above = 0 枚 means the OS temp root gained nothing"
say "b_root_entries=$(ls -1 "$ROOT" 2>/dev/null | wc -l | tr -d ' ') b_root_names=$(ls -1 "$ROOT" 2>/dev/null | tr '\n' ' ')"
say "b_root_full_paths:"
ls -1 "$ROOT" 2>/dev/null | while read -r e; do say "  $ROOT/$e"; done

# ---------- name-proven ruler: CreationTime of everything this run created ----------
powershell -NoProfile -ExecutionPolicy Bypass -File "$PSR" -Out "$D/20-roster-final.txt" -Round "final" > "$D/20-roster-final.txt.meta" 2> "$D/20-roster-final.err"
say "final_roster_rc=$? final_count=$(grep -o 'count=[0-9]*' "$D/20-roster-final.txt.meta" | head -1)"
say "run_done=1"
