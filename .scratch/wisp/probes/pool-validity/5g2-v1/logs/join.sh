#!/bin/bash
# Join: ticket | tier | m31(head24 sig) | m24(signed_files) | m33(ident) | m32(orch_self/leg) | era(orch/leg)
cd "D:/work/workspace/projects plans/Wisp" || exit 1
G=.scratch/wisp/probes/pool-validity/5g2/logs
L=.scratch/wisp/probes/pool-validity/5g2-v1/logs
: > "$L/join.tsv"
while read -r T TIER; do
  m31=$(awk -F'|' -v n="$T" '$1+0==n{ $1=""; sub(/^ *\|?/,""); print substr($0, index($0,"sig=")) }' "$G/m31-sign.txt" | cut -c1-40)
  [ -z "$m31" ] && m31="(no-line)"
  m24=$(awk -F'|' -v n="$T" '$1+0==n{print $2}' "$G/m24-sign.tsv")
  m33=$(awk -F'|' -v n="$T" '$1+0==n{print $2"/"$3"/"$4}' "$G/m33-ident.tsv")
  m32=$(awk -F'|' -v n="$T" '$1+0==n{print $3" "$4}' "$G/m32-who.tsv")
  era=$(awk -v n="$T" '$1+0==n{print $2" "$3}' "$G/era.tsv")
  printf '%s\t%s\t%s\t%s\t%s\t%s\t%s\n' "$T" "$TIER" "$m31" "$m24" "$m33" "$m32" "$era" >> "$L/join.tsv"
done < "$L/sec2-tiers.tsv"
echo "rows=$(wc -l < "$L/join.tsv")"
echo "--- 待定 with m32 authorship (orch_self>0 or leg>0) ---"
awk -F'\t' '$2=="待定"{ if ($6 ~ /orch_self_files=[1-9]/ || $6 ~ /leg_files=[1-9]/) c++ } END{print c+0}' "$L/join.tsv"
echo "--- 待定 total ---"
awk -F'\t' '$2=="待定"{c++} END{print c+0}' "$L/join.tsv"
echo "--- 待定 with m32 NEITHER (acc>0 but orch=0 leg=0) ---"
awk -F'\t' '$2=="待定" && $6 ~ /orch_self_files=0/ && $6 ~ /leg_files=0/{c++} END{print c+0}' "$L/join.tsv"
echo "--- 甲 with m31 NOSIG (should be 0) ---"
awk -F'\t' '$2=="甲" && $3 ~ /NOSIG/{print $1}' "$L/join.tsv" | tr '\n' ' '
echo "--- 待定 with era orch/leg nonzero ---"
awk -F'\t' '$2=="待定" && ($7 ~ /orchstyle=[1-9]/ || $7 ~ /legstyle=[1-9]/){c++} END{print c+0}' "$L/join.tsv"
