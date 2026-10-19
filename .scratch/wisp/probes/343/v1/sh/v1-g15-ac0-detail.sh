#!/bin/sh
set -u
C=/d/tmp/wisp/20261019/343v1/logs/census
{
echo "== whole-package DELIVERY: SA RED-elsewhere (4) =="
awk '$1=="SA" && $3=="RED-elsewhere"' $C/v2-full-delivery.txt
echo
echo "== whole-package DELIVERY: SA RED-nowhere (13) =="
awk '$1=="SA" && $3=="RED-nowhere"' $C/v2-full-delivery.txt
echo
echo "== whole-package DELIVERY: SW (bounded wordlist a[A-Z]..) all lines =="
awk '$1=="SW"' $C/v2-full-delivery.txt
echo
echo "== whole-package START: SW RED-nowhere =="
awk '$1=="SW" && $3=="RED-nowhere"' $C/v3-full-start.txt
echo
echo "== unique SA RED names (delivery, whole pkg), with file =="
awk '$1=="SA" && ($3=="RED-elsewhere"||$3=="RED-nowhere"){print $2, "|", $6, $7}' $C/v2-full-delivery.txt | sort -u
echo
echo "== BEXCL (builtin/keyword call-shape noise) whole pkg delivery =="
awk '$1=="BEXCL"{print $2}' $C/v2-full-delivery.txt | sort | uniq -c
} > $C/ac0-detail.txt 2>&1
cat $C/ac0-detail.txt
