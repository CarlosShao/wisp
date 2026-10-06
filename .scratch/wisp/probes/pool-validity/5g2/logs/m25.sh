#!/bin/bash
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
echo "== A: distinct bracket tokens across the 52 done tickets =="
: > $L/m25-tokens.txt
while read -r N; do K=$(echo $N|sed -E 's/^0*//'); F=$(ls .scratch/wisp/issues/$K-*.md|head -1)
  grep -ohE '^[[:space:]]*[-*] \[[^]]{0,3}\]' "$F" | sed -E 's/^[[:space:]]*//' | sort | uniq -c | tr '\n' ';' >> $L/m25-tokens.txt
  echo " << $K" >> $L/m25-tokens.txt
done < $L/roster52.txt
grep -ohE '\[[^]]{0,3}\]' $L/m25-tokens.txt | sort | uniq -c | sort -rn
echo "== B: AC coverage: distinct AC# on face vs in its acc files =="
: > $L/m25-ac.tsv
while read -r N; do K=$(echo $N|sed -E 's/^0*//'); F=$(ls .scratch/wisp/issues/$K-*.md|head -1)
  FACE=$(grep -ohE 'AC ?#[0-9]+[a-z]?' "$F" | grep -oE '[0-9]+' | sort -n -u | tail -1)
  FCB=$(grep -cE '^[[:space:]]*[-*] \[' "$F")
  AF=$(awk -F'|' -v n="$K" '$1==n{print substr($0,index($0,"|")+1)}' $L/k-accfiles.txt | tr ' ' '\n' | grep -E '^(docs|\.scratch)')
  COV=0; for X in $AF; do [ -f "$X" ] || continue; V=$(grep -ohE 'AC ?#[0-9]+[a-z]?' "$X" | grep -oE '^[0-9]+' | sort -n -u | tr '\n' ',' ); COV="$COV$V"; done
  UNIQ=$(echo "$COV" | tr ',' '\n' | grep -E '^[0-9]+$' | sort -n -u | tr '\n' ',')
  echo "$K|face_maxAC=$FACE|face_boxes=$FCB|acc_ACset=$UNIQ" >> $L/m25-ac.tsv
done < $L/roster52.txt
echo "-- tickets whose acc AC set top < face max AC (coverage gap candidates):"
awk -F'|' '{split($2,a,"=");split($4,c,"="); m=a[2]+0; if($3!=""){ } ; s=$5; gsub("AC.*","",$s); n=0; last=0; split($5,x,","); for(i in x){ if(x[i]!=""){n++; if(x[i]+0>last) last=x[i]+0} } printf "%s F=%s A=%s %s\n",$1,m,last, (last<m?"GAP":"")}' $L/m25-ac.tsv | grep GAP
echo "== C: tickets with no Progress log section =="
while read -r N; do K=$(echo $N|sed -E 's/^0*//'); F=$(ls .scratch/wisp/issues/$K-*.md|head -1); PL=$(grep -cE '^#{2,3} *(Progress log|进展|进度)' "$F"); [ "$PL" = 0 ] && echo "NOPL $K $F"; done < $L/roster52.txt
echo "== D: tickets with no Evidence/凭据/归口 section on face =="
while read -r N; do K=$(echo $N|sed -E 's/^0*//'); F=$(ls .scratch/wisp/issues/$K-*.md|head -1); E=$(grep -cE '^#{2,3} .*(Evidence|凭据|裁决|验收|归口)' "$F"); [ "$E" = 0 ] && echo "NOEVIDSEC $K acc_cited=$(grep -cE 'docs/evidence' $F)"; done < $L/roster52.txt
