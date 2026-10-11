#!/usr/bin/env bash
# 308-r2 verification #2 (orchestrator): before/after roster pair, SAME surface type.
# Two export trees (git archive, whole repo, go.mod/go.sum included), one at the
# dispatch baseline 202ed4a5 and one at HEAD daa2f9be. Ruler = go test ./internal/risk/
# -count=1 -v. Then name+colour roster diff and the dropped-branch sentinel count.
set -uo pipefail
cd "$(git rev-parse --show-toplevel)"
LOG=.scratch/wisp/probes/orch/20261011/logs/orch-A848-308r2-roster.txt
mkdir -p "$(dirname "$LOG")"
: > "$LOG"
A=$(mktemp -d); B=$(mktemp -d)
echo "treeA=$A treeB=$B" >> "$LOG"

for pair in "202ed4a5:$A" "daa2f9be:$B"; do
  sha=${pair%%:*}; dir=${pair##*:}
  git archive --format=tar "$sha" | tar -x -C "$dir"
  echo "$sha go_files=$(find "$dir" -name '*.go' | wc -l) has_go.mod=$([ -f "$dir/go.mod" ] && echo yes || echo no)" >> "$LOG"
done

for tag in A B; do
  dir=$( [ "$tag" = A ] && echo "$A" || echo "$B" )
  ( cd "$dir" && GOFLAGS= go test ./internal/risk/ -count=1 -v ) > "$LOG.$tag" 2>&1
  echo "rc_$tag=$?" >> "$LOG"
done

for tag in A B; do
  F="$LOG.$tag"
  echo "--- $tag ---" >> "$LOG"
  echo "RUN=$(grep -c '^=== RUN' "$F") toplevel_PASS=$(grep -c '^--- PASS' "$F") any_PASS=$(grep -c '\-\-\- PASS' "$F") SKIP=$(grep -c '^--- SKIP' "$F") FAIL=$(grep -c '^--- FAIL' "$F") panic=$(grep -c '^panic:' "$F")" >> "$LOG"
  echo "sentinel_dropped_the_two=$(grep -c 'dropped the two volume-derived' "$F")" >> "$LOG"
  echo "sentinel_cannot_be_constructed=$(grep -c 'cannot be constructed' "$F")" >> "$LOG"
  grep '^--- SKIP' "$F" >> "$LOG"
  grep -E '^(--- FAIL|ok |FAIL)' "$F" >> "$LOG"
  grep -oE '^--- (PASS|FAIL|SKIP): [A-Za-z0-9_/]+' "$F" | sed 's/^--- [A-Z]*: //' | sort > "$LOG.$tag.names"
  grep -E '^(=== RUN|--- (PASS|FAIL|SKIP)): ' "$F" | grep -E '^--- ' | sort > "$LOG.$tag.roster"
done

echo "=== roster diff A vs B (empty = identical name+colour) ===" >> "$LOG"
diff "$LOG.A.roster" "$LOG.B.roster" >> "$LOG" 2>&1
echo "diff_rc=$?" >> "$LOG"
echo "namesA=$(wc -l < "$LOG.A.names") namesB=$(wc -l < "$LOG.B.names")" >> "$LOG"
cat "$LOG"
echo "=== volumeless case result in each tree ==="
grep -E 'TestVolumeInvariantSpellingsVolumeLess308|TestClassifySpellingInvariance|TestClassifyAnchorSpellingIsNotVerdict' "$LOG.A" "$LOG.B" | grep '^---' | cat
