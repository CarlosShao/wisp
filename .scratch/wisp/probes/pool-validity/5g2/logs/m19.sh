#!/bin/bash
# M19b: final per-ticket table with numeric-normalised keys -> logs/m19-table.tsv
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
norm(){ sed -E 's/^0*([0-9]+)\|/\1|/'; }
awk -F'|' '{for(i=1;i<=NF;i++) if($i ~ /^tier=/){split($i,x,"="); print $1"|"x[2]}}' "$L/m9-rows.txt" | norm > "$L/k-tier.txt"
awk -F'|' '{ printf "%s|%s|%s|%s|%s\n", $1, $4, $5, $8, $9 }' "$L/m9-rows.txt" | norm > "$L/k-nameidx.txt"
norm < "$L/m14-boxes.txt" > "$L/k-boxes.txt"
norm < "$L/m16-strict.txt" > "$L/k-strict.txt"
norm < "$L/m16-verdict.tsv" > "$L/k-verdict.txt"
: > "$L/k-files.txt"; : > "$L/k-accfiles.txt"; : > "$L/k-impfiles.txt"
while read -r N; do
  K=$(echo "$N" | sed -E 's/^0*//')
  FILE=$(ls .scratch/wisp/issues/"$K"-*.md 2>/dev/null | head -1)
  echo "$K|$(basename "$FILE")" >> "$L/k-files.txt"
  ACC=$(awk -F'|' -v n="$K" '$1==n{print $2}' "$L/m3-list.txt" | while read -r F; do
        R=$(awk -F'|' -v f="$F" '$1==f{print $2}' "$L/m12-role.txt"); case "$R" in ACCHEAD|VLEG) echo "$F";; esac
      done | sort -u | tr '\n' ' ')
  IMP=$(awk -F'|' -v n="$K" '$1==n{print $2}' "$L/m3-list.txt" | while read -r F; do
        R=$(awk -F'|' -v f="$F" '$1==f{print $2}' "$L/m12-role.txt"); case "$R" in IMPHEAD|RLEG) echo "$F";; esac
      done | sort -u | tr '\n' ' ')
  echo "$K|$ACC" >> "$L/k-accfiles.txt"; echo "$K|$IMP" >> "$L/k-impfiles.txt"
done < "$L/roster52.txt"
rm -f "$L/m19-table.tsv"
awk -v OUT="$L/m19-table.tsv" -f "$L/join.awk" "$L/k-tier.txt" "$L/k-nameidx.txt" "$L/k-boxes.txt" "$L/k-strict.txt" "$L/k-verdict.txt" "$L/k-files.txt" "$L/k-accfiles.txt" "$L/k-impfiles.txt"
echo "M19=$(grep -c . "$L/m19-table.tsv")"
awk -F'|' '{for(i=1;i<=NF;i++) if($i ~ /^tier=/) print substr($i,6)}' "$L/m19-table.tsv" | sort | uniq -c
echo "NOACC=$(awk -F'|' '$8=="accfiles="{printf "%s ",$1}' $L/m19-table.tsv)"
