#!/bin/sh
# 161-r3 / ticket 161 AC#7② - "the repo-level rulers run LAST, and every line they
# print is attributed to the person who made it".
#
# WHAT HAPPENED, twice, and why a script exists for it
# ---------------------------------------------------------------------------
# The 票 158 acceptance bench measured `gofumpt -l . tools/d22scan tools/mockllm`
# empty at 20:31 (docs/evidence/s1/158-gate-scope-blind-spot-r1-accept-r1.md §1.3,
# artifacts probes/158/accept-r1/gofumpt-{full,narrow}.txt = 0 bytes) and graded
# AC#5③ on it. At 20:49 that same program committed the three mutation fixtures
# that make the reading non-empty (506cbae: mut/guard.no{1,2,3}.go). Nothing was
# misread at either moment - the ORDER had no gate: the ruler ran before the
# artifacts landed, and nothing could tell the next person whether a red line was
# (a) a committed artifact CI will really see, or (b) somebody's bench file that
# does not exist in a CI checkout at all.
#
# THAT SECOND DISTINCTION IS WHAT THIS SCRIPT ADDS. `gofumpt -l .` walks the
# working tree and knows nothing about git (measured on 2026-09-26: rc=2 / 8 lines
# here, of which ZERO are tracked). So on a shared tree it cannot answer "would CI
# go red", and the tempting wrong answers are both worse than the question:
# deleting someone else's fixtures, or reading a local number and calling it the
# gate's verdict.
#
# USAGE - and the ordering rule in one sentence
# ---------------------------------------------------------------------------
# Run this AFTER every artifact of the current ticket is committed. That IS the
# rule AC#7② asks for; the reason it is executable and not a slogan is the
# `git ls-files` classification below, which only makes sense on a tree where
# "mine" has already become "tracked".
#
#	sh .scratch/wisp/probes/161/r3/gate-order.sh
#
# Exit codes:
#	0  no TRACKED offender in any of the three rulers (this is the CI-faithful read)
#	1  a TRACKED offender exists - CI would go red on the pushed bytes
#	2  a ruler could not run at all (missing binary, wrong directory, ...)
# UNTRACKED offenders never change the exit code. They are printed, with the
# probe directory that owns them, so the next agent does not "fix" a red that is
# not theirs and does not delete a fixture that is somebody else's evidence.
#
# WHAT DELETING THIS SCRIPT'S KEY INGREDIENT LETS LEAK AGAIN (AC#7② asks for this
# sentence in the script itself, not in prose)
# ---------------------------------------------------------------------------
# Remove the classify-tracked-or-not step (the `is_tracked` function and its two
# call sites) and the output collapses to "gofumpt rc=2, 8 lines" - which is
# indistinguishable from the reading that invalidated the 票 158 AC#5③ verdict.
# Concretely, three failures come back:
#	1. a committed dirty .go cannot be told from a bench fixture, so a ticket can
#	   be graded green on "these 8 lines aren't mine" or red on "someone's fixture
#	   is dirty", and both sentences look like evidence;
#	2. rc=2 from gofumpt means "a file would not parse", NOT "some files need
#	   formatting" - without the per-line dump nobody is told which of the two
#	   happened (measured: 5 of today's 8 lines are parse errors of that shape);
#	3. the "run last" order itself becomes unfalsifiable: nothing records that the
#	   rulers ran after the commits, so the next program can read a stale line count
#	   again exactly like the one that expired at 20:49.

set -eu

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
root=$(CDPATH= cd -- "$here/../../../.." && pwd)
cd "$root"

GOPATH_BIN=$(go env GOPATH)
GOFUMPT="$GOPATH_BIN/bin/gofumpt"
if [ -f "$GOPATH_BIN/bin/gofumpt.exe" ]; then GOFUMPT="$GOPATH_BIN/bin/gofumpt.exe"; fi
[ -x "$GOFUMPT" ] || { echo "gate-order.sh: no gofumpt at $GOFUMPT - refusing to answer with a ruler that did not run" >&2; exit 2; }

TRACKED_TOTAL=0

# is_tracked <repo-relative path> -> 0 if git's index holds it
# This is THE ingredient. Deleting it makes every reading below unattributable.
is_tracked() {
    git ls-files --error-unmatch -- "$1" >/dev/null 2>&1
}

