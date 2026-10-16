#!/bin/bash
# 317-v1 instrument for AC#2 judgement (2): the "an explicit WISP_TEST_GOTMP_ROOT is
# still taken verbatim" branch. Physical mutation is NOT used here - this runs the
# SHIPPED scripts/portable-tests.sh in the exported clean tree (8235e683) with a
# posix-shaped override injected, exactly like logs/17 did, and checks
#   (a) the override reaches GOTMPDIR byte-for-byte (operator spelling not rewritten),
#   (b) the four 116/105 cases still go red under it => the reproduction channel and the
#       cases' sensitivity were not dismantled by the fix.
# Usage: bash r6-override.sh <tree-dir> <out-log> <override-value>
set -u
tree="$1"
out="$2"
ovr="$3"
{
    echo "== 317-v1 posix-override boundary run =="
    echo "date=$(date '+%Y-%m-%d %H:%M:%S %z')"
    echo "tree=$tree (shipped, unmutated)"
    echo "== WISP_TEST_GOTMP_ROOT as passed, byte-exact (od -c) =="
    printf '%s' "$ovr" | od -c
    echo "== rendered: [$ovr]"
    cd "$tree" || exit 9
    echo "== command: env -u GOTMPDIR WISP_TEST_GOTMP_ROOT=<as above> bash scripts/portable-tests.sh ./internal/risk/ =="
    env -u GOTMPDIR "WISP_TEST_GOTMP_ROOT=$ovr" bash scripts/portable-tests.sh ./internal/risk/
    echo "HARNESS_RC=$?"
} >"$out" 2>&1
echo "log=$out"
grep -n 'GOTMPDIR=' "$out" | head -3 | cut -c1-240
echo "RUN=$(grep -c '^=== RUN' "$out") TOPPASS=$(grep -c '^--- PASS' "$out") TOPFAIL=$(grep -c '^--- FAIL' "$out") TOPSKIP=$(grep -c '^--- SKIP' "$out")"
grep -n 'HARNESS_RC' "$out"
grep -n '^\s*--- FAIL' "$out" | cut -c1-120
echo "-- surface scan --"
grep -n 'walk \|build failed\|no tests to run\|0xc0000135' "$out" | head -3 | cut -c1-160
