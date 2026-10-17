#!/bin/sh
# Ticket 335 verdict leg 335-v1. Directed mutations for AC#1's self-reference question
# and the symmetric-difference ruling. ONE scenario per invocation, applied ONLY to the
# PHYSICAL BYTES of the out-of-repo export tree /d/V335V1/T (git archive 23d36a81 + docs
# + manual three), restored in the same invocation through a trap, restore proven by
# hash. Never run this inside the repository.
#   usage: sh mut335v1.sh <A0|B|C|D|E|F|G>
#     A0 = no mutation, baseline nail run on the pristine tree (expect PASS)
#     B  = prod: isDisplayWashCf body token -> old 15-rune hand roster (ticket 333 shape)
#     C  = TEST self-reference shot: arm 2's actual side "washSetOver(isDisplayWashCf)"
#          rewritten to the expectation side's own expression (category-vs-category)
#     D  = C + B together: self-referenced test against a drifted product
#     E  = prod: "add one Cf member" shot  if r == 0x2061 { return true }  (no-op check)
#     F  = prod: per-rune exclusion for U+2061 (leg's shot B, re-run independently)
#     G  = prod: per-rune inclusion of U+00A0, outside the category (leg's shot C)
set -e
TREE=/d/V335V1/T
WORK=/d/V335V1/mut
REPO="/d/work/workspace/projects plans/Wisp"
LOGS="$REPO/.scratch/wisp/probes/335/v1/logs"
PROD=internal/panel/display_scrub_216.go
TEST=internal/panel/display_scrub_216_test.go
NAIL=TestTheRosterIsClosedAgainstTheRunesTheRulingLeftOut
label="$1"
cd "$TREE" || exit 99
mkdir -p "$WORK" "$LOGS" /d/V335V1/gotmp

cp "$PROD" "$WORK/pristine-prod-$label.go"
cp "$TEST" "$WORK/pristine-test-$label.go"
pre_p=$(sha256sum "$PROD" | cut -d' ' -f1)
pre_t=$(sha256sum "$TEST" | cut -d' ' -f1)
restore() {
	cp "$WORK/pristine-prod-$label.go" "$PROD"
	cp "$WORK/pristine-test-$label.go" "$TEST"
}
trap restore EXIT INT TERM

echo "SHOT=$label"
echo "PRE_HASH_PROD=$pre_p"
echo "PRE_HASH_TEST=$pre_t"

roster_mut() {
	awk 'BEGIN{while ((getline l < "/d/V335V1/mut/m-a-body.txt") > 0) body = body l "\n"}
	     /^[[:space:]]*return unicode\.Is\(unicode\.Cf, r\)[[:space:]]*$/ { printf "%s", body; next }
	     { print }' "$WORK/pristine-prod-$label.go" > "$WORK/mutated-prod-$label.go"
	cp "$WORK/mutated-prod-$label.go" "$PROD"
}
case "$label" in
A0) ;;
B) roster_mut ;;
C)
	awk '/^[[:space:]]*product := washSetOver\(isDisplayWashCf\)[[:space:]]*$/ {
			print "\tproduct := washSetOver(func(r rune) bool { return unicode.Is(unicode.Cf, r) })"; next }
		{ print }' "$WORK/pristine-test-$label.go" > "$WORK/mutated-test-$label.go"
	cp "$WORK/mutated-test-$label.go" "$TEST"
	;;
D)
	awk '/^[[:space:]]*product := washSetOver\(isDisplayWashCf\)[[:space:]]*$/ {
			print "\tproduct := washSetOver(func(r rune) bool { return unicode.Is(unicode.Cf, r) })"; next }
		{ print }' "$WORK/pristine-test-$label.go" > "$WORK/mutated-test-$label.go"
	cp "$WORK/mutated-test-$label.go" "$TEST"
	roster_mut
	;;
E)
	awk '/^func isDisplayWashCf\(r rune\) bool \{$/ {
			print; print "\tif r == 0x2061 { return true }"; next }
		{ print }' "$WORK/pristine-prod-$label.go" > "$WORK/mutated-prod-$label.go"
	cp "$WORK/mutated-prod-$label.go" "$PROD"
	;;
F)
	awk '/^func isDisplayWashCf\(r rune\) bool \{$/ {
			print; print "\tif r == 0x2061 { return false }"; next }
		{ print }' "$WORK/pristine-prod-$label.go" > "$WORK/mutated-prod-$label.go"
	cp "$WORK/mutated-prod-$label.go" "$PROD"
	;;
G)
	awk '/^func isDisplayWashCf\(r rune\) bool \{$/ {
			print; print "\tif r == 0x00A0 { return true }"; next }
		{ print }' "$WORK/pristine-prod-$label.go" > "$WORK/mutated-prod-$label.go"
	cp "$WORK/mutated-prod-$label.go" "$PROD"
	;;
*) echo "unknown shot $label"; exit 98 ;;
esac

echo "MUT_HASH_PROD=$(sha256sum "$PROD" | cut -d' ' -f1)"
echo "MUT_HASH_TEST=$(sha256sum "$TEST" | cut -d' ' -f1)"
echo "PROD_CHANGED_LINES=$(diff "$WORK/pristine-prod-$label.go" "$PROD" | grep -c '^[<>]' || true)"
echo "TEST_CHANGED_LINES=$(diff "$WORK/pristine-test-$label.go" "$TEST" | grep -c '^[<>]' || true)"

export TMP=/d/V335V1/gotmp TEMP=/d/V335V1/gotmp GOTMPDIR=/d/V335V1/gotmp
set +e
go test -v ./internal/panel/ -count=1 -timeout 30m -run "$NAIL" >"$WORK/out-$label.log" 2>&1
rc_test=$?
set -e
echo "TEST_RC=$rc_test"
cp "$WORK/out-$label.log" "$LOGS/mut-$label.log"

restore
trap - EXIT INT TERM
post_p=$(sha256sum "$PROD" | cut -d' ' -f1)
post_t=$(sha256sum "$TEST" | cut -d' ' -f1)
echo "RESTORE_HASH_PROD=$post_p"
echo "RESTORE_HASH_TEST=$post_t"
[ "$pre_p" = "$post_p" ] && echo "RESTORE_PROD=YES" || echo "RESTORE_PROD=NO"
[ "$pre_t" = "$post_t" ] && echo "RESTORE_TEST=YES" || echo "RESTORE_TEST=NO"
echo "RUN=$(grep -c '^=== RUN' "$WORK/out-$label.log" || true)"
echo "PASS_TOP=$(grep -c '^--- PASS' "$WORK/out-$label.log" || true)"
echo "FAIL_TOP=$(grep -c '^--- FAIL' "$WORK/out-$label.log" || true)"
echo "SKIP=$(grep -c '^--- SKIP' "$WORK/out-$label.log" || true)"
echo "ARM2_READING:"
grep 'arm 2:' "$WORK/out-$label.log" | sed 's/^/  /'
echo "RED_ASSERTION_LINES:"
grep -Eo 'display_scrub_216_test\.go:[0-9]+' "$WORK/out-$label.log" | sort -u | sed 's/^/  /'
