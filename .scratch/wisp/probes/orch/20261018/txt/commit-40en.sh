#!/bin/sh
set -e
cd "/d/work/workspace/projects plans/Wisp"
T=.scratch/wisp/probes/orch/20261018/txt
cat "$T/A1029.txt" >> docs/reports/pending-and-issues.md
cat "$T/HANDOVER-40en.md" >> docs/reports/HANDOVER.md
git add -- docs/reports/pending-and-issues.md docs/reports/HANDOVER.md "$T/A1029.txt" "$T/HANDOVER-40en.md" "$T/40en-msg.txt" "$T/338r1-head.txt" "$T/338r1-specific.txt" "$T/338r1-dispatch.txt" "$T/commit-40en.sh"
echo "=== staged audit (must be only mine) ==="
git diff --cached --name-only | grep -v -E "^(\.scratch/wisp/probes/orch/20261018/|docs/reports/(pending-and-issues|HANDOVER)\.md)" || echo "  only-mine"
git commit -F "$T/40en-msg.txt" -- docs/reports/pending-and-issues.md docs/reports/HANDOVER.md "$T/A1029.txt" "$T/HANDOVER-40en.md" "$T/40en-msg.txt" "$T/338r1-head.txt" "$T/338r1-specific.txt" "$T/338r1-dispatch.txt" "$T/commit-40en.sh"
echo "=== after ==="
git log -1 --format='%h' 
git show --name-only --format='' HEAD | grep -v -E "^(\.scratch/wisp/probes/orch/20261018/|docs/reports/)" || echo "  commit-clean: only mine"
echo "unpushed=$(git rev-list --count origin/dev..HEAD)"
