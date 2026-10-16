#!/bin/sh
# 322-r1 final readings: the container -count=2 four numbers, the -run roster (named one
# by one), and cmd/wisp's consumer colour on the HOST.
#
# cmd/wisp is NOT runnable in this container at all: CGO_ENABLED=0 makes
# github.com/k2-fsa/sherpa-onnx-go-linux fail its build constraints ("build constraints
# exclude all Go files"), and the package then answers [setup failed]. That red is byte-
# identical on the base and the fixed tree, so it says nothing about ticket 322 - it is the
# same tooling condition the ledger already booked for container builds. The host run below
# is therefore the colour that can actually be measured, and the Linux tier of cmd/wisp is
# left named as unmeasured rather than guessed.
REPO='D:/work/workspace/projects plans/Wisp'
LOGS="$REPO/.scratch/wisp/probes/322/r1/logs"
S=/d/tmp/probe322r1
FIXED=$(ls -dt "$S"/fixed.* | head -1)
BASE=$(ls -dt "$S"/base.* | head -1)
: >"$LOGS/final-summary.txt"
say() { echo "$@" | tee -a "$LOGS/final-summary.txt"; }
export GOTMPDIR='C:/Users/swq/AppData/Local/wisp-test-gotmp/322r1gotmp'
mkdir -p "$GOTMPDIR"

# ---------- container: -count=2 whole package on the delivered bytes ----------
rc=0
MSYS_NO_PATHCONV=1 docker run --rm -e CGO_ENABLED=0 -e GOPROXY=off -e GOFLAGS=-mod=mod \
	-v "$FIXED":/work -v probe322r1build:/root/.cache/go-build \
	-v /d/work/base/gopath/pkg/mod:/gomod:ro -w /work golang:1.27 \
	sh -c 'export GOMODCACHE=/gomod GOTMPDIR=/tmp/got2; mkdir -p "$GOTMPDIR"; go test -count=2 -v ./internal/winsec/ > /work/c2.log 2>&1; rc=$?; echo "go-test rc=$rc" >> /work/c2.log; exit $rc' \
	>"$LOGS/c2.docker.log" 2>&1 || rc=$?
cp "$FIXED/c2.log" "$LOGS/c2.body.log" 2>/dev/null || say "NO c2 BODY"
B="$LOGS/c2.body.log"
say "=== container -count=2 -v ./internal/winsec/ (delivered bytes) docker-and-test rc=$rc"
say "    tree=$FIXED  .go=$(find "$FIXED" -name '*.go' | wc -l | tr -d ' ')  winsec/*.go=$(find "$FIXED/internal/winsec" -name '*.go' | wc -l | tr -d ' ')"
say "    RUN=$(grep -c '=== RUN' "$B" || true) topPASS=$(grep -c '^--- PASS' "$B" || true) topFAIL=$(grep -c '^--- FAIL' "$B" || true) topSKIP=$(grep -c '^--- SKIP' "$B" || true)"
say "    allrows PASS=$(grep -c -- '--- PASS' "$B" || true) FAIL=$(grep -c -- '--- FAIL' "$B" || true) SKIP=$(grep -c -- '--- SKIP' "$B" || true)"
say "    AC#1 paired rows (both counts, -count=2 means two samples):"
grep -h 'AC#1 paired reading' "$B" | sed 's/^/      /' | tee -a "$LOGS/final-summary.txt"
say "    tooling-word check inside this log: walk=$(grep -ci 'walk ' "$B" || true) buildfailed=$(grep -c 'build failed' "$B" || true) notestfiles=$(grep -c 'no test files' "$B" || true) setupfailed=$(grep -c 'setup failed' "$B" || true)"

# ---------- container: roster, named one by one ----------
for pat in '322r1$' 'TestAC[1-4]POSIX'; do
	name=roster-$(echo "$pat" | tr -cd 'A-Za-z0-9')
	rc=0
	MSYS_NO_PATHCONV=1 docker run --rm -e CGO_ENABLED=0 -e GOPROXY=off -e GOFLAGS=-mod=mod \
		-v "$FIXED":/work -v probe322r1build:/root/.cache/go-build \
		-v /d/work/base/gopath/pkg/mod:/gomod:ro -w /work golang:1.27 \
		sh -c 'export GOMODCACHE=/gomod GOTMPDIR=/tmp/got3; mkdir -p "$GOTMPDIR"; go test -count=1 -v -run "'"$pat"'" ./internal/winsec/ > /work/'"$name"'.log 2>&1; rc=$?; echo "go-test rc=$rc" >> /work/'"$name"'.log; exit $rc' \
		>"$LOGS/$name.docker.log" 2>&1 || rc=$?
	cp "$FIXED/$name.log" "$LOGS/$name.body.log" 2>/dev/null || say "NO BODY $name"
	B2="$LOGS/$name.body.log"
	say "=== roster runfilter='$pat' docker-and-test rc=$rc"
	say "    RUN=$(grep -c '=== RUN' "$B2" || true) topPASS=$(grep -c '^--- PASS' "$B2" || true) topFAIL=$(grep -c '^--- FAIL' "$B2" || true) topSKIP=$(grep -c '^--- SKIP' "$B2" || true)"
	say "    named top-level cases:"
	grep -- '^--- ' "$B2" | sed 's/^/      /' | tee -a "$LOGS/final-summary.txt"
done

# ---------- host: cmd/wisp consumer colour, base vs fixed ----------
hostrun() {
	label=$1
	tree=$2
	rc=0
	(cd "$tree" && go test -count=1 -v -run 'TestAC3POSIXSecretRouteLinkInsideItsDataRootStillRefused119|TestAC1POSIXSecretRouteSymlinkedHomeBecomesSealable119|TestAC1POSIXSecretRouteSymlinkedXDGConfigHomeBecomesSealable119' ./cmd/wisp/) >"$LOGS/$label.body.log" 2>&1 || rc=$?
	say "=== host $label go-test rc=$rc tree=$tree"
	say "    $(grep -E '^(ok|FAIL|--- )' "$LOGS/$label.body.log" | tr '\n' '|')"
}
hostrun cmdwisp-host-base "$BASE"
hostrun cmdwisp-host-fixed "$FIXED"
say "### final readings done $(date '+%F %T %z')"
