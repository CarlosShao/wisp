#!/bin/bash
# M5: join -> per-ticket classification (A1 from m3+m4 role census, A2 from ticket progress log).
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
: > "$L/m5.tsv"
while read -r N; do
  [ -z "$N" ] && continue
  # A1: evidence files whose name carries this ticket number, role-split by header marker
  ACC=0; IMP=0; UNK=0; S1F=0; OTHF=0
  while IFS='|' read -r F REST; do
    [ -z "$F" ] && continue
    LN=$(echo "$REST" | grep -oE 'ln=[0-9]+' | cut -d= -f2)
    ROW=$(grep -m1 -F "$F|" "$L/m4.tsv")
    HA=$(echo "$ROW" | grep -oE 'hacc=[0-9]+' | cut -d= -f2)
    HI=$(echo "$ROW" | grep -oE 'himp=[0-9]+' | cut -d= -f2)
    case "$F" in docs/evidence/s1/*) S1F=$((S1F+1));; *) OTHF=$((OTHF+1));; esac
    if [ "${HA:-0}" != "0" ]; then ACC=$((ACC+1));
    elif [ "${HI:-0}" != "0" ]; then IMP=$((IMP+1));
    else UNK=$((UNK+1)); fi
  done < <(grep "^$N|" "$L/m3-list.txt")
  # A2: own ticket progress-log markers
  PD=$(grep -m1 "^$N|" "$L/passD.tsv")
  echo "$N|$PD|A1acc=$ACC|A1imp=$IMP|A1unk=$UNK|s1=$S1F|other=$OTHF" >> "$L/m5.tsv"
done < "$L/roster52.txt"
echo "M5=$(wc -l < "$L/m5.tsv")"
echo "--- tickets with ZERO acc-shaped files:"; awk -F'|' '{if($0 ~ /A1acc=0[^0-9]/ || $0 ~ /A1acc=0\|/) print $1}' "$L/m5.tsv" | tr '\n' ' '; echo
echo "--- tickets with UNK>0:"; awk -F'|' '$0 ~ /A1unk=[1-9]/ {print $1}' "$L/m5.tsv" | tr '\n' ' '; echo
echo "--- tickets with IMP>0 && ACC=0:"; awk -F'|' '$0 ~ /A1imp=[1-9]/ && $0 ~ /A1acc=0[^0-9]/ {print $1}' "$L/m5.tsv" | tr '\n' ' '; echo
echo "--- boxes=0 tickets:"; awk -F'|' '$6=="boxes=0"{print $1}' "$L/passD.tsv" | tr '\n' ' '; echo
echo "--- anyind>unind tickets:"; awk -F'|' '{gsub("[a-z]+=","",$4);gsub("[a-z]+=","",$3); if($4+0>$3+0) print $1":"$3"/"$4}' "$L/passD.tsv" | tr '\n' ' '; echo
echo "--- unchecked>0 tickets:"; awk -F'|' '{v=$4; gsub("[a-z]+=","",v); if(v+0>0) print $1}' "$L/passD.tsv" | tr '\n' ' '; echo
