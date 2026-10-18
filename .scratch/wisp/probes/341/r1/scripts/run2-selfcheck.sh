#!/bin/bash
# 341-r1: format + shape self-check of the new test file, in the shared worktree.
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 2
f=internal/config/writeguard_unnailed_341_test.go
out=.scratch/wisp/probes/341/r1/logs/02-newfile-selfcheck.txt
{
  echo "[ruler] wc -l / CR bytes / gofumpt -l / gofumpt -d  on $f (shared worktree, read-only for the repo)"
  echo "lines=$(wc -l < "$f")"
  echo "bytes=$(wc -c < "$f")"
  cr=$(tr -dc '\r' < "$f" | wc -c)
  echo "CR bytes = $cr  (want 0: 本机 6 枚 CRLF 存量件是别腿的，我的新字节必须 LF)"
  echo "func Test count = $(grep -c '^func Test' "$f")"
  echo
  echo "[gofumpt -l] (empty output + rc=0 = clean)"
  /d/work/base/gopath/bin/gofumpt.exe -l "$f"
  echo "rc=$?"
  echo
  echo "[gofumpt -d] (any diff = I must fix the file before committing)"
  /d/work/base/gopath/bin/gofumpt.exe -d "$f"
  echo "rc_diff=$?"
  echo
  echo "[non-ASCII bytes in the new file] = $(tr -dc '\200-\377' < "$f" | wc -c)  (want 0: 本仓注释一律 ASCII)"
} > "$out" 2>&1
echo "written $(wc -c < "$out") bytes"
