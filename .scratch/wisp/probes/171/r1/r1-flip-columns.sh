#!/bin/sh
# r1-flip-columns.sh - ticket 171, leg 171-r1, AC#9-1: does the door check the 基线过期 row's
# OWN NUMBERS, or only that a row exists?
#
# THE CLAIM UNDER TEST (ticket 171 AC#9-1, inherited from 171-v2's ③-limit)
#   The new cell of flip-declaration.sh was said to check only that the 基线过期 table EXISTS,
#   names the leg and counts one row - its 实测 / 差 columns were unread. A row that falsifies
#   those two columns therefore still passed. AC#9-1 asks for the columns to be checked too,
#   and states the criterion explicitly: "when a column is WRONG the cell must go RED" -
#   NOT "one more line gets printed".
#
# WHY THIS SCRIPT RUNS ONLY THE (stale) CELL INSTEAD OF THE WHOLE DOOR
# ---------------------------------------------------------------------------
# The door's other ten cells each run the full ruler over a real commit. Measured on this
# host on 2026-10-21: one full ruler pass over the 935-file tree takes 7m20s, so the door's
# eleven passes cost about 80 minutes, and a three-door contrast about four hours. That cost
# buys nothing for THIS question: AC#9-1 is entirely decided by the (stale) cell, which runs
# the ruler exactly ONCE. So this script extracts the door's (stale) cell VERBATIM from the
# door's own bytes and drives it directly, and it asserts that extraction: the cell must
# start at "== (stale)" and end at "== (restore)", and the extracted text must differ between
# the unfixed and the fixed door. A hand-written re-implementation of the cell would be a
# second ruler ("the instrument and the thing it measures must not be two programs"), so the
# extraction is the door's own text, not a paraphrase.
#
# WHAT IS ASSERTED, and how each half can fail
#   (1) UNFIXED door + a falsified row  -> the cell must NOT notice (row-claim GREEN,
#       columns-claim NO-CLAIM). If it DOES notice, the hole was already closed by somebody
#       else and this leg would be re-doing landed work - that is a finding, not a pass.
#   (2) FIXED door + the SAME falsified row -> columns-claim MUST be RED. This is AC#9-1.
#   (3) FIXED door + the HONEST ruler -> columns-claim MUST stay GREEN. Without (3) a cell
#       that always reddens would satisfy (2) while being the constant-red shape A320/A321
#       retired - and this repo has refused that shape twice.
#   (4) The two doors must differ as whole texts AND the fixed one must carry the column
#       assertion; otherwise (1)+(2) is one door read twice.
#
# HARD EXIT RULES
#   rc=0  every assertion held
#   rc=1  an assertion did not hold (a finding, never a knob to turn)
#   rc=2  the instrument could not run - no verdict, never read it as green
#
# Usage: sh r1-flip-columns.sh [unfixed|fixed|both]      (default both)
set -eu

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)

d=$here
while [ ! -e "$d/.git" ]; do
    p=$(CDPATH= cd -- "$d/.." && pwd)
    [ "$p" = "$d" ] && { echo "r1-flip-columns.sh: no .git above $here - no verdict" >&2; exit 2; }
    d=$p
done
root=$d

gate_rel=.scratch/wisp/probes/154/gate-clauses.sh
flip_rel=.scratch/wisp/probes/161/r6/flip-declaration.sh
MODE="${1:-both}"
case $MODE in unfixed | fixed | both) ;; *) echo "r1-flip-columns.sh: mode must be unfixed|fixed|both, got '$MODE'" >&2; exit 2 ;; esac

# The anchor the ruler reads. Pinned: another leg may commit under a run this long, and an
# unpinned anchor would make two readings answer about two different trees. f4d64fcc is
# 171-r3 step 1, a commit this very ticket owns.
ANCHOR="${ANCHOR:-f4d64fcc}"
git -C "$root" cat-file -e "${ANCHOR}^{commit}" 2>/dev/null || {
    echo "r1-flip-columns.sh: pinned anchor $ANCHOR is not a commit here - no verdict" >&2; exit 2; }

