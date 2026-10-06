#!/bin/bash
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
# clean self-grade cell count: markdown table rows whose grade column says self-report only
: > "$L/m26-selfcells.tsv"
while read -r N; do
  K=$(echo "$N" | sed -E 's/^0*//'); SC=0; IC=0; AC=0
  for F in $(awk -F'|' -v n="$K" '$1==n{print substr($0,index($0,"|")+1)}' "$L/k-accfiles.txt" | tr ' ' '\n' | grep -E '^(docs|\.scratch)'); do
    [ -f "$F" ] || continue
    AC=$((AC+1))
    s=$(grep -cE '^\|.*((仅自述|仅代理自述|不背书))' "$F")
    i=$(grep -cE '^\|.*独立复现' "$F")
    SC=$((SC+s)); IC=$((IC+i))
  done
  echo "$K|accfiles=$AC|selfcells=$SC|indep_cells=$IC" >> "$L/m26-selfcells.tsv"
done < "$L/roster52.txt"
echo "SELFCELLS_GT0: $(awk -F'|' '{split($3,b,"="); if(b[2]+0>0) printf "%s(%s) ",$1,b[2]}' $L/m26-selfcells.tsv)"
echo "ZEROBOTH(accfiles>0 && selfcells==0 && indep==0): $(awk -F'|' '{split($4,c,"=");split($3,b,"=");split($2,a,"="); if(a[2]+0>0&&b[2]+0==0&&c[2]+0==0) printf "%s ",$1}' $L/m26-selfcells.tsv)"
echo "TOTALSELF=$(awk '{split($3,a,"=");s+=a[2]}END{print s}' $L/m26-selfcells.tsv) TOTALINDEP=$(awk '{split($4,a,"=");s+=a[2]}END{print s}' $L/m26-selfcells.tsv)"
# one-line-per-ticket compact reference for authoring section 2
: > "$L/m26-ref.txt"
while read -r N; do
  K=$(echo "$N" | sed -E 's/^0*//')
  ROW=$(grep -m1 "^$K|" "$L/m19-table.tsv")
  BASE=$(echo "$ROW" | cut -d'|' -f2)
  UN=$(awk -F'|' -v n="$K" '$1==n{split($3,x,"="); print x[2]}' "$L/k-strict.txt")
  BX=$(awk -F'|' -v n="$K" '$1==n{split($4,y,"="); print y[2]}' "$L/k-strict.txt")
  AF=$(awk -F'|' -v n="$K" '$1==n{print substr($0,index($0,"|")+1)}' "$L/k-accfiles.txt" | tr ' ' '\n' | grep -E '^(docs|\.scratch)' | sed -E 's|docs/evidence/s1/|s1/|; s|\.scratch/wisp/probes/|P/|' | paste -sd, -)
  IF=$(awk -F'|' -v n="$K" '$1==n{print substr($0,index($0,"|")+1)}' "$L/k-impfiles.txt" | tr ' ' '\n' | grep -E '^(docs|\.scratch)' | sed -E 's|docs/evidence/s1/|s1/|; s|\.scratch/wisp/probes/|P/|' | paste -sd, -)
  SC=$(grep -m1 "^$K|" "$L/m26-selfcells.tsv" | cut -d'|' -f3,4 | paste -sd' ' -)
  echo "$K|$BASE|cbx=$UN/$BX|acc=$AF|imp=$IF|$SC" >> "$L/m26-ref.txt"
done < <(sort -n -u "$L/roster52.txt")
wc -l < "$L/m26-ref.txt"
