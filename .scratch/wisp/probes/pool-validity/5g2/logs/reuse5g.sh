#!/bin/bash
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
awk 'BEGIN{t=""} /^==== TICKET /{ if(t!=""){print t, c} t=$3; c=0; next } /lines=/{c++} END{if(t!="")print t, c}' \
  .scratch/wisp/probes/pool-validity/5g/logs/roles.tsv > "$L/5g-counts.txt"
echo "rows=$(wc -l < "$L/5g-counts.txt")"
echo "zero-file tickets:"; awk '$2+0==0{printf "%s ", $1}' "$L/5g-counts.txt"; echo
echo "nonzero=$(awk '$2+0>0{n++}END{print n+0}' "$L/5g-counts.txt")"
