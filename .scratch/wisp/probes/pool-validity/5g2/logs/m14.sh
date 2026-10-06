#!/bin/bash
# M14: clean box census over all 94 done tickets + location census for the 52.
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
: > "$L/m14-boxes.txt"
for F in .scratch/wisp/issues/*-done.md; do
  N=$(basename "$F" | sed -E 's/^([0-9]+)-.*/\1/')
  T=$(grep -c '^[[:space:]]*- \[' "$F")
  U=$(grep -c '^- \[ \]' "$F")
  A=$(grep -c '^[[:space:]]*- \[ \]' "$F")
  D=$(grep -c '^[[:space:]]*- \[[xX]\]' "$F")
  ANY=$(grep -cE '\[[ xX]\]' "$F")
  echo "$N|boxes=$T|unind=$U|anyind=$A|done=$D|anybracket=$ANY" >> "$L/m14-boxes.txt"
done
echo "rows=$(wc -l < "$L/m14-boxes.txt")"
echo "ZEROBOX_ALL94: $(awk -F'|' '{split($2,a,"=");if(a[2]+0==0) printf "%s ",$1}' $L/m14-boxes.txt)"
echo "ZEROBRACKET_ALL94: $(awk -F'|' '{split($5,e,"=");if(e[2]+0==0) printf "%s ",$1}' $L/m14-boxes.txt)"
echo "ANYIND_NE_UNIND_ALL94: $(awk -F'|' '{split($3,b,"=");split($4,c,"=");if(b[2]+0!=c[2]+0) printf "%s(%s/%s) ",$1,b[2],c[2]}' $L/m14-boxes.txt)"
echo "UNCHECKED_GT0_ALL94: $(awk -F'|' '{split($4,c,"=");if(c[2]+0>0) printf "%s(%s) ",$1,c[2]}' $L/m14-boxes.txt)"
echo "UNCHECKED_GT0_ROSTER52: $(while read -r N; do grep -m1 "^$N|" $L/m14-boxes.txt | awk -F'|' -v n="$N" '{split($4,c,"=");if(c[2]+0>0) printf "%s(%s) ",n,c[2]}'; done < $L/roster52.txt)"
echo "ZEROBOX_ROSTER52: $(while read -r N; do grep -m1 "^$N|" $L/m14-boxes.txt | awk -F'|' -v n="$N" '{split($2,a,"=");if(a[2]+0==0) printf "%s ",n}'; done < $L/roster52.txt)"
echo "== location census =="
bash "$L/m13.sh" 2>&1 | grep -E '^M13|^ACC_ONLY_IN_PROBES|^ACC_NOWHERE'
