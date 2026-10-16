#!/bin/bash
# 317-v1 instrument for the FOURTH same-shape place (tools/d22scan/runtests.sh:92,
# mixed form today: native base + slash tail). Measures, at the file's real CI call
# sites, whether the mixed shape detonates:
#   S1 = the call site CI actually uses (ci.yml:831): ./internal/risk/ -run TestPathResolverJunctionWindows
#        (link/target built through filepath.Join => expected to survive a mixed root)
#   S2 = the same runner, same root, pointed at the four 116/105 names
        #        (expected to go red => the FORM detonates, only the SCOPE hides it)
# Usage: bash r8-runtests.sh <tree-dir> <out-log>
set -u
tree="$1"
out="$2"
{
    echo "== 317-v1 runtests.sh mixed-root measurements (shipped, unmutated tree) =="
    echo "date=$(date '+%Y-%m-%d %H:%M:%S %z')"
    echo "tree=$tree"
    echo "--- the composition line as shipped, byte-exact ---"
    sed -n '92p' "$tree/tools/d22scan/runtests.sh" | od -c | head -8
    cd "$tree" || exit 9
    echo "### S1: runtests.sh ./internal/risk/ -run TestPathResolverJunctionWindows"
    env -u GOTMPDIR -u WISP_TEST_GOTMP_ROOT bash tools/d22scan/runtests.sh ./internal/risk/ -run TestPathResolverJunctionWindows >"$out.s1" 2>&1
    echo "S1_RC=$?"
    tail -12 "$out.s1"
    echo "### S2: runtests.sh ./internal/risk/ -run <the four 116/105 names>"
    env -u GOTMPDIR -u WISP_TEST_GOTMP_ROOT bash tools/d22scan/runtests.sh ./internal/risk/ -run 'TestSyncAncestorActableLegFailsClosedOnMovedAncestor116|TestSyncFirstActableLegStillRefusesHomeRewriteAtVerdict116|TestSyncAncestorLegReResolvesThroughJunctionToRealTree|TestSyncAncestorLegNotExemptedStaysFailClosed' 2>&1
    echo "S2_RC=$?"
} >"$out" 2>&1
echo "log=$out"
grep -n 'GOTMPDIR=' "$out" | head -4 | cut -c1-220
grep -n 'S1_RC\|S2_RC\|top-level: PASS' "$out" | cut -c1-260
grep -n 'Invalid switch\|precondition:' "$out" | head -6 | cut -c1-200
