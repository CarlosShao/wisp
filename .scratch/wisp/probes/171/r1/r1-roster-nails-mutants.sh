#!/bin/sh
# make-mutants.sh - builds the mutants for ticket 171 AC#3's roster nails.
# Run from the repo root with sh. Writes ONLY under probes/171/r1/** (this leg's
# write surface). Each mutant is a byte-for-byte copy of the fixed attrib.sh with
# exactly ONE line changed, and the diff is printed so the mutation is proven on
# disk rather than claimed in prose.
#
#   M1 named-deleted      delete the named 777 must-ring case        -> ring-count + named-missing
#   M2 named-substituted  replace the 777 case with a DIFFERENT name -> named-missing ONLY
#                         (this is the mutant that proves the count nail alone is not
#                          enough: 5 ring cases are still registered, the COUNT is right,
#                          and the coverage the count stood for is gone)
#   M3 silent-deleted     delete the must-stay-silent 169 case       -> silent-count only
#
# Each mutant must keep `sh -n`/`bash -n` at 0 - a mutant that does not parse would
# "go red" for the wrong reason and prove nothing about the assertion.
set -eu

S=.scratch/wisp/probes/161/r5/attrib.sh
D=.scratch/wisp/probes/171/r1

[ -f "$S" ] || { echo "make-mutants.sh: $S not found - no verdict" >&2; exit 2; }

# The three source lines, located by content (never by a hardcoded line number: the
# file is edited by other legs too, and a stale number would mutate the wrong line).
L777=$(grep -n "st_case 1 UNATTRIBUTABLE '.scratch/wisp/probes/777/second-fake-ticket.go'" "$S" | cut -d: -f1)
LSILENT=$(grep -n "st_case 0 ticket-169 '.scratch/wisp/probes/169/a-bench-sample.go'" "$S" | cut -d: -f1)
[ -n "$L777" ] || { echo "make-mutants.sh: the 777 must-ring case is not in $S - no verdict" >&2; exit 2; }
[ -n "$LSILENT" ] || { echo "make-mutants.sh: the 169 must-stay-silent case is not in $S - no verdict" >&2; exit 2; }
# Exactly one hit each, or the "one line changed" claim is false.
[ "$(printf '%s\n' "$L777" | grep -c .)" = 1 ] || { echo "make-mutants.sh: the 777 line matched more than once - no verdict" >&2; exit 2; }
[ "$(printf '%s\n' "$LSILENT" | grep -c .)" = 1 ] || { echo "make-mutants.sh: the 169 line matched more than once - no verdict" >&2; exit 2; }
echo "make-mutants.sh: source lines - 777 must-ring at :$L777, 169 must-stay-silent at :$LSILENT"

# --- M1: delete the named must-ring case ------------------------------------
sed "${L777}d" "$S" >"$D/attrib-mut-named-deleted.sh"

# --- M2: substitute a DIFFERENT name into that slot (count unchanged) -------
sed "${L777}s|.*|    st_case 1 UNATTRIBUTABLE '.scratch/wisp/probes/999/another-fake-ticket.go'|" \
    "$S" >"$D/attrib-mut-named-substituted.sh"

# --- M3: delete the must-stay-silent case ----------------------------------
sed "${LSILENT}d" "$S" >"$D/attrib-mut-silent-deleted.sh"

fail=0
# check_mutant <file> <expected removed lines> <expected added lines>
# Counted as two separate numbers, because a pure deletion is one `<-only` hunk (0
# added lines) while a substitution is one removed + one added. A single "diff lines"
# total cannot tell those apart, and 1 would have to be accepted as either - which is
# exactly the count-blind reading these mutants exist to refuse.
check_mutant() {
    m=$1
    want_del=$2
    want_add=$3
    if sh -n "$D/$m" 2>/dev/null; then
        echo "  ok    sh -n $m rc=0"
    else
        echo "  FAIL  sh -n $m did not return 0 - this mutant would go red for the wrong reason"
        fail=1
    fi
    got_del=$(diff "$S" "$D/$m" | grep -c '^<' || true)
    got_add=$(diff "$S" "$D/$m" | grep -c '^>' || true)
    if [ "$got_del" = "$want_del" ] && [ "$got_add" = "$want_add" ]; then
        echo "  ok    $m differs on $got_del removed / $got_add added line(s), as intended"
        diff "$S" "$D/$m" | sed 's/^/        /'
    else
        echo "  FAIL  $m differs on $got_del removed / $got_add added line(s), expected $want_del/$want_add - no verdict"
        diff "$S" "$D/$m" | sed 's/^/        /'
        fail=1
    fi
}

check_mutant attrib-mut-named-deleted.sh 1 0
check_mutant attrib-mut-named-substituted.sh 1 1
check_mutant attrib-mut-silent-deleted.sh 1 0

# The positive control for this whole generator: an UNMUTATED copy must diff to zero,
# otherwise "differs on one line" would be satisfied by a file that differs anywhere.
cp "$S" "$D/attrib-mut-control-identical.sh"
if [ -z "$(diff "$S" "$D/attrib-mut-control-identical.sh")" ]; then
    echo "  ok    positive control: an unmutated copy diffs to 0 lines against $S"
else
    echo "  FAIL  positive control: the unmutated copy already differs - the diff ruler is broken"
    fail=1
fi

[ "$fail" = 0 ] && { echo "make-mutants.sh: GREEN - three mutants + one control built"; exit 0; }
echo "make-mutants.sh: RED - see the FAIL lines above"
exit 1
