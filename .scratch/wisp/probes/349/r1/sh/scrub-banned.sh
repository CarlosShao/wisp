#!/bin/sh
# 被禁字面字符（码位 U+00A7＝分节符）在我写面（probes/349/r1/**＋票尾 bullet）里一律清成
# 词面 "sec "（⛔ 写出被禁字符本身；⚠ 与主题常用间隔号 U+00B7 是两枚不同字符、本脚本⛔ 动 U+00B7）。
# ★产码 .go 一字不碰：我已另用 check-banned-in-added.sh 证 起手 blob ⇄ 交件 blob 该码位计数 delta=0
#   （.go 里那些 sec 是别腿/历史存量，⛔ 我引入、⛔ 我顺手改）。
# ⛔ 删任何件（临时件只建不删）＝就地改写同一枚件。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 90
R=.scratch/wisp/probes/349/r1
echo "== before =="
find "$R" -type f | sort | while read -r f; do
  n=$(perl -CSD -0ne 'print scalar(() = /\x{00A7}/g)' "$f")
  [ "${n:-0}" != "0" ] && printf '  %-64s U00A7=%s\n' "$f" "$n"
done
echo "== scrub U+00A7 -> 'sec ' in every probe text file (⛔ 动 .go、⛔ 动 U+00B7) =="
find "$R" -type f | sort | while read -r f; do
  perl -CSD -i -pe 's/\x{00A7}\s*/sec /g' "$f"
done
echo "== after (want none listed) =="
find "$R" -type f | sort | while read -r f; do
  n=$(perl -CSD -0ne 'print scalar(() = /\x{00A7}/g)' "$f")
  [ "${n:-0}" != "0" ] && printf '  REMAIN %-60s U00A7=%s\n' "$f" "$n"
done
echo "(nothing listed above = scrubbed)"
echo "== re-count zero-byte + .out (want none) =="
find "$R" -type f -size 0 | sort
find "$R" -type f -name '*.out' | sort
echo "== U+00B7 still intact? (positive control that I did NOT over-scrub the wrong char) =="
perl -CSD -0ne 'print scalar(() = /\x{00B7}/g)' "$R/report.md"
echo ""
