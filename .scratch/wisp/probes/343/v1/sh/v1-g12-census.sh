#!/bin/sh
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 90
C=/d/tmp/wisp/20261019/343v1/logs/census
mkdir -p $C
BIN=/d/tmp/wisp/20261019/343v1/census/census343v1.exe
{
echo "=== run 1: whole internal/risk package (delivery side) ==="
$BIN -pkg internal/risk -roots internal,cmd,tools
} > $C/v2-full-delivery.txt 2> $C/v2-full-delivery.err
echo "RC1=$?"
{
echo "=== run 2: blacklist.go only (the ticket's target file, delivery side) ==="
$BIN -pkg internal/risk -roots internal,cmd,tools -file blacklist.go
} > $C/v2-blacklist-delivery.txt 2> $C/v2-blacklist-delivery.err
echo "RC2=$?"
wc -l $C/v2-full-delivery.txt $C/v2-full-delivery.err $C/v2-blacklist-delivery.txt $C/v2-blacklist-delivery.err
echo "=== counters delivery / blacklist.go ==="
grep '^shape=' $C/v2-blacklist-delivery.txt
echo "=== counters delivery / whole package ==="
grep '^shape=' $C/v2-full-delivery.txt
echo DONE
