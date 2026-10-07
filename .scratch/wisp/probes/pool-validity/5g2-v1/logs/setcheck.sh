#!/bin/bash
# Ruler: which of the 52-line files in 5g2/logs carry the SAME ticket-number set as roster52.txt ?
# f52-all.txt lines look like: "52 <path>"
G=".scratch/wisp/probes/pool-validity/5g2/logs"
L=".scratch/wisp/probes/pool-validity/5g2-v1/logs"
sort -n "$G/roster52.txt" | uniq > "$L/set-roster.txt"
wc -l < "$L/set-roster.txt" > "$L/set-roster-n.txt"
: > "$L/setcompare2.tsv"
while read -r cnt p; do
  [ -f "$p" ] || continue
  name=$(basename "$p")
  a=$(awk -F'|' '{v=$1; gsub(/[^0-9]/,"",v); if(v!="") print v+0}' "$p" | sort -n | uniq)
  na=$(printf '%s\n' "$a" | grep -c '^[0-9]')
  b=$(awk -F'|' '{v=$NF; gsub(/[^0-9]/,"",v); if(v!="") print v+0}' "$p" | sort -n | uniq)
  nb=$(printf '%s\n' "$b" | grep -c '^[0-9]')
  c=$(grep -oE '[0-9]{1,3}$' "$p" | sort -n | uniq)
  nc=$(printf '%s\n' "$c" | grep -c '^[0-9]')
  best="none"; nbest=0; verdict="NO_TICKET_COL"
  for pair in "A:$na:$a" "B:$nb:$b" "C:$nc:$c"; do
    tag=${pair%%:*}; r1=${pair#*:}; k=${r1%%:*}; body=${r1#*:}
    if [ "$k" -ge 40 ] && [ "$k" -ge "$nbest" ]; then
      best=$tag; nbest=$k
      if printf '%s\n' "$body" | diff -q - "$L/set-roster.txt" >/dev/null; then verdict="SAME52"; else verdict="DIFF"; fi
    fi
  done
  printf '%s\t%s\t%s\t%s\n' "$name" "$nbest" "$best" "$verdict" >> "$L/setcompare2.tsv"
done < <(awk '$1==52' "$L/f52-all.txt")
echo "TOTAL 52-line files tested: $(wc -l < "$L/setcompare2.tsv")"
echo "SAME52: $(grep -c 'SAME52' "$L/setcompare2.tsv")"
echo "DIFF  : $(grep -c 'DIFF' "$L/setcompare2.tsv")"
echo "NO_TICKET_COL: $(grep -c 'NO_TICKET_COL' "$L/setcompare2.tsv")"
echo "--- roster52 unique tickets: $(cat "$L/set-roster-n.txt") ---"
echo "--- DIFF list ---"; grep 'DIFF' "$L/setcompare2.tsv" | cut -f1 | tr '\n' ' '
echo; echo "--- NO_TICKET_COL list ---"; grep 'NO_TICKET_COL' "$L/setcompare2.tsv" | cut -f1 | tr '\n' ' '
