#!/bin/sh
# r1-two-shots.sh - ticket 171, leg 171-r1, AC#1: the trailing-comment hole, measured
# as the TWO-SHOT CONTRAST the dispatch demands.
#
# THE TWO SHOTS (dispatch wording: 一发"只有同行注释" / 一发"注释＋真违规")
#   T1  the open-side word root lives ONLY in a trailing comment on a line of real code.
#       AC#1 says it must NOT be counted. The unmodified ruler counts it - that IS the hole.
#   T2  a REAL opening call plus a closing side that exists only in a trailing comment.
#       The unmodified ruler reads the pair as CLOSED and stays silent on a real leak;
#       the fixed ruler must NAME it. This shot is the positive control: the thing that
#       is known to be wrong must ring once the ruler can see it.
#
# THE THIRD SHOT (T3, a control this leg added so "silent" cannot be vacuous)
#   T3  a real opening call and NO closing mention at all, not even in a comment.
#       Both rulers must count it. Without T3, "T1 went 1 -> 0" would be indistinguishable
#       from "this ruler stopped seeing anything", which is the 恒真 direction.
#
# WHY THE SAMPLES LIVE IN A SYNTHETIC TREE OUTSIDE THE REPO
# Both leg families read `git grep "$A"` over internal/**/*.go + cmd/**/*.go of a COMMIT;
# a sample dropped under .scratch/wisp/probes/171/** is invisible to them twice over.
# So the samples are committed into a throwaway repository under $TMPDIR (never inside
# this repo: no worktree, no checkout, no branch switch, no relaxed ':!' filter).
# Precedent: probes/171/r3/r3-trailing-comment.sh:16-22, probes/171/r1/r1-pair-shapes.sh:14-20.
#
# WHY TWO RULERS ON EVERY SHOT
# "未修码上响不响" is only a judgement if the pre-fix bytes produce that reading. The
# unmodified ruler is a pinned snapshot taken from `git show 2a5489d3^:` - the parent of
# the commit that landed the AC#1 fix (171-r3 step 2, 2026-10-03). Taking it from HEAD
# would compare the fixed ruler with itself and the whole contrast would be worthless.
#
# HARD EXIT RULES
#   rc=0  every assertion held
#   rc=1  an assertion did not hold -> the hole is open, or the fix moved a reading it
#         had no business moving (a finding, never a knob to turn)
#   rc=2  the instrument could not run -> no verdict, never read it as green
#
# Usage: sh r1-two-shots.sh [unfixed|fixed|both]      (default both)
#   both (default) asserts BOTH rulers, which is what the two-shot contrast means.
#   unfixed / fixed assert one ruler only - useful when only one side can be built.
set -eu

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)

d=$here
while [ ! -e "$d/.git" ]; do
    p=$(CDPATH= cd -- "$d/.." && pwd)
    [ "$p" = "$d" ] && { echo "r1-two-shots.sh: no .git above $here - refusing to answer with a ruler that did not run" >&2; exit 2; }
    d=$p
done
root=$d

gate_rel=.scratch/wisp/probes/154/gate-clauses.sh
# The commit whose PARENT is the last unmodified ruler. 2a5489d3 is 171-r3 step 2
# ("AC#1 numerator covers both leg families"), 2026-10-03 12:43 +0800.
FIX_COMMIT="${FIX_COMMIT:-2a5489d3}"
GATE_UNFIXED="${GATE_UNFIXED:-$here/gate-unfixed-2a5489d3-parent.sh}"
GATE_FIXED="${GATE_FIXED:-$root/$gate_rel}"
MODE="${1:-both}"
case $MODE in unfixed | fixed | both) ;; *) echo "r1-two-shots.sh: mode must be unfixed|fixed|both, got '$MODE'" >&2; exit 2 ;; esac

for f in "$GATE_UNFIXED" "$GATE_FIXED"; do
    [ -f "$f" ] || { echo "r1-two-shots.sh: ruler $f not found - no verdict" >&2; exit 2; }
done
if cmp -s "$GATE_UNFIXED" "$GATE_FIXED"; then
    echo "r1-two-shots.sh: the unfixed snapshot and the working-copy ruler are byte-identical -" >&2
    echo "  the contrast would be a ruler compared with itself, which is not a reading. no verdict." >&2
    exit 2
fi
command -v git >/dev/null 2>&1 || { echo "r1-two-shots.sh: no git" >&2; exit 2; }
command -v mktemp >/dev/null 2>&1 || { echo "r1-two-shots.sh: no mktemp" >&2; exit 2; }

