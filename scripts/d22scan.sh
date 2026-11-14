#!/bin/sh
# scripts/d22scan.sh - the only supported way to run the D22 gate (ticket 67 AC#2).
#
# Why a wrapper exists: tools/d22scan is its OWN Go module, so from the repo
# root `go run ./tools/d22scan -root .` fails inside the ROOT module and the
# scanner never starts - and `cd tools/d22scan && go run .` (default -root .)
# walks a tree with no internal/ or cmd/, so every ban examines zero files.
# Both forms used to look like success. The scanner now refuses to print a
# verdict it did not compute (see checkRoot in tools/d22scan/main.go), and this
# script removes the path choices that made the mistake possible: it derives
# the repository root from its own location, so it works from any directory.
#
# Both commands are load-bearing and run in this order on purpose:
#   1. runtests.sh -C tools/d22scan ./... - the seeded-violation positive
#      control. It proves the gate CAN go red, so a later "clean" cannot mean
#      "the gate is blind".
#   2. go run . -root - the real scan of the working tree.
# Since ticket 71 AC#5 they BOTH run even when step 1 fails: the order is still
# control-first (a clean control is what licenses trusting step 2), but step 2's
# verdict is no longer allowed to go missing because step 1 was red - see the
# WARNING block below, which also records the measurement that forced it.
#
# Why step 1 goes through tools/d22scan/runtests.sh and not bare `go test`
# (ticket 99 AC#2, measured in a /tmp snapshot of HEAD, not theorized):
# several tests here - TestScannerSelfScanOfRealRepoIsGreen above all - read the
# repository tree at RUNTIME, and the go test cache only records build inputs,
# so a violation added to internal/ or frontend/ cannot invalidate them. Run 1
# of the script passed, a `frontend/src/app.js:1 approval.decide` was then
# planted, run 2 of step 1 printed:
#
#	ok  	github.com/CarlosShao/wisp/tools/d22scan	(cached)
#
# with exit code 0 - the positive control endorsed a tree it had never looked
# at again. Step 2 still caught it, so the script as a whole stayed red, but a
# control that has to be trusted cannot be the one served from a cache.
# runtests.sh forces -count=1 (and -v, and "SKIP is not a pass"), which is the
# SAME instrument CI's "D22 scanner positive control" step already calls, so the
# two call sites no longer hold two different rules for one job (ticket 71
# AC#3's lesson: bare `go test` is not a gate).
#
# No skippable step, no `|| true`, no continue-on-error (D22 run-away mode 6).
set -eu

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
root=$(CDPATH= cd -- "$here/.." && pwd)

if [ ! -f "$root/go.mod" ]; then
    echo "d22scan.sh: no go.mod at derived repo root $root - refusing to scan" >&2
    exit 2
fi

cd "$root/tools/d22scan"

echo "d22scan.sh: positive control - runtests.sh -C tools/d22scan ./..."
rc_control=0
sh "$root/tools/d22scan/runtests.sh" -C tools/d22scan ./... || rc_control=$?

# WHY THE SCAN STILL RUNS WHEN THE CONTROL WENT RED (ticket 71 AC#5, measured
# 2026-11-14 in an out-of-repo export of HEAD, not theorized): under plain
# `set -e` the first red ENDED this script, so a tree that actually carried a
# violation produced a CI step whose entire output was the control's failure -
# the scan never ran, and the one line that locates the offending file never
# existed. That is the same disease this ticket exists for ("the step after the
# first failure is skipped", D22 run-away mode 6), one level down: the wrapper
# was hiding its own gate's verdict behind its own guard. Running both and
# reporting both removes the hole without loosening anything - there is still no
# skippable step, no `|| true`, no continue-on-error, and the exit status below
# is the FIRST non-zero of the two, so no red can become green and the control's
# failure keeps exactly the exit code CI saw before. The only added information
# is the scan's verdict, which is labelled untrusted when the control is red
# because a broken instrument can also scan wrongly.
if [ "$rc_control" -ne 0 ]; then
    echo "d22scan.sh: WARNING - the positive control failed (rc=$rc_control); the scan below runs anyway so its verdict is not missing from this log, but read it as untrusted rather than as absent" >&2
fi

echo "d22scan.sh: scan of $root"
rc_scan=0
go run . -root "$root" || rc_scan=$?

if [ "$rc_control" -ne 0 ]; then
    echo "d22scan.sh: FAIL - control rc=$rc_control, scan rc=$rc_scan" >&2
    exit "$rc_control"
fi
if [ "$rc_scan" -ne 0 ]; then
    echo "d22scan.sh: FAIL - control rc=0, scan rc=$rc_scan" >&2
    exit "$rc_scan"
fi
echo "d22scan.sh: OK - control rc=0, scan rc=0"
