#!/bin/sh
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 90
export LC_ALL=C.UTF-8
L=/d/tmp/wisp/20261019/343v1/logs/leg-hygiene
mkdir -p $L
{
echo "== line counts of the leg's authored positive files (cap = 160) =="
for f in .scratch/wisp/probes/343/r1/report.md .scratch/wisp/probes/343/r1/report-sec-b.md \
         .scratch/wisp/probes/343/r1/report-sec-c.md .scratch/wisp/probes/343/r1/logs/01-ac4-color-pair-two-platform.txt \
         .scratch/wisp/probes/343/r1/logs/02-ac1-comment-face.txt .scratch/wisp/probes/343/r1/logs/03-ac0-same-family-census.txt \
         .scratch/wisp/probes/343/r1/logs/04-gates-shared-tree.txt; do
  printf '%6s %s\n' "$(wc -l < "$f")" "$f"
done
echo
echo "== mojibake scan over the leg's tracked files: U+FFFD and literal '??' runs =="
printf 'FFFD_hits='; grep -rl $'\xef\xbf\xbd' .scratch/wisp/probes/343/r1 .scratch/wisp/issues/343-* 2>/dev/null | wc -l
printf 'question-pair_hits='; grep -rn '??' .scratch/wisp/probes/343/r1/*.md 2>/dev/null | wc -l
printf 'non-utf8_files='; for f in $(git ls-files -- .scratch/wisp/probes/343/r1); do iconv -f UTF-8 -t UTF-8 "$f" >/dev/null 2>&1 || printf '%s ' "$f"; done; echo
echo "== banned literal char (U+00A7) over the leg's write surface (positive control HANDOVER=199) =="
printf 'HANDOVER='; grep -c $'\xc2\xa7' docs/reports/HANDOVER.md || true
printf 'leg_r1='; grep -rl $'\xc2\xa7' .scratch/wisp/probes/343/r1 2>/dev/null | wc -l
printf 'ticket343='; grep -c $'\xc2\xa7' .scratch/wisp/issues/343-*.md || true
echo "== zero-byte tracked files anywhere under probes/343 (both legs) =="
find .scratch/wisp/probes/343 -type f -size 0 -print | sed 's/^/ZERO /'
echo "== .out extension among tracked files under probes/343 =="
git ls-files -- .scratch/wisp/probes/343 | grep -c '\.out$' || true
echo "== ticket tail: my bullet will be the only append after 137044b6; current tail lines =="
tail -2 .scratch/wisp/issues/343-*.md | cut -c1-120
echo DONE3
} > $L/leg-hygiene.log 2>&1
cat $L/leg-hygiene.log
