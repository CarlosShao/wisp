#!/bin/bash
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
# B: who signed the acceptance artifacts (named leg id vs orchestrator)
: > $L/m21-acceptor.tsv
while read -r N; do
  K=$(echo $N | sed -E 's/^0*//')
  AF=$(awk -F'|' -v n="$K" '$1==n{print substr($0, index($0,"|")+1)}' $L/k-accfiles.txt)
  NF=0; OC=0; LG=0
  for F in $(echo "$AF" | tr ' ' '\n' | grep -E '^(docs|\.scratch)'); do
    [ -f "$F" ] || continue
    NF=$((NF+1))
    a=$(grep -c '编排者' "$F"); [ "$a" -gt 0 ] && OC=$((OC+1))
    c=$(grep -cE '本程代号|派单' "$F"); [ "$c" -gt 0 ] && LG=$((LG+1))
  done
  echo "$K|accfiles=$NF|orch_named_files=$OC|dispatch_named_files=$LG" >> $L/m21-acceptor.tsv
done < $L/roster52.txt
echo "B1_no_dispatch_ref=$(awk -F'|' '$4 ~ /=0$/{printf \"%s \",$1}' $L/m21-acceptor.tsv)"
echo "B2_orch_named_only=$(awk -F'|' '$3 ~ /[1-9]$/ && $4 ~ /=0$/{printf \"%s \",$1}' $L/m21-acceptor.tsv)"
echo "B3_accfiles0=$(awk -F'|' '$2=="accfiles=0"{printf \"%s \",$1}' $L/m21-acceptor.tsv)"
echo "== C: 104 / 110 unchecked box lines with context =="
for N in 104 110; do F=$(ls .scratch/wisp/issues/$N-*.md|head -1); echo "-- $N $F"; grep -nE '^[[:space:]]*- \[ \]' "$F" | cut -c1-120; done
echo "== D: thin acceptance artifacts, acceptor identity =="
for F in docs/evidence/s1/74-adversarial-acceptance.md docs/evidence/s1/76-adversarial-acceptance.md docs/evidence/s1/110-adversarial-acceptance.md docs/evidence/s1/67-adversarial-acceptance.md docs/evidence/s1/83-adversarial-acceptance.md; do echo "-- $F lines=$(wc -l < $F)"; grep -m2 -E '验收人|编排者|本程代号|派单|档位' $F | cut -c1-110; done
