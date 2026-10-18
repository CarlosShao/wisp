#!/bin/sh
set -e
cd "/d/work/workspace/projects plans/Wisp"
T=.scratch/wisp/probes/orch/20261018/txt

# 1) 台账追加
cat "$T/A1037.txt" >> docs/reports/pending-and-issues.md

# 2) 票 336 尾追加编排者收件 bullet（⛔ 改判据原文；336-v1 自己那条 bullet 已在盘上）
F336=$(ls .scratch/wisp/issues/336-*.md)
BEFORE336=$(wc -l < "$F336")
cat "$T/336-tail-receipt.md" >> "$F336"
echo "336 lines $BEFORE336 -> $(wc -l < "$F336")"

# 3) 两枚改名（各一笔带新旧两枚 pathspec）
F334=$(ls .scratch/wisp/issues/334-*.md)
N336=$(echo "$F336" | sed 's/\.md$/-done.md/')
N334=$(echo "$F334" | sed 's/\.md$/-done.md/')
git mv "$F336" "$N336"
git mv "$F334" "$N334"

git add -- docs/reports/pending-and-issues.md "$F336" "$N336" "$F334" "$N334" \
  "$T/A1037.txt" "$T/336-tail-receipt.md" "$T/close336-msg.txt" "$T/close336-334.sh" \
  .scratch/wisp/issues/343-blacklist-comment-names-aEquiv-a-symbol-that-does-not-exist-in-the-package.md \
  .scratch/wisp/issues/342-privatefile-open-follows-a-swapped-ancestor-on-both-platforms-no-ticket-owns-it.md \
  "$T/common-20261018.txt"
echo "=== staged audit ==="
git diff --cached --name-only | grep -vE '^(\.scratch/wisp/(issues|probes/orch/20261018/)|docs/reports/pending-and-issues\.md)' || echo "  only-mine"

git commit -F "$T/close336-msg.txt" -- docs/reports/pending-and-issues.md "$F336" "$N336" "$F334" "$N334" \
  "$T/A1037.txt" "$T/336-tail-receipt.md" "$T/close336-msg.txt" "$T/close336-334.sh" \
  .scratch/wisp/issues/343-blacklist-comment-names-aEquiv-a-symbol-that-does-not-exist-in-the-package.md \
  .scratch/wisp/issues/342-privatefile-open-follows-a-swapped-ancestor-on-both-platforms-no-ticket-owns-it.md \
  "$T/common-20261018.txt" 2>&1 | tail -8
echo "=== after: uniqueness + boxes ==="
git log -1 --format='%h'
for n in 322 334 336; do F=$(ls .scratch/wisp/issues/${n}-*.md); echo "$n files=$(git ls-tree HEAD .scratch/wisp/issues/ --name-only | grep -c "issues/${n}-") open=$(grep -cE '^[[:space:]]*- \[ \] \*\*AC#' "$F") done=$(grep -cE '^[[:space:]]*- \[x\] \*\*AC#' "$F")"; done
echo "unpushed=$(git rev-list --count origin/dev..HEAD)"
