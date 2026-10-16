#!/bin/sh
# 322-v1 (acceptance leg, read-only): POSIX AC#1 paired reading + teeth-out + reverse-shape
# mutations, PHYSICALLY inside git-archive export trees outside the repo, run in golang:1.27
# containers. Never -overlay (these rulers read the disk at run time; overlay stays green = fake).
# Mechanism lineage copied from recorded accidents in 322-r1/paired.sh + 322-r1b/mut-posix.sh:
# base must come from 8baec646^ (HEAD already carries the fix), test files are copied into the
# base tree plus a pre-ticket stub so they compile against un-fixed code; mktemp trees are kept;
# sha256 before/after every physical edit; .go counted before any 0 is trusted; rc from $?;
# GOPROXY=off + the read-only host module-cache volume so "[setup failed]" cannot fake a red.
set -e
REPO='D:/work/workspace/projects plans/Wisp'
S=/d/tmp/probe322v1
LOGS="$REPO/.scratch/wisp/probes/322/v1/logs"
export MSYS_NO_PATHCONV=1
mkdir -p "$S" "$LOGS"
HEAD=$(git -C "$REPO" rev-parse HEAD)
NEWFILES='internal/winsec/openleaf_322_other_test.go internal/winsec/openleaf_flag_322_other_test.go internal/winsec/openleaf_overwrite_322_test.go internal/winsec/openleaf_guard_322_windows_test.go'
PAT='TestAC1POSIXOpenThroughASwappedLeafLinkNeverTruncatesAForeignFile322r1|TestAC322LeafFlagRefusesALeafLinkWithoutTouchingTheTarget322r1|TestAC2OverwritePromiseHoldsForEveryProductionCaller322r1'
: >"$LOGS/mut-posix-v1-summary.txt"
say() { echo "$@" | tee -a "$LOGS/mut-posix-v1-summary.txt"; }
say "### 322-v1 POSIX teeth run, started $(date '+%F %T %z') HEAD=$HEAD"

newtree() { # <label> <rev>
	tree=$(mktemp -d "$S/$1.XXXXXX")
	( cd "$REPO" && git archive "${2:-HEAD}" | tar -x -C "$tree" )
	for f in $NEWFILES; do cp "$REPO/$f" "$tree/$f"; done
	printf 'TREE %s (rev=%s): .go=%s winsec/*.go=%s\n' "$tree" "${2:-HEAD}" \
		"$(find "$tree" -name '*.go' | wc -l | tr -d ' ')" \
		"$(find "$tree/internal/winsec" -name '*.go' | wc -l | tr -d ' ')" | tee -a "$LOGS/mut-posix-v1-summary.txt"
	return 0
}

anchor() { # <tree> <file> <pattern> <label>
	n=$(grep -c -- "$3" "$1/$2" || true)
	[ "$n" = "1" ] || { say "ANCHOR NOT UNIQUE for $4: '$3' hit $n time(s) in $2 - ABORT"; exit 3; }
	say "ANCHOR $4: $(grep -n -- "$3" "$1/$2" | head -1)"
}

container() { # <label> <tree> <pkg-or-run args...>
	label=$1; tree=$2; shift 2
	rc=0
	docker run --rm -e CGO_ENABLED=0 -e GOPROXY=off -e GOFLAGS=-mod=mod \
		-v "$tree":/work -v probe322v1build:/root/.cache/go-build \
		-v probe322mod:/gomod:ro -w /work golang:1.27 \
		sh -c 'export GOMODCACHE=/gomod GOTMPDIR=/tmp/gotmp322v1; mkdir -p "$GOTMPDIR"; go test -count=1 -v '"$*"' ./internal/winsec/ > /work/'"$label"'.log 2>&1; rc=$?; echo "go-test rc=$rc" >> /work/'"$label"'.log; exit $rc' \
		>"$LOGS/$label.log" 2>&1 || rc=$?
	cp "$tree/$label.log" "$LOGS/$label.body.log" 2>/dev/null || say "NO BODY FILE at $tree/$label.log"
	b="$LOGS/$label.body.log"
	{
		echo "=== $label rc=$rc"
		printf 'RUN=%s topPASS=%s topFAIL=%s topSKIP=%s\n' \
			"$(grep -c '^=== RUN' "$b" || true)" "$(grep -c '^--- PASS' "$b" || true)" \
			"$(grep -c '^--- FAIL' "$b" || true)" "$(grep -c '^--- SKIP' "$b" || true)"
		echo "tooling-word check (any hit = table fault, not a reading): walk=$(grep -ci walk "$b" || true) buildfailed=$(grep -ci 'build failed' "$b" || true) noteststorun=$(grep -ci 'no tests to run' "$b" || true) setupfailed=$(grep -ci 'setup failed' "$b" || true)"
		grep -h 'AC#1 paired reading' "$b" || echo "NO AC#1 PAIRED ROW in $label"
		echo "--- top-level result lines:"
		grep -hE '^--- (FAIL|PASS|SKIP)' "$b" || true
		echo "--- RED lines verbatim:"
		grep -h 'RED' "$b" || echo "(no RED line)"
	} | tee -a "$LOGS/mut-posix-v1-summary.txt"
}

