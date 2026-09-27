#!/bin/sh
# faketree.sh - ticket 161 AC#6 cell 2, form (甲): the SYNTHETIC-TREE proof, which AC#6
# names the MAIN evidence ("（甲）是主证：它连 pathspec 过滤与剔注释那几味一起吃住了").
# Written by 161-r6, 09-27.
#
# WHY THE TREE HAS TO BE OUTSIDE THE REPO
# pair()/run() read `git grep "$A"`, and $A's default is `git rev-parse HEAD` - a COMMIT,
# not the working tree - with the range pinned to internal/**/*.go + cmd/**/*.go. A sample
# dropped into .scratch/wisp/probes/161/** is therefore invisible twice over (not in that
# commit, not in those pathspecs). That is why the ticket's original wording "put a sample
# in probes/161/** and watch it ring" can never be done literally, and why the fix is a
# second tree that satisfies both conditions without touching this repo:
#   - it lives under $TMPDIR, never inside the repo directory (no worktree, no checkout,
#     no relaxed ':!' filter, no change to $A's default);
#   - the legs run with the SAME pathspec strings they use here, so the run exercises the
#     real ':!*_test.go' / ':!<definition file>' filters and the comment strip;
#   - the gate is invoked by absolute path with cwd inside the synthetic tree, so
#     `git rev-parse --show-toplevel` and `git rev-parse HEAD` resolve to THAT tree.
#     The instrument is unchanged; only the tree it points at differs.
#
# THREE PHASES, because "0 unpaired" has two meanings and only one ruler line tells them
# apart (pair() prints `# git grep rc=` for exactly that):
#   (1) only-open  -> G6 main and G7 main must ring, and ring in the FORWARD direction.
#                     That direction is the one this repo's own production range cannot
#                     produce at all (its only open verb sits in the excluded definition
#                     file), which is precisely why the repo's 0 must not be read as green.
#   (2) open+close -> both back to 0 while the roster stays non-empty (`git grep rc=0`):
#                     a real "nothing unpaired".
#   (3) sample gone -> still 0 but `# git grep rc=1`: "the range contains no such thing".
# Deleting inside this synthetic tree is the ONE deletion this dispatch allows.
#
# HARD EXIT RULES
#   rc=0  all three phases behaved
#   rc=1  a phase did not behave -> the new legs are decoration, or the ruler lies
#   rc=2  the instrument could not run
#
# A NOTE ON HOW THIS SCRIPT ASSERTS (found the hard way, 09-27): the first draft wrote
# `test "$g6b" = 0; check "..." $?`. Under `set -e` the failing `test` kills the script
# before `check` ever prints, so a RED run looked like a truncated log rather than a
# verdict. Every condition now goes through the helpers below, which never let a failing
# test be a bare command.
set -eu

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)

d=$here
while [ ! -e "$d/.git" ]; do
    p=$(CDPATH= cd -- "$d/.." && pwd)
    [ "$p" = "$d" ] && { echo "faketree.sh: no .git above $here - no ruler ran" >&2; exit 2; }
    d=$p
done
root=$d
gate=$root/.scratch/wisp/probes/154/gate-clauses.sh
[ -f "$gate" ] || { echo "faketree.sh: gate missing: $gate" >&2; exit 2; }
command -v git >/dev/null 2>&1 || { echo "faketree.sh: no git on PATH" >&2; exit 2; }
command -v mktemp >/dev/null 2>&1 || { echo "faketree.sh: no mktemp" >&2; exit 2; }

logdir="$here/logs"
mkdir -p "$logdir"

# Outside the repo directory, on purpose. Nothing under $root is written by this script.
tree=$(mktemp -d "${TMPDIR:-/tmp}/wisp-161-r6-faketree.XXXXXX")
export GIT_AUTHOR_NAME=wisp161r6 GIT_AUTHOR_EMAIL=161-r6@invalid
export GIT_COMMITTER_NAME=wisp161r6 GIT_COMMITTER_EMAIL=161-r6@invalid

fail=0
say() { printf '%s\n' "$1"; }

ok() { say "  ok    $1"; }
bad() { say "  FAIL  $1"; fail=1; }

# sec <prefix> <logfile> -> one leg's section, up to the next "## "
sec() {
    awk -v want="$1" '
        /^## / { if (index($0, want) == 1) ins = 1; else if (ins) ins = 0 }
        ins { print }
    ' "$2"
}

