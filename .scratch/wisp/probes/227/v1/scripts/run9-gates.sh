#!/bin/bash
# 227-v1 gate four numbers, shared working tree, AFTER unsetting GOTMPDIR/TMPDIR (票 311 AC#1 shape B GUARD)
cd '/d/work/workspace/projects plans/Wisp' || exit 2
unset GOTMPDIR
unset TMPDIR
out='.scratch/wisp/probes/227/v1/logs/10-gates-final.txt'
{
  echo "gates run in SHARED working tree; HEAD=$(git rev-parse --short HEAD); GOTMPDIR/TMPDIR unset before the four gates"
  echo "--- 1) gofumpt -l <leg's two files> (want empty output) ---"
  /d/work/base/gopath/bin/gofumpt.exe -l internal/config/writeguard.go internal/config/writeguard_handedit_227_test.go
  echo "gofumpt rc=$?"
  echo "--- 1b) gofumpt -l internal cmd (existing CRLF stock, named only, NOT mine) ---"
  /d/work/base/gopath/bin/gofumpt.exe -l internal cmd
  echo "gofumpt-all rc=$?"
  echo "--- 2) sh scripts/d22scan.sh ---"
  sh scripts/d22scan.sh
  echo "d22scan rc=$?"
  echo "--- 3) sh scripts/check-path-length-budget.sh ---"
  sh scripts/check-path-length-budget.sh
  echo "pathlen rc=$?"
  echo "--- 4) go vet ./internal/config/ (shared working tree) ---"
  go vet ./internal/config/
  echo "vet rc=$?"
} > "$out" 2>&1
echo "written bytes=$(wc -c < "$out")"
