#!/bin/sh
# 交件前最后一把全写面卫生尺（三形对撞）：
#   形①＝码位 U+00A7（被禁字面字符）＝每枚件 want 0
#   形②＝双重编码乱码标记（我踩过的那形：第/节的 UTF-8 字节再被编码一次）＝want 0
#   形③＝零字节件／.out 扩展名＝want 0
# 正控＝同一把码位尺打 docs/reports/HANDOVER.md（非零）。
# ⛔ 本脚本⛔ 写出被禁字符本身（只写码位名），⛔ 写出乱码字面（只写它们的 hex 字节形）。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 90
R=.scratch/wisp/probes/349/r1
echo "== form1: codepoint U+00A7 per file (want all 0) =="
find "$R" -type f | sort | while read -r f; do
  n=$(perl -CSD -0ne 'print scalar(() = /\x{00A7}/g)' "$f")
  [ "${n:-0}" != "0" ] && printf '  HIT %-64s %s\n' "$f" "$n"
done
echo "  (no HIT lines = 0 hits)"
echo "== form2: double-encoded markers, byte form c3 a7 c2 ac / c3 a8 c2 8a (want all 0) =="
find "$R" -type f | sort | while read -r f; do
  n=$(grep -acP '\xc3\xa7\xc2\xac|\xc3\xa8\xc2\x8a' "$f" || true)
  [ "${n:-0}" != "0" ] && printf '  HIT %-64s %s\n' "$f" "$n"
done
echo "  (no HIT lines = 0 hits)"
echo "== form2b: replacement-character / lone 0xEF 0xBF 0xBD (want 0) =="
find "$R" -type f | sort | while read -r f; do
  n=$(grep -acP '\xef\xbf\xbd' "$f" || true)
  [ "${n:-0}" != "0" ] && printf '  HIT %-64s %s\n' "$f" "$n"
done
echo "  (no HIT lines = 0 hits)"
echo "== form3: zero-byte files / .out extension =="
find "$R" -type f -size 0 | sort
find "$R" -type f -name '*.out' | sort
echo "  (nothing listed = rule satisfied)"
echo "== positive controls =="
printf '  HANDOVER.md U+00A7 = %s (want nonzero)\n' "$(perl -CSD -0ne 'print scalar(() = /\x{00A7}/g)' docs/reports/HANDOVER.md)"
printf '  own-file byte-count nonzero check: files=%s nonzero=%s\n' \
  "$(find "$R" -type f | wc -l)" "$(find "$R" -type f -size +0c | wc -l)"
echo "== every .go file I own? (report face only, ⛔ any .go in probes) =="
find "$R" -type f -name '*.go' | sort
echo "  (the census instrument ships as *.go.txt = deliberately NOT a .go, ⛔ enter any build)"
echo "== report line counts vs the 160-line hat =="
wc -l "$R/report.md" "$R/report-sec-b.md"
echo "== ticket box ruler (want 7 unchecked / 0 checked) =="
T=.scratch/wisp/issues/349-blacklist-provenance-comments-name-undefined-symbols-fires-taskscope-cleanup.md
printf 'unchecked=%s checked=%s\n' "$(grep -c '^- \[ \] \*\*AC#' "$T")" "$(grep -c '^- \[x\]' "$T" || true)"
