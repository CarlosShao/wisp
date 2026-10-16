#!/usr/bin/env bash
# 编排者翻框（机械 sed，只改 `- [ ]` → `- [x]`，逐枚指名 AC#0..AC#5，⛔ 改判据文字一字）
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 9
T=$(ls .scratch/wisp/issues/326-*.md | head -1)
echo "T=$T"
尺() { echo "before $(basename "$1"): unchecked=$(grep -c '^[[:space:]]*- \[ \]' "$1") checked=$(grep -c '^[[:space:]]*- \[x\]' "$1")"; }
尺 "$T"
for n in 0 1 2 3 4 5; do
  sed -i "s/^- \[ \] \*\*AC#${n} /- [x] **AC#${n} /" "$T"
  sed -i "s/^- \[ \] \*\*AC#${n}\*\*/- [x] **AC#${n}**/" "$T"
done
echo "--- 翻后 ---"
尺 "$T"
echo "--- diff 只应见 [ ]→[x] ---"
git diff --unified=0 -- "$T" | grep -E '^[-+]- \[' | cut -c1-90
echo "numstat=$(git diff --numstat -- "$T" | cat)"
