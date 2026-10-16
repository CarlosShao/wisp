#!/usr/bin/env bash
# 编排者翻框（机械 sed，只改 `- [ ]` → `- [x]`，逐枚指名 AC，⛔ 改判据文字一字）
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 9
T120=$(ls .scratch/wisp/issues/120-*.md | head -1)
T270=$(ls .scratch/wisp/issues/270-*.md | head -1)
echo "T120=$T120"; echo "T270=$T270"

尺() {
  echo "before $(basename "$1"): unchecked=$(grep -c '^[[:space:]]*- \[ \]' "$1") checked=$(grep -c '^[[:space:]]*- \[x\]' "$1")"
}
尺 "$T120"; 尺 "$T270"

# 票 120：五格全翻（凭据＝非实现者腿 120-v2 ＋我 11:5x 逐字重放它三把尺全部对上）
for n in 1 2 3 4 5; do
  sed -i "s/^- \[ \] \*\*AC#${n}\*\*/- [x] **AC#${n}**/" "$T120"
  sed -i "s/^- \[ \] \*\*AC#${n} \*\*/- [x] **AC#${n} **/" "$T120"
done
# 票 270：AC#1..AC#5 翻（AC#0 ⛔ 翻：裁决腿具名不裁、无凭据）
for n in 1 2 3 4 5; do
  sed -i "s/^- \[ \] \*\*AC#${n} /- [x] **AC#${n} /" "$T270"
  sed -i "s/^- \[ \] \*\*AC#${n}\*\*/- [x] **AC#${n}**/" "$T270"
done

echo "--- 翻后 ---"
尺 "$T120"; 尺 "$T270"
echo "--- 逐枚点名改了哪几行（diff 只应见 [ ]→[x]）---"
git diff --unified=0 -- "$T120" "$T270" | grep -E '^[-+]#|^[-+]- \[' | cut -c1-70
echo "changed_lines=$(git diff --numstat -- "$T120" "$T270" | cat)"
