#!/bin/bash
# M2 role census over the verdict-shaped candidate files collected by passC.
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
Out=$L/m2.tsv; : > "$Out"
cut -d'|' -f2 "$L/passC-verdictfiles.txt" | grep -E '^(docs|\.scratch)' | sort -u > "$L/m2-files.txt"
echo "uniquefiles=$(wc -l < "$L/m2-files.txt")"
while read -r F; do
  [ -f "$F" ] || { echo "$F|MISSING" >> "$Out"; continue; }
  LN=$(wc -l < "$F")
  RACC=$(grep -cE '非实现者|对抗验收|验收腿|验收程|验收方|裁决腿|裁决表|终裁|acceptor|独立复算|独立' "$F")
  RIMP=$(grep -cE '实现腿|实现程|落地腿|落地写腿|写码腿|写码位|实现者|实现方' "$F")
  LEG=$(grep -oE '[0-9]{2,3}-[a-z]{1,6}[0-9]{0,2}' "$F" | sort -u | tr '\n' ',' | cut -c1-40)
  VERD=$(grep -ohE 'PASS WITH CONDITIONS|PASS|FAIL|不通过|退回' "$F" | sort | uniq -c | tr '\n' ' ' | cut -c1-50)
  HEAD=$(grep -m1 '^#' "$F" | cut -c1-55)
  printf '%s|ln=%s|acc=%s|imp=%s|%s|%s|%s\n' "$F" "$LN" "$RACC" "$RIMP" "$LEG" "$VERD" "$HEAD" >> "$Out"
done < "$L/m2-files.txt"
echo "M2=$(wc -l < "$Out")"
