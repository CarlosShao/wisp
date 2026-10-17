#!/usr/bin/env bash
# 330-r2: the four gates, run ONLY in the shared worktree (a partial export tree
# makes d22scan refuse with exit 2, which is a GUARD range refusal, not a red).
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9
LOG=.scratch/wisp/probes/330/r2/logs/gates.log
{
  echo "### gate 1: gofumpt -l internal/agent/  (empty list = clean)"
  /d/work/base/gopath/bin/gofumpt.exe -l internal/agent/
  echo "gofumpt_rc=$?"
  echo
  echo "### gate 2: sh scripts/d22scan.sh"
  sh scripts/d22scan.sh > /d/tmp/w330-d22scan-2.log 2>&1
  echo "d22scan_sh_rc=$?"
  grep -E '^\s*(scan_test\.go:[0-9]+: repo HEAD violates|--- FAIL)' /d/tmp/w330-d22scan-2.log | head -10
  echo "  findings above are cmd/wisp bytes of the concurrent leg 329-r1, not this leg's:"
  grep -cE 'internal/agent/plan_present_crowded_world_330r2_test.go' /d/tmp/w330-d22scan-2.log | sed 's/^/  mentions_of_this_legs_new_file=/'
  echo
  echo "### gate 2b: the scan proper (d22scan.sh step 2, run standalone)"
  ( cd tools/d22scan && go run . -root "D:/work/workspace/projects plans/Wisp" > /d/tmp/w330-d22-real-2.log 2>&1 )
  echo "d22scan_real_rc=$?"
  grep -E '^(cmd|internal|frontend|design)/' /d/tmp/w330-d22-real-2.log | head -5
  tail -1 /d/tmp/w330-d22-real-2.log
  echo
  echo "### gate 3: sh scripts/check-path-length-budget.sh"
  sh scripts/check-path-length-budget.sh > /d/tmp/w330-len-2.log 2>&1
  echo "pathlen_rc=$?"
  grep -E 'VERDICT|denominator read' /d/tmp/w330-len-2.log | tail -2
  echo
  echo "### gate 4: go vet ./internal/agent/"
  go vet ./internal/agent/ > /d/tmp/w330-vet-2.log 2>&1
  echo "vet_rc=$?"
  wc -c < /d/tmp/w330-vet-2.log | sed 's/^/vet_output_bytes=/'
} 2>&1 | tee "$LOG"
