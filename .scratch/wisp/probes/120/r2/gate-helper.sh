#!/bin/sh
# 120-r2 gate helper: the ticket-113 acceptance roster, named one by one, in the
# container. Written by the 120-r2 leg; logs land in .scratch/wisp/probes/120/r2/.
set -e
L=.scratch/wisp/probes/120/r2/120-postfix-gate-count2-v.log
echo "unique-run-names=$(sed -n 's/^=== RUN[[:space:]]*//p' "$L" | sort -u | wc -l)"
echo "count-per-120-name:"
sed -n 's/^=== RUN[[:space:]]*\(TestPOSIX120[A-Za-z0-9_]*\).*/\1/p' "$L" | sort | uniq -c
echo "count-per-ac-posix-name:"
sed -n 's/^=== RUN[[:space:]]*\(TestAC[0-9]POSIX[A-Za-z0-9_]*\).*/\1/p' "$L" | sort | uniq -c