# logdir = this script's own directory, never an inherited CWD. The 09-27 19:1x lesson:
# one earlier pass took logdir from CWD and left 8 untracked part*.txt at the repo root.
logdir="$here/logs"
mkdir -p "$logdir"
work=$(mktemp -d "${TMPDIR:-/tmp}/wisp-171-r1-shots.XXXXXX")
tree=$(mktemp -d "${TMPDIR:-/tmp}/wisp-171-r1-shots-tree.XXXXXX")

export GIT_AUTHOR_NAME=wisp171r1 GIT_AUTHOR_EMAIL=171-r1@invalid
export GIT_COMMITTER_NAME=wisp171r1 GIT_COMMITTER_EMAIL=171-r1@invalid
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

run_ruler() { # <unfixed|fixed> <state-label> -> echoes the ruler's exit code
    which=$1
    script=$GATE_FIXED
    [ "$which" = unfixed ] && script=$GATE_UNFIXED
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
        [ -f "$src" ] || { echo "r1-two-shots.sh: fixture $src missing - no verdict" >&2; exit 2; }
        cp "$src" "$tree/internal/fake171/$f.go"
    done
    git -C "$tree" add -A -- internal/fake171 >/dev/null
    git -C "$tree" commit -q -m "171-r1: state $label" || { echo "r1-two-shots.sh: commit failed in synthetic tree" >&2; exit 2; }
    say "-- state '$label' committed at $(git -C "$tree" rev-parse --short HEAD) in $tree (outside $root)"
}

# gofumpt is installed under $(go env GOPATH)/bin, which is NOT on this shell's PATH -
# attrib.sh:139-142 derives it the same way. Resolving it here rather than writing
# "not on PATH" (171-r3 had to write that, and a named skip is still a skip).
GOFUMPT_BIN="${GOFUMPT_BIN:-}"
if [ -z "$GOFUMPT_BIN" ] && command -v gofumpt >/dev/null 2>&1; then GOFUMPT_BIN=$(command -v gofumpt); fi
if [ -z "$GOFUMPT_BIN" ] && command -v go >/dev/null 2>&1; then
    _gp=$(go env GOPATH 2>/dev/null || true)
    for c in "$_gp/bin/gofumpt.exe" "$_gp/bin/gofumpt"; do
        [ -x "$c" ] && { GOFUMPT_BIN=$c; break; }
    done
fi

