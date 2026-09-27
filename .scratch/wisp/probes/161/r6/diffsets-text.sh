#!/bin/sh
# diffsets-text.sh - ticket 161 AC#6 cell 2, form (乙): the TEXT-FED proof.
# Written by 161-r6, 09-27.
#
# WHAT IT PROVES, AND WHAT IT DELIBERATELY DOES NOT PROVE
# AC#6 names (甲) the synthetic tree as the MAIN evidence and (乙) as arithmetic-only:
# "（乙）只证差集算术本身不撒谎". So this harness feeds hand-written rosters straight into
# the same diffsets() that pair() calls (that is why AC#6 cell 2 required the差集 to be
# factored out into a function - otherwise the sample and the reading would be produced by
# two different rulers, and the proof would prove nothing).
# It does NOT prove anything about Wisp's own code: no git, no anchor, no pathspec range
# is consulted here. The rosters are the same file:line:text shape pair() prints, so the
# comment-stripping line is exercised too, and one of the four cases below exists precisely
# to show that (a family mentioned only inside a // comment must not count as an open).
#
# USAGE
#   sh .scratch/wisp/probes/161/r6/diffsets-text.sh
# Samples are written next to this script (created, never deleted) so the next run can
# diff them: .scratch/wisp/probes/161/r6/samples/*.txt
#
# HARD EXIT RULES: rc=0 all four cases behaved, rc=1 one did not, rc=2 could not run.
set -eu

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
gate=$here/../../154/gate-clauses.sh
[ -f "$gate" ] || gate=$(git -C "$here" rev-parse --show-toplevel)/.scratch/wisp/probes/154/gate-clauses.sh
[ -f "$gate" ] || { echo "diffsets-text.sh: gate-clauses.sh not found - no verdict" >&2; exit 2; }

sam=$here/samples
mkdir -p "$sam"
logdir=$here/logs
mkdir -p "$logdir"

# ---- the four hand-written rosters -------------------------------------------
# ⚠ THE SHAPE IS NOT COSMETIC (measured while writing this harness, 09-27):
# pair()/diffsets() read the file out of a git-grep line with `cut -d: -f2`, and git grep
# prefixes every hit with `<commit>:`. Feed it `path:line:text` and field 2 is the LINE
# NUMBER - the ruler then prints "9" as an unpaired file and both directions come back
# non-empty. The first draft of these samples did exactly that and 乙-2 "failed"; the
# arithmetic was fine, the sample was lying about the format. Real git grep output
# (`git grep -nEw X <anchor> -- ...`) is 4 fields, so these samples are 4 fields too.
#
# (1) opened, never closed in that file -> forward MUST name the file
cat > "$sam/open-only.txt" <<'TXT'
c0ffee0000000000000000000000000000000000:internal/fake/x.go:9:	t.OpenTask("task-1")
c0ffee0000000000000000000000000000000000:internal/fake/y.go:4:	s.CloseTask("task-1")
TXT
# (2) opened and closed in the SAME file -> both directions MUST be empty
cat > "$sam/open-then-close.txt" <<'TXT'
c0ffee0000000000000000000000000000000000:internal/fake/x.go:9:	t.OpenTask("task-1")
c0ffee0000000000000000000000000000000000:internal/fake/x.go:11:	defer t.CloseTask("task-1")
c0ffee0000000000000000000000000000000000:internal/fake/y.go:4:	s.OpenTask("task-2")
c0ffee0000000000000000000000000000000000:internal/fake/y.go:8:	s.CloseTask("task-2")
TXT
# (3) the only open verb sits in a comment -> forward MUST be empty.
# This is the one taste of (甲)'s main evidence that (乙) can reproduce on its own:
# pair() strips pure-comment call lines before求差集, and this mode uses the same strip.
cat > "$sam/comment-only.txt" <<'TXT'
c0ffee0000000000000000000000000000000000:internal/fake/z.go:3:// OpenTask opens the C25 taint scope for one task.
c0ffee0000000000000000000000000000000000:internal/fake/z.go:9:	s.CloseTask("task-1")
TXT
# (4) the second family: Defer registered, no DisposalScope anywhere in the file
cat > "$sam/defer-no-scope.txt" <<'TXT'
c0ffee0000000000000000000000000000000000:internal/fake/w.go:12:	s.Defer(func() {})
c0ffee0000000000000000000000000000000000:internal/fake/w.go:13:	s.DeferNamed("cleanup", func() {})
c0ffee0000000000000000000000000000000000:internal/fake/v.go:7:	scope := plugin.NewDisposalScope("task", ctx)
c0ffee0000000000000000000000000000000000:internal/fake/v.go:8:	scope.Defer(func() {})
TXT

fail=0

fed() { # <label> <open> <close> <roster> <expect: forward|empty>
    label=$1
    open=$2
    close=$3
    roster=$4
    expect=$5
    out="$logdir/diffsets-text-$(printf '%s' "$label" | tr -c 'A-Za-z0-9' '-').txt"
    printf '%s\n' "== $label"
    DIFFSETS_OPEN="$open" DIFFSETS_CLOSE="$close" sh "$gate" --diffsets <"$roster" >"$out" 2>&1
    grep -v '^# 票 154\|^# 生成时刻\|^# 下一程复算' "$out" | sed -e 's/^/   /'
    # the forward block = the lines between "# 正向" and "# 反向", minus the "（空）" marker
    fwd=$(awk '/^# 正向/{f=1; next} /^# 反向/{f=0} f && /^#   /' "$out" | grep -v '（空）' || true)
    case $expect in
    forward)
        if printf '%s' "$fwd" | grep -q .; then
            printf '  ok    forward is NOT empty (that is the ring this mode is supposed to show)\n'
        else
            printf '  FAIL  forward came back empty - diffsets() cannot do the one thing it is for\n'
            fail=1
        fi
        ;;
    empty)
        if printf '%s' "$fwd" | grep -q .; then
            printf '  FAIL  forward is not empty, but the sample is paired by construction:\n%s\n' "$fwd"
            fail=1
        else
            printf '  ok    forward is empty (the arithmetic does not invent unpaired files)\n'
        fi
        ;;
    *) echo "fed: unknown expectation $expect" >&2; exit 2 ;;
    esac
}

fed "乙-1 OpenTask/CloseTask, opened never closed" OpenTask CloseTask "$sam/open-only.txt" forward
fed "乙-2 OpenTask/CloseTask, opened then closed in the same file" OpenTask CloseTask "$sam/open-then-close.txt" empty
fed "乙-3 OpenTask/CloseTask, open verb only inside a comment" OpenTask CloseTask "$sam/comment-only.txt" empty
fed "乙-4 Defer/DisposalScope, Defer with no scope token in the file" 'Defer(Named)?' '(New)?DisposalScope' "$sam/defer-no-scope.txt" forward

printf '\n'
printf 'diffsets-text.sh: samples kept (created, never deleted): %s\n' "$sam"
if [ "$fail" = 0 ]; then
    printf 'diffsets-text.sh: GREEN - the差集 arithmetic behind both new legs answers the two ways it must, and eats the comment strip on the way.\n'
    exit 0
fi
printf 'diffsets-text.sh: RED - diffsets() lied about one of the hand-written rosters.\n'
exit 1
