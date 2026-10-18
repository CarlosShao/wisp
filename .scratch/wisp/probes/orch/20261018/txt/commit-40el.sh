#!/bin/sh
set -e
cd "/d/work/workspace/projects plans/Wisp"
T=.scratch/wisp/probes/orch/20261018/txt
cat "$T/A1027.txt" >> docs/reports/pending-and-issues.md
cat "$T/HANDOVER-40el.md" >> docs/reports/HANDOVER.md
echo "ledger=$(wc -l < docs/reports/pending-and-issues.md) handover=$(wc -l < docs/reports/HANDOVER.md)"
git add -- \
  docs/reports/pending-and-issues.md \
  docs/reports/HANDOVER.md \
  .scratch/wisp/issues/328-l2-card-corr-shape-never-matches-the-roster-join-key.md \
  .scratch/wisp/issues/328-l2-card-corr-shape-never-matches-the-roster-join-key-done.md \
  .scratch/wisp/issues/339-ticket334-anchor-miss-instrument-runs-on-linux-but-its-fake-home-is-windows-only.md \
  .scratch/wisp/issues/340-stacked-queue-rewrite-of-a-card-corr-has-no-real-producer-test.md \
  "$T/A1027.txt" "$T/HANDOVER-40el.md" "$T/40el-msg.txt" "$T/328-tail-bullet.md" \
  "$T/v1-verify-boundary.sh" "$T/r3c-verify-m1.sh" "$T/linuxred334.sh" "$T/flip328.sh" \
  "$T/close328.sh" "$T/stage328.sh" "$T/flip328-before.sh"
echo "=== staged ==="
git diff --cached --name-status
echo "=== commit ==="
git commit -F "$T/40el-msg.txt" -- \
  docs/reports/pending-and-issues.md \
  docs/reports/HANDOVER.md \
  .scratch/wisp/issues/328-l2-card-corr-shape-never-matches-the-roster-join-key.md \
  .scratch/wisp/issues/328-l2-card-corr-shape-never-matches-the-roster-join-key-done.md \
  .scratch/wisp/issues/339-ticket334-anchor-miss-instrument-runs-on-linux-but-its-fake-home-is-windows-only.md \
  .scratch/wisp/issues/340-stacked-queue-rewrite-of-a-card-corr-has-no-real-producer-test.md \
  "$T/A1027.txt" "$T/HANDOVER-40el.md" "$T/40el-msg.txt" "$T/328-tail-bullet.md" \
  "$T/v1-verify-boundary.sh" "$T/r3c-verify-m1.sh" "$T/linuxred334.sh" "$T/flip328.sh" \
  "$T/close328.sh" "$T/stage328.sh" "$T/flip328-before.sh"
echo "=== after ==="
git log -1 --format='%h %s' | cut -c1-90
git ls-tree HEAD .scratch/wisp/issues/ --name-only | grep -c "issues/328-" || true
git rev-list --count origin/dev..HEAD
git status --porcelain -- .scratch/wisp/issues/ docs/reports/ | head -5
