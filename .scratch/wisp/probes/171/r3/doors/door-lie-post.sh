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

gate_rel=.scratch/wisp/probes/171/r3/fixtures/gate-lie-columns.sh
[ -f "$root/$gate_rel" ] || { echo "flip-declaration.sh: gate $gate_rel not found - no verdict" >&2; exit 2; }

# The anchor is taken the way the gate takes it (a commit; $A's default is HEAD and
# this script must not change that), and is passed explicitly so a later agent
# committing underneath a run cannot make two readings incomparable.
anchor=4f2a777b60770744adabbd92144d261812a025a0

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

# ---- ring -> quiet : legs that measured "rings today" -------------------------
flip "G2 ring->quiet (run-shaped; its roster is non-empty today)" nonzero \
    's/^want G2 ring$/want G2 quiet/'
flip "G5neg ring->quiet (pair-shaped; 5 unpaired by design today)" nonzero \
    's/^want G5neg ring$/want G5neg quiet/'
flip "G7 ring->quiet (new AC#6 leg; 3 unpaired today, reverse direction only)" nonzero \
    's/^want G7 ring rev$/want G7 quiet rev/'

# ---- quiet -> ring : legs that measured "silent today" -----------------------
flip "G1 quiet->ring (run-shaped; zero hits today)" nonzero \
    's/^want G1 quiet$/want G1 ring/'
flip "G7pos quiet->ring (new AC#6 positive control; 0 unpaired today)" nonzero \
    's/^want G7pos quiet rev$/want G7pos ring rev/'

# ---- two at once: the rc is the COUNT of liars, not a flag --------------------
flip "G1 quiet->ring AND G4 quiet->ring together (expect rc to be 2, not 1)" 2 \
    's/^want G1 quiet$/want G1 ring/' 's/^want G4 quiet$/want G4 ring/'

# ---- NEW CELL, added by 171-r2 (ticket 171 AC#8-2). The nine cells above are untouched. ----
# ---- EXTENDED by 171-r3 (ticket 171 AC#9-1): the cell now checks the row's OWN NUMBERS. ----
# SHR (measured BELOW the registered want_n) is allowed to stay out of the exit code - a reading
# that got better is not a broken promise (that is the constant-red shape A320/A321 retired).
# But it may no longer be only an inline note in the aggregate list, which is human-eyes-only:
# the gate must print a standalone 基线过期 table, and this door must check the table exists
# AND says true things - 171-v2's ③-limit (docs/evidence/s1/171-run-legs-accept-v2.md §14) measured
# that a row falsifying its 实测/差 columns still passed, i.e. the cell was checking the printout.
# Pre-fix readings (same anchor, this very sed): row said 实测=9枚 差=0枚 with 基线=9枚 and
# 枚数=1, the cell said ok - that is the hole; 台件与读数＝.scratch/wisp/probes/171/r3/.
#
# WHY THE BASELINE IS LOOKED UP BY LEG INSTEAD OF HARDCODED AS "want_n 8"
# Ticket 171 AC#1 changed what counts as a call site, so a leg's registered baseline is now a
# number this door must READ, not a number it may remember. A hardcoded digit would either stop
# matching (the cell then dies with "changed not one byte") or - worse - start matching some
# OTHER leg's want_n line and flip the wrong leg. Both shapes turn this cell into decoration the
# next time anyone re-baselines, which is the failure mode the whole ticket exists to close.
printf '== (stale) raise one leg baseline above its own measurement: rc must stay 0 AND a table must appear with its own numbers right\n'
seq=$((seq + 1))
stale="$work/gate-$seq-stale.sh"
stale_leg=G5neg
want_ln=$(grep -nE "^want $stale_leg (ring|quiet)" "$pristine" | head -1 | cut -d: -f1 || true)
base_n=''
[ -n "$want_ln" ] && base_n=$(awk -v ln="$((want_ln + 1))" 'NR == ln && /^want_n [0-9]+$/ { sub(/^want_n /, ""); print }' "$pristine") || true
case $base_n in
'' | *[!0-9]*)
    printf '  FAIL  leg %s has no want_n line right after its want line (want line=%s) - this cell cannot answer, which is not a pass\n' "$stale_leg" "${want_ln:-none}"
    fail=1
    base_n=''
    ;;
