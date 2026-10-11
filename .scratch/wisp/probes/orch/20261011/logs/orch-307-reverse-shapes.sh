#!/bin/sh
# Orchestrator's own reverse-shape (mutation) reruns for ticket 307 AC#2, on MY out-of-repo export tree.
set -u
TREE="/c/Users/swq/tmp/wisp-307-orch-tree"
ALIAS="/c/Users/swq/tmp/riskattrib1/RunnerTempAreaForAttribution1"
ALIAS_SHORT='C:\Users\swq\tmp\RISKAT~1\RUNNER~1'
SRC="$TREE/internal/risk/syncdirs.go"
LOG="/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/orch/20261011/logs/orch-307-reverse-shapes.txt"
cp "$SRC" "$SRC.orchpristine"
ORIG_SHA=$(sha256sum "$SRC.orchpristine" | cut -d' ' -f1)

run_aliased () {
  cd "$TREE" || return 9
  TEMP="$ALIAS_SHORT" TMP="$ALIAS_SHORT" TMPDIR="$ALIAS_SHORT" \
    go test ./internal/risk/ -count=1 -timeout 600s -v > "$LOG.raw-$1" 2>&1
  echo "rc=$?"
  echo "PASS=$(grep -c '^    --- PASS\|^--- PASS' "$LOG.raw-$1") FAIL=$(grep -c '^    --- FAIL\|^--- FAIL' "$LOG.raw-$1") SKIP=$(grep -c '^    --- SKIP\|^--- SKIP' "$LOG.raw-$1")"
  grep -E '^\s*--- FAIL' "$LOG.raw-$1" | sed 's/^/    /'
}

{
  echo "# orchestrator reverse-shape reruns start $(date '+%H:%M:%S') tree=$TREE pristine_sha256=$ORIG_SHA"

  echo "== shape A: BOTH branches back to isUnder (pre-fix single-sided fold) =="
  sed -i -e 's/aUnder(candForms, anchorForms(e\.canon))/isUnder(cand, e.canon)/' \
         -e 's/aUnder(candForms, anchorForms(s\.home))/isUnder(cand, s.home)/' \
         -e 's/^\tcandForms := formsOf(cand)$/\tcandForms := formsOf(cand)\n\t_ = candForms/' "$SRC"
  echo "hits_roots=$(grep -c 'isUnder(cand, e.canon)' $SRC) hits_home=$(grep -c 'isUnder(cand, s.home)' $SRC)"
  run_aliased A

  echo "== restore =="
  cp "$SRC.orchpristine" "$SRC"
  echo "sha_restored=$(sha256sum "$SRC" | cut -d' ' -f1) equal=$([ "$(sha256sum "$SRC" | cut -d' ' -f1)" = "$ORIG_SHA" ] && echo yes)"

  echo "== shape C: ONLY the root-membership branch back to isUnder (:434) =="
  sed -i -e 's/aUnder(candForms, anchorForms(e\.canon))/isUnder(cand, e.canon)/' "$SRC"
  echo "hits_roots=$(grep -c 'isUnder(cand, e.canon)' $SRC) hits_home=$(grep -c 'isUnder(cand, s.home)' $SRC)"
  run_aliased C

  echo "== restore =="
  cp "$SRC.orchpristine" "$SRC"
  echo "sha_restored=$(sha256sum "$SRC" | cut -d' ' -f1) equal=$([ "$(sha256sum "$SRC" | cut -d' ' -f1)" = "$ORIG_SHA" ] && echo yes)"

  echo "== positive control: pristine tree, aliased surface =="
  run_aliased P

  echo "# done $(date '+%H:%M:%S')"
} >> "$LOG" 2>&1
