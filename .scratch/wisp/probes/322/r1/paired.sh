#!/bin/sh
# 322-r1 paired readings driver.
#
# WHAT EACH RUN MEANS
#   base     = git archive HEAD (the ticket-322 legs DO NOT EXIST yet) + this leg's test
#              files + a stub giving the pre-ticket semantics (oNoFollowPlatform = 0,
#              preOpenLeafGuard answers "no link"). This is the teeth-out half of the
#              paired reading: AC#1's numerator must be > 0 here.
#   fixed    = the same archive with winsec.go / winsec_other.go / winsec_windows.go
#              replaced by the delivered bytes. AC#1's numerator must be 0 AND the new
#              leg must have spoken (liveness > 0).
#   notrunc  = fixed with os.O_TRUNC physically deleted from winsec.go's overwrite
#              branch: AC#2's positive control, every caller pin must go red.
#
# RULES THIS SCRIPT FOLLOWS (each one is a recorded accident, not a preference)
#   - snapshots come from mktemp -d and are NEVER deleted (临时件只建不删); the作废 runs
#     of 120-r3 vanished with an rm -rf, so nothing here removes anything;
#   - the tree is exported to the D: volume (Git Bash /tmp is C: Temp and mounts empty);
#   - .go is counted before any measurement, so a 0 can never be read as a reading;
#   - mutations are PHYSICAL, with sha256 before/after and a .bak left on disk - never
#     -overlay, which is invisible to a leg that reads the disk at run time;
#   - every command writes its own log and takes rc from $? directly, never PIPESTATUS;
#   - GOPROXY=off + the host module cache, because a container without network makes a
#     missing module look like a test failure.
set -e

REPO='D:/work/workspace/projects plans/Wisp'
S=/d/tmp/probe322r1
LOGS="$REPO/.scratch/wisp/probes/322/r1/logs"
export MSYS_NO_PATHCONV=1
mkdir -p "$S" "$LOGS"

PROD='internal/winsec/winsec.go internal/winsec/winsec_other.go internal/winsec/winsec_windows.go'
NEWFILES='internal/winsec/openleaf_322_other_test.go internal/winsec/openleaf_overwrite_322_test.go internal/winsec/openleaf_guard_322_windows_test.go'

# mkstub <tree> : give the pre-ticket tree the two symbols the new tests reference, with
# their pre-ticket VALUES. One file per platform so neither tree carries a dead import.
mkstub() {
	tree=$1
	mkdir -p "$tree/internal/winsec"
	cat >"$tree/internal/winsec/zz322r1_stub_pre_ticket_other.go" <<'STUB'
//go:build !windows

package winsec

// zz322r1_stub_pre_ticket_other.go restores the POSIX semantics of the two symbols the
// ticket 322 tests reference to their PRE-TICKET values, so the same test files can be
// run against the un-fixed tree. oNoFollowPlatform = 0 is exactly what winsec.go did at
// this anchor: os.OpenFile(resolved, flags, perm), following any leaf link.
const oNoFollowPlatform = 0

func preOpenLeafGuard(string) error { return nil }
STUB
	cat >"$tree/internal/winsec/zz322r1_stub_pre_ticket_windows.go" <<'STUB'
//go:build windows

package winsec

// The Windows half of the same restoration: no leaf probe at all, which is what the
// shipped platform did before this ticket.
const oNoFollowPlatform = 0

func preOpenLeafGuard(string) error { return nil }
STUB
}

# newtree <name> : export HEAD into a fresh directory that is never removed.
newtree() {
	name=$1
	tree=$(mktemp -d "$S/$name.XXXXXX")
	(
		cd "$REPO" && git archive HEAD | tar -x -C "$tree"
	)
	for f in $NEWFILES; do
		cp "$REPO/$f" "$tree/$f"
	done
	echo "$tree"
}

# countgo <tree> : prove the export is not empty before trusting any number from it.
countgo() {
	tree=$1
	all=$(find "$tree" -name '*.go' | wc -l | tr -d ' ')
	ws=$(find "$tree/internal/winsec" -name '*.go' | wc -l | tr -d ' ')
	printf 'TREE %s: .go=%s internal/winsec/*.go=%s\n' "$tree" "$all" "$ws"
	[ "$all" -gt 500 ] || {
		echo "EMPTY EXPORT: .go=$all, refusing to read anything from $tree" >&2
		exit 4
	}
}

