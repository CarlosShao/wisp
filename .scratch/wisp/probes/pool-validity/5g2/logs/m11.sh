#!/bin/bash
# M11: (1) per-ticket ACC-file verdict aggregation, (2) box-less ticket census over all 94,
# (3) ledger claims for the 52 roster numbers (counts only into context).
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
# 1) verdict content of the acc-shaped files, per ticket
awk -F'|' '{print $1"|"$2}' "$L/m8.tsv" > "$L/m11-role.txt"
awk -F'|' '{v=""; for(i=2;i<=NF;i++) if($i ~ /^[0-9 ]+[A-Za-z]/) v=v $i; print $1"|ln="substr($2,4)"|verdict_raw="v}' "$L/m4.tsv" > /dev/null
awk -F'|' '
NR==FNR{ ln[$1]=substr($2,4); vr[$1]=substr($8,1,44); next }
' "$L/m4.tsv" "$L/m4.tsv"
: > "$L/m11-verdict.tsv"
while read -r N; do
  [ -z "$N" ] && continue
  ACCF=$(awk -F'|' -v n="$N" '$1==n{print $2}' "$L/m3-list.txt" | while read -r F; do
      R=$(awk -F'|' -v f="$F" '$1==f{print $2}' "$L/m11-role.txt")
      case "$R" in ACCHEAD|VLEG) echo "$F";; esac
    done)
  NA=$(printf '%s\n' "$ACCF" | grep -c .)
  VTOT=0; LMAX=0
  for F in $ACCF; do
    V=$(grep -ohE 'PASS WITH CONDITIONS|PASS|FAIL|不通过|退回|拒绝' "$F" | wc -l)
    VTOT=$((VTOT+V))
    LN=$(wc -l < "$F"); [ "$LN" -gt "$LMAX" ] && LMAX=$LN
  done
  echo "$N|accfiles=$NA|accverdicttokens=$VTOT|accmaxlines=$LMAX" >> "$L/m11-verdict.tsv"
done < "$L/roster52.txt"
echo "M11V=$(wc -l < "$L/m11-verdict.tsv")"
echo "ZEROVERDICT=$(awk -F'|' '$2 ~ /accverdicttokens=0/{printf "%s ",$1}' $L/m11-verdict.tsv)"
echo "ACCFILES0=$(awk -F'|' '$2=="accfiles=0"{printf "%s ",$1}' $L/m11-verdict.tsv)"
# 2) box-less census across all 94 done tickets
: > "$L/m11-boxes.txt"
for F in .scratch/wisp/issues/*-done.md; do
  N=$(basename "$F" | sed -E 's/^([0-9]+)-.*/\1/')
  A=$(grep -c '^- \[ \]' "$F"); B=$(grep -c '^[[:space:]]*- \[ \]' "$F")
  T=$(grep -c '^[[:space:]]*- \[' "$F")
  ANY=$(grep -c '\[ \]' "$F")
  echo "$N|unind=$A|anyind=$B|boxes=$T|bracketempty=$ANY" >> "$L/m11-boxes.txt"
done
echo "BOXES0_all94=$(awk -F'|' '$3=="boxes=0"{printf "%s ",$1}' $L/m11-boxes.txt)"
echo "BRACKET0_all94=$(awk -F'|' '$4=="bracketempty=0"{printf "%s ",$1}' $L/m11-boxes.txt)"
echo "UNIND_NEQ_ANYIND=$(awk -F'|' '{split($2,a,"=");split($3,b,"="); if(a[2]!=b[2]) printf "%s(%s/%s) ",$1,a[2],b[2]}' $L/m11-boxes.txt)"
# 3) ledger claims
P=docs/reports/pending-and-issues.md
echo "ledger_exists=$(test -f $P && echo yes) lines=$(wc -l < $P)"
: > "$L/m11-ledger.tsv"
while read -r N; do
  [ -z "$N" ] && continue
  C=$(grep -cE "(票 ?$N([^0-9]|$)|[Tt]icket ?$N([^0-9]|$)|# ?$N([^0-9]|$))" "$P")
  CD=$(grep -E "(票 ?$N([^0-9]|$)|[Tt]icket ?$N([^0-9]|$)|# ?$N([^0-9]|$))" "$P" | grep -ciE '已结案|已验收|验收|结案|PASS|归口|证据|凭据')
  echo "$N|ledgerlines=$C|of_which_close_or_accept_words=$CD" >> "$L/m11-ledger.tsv"
done < "$L/roster52.txt"
echo "LEDGER0=$(awk -F'|' '$2=="ledgerlines=0"{printf "%s ",$1}' $L/m11-ledger.tsv)"
echo "LEDGER_NOCLOSEWORDS=$(awk -F'|' '$2!="ledgerlines=0" && $3 ~ /=0$/{printf "%s ",$1}' $L/m11-ledger.tsv)"
