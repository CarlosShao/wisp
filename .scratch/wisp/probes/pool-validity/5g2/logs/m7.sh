#!/bin/bash
# M7: single-pass body scan of the whole scope for the 52 roster numbers, so that evidence
# filed under ANOTHER ticket's number (cross-filing / 归口 elsewhere) becomes visible.
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
PAT=$(paste -sd'|' "$L/roster52.txt")
{ find docs/evidence docs/reports -name '*.md'; find .scratch/wisp/probes -name '*.md' \
    -not -path '.scratch/wisp/probes/232/*' -not -path '.scratch/wisp/probes/frontend/*' \
    -not -path '.scratch/wisp/probes/design/*' -not -path '.scratch/wisp/probes/pool-validity/5g/*'; } \
  | sort -u > "$L/m7-scope.txt"
echo "scope=$(wc -l < "$L/m7-scope.txt")"
REGEX="票 ?(${PAT})([^0-9]|$)|[Tt]icket ?(${PAT})([^0-9]|$)"
xargs -a "$L/m7-scope.txt" grep -oHE "$REGEX" 2>/dev/null \
  | sed -E 's/:[0-9]+:/|/' > "$L/m7-raw.txt"
echo "raw_hits=$(wc -l < "$L/m7-raw.txt")"
cut -d'|' -f1 "$L/m7-raw.txt" | sort | uniq -c | sort -rn | head -30 > "$L/m7-filefreq.txt"
while read -r N; do
  [ -z "$N" ] && continue
  CNT=$(awk -F'|' -v n="$N" '$1==n' "$L/m7-raw.txt" | wc -l)
  VF=$(awk -F'|' -v n="$N" '$1==n{print $2}' "$L/m7-raw.txt" | sort -u | grep -ciE 'adversarial|accept|verdict|裁决|验收|判定|终裁|归口|backfill|audit' | tr -d '\n')
  ALL=$(awk -F'|' -v n="$N" '$1==n{print $2}' "$L/m7-raw.txt" | sort -u | wc -l | tr -d '\n')
  NAMEIDX=$(awk -F'|' -v n="$N" '$1==n' "$L/m3-list.txt" | wc -l | tr -d '\n')
  echo "$N|bodymentions=$CNT|files_mentioning=$ALL|of_which_verdictshaped=$VF|in_name_index=$NAMEIDX" >> "$L/m7.tsv"
done < "$L/roster52.txt"
echo "M7=$(wc -l < "$L/m7.tsv")"
awk -F'|' '{split($5,a,"=");split($4,b,"="); if(a[2]+0==0 && b[2]+0>0) printf "%s(%s) ",$1,b[2]}' "$L/m7.tsv"; echo "<<BODY-ONLY-CITATIONS"
awk -F'|' '{split($5,a,"="); if(a[2]+0==0 && $2=="bodymentions=0") printf "%s ",$1}' "$L/m7.tsv"; echo "<<NOBODY_NONAME"
