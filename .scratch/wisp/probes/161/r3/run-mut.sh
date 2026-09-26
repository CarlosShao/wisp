#!/bin/sh
# 161-r3 cell-1 runner: run the three ticket-158 mutation overlays that consume
# .scratch/wisp/probes/158/accept-r1/mut/guard.no{1,2,3}.go and archive the raw
# `go test -v` output.  Usage: sh run-mut.sh <label>
#
# Standing instrument pitfall (measured in this repo, ticket 158 acceptance): -overlay
# and -cover* must NEVER be combined -- the overlay is silently ignored.  This script
# passes no -cover flag at all.
set -eu

label="${1:?usage: run-mut.sh <label>}"
# Repo root via git, not via $0: $0 is relative when the caller does `sh path/to/run-mut.sh`
# from the repo root, and three `..` hops from .scratch/wisp/probes/161/r3 lands in .scratch/wisp.
root=$(git rev-parse --show-toplevel)
cd "$root"

ovldir=.scratch/wisp/probes/158/accept-r1/ovl
outdir=.scratch/wisp/probes/161/r3/logs

for ovl in r5-m1no1 r6-m1no2 r4-m2no3; do
  out="$outdir/$label-$ovl.log"
  echo "=== RUN $label $ovl -> $out"
  # shellcheck disable=SC2046
  rc=0
  go test -count=1 -v -overlay="$ovldir/$ovl.json" ./internal/tools/ >"$out" 2>&1 || rc=$?
  printf '%s\n' "$rc" >"$outdir/$label-$ovl.rc"
  grep -c '^=== RUN' "$out" | sed 's/^/RUN=/' || true
  grep -c '^--- PASS' "$out" | sed 's/^/PASS=/' || true
  grep -c '^--- FAIL' "$out" | sed 's/^/FAIL=/' || true
  grep -c '^--- SKIP' "$out" | sed 's/^/SKIP=/' || true
done
