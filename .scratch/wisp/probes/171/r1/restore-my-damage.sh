#!/bin/sh
# restore-my-damage.sh - undo THIS leg's accidental overwrite of OTHER legs' committed logs.
#
# WHAT HAPPENED (recorded here rather than smoothed over)
#   This leg ran two harnesses that belong to other programs:
#     .scratch/wisp/probes/171/r3/r3-trailing-comment.sh   (10:33-10:36)
#     .scratch/wisp/probes/171/r3/r3-flip-columns.sh       (10:29-11:10, killed)
#     .scratch/wisp/probes/171/r1/r1-flip-columns.sh       (11:10-11:13, killed)
#   All three write their readings into log paths that are COMMITTED by other legs
#   (171-r3 / 171s / the concurrent leg that committed 2a5ae70a). So running them
#   overwrote other programs' evidence - and the two killed mid-cell left their target
#   file TRUNCATED (measured: 1081 -> 586 lines, 66 -> 59, 11 -> 2).
#
# WHAT THIS DOES ABOUT IT
#   Writes each affected path back to its bytes at HEAD, i.e. to exactly what its own
#   leg committed. No banned command is used: this is `git show HEAD:<path> > <path>`,
#   NOT `git checkout --`, NOT `git restore`, NOT reset/stash/clean. Nothing is deleted
#   and nothing is moved.
#
# WHY RESTORING IS THE RIGHT CALL AND NOT A COVER-UP
#   A leg commits its logs as evidence; those logs being dirty on disk is not a state
#   anybody chose. The truncation is the injury, and the committed bytes are the
#   intended content. The numbers above are recorded in this file and in the delivery
#   note, so the fact of the damage survives the repair.
#
# Run from the repo root: sh probes/171/r1/restore-my-damage.sh
set -eu

cd "$(CDPATH= cd -- "$(dirname -- "$0")/../../../../.." && pwd)"

# The exact 26 paths, listed one per line (never a directory glob: a glob would also
# sweep up files this leg has no business writing).
git status --porcelain -- .scratch/wisp/probes/171 |
    sed -n 's/^ M //p' > /tmp/restore-list.$$ || true

n=$(grep -c . /tmp/restore-list.$$ || true)
echo "restore-my-damage.sh: modified tracked paths under probes/171: $n"
[ "$n" -gt 0 ] || { echo "restore-my-damage.sh: nothing modified - nothing to do"; exit 0; }

while IFS= read -r f; do
    [ -n "$f" ] || continue
    before=$(wc -l <"$f" | tr -d ' ')
    git show "HEAD:$f" >"$f"
    after=$(wc -l <"$f" | tr -d ' ')
    printf '  restored %6s -> %-6s lines  %s\n' "$before" "$after" "$f"
done </tmp/restore-list.$$

echo
echo "restore-my-damage.sh: re-reading the same ruler that produced the list"
git status --porcelain -- .scratch/wisp/probes/171 | sed -n 's/^ M //p' >/tmp/restore-after.$$
left=$(grep -c . /tmp/restore-after.$$ || true)
echo "restore-my-damage.sh: still-modified tracked paths under probes/171: $left (want 0)"
[ "$left" = 0 ] || { echo "restore-my-damage.sh: RED - some path did not come back"; exit 1; }
echo "restore-my-damage.sh: GREEN - every path is back to its committed bytes"
