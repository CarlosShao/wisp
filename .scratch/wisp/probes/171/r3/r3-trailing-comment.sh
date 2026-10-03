#!/bin/sh
# r3-trailing-comment.sh - ticket 171 r3 (写码位·仪器程), AC#1: the trailing-comment hole,
# measured on BOTH leg families (pair + run) in the same pass, as the ticket's 09-27 19:1x
# scope note demands ("分子口径要一次覆盖两族腿（pair ＋ run），别只修 pair 那一侧").
#
# WHAT IS BEING JUDGED
#   S1  a file whose only mention of a word root sits at the END of a line of real code.
#       AC#1 says: must NOT be counted. Unmodified code counts it (that is the hole).
#   S2  a file with a real opening call and no closing side.
#       AC#1 says: must be counted, before and after. (the positive control)
#   S3  the other face: a real opening call whose closing side exists only as a trailing
#       comment. Unmodified code reads it as PAIRED - a leak standing silent.
#   S4  word roots only in a whole-line comment. pair() drops it today, run() counts it
#       today - the run side of the same hole, which is why this file exists at all.
#
# WHY THE SAMPLES LIVE IN A SYNTHETIC TREE OUTSIDE THE REPO
# Both legs read `git grep "$A"` over internal/**/*.go + cmd/**/*.go of a COMMIT. A sample
# dropped under .scratch/wisp/probes/171/** is invisible twice over. So the samples are
# committed into a throwaway repository under $TMPDIR - never inside this repo: no worktree,
# no checkout, no branch switch, no relaxed ':!' filter. Precedent: 171-r1's r1-pair-shapes.sh
# (:14-:20) and 161-r6's faketree.sh. Deleting files INSIDE that synthetic tree is the only
# deletion this script performs; nothing in the repo is deleted or moved.
#
# WHY TWO RULERS ON EVERY STATE
# "未修码上响不响" is only a judgement if the reading comes from the pre-fix bytes. The
# unmodified ruler is the committed snapshot gate-prefix-f78d3cb.sh (git show HEAD of the
# gate, taken before this leg edited it); the other is the working copy. Same tree, same
# state, same pathspecs - the only difference is the ruler.
#
# HARD EXIT RULES
#   rc=0  every assertion held
#   rc=1  an assertion did not hold -> the hole is still open, or the fix moved a reading
#         it had no business moving (a finding, never a knob to turn)
#   rc=2  the instrument could not run -> no verdict, never read it as green
#
# Usage: sh r3-trailing-comment.sh [pre|post]      (default post)
#   pre  = assert only the unmodified-code expectations (run this BEFORE fixing the ruler)
#   post = assert unmodified-code expectations AND the fixed-ruler expectations
set -eu

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)

d=$here
while [ ! -e "$d/.git" ]; do
    p=$(CDPATH= cd -- "$d/.." && pwd)
    [ "$p" = "$d" ] && { echo "r3-trailing-comment.sh: no .git above $here - refusing to answer with a ruler that did not run" >&2; exit 2; }
    d=$p
done
root=$d

gate_rel=.scratch/wisp/probes/154/gate-clauses.sh
GATE_PRE="${GATE_PRE:-$here/gate-prefix-f78d3cb.sh}"
GATE_POST="${GATE_POST:-$root/$gate_rel}"
MODE="${1:-post}"
case $MODE in pre | post) ;; *) echo "r3-trailing-comment.sh: mode must be pre|post, got '$MODE'" >&2; exit 2 ;; esac

for f in "$GATE_PRE" "$GATE_POST"; do
    [ -f "$f" ] || { echo "r3-trailing-comment.sh: ruler $f not found - no verdict" >&2; exit 2; }
done
command -v git >/dev/null 2>&1 || { echo "r3-trailing-comment.sh: no git" >&2; exit 2; }
command -v mktemp >/dev/null 2>&1 || { echo "r3-trailing-comment.sh: no mktemp" >&2; exit 2; }

# logdir = this script's own directory, never an inherited CWD (the 09-27 19:1x lesson:
# one earlier pass took logdir from CWD and left 8 untracked part*.txt at the repo root).
logdir="$here/logs"
mkdir -p "$logdir"
work=$(mktemp -d "${TMPDIR:-/tmp}/wisp-171-r3.XXXXXX")
tree=$(mktemp -d "${TMPDIR:-/tmp}/wisp-171-r3-tree.XXXXXX")