esac
if [ -n "$base_n" ]; then
    new_n=$((base_n + 1))
    awk -v ln="$((want_ln + 1))" -v v="$new_n" 'NR == ln { print "want_n " v; next } { print }' "$pristine" >"$stale"
    if cmp -s "$pristine" "$stale"; then
        printf '  FAIL  want_n %s -> %s for leg %s changed not one byte - that is not a reading\n' "$base_n" "$new_n" "$stale_leg"
        fail=1
    else
        printf '   the baseline line this run differs by (leg %s, looked up by leg number, not by a remembered digit):\n' "$stale_leg"
        diff "$pristine" "$stale" | grep -E '^[<>] want_n ' | sed -e 's/^/   /'
        got=$(run_gate "$stale" "$logdir/flip-$seq.txt")
        expect "$stale_leg 基线 $base_n->$new_n still exits 0 (变好了不算言行不一)" "$got" zero
        if grep -q '^## 基线过期' "$logdir/flip-$seq.txt" &&
            grep -q "^# STALE 腿=$stale_leg " "$logdir/flip-$seq.txt" &&
            grep -q '^# 基线过期枚数＝1' "$logdir/flip-$seq.txt"; then
            printf '  ok    the stale baseline is its own readable table row: %s\n' \
                "$(grep -m1 '^# STALE ' "$logdir/flip-$seq.txt")"
        else
            printf '  FAIL  no standalone 基线过期 table naming %s - SHR is still human-eyes-only\n' "$stale_leg"
            fail=1
        fi
        # ---- AC#9-1: the two columns the row prints must themselves be checked ----------------
        # 实测 must equal what the SAME ruler measured for this leg in the untouched baseline run
        # (flip-baseline.txt, same anchor, one line up), and 差 must equal 基线 minus 实测. Anything
        # else is a row that exists and lies.
        m0=$(grep -E "^# (ok|SHR|BAD) *腿=$stale_leg " "$logdir/flip-baseline.txt" | tail -1 |
            sed -e 's/.*实测=\([0-9][0-9]*\)枚.*/\1/')
        case $m0 in
        '' | *[!0-9]*)
            printf '  FAIL  the row prints its own numbers - cannot read %s measured from the baseline run, so the columns cannot be checked\n' "$stale_leg"
            fail=1
            ;;
        *)
            want_row="基线=${new_n}枚 实测=${m0}枚 差=$((new_n - m0))枚"
            row=$(grep -m1 "^# STALE 腿=$stale_leg " "$logdir/flip-$seq.txt" || true)
            if [ -n "$row" ] && printf '%s' "$row" | grep -qF "$want_row"; then
                printf '  ok    the row prints its own numbers: %s\n' "$row"
            else
                printf '  FAIL  the row prints its own numbers - row got [%s] want [%s]\n' "$row" "$want_row"
                fail=1
            fi
            ;;
        esac
    fi
fi

# ---- restore ------------------------------------------------------------------
printf '== (restore) pristine bytes again, same anchor\n'
rest_rc=$(run_gate "$pristine" "$logdir/flip-restored.txt")
expect "nothing flipped any more" "$rest_rc" zero
if [ "$rest_rc" = "$base_rc" ]; then
    printf '  ok    restore matches the baseline (rc=%s)\n' "$rest_rc"
else
    printf '  FAIL  restore rc=%s != baseline rc=%s - the two runs are not comparable\n' "$rest_rc" "$base_rc"
    fail=1
fi

printf '\n'
if [ "$fail" = 0 ]; then
    printf 'flip-declaration.sh: GREEN - the aggregate exit code tracks the declarations in both directions and returns to 0 (baseline rc=%s, 14 legs booked).\n' "$base_rc"
    printf 'flip-declaration.sh: this is provably NOT "any leg rings => non-zero": the baseline books 7 ringing legs (G2 G5 G5neg G6 G6neg G7 G7neg) and 7 silent ones (G1 G1b G3 G4 G5pos G6pos G7pos) and still exits 0.\n'
    exit 0
fi
printf 'flip-declaration.sh: RED - at least one assertion above failed. A declaration that stopped matching today is a finding for the next slice, not a knob to turn.\n'
exit 1
