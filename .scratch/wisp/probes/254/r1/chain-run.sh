#!/usr/bin/env bash
# .scratch/wisp/probes/254/r1/chain-run.sh - drives the REAL CI chain
# (scripts/winsec-tests.sh -> scripts/portable-tests.sh) on the carrier's fake go,
# so ticket 254 can read what the route that ci.yml:439 takes actually reports.
#
# This is an evidence probe, not the standing carrier: the durable cases live in
# scripts/portable-tests-selftest.sh. Temp files are created, never deleted.
#
# Usage: bash chain-run.sh <clean|split|widened> <label>
set -u -o pipefail

probe=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
root=$(CDPATH= cd -- "$probe/../../../../.." && pwd) # r1/254/probes/wisp/.scratch -> repo root
[ -f "$root/scripts/winsec-tests.sh" ] || { echo "chain-run.sh: bad root $root" >&2; exit 9; }

mode=${1:?mode}
label=${2:-run}
out="$probe/logs"
mkdir -p "$out/bin"
cp "$root/scripts/testdata/portable-tests/go" "$out/bin/go"
chmod +x "$out/bin/go" 2>/dev/null || true

# The three pin files are read out of the script under test, not written here.
awk -v q="'" 'index($0,"winsec_pin="q)==1{f=1;next} f{if($0==q){f=0;next} print}' \
    "$root/scripts/portable-tests.sh" >"$out/winsec-pin.txt"
awk -v q="'" 'index($0,"core_pin="q)==1{f=1;next} f{if($0==q){f=0;next} print}' \
    "$root/scripts/portable-tests.sh" >"$out/core-pin.txt"
IFS= read -r wpath <"$out/winsec-pin.txt"
{ cat "$out/winsec-pin.txt"; printf '%s/acl\n' "$wpath"; } >"$out/winsec-pin-split.txt"

pin="$out/core-pin.txt"
args=()
case $mode in
clean)   pin="$out/winsec-pin.txt" ;;
split)   pin="$out/winsec-pin-split.txt" ;;
widened) pin="$out/core-pin.txt"; args=("./internal/winsec/" "./internal/agent/") ;;
*) echo "chain-run.sh: unknown mode $mode" >&2; exit 9 ;;
esac

log="$out/$label-$mode.txt"
rc=0
( cd "$root" && PATH="$out/bin:$PATH" FAKEGO_PIN_FILE="$pin" FAKEGO_SCOPE_FILTER=1 \
    FAKEGO_GOOS=windows bash scripts/winsec-tests.sh ${args[@]+"${args[@]}"} ) >"$log" 2>&1 || rc=$?

audit=$(grep -c "IS the winsec tier" "$log" || true)
guardc=$(grep -c "GUARD C - scope mode=winsec" "$log" || true)
scopeglob=$(grep -cE '^portable-tests\.sh: platform=windows scope=\[\./internal/winsec/\.\.\.\]$' "$log" || true)
okline=$(grep -cE '^(ok|FAIL)[[:space:]]+github\.com/CarlosShao/wisp/internal/winsec[[:space:]]' "$log" || true)
four=$(grep -c 'four numbers' "$log" || true)
printf 'chain %s mode=%s rc=%s auditline=%s guardc=%s scopeglob=%s resultline=%s fournumbers=%s\n' \
    "$label" "$mode" "$rc" "$audit" "$guardc" "$scopeglob" "$okline" "$four"
printf 'log=%s\n' "$log"