export GIT_AUTHOR_NAME=wisp171r3 GIT_AUTHOR_EMAIL=171-r3@invalid
export GIT_COMMITTER_NAME=wisp171r3 GIT_COMMITTER_EMAIL=171-r3@invalid
git -C "$tree" init -q .
mkdir -p "$tree/internal/fake171"

fail=0
say() { printf '%s\n' "$1"; }
ok() { say "  ok    $1"; }
bad() { say "  FAIL  $1"; fail=1; }

sec() { # <section-prefix> <logfile> -> that section, up to the next "## "
    awk -v want="$1" '
        /^## / { if (index($0, want) == 1) ins = 1; else if (ins) ins = 0 }
        ins { print }
    ' "$2"
}

metric() { # <section-prefix> <log> <line-glob> -> the integer on that line, or NO-READING
    raw=$(sec "$1" "$2" | grep -m1 "$3" | sed -e 's/[^0-9]*\([0-9][0-9]*\).*/\1/' || true)
    case $raw in
    '' | *[!0-9]*) printf 'NO-READING' ;;
    *) printf '%s' "$raw" ;;
    esac
}
unpaired_of() { metric "$1" "$2" '^# 未成对枚数＝'; }
hits_of() { metric "$1" "$2" '^# 命中行数＝'; }
names_of() { sec "$1" "$2" | grep '^#   UNPAIRED ' | sed -e 's/(开方调用点=[0-9]*).*//' || true; }

run_ruler() { # <pre|post> <state-label> -> echoes the ruler's exit code, log to $logdir
    which=$1
    script=$GATE_PRE
    [ "$which" = post ] && script=$GATE_POST
    log="$logdir/gate-$which-$2.txt"
    rc=0
    (cd "$tree" && sh "$script") >"$log" 2>&1 || rc=$?
    case $rc in
    *[!0-9]*) echo "run_ruler: exit code '$rc' is not a number - the runner broke, no verdict" >&2; exit 2 ;;
    esac
    printf '%s' "$rc"
}

