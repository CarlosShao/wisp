#!/bin/sh
# 325-r1 AC#2 counter-control: physically fold the exported (fixed) tree's line back
# into the mixed shape, prove the four internal/risk junction cases go red, restore the
# exact committed bytes, prove green again, and print the hashes at each step.
# Runs only on an OUTSIDE-THE-REPO git archive export. No repo file is touched here.
set -u
LOGD="/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/325/r1/logs"
REPO="/d/work/workspace/projects plans/Wisp"
F=tools/d22scan/runtests.sh
T=/c/tmp/wisp325r1/revert
RUN='^(TestSyncAncestorActableLegFailsClosedOnMovedAncestor116|TestSyncFirstActableLegStillRefusesHomeRewriteAtVerdict116|TestSyncAncestorLegReResolvesThroughJunctionToRealTree|TestSyncAncestorLegNotExemptedStaysFailClosed)$'

rm -rf "$T"
cp -a /c/tmp/wisp325r1/post "$T"
cd "$T" || exit 9

echo "== step 1: exported file as committed =="
sha256sum "$F" | tee "$LOGD/a2-90-sha-committed.txt"
sed -n '115p' "$F" | od -c | sed -n '1,12p'

echo "== step 2: physically fold :115 back to the mixed shape =="
sed -i 's|}\\\\d22scan-|}/d22scan-|' "$F"
sed -n '115p' "$F" | od -c | sed -n '1,12p'
sha256sum "$F" | tee "$LOGD/a2-91-sha-reverted-mixed.txt"
# anti-tautology: the mutation MUST have landed, or step 3 would be a green nobody
# earned. Counting the mixed spelling and the native spelling on the one line.
echo "mixed-shape-hits=$(grep -c '"}\/d22scan-' "$F") native-shape-hits=$(grep -c '}\\\\d22scan-' "$F")"

echo "== step 3: run the four through the mixed-shape runner (expect red) =="
env -u GOTMPDIR -u WISP_TEST_GOTMP_ROOT sh "$F" ./internal/risk/ -run "$RUN" \
    > "$LOGD/a2-revert-mixed-red.txt" 2>&1
echo "step3-rc=$?"

echo "== step 4: restore the committed bytes =="
git -C "$REPO" show 428ef21e:"$F" > "$F"
sha256sum "$F" | tee "$LOGD/a2-92-sha-restored.txt"
sed -n '115p' "$F" | od -c | sed -n '1,12p'
echo "after-restore: mixed-shape-hits=$(grep -c '"}\/d22scan-' "$F") native-shape-hits=$(grep -c '}\\\\d22scan-' "$F")"

echo "== step 5: rerun the four (expect green) =="
env -u GOTMPDIR -u WISP_TEST_GOTMP_ROOT sh "$F" ./internal/risk/ -run "$RUN" \
    > "$LOGD/a2-restore-green-four.txt" 2>&1
echo "step5-rc=$?"
