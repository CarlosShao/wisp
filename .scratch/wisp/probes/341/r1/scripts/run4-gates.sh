#!/bin/bash
# 341-r1 gates, shared worktree ONLY (dispatch AC#4). unset GOTMPDIR/TMPDIR first (票 311 AC#1 shape B GUARD).
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 2
unset GOTMPDIR TMPDIR
out=.scratch/wisp/probes/341/r1/logs/04-gates-start.txt
f=internal/config/writeguard_unnailed_341_test.go
{
  echo "[gate 1] gofumpt -l <my file>  (empty + rc=0 = clean)"
  /d/work/base/gopath/bin/gofumpt.exe -l "$f"; echo "rc=$?"
  echo
  echo "[gate 1b] the pre-existing CRLF inventory, named only, ⛔ touched by this leg"
  /d/work/base/gopath/bin/gofumpt.exe -l internal/config 2>/dev/null | sed 's/^/    /'; echo "rc=$?"
  echo
  echo "[gate 2] go vet ./internal/config/"
  go vet ./internal/config/ 2>&1 | tail -20; echo "rc=${PIPESTATUS[0]}"
  echo
  echo "[gate 3] sh scripts/d22scan.sh"
  sh scripts/d22scan.sh 2>&1 | tail -12; echo "rc=${PIPESTATUS[0]}"
  echo
  echo "[gate 4] sh scripts/check-path-length-budget.sh"
  sh scripts/check-path-length-budget.sh 2>&1 | tail -12; echo "rc=${PIPESTATUS[0]}"
} > "$out" 2>&1
echo "written $(wc -c < "$out") bytes"
