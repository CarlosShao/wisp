#!/bin/sh
# r3-flip-columns.sh - ticket 171 r3, AC#9-1: does the new cell of
# .scratch/wisp/probes/161/r6/flip-declaration.sh check the CONTENT of the 基线过期 row,
# or only that a row exists?
#
# THE SHAPE UNDER TEST (171-v2's ③-limit, docs/evidence/s1/171-run-legs-accept-v2.md:142)
#   Raise one leg's want_n above its own measurement -> the gate prints one STALE row whose
#   实测 / 差 columns are FALSIFIED on purpose (实测 printed as the baseline, 差 printed as 0)
#   while the row itself, the leg name and the table count stay honest. A cell that greps
#   only "table exists + names the leg + count is 1" reads that as GREEN. That is the hole.
#
# WHAT MUST BE TRUE AFTER THE FIX
#   the same falsified row -> the cell says FAIL (the two columns are checked), and
#   the honest gate (unmutated working copy) -> the same cell still says ok.
#
# WHY THE COPIES OF THE DOOR LIVE UNDER probes/171/r3/
# flip-declaration.sh resolves the repo root by walking up from its own directory and takes
# its log directory from the same place. A copy kept in this leg's own directory therefore
# (a) writes its logs into this leg's logs/, never into another leg's, and (b) needs exactly
# two lines patched to aim at a fixture ruler instead of the tracked one:
#   gate_rel=...  anchor=...
# Both patches are asserted to change bytes; a sed that matches nothing is not a reading
# (the 恒真检 this repo has already refused twice). The tracked door at
# .scratch/wisp/probes/161/r6/flip-declaration.sh is only ever READ here.
#
# HARD EXIT RULES
#   rc=0  every assertion held     rc=1  an assertion did not hold     rc=2 the instrument broke
#
# Usage: sh r3-flip-columns.sh [pre|post]     (default post)
#   pre  = run the unmodified door against the falsified row (this is the "未修码上" reading)
#   post = unmodified door + fixed door against the falsified row, plus the fixed door against
#          the honest working-copy ruler as the positive control
set -eu

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)

d=$here
while [ ! -e "$d/.git" ]; do
    p=$(CDPATH= cd -- "$d/.." && pwd)
    [ "$p" = "$d" ] && { echo "r3-flip-columns.sh: no .git above $here - refusing to answer with a ruler that did not run" >&2; exit 2; }
    d=$p
done
root=$d

gate_rel=.scratch/wisp/probes/154/gate-clauses.sh
flip_rel=.scratch/wisp/probes/161/r6/flip-declaration.sh
MODE="${1:-post}"
case $MODE in pre | post) ;; *) echo "r3-flip-columns.sh: mode must be pre|post, got '$MODE'" >&2; exit 2 ;; esac

# the tree both runs read. Pinned, because two other legs are committing to dev under a run
# like this one (see r3/impl.md 0.2) - an unpinned anchor would make the two doors answer
# about two different trees and the comparison would be worth nothing.
ANCHOR="${ANCHOR:-4f2a777b60770744adabbd92144d261812a025a0}"
git -C "$root" cat-file -e "${ANCHOR}^{commit}" 2>/dev/null || {
    echo "r3-flip-columns.sh: pinned anchor $ANCHOR is not a commit here - no verdict" >&2; exit 2; }

[ -f "$root/$gate_rel" ] || { echo "r3-flip-columns.sh: gate missing - no verdict" >&2; exit 2; }
command -v mktemp >/dev/null 2>&1 || { echo "r3-flip-columns.sh: no mktemp" >&2; exit 2; }

logdir="$here/logs"
fixtures="$here/fixtures"
# patched copies of the door must live INSIDE the repo: flip-declaration.sh resolves both the
# repo root and its own log directory by walking up from its own file, so a copy sitting in
# $TMPDIR would refuse to run ("no .git above") - and a copy kept under this leg's directory
# writes its per-flip logs here, never into probes/161/r6/logs/ where another leg is reading.
door_dir="$here/doors"
mkdir -p "$logdir" "$fixtures" "$door_dir"
work=$(mktemp -d "${TMPDIR:-/tmp}/wisp-171-r3-flip.XXXXXX")

fail=0
say() { printf '%s\n' "$1"; }
ok() { say "  ok    $1"; }
bad() { say "  FAIL  $1"; fail=1; }

# ---------------------------------------------------------------- 1. the falsified ruler
lie="$fixtures/gate-lie-columns.sh"
stale_ln=$(grep -n 'STALE="\$STALE# STALE 腿=' "$root/$gate_rel" | cut -d: -f1 | head -1)
[ -n "$stale_ln" ] || { echo "r3-flip-columns.sh: no STALE row builder in the gate - the AC#8-2 cell has nothing to check, no verdict" >&2; exit 2; }
awk -v ln="$stale_ln" '
    NR == ln {
        print "\t\tSTALE=\"$STALE# STALE 腿=$LEG 声明=$EXPECT 基线=${COUNT}枚 实测=${COUNT}枚 差=0枚 处置=把 want_n 核下来；样本一枚没变好而是被摘掉的，先回答谁摘的"
        next
    }
    { print }
