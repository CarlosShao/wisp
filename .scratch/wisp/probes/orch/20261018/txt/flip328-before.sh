#!/bin/sh
# 编排者翻框（票 328）：凭据＝非实现者终裁腿 328-v1 判语件 + 我自己复跑的决定性那一发
F=$(ls .scratch/wisp/issues/328-*.md)
echo "ticket=$F"
echo "--- AC box lines (line numbers) ---"
grep -n '^[[:space:]]*- \[' "$F"
echo "--- before ---"
echo "open=$(grep -c '^[[:space:]]*- \[ \]' "$F") done=$(grep -c '^[[:space:]]*- \[x\]' "$F")"