# containerrun <label> <tree> <runfilter> : one linux container, same image, same volume
# style as 120-r3, writing <label>.log and a three-line verdict block.
containerrun() {
	label=$1
	tree=$2
	pat=$3
	out="$LOGS/$label.log"
	rc=0
	docker run --rm -e CGO_ENABLED=0 -e GOPROXY=off -e GOFLAGS=-mod=mod \
		-v "$tree":/work -v probe322r1build:/root/.cache/go-build \
		-v probe322mod:/gomod:ro -w /work golang:1.27 \
		sh -c 'export GOMODCACHE=/gomod GOTMPDIR=/tmp/gotmp322r1; mkdir -p "$GOTMPDIR"; go test -count=1 -v -run "'"$pat"'" ./internal/winsec/ > /work/'"$label"'.log 2>&1; rc=$?; echo "go-test rc=$rc" >> /work/'"$label"'.log; exit $rc' \
		>"$out" 2>&1 || rc=$?
	body="$tree/$label.log"
	cp "$body" "$LOGS/$label.body.log" 2>/dev/null || echo "NO BODY FILE in $body (see $out for the docker error)"
	{
		echo "=== $label runfilter='$pat' docker-and-test rc=$rc"
		printf 'RUN=%s (=== RUN rows, subtests included)\n' "$(grep -c '=== RUN' "$LOGS/$label.body.log" 2>/dev/null || true)"
		printf 'top-level (anchored ^---): PASS=%s FAIL=%s SKIP=%s\n' \
			"$(grep -c -- '^--- PASS' "$LOGS/$label.body.log" 2>/dev/null || true)" \
			"$(grep -c -- '^--- FAIL' "$LOGS/$label.body.log" 2>/dev/null || true)" \
			"$(grep -c -- '^--- SKIP' "$LOGS/$label.body.log" 2>/dev/null || true)"
		printf 'all-rows (unanchored): PASS=%s FAIL=%s SKIP=%s\n' \
			"$(grep -c -- '--- PASS' "$LOGS/$label.body.log" 2>/dev/null || true)" \
			"$(grep -c -- '--- FAIL' "$LOGS/$label.body.log" 2>/dev/null || true)" \
			"$(grep -c -- '--- SKIP' "$LOGS/$label.body.log" 2>/dev/null || true)"
		grep -h 'AC#1 paired reading' "$LOGS/$label.body.log" 2>/dev/null || echo "NO AC#1 PAIRED ROW in $label"
		grep -h "go-test rc=" "$LOGS/$label.body.log" 2>/dev/null || true
	} | tee -a "$LOGS/summary.txt"
}

: >"$LOGS/summary.txt"
echo "### 322-r1 paired readings, started $(date '+%F %T %z')" | tee "$LOGS/summary.txt"
git -C "$REPO" rev-parse HEAD >"$LOGS/head.txt" 2>&1
echo "HEAD=$(cat "$LOGS/head.txt")" | tee -a "$LOGS/summary.txt"

# ---------- tree 1: base (pre-ticket) ----------
BASE=$(newtree base)
mkstub "$BASE"
countgo "$BASE" | tee -a "$LOGS/summary.txt"

# ---------- tree 2: fixed (delivered bytes) ----------
FIXED=$(newtree fixed)
for f in $PROD; do
	cp "$REPO/$f" "$FIXED/$f"
done
countgo "$FIXED" | tee -a "$LOGS/summary.txt"
sha256sum $PROD | sed "s#^#DELIVERED-BYTES #" | tee -a "$LOGS/summary.txt"
for f in $PROD; do
	a=$(sha256sum "$REPO/$f" | cut -d' ' -f1)
	b=$(sha256sum "$FIXED/$f" | cut -d' ' -f1)
	[ "$a" = "$b" ] || {
		echo "SNAPSHOT MISMATCH $f" >&2
		exit 5
	}
done
echo "SNAPSHOT MATCHES REPO for all three product files" | tee -a "$LOGS/summary.txt"

# ---------- tree 3: notrunc (AC#2 positive control) ----------
NOTR=$(newtree notrunc)
for f in $PROD; do
	cp "$REPO/$f" "$NOTR/$f"
done
W=internal/winsec/winsec.go
cp "$NOTR/$W" "$NOTR/$W.322r1-notrunc.bak"
echo "NOTRUNC-BYTES before: $(sha256sum "$NOTR/$W" | cut -d' ' -f1)" | tee -a "$LOGS/summary.txt"
start=$(grep -n '^func PrivateFile(' "$NOTR/$W" | cut -d: -f1)
line=$(awk -v s="$start" 'NR>=s && NR<=s+6 && index($0,"os.O_WRONLY|os.O_CREATE|os.O_TRUNC")>0 {print NR; exit}' "$NOTR/$W")
[ -n "$line" ] || {
	echo "MUTATION DID NOT LAND: no O_TRUNC line inside func PrivateFile ($start..$((start + 6)))" >&2
	exit 3
}
echo "LANDS: tooth notrunc targets line $line, before: $(sed -n "${line}p" "$NOTR/$W")" | tee -a "$LOGS/summary.txt"
sed -i "${line}s/|os.O_TRUNC//" "$NOTR/$W"
echo "LANDS: after: $(sed -n "${line}p" "$NOTR/$W")" | tee -a "$LOGS/summary.txt"
n=$(grep -c 'os.O_TRUNC' "$NOTR/$W" || true)
echo "LANDS: isolated - $n line(s) still carry os.O_TRUNC in $W (winsec.go keeps one in the doc comment)" | tee -a "$LOGS/summary.txt"
countgo "$NOTR" | tee -a "$LOGS/summary.txt"

# ---------- container runs, one table ----------
PAT1='TestAC1POSIXOpenThroughASwappedLeafLinkNeverTruncatesAForeignFile322r1|TestAC322LeafFlagRefusesALeafLinkWithoutTouchingTheTarget322r1|TestAC2OverwritePromiseHoldsForEveryProductionCaller322r1'
containerrun base-ac1 "$BASE" "$PAT1"
containerrun fixed-ac1 "$FIXED" "$PAT1"
containerrun notrunc-ac2 "$NOTR" 'TestAC2OverwritePromiseHoldsForEveryProductionCaller322r1'
containerrun fixed-full "$FIXED" '.'

echo "### done $(date '+%F %T %z')" | tee -a "$LOGS/summary.txt"
echo "TREES (kept on disk, never removed): base=$BASE fixed=$FIXED notrunc=$NOTR" | tee -a "$LOGS/summary.txt"
