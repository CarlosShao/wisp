#!/bin/bash
# 317-v1 (verifier leg) harness-side instrument. Runs the SHIPPED windows entry
# scripts/portable-tests.sh against ./internal/risk/ inside an exported clean tree
# (git archive 8235e683), never in the shared working tree, and never with an
# overlay. It only reads; the tree is the thing that gets mutated, outside the repo.
#
# Usage: bash r3-harness.sh <tree-dir> <out-log>
# Precondition set by the caller: GOTMPDIR and WISP_TEST_GOTMP_ROOT are drained.
set -u
tree="$1"
out="$2"
{
    echo "== 317-v1 harness run =="
    echo "date=$(date '+%Y-%m-%d %H:%M:%S %z')"
    echo "tree=$tree"
    echo "== the two composition lines as they now stand in this tree (byte-exact) =="
    grep -n 'gotmp_base=\|gotmp_run=' "$tree/scripts/portable-tests.sh" | head -6
    echo "== env drained check (must print nothing) =="
    env | grep -E '^(GOTMPDIR|WISP_TEST_GOTMP_ROOT)=' || echo "(both drained)"
    cd "$tree" || exit 9
    echo "== command: bash scripts/portable-tests.sh ./internal/risk/ =="
    env -u GOTMPDIR -u WISP_TEST_GOTMP_ROOT bash scripts/portable-tests.sh ./internal/risk/
    echo "HARNESS_RC=$?"
} >"$out" 2>&1
echo "log=$out"
grep -n 'GOTMPDIR=' "$out" | head -4 | cut -c1-260
echo "RUN=$(grep -c '^=== RUN' "$out") TOPPASS=$(grep -c '^--- PASS' "$out") TOPFAIL=$(grep -c '^--- FAIL' "$out") TOPSKIP=$(grep -c '^--- SKIP' "$out")"
grep -n 'HARNESS_RC' "$out"
echo "-- surface-failure scan (walk/build failed/no tests to run/0xc0000135) --"
grep -n 'walk \|build failed\|no tests to run\|0xc0000135\|cannot find' "$out" | head -5 | cut -c1-200
echo "-- ticket-317 symptom lines, first 6 --"
grep -n 'Invalid switch\|precondition:' "$out" | head -6 | cut -c1-190
