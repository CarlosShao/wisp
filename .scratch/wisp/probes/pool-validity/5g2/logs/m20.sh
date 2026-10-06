#!/bin/bash
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
echo "== ruler-pair check on 3 tickets =="
for N in 07 92 115; do F=$(ls .scratch/wisp/issues/$N-*.md|head -1); echo "$N unind=$(grep -c '^- \[ \]' $F) anyind=$(grep -c '^[[:space:]]*- \[ \]' $F) boxes=$(grep -c '^[[:space:]]*- \[' $F)"; done
echo "== all-94 box census, both rulers =="
: > $L/m20-boxes94.txt
for F in .scratch/wisp/issues/*-done.md; do N=$(basename $F|sed -E 's/^([0-9]+)-.*/\1/')
  U=$(grep -c '^- \[ \]' "$F"); A=$(grep -c '^[[:space:]]*- \[ \]' "$F"); T=$(grep -c '^[[:space:]]*- \[' "$F"); K=$(grep -cE '^[[:space:]]*[-*] \[[xX]\]' "$F")
  echo "$N|unind=$U|anyind=$A|boxes=$T|checked=$K" >> $L/m20-boxes94.txt; done
echo "diff_rulers: $(awk -F'|' '{split($2,a,"=");split($3,b,"="); if(a[2]!=b[2]) printf "%s(%s/%s) ",$1,a[2],b[2]}' $L/m20-boxes94.txt)"
echo "boxes0_94: $(awk -F'|' '{split($4,c,"="); if(c[2]==0) printf "%s ",$1}' $L/m20-boxes94.txt)"
echo "unchecked_52: $(while read -r N; do grep -m1 "^$(echo $N|sed -E 's/^0*//')|" $L/m20-boxes94.txt | awk -F'|' -v n=$N '{split($3,b,"="); if(b[2]+0>0) printf "%s(%s) ",n,b[2]}'; done < $L/roster52.txt)"
echo "== thin acceptance artifacts (<60 lines) among name-index ACC files =="
while read -r N; do K=$(echo $N|sed -E 's/^0*//')
  awk -F'|' -v n="$K" '$1==n{print $2}' $L/m3-list.txt | while read -r F; do
    R=$(grep -m1 -F "$F|" $L/m12-role.txt | cut -d'|' -f2)
    case "$R" in ACCHEAD|VLEG) LN=$(wc -l < "$F"); [ "$LN" -lt 60 ] && echo "$K $F lines=$LN";; esac
  done
done < $L/roster52.txt > $L/m20-thin.txt; cat $L/m20-thin.txt
echo "== 157 acc-file verdict vocabulary =="
for F in docs/evidence/s1/157-*.md; do echo "$F lines=$(wc -l < $F) PASS=$(grep -o 'PASS' $F|wc -l) verdictwords=$(grep -cE '成立|不成立|恒真|判|结论|红|绿' $F)"; done
echo "== 76 acc file head =="; head -8 docs/evidence/s1/76-adversarial-acceptance.md | cut -c1-100; echo "lines=$(wc -l < docs/evidence/s1/76-adversarial-acceptance.md)"
