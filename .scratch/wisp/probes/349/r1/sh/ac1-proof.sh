#!/bin/sh
# 349-r1 AC#1 逐枚交付尺（对一枚被动的 .go 文件跑全套，⛔ 只看"diff 像是注释"这句话）：
#   1) git diff -U0：`^[+-]` 行逐枚以 `//` 开头；非注释行数（want 0）
#   2) 起手 blob（ANCHOR）⇄ 交件 blob：GS1／GS2／GB 三把剥注释尺 sha256 逐枚相同＋cmp
#   3) 该枚被改的名字在本包 git grep 终态（带正控）
# 用法：sh ac1-proof.sh <anchor> <file> <dead-name> <positive-control-name> <logname>
set -u
REPO='/d/work/workspace/projects plans/Wisp'
cd "$REPO" || exit 90
ANCHOR="$1"; F="$2"; DEAD="$3"; POS="$4"; LOGNAME="$5"
T=/d/tmp/wisp/20261019/349r1
LOG="$T/logs/$LOGNAME"
mkdir -p "$T/logs" "$T/before" "$T/after" || exit 91
B=$(basename "$F")
git cat-file blob "$ANCHOR:$F" > "$T/before/$B" || exit 92
cp -f "$F" "$T/after/$B" || exit 93
{
echo "############ AC#1 proof: file=$F anchor=$ANCHOR ############"
echo "== 0. blob hashes: before / after =="
git rev-parse "$ANCHOR:$F"; git hash-object "$F"
echo "== 1. git diff -U0 changed lines (each must start with //) =="
git diff -U0 -- "$F" | grep -E '^[+-]' | grep -vE '^(\+\+\+|---)'
echo "-- count of changed lines: total / non-comment (want non-comment = 0) --"
git diff -U0 -- "$F" | grep -E '^[+-]' | grep -vE '^(\+\+\+|---)' | wc -l
git diff -U0 -- "$F" | grep -E '^[+-]' | grep -vE '^(\+\+\+|---)' | grep -vcE '^[+-][[:space:]]*//' || true
echo "== 1b. git numstat (added/deleted lines, this file) =="
git diff --numstat "$ANCHOR" -- "$F"
echo "== 2. THREE stripped-executable gauges, BOTH sides =="
for side in before after; do
  sh .scratch/wisp/probes/349/r1/sh/strip-sha.sh "$T/$side/$B" "$side"
done
sed -e 's://.*::' "$T/before/$B" | tr -d ' \t\r' | sed '/^$/d' > "$T/GS1-before.txt"
sed -e 's://.*::' "$T/after/$B"  | tr -d ' \t\r' | sed '/^$/d' > "$T/GS1-after.txt"
cmp "$T/GS1-before.txt" "$T/GS1-after.txt" && echo "GS1_CMP=IDENTICAL" || echo "GS1_CMP=DIFFERENT"
grep -v '^[[:space:]]*//' "$T/before/$B" > "$T/GB-before.txt"
grep -v '^[[:space:]]*//' "$T/after/$B"  > "$T/GB-after.txt"
cmp "$T/GB-before.txt" "$T/GB-after.txt" && echo "GB_CMP=IDENTICAL" || echo "GB_CMP=DIFFERENT"
echo "== 3. the dead name in internal/risk, at ANCHOR vs on disk (delivery) =="
echo "-- ANCHOR side (expect the old bare-name hit to exist there) --"
git grep -n -F "$DEAD" "$ANCHOR" -- internal/risk || echo "ZERO-HITS-at-anchor"
echo "-- POSITIVE CONTROL (a name that must hit): $POS --"
git grep -c -F "$POS" "$ANCHOR" -- internal/risk || echo "ZERO-HITS-CONTROL-BROKEN"
echo "== 3b. raw-byte and codepoint hygiene for the on-disk file =="
printf 'CR_count=%s bytes=%s lines=%s maxlen=%s\n' \
  "$(tr -cd '\r' < "$F" | wc -c)" "$(wc -c < "$F")" "$(wc -l < "$F")" \
  "$(awk '{ if (length > m) m = length } END { print m }' "$F")"
} > "$LOG" 2>&1
echo "wrote=$LOG rc=$?"
cat "$LOG"
