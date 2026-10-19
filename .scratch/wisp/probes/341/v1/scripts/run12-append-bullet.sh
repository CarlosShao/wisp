#!/usr/bin/env bash
# 341-v1: append exactly ONE bullet to the tail of ticket 341 (append-only, ⛔ 改原句/票名/Status/翻框).
R="/d/work/workspace/projects plans/Wisp"; cd "$R" || exit 9
T=".scratch/wisp/issues/341-four-hand-edit-shapes-still-unnailed-and-two-226-sentences-overclaim.md"
B="/d/tmp/wisp/20261018/341v1/bullet-341-append.txt"

L0=$(wc -l < "$T"); echo "lines before = $L0"
echo "unchecked before = $(grep -c '^- \[ \] \*\*AC#' "$T" || true)  checked before = $(grep -c '^- \[x\]' "$T" || true)"
echo "last line before (tail 60 chars):"; tail -1 "$T" | tail -c 60; echo

cp -f "$T" "/d/tmp/wisp/20261018/341v1/ticket-before.md"
cat "$B" >> "$T"

echo "lines after  = $(wc -l < "$T")"
echo "unchecked after = $(grep -c '^- \[ \] \*\*AC#' "$T" || true)  checked after = $(grep -c '^- \[x\]' "$T" || true)"
echo "diff vs before (want one added line, zero removed):"
diff "/d/tmp/wisp/20261018/341v1/ticket-before.md" "$T" | sed 's/^/   /' | cut -c1-120 | head -6
echo "added lines   = $(diff "/d/tmp/wisp/20261018/341v1/ticket-before.md" "$T" | grep -c '^> ' || true)"
echo "removed lines = $(diff "/d/tmp/wisp/20261018/341v1/ticket-before.md" "$T" | grep -c '^< ' || true)"
echo "my bullet count = $(grep -c '终裁腿 341-v1 判语' "$T" || true)   leg bullet count = $(grep -c '落地腿 341-r1 交件' "$T" || true)"
echo "title untouched = $(head -1 "$T" | cut -c1-30)..."
echo "status untouched = $(sed -n '3p' "$T" | cut -c1-30)..."
