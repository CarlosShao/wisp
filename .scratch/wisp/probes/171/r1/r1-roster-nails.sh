#!/bin/sh
# r1-roster-nails.sh - ticket 171 AC#3, the roster nails in
# .scratch/wisp/probes/161/r5/attrib.sh --self-test, and the mutants that give them teeth.
#
# THE INJURY THIS ASSERTS AGAINST (the reason the dispatch asked for it)
#   Before this leg, attrib.sh's --self-test counted exactly one thing - `cases=$ST_RUN`
#   - and printed it. Not one line of it asserted HOW MANY must-ring cases were
#   registered, and not one named a sample that had to be there. attrib.sh's own :313
#   records the pass that commented a print block out and read `cases=8 failures=0`
#   GREEN: a number that is printed is a number nobody checks. So the fix is not "put
#   the self-test in CI" (that would wire a count-only self-check to an alarm); the fix
#   is to give the self-test a roster it can fail against, HERE, and then wire it.
#
# THE THREE NAILS (each nailed separately, on purpose)
#   (1) must-RING case count   == ST_RING_EXPECT     (5)
#   (2) must-SILENT case count == ST_SILENT_EXPECT   (3)
#   (3) every NAMED must-ring sample is present      (999, 777, not_a_bench_path)
#   One total would let a ring case be traded for a silent case and read identical.
#   (3) alone survives a substitution that leaves the count intact.
#
# WHAT THIS SCRIPT RUNS (all four are the real instrument, never a copy of the logic)
#   GREEN   the fixed ruler as-is                       -> rc=0, cases=10 failures=0
#   M1      the named 777 must-ring line deleted        -> rc=1, names the 777 sample
#   M2      the 777 line replaced by a DIFFERENT sample -> rc=1, count still 5,
#           ONLY the named nail fires (this is the mutant that proves (1) is not enough)
#   M3      the must-stay-silent 169 line deleted       -> rc=1, silent-count is the nail
#   RESTORE the fixed ruler again                       -> rc=0 (restore reads green)
#
# HARD EXIT RULES
#   rc=0  every reading above held
#   rc=1  a reading did not hold -> the nails are decoration, or a mutant went red for
#         the wrong reason
#   rc=2  the instrument could not run -> no verdict, never read as green
#
# Usage: sh r1-roster-nails.sh            (builds its mutants first, then judges them)
set -eu

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)

d=$here
while [ ! -e "$d/.git" ]; do
    p=$(CDPATH= cd -- "$d/.." && pwd)
    [ "$p" = "$d" ] && { echo "r1-roster-nails.sh: no .git above $here - refusing to answer with a ruler that did not run" >&2; exit 2; }
    d=$p
done
root=$d
cd "$root"

S=.scratch/wisp/probes/161/r5/attrib.sh
[ -f "$S" ] || { echo "r1-roster-nails.sh: $S not found - no verdict" >&2; exit 2; }

# logdir = this script's own directory, never an inherited CWD (the 09-27 19:1x
# lesson: a pass that took logdir from CWD left 8 untracked files at the repo root).
logdir="$here/logs"
mkdir -p "$logdir"

# Build the mutants with the named generator; it refuses (rc=2) if a source line it
# locates by content is absent or ambiguous, so a stale mutation cannot be judged.
sh "$here/r1-roster-nails-mutants.sh" >"$logdir/nails-mutants-build.txt" 2>&1 || {
    rc=$?
    echo "r1-roster-nails.sh: the mutant generator failed (rc=$rc) - no verdict" >&2
    tail -20 "$logdir/nails-mutants-build.txt" >&2
    exit 2
}

fail=0
say() { printf '%s\n' "$1"; }
ok() { say "  ok    $1"; }
bad() { say "  FAIL  $1"; fail=1; }

# run_one <label> <script> <logfile> -> echoes the exit code, never dies on != 0
run_one() {
    rc=0
    sh "$2" --self-test >"$logdir/$3" 2>&1 || rc=$?
    case $rc in
    *[!0-9]*) echo "run_one: exit code '$rc' is not a number - the runner broke, no verdict" >&2; exit 2 ;;
    esac
    printf '%s' "$rc"
}

# summary_of <logfile> -> the `cases=N failures=M` line, verbatim
summary_of() { grep -m1 '^attrib.sh: --self-test cases=' "$1" || echo 'NO-SUMMARY-LINE'; }

say "r1-roster-nails.sh: ruler = $S"
say "r1-roster-nails.sh: logs  = $logdir"
say ""

# ---------------------------------------------------------------- GREEN ------
say "===== GREEN: the fixed ruler, unmutated ====="
rc=$(run_one GREEN "$S" nails-green.txt)
say "  rc=$rc   $(summary_of "$logdir/nails-green.txt")"
[ "$rc" = 0 ] && ok "the fixed ruler is GREEN (rc=0)" || bad "the fixed ruler is not green (rc=$rc) - nothing below can be judged"
grep -q 'PASS  want ring=5 silent=3 named=3  got ring=5 silent=3' "$logdir/nails-green.txt" &&
    ok "the roster nail PASS line is present and reads ring=5 silent=3" ||
    bad "the roster nail PASS line is missing or reads a different number"
say ""

