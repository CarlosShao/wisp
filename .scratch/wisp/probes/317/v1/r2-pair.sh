#!/bin/bash
# 317-v1 (verifier leg) instrument. Reads nothing in the repo except the exported
# clean tree at 8235e683. Does NOT modify any production file.
#
# Usage: bash r2-pair.sh '<GOTMPDIR value exactly as passed>' <out-log-path>
#
# Why: AC#1 judgement point (2) - the two-form contrast must differ ONLY in the
# separator. The implementing leg's pair used two different leaf names
# (C:/tmp/wisp317r1gotmp/posix vs C:\tmp\wisp317r1gotmp\nat), so this run pins the
# leaf name (C:\tmp\wisp317v1pair vs C:/tmp/wisp317v1pair = same directory, same
# spelling length, separators only) and prints the value byte-exactly.
set -u
tree="D:/tmp/wisp317v1/tree"
out="$2"

export GOTMPDIR="$1"

{
    echo "== 317-v1 matched-pair direct run =="
    echo "date=$(date '+%Y-%m-%d %H:%M:%S %z')"
    echo "== GOTMPDIR as passed, byte-exact (od -c) =="
    printf '%s' "$GOTMPDIR" | od -c
    echo "== rendered: [$GOTMPDIR]"
    case "$GOTMPDIR" in
        [A-Za-z]:\\*) echo "== leading-prefix form = NATIVE (drive + backslash)" ;;
        [A-Za-z]:/*)  echo "== leading-prefix form = POSIX (drive + slash)" ;;
        *)            echo "== leading-prefix form = other" ;;
    esac
    case "$GOTMPDIR" in
        */*) echo "== contains a slash anywhere: YES" ;;
        *)   echo "== contains a slash anywhere: NO" ;;
    esac
    mkdir -p "$GOTMPDIR"
    echo "mkdir -p rc=$?"
    cd "$tree" || exit 9
    echo "== command: go test ./internal/risk/ -count=1 -v -run <the four 116/105 names> =="
    go test ./internal/risk/ -count=1 -v -run 'TestSyncAncestorActableLegFailsClosedOnMovedAncestor116|TestSyncFirstActableLegStillRefusesHomeRewriteAtVerdict116|TestSyncAncestorLegReResolvesThroughJunctionToRealTree|TestSyncAncestorLegNotExemptedStaysFailClosed'
    echo "GO_TEST_RC=$?"
} >"$out" 2>&1

echo "log=$out"
echo "RUN=$(grep -c '^=== RUN' "$out") TOPPASS=$(grep -c '^--- PASS' "$out") TOPFAIL=$(grep -c '^--- FAIL' "$out") TOPSKIP=$(grep -c '^--- SKIP' "$out")"
grep -n 'GO_TEST_RC' "$out"
grep -n 'Invalid switch\|precondition:' "$out" | cut -c1-200