gofumpt_check() { # the ticket's 派单前置条件 ② applies to throwaway samples too
    if [ -z "$GOFUMPT_BIN" ]; then
        bad "gofumpt could not be located (not on PATH, not under GOPATH/bin) - the samples' formatting is UNCHECKED, which is not a pass"
        return 0
    fi
    # WHY A COPY WITH LF ENDINGS: the samples are committed in a throwaway tree, and on this
    # host git rewrites them to CRLF on checkout (the "LF will be replaced by CRLF" warning
    # above). gofumpt -l then flags every file for the line endings alone, which is a
    # property of the throwaway tree, not of the sample this leg wrote. So the formatting
    # reading is taken on an LF copy of the same bytes; the CRLF condition is stated, not
    # silently swallowed. What this still proves: the sample's Go source is gofumpt-clean.
    fmtcopy="$work/fmtcheck-$1"
    rm -rf "$fmtcopy"
    mkdir -p "$fmtcopy"
    for g in "$tree"/internal/fake171/*.go; do
        [ -f "$g" ] || continue
        tr -d '\r' <"$g" >"$fmtcopy/$(basename "$g")"
    done
    out=$("$GOFUMPT_BIN" -l "$fmtcopy" 2>/dev/null || true)
    if [ -n "$out" ]; then
        bad "the samples this script generated are not gofumpt-clean: $out"
    else
        ok "generated samples are gofumpt-clean ($("$GOFUMPT_BIN" --version), read on an LF copy: the throwaway tree's CRLF rewriting is a property of the tree, not of the sample)"
    fi
}

say "r1-two-shots.sh: mode=$MODE"
say "r1-two-shots.sh: unfixed ruler = $GATE_UNFIXED (from 2a5489d3^, the parent of the AC#1 fix)"
say "r1-two-shots.sh: fixed   ruler = $GATE_FIXED"
say "r1-two-shots.sh: synthetic tree (deliberately OUTSIDE $root) at $tree"
say ""

G6='## G6 主尺'
G2='## G2 '

# ------------------------------------------------------------------ shot T1 --------
say "===== T1: word roots ONLY in a trailing comment (must NOT be counted) ====="
commit_state T1 t1_only_trailing_comment
gofumpt_check T1
run_ruler unfixed T1 >/dev/null || true
run_ruler fixed T1 >/dev/null || true
P_UN=$(unpaired_of "$G6" "$logdir/gate-unfixed-T1.txt")
P_HIT=$(hits_of "$G2" "$logdir/gate-unfixed-T1.txt")
Q_UN=$(unpaired_of "$G6" "$logdir/gate-fixed-T1.txt")
Q_HIT=$(hits_of "$G2" "$logdir/gate-fixed-T1.txt")
say "  pair G6 未成对枚数：unfixed=$P_UN   fixed=$Q_UN"
say "  run  G2 命中行数 ：unfixed=$P_HIT   fixed=$Q_HIT"
say "  pair G6 点名名册（unfixed）："; names_of "$G6" "$logdir/gate-unfixed-T1.txt" | sed -e 's/^/     /'
say "  pair G6 点名名册（fixed）  ："; names_of "$G6" "$logdir/gate-fixed-T1.txt" | sed -e 's/^/     /'
if [ "$MODE" != fixed ]; then
    case $P_UN in
    1 | [2-9] | [1-9][0-9]*) ok "T1 unfixed: the trailing-comment file IS counted (未成对=$P_UN 枚) <- the hole is real today" ;;
    *) bad "T1 unfixed: expected the comment file to be counted, got $P_UN" ;;
    esac
    case $P_HIT in
    1 | [2-9] | [1-9][0-9]*) ok "T1 unfixed: the trailing comment is a run hit (命中行数=$P_HIT) <- same hole, run leg" ;;
    *) bad "T1 unfixed: expected a run hit from the comment, got $P_HIT" ;;
    esac
fi
if [ "$MODE" != unfixed ]; then
    [ "$Q_UN" = 0 ] && ok "T1 fixed: the trailing-comment file is NOT counted (未成对=0)" || bad "T1 fixed: expected 0 未成对, got $Q_UN"
    [ "$Q_HIT" = 0 ] && ok "T1 fixed: the trailing comment is NOT a run hit (命中行数=0)" || bad "T1 fixed: expected 0 命中行数, got $Q_HIT"
fi
say ""

# ------------------------------------------------------------------ shot T2 --------
say "===== T2: trailing comment PLUS a real violation (the positive control must RING) ====="
commit_state T2 t2_comment_plus_real_violation
gofumpt_check T2
run_ruler unfixed T2 >/dev/null || true
run_ruler fixed T2 >/dev/null || true
P_UN=$(unpaired_of "$G6" "$logdir/gate-unfixed-T2.txt")
P_HIT=$(hits_of "$G2" "$logdir/gate-unfixed-T2.txt")
Q_UN=$(unpaired_of "$G6" "$logdir/gate-fixed-T2.txt")
Q_HIT=$(hits_of "$G2" "$logdir/gate-fixed-T2.txt")
say "  pair G6 未成对枚数：unfixed=$P_UN   fixed=$Q_UN"
say "  run  G2 命中行数 ：unfixed=$P_HIT   fixed=$Q_HIT"
say "  pair G6 点名名册（unfixed）："; names_of "$G6" "$logdir/gate-unfixed-T2.txt" | sed -e 's/^/     /'
say "  pair G6 点名名册（fixed）  ："; names_of "$G6" "$logdir/gate-fixed-T2.txt" | sed -e 's/^/     /'
if [ "$MODE" != fixed ]; then
    [ "$P_UN" = 0 ] && ok "T2 unfixed: the leak is SILENT (未成对=0 - the comment bought it a close) <- the hole, dangerous direction" ||
        bad "T2 unfixed: expected the silent leak (0), got $P_UN"
    case $P_HIT in
    1 | [2-9] | [1-9][0-9]*) ok "T2 unfixed: the run leg already counts the line (命中行数=$P_HIT)" ;;
    *) bad "T2 unfixed: expected a run hit, got $P_HIT" ;;
    esac
fi
if [ "$MODE" != unfixed ]; then
    [ "$Q_UN" = 1 ] && ok "T2 fixed: the leak is NAMED (未成对=1) <- the positive control rang" || bad "T2 fixed: the leak must be named (1), got $Q_UN"
    # The run leg counts LINES, not files: this sample has exactly one line whose CODE
    # carries an open-side word root (the call). Its one close-side mention is a trailing
    # comment, so a fixed ruler must read 1. This is the number 171-r3 measured for the same
    # shape (s3_fake_close_comment: unfixed 2 -> fixed 1).
    [ "$Q_HIT" = 1 ] && ok "T2 fixed: the run numerator is 1 (the call), not 2" || bad "T2 fixed: expected 1 命中行数, got $Q_HIT"
fi
say ""

# ------------------------------------------------------------------ shot T3 --------
say "===== T3 (this leg's control): a real violation and NO close mention at all - both rulers must count it ====="
say "  (without T3, 'T1 went 1 -> 0' would be indistinguishable from 'this ruler stopped seeing anything')"
commit_state T3 t3_real_violation_no_mention
gofumpt_check T3
run_ruler unfixed T3 >/dev/null || true
run_ruler fixed T3 >/dev/null || true
P_UN=$(unpaired_of "$G6" "$logdir/gate-unfixed-T3.txt")
Q_UN=$(unpaired_of "$G6" "$logdir/gate-fixed-T3.txt")
P_HIT=$(hits_of "$G2" "$logdir/gate-unfixed-T3.txt")
Q_HIT=$(hits_of "$G2" "$logdir/gate-fixed-T3.txt")
say "  pair G6 未成对枚数：unfixed=$P_UN   fixed=$Q_UN"
say "  run  G2 命中行数 ：unfixed=$P_HIT   fixed=$Q_HIT"
[ "$P_UN" = 1 ] && ok "T3 unfixed: the bare violation is counted (未成对=1)" || bad "T3 unfixed: expected 1 未成对, got $P_UN"
[ "$Q_UN" = 1 ] && ok "T3 fixed: the bare violation is STILL counted (未成对=1) <- the fix did not blind the ruler" || bad "T3 fixed: expected 1 未成对, got $Q_UN"
say ""

# ------------------------------------------------------------------ roster diff ----
# The two-shot contrast would mean nothing if the fix had also MOVED the scope. This
# checks that the pathspec and leg-declaration lines of both rulers are identical, and it
# carries its own positive control: the two rulers' whole texts MUST differ, so a grep
# that matches nothing on both sides cannot be read as "the scopes agree".
say "===== scope control: the pathspec lines of the two rulers must be identical ====="
scope_of() { grep -E '^[[:space:]]+.[^#]*"\$GO" "\$GO2"' "$1" | sort -u; }   # just this line, no -A
su=$(scope_of "$GATE_UNFIXED")
sf=$(scope_of "$GATE_FIXED")
nu=$(printf '%s\n' "$su" | grep -c . || true)
nf=$(printf '%s\n' "$sf" | grep -c . || true)
say "  pathspec-carrying leg lines: unfixed=$nu  fixed=$nf"
# Control 1: the ruler texts themselves must differ (otherwise pre/post is a self-comparison).
cmp -s "$GATE_UNFIXED" "$GATE_FIXED" && bad "the two rulers are byte-identical - there is no pre/post contrast" ||
    ok "the two rulers differ as whole texts ($(wc -l <"$GATE_UNFIXED" | tr -d ' ') vs $(wc -l <"$GATE_FIXED" | tr -d ' ') lines)"
# Control 2: the scope grep must find something, or "equal" is vacuous.
if [ "$nu" -ge 1 ] && [ "$nf" -ge 1 ] && [ "$su" = "$sf" ]; then
    ok "the pathspec tails of both rulers are identical ($nu leg lines) <- the fix did not move the scope"
else
    bad "the two rulers disagree on pathspec tails, or the grep found nothing (unfixed=$nu fixed=$nf)"
    printf '%s\n' "$su" | sed -e 's/^/     unfixed: /'
    printf '%s\n' "$sf" | sed -e 's/^/     fixed  : /'
fi
say ""

say "r1-two-shots.sh: aggregate exit codes seen (INFORMATIONAL ONLY - the synthetic tree holds no"
say "  registered samples, so want/want_n noise there is expected and is NOT what AC#1 judges):"
for st in T1 T2 T3; do
    say "   state $st  unfixed rc=$(grep -m1 '^# 聚合退码＝' "$logdir/gate-unfixed-$st.txt" 2>/dev/null || echo NO-READING)"
    say "   state $st  fixed   rc=$(grep -m1 '^# 聚合退码＝' "$logdir/gate-fixed-$st.txt" 2>/dev/null || echo NO-READING)"
done
say ""
if [ "$fail" = 0 ]; then
    say "r1-two-shots.sh: GREEN - mode=$MODE, every assertion above held."
    exit 0
fi
say "r1-two-shots.sh: RED - at least one assertion above failed."
exit 1
