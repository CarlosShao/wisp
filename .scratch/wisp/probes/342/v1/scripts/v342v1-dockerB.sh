#!/usr/bin/env bash
# 342-v1 container driver, B variant (⛔ GOTMPDIR/TMPDIR bind mount).
# WHY THIS VARIANT EXISTS, as a finding not a shrug: my A variant mounted
# -v $S/gotmp:/gotmp and set GOTMPDIR=/gotmp for every container run. On that face the
# pristine baseline treeC read 71/41/2/0, with TWO reds (
# TestAC1POSIXSealHandleNameCarriesAnotherRealObject120r3 and
# TestPOSIX120SealWindowAncestorRelativeLinkIsClosed). Both are mode-bit/liveness tests,
# and a Docker-Desktop bind mount of a D: directory is a windowsdrv/9p share that does
# not store POSIX mode bits -- which is exactly the "filesystem accepts chmod and reports
# something else" case this package argues about (winsec_other.go:44-49). So that red is
# a BENCH defect of this leg, not a reading about ticket 342's window.
# The 342-a1 driver (scripts/driver-g3.sh, its own bytes) set only GOMODCACHE/GOCACHE and
# let TMPDIR default to the container's own /tmp on overlayfs. This B variant matches that
# config so colour faces are comparable, and the A/B pair above is the evidence.
# Usage: bash v342v1-dockerB.sh <tree> <sh-command> <log-name>
set -u
S=/d/tmp/wisp/20261018/342v1
IMG=golang:1.27
TREE="$1"
CMD="$2"
LOG="$3"
mkdir -p "$S/logs" "$S/gobuild"
MSYS_NO_PATHCONV=1 docker run --rm \
  -e CGO_ENABLED=0 -e GOPROXY=off -e GOFLAGS=-mod=mod \
  -v "$S/$TREE:/src" -v /d/work/base/gopath/pkg/mod:/gomod:ro \
  -v "$S/gobuild:/gobuild" -w /src "$IMG" \
  sh -c "GOMODCACHE=/gomod GOCACHE=/gobuild $CMD" \
  > "$S/logs/$LOG" 2>&1
echo "rc=$? tree=$TREE log=$LOG"
echo "RUN=$(grep -c '^=== RUN' "$S/logs/$LOG" || true) PASS_TOP=$(grep -c '^--- PASS' "$S/logs/$LOG" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$S/logs/$LOG" || true) SKIP=$(grep -c '^--- SKIP' "$S/logs/$LOG" || true)"
