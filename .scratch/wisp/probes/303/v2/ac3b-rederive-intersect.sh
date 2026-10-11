#!/bin/sh
# 303-v2 — AC#3b 需求③ 自取证：从腿件 303-v1 的原始逐发红名册 names-*.txt 自己重算交集作差
# ⛔ 不跑 go test（编排者车道）；本脚本只对已落盘的 stdout 名册文本做 comm/sort（非编译面）
set -u
SRC=".scratch/wisp/probes/303/v1/logs"
OUT=".scratch/wisp/probes/303/v2/logs"

# 只取顶层 FAIL 名（剥前缀），逐枚 sort
for f in before-1 before-2 after-1 after-2; do
  grep '^--- FAIL: ' "$SRC/names-$f.txt" | sed 's/^--- FAIL: //' | sort -u > "$OUT/re-$f-fail.txt"
done

# 我自己重算的改前交集 / 改后交集
comm -12 "$OUT/re-before-1-fail.txt" "$OUT/re-before-2-fail.txt" > "$OUT/re-before-intersect.txt"
comm -12 "$OUT/re-after-1-fail.txt"  "$OUT/re-after-2-fail.txt"  > "$OUT/re-after-intersect.txt"

# 新增红 = 改后交集 \ 改前交集
comm -13 "$OUT/re-before-intersect.txt" "$OUT/re-after-intersect.txt" > "$OUT/re-new-red.txt"
# 被修好 = 改前交集 \ 改后交集
comm -23 "$OUT/re-before-intersect.txt" "$OUT/re-after-intersect.txt" > "$OUT/re-fixed-red.txt"

# 单发红（两发并集 \ 交集）= 交集尺会藏住的那一枚，必须具名报出
cat "$OUT/re-after-1-fail.txt" "$OUT/re-after-2-fail.txt" | sort -u > "$OUT/re-after-union.txt"
comm -23 "$OUT/re-after-union.txt" "$OUT/re-after-intersect.txt" > "$OUT/re-after-single-shot-red.txt"

echo "re-before-intersect=$(wc -l < "$OUT/re-before-intersect.txt")"
echo "re-after-intersect=$(wc -l < "$OUT/re-after-intersect.txt")"
echo "re-new-red=$(wc -l < "$OUT/re-new-red.txt")"
echo "re-fixed-red=$(wc -l < "$OUT/re-fixed-red.txt")"
echo "re-after-single-shot-red=$(wc -l < "$OUT/re-after-single-shot-red.txt")"
echo "--- re-new-red 内容 ---"; cat "$OUT/re-new-red.txt"
echo "--- re-fixed-red 内容 ---"; cat "$OUT/re-fixed-red.txt"
echo "--- re-after-single-shot-red 内容 ---"; cat "$OUT/re-after-single-shot-red.txt"
echo "rc=$?"