# ------------------------------------------------------------------ M1 -------
say "===== M1: the named 777 must-ring case DELETED (the dispatch's mutant) ====="
rc=$(run_one M1 "$here/attrib-mut-named-deleted.sh" nails-mut-named-deleted.txt)
say "  rc=$rc   $(summary_of "$logdir/nails-mut-named-deleted.txt")"
[ "$rc" = 1 ] && ok "M1 goes RED (rc=1)" || bad "M1 must go red, got rc=$rc"
grep -q 'ring-count: registered must-RING cases=4, this file declares 5' "$logdir/nails-mut-named-deleted.txt" &&
    ok "M1 trips the ring-count nail (4 registered, 5 declared)" ||
    bad "M1 did not trip the ring-count nail"
grep -q '^           \.scratch/wisp/probes/777/second-fake-ticket\.go$' "$logdir/nails-mut-named-deleted.txt" &&
    ok "M1 names the missing sample verbatim: .scratch/wisp/probes/777/second-fake-ticket.go" ||
    bad "M1 did not name the missing sample verbatim"
say "  the red rows, verbatim:"
grep -e 'FAIL  roster assertion' -A6 "$logdir/nails-mut-named-deleted.txt" | sed 's/^/     /'
say ""

# ------------------------------------------------------------------ M2 -------
say "===== M2: the 777 case SUBSTITUTED by a different sample (count unchanged) ====="
say "  (this is the mutant that proves the count nail alone is insufficient)"
rc=$(run_one M2 "$here/attrib-mut-named-substituted.sh" nails-mut-named-substituted.txt)
say "  rc=$rc   $(summary_of "$logdir/nails-mut-named-substituted.txt")"
[ "$rc" = 1 ] && ok "M2 goes RED (rc=1)" || bad "M2 must go red, got rc=$rc"
if grep -q 'ring-count:' "$logdir/nails-mut-named-substituted.txt"; then
    bad "M2 tripped the count nail too - then M2 does not isolate the named nail"
else
    ok "M2 leaves the ring COUNT intact (5) - only the named nail can catch it"
fi
grep -q 'named-missing: a NAMED must-ring sample is not among the cases that ran' "$logdir/nails-mut-named-substituted.txt" &&
    ok "M2 trips the named nail" ||
    bad "M2 did not trip the named nail"
grep -q '^           \.scratch/wisp/probes/777/second-fake-ticket\.go$' "$logdir/nails-mut-named-substituted.txt" &&
    ok "M2 names the substituted-out sample verbatim" ||
    bad "M2 did not name the substituted-out sample"
say "  the red rows, verbatim:"
grep -e 'FAIL  roster assertion' -A6 "$logdir/nails-mut-named-substituted.txt" | sed 's/^/     /'
say ""

# ------------------------------------------------------------------ M3 -------
say "===== M3: the must-STAY-SILENT 169 case DELETED ====="
rc=$(run_one M3 "$here/attrib-mut-silent-deleted.sh" nails-mut-silent-deleted.txt)
say "  rc=$rc   $(summary_of "$logdir/nails-mut-silent-deleted.txt")"
[ "$rc" = 1 ] && ok "M3 goes RED (rc=1)" || bad "M3 must go red, got rc=$rc"
grep -q 'silent-count: registered must-STAY-SILENT cases=2, this file declares 3' "$logdir/nails-mut-silent-deleted.txt" &&
    ok "M3 trips the silent-count nail (2 registered, 3 declared)" ||
    bad "M3 did not trip the silent-count nail"
say "  the red rows, verbatim:"
grep -e 'FAIL  roster assertion' -A3 "$logdir/nails-mut-silent-deleted.txt" | sed 's/^/     /'
say ""

# --------------------------------------------------------------- RESTORE -----
say "===== RESTORE: the fixed ruler again, after the mutants ====="
rc=$(run_one RESTORE "$S" nails-restore.txt)
say "  rc=$rc   $(summary_of "$logdir/nails-restore.txt")"
[ "$rc" = 0 ] && ok "the restore reads GREEN again (rc=0) - the red above was the mutation, not a broken tree" ||
    bad "the restore is not green (rc=$rc)"
say ""

# -------------------------------------------------- negative control ---------
# The last thing that must be true: the GREEN run and the M1 run are the SAME
# instrument on the SAME tree, so their summaries must differ in exactly the way the
# mutation says. If GREEN and M1 ever read identically, the harness is comparing a
# file with itself and every ok above is worthless.
say "===== control: GREEN and M1 must not read the same ====="
if cmp -s "$logdir/nails-green.txt" "$logdir/nails-mut-named-deleted.txt"; then
    bad "the GREEN log and the M1 log are byte-identical - this harness compared a ruler with itself"
else
    ok "the GREEN log and the M1 log differ ($(grep -c . "$logdir/nails-green.txt" || true) vs $(grep -c . "$logdir/nails-mut-named-deleted.txt" || true) lines)"
fi
say ""

if [ "$fail" = 0 ]; then
    say "r1-roster-nails.sh: GREEN - the roster nails hold on the fixed ruler, and all three mutants go red for the nail they were built to trip."
    exit 0
fi
say "r1-roster-nails.sh: RED - at least one reading above did not hold."
exit 1
