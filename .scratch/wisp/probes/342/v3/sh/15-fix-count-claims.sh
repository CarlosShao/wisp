#!/bin/sh
# 342-v3: make every file-count claim in the write face equal the freshly measured counts.
# Ruler = `find <root> -type f | wc -l`; the same ruler's full output is archived as
# .scratch/wisp/probes/342/v3/logs/30-writeface-hygiene-final.txt (and a copy in the 348 root).
set -e
cd "/d/work/workspace/projects plans/Wisp"
c342=$(find .scratch/wisp/probes/342/v3 -type f | wc -l)
c348=$(find .scratch/wisp/probes/348/v2 -type f | wc -l)
echo "MEASURED probes/342/v3=$c342 probes/348/v2=$c348"
FILES="docs/evidence/s1/342-v3-verdict.md docs/evidence/s1/348-v2-verdict.md .scratch/wisp/probes/342/v3/report.md .scratch/wisp/probes/348/v2/readings-348.md /d/tmp/wisp/20261020/342v3/sh/tail-342.txt /d/tmp/wisp/20261020/342v3/sh/tail-348.txt"
for f in $FILES; do
  perl -i -pe "s/（41 枚/（$c342 枚/g; s/（9 枚同）/（$c348 枚同）/g; s/（12 枚同）/（$c348 枚同）/g; s/（12 枚：/（$c348 枚：/g; s/（42 枚）/（$c342 枚）/g;" "$f"
done
echo "-- count claims left (39 枚 = 342-r2's own bundle = a verified claim about ANOTHER leg, correct as-is) --"
for f in $FILES; do echo "--- $f"; grep -o '（[0-9]\{1,3\} 枚[^）]\{0,10\}' "$f" || true; done