leg_field() { # <prefix> <log> <line-prefix> -> the integer on that line, or "NO-READING"
    raw=$(sec "$1" "$2" | grep -m1 "^$3" | sed -e "s/^$3\([0-9][0-9]*\).*$/\1/" || true)
    case $raw in
    '' | *[!0-9]*) printf 'NO-READING' ;;
    *) printf '%s' "$raw" ;;
    esac
}

leg_rc() { leg_field "$1" "$2" '# rc='; }
leg_grep_rc() { leg_field "$1" "$2" '# git grep rc='; }

fwd_lines() { sec "$1" "$2" | grep '^#   UNPAIRED ' || true; }

assert_num() { # <label> <value> <mode: ge1|eq0> <meaning>
    if [ "$2" = NO-READING ]; then
        bad "$1 - the leg printed no rc line at all ($4)"
        return 0
    fi
    case $2:$3 in
    *:ge1) [ "$2" -ge 1 ] && ok "$1 -> $2 ($4)" || bad "$1 -> $2, expected >= 1 ($4)" ;;
    *:eq0) [ "$2" -eq 0 ] && ok "$1 -> $2 ($4)" || bad "$1 -> $2, expected 0 ($4)" ;;
    *) bad "$1 - unknown mode $3" ;;
    esac
}

assert_has() { # <label> <file> <grep-pattern>
    if grep -q "$3" "$2"; then ok "$1"; else bad "$1 (pattern '$3' absent from $2)"; fi
}

assert_body() { # <label> <section-lines> <mode: nonempty|empty>
    if printf '%s' "$2" | grep -q .; then
        if [ "$3" = nonempty ]; then ok "$1"; else bad "$1 - the leg named files where none should be: $2"; fi
    else
        if [ "$3" = empty ]; then ok "$1"; else bad "$1 - forward list came back empty"; fi
    fi
}

say "faketree.sh: synthetic tree (deliberately OUTSIDE $root) at $tree"

# ---- the fake source, in three states ----------------------------------------
cat > "$tree/README-what-this-is.txt" <<'TXT'
Synthetic tree built by .scratch/wisp/probes/161/r6/faketree.sh (ticket 161 AC#6 cell 2,
form (甲)). It is OUTSIDE the Wisp repository on purpose: the gate reads a commit and a
fixed pathspec range, so a sample can only prove a leg rings by being both committed and
inside internal/**/*.go / cmd/**/*.go. Nothing here is a claim about Wisp's own code.
TXT

mkdir -p "$tree/internal/fake"
git -C "$tree" init -q .

# phase 1: every OPEN verb present, every CLOSE/scope token absent.
# G7's close token must be absent entirely - neither "DisposalScope" nor "NewDisposalScope"
# appears in this file, so the forward reading is not an artefact of the (New)? alternation.
cat > "$tree/internal/fake/x.go" <<'GO'
package fake

// Phase 1 of ticket 161 AC#6 (甲): the open side of both new families is here and the
// close side is nowhere. This is the shape the two new legs are supposed to catch.
type scopeish struct{}

func (scopeish) Defer(fn func()) bool { _ = fn; return true }

func (scopeish) DeferNamed(name string, fn func()) bool { _ = name; _ = fn; return true }

type taskish struct{ id string }

func (t taskish) OpenTask(id string) {}

func Run(t taskish, s scopeish) {
	t.OpenTask("task-1")
	s.Defer(func() {})
}
GO
git -C "$tree" add README-what-this-is.txt internal/fake/x.go
if git -C "$tree" commit -q -m "phase1: open without close (161 AC#6 jia)"; then
    say "faketree.sh: phase 1 committed at $(git -C "$tree" rev-parse --short HEAD) - the gate's \$A resolves to this tree because cwd is inside it"
else
    echo "faketree.sh: git commit failed in the synthetic tree - no verdict" >&2
    exit 2
fi

log1="$logdir/faketree-phase1.txt"
if (cd "$tree" && sh "$gate") >"$log1" 2>&1; then :; fi
g6=$(leg_rc '## G6 主尺 OpenTask' "$log1")
g7=$(leg_rc '## G7 主尺 Defer' "$log1")
say "== phase 1 (open, never closed)"
sec '## G6 主尺 OpenTask' "$log1" | grep -E '^#   UNPAIRED|^# 未成对枚数|^# 反向未成对枚数|^# git grep rc' | sed -e 's/^/   /'
sec '## G7 主尺 Defer' "$log1" | grep -E '^#   UNPAIRED|^# 未成对枚数|^# 反向未成对枚数|^# git grep rc' | sed -e 's/^/   /'
assert_num "G6 main rings in the synthetic tree" "$g6" ge1 ">=1 unpaired"
assert_num "G7 main rings in the synthetic tree" "$g7" ge1 ">=1 unpaired"
assert_body "G6 rings in the FORWARD direction (open-without-close)" "$(fwd_lines '## G6 主尺 OpenTask' "$log1")" nonempty
assert_body "G7 rings in the FORWARD direction (Defer with no scope token)" "$(fwd_lines '## G7 主尺 Defer' "$log1")" nonempty
assert_has "the named file is the one this harness just created" "$log1" 'internal/fake/x.go'

