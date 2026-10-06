#!/bin/bash
# M9: final join -> one compact ASCII row per ticket with tier + off-index verdict-shaped evidence.
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
# role lookup
awk -F'|' '{print $1"|"$2}' "$L/m8.tsv" > "$L/m9-role.txt"
# per-ticket: name-index acc/imp/unk
awk -F'|' '
NR==FNR{ split($0,a,"|"); role[a[1]]=a[2]; next }
{ n=$1; f=$2; r=(f in role)?role[f]:"MISSING";
  cnt[n]++;
  if(r ~ /^ACC|^VLEG/) acc[n]++;
  else if(r ~ /^IMP|^RLEG/) imp[n]++;
  else unk[n]++;
}
END{ for(k in cnt) print k"|"acc[k]+0"|"imp[k]+0"|"unk[k]+0"|"cnt[k]+0 }' "$L/m9-role.txt" "$L/m3-list.txt" > "$L/m9-nameidx.txt"
# off-index verdict-shaped files per ticket
: > "$L/m9-offidx.txt"
while read -r N; do
  [ -z "$N" ] && continue
  VF=$(awk -F'|' -v n="$N" '$1==n && $2=="OFFINDEX"{print $3}' "$L/m7-offindex.txt" | grep -ciE 'adversarial|accept|verdict|recheck|audit|裁决|验收|终裁|归口|backfill|gap|缺口')
  TOT=$(awk -F'|' -v n="$N" '$1==n && $2=="OFFINDEX"{print $3}' "$L/m7-offindex.txt" | sort -u | wc -l)
  echo "$N|offv=$VF|offtot=$TOT" >> "$L/m9-offidx.txt"
done < "$L/roster52.txt"
# final rows
: > "$L/m9-rows.txt"
while IFS='|' read -r N REST; do
  [ -z "$N" ] && continue
  [ -n "$N" ] || continue
  NA=$(awk -F'|' -v n="$N" '$1==n{print $2}' "$L/m9-nameidx.txt"); NI=$(awk -F'|' -v n="$N" '$1==n{print $3}' "$L/m9-nameidx.txt")
  NU=$(awk -F'|' -v n="$N" '$1==n{print $4}' "$L/m9-nameidx.txt"); NT=$(awk -F'|' -v n="$N" '$1==n{print $5}' "$L/m9-nameidx.txt")
  OV=$(awk -F'|' -v n="$N" '$1==n{print $2}' "$L/m9-offidx.txt" | cut -d= -f2)
  OT=$(awk -F'|' -v n="$N" '$1==n{print $3}' "$L/m9-offidx.txt" | cut -d= -f2)
  FILE=$(echo "$REST" | cut -d'|' -f1)
  UN=$(echo "$REST" | grep -oE 'unind=[0-9]+' | cut -d= -f2); AN=$(echo "$REST" | grep -oE 'anyind=[0-9]+' | cut -d= -f2)
  DN=$(echo "$REST" | grep -oE 'done=[0-9]+' | cut -d= -f2); BX=$(echo "$REST" | grep -oE 'boxes=[0-9]+' | cut -d= -f2)
  PV=$(echo "$REST" | grep -oE 'pl_verdict=[0-9]+' | cut -d= -f2); PI=$(echo "$REST" | grep -oE 'pl_impl=[0-9]+' | cut -d= -f2)
  RF=$(echo "$REST" | grep -oE 'refmd=[0-9]+' | cut -d= -f2)
  if [ "${NA:-0}" -gt 0 ]; then T=A; elif [ "${OV:-0}" -gt 0 ]; then T=B; elif [ "${NI:-0}" -gt 0 ] || [ "${NU:-0}" -gt 0 ]; then T=C; elif [ "${PI:-0}" -gt 0 ] || [ "${RF:-0}" -gt 0 ]; then T=D; else T=E; fi
  printf '%s|%s|tier=%s|nA=%s nI=%s nU=%s nT=%s|offv=%s offtot=%s|un=%s any=%s done=%s boxes=%s|pv=%s pi=%s ref=%s|%s\n' \
    "$N" "$FILE" "$T" "${NA:-0}" "${NI:-0}" "${NU:-0}" "${NT:-0}" "${OV:-0}" "${OT:-0}" "${UN:-?}" "${AN:-?}" "${DN:-?}" "${BX:-?}" "${PV:-?}" "${PI:-?}" "${RF:-?}" "$L" >> "$L/m9-rows.txt"
done < <(awk -F'|' '{n=$1; rest=""; for(i=2;i<=NF;i++) rest=rest (i>2?"|":"") $i; print n"|"rest}' "$L/passD.tsv" | while IFS='|' read -r N F; do
   PD=$(grep -m1 "^$N|" .scratch/wisp/probes/pool-validity/5g2/logs/passD.tsv); echo "$PD"; done)
echo "ROWS=$(grep -c . "$L/m9-rows.txt")"
grep -oE 'tier=[A-Z]' "$L/m9-rows.txt" | sort | uniq -c
