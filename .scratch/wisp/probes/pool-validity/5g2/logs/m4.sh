#!/bin/bash
# M4: role census over the M3 filename-index candidates (97 files).
# Markers are read from the whole file, plus the header block separately.
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
cut -d'|' -f2 "$L/m3-list.txt" | sort -u > "$L/m4-files.txt"
: > "$L/m4.tsv"
while read -r F; do
  [ -f "$F" ] || { echo "$F|MISSING" >> "$L/m4.tsv"; continue; }
  LN=$(wc -l < "$F")
  HB=$(head -20 "$F")
  H_ACC=$(printf '%s\n' "$HB" | grep -cE '非实现者|对抗验收|验收腿|验收程|验收方|裁决腿|裁决表|终裁|独立复算|acceptor|independent')
  H_IMP=$(printf '%s\n' "$HB" | grep -cE '实现腿|实现程|落地腿|落地写腿|写码腿|写码位|实现件|实现方|实现者')
  W_ACC=$(grep -cE '非实现者|对抗验收|验收腿|验收程|裁决腿|终裁|acceptor' "$F")
  W_IMP=$(grep -cE '实现腿|实现程|落地腿|写码位|写码腿|实现件' "$F")
  VERD=$(grep -ohE 'PASS WITH CONDITIONS|PASS\b|FAIL\b|不通过|退回|拒绝' "$F" | sort | uniq -c | tr -d '\n' | cut -c1-45)
  HT=$(grep -m1 '^#' "$F" | cut -c1-60)
  printf '%s|ln=%s|hacc=%s|himp=%s|wacc=%s|wimp=%s|%s|%s\n' "$F" "$LN" "$H_ACC" "$H_IMP" "$W_ACC" "$W_IMP" "$VERD" "$HT" >> "$L/m4.tsv"
done < "$L/m4-files.txt"
echo "M4=$(wc -l < "$L/m4.tsv")"
awk -F'|' '{split($3,a,"=");split($4,b,"="); if(a[2]>0) acc++; else if(b[2]>0) imp++; else oth++} END{print "header-acc="acc" header-imp="imp" header-neither="oth}' "$L/m4.tsv"
