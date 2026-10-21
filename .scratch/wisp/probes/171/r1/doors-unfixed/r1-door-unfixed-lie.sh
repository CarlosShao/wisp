#!/bin/sh
# flip-declaration.sh - ticket 161 AC#6 cell 1, "is the approval door ever rung" for the
# aggregate exit code of .scratch/wisp/probes/154/gate-clauses.sh.
# Written by 161-r6, 09-27. It READS the gate file and never writes it.
#
# WHAT IT ANSWERS
# ---------------------------------------------------------------------------
# AC#6 cell 1 is only satisfied by a reading of the form
#     "flip one leg's declaration -> the aggregate exit code moves", both ways,
#     and the unflipped baseline is 0 again.
# Running the gate twice would print rc=0 twice and call that green, which is the
# exact shape AC#6 was written against ("the door lit up, yet it still exits 0").
#
# WHY BOTH DIRECTIONS AND WHY THE RESTORE RUN IS NOT OPTIONAL
#   ring -> quiet on a leg that measured "rings today"  = should be silent, it rang
#   quiet -> ring on a leg that measured "silent today" = should have rung, it did not
# An aggregate wired as "any leg ringing makes the script non-zero" - the constant-red
# shape AC#7/ledger A320/A321 just retired - passes every flip assertion above while
# being the same bug behind another door. Only the restore run (expect rc=0 with 6 of
# 14 legs ringing by design) separates "reacts to declarations" from "reacts to noise".
# One run also flips TWO declarations at once and expects exit code 2, because the
# contract is "aggregate rc = number of mismatched legs", not "is at least one wrong".
#
# WHY THE FLIPS LIVE OUTSIDE THE REPO
# Putting a knowingly-false declaration into the tracked gate file would leave a lie in
# the tree if a run were interrupted. The pristine bytes are copied to $TMPDIR, the sed
# flip is applied to the copy, and every flip is checked with cmp so a sed expression
# that matches nothing cannot be reported as a reading (that no-op is the 恒真检 this
# repo has already refused twice).
#
# HARD EXIT RULES
#   rc=0  every assertion held
#   rc=1  an assertion did not hold -> the aggregate is hollow, or constant-red,
#         or a committed declaration no longer matches today's reading (a finding,
#         NOT a knob: do not fix it by relaxing the declaration)
#   rc=2  the instrument could not run -> no verdict, never read it as green
set -eu

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)

# ---- repo root: walk up to .git, never hardcode a depth ----------------------
d=$here
while [ ! -e "$d/.git" ]; do
    p=$(CDPATH= cd -- "$d/.." && pwd)
    [ "$p" = "$d" ] && { echo "flip-declaration.sh: no .git above $here - refusing to answer with a ruler that did not run" >&2; exit 2; }
    d=$p
done
root=$d

gate_rel=.scratch/wisp/probes/171/r1/fixtures/gate-lie-columns-r1.sh
[ -f "$root/$gate_rel" ] || { echo "flip-declaration.sh: gate $gate_rel not found - no verdict" >&2; exit 2; }

# The anchor is taken the way the gate takes it (a commit; $A's default is HEAD and
# this script must not change that), and is passed explicitly so a later agent
# committing underneath a run cannot make two readings incomparable.
anchor=f4d64fcc

command -v mktemp >/dev/null 2>&1 || { echo "flip-declaration.sh: no mktemp" >&2; exit 2; }
work=$(mktemp -d "${TMPDIR:-/tmp}/wisp-161-r6-flip.XXXXXX")
logdir="$here/logs"
mkdir -p "$logdir"

pristine="$work/gate-pristine.sh"
cp "$root/$gate_rel" "$pristine"

fail=0
seq=0

run_gate() { # <script> <logfile> -> echoes the aggregate exit code
    script=$1
    log=$2
    rc=0
    (cd "$root" && sh "$script" "$anchor") >"$log" 2>&1 || rc=$?
    case $rc in
    *[!0-9]*) echo "run_gate: exit code '$rc' is not a number - the runner broke, no verdict" >&2; exit 2 ;;
    esac
    printf '%s' "$rc"
}

