#!/bin/bash
# 341-r1: verify the three shapes of one byte, end to end, for my new test file.
# Why: git's core.autocrlf warns "LF will be replaced by CRLF the next time Git touches it", which means the
# BLOB is stored LF while a checkout would lay CRLF down. The dispatch requires my new bytes to be LF - the
# question this gauge answers is what the committed blob actually is, and whether gofumpt is clean against
# both the worktree bytes and a fresh checkout's bytes (⛔ touching another leg's files either way).
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 2
unset GOTMPDIR TMPDIR
out=.scratch/wisp/probes/341/r1/logs/16-eol-proof.txt
f=internal/config/writeguard_unnailed_341_test.go
# the commit that added MY file, found by the strict subject-prefix ruler (HEAD~1 is not safe: other legs land between).
myc=$(git log --perl-regexp --grep='^票 341 落地腿 341-r1·' --format='%H' --reverse -- "$f" | head -1)
{
  echo "[config] core.autocrlf=$(git config --get core.autocrlf || echo unset) ; eol=$(git config --get core.eol || echo unset)"
  echo "[.gitattributes?] $(git ls-files -- .gitattributes | head -1)  (empty = none tracked)"
  echo "[commit that added the file, strict ruler] = $myc"
  echo
  echo "[CR bytes in the WORKTREE file]  = $(tr -dc '\r' < "$f" | wc -c)"
  echo "[CR bytes in the COMMITTED BLOB] = $(git show "$myc:$f" | tr -dc '\r' | wc -c)"
  echo "[blob id] $(git rev-parse "$myc:$f")"
  echo
  echo "[gofumpt -l against the worktree bytes]"; /d/work/base/gopath/bin/gofumpt.exe -l "$f"; echo "  rc=$?"
  echo
  echo "[what a fresh checkout would contain: expand into a numbered scratch dir, measure, do not delete]"
  sc=/d/tmp/wisp/20261018/341r1/co_v1
  mkdir -p "$sc" || exit 2
  git archive "$myc" internal/config/writeguard_unnailed_341_test.go | tar -x -C "$sc" || echo "  archive rc=$?"
  co="$sc/internal/config/writeguard_unnailed_341_test.go"
  echo "  checked-out CR bytes = $(tr -dc '\r' < "$co" | wc -c)  bytes=$(wc -c < "$co")"
  echo "  gofumpt -l on the checked-out copy:"; /d/work/base/gopath/bin/gofumpt.exe -l "$co"; echo "    rc=$?"
  echo
  echo "[the 227 recipe file, for comparison: what does ITS committed blob look like?]"
  echo "  worktree CR=$(tr -dc '\r' < internal/config/writeguard_handedit_227_test.go | wc -c) blob CR=$(git show HEAD:internal/config/writeguard_handedit_227_test.go | tr -dc '\r' | wc -c)"
  echo
  echo "[named inventory: which files in internal/config are gofumpt-dirty in the worktree right now]"
  /d/work/base/gopath/bin/gofumpt.exe -l internal/config | sed 's/^/    /'; echo "  rc=$?"
} > "$out" 2>&1
echo "written $(wc -c < "$out") bytes"