[ -f "$root/$gate_rel" ] || { echo "r1-flip-columns.sh: gate $gate_rel missing - no verdict" >&2; exit 2; }
[ -f "$root/$flip_rel" ] || { echo "r1-flip-columns.sh: door $flip_rel missing - no verdict" >&2; exit 2; }
command -v mktemp >/dev/null 2>&1 || { echo "r1-flip-columns.sh: no mktemp" >&2; exit 2; }
command -v awk >/dev/null 2>&1 || { echo "r1-flip-columns.sh: no awk" >&2; exit 2; }

logdir="$here/logs"
fixtures="$here/fixtures"
mkdir -p "$logdir" "$fixtures"
work=$(mktemp -d "${TMPDIR:-/tmp}/wisp-171-r1-flipcols.XXXXXX")

fail=0
say() { printf '%s\n' "$1"; }
ok() { say "  ok    $1"; }
bad() { say "  FAIL  $1"; fail=1; }

# ---------------------------------------------------------- 1. the falsified ruler --------
# Exactly one line changes: the STALE row is printed with 实测=<baseline> and 差=0 instead of
# the measurement actually taken. Leg name, baseline, the row's existence and the table's row
# count all stay honest - which is precisely the shape that defeated the pre-fix cell.
honest_gate="$root/$gate_rel"
lie_gate="$fixtures/gate-lie-columns-r1.sh"
stale_ln=$(grep -n 'STALE="\$STALE# STALE 腿=' "$honest_gate" | cut -d: -f1 | head -1)
[ -n "$stale_ln" ] || { echo "r1-flip-columns.sh: no STALE row builder in the gate - AC#9 has nothing to check, no verdict" >&2; exit 2; }
awk -v ln="$stale_ln" '
    NR == ln {
        print "\t\tSTALE=\"$STALE# STALE 腿=$LEG 声明=$EXPECT 基线=${COUNT}枚 实测=${COUNT}枚 差=0枚 处置=把 want_n 核下来；样本一枚没变好而是被摘掉的，先回答谁摘的"
        next
    }
    { print }
' "$honest_gate" >"$lie_gate"
[ -s "$lie_gate" ] || { echo "r1-flip-columns.sh: could not write the falsified ruler - no verdict" >&2; exit 2; }
if cmp -s "$honest_gate" "$lie_gate"; then
    echo "r1-flip-columns.sh: the falsification changed not one byte - that is not a reading, no verdict" >&2
    exit 2
fi
n_diff=$(diff "$honest_gate" "$lie_gate" | grep -c '^[<>]')
say "falsified ruler: $lie_gate (from $gate_rel:$stale_ln; diff lines=$n_diff, expected 2)"
[ "$n_diff" = 2 ] || { echo "r1-flip-columns.sh: the falsification touched $n_diff lines instead of one pair - no verdict" >&2; exit 2; }
say "  the row it prints:"
grep -o 'STALE 腿=[^"]*' "$lie_gate" | sed -e 's/^/     /'
say ""

# ---------------------------------------------------------- 2. the two doors --------------
# The unmodified door is a BYTE SNAPSHOT of the parent of the commit that landed the fix
# (2a5489d3 = 171-r3 step 2, which carries BOTH AC#1 and AC#9-1). Taking it from HEAD would
# stop being the unmodified door the moment the fix is committed, and the contrast would
# silently become one door compared with itself.
door_pre="$here/flip-unfixed-2a5489d3-parent.sh"
if [ ! -f "$door_pre" ]; then
    git -C "$root" show "2a5489d3^:$flip_rel" >"$door_pre"
    say "   (first run: the unmodified door was taken from 2a5489d3^ into $door_pre)"
fi
[ -s "$door_pre" ] || { echo "r1-flip-columns.sh: the unmodified-door snapshot is empty - no verdict" >&2; exit 2; }
door_post="$root/$flip_rel"
say "unmodified door = $door_pre"
say "working-copy door = $door_post"