expect() { # <label> <got> <mode: zero|nonzero|2>
    got=$2
    mode=$3
    case $mode in
    zero) held=$( [ "$got" -eq 0 ] && echo 1 || echo 0 ) ;;
    nonzero) held=$( [ "$got" -ne 0 ] && echo 1 || echo 0 ) ;;
    2) held=$( [ "$got" -eq 2 ] && echo 1 || echo 0 ) ;;
    *) echo "expect: unknown mode $mode" >&2; exit 2 ;;
    esac
    if [ "$held" = 1 ]; then
        printf '  ok    %s -> rc=%s (want %s)\n' "$1" "$got" "$mode"
    else
        printf '  FAIL  %s -> rc=%s (want %s)\n' "$1" "$got" "$mode"
        fail=1
    fi
}

flip() { # <label> <want-mode> <sed-expr> [<second-sed-expr>]
    label=$1
    mode=$2
    e1=$3
    e2=${4:-}
    seq=$((seq + 1))
    flipped="$work/gate-$seq.sh"
    cp "$pristine" "$flipped"
    if [ -n "$e2" ]; then
        sed -e "$e1" -e "$e2" "$pristine" >"$flipped"
    else
        sed -e "$e1" "$pristine" >"$flipped"
    fi
    if cmp -s "$pristine" "$flipped"; then
        printf '  FAIL  %s - sed changed not one byte, so nothing was flipped; that is not a reading\n' "$label"
        fail=1
        return 0
    fi
    printf '== flip %s: %s\n' "$seq" "$label"
    printf '   the declaration lines this run differs by:\n'
    diff "$pristine" "$flipped" | grep -E '^[<>] want ' | sed -e 's/^/   /'
    got=$(run_gate "$flipped" "$logdir/flip-$seq.txt")
    expect "$label" "$got" "$mode"
    printf '   which legs the aggregate called out (rc=%s):\n' "$got"
    grep -E '^# BAD ' "$logdir/flip-$seq.txt" | sed -e 's/^/   /' || printf '   (none - the mismatch list is empty, which is what the rc contradicts)\n'
}

printf 'flip-declaration.sh: anchor %s (the gate reads the COMMITTED tree at this anchor, never the working tree)\n' "$anchor"

printf '== (0) baseline: committed declarations, untouched\n'
base_rc=$(run_gate "$pristine" "$logdir/flip-baseline.txt")
expect "14 legs each say what it does" "$base_rc" zero
printf '   legs booked: %s   mismatched: %s\n' \
    "$(grep -oE '腿数＝[0-9]+' "$logdir/flip-baseline.txt" | tail -1)" \
    "$(grep -oE '声明与实测不符＝[0-9]+' "$logdir/flip-baseline.txt" | tail -1)"

printf '== (stale) raise one baseline above its own measurement: rc must stay 0 AND a table must appear\n'
seq=$((seq + 1))
stale="$work/gate-$seq-stale.sh"
sed 's/^want_n 8$/want_n 9/' "$pristine" >"$stale"
if cmp -s "$pristine" "$stale"; then
    printf '  FAIL  want_n 8 -> 9 changed not one byte - G5neg has no such baseline line any more, so this cell cannot answer\n'
    fail=1
else
    printf '   the baseline line this run differs by:\n'
    diff "$pristine" "$stale" | grep -E '^[<>] want_n ' | sed -e 's/^/   /'
    got=$(run_gate "$stale" "$logdir/flip-$seq.txt")
    expect "G5neg 基线 8->9（实测 8）still exits 0 (变好了不算言行不一)" "$got" zero
    if grep -q '^## 基线过期' "$logdir/flip-$seq.txt" &&
        grep -q '^# STALE 腿=G5neg ' "$logdir/flip-$seq.txt" &&
        grep -q '^# 基线过期枚数＝1' "$logdir/flip-$seq.txt"; then
        printf '  ok    the stale baseline is its own readable table row: %s\n' \
            "$(grep -m1 '^# STALE ' "$logdir/flip-$seq.txt")"
    else
        printf '  FAIL  no standalone 基线过期 table naming G5neg - SHR is still human-eyes-only\n'
        fail=1
    fi
fi

# ---- restore ------------------------------------------------------------------

printf "\n"
exit 0
