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
#   1. go test ./...  - the seeded-violation positive control. It proves the
#      gate CAN go red, so a later "clean" cannot mean "the gate is blind".
#   2. go run . -root - the real scan of the working tree.
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

echo "d22scan.sh: positive control - go test ./... (tools/d22scan)"
go test ./...

echo "d22scan.sh: scan of $root"
go run . -root "$root"
