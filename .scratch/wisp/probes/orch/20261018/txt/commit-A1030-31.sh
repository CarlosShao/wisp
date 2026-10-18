#!/bin/sh
set -e
cd "/d/work/workspace/projects plans/Wisp"
T=.scratch/wisp/probes/orch/20261018/txt
F338=$(ls .scratch/wisp/issues/338-*.md)
cat "$T/A1030.txt" >> docs/reports/pending-and-issues.md
cat "$T/A1031.txt" >> docs/reports/pending-and-issues.md
git add -- docs/reports/pending-and-issues.md "$F338" "$T/A1030.txt" "$T/A1031.txt" "$T/40ep-msg.txt" "$T/recount-322r2.sh" "$T/338r1-amend-01.txt" "$T/commit-A1030.sh" "$T/commit-A1030-31.sh"
echo "=== staged audit ==="
git diff --cached --name-only | grep -v -E '^(\.scratch/wisp/probes/orch/20261018/|\.scratch/wisp/issues/338-|docs/reports/pending-and-issues\.md)$' || echo "  only-mine"
git commit -F "$T/40ep-msg.txt" -- docs/reports/pending-and-issues.md "$F338" "$T/A1030.txt" "$T/A1031.txt" "$T/40ep-msg.txt" "$T/recount-322r2.sh" "$T/338r1-amend-01.txt" "$T/commit-A1030.sh" "$T/commit-A1030-31.sh"
echo "=== self-check ==="
git show --name-only --format='' HEAD | grep -v -E '^(\.scratch/wisp/probes/orch/20261018/|\.scratch/wisp/issues/338-|docs/reports/pending-and-issues\.md)$' || echo "  commit-clean: only mine"
git log -1 --format='%h'
echo "unpushed=$(git rev-list --count origin/dev..HEAD)"
