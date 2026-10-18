#!/bin/sh
set -e
cd "/d/work/workspace/projects plans/Wisp"
T=.scratch/wisp/probes/orch/20261018/txt
cat "$T/A1033.txt" >> docs/reports/pending-and-issues.md
cat "$T/HANDOVER-40eq.md" >> docs/reports/HANDOVER.md
git add -- docs/reports/pending-and-issues.md docs/reports/HANDOVER.md "$T/A1033.txt" "$T/HANDOVER-40eq.md" "$T/40eq-msg.txt" "$T/227-verify.sh" "$T/227-verify2.sh" "$T/227-verify3.sh" .scratch/wisp/issues/341-four-hand-edit-shapes-still-unnailed-and-two-226-sentences-overclaim.md
echo "=== staged audit ==="
git diff --cached --name-only | grep -vE '^(\.scratch/wisp/probes/orch/20261018/|\.scratch/wisp/issues/341-|docs/reports/(pending-and-issues|HANDOVER)\.md)$' || echo "  only-mine"
git commit -F "$T/40eq-msg.txt" -- docs/reports/pending-and-issues.md docs/reports/HANDOVER.md "$T/A1033.txt" "$T/HANDOVER-40eq.md" "$T/40eq-msg.txt" "$T/227-verify.sh" "$T/227-verify2.sh" "$T/227-verify3.sh" .scratch/wisp/issues/341-four-hand-edit-shapes-still-unnailed-and-two-226-sentences-overclaim.md
echo "=== self-check ==="
git show --name-only --format='' HEAD | grep -vE '^(\.scratch/wisp/probes/orch/20261018/|\.scratch/wisp/issues/341-|docs/reports/(pending-and-issues|HANDOVER)\.md)$' || echo "  commit-clean: only mine"
git log -1 --format='%h'
echo "ledger=$(wc -l < docs/reports/pending-and-issues.md) handover=$(wc -l < docs/reports/HANDOVER.md) unpushed=$(git rev-list --count origin/dev..HEAD)"
