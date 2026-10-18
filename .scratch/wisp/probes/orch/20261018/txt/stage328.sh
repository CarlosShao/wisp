#!/bin/sh
set -e
cd "/d/work/workspace/projects plans/Wisp"
NEW=.scratch/wisp/issues/328-l2-card-corr-shape-never-matches-the-roster-join-key-done.md
git add -- "$NEW" .scratch/wisp/issues/340-stacked-queue-rewrite-of-a-card-corr-has-no-real-producer-test.md .scratch/wisp/probes/orch/20261018/txt/328-tail-bullet.md
echo "=== staged (must show R + the two new tickets) ==="
git diff --cached --name-status
echo "=== numstat for the renamed ticket ==="
git diff --cached --numstat -- "$NEW"
echo "=== boxes in staged version ==="
git show ":$NEW" | grep -c '^[[:space:]]*- \[ \]' || true
git show ":$NEW" | grep -c '^[[:space:]]*- \[x\]' || true
echo "=== worktree clean for that file? ==="
git status --porcelain -- "$NEW" | sed 's/^/  /'
echo "=== 339/340 box counts ==="
for f in .scratch/wisp/issues/339-*.md .scratch/wisp/issues/340-*.md; do echo "  $f open=$(grep -c '^[[:space:]]*- \[ \]' "$f")"; done
