#!/bin/bash
# M7c: single-pass body scan of whole scope for the 52 roster numbers (all citation shapes),
# per-ticket rollup, plus off-name-index (cross-filed) evidence detection.
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
PAT=$(paste -sd'|' "$L/roster52.txt")
REGEX="票 ?(${PAT})([^0-9]|$)|[Tt]icket ?(${PAT})([^0-9]|$)|#(${PAT})([^0-9]|$)|/(${PAT})/|AC ?#(${PAT})([^0-9]|$)"
{ find docs/evidence docs/reports -name '*.md'; find .scratch/wisp/probes -name '*.md' \
    -not -path '.scratch/wisp/probes/232/*' -not -path '.scratch/wisp/probes/frontend/*' \
    -not -path '.scratch/wisp/probes/design/*' -not -path '.scratch/wisp/probes/pool-validity/5g/*'; } \
  | sort -u > "$L/m7-scope.txt"
echo "scope=$(wc -l < "$L/m7-scope.txt")"
xargs -a "$L/m7-scope.txt" grep -oHE "$REGEX" 2>/dev/null > "$L/m7-raw.txt"
echo "rawlines=$(wc -l < "$L/m7-raw.txt")"
awk -F: -v RN="$L/roster52.txt" '
BEGIN{ while((getline n < RN)>0) if(n!="") keep[n+0]=1 }
{ path=$1; m=$2; if(path==""||m=="") next;
  s=m;
  while(match(s,/[0-9]+/)){ d=substr(s,RSTART,RLENGTH)+0; if(d in keep){ print d "|" path } s=substr(s,RSTART+RLENGTH) }
}' RN="$L/roster52.txt" "$L/m7-raw.txt" | sort -u > "$L/m7-pairs.txt"
echo "pairs=$(wc -l < "$L/m7-pairs.txt")"
awk -F'|' 'NR==FNR{ if(FNR==FILENAME ~ /m3-list/){} ; idx[$1"|"$2]=1; next }' /dev/null /dev/null
: > "$L/m7.tsv"; : > "$L/m7-offindex.txt"
while read -r N; do
  [ -z "$N" ] && continue
  FILES=$(awk -F'|' -v n="$N" '$1==n+0{print $2}' "$L/m7-pairs.txt" | sort -u)
  ALL=$(printf '%s\n' "$FILES" | grep -c .)
  VF=$(printf '%s\n' "$FILES" | grep -ciE 'adversarial|accept|verdict|裁决|验收|终裁|归口|backfill|audit')
  NAME=$(awk -F'|' -v n="$N" '$1==n' "$L/m3-list.txt" | wc -l | tr -d '\n')
  echo "$N|files_mention=$ALL|verdictshaped=$VF|name_index=$NAME" >> "$L/m7.tsv"
  printf '%s\n' "$FILES" | while read -r F; do
    [ -z "$F" ] && continue
    grep -q "^$N|$F|" "$L/m3-list.txt" 2>/dev/null || echo "$N|$F" >> "$L/m7-offindex.txt"
  done
done < "$L/roster52.txt"
echo "M7=$(wc -l < "$L/m7.tsv") OFF=$(grep -c . "$L/m7-offindex.txt")"
echo "ZERO_NAMEIDX:"; awk -F'|' '$4=="name_index=0"{printf "%s ",$1}' "$L/m7.tsv"; echo
echo "OFFINDEX_PER_TICKET:"; cut -d'|' -f1 "$L/m7-offindex.txt" | sort | uniq -c | tr '\n' ' '; echo
