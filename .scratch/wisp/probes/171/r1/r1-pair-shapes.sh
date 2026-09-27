#!/bin/sh
# r1-pair-shapes.sh - ticket 171 r1 (写码位·仪器程), the two REVERSE JUDGEMENTS:
#   AC#6 (cell 3): does the OpenScope<->CloseScope ruler recognise the handle shape
#                  (`sc := prov.OpenScope(id)` ... `sc.Close()`), and does it still name a
#                  leak that the OLD word root could not see?
#   AC#2 (cell 2): can a SECOND unpaired file inside one already-ringing leg hide inside
#                  "non-empty", i.e. does the aggregate exit code move when 1 枚 becomes 2 枚?
#
# BOTH JUDGEMENTS NEED THE UNMODIFIED RULER. "未修码上必须不响" is only a judgement if the
# reading comes from the pre-fix bytes, so the pristine gate is taken out of a PINNED commit
# (`git show $BASE:...gate-clauses.sh`) into $TMPDIR and run there side by side with the
# working-copy gate. Same tree, same anchor, same pathspecs - the only difference is the ruler.
#
# WHY THE SAMPLES LIVE IN A SYNTHETIC TREE OUTSIDE THE REPO
# pair()/run() read `git grep "$A"` where $A defaults to a COMMIT, pinned to
# internal/**/*.go + cmd/**/*.go. A sample dropped under .scratch/wisp/probes/171/** is
# invisible twice over (not in that commit, not in those pathspecs). So the samples are
# committed into a throwaway repository under $TMPDIR - never inside this repo, no worktree,
# no checkout, no branch switch, no relaxed ':!' filter. Precedent: 161-r6's faketree.sh.
# Deleting files INSIDE that synthetic tree is the only deletion this script performs.
#
# HARD EXIT RULES
#   rc=0  every assertion held
#   rc=1  an assertion did not hold -> the cell is decoration, or the ruler lies
#   rc=2  the instrument could not run -> no verdict, never read it as green
set -eu

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)

d=$here
while [ ! -e "$d/.git" ]; do
    p=$(CDPATH= cd -- "$d/.." && pwd)
    [ "$p" = "$d" ] && { echo "r1-pair-shapes.sh: no .git above $here - refusing to answer with a ruler that did not run" >&2; exit 2; }
    d=$p
done
root=$d

gate_rel=.scratch/wisp/probes/154/gate-clauses.sh
gate_fixed="$root/$gate_rel"
[ -f "$gate_fixed" ] || { echo "r1-pair-shapes.sh: working-copy gate missing - no verdict" >&2; exit 2; }

# The pre-fix anchor: HEAD at the moment 171-r1 started work (step 0 reading).
BASE="${BASE:-19513ccefa45d51206871f00b3efa6ab19dd620a}"
git -C "$root" cat-file -e "${BASE}^{commit}" 2>/dev/null || {
    echo "r1-pair-shapes.sh: pinned pre-fix anchor $BASE is not a commit here - no verdict" >&2; exit 2; }

command -v git >/dev/null 2>&1 || { echo "r1-pair-shapes.sh: no git" >&2; exit 2; }
command -v mktemp >/dev/null 2>&1 || { echo "r1-pair-shapes.sh: no mktemp" >&2; exit 2; }

work=$(mktemp -d "${TMPDIR:-/tmp}/wisp-171-r1.XXXXXX")
logdir="$here/logs"
mkdir -p "$logdir"

gate_pristine="$work/gate-pristine.sh"
git -C "$root" show "$BASE:$gate_rel" >"$gate_pristine"
if cmp -s "$gate_pristine" "$gate_fixed"; then
    echo "r1-pair-shapes.sh: the pre-fix and working-copy gates are byte-identical - there is nothing to compare" >&2
    exit 2
fi

tree=$(mktemp -d "${TMPDIR:-/tmp}/wisp-171-r1-tree.XXXXXX")
export GIT_AUTHOR_NAME=wisp171r1 GIT_AUTHOR_EMAIL=171-r1@invalid
export GIT_COMMITTER_NAME=wisp171r1 GIT_COMMITTER_EMAIL=171-r1@invalid
git -C "$tree" init -q .
mkdir -p "$tree/internal/fake171" "$tree/cmd/wisp"