# ---------------------------------------------------------- 3. extract the (stale) cell ---
# Verbatim extraction, with the boundaries asserted. A cell that cannot be found is rc=2
# (no verdict), never a silent skip.
extract_cell() { # <door> <outfile>
    awk '/^printf .== \(stale\)/ { ins = 1 } /^printf .== \(restore\)/ { ins = 0 } ins { print }' "$1" >"$2"
    [ -s "$2" ] || { echo "r1-flip-columns.sh: the (stale) cell could not be extracted from $1 - no verdict" >&2; exit 2; }
    # control: the extraction must not be the whole file (that would mean the end marker moved)
    if [ "$(wc -l <"$2" | tr -d ' ')" -ge "$(wc -l <"$1" | tr -d ' ')" ]; then
        echo "r1-flip-columns.sh: the extracted cell is as long as the whole door - the end marker did not match, no verdict" >&2
        exit 2
    fi
}
cell_pre="$work/cell-pre.sh"
cell_post="$work/cell-post.sh"
extract_cell "$door_pre" "$cell_pre"
extract_cell "$door_post" "$cell_post"
say "extracted (stale) cell: unfixed=$(wc -l <"$cell_pre" | tr -d ' ') lines, fixed=$(wc -l <"$cell_post" | tr -d ' ') lines"
say ""

patch_cell() { # <cell> <dst> <gate-path-to-use> -> rewrites gate_rel and anchor, asserts bytes moved
    awk -v gr="$3" -v an="$ANCHOR" '
        /^gate_rel=/ { print "gate_rel=" gr; next }
        /^anchor=\$\(git -C "\$root" rev-parse HEAD\)/ { print "anchor=" an; next }
        { print }
    ' "$1" >"$2"
    if cmp -s "$1" "$2"; then
        echo "r1-flip-columns.sh: re-pointing the cell changed not one byte - it no longer has the lines this leg expects, no verdict" >&2
        exit 2
    fi
    grep -qE "^gate_rel=$3$" "$2" || { echo "r1-flip-columns.sh: gate_rel patch did not land in $2" >&2; exit 2; }
    grep -q "^anchor=$ANCHOR$" "$2" || { echo "r1-flip-columns.sh: anchor patch did not land in $2" >&2; exit 2; }
}

# The cell needs the door's own prelude AND its baseline read: the fixed cell takes the
# "measured" column from `$logdir/flip-baseline.txt`, i.e. from the run the BASELINE cell
# writes. So the baseline cell must be kept - stripping it produces a cell that cannot answer
# and reddens on an honest ruler, which is a red this leg made itself (measured: the first
# version of this script did exactly that, see logs/ and the evidence file).
# What IS stripped: the six `flip` calls (each costs a full ruler pass and none of them
# touches this question) and the restore cell.
splice_baseline_and_stale() { # <door> <dst>
    awk '
        /^# ---- ring -> quiet/ { skip = 1 }
        /^printf .== \(stale\)/ { skip = 0 }
        /^printf .== \(restore\)/ { skip = 1 }
        /^printf ....flip-declaration.sh: GREEN/ { skip = 1 }
        /^printf ....flip-declaration.sh: RED/ { skip = 1 }
        /^if \[ "\$fail" = 0 \]/ { skip = 1 }
        skip { next }
        { print }
    ' "$1" >"$2"
    [ -s "$2" ] || { echo "r1-flip-columns.sh: could not strip the door down to baseline+stale - no verdict" >&2; exit 2; }
    grep -q '== (stale)' "$2" || { echo "r1-flip-columns.sh: the (stale) cell vanished while stripping - no verdict" >&2; exit 2; }
    grep -q '== (0) baseline' "$2" || { echo "r1-flip-columns.sh: the baseline cell vanished while stripping - the stale cell would have no measured column to read, no verdict" >&2; exit 2; }
    if grep -q '^flip "G2 ring->quiet' "$2"; then
        echo "r1-flip-columns.sh: the flip cells are still present after stripping - the extraction is wrong, no verdict" >&2
        exit 2
    fi
}