' "$root/$gate_rel" >"$lie"
[ -s "$lie" ] || { echo "r3-flip-columns.sh: could not write the falsified ruler" >&2; exit 2; }
if cmp -s "$root/$gate_rel" "$lie"; then
    echo "r3-trailing-comment.sh: the falsification changed not one byte - that is not a reading" >&2
    exit 2
fi
n_diff=$(diff "$root/$gate_rel" "$lie" | grep -c '^[<>]')
say "falsified ruler: $lie (from $gate_rel:$stale_ln; diff lines=$n_diff, expected 2)"
[ "$n_diff" = 2 ] || { echo "r3-flip-columns.sh: the falsification touched $n_diff lines instead of one pair - no verdict" >&2; exit 2; }
say "   the line it now prints:"
grep -n 'STALE 腿=' "$lie" | sed -e 's/^/     /'

# ---------------------------------------------------------------- 2. the two doors
# The unmodified door is a BYTE SNAPSHOT kept in this directory (flip-prefix-4f2a777b.sh), not
# "whatever HEAD says when the script runs": the moment this leg's fix is committed, HEAD's door
# stops being the unmodified one and the whole pre/post comparison would silently become a
# comparison of a door with itself. The snapshot is created from HEAD only if it is missing,
# and every later run re-uses the same bytes - which is what makes the two readings comparable
# for whoever audits this cell.
flip_pre="$here/flip-prefix-4f2a777b.sh"
if [ ! -f "$flip_pre" ]; then
    git -C "$root" show "HEAD:$flip_rel" >"$flip_pre"
    printf '   (first run: the unmodified door was taken from HEAD into %s)\n' "$flip_pre"
fi
[ -s "$flip_pre" ] || { echo "r3-flip-columns.sh: the unmodified-door snapshot is empty - no verdict" >&2; exit 2; }
if cmp -s "$flip_pre" "$root/$flip_rel"; then
    echo "r3-flip-columns.sh: the unmodified-door snapshot and the working-copy door are byte-identical - the pre/post comparison would be a door compared with itself, no verdict" >&2
    exit 2
fi
same_unmodified=0

patch_door() { # <src> <dst> <gate_rel-to-use> -> rewrites gate_rel and anchor, asserts bytes moved
    awk -v gr="$3" -v an="$ANCHOR" '
        /^gate_rel=/ { print "gate_rel=" gr; next }
        /^anchor=\$\(git -C "\$root" rev-parse HEAD\)/ { print "anchor=" an; next }
        { print }
    ' "$1" >"$2"
    if cmp -s "$1" "$2"; then
        echo "r3-flip-columns.sh: re-pointing the door changed not one byte - the door no longer has the lines this leg expects" >&2
        exit 2
    fi
    grep -qE "^gate_rel=$3$" "$2" || { echo "r3-flip-columns.sh: gate_rel patch did not land in $2" >&2; exit 2; }
    grep -q "^anchor=$ANCHOR$" "$2" || { echo "r3-flip-columns.sh: anchor patch did not land in $2" >&2; exit 2; }
}

door_lie_pre="$door_dir/door-lie-pre.sh"
patch_door "$flip_pre" "$door_lie_pre" .scratch/wisp/probes/171/r3/fixtures/gate-lie-columns.sh

stale_cell() { # <logfile> -> that cell only, from "== (stale)" to "== (restore)"
    awk '/^== \(stale\)/ { ins = 1 } /^== \(restore\)/ { ins = 0 } ins { print }' "$1"
}
# WHY THE VERDICT IS PER-CLAIM AND NOT "did the cell print any FAIL"
# At the pinned anchor the door's baseline cell and its "rc must stay 0" expectation are already
# RED for a reason that is not this ticket: three legs call BAD because two other legs landed new
# call sites under internal/ and cmd/ (see r3/impl.md 0.2 / 2.4). A per-claim verdict is the only
# shape that separates "were the two columns checked" from "the tree moved under the door".
claim() { # <logfile> <label-substring> -> GREEN / RED / NO-CLAIM for that one assertion
    line=$(stale_cell "$1" | grep -F "$2" | head -1 || true)
    [ -n "$line" ] || { printf 'NO-CLAIM'; return; }
    case $line in
    '  ok'*) printf 'GREEN' ;;
    '  FAIL'*) printf 'RED' ;;
    *) printf 'NO-CLAIM' ;;
    esac
}
cell_verdict() { # <logfile> -> GREEN if the cell printed no FAIL, RED if it did, NO-CELL if the cell is absent
    c=$(stale_cell "$1")
    [ -n "$c" ] || { printf 'NO-CELL'; return; }
    if printf '%s\n' "$c" | grep -q '^  FAIL'; then printf 'RED'; else printf 'GREEN'; fi
}
run_door() { # <door-script> <logname>
    log="$logdir/$2.txt"
    rc=0
    (cd "$root" && sh "$1") >"$log" 2>&1 || rc=$?
    case $rc in
    *[!0-9]*) echo "run_door: exit code '$rc' is not a number - the runner broke, no verdict" >&2; exit 2 ;;
    esac
    printf '%s' "$rc"
}

