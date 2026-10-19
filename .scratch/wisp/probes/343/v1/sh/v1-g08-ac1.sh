#!/bin/sh
# AC#1 完成判据（我自己现跑，⛔ 抄腿）：
#  ① 交件侧 git grep -n aEquiv HEAD -- internal/risk ＝ 0 命中（正控：起手侧 1 命中）
#  ② git diff -U0 里每一枚 ^[+-] 行都以 // 开头 ⇒ 非注释行变化 0
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 90
export LC_ALL=C.UTF-8
LOG=/d/tmp/wisp/20261019/343v1/logs
{
echo "== ① delivery-side grep (want 0 hits) =="
git grep -n aEquiv HEAD -- internal/risk > $LOG/_grep_head.txt 2>&1; rc=$?
echo "rc=$rc hits=$(wc -l < $LOG/_grep_head.txt)"
cat $LOG/_grep_head.txt
echo "-- positive control: same grep on start anchor 29d258f9 (want 1 hit, and it is a comment line) --"
git grep -n aEquiv 29d258f9 -- internal/risk > $LOG/_grep_start.txt 2>&1; rc=$?
echo "rc=$rc hits=$(wc -l < $LOG/_grep_start.txt)"
cat $LOG/_grep_start.txt
echo "-- func-shape grep both sides --"
git grep -c 'func aEquiv' HEAD -- internal/risk || true
git grep -c 'func aEquiv' 29d258f9 -- internal/risk || true
echo
echo "== ② the single .go commit: 31d4a476 numstat + -U0 +/- line classification =="
git show --numstat --format='%H %s' 31d4a476 -- internal/risk/blacklist.go
echo "-- all ^[+-] content lines from git diff -U0 (start-anchor parent -> impl commit) --"
git diff -U0 31d4a476^ 31d4a476 -- internal/risk/blacklist.go | grep -E '^[+-]' | grep -vE '^(\+\+\+|---)' > $LOG/_diff_lines.txt
tot=$(wc -l < $LOG/_diff_lines.txt)
noncomment=$(grep -vcE '^[+-][[:space:]]*//' $LOG/_diff_lines.txt || true)
printf 'changed_lines=%s comment_lines=%s noncomment_lines=%s\n' "$tot" "$((tot-noncomment))" "$noncomment"
echo "-- verbatim --"
cat -A $LOG/_diff_lines.txt | sed 's/\$$/|EOL/' | cut -c1-200
echo
echo "== whole-commit file list for 31d4a476 (must be 1 .go) =="
git show --name-only --format= 31d4a476
echo "== full-HEAD numstat of that commit (any other file touched?) =="
git diff --numstat 31d4a476^ 31d4a476
} > $LOG/g08-ac1-gauges.log 2>&1
echo DONE
