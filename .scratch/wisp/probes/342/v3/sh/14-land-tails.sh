#!/bin/sh
# 342-v3 landing: append the two tails (exactly one bullet each), then re-measure the ruler that backs
# every 枚-count claim in the write face. Commits happen in the next script, with literal pathspecs.
set -e
cd "/d/work/workspace/projects plans/Wisp"
RULER=/d/tmp/wisp/20261020/342v3/sh

f342=$(ls .scratch/wisp/issues/342-*.md)
f348=$(ls .scratch/wisp/issues/348-*.md)
echo "TICKETS | $f342 | $f348"
echo "before: 342 lines=$(wc -l < "$f342") bytes=$(wc -c < "$f342") | 348 lines=$(wc -l < "$f348") bytes=$(wc -c < "$f348")"

cat "$RULER/tail-342.txt" >> "$f342"
cat "$RULER/tail-348.txt" >> "$f348"

echo "== numstat of the two tickets (expect additions-only, deletions 0) =="
git diff --numstat -- "$f342" "$f348"
echo "== this leg's bullets per ticket (expect 1 / 1) =="
echo "342=$(grep -c 'agent=`342-v3`' "$f342" || true) 348=$(grep -c 'agent=`342-v3`' "$f348" || true)"
echo "== frames after the append (expect 342=5/2, 348=6/1) =="
echo "342 [x]=$(grep -c '^- \[x\] \*\*AC' "$f342" || true) [ ]=$(grep -c '^- \[ \] \*\*AC' "$f342" || true)"
echo "348 [x]=$(grep -c '^- \[x\] \*\*AC' "$f348" || true) [ ]=$(grep -c '^- \[ \] \*\*AC' "$f348" || true)"
echo "== appended bullets that look like AC frames (expect 0 / 0) =="
echo "342=$(awk 'NR>68' "$f342" | grep -c '^- \[[ x]\] \*\*AC' || true) 348=$(awk 'NR>54' "$f348" | grep -c '^- \[[ x]\] \*\*AC' || true)"
echo "== write-face file counts (the ruler behind every 枚 claim) =="
c342=$(find .scratch/wisp/probes/342/v3 -type f | wc -l)
c348=$(find .scratch/wisp/probes/348/v2 -type f | wc -l)
l342=$(find .scratch/wisp/probes/342/v3/logs -type f | wc -l)
echo "probes/342/v3=$c342 probes/348/v2=$c348 logs342=$l342"
echo "== every 枚-count claim about MY OWN roots, with the measured number beside it =="
grep -o '（[0-9]\{1,3\} 枚[^）]\{0,14\}' docs/evidence/s1/342-v3-verdict.md docs/evidence/s1/348-v2-verdict.md .scratch/wisp/probes/342/v3/report.md .scratch/wisp/probes/348/v2/readings-348.md "$RULER/tail-342.txt" "$RULER/tail-348.txt" || true