# ---------- tree 1: base = the commit BEFORE 322's product bytes + pre-ticket stubs ----------
BASE=$(newtree base '8baec646^' | awk '{print $2; exit}')
cat >"$BASE/internal/winsec/zz322v1_stub_pre_ticket_other.go" <<'STUB'
//go:build !windows

package winsec

// Pre-ticket POSIX values, restored so the ticket's own test files compile against
// un-fixed code: the open carried no leaf flag and no pre-open probe existed.
const oNoFollowPlatform = 0

func preOpenLeafGuard(string) error { return nil }
STUB
container base-v1 "$BASE" "-run $PAT"

# ---------- tree 2: fixed = HEAD bytes (PAT trio) ----------
FIXED=$(newtree fixed | awk '{print $2; exit}')
for f in internal/winsec/winsec.go internal/winsec/winsec_other.go internal/winsec/winsec_windows.go; do cp "$REPO/$f" "$FIXED/$f"; done
sha256sum "$FIXED/internal/winsec/winsec.go" "$FIXED/internal/winsec/winsec_other.go" | sed 's#^#FIXED-BYTES #' | tee -a "$LOGS/mut-posix-v1-summary.txt"
container fixed-v1 "$FIXED" "-run $PAT"
# whole-package colour on the same fixed tree (container dimension of AC#4)
container fullfixed-v1 "$FIXED" ""

# ---------- tree 3: M1 teeth out (drop |oNoFollowPlatform from the open) ----------
M1=$(newtree m1-teethout | awk '{print $2; exit}')
W=internal/winsec/winsec.go
anchor "$M1" "$W" 'os\.OpenFile(resolved, flags|oNoFollowPlatform, perm)' 'M1 open line'
cp "$M1/$W" "$M1/$W.m1.bak"
say "M1-BEFORE sha=$(sha256sum "$M1/$W" | cut -d' ' -f1)"
awk '{gsub(/flags\|oNoFollowPlatform/,"flags"); print}' "$M1/$W.m1.bak" > "$M1/$W"
say "M1-AFTER line=$(sed -n '149p' "$M1/$W") hits-in-open=$(grep -c 'flags|oNoFollowPlatform' "$M1/$W" || true)"
container m1-v1 "$M1" "-run $PAT"
cp "$M1/$W.m1.bak" "$M1/$W"
say "M1-RESTORE sha=$(sha256sum "$M1/$W" | cut -d' ' -f1) (must equal M1-BEFORE)"
container m1-restore-v1 "$M1" "-run $PAT"

# ---------- tree 4: M2 reverse shape = POSIX guard refuses EVERY name ----------
M2=$(newtree m2-blanket | awk '{print $2; exit}')
O=internal/winsec/winsec_other.go
anchor "$M2" "$O" 'func preOpenLeafGuard(string) error { return nil }' 'M2 posix guard body'
cp "$M2/$O" "$M2/$O.m2.bak"
say "M2-BEFORE sha=$(sha256sum "$M2/$O" | cut -d' ' -f1)"
awk '{gsub(/func preOpenLeafGuard\(string\) error \{ return nil \}/,"func preOpenLeafGuard(path string) error { return fmt.Errorf(\"%w: MUTANT blanket refusal on every name, fixed or not: %s\", ErrUnresolvedPath, path) }"); print}' "$M2/$O.m2.bak" > "$M2/$O"
say "M2-AFTER body: $(grep -n 'func preOpenLeafGuard' "$M2/$O")"
container m2-v1 "$M2" "-run $PAT"
cp "$M2/$O.m2.bak" "$M2/$O"
say "M2-RESTORE sha=$(sha256sum "$M2/$O" | cut -d' ' -f1) (must equal M2-BEFORE)"

say "### done $(date '+%F %T %z')"
say "TREES kept (never removed): base=$BASE fixed=$FIXED m1=$M1 m2=$M2"