# phase 2: same file, both sides present.
cat > "$tree/internal/fake/x.go" <<'GO'
package fake

// Phase 2 of ticket 161 AC#6 (甲): the same file now carries both sides of each pair, so
// the two new legs must go quiet WITHOUT the roster going empty - which is what makes
// phase 3's identical 0 a different reading.
type DisposalScope struct{}

func (DisposalScope) Defer(fn func()) bool { _ = fn; return true }

func (DisposalScope) DeferNamed(name string, fn func()) bool { _ = name; _ = fn; return true }

func (DisposalScope) Dispose() {}

type taskish struct{ id string }

func (t taskish) OpenTask(id string)  {}
func (t taskish) CloseTask(id string) {}

func Run(t taskish) {
	s := &DisposalScope{}
	t.OpenTask("task-1")
	defer t.CloseTask("task-1")
	defer s.Defer(func() {})
	_ = s
}
GO
git -C "$tree" add README-what-this-is.txt internal/fake/x.go
git -C "$tree" commit -q -m "phase2: open and close in one file (161 AC#6 jia)"

log2="$logdir/faketree-phase2.txt"
if (cd "$tree" && sh "$gate") >"$log2" 2>&1; then :; fi
say "== phase 2 (opened, then closed in the same file)"
assert_num "G6 main back to 0" "$(leg_rc '## G6 主尺 OpenTask' "$log2")" eq0 "a real 0"
assert_num "G7 main back to 0" "$(leg_rc '## G7 主尺 Defer' "$log2")" eq0 "a real 0"
assert_num "G6 still sees the family (git grep rc=0, so this 0 is not an empty ruler)" "$(leg_grep_rc '## G6 主尺 OpenTask' "$log2")" eq0 "roster non-empty"
assert_num "G7 still sees the family" "$(leg_grep_rc '## G7 主尺 Defer' "$log2")" eq0 "roster non-empty"

# phase 3: the sample is gone (the ONE deletion this dispatch permits, inside a tree
# that lives outside the repository).
rm -f "$tree/internal/fake/x.go"
git -C "$tree" add README-what-this-is.txt internal/fake/x.go
git -C "$tree" commit -q -m "phase3: sample removed inside the synthetic tree (161 AC#6 jia)"

log3="$logdir/faketree-phase3.txt"
if (cd "$tree" && sh "$gate") >"$log3" 2>&1; then :; fi
say "== phase 3 (nothing of the family left in the range)"
assert_num "G6 reads 0" "$(leg_rc '## G6 主尺 OpenTask' "$log3")" eq0 "count"
assert_num "G7 reads 0" "$(leg_rc '## G7 主尺 Defer' "$log3")" eq0 "count"
g6grc=$(leg_grep_rc '## G6 主尺 OpenTask' "$log3")
g7grc=$(leg_grep_rc '## G7 主尺 Defer' "$log3")
if [ "$g6grc" = 1 ]; then
    ok "G6 prints 'git grep rc=1' - pair() calls that '射程里没这东西', not '都关好了' (got $g6grc)"
else
    bad "G6's git grep rc is $g6grc, expected 1 (the 0 would otherwise be indistinguishable from phase 2's)"
fi
if [ "$g7grc" = 1 ]; then
    ok "G7 prints the same (got $g7grc)"
else
    bad "G7's git grep rc is $g7grc, expected 1"
fi

printf '\n'
say "faketree.sh: logs kept (created, never deleted): $log1 $log2 $log3"
say "faketree.sh: synthetic tree left in place at $tree; the Wisp working tree was never written by this run"
if [ "$fail" = 0 ]; then
    say "faketree.sh: GREEN - the two new AC#6 legs name a freshly created unpaired sample, in the forward direction the repo's own production range structurally cannot produce, and go quiet again when it is paired."
    exit 0
fi
say "faketree.sh: RED - a new leg failed to ring, rang without naming the file, or could not tell a real 0 from an empty ruler."
exit 1
