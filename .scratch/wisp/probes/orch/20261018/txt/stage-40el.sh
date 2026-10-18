#!/bin/sh
set -e
cd "/d/work/workspace/projects plans/Wisp"
T=.scratch/wisp/probes/orch/20261018/txt
git add -- \
  docs/reports/pending-and-issues.md \
  docs/reports/HANDOVER.md \
  .scratch/wisp/issues/328-l2-card-corr-shape-never-matches-the-roster-join-key-done.md \
  .scratch/wisp/issues/339-ticket334-anchor-miss-instrument-runs-on-linux-but-its-fake-home-is-windows-only.md \
  .scratch/wisp/issues/340-stacked-queue-rewrite-of-a-card-corr-has-no-real-producer-test.md \
  "$T/A1027.txt" "$T/HANDOVER-40el.md" "$T/40el-msg.txt" "$T/328-tail-bullet.md" \
  "$T/v1-verify-boundary.sh" "$T/r3c-verify-m1.sh" "$T/linuxred334.sh" "$T/flip328.sh" \
  "$T/close328.sh" "$T/stage328.sh" "$T/flip328-before.sh" "$T/commit-40el.sh"
echo "=== index contents (must be ONLY my paths; nothing of other legs) ==="
git diff --cached --name-status
echo "=== any non-mine staged? ==="
git diff --cached --name-only | grep -v -E "^(docs/reports/(pending-and-issues|HANDOVER)\.md|\.scratch/wisp/(issues/328-|issues/339-|issues/340-|probes/orch/20261018/))" || echo "  none"
