#!/bin/sh
# 编排者翻框（票 328）：凭据＝非实现者 328-v1 判语件（AC#2/AC#3/AC#4 三格判"过"）＋我自己复跑的决定性那一发
F=$(ls .scratch/wisp/issues/328-*.md)
echo "ticket=$F"
echo "before: open=$(grep -c '^[[:space:]]*- \[ \]' "$F") done=$(grep -c '^[[:space:]]*- \[x\]' "$F")"
# 只翻判语点名的三行：33=AC#2、34=AC#3、35=AC#4（判据区射程，⛔ 碰 50/59/61 那些 bullet 行）
sed -i '33,35s/^- \[ \]/- [x]/' "$F"
echo "after : open=$(grep -c '^[[:space:]]*- \[ \]' "$F") done=$(grep -c '^[[:space:]]*- \[x\]' "$F")"
echo "--- numstat ---"
git diff --numstat -- "$F"
echo "--- what changed (must be exactly 3 lines, box marks only) ---"
git diff -U0 -- "$F" | grep -E '^[+-]' | grep -v '^[+-][+-]' | cut -c1-40
