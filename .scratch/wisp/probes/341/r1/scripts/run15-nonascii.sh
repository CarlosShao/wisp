#!/bin/bash
# 341-r1: does this repo already carry non-ASCII in Go comments? (decides whether my new file must be ASCII-only)
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 2
out=.scratch/wisp/probes/341/r1/logs/14-nonascii-precedent.txt
{
  echo "[files under internal/ + cmd/ whose comments carry non-ASCII bytes, at HEAD]"
  for f in $(git ls-files 'internal/**/*.go' 'cmd/**/*.go' | head -400); do
    n=$(tr -dc '\200-\377' < "$f" | wc -c)
    [ "$n" -gt 0 ] && echo "  $n bytes  $f"
  done | sort -rn | head -25
  echo
  echo "count of non-ASCII-bearing .go files = $(for f in $(git ls-files 'internal/**/*.go' 'cmd/**/*.go'); do n=$(tr -dc '\200-\377' < "$f" | wc -c); [ "$n" -gt 0 ] && echo x; done | wc -l)"
  echo "total .go files tracked              = $(git ls-files 'internal/**/*.go' 'cmd/**/*.go' | wc -l)"
  echo
  echo "[the recipe file 227 used: writeguard_handedit_227_test.go non-ASCII bytes]"
  echo "  $(tr -dc '\200-\377' < internal/config/writeguard_handedit_227_test.go | wc -c)"
  echo "[mine, for comparison]"
  echo "  $(tr -dc '\200-\377' < internal/config/writeguard_unnailed_341_test.go | wc -c)"
  echo
  echo "[d22scan verdict on the current tree, with my file in it]"
  sh scripts/d22scan.sh 2>&1 | tail -1
  echo "rc=${PIPESTATUS[0]}"
  echo
  echo "[does d22scan exempt comments for ban #8? cite its own words]"
  grep -nE 'comment|exempt|strings' tools/d22scan/main.go | head -12
} > "$out" 2>&1
echo "written $(wc -c < "$out") bytes"
