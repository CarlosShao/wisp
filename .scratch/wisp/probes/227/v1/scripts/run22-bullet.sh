#!/bin/bash
# 227-v1: append exactly ONE bullet to ticket 227 (additions only; 六格⛔ 翻框).
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 2
t='.scratch/wisp/issues/227-guarded-write-delete-shapes-and-migration-on-read-are-unasserted.md'
out='.scratch/wisp/probes/227/v1/logs/23-ticket-bullet-append.txt'
cp "$t" /d/tmp/wisp/20261018/227v1/ticket-before.txt
{
  echo "== before =="
  echo "lines=$(wc -l < "$t")"
  echo "AC-unchecked=$(grep -c '^[[:space:]]*- \[ \] \*\*AC#' "$t" || true) AC-checked=$(grep -c '^[[:space:]]*- \[x\] \*\*AC#' "$t" || true)"
  echo "bullet-count-227-r1=$(grep -c '^- \*\*2026-10-18 落地腿 227-r1' "$t" || true)"
  echo "bullet-count-227-v1=$(grep -c '^- \*\*2026-10-18 非实现者终裁腿 227-v1' "$t" || true) (want 0 before)"
} > "$out" 2>&1
cat /d/tmp/wisp/20261018/227v1/msgs/bullet227.txt >> "$t"
{
  echo "== after the append =="
  echo "lines=$(wc -l < "$t")"
  echo "AC-unchecked=$(grep -c '^[[:space:]]*- \[ \] \*\*AC#' "$t" || true) AC-checked=$(grep -c '^[[:space:]]*- \[x\] \*\*AC#' "$t" || true)"
  echo "bullet-count-227-r1=$(grep -c '^- \*\*2026-10-18 落地腿 227-r1' "$t" || true)"
  echo "bullet-count-227-v1=$(grep -c '^- \*\*2026-10-18 非实现者终裁腿 227-v1' "$t" || true) (want 1)"
  echo "== diff numstat (want 2 0 => pure addition of one bullet line + its blank line, zero deletions) =="
  git diff --numstat -- "$t"
  echo "== git diff word-diff sanity: no removed content =="
  git diff -- "$t" | grep -c '^-[^-]' || true
} >> "$out" 2>&1
echo "appended; see $out"
