#!/bin/sh
# 322-v1 (acceptance leg) POSIX rerun, second attempt. v1's container() shipped a quoting bug:
# the -run filter was concatenated INSIDE the single-quoted `sh -c` string, so the '|' between
# test names was eaten as a shell pipe by the container's sh (rc=127, "not found", "no Go files
# in /work" = table fault, zero readings counted). This version quotes the filter the way
# 322-r1b's mut-posix.sh did: -run "<pat>" as literal double quotes inside the container command.
set -e
REPO='D:/work/workspace/projects plans/Wisp'
S=/d/tmp/probe322v1
LOGS="$REPO/.scratch/wisp/probes/322/v1/logs"
export MSYS_NO_PATHCONV=1
mkdir -p "$S" "$LOGS"
HEAD=$(git -C "$REPO" rev-parse HEAD)
NEWFILES='internal/winsec/openleaf_322_other_test.go internal/winsec/openleaf_flag_322_other_test.go internal/winsec/openleaf_overwrite_322_test.go internal/winsec/openleaf_guard_322_windows_test.go'
PAT='TestAC1POSIXOpenThroughASwappedLeafLinkNeverTruncatesAForeignFile322r1|TestAC322LeafFlagRefusesALeafLinkWithoutTouchingTheTarget322r1|TestAC2OverwritePromiseHoldsForEveryProductionCaller322r1'
: >"$LOGS/mut-posix-v1b-summary.txt"
say() { echo "$@" | tee -a "$LOGS/mut-posix-v1b-summary.txt"; }
say "### 322-v1 POSIX rerun (quote-fixed), started $(date '+%F %T %z') HEAD=$HEAD"

newtree() {
	tree=$(mktemp -d "$S/$1.XXXXXX")
	( cd "$REPO" && git archive "${2:-HEAD}" | tar -x -C "$tree" )
	for f in $NEWFILES; do cp "$REPO/$f" "$tree/$f"; done
	printf 'TREE %s (rev=%s): .go=%s winsec/*.go=%s\n' "$tree" "${2:-HEAD}" \
		"$(find "$tree" -name '*.go' | wc -l | tr -d ' ')" \
		"$(find "$tree/internal/winsec" -name '*.go' | wc -l | tr -d ' ')" | tee -a "$LOGS/mut-posix-v1b-summary.txt"
	return 0
}

anchor() {
	n=$(grep -c -- "$3" "$1/$2" || true)
	[ "$n" = "1" ] || { say "ANCHOR NOT UNIQUE for $4: '$3' hit $n time(s) in $2 - ABORT"; exit 3; }
	say "ANCHOR $4: $(grep -n -- "$3" "$1/$2" | head -1)"
}

runfilter=$1 # empty means whole package
container() { # <label> <tree>
	label=$1; tree=$2; rc=0
	docker run --rm -e CGO_ENABLED=0 -e GOPROXY=off -e GOFLAGS=-mod=mod \
		-v "$tree":/work -v probe322v1build:/root/.cache/go-build \
		-v probe322mod:/gomod:ro -w /work golang:1.27 \
		sh -c 'export GOMODCACHE=/gomod GOTMPDIR=/tmp/gotmp322v1b; mkdir -p "$GOTMPDIR"; go test -count=1 -v -run "'"$runfilter"'" ./internal/winsec/ > /work/'"$label"'.log 2>&1; rc=$?; echo "go-test rc=$rc" >> /work/'"$label"'.log; exit $rc' \
		>"$LOGS/$label.log" 2>&1 || rc=$?
	cp "$tree/$label.log" "$LOGS/$label.body.log" 2>/dev/null || say "NO BODY FILE at $tree/$label.log"
	b="$LOGS/$label.body.log"
	{
		echo "=== $label rc=$rc filter='$runfilter'"
		printf 'RUN=%s topPASS=%s topFAIL=%s topSKIP=%s\n' \
			"$(grep -c '^=== RUN' "$b" || true)" "$(grep -c '^--- PASS' "$b" || true)" \
			"$(grep -c '^--- FAIL' "$b" || true)" "$(grep -c '^--- SKIP' "$b" || true)"
		echo "tooling-word check: walk=$(grep -ci walk "$b" || true) buildfailed=$(grep -ci 'build failed' "$b" || true) noteststorun=$(grep -ci 'no tests to run' "$b" || true) setupfailed=$(grep -ci 'setup failed' "$b" || true)"
		grep -h 'AC#1 paired reading' "$b" || echo "NO AC#1 PAIRED ROW in $label"
		echo "--- top-level result lines:"
		grep -hE '^--- (FAIL|PASS|SKIP)' "$b" || true
		echo "--- RED lines verbatim:"
		grep -h 'RED' "$b" || echo "(no RED line)"
	} | tee -a "$LOGS/mut-posix-v1b-summary.txt"
}

