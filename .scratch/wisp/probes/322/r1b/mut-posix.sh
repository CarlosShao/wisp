#!/bin/sh
# 322-r1b: teeth check (去牙反控) for the POSIX half of ticket 322, plus the AC#1 paired
# reading at the constant this leg actually ships (leafRaceAttempts322r1 = 3000).
#
# WHY THIS FILE EXISTS (the dead leg's own gap, named in its代提 commit):
#   322-r1 last ran the POSIX AC#1 row at attempts=800 and its newest container run was RED
#   (fixed-ac1: damagedWitnesses=0 but leafLegRefusals=0 -> instrument refused to read).
#   It then raised the constant to 3000 and died before re-running. So the shipped 3000
#   version had NO reading at all. This leg re-runs base/fixed, and adds the two shapes the
#   dispatch asks for: teeth OUT, and the REVERSE shape (blanket refusal).
#
# RULES (each is a recorded accident, copied from 322-r1/paired.sh, not a preference):
#   - mutations are PHYSICAL edits to bytes inside a git-archive export tree;
#     never -overlay (these rulers read the disk at run time, overlay stays green = fake);
#   - trees come from mktemp -d and are NEVER removed; every mutant keeps a .bak;
#   - sha256 before/after each mutation, and the restore run must hash equal;
#   - .go is counted before any number is trusted;
#   - each command writes its own log, rc comes from $? directly, never PIPESTATUS;
#   - GOPROXY=off + the read-only host module cache, because a network-less container makes a
#     missing module look like a test failure ("[setup failed]" is a table fault, not a red).
set -e

REPO='D:/work/workspace/projects plans/Wisp'
S=/d/tmp/probe322r1b
LOGS="$REPO/.scratch/wisp/probes/322/r1b/logs"
export MSYS_NO_PATHCONV=1
mkdir -p "$S" "$LOGS"
HEAD=$(git -C "$REPO" rev-parse HEAD)

NEWFILES='internal/winsec/openleaf_322_other_test.go internal/winsec/openleaf_overwrite_322_test.go internal/winsec/openleaf_guard_322_windows_test.go'
PAT322='TestAC1POSIXOpenThroughASwappedLeafLinkNeverTruncatesAForeignFile322r1|TestAC322LeafFlagRefusesALeafLinkWithoutTouchingTheTarget322r1|TestAC2OverwritePromiseHoldsForEveryProductionCaller322r1'

: >"$LOGS/mut-posix-summary.txt"
say() { echo "$@" | tee -a "$LOGS/mut-posix-summary.txt"; }
say "### 322-r1b POSIX teeth run, started $(date '+%F %T %z') HEAD=$HEAD"

newtree() { # <label> <rev> : HEAD by default; pass 8baec646^ for the pre-ticket tree
	tree=$(mktemp -d "$S/$1.XXXXXX")
	( cd "$REPO" && git archive "${2:-HEAD}" | tar -x -C "$tree" )
	for f in $NEWFILES; do cp "$REPO/$f" "$tree/$f"; done
	printf 'TREE %s (rev=%s): .go=%s winsec/*.go=%s\n' "$tree" "${2:-HEAD}" \
		"$(find "$tree" -name '*.go' | wc -l | tr -d ' ')" \
		"$(find "$tree/internal/winsec" -name '*.go' | wc -l | tr -d ' ')" | tee -a "$LOGS/mut-posix-summary.txt"
	return 0
}

# one-hit anchor check: name the file, the pattern, and refuse to mutate unless it is unique
anchor() { # <tree> <file> <grep-pattern> <label>
	n=$(grep -c -- "$3" "$1/$2" || true)
	[ "$n" = "1" ] || { say "ANCHOR NOT UNIQUE for $4: pattern '$3' hit $n time(s) in $2 - ABORT"; exit 3; }
	line=$(grep -n -- "$3" "$1/$2" | cut -d: -f1 | head -1)
	say "ANCHOR $4: $2:$line before: $(sed -n "${line}p" "$1/$2")"
}

container() { # <label> <tree> <runfilter>
	label=$1; tree=$2; pat=$3
	rc=0
	docker run --rm -e CGO_ENABLED=0 -e GOPROXY=off -e GOFLAGS=-mod=mod \
		-v "$tree":/work -v probe322r1build:/root/.cache/go-build \
		-v probe322mod:/gomod:ro -w /work golang:1.27 \
		sh -c 'export GOMODCACHE=/gomod GOTMPDIR=/tmp/gotmp322r1b; mkdir -p "$GOTMPDIR"; go test -count=1 -v -run "'"$pat"'" ./internal/winsec/ > /work/'"$label"'.log 2>&1; rc=$?; echo "go-test rc=$rc" >> /work/'"$label"'.log; exit $rc' \
		>"$LOGS/$label.log" 2>&1 || rc=$?
	cp "$tree/$label.log" "$LOGS/$label.body.log" 2>/dev/null || say "NO BODY FILE in $tree/$label.log"
	b="$LOGS/$label.body.log"
	{
		echo "=== $label rc=$rc filter='$pat'"
		printf 'RUN=%s topPASS=%s topFAIL=%s topSKIP=%s\n' \
			"$(grep -c '=== RUN' "$b" || true)" "$(grep -c '^--- PASS' "$b" || true)" \
			"$(grep -c '^--- FAIL' "$b" || true)" "$(grep -c '^--- SKIP' "$b" || true)"
		echo "tooling-word check (any hit = table fault, not a reading): walk=$(grep -ci walk "$b" || true) buildfailed=$(grep -ci 'build failed' "$b" || true) noteststorun=$(grep -ci 'no tests to run' "$b" || true) setupfailed=$(grep -ci 'setup failed' "$b" || true)"
		grep -h 'AC#1 paired reading' "$b" || echo "NO AC#1 PAIRED ROW in $label"
		grep -hE '^--- (FAIL|PASS|SKIP)' "$b" || true
		echo "--- red lines verbatim:"
		grep -h 'RED' "$b" || echo "(no RED line)"
	} | tee -a "$LOGS/mut-posix-summary.txt"
}

