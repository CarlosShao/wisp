#!/bin/sh
# Ticket 335 leg 335-r1, AC#1's paired de-fang (plus one bonus shot). ONE mutation per
# run, applied to the PHYSICAL BYTES of the out-of-repo export tree T3, restored in the
# same step through a trap, restore proven by hash. Never run this inside the repository.
#   usage: sh mut335r1.sh <A|B|C>
#     A = isDisplayWashCf's body token replaced by the 15-rune hand roster ticket 333
#         removed (content anchor: the single "return unicode.Is(unicode.Cf, r)" token)
#     B = a per-rune exception for U+2061 inserted right after the signature line of
#         isDisplayWashCf (a different content anchor AND a different line than A)
#     C = bonus, the other direction a hand list can be wrong: a per-rune exception that
#         washes U+00A0, which the category does not name
set -e
TREE=/d/V335R1/T3
WORK=/d/V335R1/mut
REPO="/d/work/workspace/projects plans/Wisp"
LOGS="$REPO/.scratch/wisp/probes/335/r1/logs"
TARGET=internal/panel/display_scrub_216.go
ROSTER='TestDisplayWash|TestEveryNamedCfRuneIsEscapedAndNeverDeleted|TestTheRosterIsClosedAgainstTheRunesTheRulingLeftOut|TestScrubForDisplay'
label="$1"
cd "$TREE" || exit 99
mkdir -p "$WORK" "$LOGS"

cp "$TARGET" "$WORK/pristine-$label.go"
pre=$(sha256sum "$TARGET" | cut -d' ' -f1)
restore() { cp "$WORK/pristine-$label.go" "$TARGET"; }
trap restore EXIT INT TERM

echo "SHOT=$label"
echo "PRE_HASH=$pre"
if [ "$label" = "A" ]; then
	awk 'BEGIN{while ((getline l < "/d/V335R1/mut/m-a-body.txt") > 0) body = body l "\n"}
	     /^[[:space:]]*return unicode\.Is\(unicode\.Cf, r\)[[:space:]]*$/ { printf "%s", body; next }
	     { print }' "$WORK/pristine-$label.go" > "$WORK/mutated-$label.go"
elif [ "$label" = "B" ]; then
	awk '/^func isDisplayWashCf\(r rune\) bool \{$/ {
			print; print "\tif r == 0x2061 { return false }"; next }
		{ print }' "$WORK/pristine-$label.go" > "$WORK/mutated-$label.go"
else
	awk '/^func isDisplayWashCf\(r rune\) bool \{$/ {
			print; print "\tif r == 0x00A0 { return true }"; next }
		{ print }' "$WORK/pristine-$label.go" > "$WORK/mutated-$label.go"
fi
cp "$WORK/mutated-$label.go" "$TARGET"
echo "MUT_HASH=$(sha256sum "$TARGET" | cut -d' ' -f1)"
echo "MUT_CHANGED_LINES=$(diff "$WORK/pristine-$label.go" "$TARGET" | grep -c '^[<>]')"
sed -n '/^func isDisplayWashCf/,/^}/p' "$TARGET" > "$WORK/body-$label.txt"
echo "MUTATED_BODY_LINES=$(wc -l < "$WORK/body-$label.txt")"

export TMP="/tmp" TEMP="/tmp"
unset GOTMPDIR
set +e
go test -v ./internal/panel/ -count=1 -timeout 30m -run "$ROSTER" \
	>"$WORK/out-$label.log" 2>&1
rc_test=$?
set -e
echo "TEST_RC=$rc_test"
cp "$WORK/out-$label.log" "$LOGS/mut-$label.log"

restore
trap - EXIT INT TERM
post=$(sha256sum "$TARGET" | cut -d' ' -f1)
echo "RESTORE_HASH=$post"
[ "$pre" = "$post" ] && echo "RESTORE_VERIFIED=YES" || echo "RESTORE_VERIFIED=NO"
echo "RUN=$(grep -c '^=== RUN' "$WORK/out-$label.log")"
echo "PASS_TOP=$(grep -c '^--- PASS' "$WORK/out-$label.log")"
echo "FAIL_TOP=$(grep -c '^--- FAIL' "$WORK/out-$label.log")"
echo "SKIP=$(grep -c '^--- SKIP' "$WORK/out-$label.log")"
echo "RED_NAMES:"
grep '^--- FAIL' "$WORK/out-$label.log" | sed 's/ (.*//; s/^/  /'
echo "RED_LINES_IN_RESHAPED_NAIL:"
grep -Eo 'display_scrub_216_test\.go:[0-9]+' "$WORK/out-$label.log" | sort -u | sed 's/^/  /'