fail=0
say() { printf '%s\n' "$1"; }
ok() { say "  ok    $1"; }
bad() { say "  FAIL  $1"; fail=1; }

sec() { # <prefix> <logfile> -> one leg's section, up to the next "## "
    awk -v want="$1" '
        /^## / { if (index($0, want) == 1) ins = 1; else if (ins) ins = 0 }
        ins { print }
    ' "$2"
}

leg_field() { # <prefix> <log> <line-prefix> -> the integer on that line, or NO-READING
    raw=$(sec "$1" "$2" | grep -m1 "^$3" | sed -e "s/^$3\([0-9][0-9]*\).*$/\1/" || true)
    case $raw in
    '' | *[!0-9]*) printf 'NO-READING' ;;
    *) printf '%s' "$raw" ;;
    esac
}
leg_rc() { leg_field "$1" "$2" '# rc='; }

names() { sec "$1" "$2" | grep '^#   UNPAIRED ' | sed -e 's/(开方调用点=[0-9]*).*//' || true; }

run_gate() { # <which: pristine|fixed> <tree> -> echoes the gate's exit code, log to $logdir/<name>
    which=$1
    script=$gate_pristine
    [ "$which" = fixed ] && script=$gate_fixed
    log="$logdir/$3.txt"
    rc=0
    (cd "$tree" && sh "$script") >"$log" 2>&1 || rc=$?
    case $rc in
    *[!0-9]*) echo "run_gate: exit code '$rc' is not a number - the runner broke, no verdict" >&2; exit 2 ;;
    esac
    printf '%s' "$rc"
}

commit_state() { # <label>
    # Explicit pathspec, never a bare "add -A ." - the repo's git rule is about not sweeping a
    # shared worktree, and a throwaway tree gets the same discipline so the habit has no hole.
    git -C "$tree" add -A -- internal/fake171 >/dev/null
    git -C "$tree" commit -q -m "171-r1: $1" || { echo "r1-pair-shapes.sh: commit failed in synthetic tree" >&2; exit 2; }
    say "-- state '$1' committed at $(git -C "$tree" rev-parse --short HEAD) (synthetic tree, outside $root)"
}

gofumpt_check() { # keeps the repo's "自家台件撞全仓门" lesson honest even for throwaway samples
    if command -v gofumpt >/dev/null 2>&1; then
        out=$(gofumpt -l "$tree/internal/fake171" "$tree/cmd/wisp" 2>/dev/null || true)
        if [ -n "$out" ]; then bad "the samples this script generated are not gofumpt-clean: $out"; else ok "generated samples are gofumpt-clean ($(gofumpt --version))"; fi
    else
        say "  --    gofumpt not on PATH; sample formatting left unchecked"
    fi
}

say "r1-pair-shapes.sh: pre-fix gate = git show $BASE:$gate_rel  (written to $gate_pristine)"
say "r1-pair-shapes.sh: fixed gate   = $gate_fixed (working copy, uncommitted at the time of writing)"
say "r1-pair-shapes.sh: synthetic tree (deliberately OUTSIDE $root) at $tree"

# ============================================================================
# PART 1 - AC#6 cell 3: three shapes, each committed, each measured by BOTH rulers
# ============================================================================
say ""
say "===== PART 1 (AC#6): does the合方 side recognise the handle shape? ====="

# state A: opened via handle, closed via handle -> a genuinely balanced file.
cat >"$tree/internal/fake171/a_paired.go" <<'GO'
package fake171

import "os"

type scopeHandle struct{ id string }

type provider struct{}

func (provider) OpenScope(id string) *scopeHandle { return &scopeHandle{id: id} }

type keeper struct{ scopes map[string]*scopeHandle }

// A_pair: keeps the handle and closes it with the handle's own Close - ticket 160's shape,
// which is also what internal/tools/bridge.go:648/:706 does today.
func A(k *keeper, f *os.File) {
	k.scopes["a"] = provider{}.OpenScope("a")
	defer f.Close()
}

