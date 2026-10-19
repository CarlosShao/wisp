#!/bin/sh
set -u
C=/d/tmp/wisp/20261019/343v1/logs/census
O=$C/ac3-strict.txt
{
echo "############ AC#3 STRICT call-shape (foo( , no space) — my own re-measurement ############"
echo "== production-only (non _test.go) row/name counts, START side =="
for sh in SA SB SBL SW; do
  rows=$(awk -v s="$sh" '$1==s' $C/v3-full-start.txt | grep -v _test.go | wc -l)
  redr=$(awk -v s="$sh" '$1==s && ($3=="RED-nowhere"||$3=="RED-elsewhere")' $C/v3-full-start.txt | grep -v _test.go | wc -l)
  redn=$(awk -v s="$sh" '$1==s && ($3=="RED-nowhere"||$3=="RED-elsewhere")' $C/v3-full-start.txt | grep -v _test.go | awk '{print $2}' | sort -u | wc -l)
  inp=$(awk -v s="$sh" '$1==s && $3=="IN-PKG"' $C/v3-full-start.txt | grep -v _test.go | awk '{print $2}' | sort -u | wc -l)
  printf 'shape=%s prodrows=%s prodREDrows=%s prodREDunique=%s prodINPKGunique=%s\n' "$sh" "$rows" "$redr" "$redn" "$inp"
done
echo
echo "== the strict-SB production RED names, START side, one per line with file:line and where-field =="
awk '$1=="SB" && ($3=="RED-nowhere"||$3=="RED-elsewhere")' $C/v3-full-start.txt | grep -v _test.go | sort -u -k2,2 -k6,6
echo
echo "== count =="
awk '$1=="SB" && ($3=="RED-nowhere"||$3=="RED-elsewhere")' $C/v3-full-start.txt | grep -v _test.go | awk '{print $2}' | sort -u | wc -l
echo "== and the leg's own claim set for comparison (28 claims / 11 dangling): what I get is above =="
echo
echo "== DEFERRED(...) marker rows (the FP class the leg named) =="
grep -n 'DEFERRED' internal/risk/*.go 2>/dev/null | grep -v '_test.go' | head -12
echo "(cwd note: this last block runs in repo root only if cd'ed; harmless if empty)"
echo AC3STRICT_DONE
} > $O 2>&1
cat $O
