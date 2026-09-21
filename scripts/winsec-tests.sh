#!/usr/bin/env bash
# scripts/winsec-tests.sh - the CI step that runs internal/winsec's OWN tests
# (ticket 110 AC#2).
#
# WHY A SEPARATE STEP EXISTS (measured, not theorized - ticket 110 AC#1):
# in run 35591482293 the windows job's test steps were step 6 (portable-tests.sh
# with ./internal/proc/ ./internal/secret/ ./internal/config/) and step 7
# (runtests.sh -run TestPathResolverJunctionWindows on ./internal/risk/). Across
# the whole run, the string `ok github.com/CarlosShao/wisp/internal/winsec`
# appeared ZERO times: the package that decides whether a path gets a real ACL
# was never itself tested on CI. Tickets 94/103/106 changed that code and their
# CI "evidence" was downstream (internal/secret calls into it) - indirect cover
# for the one package where the failure mode is the machine's own DACL, which is
# exactly the thing a developer laptop cannot reproduce.
#
# WHY IT IS NOT just one more path in the existing windows step's list:
# AC#3 asks the gate to prove it is not an empty instrument. If winsec were
# silently appended to a list three floors up in ci.yml, deleting it from that
# list would still be green - the same "配置里有一行 = 它跑过了" disease ticket
# 110 exists to close (registry A54, A61). So the scope is stated HERE, once,
# and guard 1 below makes removing it a red step rather than a shorter run.
#
# THE INSTRUMENT IS NOT NEWER OR WEAKER THAN THE EXISTING ONE: every test
# verdict is delegated to scripts/portable-tests.sh, which delegates to
# tools/d22scan/runtests.sh (ticket 71 AC#3: a non-zero `go test` exit propagates
# unchanged, any top-level SKIP is fatal, and zero PASS *and* zero FAIL is fatal).
# This script adds two guards and subtracts nothing:
#   guard 1 - the package list this step was given must name ./internal/winsec/;
#   guard 2 - the run must actually print a top-level result line for the
#             winsec package itself (a build tag that ever excludes the whole
#             package would otherwise read as green while testing nothing).
set -eu -o pipefail

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
root=$(CDPATH= cd -- "$here/.." && pwd)
cd "$root"

target=./internal/winsec/
portable="$root/scripts/portable-tests.sh"
if [ ! -f "$portable" ]; then
    echo "winsec-tests.sh: $portable is missing - this script is a scope-narrowing" \
        "shell around it and refuses to become the weaker instrument it replaces." >&2
    exit 2
fi

# The scope this step exists for, stated once. Extra packages may be appended by
# hand while debugging; the list can never lose winsec without going red.
scope=()
if [ $# -eq 0 ]; then
    scope=("$target")
else
    scope=("$@")
fi

# guard 1: an explicit failure beats a shorter run.
named=0
for p in "${scope[@]}"; do
    [ "$p" = "$target" ] && named=1
done
if [ "$named" -ne 1 ]; then
    echo "winsec-tests.sh: GUARD 1 - the package list for this step does not name" \
        "$target; it is [${scope[*]}]." >&2
    echo "winsec-tests.sh: this step IS the internal/winsec gate (ticket 110 AC#2/AC#3)." \
        "Removing the package from the list must fail the step, not shorten it." >&2
    exit 2
fi

goos=$(go env GOOS)
if [ -z "$goos" ]; then
    echo "winsec-tests.sh: go env GOOS came back empty" >&2
    exit 2
fi
if [ "$goos" != windows ]; then
    # The sealing assertions are //go:build windows files. Running this step on
    # another platform would measure a different (smaller) denominator while
    # printing the same green, so it is a mis-wired step, not a pass.
    echo "winsec-tests.sh: GUARD - this step is the windows leg of the winsec gate," \
        "GOOS=$goos. Fix the job it is wired into (ticket 110 AC#5); do not relax this line." >&2
    exit 2
fi

pkgpath=$(go list "$target")
if [ -z "$pkgpath" ]; then
    echo "winsec-tests.sh: go list $target came back empty" >&2
    exit 2
fi

echo "winsec-tests.sh: platform=$goos gate=internal/winsec scope=[${scope[*]}] package=$pkgpath"

capture=$(mktemp 2>/dev/null || echo "$root/.winsec-tests.$$.log")
rc=0
# tee keeps portable-tests.sh's own four numbers verbatim while giving this
# script the bytes guard 2 needs.
bash "$portable" "${scope[@]}" 2>&1 | tee "$capture" || rc=$?

count() { grep -c "$1" "$capture" || true; }
ran=$(count '^=== RUN')
passed=$(count '^--- PASS')
failed=$(count '^--- FAIL')
skipped=$(count '^--- SKIP')

# guard 2: the package itself must have booked a top-level result.
result_line=$(grep -E "^(ok|FAIL)[[:space:]]+$pkgpath([[:space:]]|$)" "$capture" || true)
if [ -z "$result_line" ]; then
    echo "winsec-tests.sh: GUARD 2 - the run printed NO top-level result line for" \
        "$pkgpath, so this step tested the package it is named for zero times." >&2
    echo "winsec-tests.sh: a build tag or a rename that empties the package is a" \
        "red step here (ticket 110 AC#3), not a green short-circuit." >&2
    rc=1
    result_line='<none>'
fi

echo "winsec-tests.sh: winsec result line: $result_line"
echo "winsec-tests.sh: four numbers (all from -v output): === RUN=$ran  --- PASS=$passed  --- FAIL=$failed  --- SKIP=$skipped"
rm -f "$capture"
[ "$rc" -eq 0 ] || echo "winsec-tests.sh: gate is RED (rc=$rc) for scope=[${scope[*]}]" >&2
exit "$rc"
