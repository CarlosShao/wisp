#!/bin/sh
cd "/d/work/workspace/projects plans/Wisp" || exit 1
cp /d/tmp/wisp/20261020/342v3/logs/writeface-hygiene.txt .scratch/wisp/probes/342/v3/logs/30-writeface-hygiene-final.txt
cp /d/tmp/wisp/20261020/342v3/logs/writeface-hygiene.txt .scratch/wisp/probes/348/v2/logs/04-writeface-hygiene.txt
echo "ROOT342=$(find .scratch/wisp/probes/342/v3 -type f | wc -l)"
echo "ROOT348=$(find .scratch/wisp/probes/348/v2 -type f | wc -l)"
echo "LOGS342=$(find .scratch/wisp/probes/342/v3/logs -type f | wc -l)"
echo "SH342=$(find .scratch/wisp/probes/342/v3/sh -type f | wc -l)"
echo "ZERO=$(find .scratch/wisp/probes/342/v3 .scratch/wisp/probes/348/v2 -type f -size 0 | wc -l)"
echo "OUT=$(find .scratch/wisp/probes/342/v3 .scratch/wisp/probes/348/v2 -type f -name '*.out' | wc -l)"
echo "-- claims about my own roots in every doc --"
grep -o '（[0-9]\{1,3\} 枚[^）]\{0,14\}' docs/evidence/s1/342-v3-verdict.md docs/evidence/s1/348-v2-verdict.md .scratch/wisp/probes/342/v3/report.md .scratch/wisp/probes/348/v2/readings-348.md /d/tmp/wisp/20261020/342v3/sh/tail-342.txt /d/tmp/wisp/20261020/342v3/sh/tail-348.txt || true
