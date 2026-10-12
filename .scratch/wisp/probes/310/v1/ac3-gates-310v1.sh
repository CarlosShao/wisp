#!/usr/bin/env bash
# 310-v1 AC#3: the four gates, each run myself, each self-landing rc=N.
# Spelling landmine honoured (gate-roster-304.md S1, added by the orchestrator 10-12):
#   the winsec tier is run as `bash scripts/winsec-tests.sh ./internal/winsec/`.
#   `--scope=winsec` fed to winsec-tests.sh is rejected by its own GUARD 1 with rc=2
#   (that belongs to portable-tests.sh). I record that rejection as a separate shot so the
#   rc difference itself is on disk, and I do NOT use it as a gate reading.
set -u
REPO="D:/work/workspace/projects plans/Wisp"
cd "$REPO" || exit 9
D=".scratch/wisp/probes/310/v1/logs"
M="$D/30-ac3-run.txt"
: > "$M"
say() { echo "$@" >> "$M"; }

TOT="packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0"

# ---------- baseline BEFORE, taken now (never an old number) ----------
bash scripts/portable-tests.sh --scope=census > "$D/30-census-before.txt" 2> "$D/30-census-before.err"
say "census_before_rc=$? gotmpdir_line=$(grep -c 'ticket 310 AC#1: GOTMPDIR=' "$D/30-census-before.txt" || true) degraded_line=$(grep -c 'ticket 310 AC#1: DEGRADED' "$D/30-census-before.txt" "$D/30-census-before.err" | awk -F: '{s+=$2} END{print s+0}') err_bytes=$(wc -c < "$D/30-census-before.err")"
grep -E 'packages=[0-9]+ with-zero-compiled-tests' "$D/30-census-before.txt" > "$D/30-census-before-totals.txt" || true
say "census_before_totals=$(cat "$D/30-census-before-totals.txt")"
say "census_before_totals_matches_ticket_baseline=$([ "$(cat "$D/30-census-before-totals.txt")" = "$TOT" ] && echo YES || echo NO)"

# ---------- gate 1 ----------
sh scripts/d22scan.sh > "$D/30-d22scan.txt" 2> "$D/30-d22scan.err"
say "d22scan_rc=$? out_lines=$(wc -l < "$D/30-d22scan.txt") err_bytes=$(wc -c < "$D/30-d22scan.err") gotmpdir_hits=$(grep -c 'GOTMPDIR=' "$D/30-d22scan.txt" || true) degraded_hits=$(grep -c 'DEGRADED' "$D/30-d22scan.txt" "$D/30-d22scan.err" | awk -F: '{s+=$2} END{print s+0}')"
say "d22scan_tail=$(tail -2 "$D/30-d22scan.txt" | tr '\n' '|')"

# ---------- gate 2 ----------
sh scripts/check-path-length-budget.sh --with-self-test > "$D/30-pathlen.txt" 2> "$D/30-pathlen.err"
say "pathlen_rc=$? out_lines=$(wc -l < "$D/30-pathlen.txt") err_bytes=$(wc -c < "$D/30-pathlen.err")"
say "pathlen_verdict=$(grep -m1 'VERDICT' "$D/30-pathlen.txt")"

# ---------- gate 4 (winsec tier, canonical spelling) ----------
bash scripts/winsec-tests.sh ./internal/winsec/ > "$D/30-winsec.txt" 2> "$D/30-winsec.err"
say "winsec_rc=$? err_bytes=$(wc -c < "$D/30-winsec.err")"
say "winsec_four_numbers=$(grep -m1 'four numbers' "$D/30-winsec.txt")"
say "winsec_gotmpdir_line=$(grep -m1 -o 'ticket 310 AC#1: GOTMPDIR=[^ ]*' "$D/30-winsec.txt" || true)"
say "winsec_tempdir_removeall=$(grep -c 'TempDir RemoveAll cleanup' "$D/30-winsec.txt" "$D/30-winsec.err" | awk -F: '{s+=$2} END{print s+0}')  (the companion ruler added to the canonical roster; empty above = 0 lines = self-cleanup never refused on this root shape)"
say "winsec_err_first=$(head -1 "$D/30-winsec.err")"

# ---------- the spelling landmine, recorded as its own shot (NOT a gate reading) ----------
bash scripts/winsec-tests.sh --scope=winsec > "$D/30-winsec-wrongspelling.txt" 2> "$D/30-winsec-wrongspelling.err"
say "wrongspelling_rc=$? (expect 2 = its own GUARD 1 rejects it; the --scope= layer belongs to portable-tests.sh)"
say "wrongspelling_err=$(head -2 "$D/30-winsec-wrongspelling.err" | tr '\n' '|')"

# ---------- baseline AFTER, taken now ----------
bash scripts/portable-tests.sh --scope=census > "$D/30-census-after.txt" 2> "$D/30-census-after.err"
say "census_after_rc=$? gotmpdir_line=$(grep -c 'ticket 310 AC#1: GOTMPDIR=' "$D/30-census-after.txt" || true) err_bytes=$(wc -c < "$D/30-census-after.err")"
grep -E 'packages=[0-9]+ with-zero-compiled-tests' "$D/30-census-after.txt" > "$D/30-census-after-totals.txt" || true
say "census_after_totals=$(cat "$D/30-census-after-totals.txt")"
diff "$D/30-census-before-totals.txt" "$D/30-census-after-totals.txt" > "$D/30-census-totals-diff.txt" 2>&1
say "census_totals_diff_rc=$? diff_bytes=$(wc -c < "$D/30-census-totals-diff.txt") empty above = 0 bytes = the totals row did not move"

# ---------- the roots my own shots created (I count them, I do not delete them) ----------
CDIR=$(cygpath -u "$LOCALAPPDATA")
R1="$CDIR/wisp-test-gotmp"; R2="/d/tmp/wisp-test-gotmp"
for r in "$R1" "$R2"; do
  say "root=$r child_dirs=$(ls -1 "$r" 2>/dev/null | wc -l | tr -d ' ') my_v1_child_dirs=$(ls -1d "$r"/v1-* 2>/dev/null | wc -l | tr -d ' ')"
done
powershell -NoProfile -ExecutionPolicy Bypass -File .scratch/wisp/probes/310/v1/roster-310v1.ps1 -Out "$D/30-roster-end.txt" -Round "ac3-end" > "$D/30-roster-end.txt.meta" 2> "$D/30-roster-end.err"
say "roster_end_rc=$? count=$(grep -o 'count=[0-9]*' "$D/30-roster-end.txt.meta" | head -1) err_bytes=$(wc -c < "$D/30-roster-end.err")"
say "ac3_done=1"
