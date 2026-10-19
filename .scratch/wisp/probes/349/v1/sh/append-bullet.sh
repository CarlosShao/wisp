#!/bin/sh
set -u
REPO='/d/work/workspace/projects plans/Wisp'
T="$REPO/.scratch/wisp/issues/349-blacklist-provenance-comments-name-undefined-symbols-fires-taskscope-cleanup.md"
B=/d/tmp/wisp/20261019/349v1/msg/ticket-bullet-v1.md
echo "before: bytes=$(wc -c < "$T") lines=$(wc -l < "$T") boxes_unchecked=$(grep -c '^- \[ \] \*\*AC#' "$T") boxes_checked=$(grep -c '^- \[x\] \*\*AC#' "$T")"
echo "bullet: bytes=$(wc -c < "$B") lines=$(wc -l < "$B") U00A7=$(perl -CSD -0ne 'print scalar(() = /\x{00A7}/g)' "$B")"
cp "$T" /d/tmp/wisp/20261019/349v1/msg/ticket-before.md
cat "$B" >> "$T"
echo "after: bytes=$(wc -c < "$T") lines=$(wc -l < "$T") boxes_unchecked=$(grep -c '^- \[ \] \*\*AC#' "$T") boxes_checked=$(grep -c '^- \[x\] \*\*AC#' "$T")"
echo "appended bullets at end: $(tail -2 "$T" | head -1 | cut -c1-60)"
printf 'diff stat for ticket = %s\n' "$(git -C "$REPO" diff --numstat -- "$T" | tr -s '\t' ' ')"
printf 'U00A7 in ticket after = %s\n' "$(perl -CSD -0ne 'print scalar(() = /\x{00A7}/g)' "$T")"
echo RC=$?
