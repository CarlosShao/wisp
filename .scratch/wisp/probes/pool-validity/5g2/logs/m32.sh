#!/bin/bash
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
: > "$L/m32-who.tsv"
while read -r N; do
  K=$(echo "$N" | sed -E 's/^0*//')
  SELF=0; LEG=0; NA=0
  for F in $(awk -F'|' -v n="$K" '$1==n{print substr($0,index($0,"|")+1)}' "$L/k-accfiles.txt" | tr ' ' '\n' | grep -E '^(docs|\.scratch)'); do
    [ -f "$F" ] || continue
    NA=$((NA+1))
    s=$(head -30 "$F" | grep -coE 'Acceptor:\*\* ?orchestrator|验收人\*\*：编排者|验收人：编排者|编排者亲自|编排者自裁|编排者执行|由编排者')
    g=$(head -30 "$F" | grep -coE '本程代号|验收腿|裁决腿|验收程|非实现者')
    [ "$s" -gt 0 ] && SELF=$((SELF+1))
    [ "$g" -gt 0 ] && LEG=$((LEG+1))
  done
  echo "$K|acc=$NA|orch_self_files=$SELF|leg_files=$LEG" >> "$L/m32-who.tsv"
done < <(sort -n -u "$L/roster52.txt")
echo "ORCHSELF=$(awk -F'|' '{split($3,b,"=");split($4,c,"="); if(b[2]+0>0&&c[2]+0==0) printf "%s ",$1}' $L/m32-who.tsv)"
echo "LEGONLY=$(awk -F'|' '{split($3,b,"=");split($4,c,"="); if(b[2]+0==0&&c[2]+0>0) printf "%s ",$1}' $L/m32-who.tsv)"
echo "BOTH=$(awk -F'|' '{split($3,b,"=");split($4,c,"="); if(b[2]+0>0&&c[2]+0>0) printf "%s ",$1}' $L/m32-who.tsv)"
echo "NEITHER=$(awk -F'|' '{split($2,a,"=");split($3,b,"=");split($4,c,"="); if(a[2]+0>0&&b[2]+0==0&&c[2]+0==0) printf "%s ",$1}' $L/m32-who.tsv)"
echo "NOACC=$(awk -F'|' '$2=="acc=0"{printf "%s ",$1}' $L/m32-who.tsv)"
