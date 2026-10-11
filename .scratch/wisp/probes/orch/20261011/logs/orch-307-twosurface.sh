#!/bin/sh
# Orchestrator's own two-surface rerun for ticket 307 AC#1 (export tree, new archive template).
set -u
REPO="/d/work/workspace/projects plans/Wisp"
TREE="/c/Users/swq/tmp/wisp-307-orch-tree"
ALIAS="/c/Users/swq/tmp/riskattrib1/RunnerTempAreaForAttribution1"
ALIAS_SHORT='C:\Users\swq\tmp\RISKAT~1\RUNNER~1'
LOG="$REPO/.scratch/wisp/probes/orch/20261011/logs/orch-307-twosurface.txt"

{
  echo "# orchestrator two-surface rerun start $(date '+%H:%M:%S') HEAD=$(cd "$REPO" && git rev-parse --short HEAD)"
  rm -rf "$TREE" 2>/dev/null; rm_rf_rc=$?
  mkdir -p "$TREE" || { echo "tree_mkdir_rc=$?"; exit 1; }
  echo "rm_rf_rc=$rm_rf_rc (target is this script's own out-of-repo tree $TREE)"
  cd "$REPO" || exit 1
  git archive HEAD go.mod go.sum internal cmd | tar -x -C "$TREE"
  echo "archive_rc=$?"
  ls "$TREE" | tr '\n' ' '; echo
  echo "tree_go_files=$(find "$TREE" -name '*.go' | wc -l)"

  echo "=== identity check tree vs repo (cmp -s, 3 files) ==="
  for f in internal/risk/syncdirs.go internal/risk/blacklist.go internal/risk/provenance.go; do
    if cmp -s "$TREE/$f" "$REPO/$f"; then echo "SAME $f"; else echo "DIFF $f"; fi
  done
  echo "=== internal/cmd dirty in repo (should be 0) ==="
  echo "dirty=$(git status --porcelain -- internal cmd | wc -l)"

  echo "=== surface A: ALIASED TEMP ($ALIAS_SHORT) ==="
  cd "$TREE" || exit 1
  mkdir -p "$ALIAS" || true
  TEMP="$ALIAS_SHORT" TMP="$ALIAS_SHORT" TMPDIR="$ALIAS_SHORT" \
    go test ./internal/risk/ -count=1 -timeout 600s -v > "$LOG.tmp-aliased" 2>&1
  echo "aliased_rc=$?"
  echo "aliased PASS=$(grep -c '^    --- PASS\|^--- PASS' "$LOG.tmp-aliased")"
  echo "aliased FAIL=$(grep -c '^    --- FAIL\|^--- FAIL' "$LOG.tmp-aliased")"
  echo "aliased SKIP=$(grep -c '^    --- SKIP\|^--- SKIP' "$LOG.tmp-aliased")"
  echo "aliased fail names:"
  grep -E '^\s*--- FAIL' "$LOG.tmp-aliased" | sed 's/^/  /'

  echo "=== surface B: default TEMP (no alias on this machine) ==="
  go test ./internal/risk/ -count=1 -timeout 600s -v > "$LOG.tmp-default" 2>&1
  echo "default_rc=$?"
  echo "default PASS=$(grep -c '^    --- PASS\|^--- PASS' "$LOG.tmp-default")"
  echo "default FAIL=$(grep -c '^    --- FAIL\|^--- FAIL' "$LOG.tmp-default")"
  echo "default SKIP=$(grep -c '^    --- SKIP\|^--- SKIP' "$LOG.tmp-default")"
  echo "default fail names:"
  grep -E '^\s*--- FAIL' "$LOG.tmp-default" | sed 's/^/  /'
  echo "alias short path used exists: $([ -d "$ALIAS" ] && echo yes)"
  echo "# done $(date '+%H:%M:%S')"
} >> "$LOG" 2>&1
tail -40 "$LOG"
