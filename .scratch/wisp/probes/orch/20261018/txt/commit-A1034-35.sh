#!/bin/sh
set -e
cd "/d/work/workspace/projects plans/Wisp"
T=.scratch/wisp/probes/orch/20261018/txt
cat "$T/A1034.txt" >> docs/reports/pending-and-issues.md
cat "$T/A1035.txt" >> docs/reports/pending-and-issues.md
cat "$T/HANDOVER-40er.md" >> docs/reports/HANDOVER.md
git add -- docs/reports/pending-and-issues.md docs/reports/HANDOVER.md "$T/A1034.txt" "$T/A1035.txt" "$T/HANDOVER-40er.md" "$T/40er-msg.txt" "$T/339r1-head.txt"
echo "=== staged audit ==="
git diff --cached --name-only | grep -vE '^(\.scratch/wisp/probes/orch/20261018/|docs/reports/pending-and-issues\.md|docs/reports/HANDOVER\.md)' || echo "  only-mine"
git commit -F "$T/40er-msg.txt" -- docs/reports/pending-and-issues.md docs/reports/HANDOVER.md "$T/A1034.txt" "$T/A1035.txt" "$T/HANDOVER-40er.md" "$T/40er-msg.txt" "$T/339r1-head.txt"
echo "=== self-check ==="
git show --name-only --format='' HEAD | grep -vE '^(\.scratch/wisp/probes/orch/20261018/|docs/reports/pending-and-issues\.md|docs/reports/HANDOVER\.md)' || echo "  commit-clean: only mine"
git log -1 --format='%h'
echo "ledger=$(wc -l < docs/reports/pending-and-issues.md) handover=$(wc -l < docs/reports/HANDOVER.md) unpushed=$(git rev-list --count origin/dev..HEAD)"