# norm <path-as-the-tool-printed-it> -> forward slashes, stripped of :line:col tails
norm() {
    printf '%s\n' "$1" | tr '\\' '/' | sed -e 's/:[0-9][0-9]*:[0-9][0-9]*$//' -e 's/:[0-9][0-9]*:[0-9][0-9]*: .*//'
}

# owner_hint <path> -> which bench or tree this came from, for the human line
owner_hint() {
    case $1 in
    .scratch/wisp/probes/*) printf '%s' "bench artifact: $(printf '%s' "$1" | cut -d/ -f1-5)" ;;
    frontend/*)             printf '%s' "owned by another session (ticket 169's tree)" ;;
    *)                      printf '%s' "-" ;;
    esac
}

echo "gate-order.sh: repo root = $root"
echo "gate-order.sh: $(printf '%s' "$("$GOFUMPT" --version 2>&1)") ($(date '+%Y-%m-%d %H:%M %z'))"
echo

# ---- ruler 1: gofumpt ------------------------------------------------------
echo "== ruler 1/3: $GOFUMPT -l . tools/d22scan tools/mockllm"
raw=$("$GOFUMPT" -l . tools/d22scan tools/mockllm 2>&1 || true)
if [ -z "$raw" ]; then
    echo "   empty - nothing to attribute"
else
    tracked_lines=0
    untracked_lines=0
    while IFS= read -r line; do
        [ -n "$line" ] || continue
        p=$(norm "$line")
        if is_tracked "$p"; then
            who=$(git log -1 --format='%h %ad %an' --date=format:'%m-%d %H:%M' -- "$p")
            echo "   TRACKED   $p   (last committed by: $who)   <== CI WILL BE RED ON THIS"
            tracked_lines=$((tracked_lines + 1))
        else
            echo "   UNTRACKED $p   ($(owner_hint "$p"))   <== not in a CI checkout"
            untracked_lines=$((untracked_lines + 1))
        fi
    done <<EOF
$raw
EOF
    echo "   tracked=$tracked_lines untracked=$untracked_lines (gofumpt's own rc was NOT printed: a line dump is the reading, rc=2 here means 'a file would not parse', not 'files need formatting')"
    TRACKED_TOTAL=$((TRACKED_TOTAL + tracked_lines))
fi
echo

# ---- ruler 2: the D22 scan -------------------------------------------------
echo "== ruler 2/3: sh scripts/d22scan.sh"
scan_rc=0
scan_out=$(sh scripts/d22scan.sh 2>&1) || scan_rc=$?
findings=$(printf '%s\n' "$scan_out" | grep -E '^[A-Za-z0-9_./-]+:[0-9]+: \[' || true)
if [ -z "$findings" ]; then
    echo "   0 finding lines; script rc=$scan_rc"
else
    n=0
    while IFS= read -r line; do
        [ -n "$line" ] || continue
        p=$(norm "$line")
        if is_tracked "$p"; then
            who=$(git log -1 --format='%h %ad %an' --date=format:'%m-%d %H:%M' -- "$p")
            echo "   TRACKED   $line   (last committed by: $who)   <== CI WILL BE RED ON THIS"
            n=$((n + 1))
        else
            echo "   UNTRACKED $line   ($(owner_hint "$p"))"
        fi
    done <<EOF
$findings
EOF
    TRACKED_TOTAL=$((TRACKED_TOTAL + n))
    echo "   script rc=$scan_rc, finding lines above"
fi
echo

# ---- ruler 3: go vet (both modules) ---------------------------------------
echo "== ruler 3/3: go vet ./...  (root module) and go vet ./  (tools/d22scan)"
vet_rc=0
(cd "$root" && go vet ./... 2>&1) | sed 's/^/   root: /' || vet_rc=$?
(cd "$root/tools/d22scan" && go vet ./ 2>&1) | sed 's/^/   d22scan: /' || vet_rc=$?
echo "   note: the root-module form 'go vet ./tools/d22scan/' is NOT runnable (separate"
echo "         go.mod; measured rc=1 'main module does not contain package'), so the"
echo "         in-module form above is the one CI's step 8 uses."
echo

echo "gate-order.sh: TRACKED offenders across the three rulers = $TRACKED_TOTAL"
if [ "$TRACKED_TOTAL" -gt 0 ]; then
    echo "gate-order.sh: RED - a committed artifact fails a repo-level ruler; fix the bytes, do not delete someone's bench"
    exit 1
fi
echo "gate-order.sh: clean on the committed set (UNTRACKED lines above, if any, are local bench bytes CI never sees)"
