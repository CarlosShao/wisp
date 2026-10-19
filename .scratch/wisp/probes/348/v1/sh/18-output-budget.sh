#!/bin/sh
# 342-v2 输出预算尺：正文件枚数／每件行数／逐格行数／0 字节／.out／被禁码位 U+00A7（正控 HANDOVER）。
cd '/d/work/workspace/projects plans/Wisp' || exit 9
for F in docs/evidence/s1/342-v2-verdict.md docs/evidence/s1/348-v1-verdict.md; do
  [ -f "$F" ] || { echo "MISSING $F"; continue; }
  echo "== $F LINES=$(wc -l < "$F") BYTES=$(wc -c < "$F")"
  awk '
    /^- \*\*AC/ { if (name != "") printf "   cell %s -> %d lines\n", name, cnt; name = $0; sub(/^- \*\*/, "", name); sub(/ .*/, "", name); cnt = 0; next }
    { cnt++ }
    END { if (name != "") printf "   cell %s -> %d lines\n", name, cnt }
  ' "$F"
done
echo "== 写面尺 =="
find .scratch/wisp/probes/342/v2 .scratch/wisp/probes/348/v1 docs/evidence/s1/342-v2-verdict.md docs/evidence/s1/348-v1-verdict.md -type f -size 0 | sed 's/^/ZERO_BYTE /'
find .scratch/wisp/probes/342/v2 .scratch/wisp/probes/348/v1 -name '*.out' | sed 's/^/DOTOUT /'
echo "CTA7_HITS=$(grep -rlP "\xc2\xa7" .scratch/wisp/probes/342/v2 .scratch/wisp/probes/348/v1 docs/evidence/s1/342-v2-verdict.md docs/evidence/s1/348-v1-verdict.md 2>/dev/null | wc -l)"
echo "POSITIVE_CONTROL_HANDOVER_CTA7=$(grep -cP "\xc2\xa7" docs/reports/HANDOVER.md || true)"
echo "FILES_342_v2=$(find .scratch/wisp/probes/342/v2 -type f | wc -l)"
echo "FILES_348_v1=$(find .scratch/wisp/probes/348/v1 -type f | wc -l)"