say ""
say "r3-flip-columns.sh: mode=$MODE   pinned anchor=$ANCHOR"
if [ "$same_unmodified" = 1 ]; then
    say "   NOTE  the committed door and the working-copy door are byte-identical right now,"
    say "         which is expected when this runs before the AC#9-1 fix lands."
fi
say ""
say "===== (1) the unmodified door, against the falsified 实测/差 columns ====="
rc1=$(run_door "$door_lie_pre" lie-door-unmodified)
v1=$(cell_verdict "$logdir/lie-door-unmodified.txt")
row1=$(claim "$logdir/lie-door-unmodified.txt" 'its own readable table row')
col1=$(claim "$logdir/lie-door-unmodified.txt" 'the row prints its own numbers')
say "   the cell said:"
stale_cell "$logdir/lie-door-unmodified.txt" | sed -e 's/^/     /'
say "   row-claim=$row1  columns-claim=$col1  whole-cell verdict=$v1  whole-door rc=$rc1"
if [ "$MODE" = pre ]; then
    [ "$row1" = GREEN ] && ok "unmodified door: the falsified row passes for a reading (that is AC#9-1's hole, 171-v2's ③-limit)" || bad "unmodified door: expected the row claim to be fooled, got $row1"
    [ "$col1" = NO-CLAIM ] && ok "unmodified door: it never asserts anything about the 实测/差 columns (NO-CLAIM = no such check exists)" || bad "unmodified door: unexpected columns claim '$col1'"
else
    [ "$row1" = GREEN ] && ok "unmodified door still fooled by the falsified row (the comparison is intact)" || bad "unmodified door: expected GREEN-on-lie, got $row1"
    [ "$col1" = NO-CLAIM ] && ok "unmodified door: still no assertion on the two columns" || bad "unmodified door: unexpected columns claim '$col1'"
fi
say ""

if [ "$MODE" = post ]; then
    say "===== (2) the fixed door, against the same falsified columns ====="
    door_lie_post="$door_dir/door-lie-post.sh"
    patch_door "$root/$flip_rel" "$door_lie_post" .scratch/wisp/probes/171/r3/fixtures/gate-lie-columns.sh
    rc2=$(run_door "$door_lie_post" lie-door-fixed)
    row2=$(claim "$logdir/lie-door-fixed.txt" 'its own readable table row')
    col2=$(claim "$logdir/lie-door-fixed.txt" 'the row prints its own numbers')
    stale_cell "$logdir/lie-door-fixed.txt" | sed -e 's/^/     /'
    say "   row-claim=$row2  columns-claim=$col2  whole-door rc=$rc2"
    [ "$col2" = RED ] && ok "fixed door: the falsified 实测/差 columns are now called out (AC#9-1 has teeth)" || bad "fixed door: expected the columns claim to be RED on the falsified row, got $col2"
    say ""
    say "===== (3) the fixed door, against the HONEST working-copy ruler (positive control) ====="
    door_honest="$door_dir/door-honest-post.sh"
    patch_door "$root/$flip_rel" "$door_honest" "$gate_rel"
    rc3=$(run_door "$door_honest" honest-door-fixed)
    row3=$(claim "$logdir/honest-door-fixed.txt" 'its own readable table row')
    col3=$(claim "$logdir/honest-door-fixed.txt" 'the row prints its own numbers')
    stale_cell "$logdir/honest-door-fixed.txt" | sed -e 's/^/     /'
    say "   row-claim=$row3  columns-claim=$col3  whole-door rc=$rc3"
    [ "$col3" = GREEN ] && ok "fixed door on the honest ruler: the new claim still says ok (it is not a constant red)" || bad "fixed door on the honest ruler: got $col3 - the new check accuses an innocent row"
    [ "$row3" = GREEN ] && ok "fixed door on the honest ruler: the pre-existing row claim is untouched" || bad "fixed door on the honest ruler: the row claim moved to $row3 - the AC#8-2 cell was rewritten, not extended"
fi

say ""
if [ "$fail" = 0 ]; then
    say "r3-flip-columns.sh: GREEN - mode=$MODE."
    exit 0
fi
say "r3-flip-columns.sh: RED - mode=$MODE."
exit 1
