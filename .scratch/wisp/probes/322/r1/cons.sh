#!/bin/sh
# 322-r1 AC#4 consumer colour, base vs fixed on ONE table.
#
# Why this file exists separately: the first attempt mounted the volume probe322mod at
# /gomod, which turned out to be EMPTY, so every consumer package answered
# `[setup failed]` with `mkdir /gomod/cache: read-only file system` and a "go: downloading"
# line. That is a TOOLING failure, not a colour: base and fixed failed identically, and
# per this repository's rule a red carrying `setup failed` is never counted as a tooth.
# The recipe that works is the HOST module cache mounted read-only plus GOPROXY=off.
REPO='D:/work/workspace/projects plans/Wisp'
LOGS="$REPO/.scratch/wisp/probes/322/r1/logs"
S=/d/tmp/probe322r1
FIXED=$(ls -dt "$S"/fixed.* | head -1)
BASE=$(ls -dt "$S"/base.* | head -1)
: >"$LOGS/cons-summary.txt"
say() { echo "$@" | tee -a "$LOGS/cons-summary.txt"; }

consrun() {
	label=$1
	tree=$2
	pkgs=$3
	rc=0
	MSYS_NO_PATHCONV=1 docker run --rm -e CGO_ENABLED=0 -e GOPROXY=off -e GOFLAGS=-mod=mod \
		-v "$tree":/work -v probe322r1build:/root/.cache/go-build \
		-v /d/work/base/gopath/pkg/mod:/gomod:ro -w /work golang:1.27 \
		sh -c 'export GOMODCACHE=/gomod GOTMPDIR=/tmp/got; mkdir -p "$GOTMPDIR"; go test -count=1 '"$pkgs"' > /work/'"$label"'.log 2>&1; rc=$?; echo "go-test rc=$rc" >> /work/'"$label"'.log; exit $rc' \
		>"$LOGS/$label.docker.log" 2>&1 || rc=$?
	cp "$tree/$label.log" "$LOGS/$label.body.log" 2>/dev/null || say "NO BODY for $label"
	say "=== $label rc=$rc tree=$tree"
	say "    per-package rows: $(grep -E '^(ok|FAIL|---)' "$LOGS/$label.body.log" 2>/dev/null | tr '\n' '|')"
	say "    setup-failed count=$(grep -c 'setup failed' "$LOGS/$label.body.log" 2>/dev/null || true) connection-refused count=$(grep -ci 'connection refused' "$LOGS/$label.body.log" 2>/dev/null || true)"
}

PKGS='./internal/risk ./internal/secret ./internal/config ./internal/agent ./internal/memory'
consrun cons2-base "$BASE" "$PKGS"
consrun cons2-fixed "$FIXED" "$PKGS"
# cmd/wisp is driven with a directed -run only: ticket 314 booked a whole-process death
# zone on that face and 165-r2 is live there in the working tree, so this runs inside the
# exported snapshot, which carries neither.
RUNF='-run "TestAC3POSIXSecretRouteLinkInsideItsDataRootStillRefused119|TestAC1POSIXSecretRouteSymlinkedHomeBecomesSealable119|TestAC1POSIXSecretRouteSymlinkedXDGConfigHomeBecomesSealable119" ./cmd/wisp'
consrun cons2-cmdwisp-base "$BASE" "$RUNF"
consrun cons2-cmdwisp-fixed "$FIXED" "$RUNF"
say "### consumer colour done $(date '+%F %T %z')"
