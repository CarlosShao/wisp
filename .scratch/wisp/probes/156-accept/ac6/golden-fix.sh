#!/usr/bin/env bash
# accept-r1 AC#6 arm-fix: the ticket's "任何 golden" arm, with a pathspec that CAN fire.
# My first census used bare `testdata/golden`, which as a git prefix only matches a
# TOP-LEVEL testdata/golden/ (this repo has none) => a dead arm that could only ever
# read 0. Here the same arm is asked with glob pathspecs, and a positive control
# proves the ruler fires.
set -u
ANCHOR=7ca1130
OUT=.scratch/wisp/probes/156-accept/ac6
mkdir -p "$OUT"

GOLDEN=(':(glob)**/testdata/golden/**' ':(glob)**/golden/**' ':(glob)***.golden**')

echo "== how many files each golden pathspec selects at the anchor (firing check)"
for g in "${GOLDEN[@]}"; do
  echo "$g -> $(git ls-tree -r --name-only $ANCHOR -- "$g" | wc -l) files"
done

echo
echo "== golden arms per implementation commit (10, path-picked)"
while read -r c; do
  [ -z "$c" ] && continue
  tot=0
  for g in "${GOLDEN[@]}"; do
    n=$(git diff-tree -r --no-commit-id --name-only "$c" -- "$g" | wc -l)
    tot=$((tot+n))
  done
  echo "$(echo $c | cut -c1-7) golden_hits=$tot"
done < <(awk '{print $1}' "$OUT/impl-commits.txt")

echo
echo "== POSITIVE CONTROL for the golden arm: a commit that really moved a golden"
c=$(git log --format=%H -1 -- ':(glob)**/testdata/golden/**')
echo "picked $c $(git log -1 --format=%s $c | cut -c1-50)"
for g in "${GOLDEN[@]}"; do
  echo "   $g hits=$(git diff-tree -r --no-commit-id --name-only $c -- "$g" | wc -l)"
done

echo
echo "== thresholds.go / slo-check.ps1 / allowlist / d22scan at anchor (存量, firing check)"
for p in internal/observe/thresholds.go scripts/slo-check.ps1 tools/d22scan/allowlist.txt tools/d22scan docs/PLAN.md docs/specs internal/risk internal/panel internal/agent/approval internal/observe frontend design; do
  echo "$p -> $(git ls-tree -r --name-only $ANCHOR -- $p | wc -l)"
done

echo
echo "== whole-era (9835d81..7ca1130) hits on the contract axis, EVERY impl commit already 0 above"
git diff --name-only 9835d81 7ca1130 -- docs/PLAN.md docs/specs internal/risk internal/panel internal/agent/approval internal/observe ':(glob)**/testdata/golden/**' ':(glob)**/golden/**' tools/d22scan/allowlist.txt scripts/slo-check.ps1 tools/d22scan go.mod go.sum | sed 's/^/  RANGEHIT /'
echo "  (range diff is NOT the AC#6 ruler - ticket face demands per-commit; this is a cross-check only)"
