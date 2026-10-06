#!/bin/bash
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
: > "$L/final-b.tsv"
while read -r N; do
  K=$(echo "$N" | sed -E 's/^0*//')
  ACC=$(awk -F'|' -v n="$K" '$1==n{print substr($0,index($0,"|")+1)}' "$L/k-accfiles.txt" | tr ' ' '\n' | grep -E '^(docs|\.scratch)')
  IMP=$(awk -F'|' -v n="$K" '$1==n{print substr($0,index($0,"|")+1)}' "$L/k-impfiles.txt" | tr ' ' '\n' | grep -E '^(docs|\.scratch)')
  NA=$(printf '%s\n' "$ACC" | grep -c .); NI=$(printf '%s\n' "$IMP" | grep -c .)
  ID=0; IDEN=""
  for F in $ACC; do [ -f "$F" ] || continue
    a=$(grep -coE '编排者|非实现者|验收人|裁决人|本程代号|派单 ' "$F")
    [ "$a" -gt 0 ] && { ID=$((ID+1)); IDEN="$IDEN $(basename "$F")"; }
  done
  if [ "$NA" -eq 0 ] && [ "$NI" -eq 0 ]; then B="NONE";
  elif [ "$NA" -eq 0 ] && [ "$NI" -gt 0 ]; then B="IMPL-ONLY";
  elif [ "$ID" -eq 0 ]; then B="ACC-NO-IDENTITY";
  else B="ACC-IDENTIFIED"; fi
  echo "$K|B=$B|acc=$NA|imp=$NI|identity_files=$ID" >> "$L/final-b.tsv"
done < "$L/roster52.txt"
for X in ACC-IDENTIFIED ACC-NO-IDENTITY IMPL-ONLY NONE; do echo "$X: $(grep -c "|$B=$X|" /dev/null 2>/dev/null; awk -F'|' -v x="B=$X" '$2==x{printf "%s ",$1}' $L/final-b.tsv)"; done
echo "== total =="; awk -F'|' '{split($2,a,"="); s[a[2]]++} END{for(k in s) print k" "s[k]}' $L/final-b.tsv
