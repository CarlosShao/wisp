#!/bin/bash
# M16: (1) strict checkbox census over 94 done tickets, (2) verdict-token extraction for the
# acceptance-shaped named files of the 52, (3) prior-leg disagreement rows, (4) ledger 76 lines.
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
: > "$L/m16-strict.txt"
for F in .scratch/wisp/issues/*-done.md; do
  N=$(basename "$F" | sed -E 's/^([0-9]+)-.*/\1/')
  S=$(grep -cE '^[[:space:]]*- \[[ xX]\]' "$F")
  echo "$N|strictcb=$S" >> "$L/m16-strict.txt"
done
echo "STRICTCB0: $(awk -F'|' '$2=="strictcb=0"{printf "%s ",$1}' $L/m16-strict.txt)"
echo "STRICTCB0_COUNT_ROSTER52: $(while read -r N; do grep -m1 "^$N|" $L/m16-strict.txt | awk -F'|' '{if($2=="strictcb=0") print "Y"}'; done < $L/roster52.txt | grep -c Y)"
: > "$L/m16-verdict.tsv"
while read -r N; do
  [ -z "$N" ] && continue
  PA=0; FA=0; BK=0; NW=0
  for F in $(awk -F'|' -v n="$N" '$1==n{print $2}' "$L/m3-list.txt" | while read -r X; do
      R=$(awk -F'|' -v f="$X" '$1==f{print $2}' "$L/m12-role.txt")
      case "$R" in ACCHEAD|VLEG) echo "$X";; esac
    done); do
    P=$(grep -ocE 'PASS' "$F" 2>/dev/null); P=$(grep -oE 'PASS' "$F" | wc -l)
    F2=$(grep -oE 'FAIL' "$F" | wc -l); B=$(grep -cE '退回|不通过|reject' "$F")
    PA=$((PA+P)); FA=$((FA+F2)); BK=$((BK+B)); NW=$((NW+1))
  done
  echo "$N|accfiles=$NW|pass=$PA|fail=$FA|退回or不通过=$BK" >> "$L/m16-verdict.tsv"
done < "$L/roster52.txt"
echo "LOWPASS: $(awk -F'|' '{split($3,a,"="); if(a[2]+0==0) printf "%s ",$1}' $L/m16-verdict.tsv)"
echo "HASBACK: $(awk -F'|' '{split($5,a,"="); if(a[2]+0>0) printf "%s ",$1}' $L/m16-verdict.tsv)"
echo "== prior-leg disagreement: nonimpl_marked=0 / header_nonimpl=0 / self_named low =="
awk '{split($2,a,"="); if(a[2]+0==0) printf "nonimpl0:%s ", $1}' .scratch/wisp/probes/pool-validity/5g/logs/rolecount.tsv; echo
awk '{split($2,a,"="); if(a[2]+0==0) printf "header0:%s ", $1}' .scratch/wisp/probes/pool-validity/5g/logs/header.tsv; echo
echo "== 104 / 110 unchecked box lines =="
grep -nE '^[[:space:]]*- \[ \]' .scratch/wisp/issues/104-*.md | cut -c1-120
grep -nE '^[[:space:]]*- \[ \]' .scratch/wisp/issues/110-*.md | cut -c1-120
echo "== ledger lines for 76 =="
grep -nE '(票 ?76([^0-9]|$)|[Tt]icket ?76([^0-9]|$)|# ?76([^0-9]|$))' docs/reports/pending-and-issues.md | cut -c1-120 | head -8
