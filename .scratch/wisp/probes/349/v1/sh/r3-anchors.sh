#!/bin/sh
cd '/d/work/workspace/projects plans/Wisp' || exit 9
LEGFIRST=bf7f13bce55cd91ced1db96f02dbac2ea94a405a
BASE=$(git rev-parse $LEGFIRST~1)
LEG=ed8db1ea27eb50d3057dd0cd7d5c955a13e978fc
echo "BASEANCHOR=$BASE"
git log -1 --format='%h %ad %s' --date=iso $BASE | cut -c1-140
echo "LEGDELIVER=$LEG"
echo "HEADNOW=$(git rev-parse HEAD)"
for f in internal/risk/blacklist.go internal/risk/provenance.go internal/risk/syncdirs.go; do
  b=$(git rev-parse $BASE:$f); l=$(git rev-parse $LEG:$f); h=$(git rev-parse HEAD:$f)
  printf '%s\n  BASE=%s\n  LEG  =%s\n  HEAD =%s\n  BASEeqLEG=%s  LEGeqHEAD=%s\n' "$f" "$b" "$l" "$h" \
    "$([ "$b" != "$l" ] && echo CHANGED || echo SAME)" "$([ "$l" = "$h" ] && echo YES || echo NO)"
done
echo '--- last commit touching each file before leg first (blob-identity walk back) ---'
for f in internal/risk/blacklist.go internal/risk/provenance.go internal/risk/syncdirs.go; do
  tb=$(git rev-parse $BASE:$f)
  c=$(git rev-list -1 $BASE -- $f)
  while [ -n "$c" ] && [ "$(git rev-parse $c:$f 2>/dev/null)" != "$tb" ]; do c=$(git rev-list -1 $c~1 -- $f); done
  printf '%s lastcommit=%s %s\n' "$f" "$(git log -1 --format='%h %ad' --date=format:'%m-%d %H:%M' $c)" "$(git log -1 --format='%s' $c | cut -c1-60)"
done
echo '--- commits between BASE and LEG touching internal/risk ---'
git log --format='%h %s' $BASE..$LEG -- internal/risk | cut -c1-90
echo '--- commits between LEG and HEAD touching internal/risk (should be 0) ---'
git log --format='%h %s' $LEG..HEAD -- internal/risk | cut -c1-90
echo RC=$?
