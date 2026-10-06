#!/bin/bash
# M17: final per-ticket table -> logs/m17-table.tsv (compact, full paths, ASCII only).
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
# add tier column from m9-rows into a lookup
awk -F'|' '{ for(i=1;i<=NF;i++) if($i ~ /^tier=/) print $1"|tier="substr($i,6) }' "$L/m9-rows.txt" > "$L/m17-tier.txt"
: > "$L/m17-table.tsv"
sort -n -u "$L/roster52.txt" | while read -r N; do
  FILE=$(ls .scratch/wisp/issues/"$N"-*.md 2>/dev/null | head -1)
  BASE=$(basename "$FILE")
  TIER=$(grep -m1 "^$N|" "$L/m17-tier.txt" | cut -d'|' -f2)
  NUM=$(echo "$N" | sed 's/^0*//')
  [ -z "$TIER" ] && TIER="tier=$(grep -m1 "^$NUM|" "$L/m17-tier.txt" | cut -d'|' -f2)"
  COUNTS=$(grep -m1 "^$NUM|" "$L/m9-rows.txt" | grep -oE 'nA=[0-9]+|nI=[0-9]+|nU=[0-9]+|offv=[0-9]+' | tr '\n' ' ')
  UNCHECK=$(grep -m1 "^$NUM|" "$L/m16-strict.txt" | cut -d'|' -f2)
  [ -z "$UNCHECK" ] && UNCHECK=$(grep -m1 "^$N|" "$L/m16-strict.txt" | cut -d'|' -f2)
  UNIND=$(grep -m1 "^$NUM|" "$L/m14-boxes.txt" | grep -oE 'unind=[0-9]+' | sed 's/unind=//')
  ANYIND=$(grep -m1 "^$NUM|" "$L/m14-boxes.txt" | grep -oE 'anyind=[0-9]+' | sed 's/anyind=//')
  ACCF=$(awk -F'|' -v n="$NUM" '$1==n{print $2}' "$L/m3-list.txt" | while read -r F; do
      R=$(awk -F'|' -v f="$F" '$1==f{print $2}' "$L/m12-role.txt")
      case "$R" in ACCHEAD) echo "$F";; VLEG) echo "$F";; esac
    done | sort -u | tr '\n' ',' )
  IMPF=$(awk -F'|' -v n="$NUM" '$1==n{print $2}' "$L/m3-list.txt" | while read -r F; do
      R=$(awk -F'|' -v f="$F" '$1==f{print $2}' "$L/m12-role.txt")
      case "$R" in IMPHEAD|RLEG) echo "$F";; esac
    done | sort -u | tr '\n' ',')
  LED=$(grep -m1 "^$NUM|" "$L/m11-ledger.tsv" | cut -d'|' -f2,3)
  echo "$NUM|$BASE|$TIER|$COUNTS|unchecked_strict=$UNCHECK unind=$UNIND anyind=$ANYIND|$LED|$ACCF|$IMPFL" >> "$L/m17-table.tsv"
done
echo "M17=$(wc -l < "$L/m17-table.tsv")"
awk -F'|' '{ if(length($0)>240) $0=substr($0,1,240); print $1" len="length($0) }' "$L/m17-table.tsv" | tail -3
echo "rows52=$(grep -c . "$L/m17-table.tsv")"
