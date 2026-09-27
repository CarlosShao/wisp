#!/bin/sh
# tracked-dirty-proof.sh - ticket 161-r5 cell 3.
#
# WHAT IS UNPROVEN ABOUT attrib.sh'S FORM (A)
# The r4 program registered this honestly: "（甲）的 TRACKED-DIRTY 那一发从未实测红过
# - 要证它得弄脏一枚已跟踪 .go，共享树里不该由实现程顺手做". Form (A) is the leg CI
# executes, so a leg that has never been seen red is a leg nobody has checked.
# Proving it the obvious way means dirtying somebody else's tracked file in a shared
# working tree, which is exactly what this repo forbids.
#
# THE SHAPE USED INSTEAD (and the two shapes NOT used)
# A git checkout of the repo's tracked .go SUBSET, copied OUT OF THE REPO into a
# temporary directory, given its own index, dirtied there, measured, restored,
# measured again. Explicitly avoided:
#   - a worktree or checkout INSIDE the repository (banned outright here, and it
#     would also make form (B) read someone else's bench tree),
#   - `git archive | tar -x` as "clean tree" evidence (this repo measured that
#     `* text=auto` + core.autocrlf=true rewrite line endings, so an archive-extract
#     would not be byte-comparable to the tracked set - and gofumpt cares).
# So: `git ls-files '*.go'` for the roster, `cp --parents` for the bytes, and a real
# `git init` + `git add` in the copy so that the SAME question form (A) asks - "what
# does the index say is tracked?" - has a truthful answer there too.
#
# The copy is left ON DISK when this finishes (create-only, and this program runs no
# delete commands at all). It is outside the repository; its path is printed. If it
# ever needs retiring, that is the owner's call, not something to bury in a trap.
#
# USAGE
#   sh .scratch/wisp/probes/161/r5/tracked-dirty-proof.sh [subset-size]   # default 40

set -eu

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
SUBSET=${1:-40}

# ---- derive the real repo root the same way attrib.sh does, then leave it -----
d=$here
while [ ! -e "$d/.git" ]; do
    p=$(CDPATH= cd -- "$d/.." && pwd)
    [ "$p" = "$d" ] && { echo "tracked-dirty-proof: no .git above $here" >&2; exit 2; }
    d=$p
done
repo=$d
[ -f "$repo/go.mod" ] || { echo "tracked-dirty-proof: no go.mod at $repo" >&2; exit 2; }

