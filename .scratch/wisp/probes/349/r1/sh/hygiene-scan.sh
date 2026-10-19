#!/bin/sh
# 同一枚被禁字符尺的两形对撞（U+00A7）：
#   形①＝字节形（⛔ 解码 UTF-8，把 0xA7 当字节数）＝会误报（三字节序列的续字节可能是 0xA7）
#   形②＝Unicode 形（-C 解码后数码位）＝真尺
# 本仓"被禁字面字符"的判据＝码位 U+00A7 ⇒ 认形②；形①具名摊出来防"看着 0 其实没解码"。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 90
echo "== form2 (UNICODE codepoint count, the gauge the ticket means) =="
find .scratch/wisp/probes/349/r1 -type f | sort | while read -r f; do
  n=$(perl -CSD -0ne 'print scalar(() = /\x{00A7}/g)' "$f")
  printf '%-72s U00A7=%s\n' "$f" "${n:-0}"
done
echo "== form1 (raw BYTE 0xA7 count; informational, can over-count UTF-8 continuations) =="
find .scratch/wisp/probes/349/r1 -type f | sort | while read -r f; do
  printf '%-72s byte_A7=%s\n' "$f" "$(tr -cd '\247' < "$f" | wc -c)"
done
echo "== POSITIVE CONTROL (form2 on docs/reports/HANDOVER.md, want nonzero) =="
perl -CSD -0ne 'print scalar(() = /\x{00A7}/g)' docs/reports/HANDOVER.md
echo ""
