#!/bin/sh
# Ticket 314 shape-ruler readings. Runs the staged carrier from the repo root;
# every reading in docs/evidence/s1/314-shape-ruler.md was produced by this file.
# Usage: run-ruler.sh -test.run '^Test314'
set -u
cd "$(dirname "$0")/../../../.." || exit 2
echo "pwd=$PWD"
exec sh scripts/wisp-cli-staged-tests.sh "$@"
