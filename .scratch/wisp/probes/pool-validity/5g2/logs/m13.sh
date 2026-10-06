#!/bin/bash
# M13: where do the acceptance-shaped artifacts live (docs/evidence vs .scratch)?
# Also: ticket 184 deep check + numeric-safe ledger claim parse.
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
: > "$L/m13-loc.tsv"
while read -r N; do
  [ -z "$N" ] && continue
  ACCF=$(awk -F'|' -v n="$N" '$1==n{print $2}' "$L/m3-list.txt" | while read -r F; do
      R=$(awk -F'|' -v f="$F" '$1==f{print $2}' "$L/m12-role.txt")
      case "$R" in ACCHEAD|VLEG) echo "$F";; esac
    done)
  TD=$(printf '%s\n' "$ACCF" | grep -c '^docs/evidence/s1/')
  TP=$(printf '%s\n' "$ACCF" | grep -c '^\.scratch/wisp/probes/')
  TO=$(( $(printf '%s\n' "$ACCF" | grep -c .) - TD - TP ))
  # does the strongest acc artifact name itself carry an acceptor leg id?
  echo "$N|acc_docs_evidence_s1=$TD|acc_probes=$TP|acc_other=$TO" >> "$L/m13-loc.tsv"
done < "$L/roster52.txt"
echo "M13=$(wc -l < "$L/m13-loc.tsv")"
echo "ACC_ONLY_IN_PROBES=$(awk -F'|' '{split($2,a,"=");split($3,b,"="); if(a[2]+0==0 && b[2]+0>0) printf "%s ",$1}' $L/m13-loc.tsv)"
echo "ACC_NOWHERE=$(awk -F'|' '{split($2,a,"=");split($3,b,"=");split($4,c,"="); if(a[2]+0==0&&b[2]+0==0&&c[2]+0==0) printf "%s ",$1}' $L/m13-loc.tsv)"
echo "== 184: every file in repo whose name carries 184 =="
find . -path ./.git -prune -o -name '*184*' -print | cut -c1-95
echo "== 184: mentions inside docs/evidence/s1 filenames of other numbers (body) =="
awk -F'|' '$1==184{print $2}' "$L/m7-offindex.txt" | grep -c .
echo "== 184 ticket: header + progress log verdict lines =="
F=$(ls .scratch/wisp/issues/184-*.md); echo "FILE=$F lines=$(wc -l < $F)"
grep -nE '非实现者|对抗验收|验收|裁决|PASS|FAIL|通过|不通过|结案' "$F" | tail -12 | cut -c1-150
echo "== ledger parse (numeric-safe) =="
awk -F'|' '{split($2,a,"=");split($3,b,"="); if(a[2]+0==0) printf "%s ",$1}' "$L/m11-ledger.tsv"; echo "<<LEDGER_ZERO_MENTION"
awk -F'|' '{split($2,a,"=");split($3,b,"="); if(a[2]+0>0 && b[2]+0==0) printf "%s(%s) ",$1,a[2]}' "$L/m11-ledger.tsv"; echo "<<LEDGER_MENTION_BUT_NO_CLOSEWORDS"
