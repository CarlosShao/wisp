#!/bin/sh
cd '/d/work/workspace/projects plans/Wisp' || exit 9
for p in .scratch/wisp/probes/353/v1/01-independent-readings.txt .scratch/wisp/probes/353/v1/02-ticket-draft-secC.md docs/evidence/s1/353-v1-verdict.md; do
  echo "== $p"
  git cat-file blob "HEAD:$p" 2>/dev/null | grep -n '3\.\.54\|115\.\.153\|102\.\.153\|23835\|23,835\|10\.\.31\|10\.\.18' | cut -c1-200 || true
done
