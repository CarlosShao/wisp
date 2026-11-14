#!/bin/sh
# Ticket 441 AC#1: reproduce the leaf-swap window instrument's abstention under
# CONTROLLED contention, so the red sentence it prints is a reading and not a guess.
#
# This script runs INSIDE the Linux container (golang:1.27-bookworm, the same image
# ticket 346 used) because posix_seal_window_120_other_test.go carries
# "//go:build !windows": on the Windows bench it is not compiled at all, so a
# "did not reproduce here" reading from the host would be about nothing.
#
# Contention knob = -cpu (GOMAXPROCS). The racer and the seal loop are two
# goroutines; pinning them to one P is the controlled form of "a loaded CI runner
# schedules the seal loop in the gaps", and it is a knob of the RUN, not of the
# assertion: nothing in the test file is changed by this script.
#
# Env: RUN_NAME (log suffix), COUNT (repetitions of the whole test), CPUS (-cpu).
# Usage (host): docker run --rm -v "<repo>:/src" -w /src golang:1.27-bookworm \
#                   bash -c 'sh .scratch/wisp/probes/441/repro-leafswap.sh'
set -u

run=${RUN_NAME:-r0}
count=${COUNT:-10}
cpus=${CPUS:-1}
out=${OUT:-/src/.scratch/wisp/probes/441/logs/repro-$run.log}

/usr/local/go/bin/go test ./internal/winsec \
    -run '^TestPOSIX120SealWindowLeafSwapIsClosed$' \
    -count="$count" -cpu="$cpus" -v >"$out" 2>&1
rc=$?

echo "run=$run count=$count cpus=$cpus rc=$rc log=$out"
echo "samples-with-window: $(grep -c 'AC#1 sample .* shape=leaf-swap:' "$out")"
echo "red-instrument-dead: $(grep -c 'AC#1 RED instrument is dead' "$out")"
echo "red-sample-open:     $(grep -c 'AC#1 RED sample' "$out")"
echo "test-cases-pass:     $(grep -c '^--- PASS: TestPOSIX120SealWindowLeafSwapIsClosed' "$out")"
echo "test-cases-fail:     $(grep -c '^--- FAIL: TestPOSIX120SealWindowLeafSwapIsClosed' "$out")"
