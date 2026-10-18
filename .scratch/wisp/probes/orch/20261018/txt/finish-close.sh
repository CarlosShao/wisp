#!/bin/sh
set -e
cd "/d/work/workspace/projects plans/Wisp"
T=.scratch/wisp/probes/orch/20261018/txt
N336=.scratch/wisp/issues/336-a-tier-fallback-branch-per-anchor-teeth-334-ac2-remainder-done.md
N334=.scratch/wisp/issues/334-anchor-missing-leaf-reason-names-another-path-and-ci-surface-reds-216-instrument-done.md
F343=.scratch/wisp/issues/343-blacklist-comment-names-aEquiv-a-symbol-that-does-not-exist-in-the-package.md

git add -- docs/reports/pending-and-issues.md "$N336" "$N334" "$F343" \
  "$T/A1037.txt" "$T/336-tail-receipt.md" "$T/close336-msg.txt" "$T/close336-334.sh" "$T/common-20261018.txt"

echo "=== staged audit ==="
git diff --cached --name-only | grep -vE '^(\.scratch/wisp/issues/|\.scratch/wisp/probes/orch/20261018/|docs/reports/pending-and-issues\.md)' || echo "  only-mine"
echo "=== rename now carries content? (R100 -> R0xx) ==="
git diff --cached --name-status | grep -E "^R"
echo "=== staged content checks ==="
echo "336 staged lines=$(git show ":$N336" | wc -l) open=$(git show ":$N336" | grep -cE '^[[:space:]]*- \[ \] \*\*AC#') done=$(git show ":$N336" | grep -cE '^[[:space:]]*- \[x\] \*\*AC#') status-closed=$(git show ":$N336" | sed -n '3p' | grep -c 已结案)"
echo "334 staged lines=$(git show ":$N334" | wc -l) open=$(git show ":$N334" | grep -cE '^[[:space:]]*- \[ \] \*\*AC#') done=$(git show ":$N334" | grep -cE '^[[:space:]]*- \[x\] \*\*AC#') status-closed=$(git show ":$N334" | sed -n '3p' | grep -c 已结案)"

git commit -F "$T/close336-msg.txt" -- docs/reports/pending-and-issues.md \
  .scratch/wisp/issues/336-a-tier-fallback-branch-per-anchor-teeth-334-ac2-remainder.md "$N336" \
  .scratch/wisp/issues/334-anchor-missing-leaf-reason-names-another-path-and-ci-surface-reds-216-instrument.md "$N334" \
  "$F343" \
  "$T/A1037.txt" "$T/336-tail-receipt.md" "$T/close336-msg.txt" "$T/close336-334.sh" "$T/common-20261018.txt" \
  2>&1 | tail -8
echo "=== after ==="
git log -1 --format='%h'
for n in 322 334 336; do F=$(ls .scratch/wisp/issues/${n}-*.md); echo "$n files-at-HEAD=$(git ls-tree HEAD .scratch/wisp/issues/ --name-only | grep -c "issues/${n}-") open=$(grep -cE '^[[:space:]]*- \[ \] \*\*AC#' "$F") done=$(grep -cE '^[[:space:]]*- \[x\] \*\*AC#' "$F")"; done
echo "unpushed=$(git rev-list --count origin/dev..HEAD)"
