#!/bin/bash
# M7b: second single-pass body scan with the shapes the orchestrator's ruler cannot see:
#   #NN , NN-slug , /NN- , AC#NN , 「NN 号」
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
PAT=$(paste -sd'|' "$L/roster52.txt")
REGEX="#(${PAT})([^0-9]|$)|(^|[^0-9])(${PAT})-[a-z]|/(${PAT})/|AC ?#(${PAT})([^0-9]|$)|(${PAT}) 号"
xargs -a "$L/m7-scope.txt" grep -oHE "$REGEX" 2>/dev/null | sed -E 's/:[0-9]+:/|/' > "$L/m7b-raw.txt"
echo "raw=$(wc -l < "$L/m7b-raw.txt")"
: > "$L/m7b.tsv"
while read -r N; do
  [ -z "$N" ] && continue
  L2=$(awk -F'|' -v n="$N" '$1==n{print $2}' "$L/m7b-raw.txt" | sort -u | wc -l | tr -d '\n')
  VF=$(awk -F'|' -v n="$N" '$1==n{print $2}' "$L/m7b-raw.txt" | sort -u | grep -ciE 'adversarial|accept|verdict|裁决|验收|终裁|归口|backfill|audit' | tr -d '\n')
  echo "$N|othershape_files=$L2|othershape_verdictshaped=$VF" >> "$L/m7b.tsv"
done < "$L/roster52.txt"
echo "== tickets with >=1 verdict-shaped mention in OTHER shapes but 0 in name index =="
awk -F'|' '{split($3,a,"="); if(a[2]+0>0) print $0}' "$L/m7b.tsv" | cut -d'|' -f1-3 | head -60
