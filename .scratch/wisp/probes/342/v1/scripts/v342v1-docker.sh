#!/usr/bin/env bash
# 342-v1 container driver (POSIX leg). ⛔ shared working tree. golang:1.27 with
# CGO_ENABLED=0, GOPROXY=off, host module cache mounted read-only, build cache on D:.
# Usage: bash v342v1-docker.sh <tree-dir-basename> <sh-command-inside> <log-name>
set -u
S=/d/tmp/wisp/20261018/342v1
IMG=golang:1.27
TREE="$1"
CMD="$2"
LOG="$3"
mkdir -p "$S/gobuild" "$S/gotmp"
MSYS_NO_PATHCONV=1 docker run --rm \
  -e CGO_ENABLED=0 -e GOPROXY=off -e GOFLAGS=-mod=mod \
  -v "$S/$TREE:/src" -v /d/work/base/gopath/pkg/mod:/gomod:ro \
  -v "$S/gobuild:/gobuild" -v "$S/gotmp:/gotmp" -w /src "$IMG" \
  sh -c "GOMODCACHE=/gomod GOCACHE=/gobuild GOTMPDIR=/gotmp TMPDIR=/gotmp TMP=/gotmp TEMP=/gotmp $CMD" \
  > "$S/logs/$LOG" 2>&1
echo "rc=$? tree=$TREE log=$LOG"

four () {
  echo "RUN=$(grep -c '^=== RUN' "$1" || true) PASS_TOP=$(grep -c '^--- PASS' "$1" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$1" || true) SKIP=$(grep -c '^--- SKIP' "$1" || true)"
}
echo "-- four $LOG:"; four "$S/logs/$LOG"
