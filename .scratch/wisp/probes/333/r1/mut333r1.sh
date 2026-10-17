#!/bin/sh
# Ticket 333 leg 333-r1, AC#5: ONE mutation shot on the PHYSICAL BYTES of an
# out-of-repo export tree, then restore in the same step and prove the restore by hash.
# Never run this inside the repository: the tree lives at D:\W333R1MUT.
#   usage: sh mut333r1.sh <M1|M2>
set -e
TREE=/d/W333R1MUT
WORK=/d/wisp-333r1-tmp/mut
REPO="/d/work/workspace/projects plans/Wisp"
LOGS="$REPO/.scratch/wisp/probes/333/r1/logs"
TARGET=internal/panel/display_scrub_216.go
ROSTER='TestDisplayWash|TestEveryNamedCfRuneIsEscapedAndNeverDeleted|TestNeedsDisplayWashNamesEveryWashedCfRune|TestTheWriterBranchWashesCfEvenWhenTheFastPathAlreadyFired|TestTheRosterIsClosedAgainstTheRunesTheRulingLeftOut'
label="$1"
cd "$TREE" || exit 99
mkdir -p "$WORK"

cp "$TARGET" "$WORK/pristine-$label.go"
pre=$(sha256sum "$TARGET" | cut -d' ' -f1)
echo "SHOT=$label"
echo "PRE_HASH=$pre"

if [ "$label" = "M1" ]; then
	# Landing point: the single token inside isDisplayWashCf's body (content anchor,
	# not a line number). Mutation: back to the hand roster ticket 333 removed.
	cp "$WORK/pristine-$label.go" "$WORK/base.go"
	awk 'BEGIN{while ((getline l < "/d/wisp-333r1-tmp/mut/m1-body.txt") > 0) body = body l "\n"}
	     /return unicode\.Is\(unicode\.Cf, r\)/ { printf "%s", body; next }
	     { print }' "$WORK/base.go" > "$WORK/mutated-$label.go"
else
	# Landing point: needsDisplayWash's signature line - a different function and a
	# different assertion line than M1. Mutation: the fast path answers constant false.
	cp "$WORK/pristine-$label.go" "$WORK/base.go"
	awk '/^func needsDisplayWash\(s string\) bool \{$/ { print; print "\treturn false"; next }
	     { print }' "$WORK/base.go" > "$WORK/mutated-$label.go"
fi
cp "$WORK/mutated-$label.go" "$TARGET"
muthash=$(sha256sum "$TARGET" | cut -d' ' -f1)
echo "MUT_HASH=$muthash"
echo "MUT_DIFF_FROM_PRISTINE_LINES=$(diff "$WORK/pristine-$label.go" "$TARGET" | grep -c '^[<>]')"

go test -v ./internal/panel/ -count=1 -timeout 30m -run "$ROSTER" >"$WORK/out-$label.log" 2>&1 || true
cp "$WORK/out-$label.log" "$LOGS/mut-$label.log"

cp "$WORK/pristine-$label.go" "$TARGET"
post=$(sha256sum "$TARGET" | cut -d' ' -f1)
echo "RESTORE_HASH=$post"
if [ "$pre" = "$post" ]; then echo "RESTORE_VERIFIED=YES"; else echo "RESTORE_VERIFIED=NO"; fi
echo "RUN=$(grep -c '^=== RUN' "$WORK/out-$label.log")"
echo "PASS_TOP=$(grep -c '^--- PASS' "$WORK/out-$label.log")"
echo "FAIL_TOP=$(grep -c '^--- FAIL' "$WORK/out-$label.log")"
echo "SKIP=$(grep -c '^--- SKIP' "$WORK/out-$label.log")"
echo "RED_NAMES:"
grep '^--- FAIL' "$WORK/out-$label.log" | sed 's/ (.*//; s/^/  /'
