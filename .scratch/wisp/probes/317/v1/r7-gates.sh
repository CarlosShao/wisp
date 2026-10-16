#!/bin/bash
# 317-v1 gate instrument. Runs the five gate readings the dispatch asks for and
# records rc per gate into a log. It never edits anything; it only reads.
# Usage: bash r7-gates.sh <tree-dir> <out-log>
set -u
tree="$1"
out="$2"
{
    echo "== 317-v1 gates =="
    echo "date=$(date '+%Y-%m-%d %H:%M:%S %z')"
    echo "tree=$tree"
    cd "$tree" || exit 9
    echo "--- git surface of this tree ---"
    git rev-parse --short HEAD 2>&1 | head -1
    echo "GATE bash-n-portable rc-start"
    bash -n scripts/portable-tests.sh; echo "GATE bash -n scripts/portable-tests.sh rc=$?"
    bash -n scripts/winsec-tests.sh;    echo "GATE bash -n scripts/winsec-tests.sh rc=$?"
    bash -n scripts/wisp-cli-tests.sh;  echo "GATE bash -n scripts/wisp-cli-tests.sh rc=$?"
    bash -n tools/d22scan/runtests.sh;  echo "GATE bash -n tools/d22scan/runtests.sh rc=$?"
    go build ./... 2>&1 | tail -20;     echo "GATE go build ./... rc=${PIPESTATUS[0]}"
    go vet ./internal/risk/ 2>&1 | tail -10; echo "GATE go vet ./internal/risk/ rc=${PIPESTATUS[0]}"
    echo "--- gofmt -l over the four touched scripts plus internal/risk (listed, not counted blindly) ---"
    gofmt -l internal/risk scripts tools 2>&1; echo "GATE gofmt -l rc=$? (rc is always 0; the list above is the reading)"
    echo "--- d22scan gate ---"
    sh scripts/d22scan.sh 2>&1 | tail -25; echo "GATE sh scripts/d22scan.sh rc=${PIPESTATUS[0]}"
    echo "--- path length gate ---"
    sh scripts/check-path-length-budget.sh --with-self-test 2>&1 | tail -12; echo "GATE length rc=${PIPESTATUS[0]}"
} >"$out" 2>&1
echo "log=$out"
grep -n '^GATE' "$out"
