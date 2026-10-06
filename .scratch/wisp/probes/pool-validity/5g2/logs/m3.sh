#!/bin/bash
# M3: independent filename-index census. For each roster number, find every md file in
# scope whose BASENAME carries that number as a leading token (N- / N_ / -N- / N at end).
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
find docs/evidence .scratch/wisp/probes docs/reports -name '*.md' 2>/dev/null \
  | grep -v '^\.scratch/wisp/probes/pool-validity/5g/' \
  | grep -v '^\.scratch/wisp/probes/232/' \
  | grep -v -e '^\.scratch/wisp/probes/frontend/' -e '^\.scratch/wisp/probes/design/' \
  | sort > "$L/m3-scope-files.txt"
echo "scope_files=$(wc -l < "$L/m3-scope-files.txt")"
: > "$L/m3.tsv"
while read -r N; do
  [ -z "$N" ] && continue
  FILES=$(grep -E "(^|/)${N}[-_.]|-${N}[-_.]|(^|/)${N}\.md" "$L/m3-scope-files.txt")
  CNT=$(printf '%s\n' "$FILES" | grep -c .)
  S1=$(printf '%s\n' "$FILES" | grep -c '^docs/evidence/s1/')
  OTHER=$(printf '%s\n' "$FILES" | grep -vc '^docs/evidence/s1/')
  echo "$N|files=$CNT|docs_evidence_s1=$S1|other=$OTHER" >> "$L/m3.tsv"
  printf '%s\n' "$FILES" | while read -r F; do
    [ -z "$F" ] && continue
    echo "$N|$F|ln=$(wc -l < "$F")" >> "$L/m3-list.txt.tmp"
  done
done < "$L/roster52.txt"
mv -f "$L/m3-list.txt.tmp" "$L/m3-list.txt" 2>/dev/null
echo "m3tsv=$(wc -l < "$L/m3.tsv") m3list=$(wc -l < "$L/m3-list.txt")"
echo "TOTALFILES=$(cut -d'|' -f2 "$L/m3-list.txt" | sort -u | wc -l)"