# A stripped door always ends 0 rather than carrying the real door's verdict, because the
# cells that made it red are gone. That is stated, not hidden: the verdicts this script reads
# are the two CLAIM lines, never the door's exit code.
tail_of_stripped() { printf '\nprintf "\\n"\nexit 0\n'; }

build_door() { # <src door> <dst> <gate-rel>
    stripped="$work/stripped-$(basename "$2")"
    splice_baseline_and_stale "$1" "$stripped"
    { cat "$stripped"; tail_of_stripped; } >"$2"
    patch_cell "$2" "$2.patched" "$3"
    mv "$2.patched" "$2"
    # the door resolves its root by walking up from its own file, so the copy must live in
    # the repo (it does: this leg's directory) and its logdir follows it there.
    chmod +x "$2"
}

# ---------------------------------------------------------- 4. reading the claims ---------
claim() { # <logfile> <label-substring> -> GREEN / RED / NO-CLAIM
    line=$(grep -F "$2" "$1" | head -1 || true)
    [ -n "$line" ] || { printf 'NO-CLAIM'; return; }
    case $line in
    '  ok'*) printf 'GREEN' ;;
    '  FAIL'*) printf 'RED' ;;
    *) printf 'NO-CLAIM' ;;
    esac
}

run_cell_door() { # <door> <logname> -> echoes rc
    # WHY THE DOORS RUN ONE AT A TIME, AND WHY flip-baseline.txt IS CLEARED FIRST.
    # The door derives logdir from its own directory, so all three copies in this leg's
    # directory share one logdir. The stale cell reads the MEASURED column out of
    # $logdir/flip-baseline.txt - the run the baseline cell writes - so a leftover file from
    # the previous door would let the stale cell read somebody else's measurement. Clearing it
    # first makes each door answer about its own run. Sequential execution is load-bearing for
    # the same reason and is stated here rather than left to luck.
    rm -f "$logdir/flip-baseline.txt"
    log="$logdir/$2.txt"
    rc=0
    (cd "$root" && sh "$1") >"$log" 2>&1 || rc=$?
    case $rc in
    *[!0-9]*) echo "run_cell_door: exit code '$rc' is not a number - the runner broke, no verdict" >&2; exit 2 ;;
    esac
    # Keep the baseline reading this door produced, next to its own log, before the next door
    # overwrites it - the evidence file cites both.
    if [ -f "$logdir/flip-baseline.txt" ]; then
        cp "$logdir/flip-baseline.txt" "$logdir/flip-baseline-for-$2.txt"
    fi
    printf '%s' "$rc"
}

ROW_LABEL='the stale baseline is its own readable table row'
COL_LABEL='the row prints its own numbers'

lie_rel=".scratch/wisp/probes/171/r1/fixtures/gate-lie-columns-r1.sh"
honest_rel="$gate_rel"

# ---------------------------------------------------------- (1) unfixed door --------------
if [ "$MODE" != fixed ]; then
    say "===== (1) UNFIXED door + the FALSIFIED row: the cell must NOT notice (this is the hole) ====="
    door1="$here/doors-unfixed/r1-door-unfixed-lie.sh"
    mkdir -p "$here/doors-unfixed"
    build_door "$door_pre" "$door1" "$lie_rel"
    rc=$(run_cell_door "$door1" flip-unfixed-lie)
    rowv=$(claim "$logdir/flip-unfixed-lie.txt" "$ROW_LABEL")
    colv=$(claim "$logdir/flip-unfixed-lie.txt" "$COL_LABEL")
    say "  door rc=$rc   row-claim=$rowv   columns-claim=$colv"
    grep -E '^  (ok|FAIL)|^   (the baseline line|the row it|< want_n|> want_n)' "$logdir/flip-unfixed-lie.txt" | sed -e 's/^/     /' || true
    [ "$rowv" = GREEN ] && ok "(1) the unfixed cell DOES see a row (that half was already checked)" ||
        bad "(1) expected the unfixed cell to see the row, got $rowv"
    [ "$colv" = NO-CLAIM ] && ok "(1) the unfixed cell makes NO claim about 实测/差 <- 171-v2's ③-limit reproduces today" ||
        bad "(1) expected NO-CLAIM for the columns from the unfixed door, got $colv - the hole may already be closed by somebody else"
    say ""
