#!/bin/sh
# 票 349 尾追加**一条** bullet（AC#6 只许这一枚）。守门＝落刀前现量"我这枚 bullet 的锚＝0 枚"
# （防重复追加＝编排者今天自己踩过"翻倍套娃"那形），落刀后再量一次＝1 枚。
# ⛔ 改票名／Status／任何 AC 框／⛔ 抹原句：只在文件尾加一行。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 90
T=.scratch/wisp/issues/349-blacklist-provenance-comments-name-undefined-symbols-fires-taskscope-cleanup.md
B=.scratch/wisp/probes/349/r1/msg/ticket-bullet-349.md
echo "== before: my-bullet anchor hits (want 0) =="
grep -c '落地腿 `349-r1` 交件' "$T" || true
echo "== ticket lines/bytes before =="
wc -l -c "$T"
cat "$B" >> "$T"
echo "== after: my-bullet anchor hits (want 1) =="
grep -c '落地腿 `349-r1` 交件' "$T" || true
echo "== ticket lines/bytes after (want +2 lines = 1 blank + 1 bullet) =="
wc -l -c "$T"
echo "== AC box ruler after (want unchecked=7, checked=0) =="
printf 'unchecked=%s checked=%s\n' "$(grep -c '^- \[ \] \*\*AC#' "$T")" "$(grep -c '^- \[x\]' "$T" || true)"
echo "== Status line untouched? (diff should show ONLY appended lines) =="
git diff --numstat -- "$T"
git diff -U0 -- "$T" | grep -E '^[+-]' | grep -vE '^(\+\+\+|---)' | cut -c1-60 | head -5
