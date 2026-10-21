#!/bin/sh
# r1-unfixed-probe.sh - ticket 171 AC#3: "does the mutation ring on the UNMODIFIED ruler?"
# This is the anti-decoration reading. A nail that only fires after you add it proves
# nothing about the hole; what has to be shown is that the EXACT mutation (delete one
# must-ring sample) reads GREEN, rc=0, on the ruler as committed at HEAD.
#
# Run from the repo root. Writes only under probes/171/r1/** (this leg's surface).
set -eu
D=.scratch/wisp/probes/171/r1
ANCHOR=HEAD

git show "$ANCHOR:.scratch/wisp/probes/161/r5/attrib.sh" >"$D/attrib-prefix-head.sh"
echo "== unfixed ruler extracted from $ANCHOR ($(wc -l <"$D/attrib-prefix-head.sh" | tr -d ' ') lines) =="

# NO roster assertion anywhere in the unfixed text: the positive control for "this is
# the code before the nail", read off the bytes rather than remembered.
n=$(grep -cE 'ST_RING_EXPECT|ST_SILENT_EXPECT|named-missing' "$D/attrib-prefix-head.sh" || true)
echo "== roster-assertion tokens in the unfixed ruler: $n (must be 0) =="
[ "$n" = 0 ] || { echo "r1-unfixed-probe.sh: the 'unfixed' copy already carries a roster assertion - it is not the unfixed code. no verdict." >&2; exit 2; }

# The same mutation the fixed ruler is judged on: delete the 777 must-ring case.
LN=$(grep -n "st_case 1 UNATTRIBUTABLE '.scratch/wisp/probes/777/second-fake-ticket.go'" "$D/attrib-prefix-head.sh" | cut -d: -f1)
echo "== the 777 must-ring case sits at :$LN of the unfixed ruler =="
[ -n "$LN" ] || { echo "r1-unfixed-probe.sh: no 777 case in the unfixed ruler - no verdict" >&2; exit 2; }
sed "${LN}d" "$D/attrib-prefix-head.sh" >"$D/attrib-prefix-head-noring.sh"
echo "== the mutation, as a diff (1 removed / 0 added expected) =="
echo "   removed=$(diff "$D/attrib-prefix-head.sh" "$D/attrib-prefix-head-noring.sh" | grep -c '^<' || true) added=$(diff "$D/attrib-prefix-head.sh" "$D/attrib-prefix-head-noring.sh" | grep -c '^>' || true)"
diff "$D/attrib-prefix-head.sh" "$D/attrib-prefix-head-noring.sh" | sed 's/^/   /' || true
sh -n "$D/attrib-prefix-head-noring.sh"
echo "== sh -n on the mutant: rc=0 =="
