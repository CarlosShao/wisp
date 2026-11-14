#!/bin/sh
# Ticket 444 AC#1/AC#2: make the seven race-window instruments run with the window NOT
# reached, so the sentence each one prints when the window is absent is a READING.
#
# WHERE THIS RUNS. The three POSIX instruments this script drives live in files carrying
# "//go:build !windows" / "//go:build unix", so on the Windows bench they are not compiled
# at all: a "did not reproduce here" reading from the host would be about nothing. They run
# INSIDE golang:1.27-bookworm (the image tickets 346 and 441 used). The three Windows
# instruments (ancestoropen_342_windows_test.go, ancestorunlink_348_windows_test.go,
# deleteseam_348r2_windows_test.go) are driven on the host with the same env knobs; the
# exact host commands are in docs/evidence/s1/444-instrument-abstention.md section 1.3.
#
# TWO KNOBS, BOTH OF THE RUN, NEITHER OF THE ASSERTIONS:
#   - WISP_444_ATTEMPTS=<n> shrinks an instrument's denominator (the committed default is
#     the instrument's own constant: 3000 / 1500 / 400). 1 is the smallest run that still
#     makes the route call once, and it is what AC#1's red sentences were read at.
#   - WISP_444_HOP_REFUSE=1 (deleteseam only) parks a NON-EMPTY directory where the seam's
#     flip lands, so rename(2) refuses every hop - the same branch the platform takes on CI
#     when the pinned handle turns the rename into "being used by another process".
#
# Env: INSTRUMENT=all|322|342u|348u, RUN_NAME (log suffix), COUNT, CPUS (-cpu).
# Usage (host):
#   docker run --rm -v "<repo>:/src" -w /src -e WISP_444_ATTEMPTS=1 \
#       golang:1.27-bookworm bash -c 'sh .scratch/wisp/probes/444/repro-444.sh'
set -u

run=${RUN_NAME:-r0}
count=${COUNT:-1}
cpus=${CPUS:-1}
what=${INSTRUMENT:-all}
attempts=${WISP_444_ATTEMPTS:-default}
out=${OUT:-/src/.scratch/wisp/probes/444/logs/repro-$what-$attempts-$run.log}

case "$what" in
322) sel='^TestAC1POSIXOpenThroughASwappedLeafLinkNeverTruncatesAForeignFile322r1$' ;;
342u) sel='^TestAC1POSIXPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1$' ;;
348u) sel='^TestAC2POSIXPrivateFileRaceThroughASwappedAncestorNeverDeletesAForeignFile342r1$' ;;
all) sel='^(TestAC1POSIXOpenThroughASwappedLeafLinkNeverTruncatesAForeignFile322r1|TestAC1POSIXPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1|TestAC2POSIXPrivateFileRaceThroughASwappedAncestorNeverDeletesAForeignFile342r1)$' ;;
*) echo "INSTRUMENT must be all, 322, 342u or 348u" >&2; exit 2 ;;
esac

/usr/local/go/bin/go test ./internal/winsec -run "$sel" \
    -count="$count" -cpu="$cpus" -v >"$out" 2>&1
rc=$?

echo "instrument=$what attempts=$attempts count=$count cpus=$cpus rc=$rc log=$out"
echo "abstain-lines:        $(grep -c 'ABSTAIN(444:' "$out")"
echo "controls-printed:     $(grep -c 'un-raced .*control' "$out")"
echo "red-arrival:          $(grep -c 'RED instrument: the racer' "$out")"
echo "red-guardrail:        $(grep -c 'did not keep this route.s act inside the pinned object\|did not get the refusal ticket 322\|neither kept the foreign name' "$out")"
echo "red-damage:           $(grep -c 'AC#1 RED:\|AC#2 RED:' "$out")"
echo "cases-pass:           $(grep -c '^--- PASS:' "$out")"
echo "cases-fail:           $(grep -c '^--- FAIL:' "$out")"