tmp=$(mktemp -d "${TMPDIR:-/tmp}/wisp-161-r5-trackeddirty.XXXXXX")
# mktemp resolves through $TMPDIR, which is outside the repository by construction;
# assert it rather than trusting the name, because a copy that lands back inside the
# repo would be the one banned shape this whole script exists to avoid.
case $tmp/ in
"$repo"/*) echo "tracked-dirty-proof: temp dir $tmp is INSIDE $repo - refusing, that shape is banned here" >&2; exit 2 ;;
esac
echo "tracked-dirty-proof: repo  = $repo"
echo "tracked-dirty-proof: copy  = $tmp   (outside the repo; left in place when this ends)"

cd "$repo"

# ---- build the subset roster, then copy it with its directory structure -------
git ls-files '*.go' | sed -n "1,${SUBSET}p" > "$tmp/subset-roster.txt"
N=$(grep -c . "$tmp/subset-roster.txt" || true)
echo "tracked-dirty-proof: roster = git ls-files '*.go' | sed -n '1,${SUBSET}p'  -> $N files, copied with cp --parents"
[ "$N" -gt 0 ] || { echo "tracked-dirty-proof: empty roster - a ruler handed nothing is not evidence" >&2; exit 2; }
VICTIM=$(sed -n '1p' "$tmp/subset-roster.txt")
echo "tracked-dirty-proof: victim = $VICTIM  (the file that gets dirtied, inside the copy only)"

# go.mod is required by the instrument's own root check, not by gofumpt.
# .gitattributes travels with it on purpose: without it the copy's `git add` runs
# under this machine's core.autocrlf=true and spends the run printing 40 CRLF
# warnings while claiming nothing about line endings. With it, the copy normalizes
# exactly like the repo does (`*.go text eol=lf`), so what gofumpt reads in the copy
# is byte-for-byte what it reads in a CI checkout. That matters because the very
# reason this script does not use `git archive | tar -x` is line-ending drift.
cp --parents go.mod .gitattributes "$tmp"/
while IFS= read -r f; do
    [ -n "$f" ] || continue
    cp --parents "$f" -t "$tmp"/
done < "$tmp/subset-roster.txt"

# ---- give the copy an index, so 'tracked' means what form (A) thinks it means --
cd "$tmp"
git init -q .
git add -- go.mod "$VICTIM"
while IFS= read -r f; do
    [ -n "$f" ] || continue
    git add -- "$f"
done < subset-roster.txt
echo "tracked-dirty-proof: copy's own index now tracks $(git ls-files -z '*.go' | tr -dc '\0' | wc -c | tr -d ' ') .go files"

# ---- the run helper: invokes the SAME instrument CI's new step invokes ---------
# run_a <label> prints the instrument's verdict lines and then its exit code, and
# hands that code back to the caller. Returning it explicitly is the point: a helper
# that swallowed it would turn this script's own "did the leg ring?" answer into the
# empty kind of green this ticket keeps refusing.
run_a() {
    label=$1
    rc=0
    ATTRIB_ROOT=$tmp sh "$repo/.scratch/wisp/probes/161/r5/attrib.sh" --tracked-only || rc=$?
    printf 'tracked-dirty-proof: %s rc=%s\n' "$label" "$rc"
    return "$rc"
}

echo
echo "===== 1. pristine copy -> expect GREEN (this is form (A) answering normally) ====="
rc=0
run_a "pristine" || rc=$?
echo "TRACKED_PRISTINE_RC=$rc"

echo
echo "===== 2. dirty $VICTIM inside the copy -> expect RED (TRACKED-DIRTY, rc=1) ====="
# Two-space indent in a composite literal: parses fine, and gofumpt will not accept
# the shape. Appended, never edited in place, so the pristine bytes stay recoverable
# from the repo without touching the copy's history.
printf '\nvar dirty161r5 = map[string]int{\n  "two space indent is not gofumpt": 1,\n}\n' >> "$VICTIM"
rc2=0
run_a "dirtied" || rc2=$?
echo "TRACKED_DIRTY_RC=$rc2"
if [ "$rc2" = 0 ]; then
    echo "tracked-dirty-proof: THE LEG DID NOT RING - form (A) is unproven, this is a finding, not a pass" >&2
    exit 3
fi

echo
echo "===== 3. restore the victim from the repo's tracked bytes -> expect GREEN ====="
( cd "$repo" && cp --parents "$VICTIM" -t "$tmp"/ )
cmp "$repo/$VICTIM" "$tmp/$VICTIM" && echo "tracked-dirty-proof: byte-identical to $repo/$VICTIM after restore"
rc3=0
run_a "restored" || rc3=$?
echo "TRACKED_RESTORED_RC=$rc3"

echo
echo "===== 4. the SAME command handed ZERO tracked .go files -> expect rc=2, never rc=0 ====="
# This is the leg that keeps step 2 of this script's exit code honest: a formatting
# gate whose denominator comes from `git ls-files` can be handed nothing by a typo'd
# switch (r4's own history: `git ls-files -Z` is not a switch, the pipeline still
# exited 0, and the first reading that file ever printed was vacuously green). The
# CI step added for 161-r5 runs exactly this command, so the refusal is measured
# here rather than asserted in a comment.
tmp2=$(mktemp -d "${TMPDIR:-/tmp}/wisp-161-r5-hollowguard.XXXXXX")
case $tmp2/ in
"$repo"/*) echo "tracked-dirty-proof: temp dir $tmp2 is INSIDE $repo - refusing" >&2; exit 2 ;;
esac
( cd "$tmp2" && git init -q . && cp "$repo/go.mod" . && cp "$repo/.gitattributes" . )
rc4=0
ATTRIB_ROOT=$tmp2 sh "$repo/.scratch/wisp/probes/161/r5/attrib.sh" --tracked-only || rc4=$?
echo "tracked-dirty-proof: hollow-guard rc=$rc4  (copy=$tmp2, tracked .go files in it: $(cd "$tmp2" && git ls-files -z '*.go' | tr -dc '\0' | wc -c | tr -d ' '))"
if [ "$rc4" = 0 ]; then
    echo "tracked-dirty-proof: THE HOLLOW GUARD IS GONE - a ruler handed nothing reported green, which is the shape AC#7 exists to refuse" >&2
    exit 3
fi

echo
echo "tracked-dirty-proof: readings - pristine rc=$rc, dirtied rc=$rc2, restored rc=$rc3, handed-nothing rc=$rc4"
if [ "$rc" = 0 ] && [ "$rc2" = 1 ] && [ "$rc3" = 0 ] && [ "$rc4" = 2 ]; then
    echo "tracked-dirty-proof: form (A)'s TRACKED-DIRTY leg is now MEASURED in both directions, on a tree nobody else shares, and its 0-file refusal is measured too"
    exit 0
fi
echo "tracked-dirty-proof: unexpected combination of exit codes - read the four blocks above, do not re-run this into a story" >&2
exit 1
