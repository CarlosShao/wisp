#!/bin/sh
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 90
export LC_ALL=C.UTF-8
C=/d/tmp/wisp/20261019/343v1/logs/census
echo "=== bOverrides in my go/ast census (blind-spot cross-check) ==="
grep -n 'bOverrides' $C/v3-blacklist-start.txt | head -4
echo "=== strict SB rows on blacklist.go START (expect fires + CloseScope-like only) ==="
awk '$1=="SB"' $C/v3-blacklist-start.txt
echo "=== strict SB rows on blacklist.go DELIVERY ==="
awk '$1=="SB"' $C/v2-blacklist-delivery.txt
echo "=== how many aEq(/aUnder( call sites in package code (comment says 'every A-anchor comparison below') ==="
printf 'blacklist.go aEq|aUnder call sites='; grep -c 'aEq(\|aUnder(' internal/risk/blacklist.go || true
printf '  of which in comments='; grep 'aEq(\|aUnder(' internal/risk/blacklist.go | grep -c '^\s*//' || true
echo "=== Gate() signature (bOverrides param, cited by census as IN-PKG param) ==="
grep -n 'func Gate(' internal/risk/blacklist.go | head -3
echo "=== anchorMatches call sites (the new sentence says it does NOT report miss as hit) ==="
printf 'anchorMatches( call sites='; grep -c 'anchorMatches(' internal/risk/*.go || true
echo "=== uncertainAnchorMiss call sites ==="
grep -n 'uncertainAnchorMiss(' internal/risk/*.go | grep -v _test | head -6
echo "=== aRelates body (the two-gauge structure the new comment asserts) ==="
sed -n '277,292p' internal/risk/blacklist.go
echo DONE