func (k *keeper) close(id string) {
	_ = k.scopes[id].Close()
}
GO
# The open and the handle's Close sit in two functions of the SAME file, which is exactly
# what bridge.go does (open in OpenTask :648, close in CloseTask :706): the ruler pairs per
# file, and that is the whole property it claims to see.

commit_state "A: opened via handle, handle closed in the same file (balanced)"
gofumpt_check
rcA=$(run_gate pristine "$tree" "part1-A-pristine")
rcA2=$(run_gate fixed "$tree" "part1-A-fixed")
gA=$(sec '## G5 OpenScope' "$logdir/part1-A-pristine.txt")
say "   pristine G5 names: [$(names '## G5 OpenScope' "$logdir/part1-A-pristine.txt" | tr '\n' ' ')] rc=$(leg_rc '## G5 OpenScope' "$logdir/part1-A-pristine.txt")"
say "   fixed    G5 names: [$(names '## G5 OpenScope' "$logdir/part1-A-fixed.txt" | tr '\n' ' ')] rc=$(leg_rc '## G5 OpenScope' "$logdir/part1-A-fixed.txt")"
if names '## G5 OpenScope' "$logdir/part1-A-pristine.txt" | grep -q 'a_paired.go'; then
    ok "pristine ruler names a file that DOES close the scope (this is the false red AC#6 is about)"
else
    bad "pristine ruler did not name the balanced handle file - the symptom AC#6 describes is not reproducible here"
fi
if names '## G5 OpenScope' "$logdir/part1-A-fixed.txt" | grep -q 'a_paired.go'; then
    bad "fixed ruler still names the balanced handle file - the合方 side does not know the handle shape"
else
    ok "fixed ruler goes quiet on the balanced handle file WITHOUT any pathspec being dropped"
fi
if sec '## G5 OpenScope' "$logdir/part1-A-fixed.txt" | grep -q 'CLOSE-BY-GENERIC-CLOSE .*a_paired.go'; then
    ok "fixed ruler prints who closed it (CLOSE-BY-GENERIC-CLOSE) so the coarseness is reviewable"
else
    bad "fixed ruler laundered the silent pass: it went quiet without saying that a generic Close() did it"
fi
say "   aggregate rc on state A: pristine=$rcA fixed=$rcA2"

# state B: opened via handle, never closed -> the real leak in the handle shape.
rm -f "$tree/internal/fake171/a_paired.go"
cat >"$tree/internal/fake171/b_leak.go" <<'GO'
package fake171

type scopeHandleB struct{ id string }

type providerB struct{}

func (providerB) OpenScope(id string) *scopeHandleB { return &scopeHandleB{id: id} }

type keeperB struct{ scopes map[string]*scopeHandleB }

// B_leak: opens a C25-scope-shaped handle and never closes it anywhere in this file.
func B(k *keeperB) {
	k.scopes["b"] = providerB{}.OpenScope("b")
}
GO
commit_state "B: handle-shaped OPEN, never closed (the leak AC#6 says must not become invisible)"
rcB=$(run_gate pristine "$tree" "part1-B-pristine")
rcB2=$(run_gate fixed "$tree" "part1-B-fixed")
say "   pristine G5 names: [$(names '## G5 OpenScope' "$logdir/part1-B-pristine.txt" | tr '\n' ' ')] rc=$(leg_rc '## G5 OpenScope' "$logdir/part1-B-pristine.txt")"
say "   fixed    G5 names: [$(names '## G5 OpenScope' "$logdir/part1-B-fixed.txt" | tr '\n' ' ')] rc=$(leg_rc '## G5 OpenScope' "$logdir/part1-B-fixed.txt")"
if names '## G5 OpenScope' "$logdir/part1-B-fixed.txt" | grep -q 'b_leak.go'; then
    ok "fixed ruler NAMES the handle-shaped open-without-close (AC#6 cell 3, '修完之后必须点名它')"
else
    bad "fixed ruler misses the handle-shaped leak - broadening the合方 side bought back the blind spot"
