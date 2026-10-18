#!/bin/bash
# 227-v1: find a positive-control file that really contains a slot-shaped token.
cd '/d/work/workspace/projects plans/Wisp' || exit 2
for f in docs/reports/HANDOVER.md docs/reports/pending-and-issues.md .scratch/wisp/issues/README.md docs/PLAN.md; do
  h=$(grep -cE 'TODO|FIXME|待填|占位' "$f" || true)
  echo "$f slot-hits=$h"
done
echo "--- first detail lines ---"
grep -nE 'TODO|FIXME|待填|占位' docs/reports/HANDOVER.md | head -3
grep -nE 'TODO|FIXME|待填|占位' docs/reports/pending-and-issues.md | head -3
