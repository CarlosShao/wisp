#!/bin/bash
# R1: for every 待定 ticket, does its 名下 acc file head-30 carry an ORCHESTRATOR phrase or a NAMED leg id?
# R2: re-test the 9 "DIFF" 52-line files using first whitespace token as ticket column.
cd "D:/work/workspace/projects plans/Wisp" || exit 1
G=.scratch/wisp/probes/pool-validity/5g2/logs
L=.scratch/wisp/probes/pool-validity/5g2-v1/logs
ORCH='orchestrator|编排者'
LEGID='acceptor-ticket[0-9]+|agent-ticket[0-9]+|[0-9]{2,3}-(v|r|a)[0-9]+|本程代号'
: > "$L/named41.tsv"
awk -F'\t' '$2=="待定"{print $1}' "$L/join.tsv" | while read -r T; do
  files=$(awk -F'|' -v n="$T" '$1+0==n{print substr($0,index($0,"|")+1)}' "$G/k-accfiles.txt" | tr ' ' '\n' | grep -E '^(docs|\.scratch)')
  o=0; l=0; first=""
  for F in $files; do
    [ -f "$F" ] || continue
    h=$(head -30 "$F")
    oc=$(printf '%s\n' "$h" | grep -coE "$ORCH")
    lc=$(printf '%s\n' "$h" | grep -coE "$LEGID")
    [ "$oc" -gt 0 ] && o=1
    [ "$lc" -gt 0 ] && l=1
    [ -z "$first" ] && first=$(printf '%s\n' "$h" | grep -nE "$ORCH|$LEGID" | head -1 | cut -c1-80)
  done
  cat1="NONE"; [ "$o" -eq 1 ] && [ "$l" -eq 1 ] && cat1="ORCH+LEGID"
  [ "$o" -eq 1 ] && [ "$l" -eq 0 ] && cat1="ORCHonly"
  [ "$o" -eq 0 ] && [ "$l" -eq 1 ] && cat1="LEGIDonly"
  printf '%s\t%s\t%s\n' "$T" "$cat1" "$first" >> "$L/named41.tsv"
done
echo "== 待定 named-authorship breakdown (n=41) =="
cut -f2 "$L/named41.tsv" | sort | uniq -c
echo "== ORCHonly list =="; awk -F'\t' '$2=="ORCHonly"{printf "%s ",$1}' "$L/named41.tsv"
echo; echo "== LEGIDonly list =="; awk -F'\t' '$2=="LEGIDonly"{printf "%s ",$1}' "$L/named41.tsv"
echo; echo "== NONE list =="; awk -F'\t' '$2=="NONE"{printf "%s ",$1}' "$L/named41.tsv"
echo
echo "== R2: re-test 9 DIFF files with whitespace-first-token =="
for f in 5g-counts.txt digest.txt era.tsv k-files.txt m11-verdict.tsv m36-clean.txt m5.tsv m9-offsum.txt netself.txt; do
  s=$(awk '{print $1}' "$G/$f" | grep -E '^[0-9]{1,3}$' | sort -n | uniq | diff -q - "$L/set-roster.txt" >/dev/null && echo SAME52 || echo DIFF)
  echo "$f $s"
done
