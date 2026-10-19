#!/bin/sh
set -u
C=/d/tmp/wisp/20261019/343v1/logs/census
O=$C/ac3-table2.txt
{
echo "### AC#3 table (343-v1's own gauge; production files = non _test.go; START side blob c1799cd2) ###"
for side in v3-full-start v2-full-delivery; do
  echo "--- $side ---"
  for sh in SA SB SW; do
    lines=$(awk -v s="$sh" '$1==s' $C/$side.txt | grep -v _test.go | wc -l)
    redl=$(awk -v s="$sh" '$1==s && ($3=="RED-nowhere"||$3=="RED-elsewhere")' $C/$side.txt | grep -v _test.go | wc -l)
    redn=$(awk -v s="$sh" '$1==s && ($3=="RED-nowhere"||$3=="RED-elsewhere")' $C/$side.txt | grep -v _test.go | awk '{print $2}' | sort -u | wc -l)
    nowh=$(awk -v s="$sh" '$1==s && $3=="RED-nowhere"' $C/$side.txt | grep -v _test.go | awk '{print $2}' | sort -u | wc -l)
    els=$(awk -v s="$sh" '$1==s && $3=="RED-elsewhere"' $C/$side.txt | grep -v _test.go | awk '{print $2}' | sort -u | wc -l)
    inpkg=$(awk -v s="$sh" '$1==s && $3=="IN-PKG"' $C/$side.txt | grep -v _test.go | awk '{print $2}' | sort -u | wc -l)
    printf 'shape=%s comment-rows=%s RED_rows=%s RED_unique=%s (nowhere=%s elsewhere=%s) INPKG_unique=%s\n' "$sh" "$lines" "$redl" "$redn" "$nowh" "$els" "$inpkg"
  done
done
echo
echo "### unique SA RED names, START, production ###"
awk '$1=="SA" && ($3=="RED-nowhere"||$3=="RED-elsewhere")' $C/v3-full-start.txt | grep -v _test.go | awk '{print $2}' | sort -u
echo
echo "### SW bounded wordlist: unique names seen in production comments, START ###"
awk '$1=="SW"' $C/v3-full-start.txt | grep -v _test.go | awk '{print $2}' | sort -u | tr '\n' ' '; echo
echo "### SW RED rows (either side, production) ###"
awk '$1=="SW" && ($3=="RED-nowhere"||$3=="RED-elsewhere")' $C/v3-full-start.txt | grep -v _test.go || echo NONE
echo "### SW on delivery side production, RED rows ###"
awk '$1=="SW" && ($3=="RED-nowhere"||$3=="RED-elsewhere")' $C/v2-full-delivery.txt | grep -v _test.go || echo NONE
echo
echo "### the four family members: which shape sees each (START, production rows) ###"
for nm in aEquiv fires taskScopeID disposalScope; do
  printf '%-14s' "$nm"
  for sh in SA SB SW; do
    h=$(awk -v s="$sh" -v n="$nm" '$1==s && $2==n && ($3=="RED-nowhere"||$3=="RED-elsewhere")' $C/v3-full-start.txt | grep -v _test.go | wc -l)
    printf ' %s=%s' "$sh" "$h"
  done; echo
done
echo
echo "### SB RED-nowhere production unique names (the FP class) ###"
awk '$1=="SB" && $3=="RED-nowhere"' $C/v3-full-start.txt | grep -v _test.go | awk '{print $2}' | sort -u > $C/sb-names.txt
echo "count=$(wc -l < $C/sb-names.txt)"
echo "all-lowercase-alphabetic (prose words) = $(grep -cE '^[a-z]+$' $C/sb-names.txt || true)"
echo "- the rest:"; grep -vE '^[a-z]+$' $C/sb-names.txt | tr '\n' ' '; echo
echo "- CloseScope present? $(grep -c 'CloseScope' $C/sb-names.txt || true)"
echo "- DEFERRED/Nutstore/OpenScope present? $(grep -cE 'DEFERRED|Nutstore|OpenScope' $C/sb-names.txt || true)"
echo "- full list:"; tr '\n' ' ' < $C/sb-names.txt; echo
echo
echo "### BEXCL (call-shape hits that are builtin/keyword = pure noise) production ###"
awk '$1=="BEXCL"' $C/v3-full-start.txt | grep -v _test.go | awk '{print $2}' | sort | uniq -c
echo AC3TABLE2_DONE
} > $O 2>&1
cat $O
