#!/usr/bin/env bash
# 编排者翻框（票 216 的 AC#2/AC#3/AC#4；机械 sed，⛔ 改判据文字一字）
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 9
T=$(ls .scratch/wisp/issues/216-*.md | head -1)
echo "T=$T"
echo "before: unchecked=$(grep -c '^[[:space:]]*- \[ \]' "$T") checked=$(grep -c '^[[:space:]]*- \[x\]' "$T")"
for n in 2 3 4; do
  sed -i "s/^- \[ \] \*\*AC#${n}\*\*/- [x] **AC#${n}**/" "$T"
  # AC#3 的票面写法是 `**AC#3（正控，必须有）**`，同一种翻法、另一种形状
  sed -i "s/^- \[ \] \*\*AC#${n}（/- [x] **AC#${n}（/" "$T"
done
echo "after : unchecked=$(grep -c '^[[:space:]]*- \[ \]' "$T") checked=$(grep -c '^[[:space:]]*- \[x\]' "$T")"
echo "--- 逐枚 diff（必须只见 [ ] -> [x]）"
git diff --unified=0 -- "$T" | grep -E '^[-+]- \[' 
