#!/usr/bin/env bash
# 336-v1 mandatory-(iii): the POSIX half of ticket 336 AC#2, re-measured BY THE JUDGE on
# Linux (docker golang:1.27 / go1.27.1 linux/amd64) instead of taken from the leg.
# Mounts the judge's own exported tree (named anchor e5021339) and runs the risk package
# roster. Expectation being tested: the !windows tripwire is GREEN, the windows-tagged
# teeth do NOT run on this platform (build tag), and nothing pretends to be a tooth here.
# NOTE: -w with a Windows path fails on this daemon ("working directory 'W:/' is
# invalid"), so the container cds into /w itself. First attempt kept as
# 24-posix-linux-run.raw (docker_rc=125, a platform/CLI issue, NOT a judgment reading).
set -u
REPO="D:/work/workspace/projects plans/Wisp"
D="$REPO/.scratch/wisp/probes/336/v1/logs"
TREEWIN="D:/tmp/wisp/20261018/336v1/tree-ac2"
{
echo "### judge Linux run: tree=$TREEWIN anchor=e5021339"
echo "### go version: $(docker run --rm -v "$TREEWIN:/w" golang:1.27 sh -c 'cd /w && go version' 2>&1 | tail -1)"
# ★First attempt failed on the台面, not on the judgment: passing -e GOCACHE=/tmp/gocache
# through Git Bash made MSYS rewrite the value to C:/Users/.../Temp/gocache, and the Go
# command then died with "GOMODCACHE entry is relative; must be absolute path" (rc=1, zero
# RUN lines = nothing executed). Kept as 24-posix-linux-run.log as the failure record.
# Fix = set the cache dirs INSIDE the container command so no path passes through MSYS.
docker run --rm -v "$TREEWIN:/w" golang:1.27 \
  sh -c 'export GOCACHE=/tmp/gocache GOMODCACHE=/tmp/gomod GOFLAGS=-mod=mod; cd /w && go env GOCACHE GOMODCACHE && go test ./internal/risk/ -run "TestFallbackBranchHasNoWorldOnThisPlatform336|TestAnchorFallbackTeethPerAnchor336|TestAnchorMissReasonNamesCandidateAndMissingAnchor|TestAnchorMissKeepsNonOverridableClassAndNeedL2|TestFileBaseNameClaimChangesNoClassOrReason" -count=1 -timeout 30m -v 2>&1; echo "rc_test=$?"' > "$D/24b-posix-linux-run.log" 2>&1
echo "  docker_rc=$?"
echo "### four counts on the Linux leg (RUN/PASS_TOP/FAIL_TOP/SKIP):"
printf "  RUN=%s PASS_TOP=%s FAIL_TOP=%s SKIP=%s rc=%s\n" \
  "$(grep -c '^=== RUN' "$D/24b-posix-linux-run.log" || true)" \
  "$(grep -c '^--- PASS' "$D/24b-posix-linux-run.log" || true)" \
  "$(grep -c '^--- FAIL' "$D/24b-posix-linux-run.log" || true)" \
  "$(grep -c '^--- SKIP' "$D/24b-posix-linux-run.log" || true)" \
  "$(grep -oE 'rc_test=[0-9]+' "$D/24b-posix-linux-run.log" | tail -1)"
echo "### top-level result lines:"
grep -E '^(--- |ok |FAIL|PASS)' "$D/24b-posix-linux-run.log" | sed 's/^/    /' || true
echo "### tripwire's own readouts (cand_root / verdict per placement):"
grep -E 'TEETH336-OTHER' "$D/24b-posix-linux-run.log" | sed 's/^/    /' || true
echo "### positive control that the roster really executed on Linux:"
printf "  RUN_lines_for_tripwire=%s (want >=1)  teeth_windows_lines=%s (want 0: build-tagged windows)\n" \
  "$(grep -c 'RUN   TestFallbackBranchHasNoWorldOnThisPlatform336' "$D/24b-posix-linux-run.log" || true)" \
  "$(grep -c 'TestAnchorFallbackTeethPerAnchor336' "$D/24b-posix-linux-run.log" || true)"
echo "### any compile/setup failure named (want none):"
grep -iE 'build failed|setup failed|no test files|undefined|0xc0000135|relative' "$D/24b-posix-linux-run.log" | sed 's/^/    /' || true
echo "posix_rc=0"
} 2>&1 | tee "$D/24-posix-rerun.log"
