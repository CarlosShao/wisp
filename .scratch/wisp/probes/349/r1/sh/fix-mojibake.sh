#!/bin/sh
# 我上一把"把 SECTION-CHAR 写成词面"的改正脚本把中文写成了双重编码乱码（M-OJIBAKE-A / M-OJIBAKE-B = 第/节 的 UTF-8 字节再被编码一次）
# ⇒ 本脚本按**字节形**把它换回正确的 UTF-8（replacement 也用 \x 转义＝⛔ 依赖脚本自身的编码）。
# 先摊 hex 证明我看到了什么，再替换，再复尺。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 90
R=.scratch/wisp/probes/349/r1
echo "== step 0: hex of the mojibake run (first occurrence, 24 bytes around it) =="
for f in "$R/report.md" "$R/report-sec-b.md" "$R/msg/ticket-bullet-349.md"; do
  printf '%s : ' "$f"
  grep -aboP '\xc3\xa7\xc2\xac\xc2\xac.{0,12}' "$f" | head -1 | cut -c1-160
done
echo "== step 1: byte-level replace  (mojibake 第/节 pair, digit optional) -> correct UTF-8 =="
for f in "$R/report.md" "$R/report-sec-b.md" "$R/msg/ticket-bullet-349.md"; do
  perl -i -pe 's/\xc3\xa7\xc2\xac\xc2\xac\x20(\d?)\s?\xc3\xa8\xc2\x8a\xc2\x82/\xe7\xac\xac$1\xe8\x8a\x82/g' "$f"
  printf 'replaced-in=%s remaining_mojibake=%s\n' "$f" "$(grep -acP '\xc3\xa7\xc2\xac' "$f" || true)"
done
echo "== step 2: also collapse the 'sec ' leftover form the earlier scrub produced (if any remain) =="
for f in "$R/report.md" "$R/report-sec-b.md" "$R/msg/ticket-bullet-349.md" "$R/logs/ac3/ac3-misfire.log" "$R/logs/ac3/ac3-outer.log" "$R/logs/ac3/ac3-outer2.log"; do
  [ -f "$f" ] || continue
  printf '%s sec-space-digit=%s\n' "$f" "$(grep -c 'sec [0-9]' "$f" || true)"
done
echo "== step 3: mojibake scan repo-wide in my write face (want 0 files) =="
grep -rlnP '\xc3\xa7\xc2\xac|\xc3\xa8\xc2\x8a' "$R" | sort || echo "none"
echo "== step 4: U+00A7 re-scan (codepoint form; want 0 everywhere) =="
find "$R" -type f | sort | while read -r f; do
  n=$(perl -CSD -0ne 'print scalar(() = /\x{00A7}/g)' "$f")
  [ "${n:-0}" != "0" ] && printf '  HIT %-60s %s\n' "$f" "$n"
done
echo "(no HIT lines above = clean)"