fi
if names '## G5 OpenScope' "$logdir/part1-B-pristine.txt" | grep -q 'b_leak.go'; then
    say "   --     NOTE: the pristine ruler names it TOO (its合方 root is absent here)."
    say "   --     The ticket's wording '未修码上主尺点不到它' therefore does NOT hold for this"
    say "   --     sample; see the evidence file section 3 for why (widening a合方 root can only"
    say "   --     ever reduce naming, never add it). State C is the sample where the direction holds."
else
    bad "pristine ruler missed state B too - then state B proves nothing about the fix"
fi
say "   aggregate rc on state B: pristine=$rcB fixed=$rcB2"

# state C: the shape the OLD ruler could not name at all - a discarded handle plus a
# by-id CloseScope() call for a DIFFERENT id. Word roots: open present, close present.
rm -f "$tree/internal/fake171/b_leak.go"
cat >"$tree/internal/fake171/c_discarded.go" <<'GO'
package fake171

type scopeHandleC struct{ id string }

type providerC struct{}

func (providerC) OpenScope(id string) *scopeHandleC { return &scopeHandleC{id: id} }

// legacy is the by-id closer ticket 160 removed from the real risk package; a synthetic
// tree may still spell it, and the point is what the RULER can see, not what compiles.
type legacy struct{}

func (legacy) CloseScope(id string) {}

// C_leak: throws the handle away (so nothing in this file CAN close the scope it opened)
// while a same-file CloseScope call closes somebody else's id.
func C(p providerC, l legacy) {
	_ = p.OpenScope("task-C")
	l.CloseScope("task-other")
}
GO
commit_state "C: discarded handle (_ =) plus a by-id CloseScope of another id (invisible to the old word roots)"
rcC=$(run_gate pristine "$tree" "part1-C-pristine")
rcC2=$(run_gate fixed "$tree" "part1-C-fixed")
say "   pristine G5 names: [$(names '## G5 OpenScope' "$logdir/part1-C-pristine.txt" | tr '\n' ' ')] rc=$(leg_rc '## G5 OpenScope' "$logdir/part1-C-pristine.txt")"
say "   fixed    G5 names: [$(names '## G5 OpenScope' "$logdir/part1-C-fixed.txt" | tr '\n' ' ')] rc=$(leg_rc '## G5 OpenScope' "$logdir/part1-C-fixed.txt")"
if names '## G5 OpenScope' "$logdir/part1-C-pristine.txt" | grep -q 'c_discarded.go'; then
    bad "pristine ruler already names state C - then state C cannot be the '未修码上不响' sample"
else
    ok "pristine ruler is SILENT on state C (open root + close root both present in one file)"
fi
if names '## G5 OpenScope' "$logdir/part1-C-fixed.txt" | grep -q 'c_discarded.go'; then
    ok "fixed ruler NAMES state C: a discarded handle cannot be closed by anything in the file (AC#6 cell 3's direction, both ways read)"
else
    bad "fixed ruler misses state C - the handle-discarded flavour is not wired into pair()"
fi
sec '## G5 OpenScope' "$logdir/part1-C-fixed.txt" | grep -E '^#   DISCARD-HANDLE|^#     note' | sed -e 's/^/       /'
say "   aggregate rc on state C: pristine=$rcC fixed=$rcC2"

# ============================================================================
# PART 2 - AC#2: one leg, 1 unpaired file then 2. Does the exit code move?
# ============================================================================
say ""
say "===== PART 2 (AC#2): 同一条腿内未成对从 1 枚变 2 枚 ====="

rm -f "$tree/internal/fake171/c_discarded.go"
cat >"$tree/internal/fake171/leak1.go" <<'GO'
package fake171

type scopeHandle1 struct{ id string }

type provider1 struct{}

func (provider1) OpenScope(id string) *scopeHandle1 { return &scopeHandle1{id: id} }

type keeper1 struct{ scopes map[string]*scopeHandle1 }

// Leak1: one open, no close anywhere in this file.
func Leak1(k *keeper1) {
	k.scopes["one"] = provider1{}.OpenScope("one")
}
GO
commit_state "phase1: exactly ONE unpaired file in the leg"
rc1p=$(run_gate pristine "$tree" "part2-1-pristine")
rc1f=$(run_gate fixed "$tree" "part2-1-fixed")
m1p=$(leg_rc '## G5 OpenScope' "$logdir/part2-1-pristine.txt")
m1f=$(leg_rc '## G5 OpenScope' "$logdir/part2-1-fixed.txt")
say "   phase1: G5 measures pristine=${m1p}枚 fixed=${m1f}枚 | aggregate rc pristine=${rc1p} fixed=${rc1f}"