fi

# ---------------------------------------------------------- (2)+(3) fixed door -----------
if [ "$MODE" != unfixed ]; then
    say "===== (2) FIXED door + the SAME falsified row: the cell MUST go RED (this is the criterion) ====="
    door2="$here/doors-fixed/r1-door-fixed-lie.sh"
    mkdir -p "$here/doors-fixed"
    build_door "$door_post" "$door2" "$lie_rel"
    rc=$(run_cell_door "$door2" flip-fixed-lie)
    rowv=$(claim "$logdir/flip-fixed-lie.txt" "$ROW_LABEL")
    colv=$(claim "$logdir/flip-fixed-lie.txt" "$COL_LABEL")
    say "  door rc=$rc   row-claim=$rowv   columns-claim=$colv"
    grep -E '^  (ok|FAIL)|^   (the baseline line|the row it|< want_n|> want_n)' "$logdir/flip-fixed-lie.txt" | sed -e 's/^/     /' || true
    [ "$colv" = RED ] && ok "(2) the fixed cell goes RED on a row whose 实测/差 are falsified <- AC#9-1's criterion holds" ||
        bad "(2) expected RED for the columns from the fixed door, got $colv - the cell is checking the printout, not the numbers"
    say ""

    say "===== (3) FIXED door + the HONEST row: the cell must stay GREEN (anti-constant-red control) ====="
    door3="$here/doors-fixed/r1-door-fixed-honest.sh"
    build_door "$door_post" "$door3" "$honest_rel"
    rc=$(run_cell_door "$door3" flip-fixed-honest)
    rowv=$(claim "$logdir/flip-fixed-honest.txt" "$ROW_LABEL")
    colv=$(claim "$logdir/flip-fixed-honest.txt" "$COL_LABEL")
    say "  door rc=$rc   row-claim=$rowv   columns-claim=$colv"
    grep -E '^  (ok|FAIL)|^   (the baseline line|the row it|< want_n|> want_n)' "$logdir/flip-fixed-honest.txt" | sed -e 's/^/     /' || true
    [ "$colv" = GREEN ] && ok "(3) the fixed cell stays GREEN on an honest row <- it is checking the numbers, not always ringing" ||
        bad "(3) expected GREEN for the columns on the honest row, got $colv - the cell became constant-red, which A320/A321 retired"
    say ""
fi

# ---------------------------------------------------------- (4) door contrast --------------
say "===== (4) control: the two doors must differ, and the fix must be in the difference ====="
cmp -s "$door_pre" "$door_post" && bad "(4) the two doors are byte-identical - there is no contrast to read" ||
    ok "(4) the two doors differ ($(wc -l <"$door_pre" | tr -d ' ') vs $(wc -l <"$door_post" | tr -d ' ') lines)"
np=$(grep -c 'the row prints its own numbers' "$door_pre" || true)
nf=$(grep -c 'the row prints its own numbers' "$door_post" || true)
if [ "$nf" -ge 1 ] && [ "$np" = 0 ]; then
    ok "(4) only the fixed door carries the column assertion (unfixed=$np hits, fixed=$nf hits)"
else
    bad "(4) the column assertion is not exclusive to the fixed door (unfixed=$np, fixed=$nf)"
fi
say ""

say "r1-flip-columns.sh: NOTE the verdicts above are taken PER CLAIM. The stripped door always"
say "  exits 0 because the cells that made the real door red (baseline / flips 1-6 / restore) are"
say "  removed on purpose; this leg judges only the two (stale)-cell assertions."
say ""
if [ "$fail" = 0 ]; then
    say "r1-flip-columns.sh: GREEN - mode=$MODE, every assertion above held."
    exit 0
fi
say "r1-flip-columns.sh: RED - at least one assertion above failed."
exit 1