commit_state() { # <label> <fixture-basename...>
    label=$1
    shift
    rm -f "$tree"/internal/fake171/*.go
    for f in "$@"; do
        src="$here/fixtures/$f.go.txt"
        [ -f "$src" ] || { echo "r3-trailing-comment.sh: fixture $src missing - no verdict" >&2; exit 2; }
        cp "$src" "$tree/internal/fake171/$f.go"
    done
    git -C "$tree" add -A -- internal/fake171 >/dev/null
    git -C "$tree" commit -q -m "171-r3: state $label" || { echo "r3-trailing-comment.sh: commit failed in synthetic tree" >&2; exit 2; }
    say "-- state '$label' committed at $(git -C "$tree" rev-parse --short HEAD) in $tree (outside $root)"
}

gofumpt_check() { # the ticket's 派单前置条件 ② applies to throwaway samples too
    if command -v gofumpt >/dev/null 2>&1; then
        out=$(gofumpt -l "$tree/internal/fake171" 2>/dev/null || true)
        if [ -n "$out" ]; then bad "the samples this script generated are not gofumpt-clean: $out"; else ok "generated samples are gofumpt-clean ($(gofumpt --version))"; fi
    else
        say "  --    gofumpt not on PATH; sample formatting left unchecked (named, not guessed)"
    fi
}

say "r3-trailing-comment.sh: mode=$MODE"
say "r3-trailing-comment.sh: unmodified ruler = $GATE_PRE"
say "r3-trailing-comment.sh: working-copy ruler = $GATE_POST"
if cmp -s "$GATE_PRE" "$GATE_POST"; then
    say "   NOTE  the two rulers are byte-identical right now, so the 'fixed' column below is the SAME"
    say "         ruler read a second time. That is expected when this script is run before the AC#1 fix"
    say "         lands; it is not a comparison and is not reported as one."
fi
say "r3-trailing-comment.sh: synthetic tree (deliberately OUTSIDE $root) at $tree"
say ""

G6='## G6 主尺'
G2='## G2 '

# ---------------------------------------------------------------- state S1 -----------
say "===== S1: word root only in a TRAILING comment (AC#1's 'mention is not a call') ====="
commit_state S1 s1_only_trailing_comment
gofumpt_check
rc_pre=$(run_ruler pre S1)
rc_post=$(run_ruler post S1)
P_UN=$(unpaired_of "$G6" "$logdir/gate-pre-S1.txt")
P_HIT=$(hits_of "$G2" "$logdir/gate-pre-S1.txt")
Q_UN=$(unpaired_of "$G6" "$logdir/gate-post-S1.txt")
Q_HIT=$(hits_of "$G2" "$logdir/gate-post-S1.txt")
say "  pair  G6 未成对枚数：unmodified=$P_UN   working-copy=$Q_UN"
say "  run   G2 命中行数 ：unmodified=$P_HIT   working-copy=$Q_HIT"
say "  pair  G6 点名名册（unmodified）："
names_of "$G6" "$logdir/gate-pre-S1.txt" | sed -e 's/^/     /'
[ "$MODE" = post ] && say "  pair  G6 点名名册（working-copy）：" && names_of "$G6" "$logdir/gate-post-S1.txt" | sed -e 's/^/     /'
# unmodified code: counted (this is the reading AC#1 says must exist before the fix)
case $P_UN in 1 | [2-9] | [1-9][0-9]*) ok "unmodified: the trailing-comment file IS counted (未成对=$P_UN 枚)" ;; *) bad "unmodified: expected the comment file to be counted, got $P_UN" ;; esac
case $P_HIT in 1 | [2-9] | [1-9][0-9]*) ok "unmodified: the trailing-comment file IS a run hit (命中行数=$P_HIT)" ;; *) bad "unmodified: expected a run hit from the comment file, got $P_HIT" ;; esac
if [ "$MODE" = post ]; then
    [ "$Q_UN" = 0 ] && ok "fixed: the trailing-comment file is NOT counted (未成对=0)" || bad "fixed: expected 0 未成对, got $Q_UN"
    [ "$Q_HIT" = 0 ] && ok "fixed: the trailing-comment file is NOT a run hit (命中行数=0)" || bad "fixed: expected 0 命中行数, got $Q_HIT"
fi
say ""

# ---------------------------------------------------------------- state S2 -----------
say "===== S2: a REAL opening call (AC#1's 'must be counted', both rulers) ====="
commit_state S2 s2_real_call
rc_pre=$(run_ruler pre S2)
rc_post=$(run_ruler post S2)
P_UN=$(unpaired_of "$G6" "$logdir/gate-pre-S2.txt")
P_HIT=$(hits_of "$G2" "$logdir/gate-pre-S2.txt")
Q_UN=$(unpaired_of "$G6" "$logdir/gate-post-S2.txt")
Q_HIT=$(hits_of "$G2" "$logdir/gate-post-S2.txt")
say "  pair  G6 未成对枚数：unmodified=$P_UN   working-copy=$Q_UN"
say "  run   G2 命中行数 ：unmodified=$P_HIT   working-copy=$Q_HIT"
names_of "$G6" "$logdir/gate-pre-S2.txt" | sed -e 's/^/     /'
if [ "$MODE" = post ]; then names_of "$G6" "$logdir/gate-post-S2.txt" | sed -e 's/^/     /'; fi
[ "$P_UN" = 1 ] && ok "unmodified: the real call is counted (未成对=1)" || bad "unmodified: expected 1 未成对, got $P_UN"
[ "$P_HIT" = 1 ] && ok "unmodified: the real call is a run hit (命中行数=1)" || bad "unmodified: expected 1 命中行数, got $P_HIT"
if [ "$MODE" = post ]; then
    [ "$Q_UN" = 1 ] && ok "fixed: the real call is STILL counted (未成对=1)" || bad "fixed: the real call must still count, got $Q_UN"
    [ "$Q_HIT" = 1 ] && ok "fixed: the real call is STILL a run hit (命中行数=1)" || bad "fixed: the real call must still hit, got $Q_HIT"
fi
say ""

# ---------------------------------------------------------------- state S3 -----------
say "===== S3: real opening call, closing side only in a trailing comment ====="
commit_state S3 s3_fake_close_comment
rc_pre=$(run_ruler pre S3)
rc_post=$(run_ruler post S3)
P_UN=$(unpaired_of "$G6" "$logdir/gate-pre-S3.txt")
P_HIT=$(hits_of "$G2" "$logdir/gate-pre-S3.txt")
Q_UN=$(unpaired_of "$G6" "$logdir/gate-post-S3.txt")
Q_HIT=$(hits_of "$G2" "$logdir/gate-post-S3.txt")
say "  pair  G6 未成对枚数：unmodified=$P_UN   working-copy=$Q_UN"
say "  run   G2 命中行数 ：unmodified=$P_HIT   working-copy=$Q_HIT"
say "  pair  G6 点名名册（unmodified）："; names_of "$G6" "$logdir/gate-pre-S3.txt" | sed -e 's/^/     /'
if [ "$MODE" = post ]; then say "  pair  G6 点名名册（working-copy）："; names_of "$G6" "$logdir/gate-post-S3.txt" | sed -e 's/^/     /'; fi
[ "$P_UN" = 0 ] && ok "unmodified: the leak is SILENT (未成对=0, the comment bought it a close)" || bad "unmodified: expected the silent leak, got $P_UN"
[ "$P_HIT" = 2 ] && ok "unmodified: the run numerator counts 2 lines, one of them a comment" || bad "unmodified: expected 2 命中行数, got $P_HIT"
if [ "$MODE" = post ]; then
    [ "$Q_UN" = 1 ] && ok "fixed: the leak is now NAMED (未成对=1)" || bad "fixed: the leak must be named, got $Q_UN"
    [ "$Q_HIT" = 1 ] && ok "fixed: the run numerator is 1 (the call), not 2" || bad "fixed: expected 1 命中行数, got $Q_HIT"
fi
say ""

# ---------------------------------------------------------------- state S4 -----------
say "===== S4: word roots only in a WHOLE-LINE comment (the run side of the hole) ====="
commit_state S4 s4_full_line_comment
rc_pre=$(run_ruler pre S4)
rc_post=$(run_ruler post S4)
P_UN=$(unpaired_of "$G6" "$logdir/gate-pre-S4.txt")
P_HIT=$(hits_of "$G2" "$logdir/gate-pre-S4.txt")
Q_UN=$(unpaired_of "$G6" "$logdir/gate-post-S4.txt")
Q_HIT=$(hits_of "$G2" "$logdir/gate-post-S4.txt")
say "  pair  G6 未成对枚数：unmodified=$P_UN   working-copy=$Q_UN"
say "  run   G2 命中行数 ：unmodified=$P_HIT   working-copy=$Q_HIT"
[ "$P_UN" = 0 ] && ok "unmodified: pair() already drops the whole-line comment (未成对=0)" || bad "unmodified: expected 0 未成对, got $P_UN"
[ "$P_HIT" = 1 ] && ok "unmodified: run() counts that comment line as a hit (命中行数=1)" || bad "unmodified: expected 1 命中行数, got $P_HIT"
if [ "$MODE" = post ]; then
    [ "$Q_UN" = 0 ] && ok "fixed: pair() still drops it (未成对=0)" || bad "fixed: pair() changed shape here, got $Q_UN"
    [ "$Q_HIT" = 0 ] && ok "fixed: run() no longer counts the comment (命中行数=0)" || bad "fixed: expected 0 命中行数, got $Q_HIT"
fi
say ""

say "r3-trailing-comment.sh: aggregate exit codes seen (informational; the synthetic tree holds no"
say "  registered samples, so want/want_n noise there is expected and is NOT what AC#1 judges):"
for st in S1 S2 S3 S4; do
    say "   state $st  unmodified rc=$(grep -m1 '^# 聚合退码＝' "$logdir/gate-pre-$st.txt" || echo NO-READING)  working-copy rc=$(grep -m1 '^# 聚合退码＝' "$logdir/gate-post-$st.txt" || echo NO-READING)"
done
say ""
if [ "$fail" = 0 ]; then
    say "r3-trailing-comment.sh: GREEN - mode=$MODE, every assertion above held."
    exit 0
fi
say "r3-trailing-comment.sh: RED - at least one assertion above failed."
exit 1
