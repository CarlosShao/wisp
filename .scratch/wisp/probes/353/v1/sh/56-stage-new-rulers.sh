#!/bin/sh
# 353-v1 R56 = 把三枚新尺（sh/53|54|55）复制进写面并逐枚量字节
cd '/d/work/workspace/projects plans/Wisp' || exit 9
D=.scratch/wisp/probes/353/v1/sh
cp -f /d/tmp/wisp/20261020/353v1/sh/53-find-ranges.sh "$D/" 2>/dev/null || true
cp -f /d/tmp/wisp/20261020/353v1/sh/54-alt-negatives.sh "$D/" 2>/dev/null || true
cp -f /d/tmp/wisp/20261020/353v1/sh/55-scan-stale.sh "$D/" 2>/dev/null || true
for f in "$D/53-find-ranges.sh" "$D/54-alt-negatives.sh" "$D/55-scan-stale.sh"; do
  printf '%-46s %s B\n' "$f" "$(wc -c < "$f" 2>/dev/null || echo MISSING)"
done
echo "sh/ 总枚数=$(ls -1 "$D" | wc -l)"
echo "0 字节枚数=$(find "$D" -type f -size 0 | wc -l)"
echo ".out 枚数=$(ls -1 "$D" | grep -c '\.out$' || true)"