cat >"$tree/internal/fake171/leak2.go" <<'GO'
package fake171

type scopeHandle2 struct{ id string }

type provider2 struct{}

func (provider2) OpenScope(id string) *scopeHandle2 { return &scopeHandle2{id: id} }

type keeper2 struct{ scopes map[string]*scopeHandle2 }

// Leak2: a SECOND open-without-close in a second file, same leg, same range.
func Leak2(k *keeper2) {
	k.scopes["two"] = provider2{}.OpenScope("two")
}
GO
commit_state "phase2: TWO unpaired files in the SAME leg (1 枚 -> 2 枚)"
gofumpt_check
rc2p=$(run_gate pristine "$tree" "part2-2-pristine")
rc2f=$(run_gate fixed "$tree" "part2-2-fixed")
m2p=$(leg_rc '## G5 OpenScope' "$logdir/part2-2-pristine.txt")
m2f=$(leg_rc '## G5 OpenScope' "$logdir/part2-2-fixed.txt")
say "   phase2: G5 measures pristine=${m2p}枚 fixed=${m2f}枚 | aggregate rc pristine=${rc2p} fixed=${rc2f}"
say "   which legs the FIXED aggregate calls out in phase2:"
grep -E '^# (BAD|SHR) ' "$logdir/part2-2-fixed.txt" | sed -e 's/^/       /' || say "       (none)"

if [ "$m1p" = 1 ] && [ "$m2p" = 2 ]; then
    ok "pristine ruler does see 1 枚 then 2 枚 - the reading is there, only the verdict hides it"
else
    bad "pristine ruler's per-leg counts came back ${m1p}/${m2p}, expected 1/2 - the two runs are not comparable"
fi
if [ "$rc1p" = "$rc2p" ]; then
    ok "UNMODIFIED CODE: 1 枚 -> 2 枚 moves the aggregate exit code NOT AT ALL (rc=${rc1p} both) - that is ticket 171's hole 3, reproduced"
else
    bad "UNMODIFIED CODE already reacted to the second unpaired (rc ${rc1p} -> ${rc2p}) - AC#2 would then be decoration; report, do not celebrate"
fi
if [ "$m1f" = 1 ] && [ "$m2f" = 2 ]; then
    ok "fixed ruler measures 1 枚 then 2 枚 on the same leg"
else
    bad "fixed ruler's per-leg counts came back ${m1f}/${m2f}, expected 1/2"
fi
if [ "$rc1f" != "$rc2f" ]; then
    ok "FIXED CODE: the same jump moves the aggregate exit code (${rc1f} -> ${rc2f}) - the new second file cannot hide in 'non-empty' any more"
else
    bad "FIXED CODE did not react (rc=${rc1f} both) - AC#2's 枚数 judgement is decoration"
fi
if sec '## G5 OpenScope' "$logdir/part2-2-fixed.txt" | grep -q '^# 基线枚数＝'; then
    ok "the leg prints its own baseline next to its measurement (声明与实测同一行可 diff)"
else
    bad "the leg does not print its baseline - the count claim is unsourced"
fi

# ---- keep the artefacts, print where they are -------------------------------
say ""
say "r1-pair-shapes.sh: logs kept (created, never deleted): $logdir/part1-*.txt $logdir/part2-*.txt"
say "r1-pair-shapes.sh: pristine gate copy: $gate_pristine"
say "r1-pair-shapes.sh: synthetic tree left in place at $tree; this repo's working tree was never written by this run"
if [ "$fail" = 0 ]; then
    say "r1-pair-shapes.sh: GREEN - AC#6's合方 side reads the handle shape (and names a leak the old roots could not), AC#2's exit code tracks 枚数 where the old one could not."
    exit 0
fi
say "r1-pair-shapes.sh: RED - at least one assertion above failed."
exit 1
