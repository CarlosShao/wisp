#!/bin/bash
# M8: refined role census on the name-index candidates: english+chinese markers + leg-id convention.
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
ACC_RE='非实现者|对抗验收|验收腿|验收程|验收方|裁决腿|裁决表|终裁|独立复算|独立|acceptance|adversarial|acceptor|verdict|recheck|独立验收|缺口审计|复核'
IMP_RE='实现腿|实现程|落地腿|落地写腿|写码腿|写码位|实现件|实现方|实现者|self-report'
: > "$L/m8.tsv"
awk -F'|' '{print $2}' "$L/m3-list.txt" | sort -u | while read -r F; do
  [ -f "$F" ] || { echo "$F|MISSING" >> "$L/m8.tsv"; continue; }
  B=$(basename "$F")
  H=$(head -24 "$F")
  HA=$(printf '%s\n' "$H" | grep -ciE "$ACC_RE"); HI=$(printf '%s\n' "$H" | grep -ciE "$IMP_RE")
  NA=$(echo "$B" | grep -ciE "$ACC_RE")
  LEGV=$(echo "$B" | grep -oE '\-v[0-9]+|\-a[0-9]+' | head -1)
  LEGR=$(echo "$B" | grep -oE '\-r[0-9]+' | head -1)
  ROLE=UNK; [ "$NA" -gt 0 ] && ROLE=ACCNAME; [ "$HA" -gt 0 ] && ROLE=ACCHEAD
  [ "$ROLE" = UNK ] && [ "$HI" -gt 0 ] && ROLE=IMPHEAD
  [ "$ROLE" = UNK ] && [ -n "$LEGV" ] && ROLE=VLEG
  [ "$ROLE" = UNK ] && [ -n "$LEGR" ] && ROLE=RLEG
  printf '%s|%s|%s|%s|%s\n' "$F" "$ROLE" "hacc=$HA" "himp=$HI" "$(grep -m1 '^#' "$F" | cut -c1-40)" >> "$L/m8.tsv"
done
echo "M8=$(wc -l < "$L/m8.tsv")"; cut -d'|' -f2 "$L/m8.tsv" | sort | uniq -c
awk -F'|' 'NR==FNR{role[$1]=$2; next}
 { n=$1; f=$2; r=(f in role)?role[f]:"MISSING";
   cnt[n]++; if(r ~ /^ACC|^VLEG/) acc[n]++; if(r ~ /^IMP/) imp[n]++; if(r=="UNK"||r=="MISSING") unk[n]++;
 }
 END{ while((getline line < LS)>0){ split(line,b,"|"); if(b[1]=="")continue;
   printf "%s|NAMEACC=%d|NAMEIMP=%d|NAMEUNK=%d|NAMEFILES=%d\n", b[1], acc[b[1]]+0, imp[b[1]]+0, unk[b[1]]+0, cnt[b[1]]+0 } }' \
 LS="$L/passD.tsv" "$L/m8.tsv" "$L/m3-list.txt" > "$L/m8-roll.tsv"
echo "M8ROLL=$(wc -l < "$L/m8-roll.tsv")"
echo "<< no acc-shaped name file:"; awk -F'|' '$2=="NAMEACC=0"{printf "%s ",$1}' "$L/m8-roll.tsv"; echo
echo "<< has imp-shaped:"; awk -F'|' '$3!="NAMEIMP=0"{printf "%s(%s) ",$1,$3}' "$L/m8-roll.tsv"; echo
echo "<< has unk:"; awk -F'|' '$4!="NAMEUNK=0"{printf "%s(%s) ",$1,$4}' "$L/m8-roll.tsv"; echo
