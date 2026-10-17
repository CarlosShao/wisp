#!/usr/bin/env bash
# 编排者翻票 327 的 AC 框（机械 sed，按行号＋框形定位；⛔ 改判据文字一字）
set -eu
cd "/d/work/workspace/projects plans/Wisp"
F=".scratch/wisp/issues/327-display-boundary-washes-cc-only-so-cf-invisible-chars-still-pass-through.md"
echo "before: unchecked=$(grep -cE '^[[:space:]]*- \[ \]' "$F") checked=$(grep -cE '^[[:space:]]*- \[x\]' "$F")"
cp "$F" "/d/work/workspace/wisp-327v1-orch-094846/327.before-flip.md"
# 行号 25..29 ＝ AC#0 AC#1 AC#2 AC#3 AC#4；行号 30 ＝ AC#5（只在颜色到手后由第二个参数翻）
for n in 25 26 27 28 29; do
  line="$(sed -n "${n}p" "$F")"
  case "$line" in
    "- [ ] **"*) sed -i "${n}s/^- \[ \]/- [x]/" "$F" ;;
    *) echo "SKIP_LINE_$n_not_a_box: $line" ; exit 1 ;;
  esac
done
if [ "${1:-}" = "with-ac5" ]; then
  line="$(sed -n 30p "$F")"
  case "$line" in
    "- [ ] **"*) sed -i "30s/^- \[ \]/- [x]/" "$F" ;;
    *) echo "SKIP_LINE_30_not_a_box: $line"; exit 1 ;;
  esac
fi
echo "after: unchecked=$(grep -cE '^[[:space:]]*- \[ \]' "$F") checked=$(grep -cE '^[[:space:]]*- \[x\]' "$F")"
echo "=== diff 只见框形？ ==="
diff "/d/work/workspace/wisp-327v1-orch-094846/327.before-flip.md" "$F" | grep -E '^[<>]' | cut -c1-60