step_base() {
	BASE=$(newtree base1b '8baec646^' | awk '{print $2; exit}')
	cat >"$BASE/internal/winsec/zz322v1b_stub_pre_ticket_other.go" <<'STUB'
//go:build !windows

package winsec

// Pre-ticket POSIX values, restored so the ticket's own test files compile against
// un-fixed code: the open carried no leaf flag and no pre-open probe existed.
const oNoFollowPlatform = 0

func preOpenLeafGuard(string) error { return nil }
STUB
	container base1b "$BASE"
}

step_fixed() {
	FIXED=$(newtree fixed1b | awk '{print $2; exit}')
	sha256sum "$FIXED/internal/winsec/winsec.go" | sed 's#^#FIXED-BYTES #' | tee -a "$LOGS/mut-posix-v1b-summary.txt"
	container fixed1b "$FIXED"
}

step_m1() {
	M1=$(newtree m1b-teethout | awk '{print $2; exit}')
	W=internal/winsec/winsec.go
	anchor "$M1" "$W" 'os\.OpenFile(resolved, flags|oNoFollowPlatform, perm)' 'M1 open line'
	cp "$M1/$W" "$M1/$W.m1b.bak"
	say "M1-BEFORE sha=$(sha256sum "$M1/$W" | cut -d' ' -f1)"
	awk '{gsub(/flags\|oNoFollowPlatform/,"flags"); print}' "$M1/$W.m1b.bak" > "$M1/$W"
	say "M1-AFTER open-line=$(sed -n '149p' "$M1/$W") hits=$(grep -c 'flags|oNoFollowPlatform' "$M1/$W" || true)"
	container m1b "$M1"
	cp "$M1/$W.m1b.bak" "$M1/$W"
	say "M1-RESTORE sha=$(sha256sum "$M1/$W" | cut -d' ' -f1) (must equal M1-BEFORE)"
}

step_m2() {
	M2=$(newtree m2b-blanket | awk '{print $2; exit}')
	O=internal/winsec/winsec_other.go
	anchor "$M2" "$O" 'func preOpenLeafGuard(string) error { return nil }' 'M2 posix guard body'
	cp "$M2/$O" "$M2/$O.m2b.bak"
	say "M2-BEFORE sha=$(sha256sum "$M2/$O" | cut -d' ' -f1)"
	awk '{gsub(/func preOpenLeafGuard\(string\) error \{ return nil \}/,"func preOpenLeafGuard(path string) error { return fmt.Errorf(\"%w: MUTANT blanket refusal on every name, fixed or not: %s\", ErrUnresolvedPath, path) }"); print}' "$M2/$O.m2b.bak" > "$M2/$O"
	say "M2-AFTER body=$(grep -n 'func preOpenLeafGuard' "$M2/$O")"
	container m2b "$M2"
	cp "$M2/$O.m2b.bak" "$M2/$O"
	say "M2-RESTORE sha=$(sha256sum "$M2/$O" | cut -d' ' -f1) (must equal M2-BEFORE)"
}

case "${STEP:-all}" in
	base) step_base ;;
	fixed) step_fixed ;;
	m1) step_m1 ;;
	m2) step_m2 ;;
	full) runfilter=''; FIXED=$(newtree full1b | awk '{print $2; exit}'); container full1b "$FIXED" ;;
	all) step_base; step_fixed; step_m1; step_m2 ;;
esac
say "### done $(date '+%F %T %z')"
