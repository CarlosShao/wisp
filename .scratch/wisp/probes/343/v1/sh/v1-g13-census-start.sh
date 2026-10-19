#!/bin/sh
set -u
C=/d/tmp/wisp/20261019/343v1/logs/census
BIN=/d/tmp/wisp/20261019/343v1/census/census343v1.exe
# same gauge, run inside my START tree (blacklist.go blob c1799cd2 = the un-fixed comment)
{
echo "=== run: blacklist.go only, START side (treeBase) ==="
( cd /d/tmp/wisp/20261019/343v1/treeBase && $BIN -pkg internal/risk -roots internal,cmd,tools -file blacklist.go )
} > $C/v3-blacklist-start.txt 2> $C/v3-blacklist-start.err
echo "RC3=$?"
{
echo "=== run: whole package, START side (treeBase) ==="
( cd /d/tmp/wisp/20261019/343v1/treeBase && $BIN -pkg internal/risk -roots internal,cmd,tools )
} > $C/v3-full-start.txt 2> $C/v3-full-start.err
echo "RC4=$?"
echo "== start-side counters (blacklist.go) =="
grep '^shape=' $C/v3-blacklist-start.txt
echo "== start-side counters (whole pkg) =="
grep '^shape=' $C/v3-full-start.txt
echo "== the aEquiv line on start side, all shapes =="
grep -n 'aEquiv' $C/v3-blacklist-start.txt || echo "no aEquiv line"
echo "== header lines =="
head -1 $C/v3-blacklist-start.txt
head -1 $C/v3-blacklist-delivery.txt
wc -l $C/v3-blacklist-start.txt $C/v3-blacklist-start.err $C/v3-full-start.txt $C/v3-full-start.err
echo DONE
