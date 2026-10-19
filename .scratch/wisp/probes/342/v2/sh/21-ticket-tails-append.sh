#!/bin/sh
# 342-v2 票尾追加尺：两票尾**各一条**只追加 bullet（⛔ 抹原句、⛔ 翻框、⛔ 改票名/Status/-done）。
# 纯追加证明＝git diff --numstat 删除行数 0 ＋ 框计数与起手逐字同 ＋ before/after sha256 前缀不变（只增长）。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 9
B=/d/tmp/wisp/20261019/342v2
T342='.scratch/wisp/issues/342-privatefile-open-follows-a-swapped-ancestor-on-both-platforms-no-ticket-owns-it.md'
T348='.scratch/wisp/issues/348-refusal-path-remove-follows-an-ancestor-link-and-deletes-a-foreign-file.md'
for T in "$T342" "$T348"; do
  echo "BEFORE $T LINES=$(wc -l < "$T") BYTES=$(wc -c < "$T") SHA16=$(sha256sum "$T" | cut -c1-16)"
  tail -c 1 "$T" | od -An -c | tr -d ' ' | sed 's/^/  last_byte=/'
done
# 正控尺＝票框计数 before
echo "BEFORE T342 x=$(grep -c '^- \[x\]' "$T342") open=$(grep -c '^- \[ \]' "$T342")"
echo "BEFORE T348 x=$(grep -c '^- \[x\]' "$T348") open=$(grep -c '^- \[ \]' "$T348")"
cat "$B/bullet-342.txt" >> "$T342"
cat "$B/bullet-348.txt" >> "$T348"
for T in "$T342" "$T348"; do
  echo "AFTER $T LINES=$(wc -l < "$T") BYTES=$(wc -c < "$T")"
done
echo '### 纯追加证明（尺＝git diff --numstat；删除列须=0）'
git diff --numstat -- "$T342" "$T348"
echo "DELETED_LINES_TOTAL=$(git diff --numstat -- "$T342" "$T348" | awk '{s+=$2} END{print s+0}')"
echo '### 逐字未动的旧行（前 61/47 行 before⇄after 全同＝diff 里只该有 + 行）'
git diff -- "$T342" "$T348" | grep -c '^-[^-]' || true
echo "AFTER T342 x=$(grep -c '^- \[x\]' "$T342") open=$(grep -c '^- \[ \]' "$T342")"
echo "AFTER T348 x=$(grep -c '^- \[x\]' "$T348") open=$(grep -c '^- \[ \]' "$T348")"
echo '### 票名/Status 未动（第一行＋Status 行 before⇄after 同串）'
head -1 "$T342"; sed -n '3p' "$T342" | cut -c1-60
head -1 "$T348"; sed -n '3p' "$T348" | cut -c1-60
