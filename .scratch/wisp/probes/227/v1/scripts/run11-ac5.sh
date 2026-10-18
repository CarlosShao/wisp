#!/bin/bash
# 227-v1 AC#5 partial: whole ./internal/... on the delivery export tree (cmd/wisp half = lane-gated, NOT run here)
set -u
root='D:/tmp/wisp/20261018/227v1'
out='/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/227/v1/logs'
export GOTMPDIR="$root/gotmp"
unset TMPDIR
cd "$root/deliv" || exit 2
pkgs=$(go list ./internal/...)
echo "pkgcount=$(echo "$pkgs" | wc -l)" > "$out/12-ac5-internal-all.txt"
go test ./internal/... -count=1 -timeout 30m -v > "$out/12-ac5-internal-all-full.txt" 2>&1
rc=$?
RUN=$(grep -c '^=== RUN' "$out/12-ac5-internal-all-full.txt" || true)
PASSTOP=$(grep -c '^--- PASS' "$out/12-ac5-internal-all-full.txt" || true)
FAILTOP=$(grep -c '^--- FAIL' "$out/12-ac5-internal-all-full.txt" || true)
SKIP=$(grep -c '^--- SKIP' "$out/12-ac5-internal-all-full.txt" || true)
{
  echo "tree=deliv(HEAD) scope=./internal/... (NOT ./...; cmd/wisp half deliberately NOT run = 车道归编排者)"
  echo "AC#5 rc=$rc RUN=$RUN PASS_TOP=$PASSTOP FAIL_TOP=$FAILTOP SKIP=$SKIP"
  echo "---- named FAIL roster ----"
  grep '^--- FAIL' "$out/12-ac5-internal-all-full.txt" || echo NONE
  echo "---- package lines ----"
  grep -E '^(ok|FAIL|---)\s' "$out/12-ac5-internal-all-full.txt" | grep -E '^(ok|FAIL)' | sed 's/[0-9.]*s$//' | sort | uniq -c | head -40
  echo "---- build/setup failure sniff (台面故障 must not count) ----"
  grep -nE 'build failed|setup failed|no test files|0xc0000135|cannot find package' "$out/12-ac5-internal-all-full.txt" | head -20 || echo "NONE"
} >> "$out/12-ac5-internal-all.txt"
echo "rc=$rc"
