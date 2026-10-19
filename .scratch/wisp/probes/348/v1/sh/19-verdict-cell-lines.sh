#!/bin/sh
# 判语件逐格尺：每格＝从 `- **AC#n` 起，到下一格标题／下一节标题前；⛔ 含段间空行。帽＝每件≤130 行、逐格≤8 行。
cd '/d/work/workspace/projects plans/Wisp' || exit 9
for F in docs/evidence/s1/342-v2-verdict.md docs/evidence/s1/348-v1-verdict.md; do
  echo "== $F LINES=$(wc -l < "$F") BYTES=$(wc -c < "$F")"
  awk '
    /^- \*\*AC#/ { if (cell != "") printf "   %s -> %d lines\n", cell, cnt; cell = substr($0, 5, 6); cnt = 1; insec = 1; next }
    /^## / { if (cell != "") { printf "   %s -> %d lines\n", cell, cnt; cell = "" } ; next }
    { if (cell != "") cnt++ }
    END { if (cell != "") printf "   %s -> %d lines\n", cell, cnt }
  ' "$F"
done
