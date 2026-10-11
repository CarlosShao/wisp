#!/usr/bin/env bash
# 308-r2 verification #1: comment-only proof, re-run by the orchestrator.
# Ruler = strip whole-line comments from both the baseline blob (202ed4a5) and the
# current worktree file, then cmp. If they are identical, no executable line moved.
set -uo pipefail
cd "$(git rev-parse --show-toplevel)"
BASE=202ed4a5
OUT=.scratch/wisp/probes/orch/20261011/logs/orch-A848-308r2-strip.txt
TMP=$(mktemp -d)
mkdir -p "$(dirname "$OUT")"
: > "$OUT"
FILES="internal/risk/pathresolver_anchor_spelling_windows_test.go internal/risk/pathresolver_volumeless_spelling_308_test.go"

for f in $FILES; do
  b="$TMP/$(basename "$f").base"
  w="$TMP/$(basename "$f").work"
  git show "$BASE:$f" > "$b" 2>>"$OUT"
  cp "$f" "$w"
  sb="$TMP/$(basename "$f").base.strip"
  sw="$TMP/$(basename "$f").work.strip"
  # strip: drop lines whose first non-space chars are // ; also drop blank-only lines
  grep -vE '^[[:space:]]*//' "$b" | grep -vE '^[[:space:]]*$' > "$sb"
  grep -vE '^[[:space:]]*//' "$w" | grep -vE '^[[:space:]]*$' > "$sw"
  echo "=== $f ===" >> "$OUT"
  echo "lines base_total=$(wc -l < "$b") work_total=$(wc -l < "$w")" >> "$OUT"
  echo "lines base_code=$(wc -l < "$sb") work_code=$(wc -l < "$sw")" >> "$OUT"
  echo "blockcomment_hits_base=$(grep -c '/\*' "$b" || true) work=$(grep -c '/\*' "$w" || true)" >> "$OUT"
  echo "cr_count base=$(tr -dc '\r' < "$b" | wc -c) work=$(tr -dc '\r' < "$w" | wc -c)" >> "$OUT"
  if cmp -s "$sb" "$sw"; then echo "STRIPPED_CMP=IDENTICAL" >> "$OUT"; else echo "STRIPPED_CMP=DIFF" >> "$OUT"; diff "$sb" "$sw" | head -20 >> "$OUT"; fi
  echo "md5 base_code=$(md5sum "$sb" | cut -c1-12) work_code=$(md5sum "$sw" | cut -c1-12)" >> "$OUT"
done

echo "=== production files touched by this leg? (must be empty) ===" >> "$OUT"
git diff --name-only "$BASE" HEAD -- internal/ cmd/ tools/ >> "$OUT"
echo "=== full changed-file list base..HEAD ===" >> "$OUT"
git diff --name-only "$BASE" HEAD >> "$OUT"
echo "=== sha256 of the two production files at HEAD vs baseline ===" >> "$OUT"
for p in internal/risk/pathresolver.go internal/tools/volume_windows.go; do
  git show "$BASE:$p" > "$TMP/p.base"
  git show "HEAD:$p" > "$TMP/p.head"
  echo "$p base=$(sha256sum "$TMP/p.base" | cut -c1-8) head=$(sha256sum "$TMP/p.head" | cut -c1-8) $(cmp -s "$TMP/p.base" "$TMP/p.head" && echo SAME || echo DIFF)" >> "$OUT"
done
cat "$OUT"
