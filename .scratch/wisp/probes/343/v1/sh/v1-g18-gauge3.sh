#!/bin/sh
# AC#0 third gauge, INDEPENDENT IMPLEMENTATION PATH: pure shell (no go/ast, no perl),
# so it cannot share a bug with my census.exe.
# For every identifier-looking token appearing in a COMMENT of blacklist.go, ask:
#   does the package define that name anywhere (func/method/type/const/var/param/field/local)?
# Definition search = whole internal/risk/*.go with a declaration-shaped regex.
set -u
cd '/d/tmp/wisp/20261019/343v1' || exit 90
SRC=treeBase/internal/risk/blacklist.go
SRC2=treeDel/internal/risk/blacklist.go
PKG='/d/work/workspace/projects plans/Wisp/internal/risk'
OUT=logs/census/gauge3
mkdir -p $OUT

toks () { # comment text -> unique identifier tokens (>=2 chars)
  sed -n 's,^[[:space:]]*//,,p' "$1" | grep -oE '[A-Za-z_][A-Za-z0-9_]{1,}' | sort -u
}

defined_in_pkg () { # $1=name ; declaration-shaped grep over the package
  grep -hE "(^|[[:space:]])(func([[:space:]]+\([^)]*\)[[:space:]]+)?|type|const|var)[[:space:]]+$1\b" \
    "$PKG"/*.go 2>/dev/null | head -1
}

for side in START DELIVERY; do
  if [ "$side" = START ]; then F=$SRC; else F=$SRC2; fi
  {
    echo "=== gauge3 (shell-only), side=$side ==="
    echo "-- tokens named in comments, checked against package declarations --"
    for t in $(toks "$F"); do
      if [ -n "$(defined_in_pkg "$t")" ]; then
        printf 'DEFINED   %s\n' "$t"
      else
        printf 'NOT-DEF   %s\n' "$t"
      fi
    done
  } > $OUT/gauge3-$side.txt 2>&1
done
{
echo "== symdiff of NOT-DEF sets (side START vs side DELIVERY) =="
grep '^NOT-DEF' $OUT/gauge3-START.txt | sort > $OUT/nd-start.txt
grep '^NOT-DEF' $OUT/gauge3-DELIVERY.txt | sort > $OUT/nd-del.txt
wc -l $OUT/nd-start.txt $OUT/nd-del.txt
diff $OUT/nd-start.txt $OUT/nd-del.txt || true
echo "== count of NOT-DEF tokens per side (upper bound: prose words included) =="
printf 'START=%s DELIVERY=%s\n' "$(wc -l < $OUT/nd-start.txt)" "$(wc -l < $OUT/nd-del.txt)"
echo "-- the camel-ish subset only (the ticket's shape), both sides --"
for f in nd-start nd-del; do
  printf '%s camelish=' "$f"
  awk '{print $2}' $OUT/$f.txt | grep -E '[a-z][A-Z]' | tr '\n' ' '
  echo
done
echo "-- the a[A-Z] wordlist subset, both sides --"
for f in nd-start nd-del; do
  printf '%s aword=' "$f"
  awk '{print $2}' $OUT/$f.txt | grep -E '^a[A-Z]' | tr '\n' ' '
  echo
done
echo DONE3
} >> $OUT/gauge3-summary.txt 2>&1
cat $OUT/gauge3-summary.txt
