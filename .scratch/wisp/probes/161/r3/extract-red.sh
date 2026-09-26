#!/bin/sh
# 161-r3 cell-1: extract from a `go test -v` log the four numbers, the FAIL roster and the
# verbatim red sentences = the contiguous block of indented "file.go:NN: message" lines that
# sit directly above each "--- FAIL:" line (so diagnostics logged by PASSING tests are not
# counted as red sentences).
# Usage: sh extract-red.sh <log> > <out>
set -eu
log="${1:?usage: extract-red.sh <log>}"
printf 'RUN=%s\n'  "$(grep -c '^=== RUN'   "$log" || true)"
printf 'PASS=%s\n' "$(grep -c '^--- PASS'  "$log" || true)"
printf 'FAIL=%s\n' "$(grep -c '^--- FAIL'  "$log" || true)"
printf 'SKIP=%s\n' "$(grep -c '^--- SKIP'  "$log" || true)"
echo "--- failing test names (sorted) ---"
grep '^--- FAIL' "$log" | sed 's/ (.*//' | sort
echo "--- red sentences (verbatim, block above each --- FAIL:) ---"
awk '
  /^[[:space:]]+/ { buf = buf $0 "\n"; next }
  /^--- FAIL/ { print ">>> " $0; printf "%s", buf; buf=""; next }
  { buf="" }
' "$log"
