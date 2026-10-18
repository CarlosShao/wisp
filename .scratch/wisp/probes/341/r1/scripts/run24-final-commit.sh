#!/bin/bash
# 341-r1: the delivery commit (fifth). Literal pathspecs in BOTH `add` and `commit` - ⛔ anything derived
# from the shared index (the dispatch's A1027 trap; this leg saw foreign staged files in the index and its
# literal pathspecs kept them out, logged in logs/18 [2]).
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 2
P=.scratch/wisp/probes/341/r1
git add -- \
  "$P/report.md" \
  "$P/logs/20-selfproof-final.txt" \
  "$P/logs/21-selfproof-corrections.txt" \
  "$P/scripts/msg-6-selfproof.txt" \
  "$P/scripts/run21-commit.sh" \
  "$P/scripts/run22-selfproof.sh" \
  "$P/scripts/run23-corrections.sh" 
  "$P/scripts/run24-final-commit.sh"
git commit -F "$P/scripts/msg-6-selfproof.txt" -- \
  "$P/report.md" \
  "$P/logs/20-selfproof-final.txt" \
  "$P/logs/21-selfproof-corrections.txt" \
  "$P/scripts/msg-6-selfproof.txt" \
  "$P/scripts/run21-commit.sh" \
  "$P/scripts/run22-selfproof.sh" \
  "$P/scripts/run23-corrections.sh" 
  "$P/scripts/run24-final-commit.sh"
echo "commit rc=$?"
echo '---- final roster, strict subject-prefix ruler, with per-commit out-of-scope count ----'
for c in $(git log --perl-regexp --grep='^票 341 落地腿 341-r1·' --format='%H' --reverse); do
  oos=$(git show --name-only --format= "$c" | sed '/^$/d' | grep -vcE '^(internal/config/[^/]*_test\.go|\.scratch/wisp/probes/341/r1/|\.scratch/wisp/issues/341-)' || true)
  prod=$(git diff --numstat "$c^" "$c" | grep -E '\.go$' | grep -vE '_test\.go$' | wc -l)
  echo "  $(git rev-parse --short "$c")  paths=$(git show --name-only --format= "$c" | sed '/^$/d' | wc -l)  out-of-scope=$oos  prod-go=$prod"
done
echo "count=$(git log --perl-regexp --grep='^票 341 落地腿 341-r1·' --format='%H' | wc -l)"
echo '---- my face clean? ----'
git status --porcelain -- "$P" internal/config .scratch/wisp/issues/341-four-hand-edit-shapes-still-unnailed-and-two-226-sentences-overclaim.md
echo '[end status: empty = everything landed]'
echo '---- ticket boxes + line count ----'
T=.scratch/wisp/issues/341-four-hand-edit-shapes-still-unnailed-and-two-226-sentences-overclaim.md
echo "  unchecked=$(grep -c '^[[:space:]]*- \[ \] \*\*AC#' "$T" || true) checked=$(grep -c '^[[:space:]]*- \[x\] \*\*AC#' "$T" || true) lines=$(wc -l < "$T")"
echo '---- forbidden production blobs (anchor 98dbf20a vs HEAD) ----'
for p in internal/config/writeguard.go internal/config/loader.go internal/config/migrate.go; do
  echo "    $p  $(git rev-parse "98dbf20a:$p" | cut -c1-12) vs $(git rev-parse "HEAD:$p" | cut -c1-12)"
done
