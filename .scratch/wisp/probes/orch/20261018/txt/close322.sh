#!/bin/sh
set -e
cd "/d/work/workspace/projects plans/Wisp"
T=.scratch/wisp/probes/orch/20261018/txt
F=$(ls .scratch/wisp/issues/322-*.md)
N=$(ls .scratch/wisp/issues/342-*.md | xargs -n1 basename)
echo "342 name len=${#N}"
unset GOTMPDIR TMPDIR
sh scripts/check-path-length-budget.sh | tail -2
echo "--- flip AC#4 (line 33 only) ---"
echo "before: open=$(grep -c '^[[:space:]]*- \[ \]' "$F") done=$(grep -c '^[[:space:]]*- \[x\]' "$F")"
sed -i '33s/^- \[ \]/- [x]/' "$F"
echo "after : open=$(grep -c '^[[:space:]]*- \[ \]' "$F") done=$(grep -c '^[[:space:]]*- \[x\]' "$F")"
grep -n '^[[:space:]]*- \[' "$F" | sed -n '4,6p' | cut -c1-40
echo "--- append close bullet ---"
cat "$T/322-tail-close.md" >> "$F"
NEW=$(echo "$F" | sed 's/\.md$/-done.md/')
git mv "$F" "$NEW"
echo "renamed to $NEW"
git add -- "$F" "$NEW" "$T/322-tail-close.md" .scratch/wisp/issues/342-privatefile-open-follows-a-swapped-ancestor-on-both-platforms-no-ticket-owns-it.md
echo "=== numstat of ticket ==="
git diff --cached --numstat -- "$NEW"
echo "=== uniqueness (must be 1) ==="
git ls-tree HEAD .scratch/wisp/issues/ --name-only | grep -c 'issues/322-' || true
echo "=== staged audit ==="
git diff --cached --name-only | grep -vE '^(\.scratch/wisp/issues/(322-|342-)|\.scratch/wisp/probes/orch/20261018/)' || echo "  only-mine"
