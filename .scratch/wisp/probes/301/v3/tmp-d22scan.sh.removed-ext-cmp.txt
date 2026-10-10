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
sh "$root/tools/d22scan/runtests.sh" -C tools/d22scan ./...

echo "d22scan.sh: scan of $root"
go run . -root "$root"
