#!/bin/sh
set -u
C=/d/tmp/wisp/20261019/343v1/logs/census
D=/d/tmp/wisp/20261019/343v1/logs/census
{
echo "############ AC#0 same-family census: blacklist.go, three shapes, both sides ############"
echo "== START side (blob c1799cd2) RED-nowhere, per shape =="
awk '$3=="RED-nowhere"{print}' $D/v3-blacklist-start.txt | sort -u
echo
echo "== DELIVERY side (blob 62f2dddb) RED-nowhere, per shape =="
awk '$3=="RED-nowhere"{print}' $C/v2-blacklist-delivery.txt | sort -u
echo
echo "== DELIVERY side RED-elsewhere (shape SB) = the FP class for a package-scoped gate =="
awk '$3=="RED-elsewhere"{print}' $C/v2-blacklist-delivery.txt | sort -u
echo
echo "== union sets =="
awk '$3=="RED-nowhere"{print $1,$2,$5":"$6}' $D/v3-blacklist-start.txt | sort -u > $C/set-start-red.txt
awk '$3=="RED-nowhere"{print $1,$2,$5":"$6}' $C/v2-blacklist-delivery.txt | sort -u > $C/set-del-red.txt
wc -l $C/set-start-red.txt $C/set-del-red.txt
echo "-- symdiff start vs delivery (only aEquiv lines expected) --"
diff $C/set-start-red.txt $C/set-del-red.txt || true
echo "-- SA-shape RED-nowhere names, delivery --"
awk '$1=="SA" && $3=="RED-nowhere"{print}' $C/v2-blacklist-delivery.txt | sort -u || echo none
echo "-- SW-shape RED-nowhere names, delivery --"
awk '$1=="SW" && $3=="RED-nowhere"{print}' $C/v2-blacklist-delivery.txt | sort -u || echo none
echo "-- SB-shape RED-nowhere names, delivery (the fires() family) --"
awk '$1=="SB" && $3=="RED-nowhere"{print $2, $5":"$6}' $C/v2-blacklist-delivery.txt | sort -u
} > $C/ac0-summary.txt 2>&1
echo DONE
cat $C/ac0-summary.txt