# ---------- tree 1: base = the commit BEFORE 322-r1's product bytes + pre-ticket symbols ----------
# ⚠ 8baec646^ and NOT HEAD: HEAD already carries the leg's product files, so exporting HEAD and
# adding a stub would be a duplicate declaration, and exporting HEAD alone would be "fixed".
BASE=$(newtree base '8baec646^' | awk '{print $2; exit}')
mkdir -p "$BASE/internal/winsec"
cat >"$BASE/internal/winsec/zz322r1b_stub_pre_ticket_other.go" <<'STUB'
//go:build !windows

package winsec

// Pre-ticket POSIX values, restored so the ticket's own test files compile and run
// against un-fixed code: the open carried no leaf flag and no pre-open probe existed.
const oNoFollowPlatform = 0

func preOpenLeafGuard(string) error { return nil }
STUB

# ---------- tree 2: fixed = HEAD bytes ----------
FIXED=$(newtree fixed | awk '{print $2; exit}')
for f in internal/winsec/winsec.go internal/winsec/winsec_other.go internal/winsec/winsec_windows.go; do
	cp "$REPO/$f" "$FIXED/$f"
done
sha256sum "$FIXED/internal/winsec/winsec.go" "$FIXED/internal/winsec/winsec_other.go" | sed 's#^#FIXED-BYTES #' | tee -a "$LOGS/mut-posix-summary.txt"

container base-322  "$BASE"  "$PAT322"
container fixed-322 "$FIXED" "$PAT322"

# ---------- tree 3: M1 teeth out (drop |oNoFollowPlatform from the open) ----------
M1=$(newtree m1-teethout | awk '{print $2; exit}')
for f in internal/winsec/winsec.go internal/winsec/winsec_other.go internal/winsec/winsec_windows.go; do cp "$REPO/$f" "$M1/$f"; done
W=internal/winsec/winsec.go
anchor "$M1" "$W" 'os\.OpenFile(resolved, flags|oNoFollowPlatform, perm)' 'M1 open line'
cp "$M1/$W" "$M1/$W.m1.bak"
say "M1-BEFORE sha=$(sha256sum "$M1/$W" | cut -d' ' -f1)"
awk '{gsub(/flags\|oNoFollowPlatform/,"flags"); print}' "$M1/$W.m1.bak" > "$M1/$W"
say "M1-AFTER  line: $(sed -n "${line}p" "$M1/$W")"
say "M1-AFTER  sha=$(sha256sum "$M1/$W" | cut -d' ' -f1)  hits-of-const-in-open=$(grep -c 'flags|oNoFollowPlatform' "$M1/$W" || true)"
container m1-322 "$M1" "$PAT322"
cp "$M1/$W.m1.bak" "$M1/$W"
say "M1-RESTORE sha=$(sha256sum "$M1/$W" | cut -d' ' -f1) (must equal M1-BEFORE)"
container m1-restore-322 "$M1" "$PAT322"

# ---------- tree 4: M2 reverse shape = the probe refuses EVERYTHING ----------
M2=$(newtree m2-blanket | awk '{print $2; exit}')
for f in internal/winsec/winsec.go internal/winsec/winsec_other.go internal/winsec/winsec_windows.go; do cp "$REPO/$f" "$M2/$f"; done
O=internal/winsec/winsec_other.go
anchor "$M2" "$O" 'func preOpenLeafGuard(string) error { return nil }' 'M2 posix probe body'
cp "$M2/$O" "$M2/$O.m2.bak"
say "M2-BEFORE sha=$(sha256sum "$M2/$O" | cut -d' ' -f1)"
awk '{gsub(/func preOpenLeafGuard\(string\) error \{ return nil \}/,"func preOpenLeafGuard(path string) error { return fmt.Errorf(\"%w: MUTANT blanket refusal on every name, fixed or not: %s\", ErrUnresolvedPath, path) }"); print}' "$M2/$O.m2.bak" > "$M2/$O"
say "M2-AFTER  body: $(grep -n 'func preOpenLeafGuard' "$M2/$O")"
say "M2-AFTER  sha=$(sha256sum "$M2/$O" | cut -d' ' -f1)"
container m2-322 "$M2" "$PAT322"

say "### done $(date '+%F %T %z')"
say "TREES kept (never removed): base=$BASE fixed=$FIXED m1=$M1 m2=$M2"
