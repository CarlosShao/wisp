#!/bin/sh
set -e
cd "/d/work/workspace/projects plans/Wisp"
OLD=$(ls .scratch/wisp/issues/328-*.md)
NEW=$(echo "$OLD" | sed 's/\.md$/-done.md/')
echo "old=$OLD"
echo "new=$NEW"
echo "name-len new=$(basename "$NEW" | wc -c)"
# 1) 票尾追加一条 bullet（只追加）
cat .scratch/wisp/probes/orch/20261018/txt/328-tail-bullet.md >> "$OLD"
echo "after append: open=$(grep -c '^[[:space:]]*- \[ \]' "$OLD") done=$(grep -c '^[[:space:]]*- \[x\]' "$OLD") lines=$(wc -l < "$OLD")"
# 2) 改名：一笔带新旧两枚 pathspec
git mv "$OLD" "$NEW"
git add -- "$OLD" "$NEW" .scratch/wisp/issues/340-stacked-queue-rewrite-of-a-card-corr-has-no-real-producer-test.md .scratch/wisp/probes/orch/20261018/txt/328-tail-bullet.md
echo "=== ls-tree uniqueness (must be exactly 1) ==="
git ls-tree HEAD .scratch/wisp/issues/ --name-only | grep -c '^\.scratch/wisp/issues/328-' || true
echo "=== staged now ==="
git diff --cached --name-status
