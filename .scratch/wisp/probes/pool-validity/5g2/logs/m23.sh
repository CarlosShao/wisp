#!/bin/bash
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
echo "== 1: 84 off-index files: do any claim to BE ticket 84's acceptance table? =="
for F in $(awk -F'|' '$1==84{print $2}' "$L/m7-offindex.txt" | sort -u); do
  H=$(grep -m1 '^#' "$F" 2>/dev/null | cut -c1-60)
  M=$(grep -cE '票 ?84' "$F" 2>/dev/null)
  echo "84|$F|lines=$(wc -l < "$F")|mentions84=$M|$H"
done | head -12
echo "== 2: 184 dispatch file exists? =="
ls -1 .scratch/wisp/dispatches/ 2>/dev/null | grep -c '184'
echo "== 3: legend-aware net self-report cells =="
: > "$L/m23-selfcells.tsv"
while read -r N; do
  K=$(echo "$N" | sed -E 's/^0*//'); REAL=0; DET=""
  for F in $(awk -F'|' -v n="$K" '$1==n{print substr($0,index($0,"|")+1)}' "$L/k-accfiles.txt" | tr ' ' '\n' | grep -E '^(docs|\.scratch)'); do
    [ -f "$F" ] || continue
    S=$(grep -coE '仅自述|仅代理自述|不背书' "$F")
    LEG=$(grep -cE '图例|〔独立复现〕／|〔独立复现〕 /' "$F")
    NET=$((S - LEG)); [ "$NET" -lt 0 ] && NET=0
    [ "$NET" -gt 0 ] && { REAL=$((REAL+NET)); DET="$DET $(basename "$F"):$NET"; }
  done
  echo "$K|net_self=$REAL|$DET" >> "$L/m23-selfcells.tsv"
done < "$L/roster52.txt"
echo "NETSELF_GT0: $(awk -F'|' '{split($2,a,"="); if(a[2]+0>0) printf "%s(%s) ",$1,substr(a[2],0,0)substr($0,index($0,"net_self=")+9,3)}' "$L/m23-selfcells.tsv" | tr -s ' ')"
echo "== 4: per-ticket acceptance artifact count + which dir =="
awk -F'|' '{split($2,a,"="); if(a[2]+0==0) printf "%s_ZEROACCFILE ",$1}' "$L/m22-selfgrade.tsv"; echo
echo "== 5: 152/154 high offv: are those cross-refs or real filings? =="
for N in 152 154; do echo "-- $N top offindex files"; awk -F'|' -v n="$N" '$1==n{print $2}' "$L/m7-offindex.txt" | grep -E '^docs/evidence' | head -4; done
