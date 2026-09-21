#!/bin/sh
# tools/d22scan/runtests.sh - a `go test` runner that refuses to read "no tests
# ran" as a pass (ticket 71 AC#3).
#
# WHY THIS EXISTS (measured on HEAD 2026-09-21, not theorized):
#
#	$ go test ./internal/proc/ -run TestNoSuchThingExists -count=1 -v
#	testing: warning: no tests to run
#	PASS
#	ok  	github.com/CarlosShao/wisp/internal/proc	0.033s [no tests to run]
#	exit 0
#
# Go's own semantics: a -run pattern matching nothing still prints PASS and
# exits 0, and a skipped test prints `ok` too. Two of CI's steps are exactly
# this shape (`-run TestPathResolverJunctionWindows`, `-run TestLayoutForTestEnv`),
# so renaming either test turns the step into a permanent green no-op - the same
# "empty instrument" disease tools/d22scan got a self-report for, in the test
# layer instead of the scanner layer. We cannot change Go; we can stop calling
# bare `go test` a gate.
#
# Rules this enforces, all fatal:
#   1. zero top-level `--- PASS`/`--- FAIL` lines  -> the pattern matched nothing
#      or no test binary ran at all.
#   2. any `--- SKIP`                             -> SKIP is NOT a pass. A test
#      that legitimately cannot run on this platform must be moved out of the
#      step's scope, not silenced here (that would be lowering a threshold).
#   3. a non-zero `go test` exit code             -> propagated unchanged.
#
# Usage:
#	sh tools/d22scan/runtests.sh [go test flags and packages...]
#	sh tools/d22scan/runtests.sh -C tools/d22scan ./...     # run inside another module
#
# -C <dir> is relative to the repository root this script derives from its own
# location (ticket 67 AC#2's lesson: tools/* are separate Go modules, and a
# `go run ./tools/...` from the root silently scans nothing - so the directory
# is stated once here and echoed in every verdict line, never guessed).
#
# No `|| true`, no skippable step, output is echoed verbatim before the asserts.
set -eu

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
root=$(CDPATH= cd -- "$here/../.." && pwd)

cwd="$root"
if [ "${1:-}" = "-C" ]; then
    [ $# -ge 3 ] || { echo "runtests.sh: -C needs a directory and then the go test args" >&2; exit 2; }
    cwd=$2
    shift 2
    if [ "${1:-}" = "-C" ]; then
        echo "runtests.sh: -C given twice - one module per invocation, on purpose" >&2
        exit 2
    fi
    case $cwd in
        /*|*:*) ;; # already absolute (posix path or windows drive letter)
        *) cwd="$root/$cwd" ;;
    esac
fi
if [ ! -f "$cwd/go.mod" ]; then
    echo "runtests.sh: no go.mod in $cwd - refusing to run tests against a directory" \
        "that is not a module root (a wrong -C is how a green no-op is born)" >&2
    exit 2
fi
if [ $# -eq 0 ]; then
    echo "runtests.sh: no packages given; pass e.g. ./internal/proc/ -run TestX" >&2
    exit 2
fi

out=$(mktemp 2>/dev/null || echo "$cwd/.runtests.$$.log")
cd "$cwd"

# -v is forced: without it Go prints no per-test lines at all and the counts
# below would be vacuous. -count=1 is forced for the same reason a cached result
# must never satisfy a gate.
set +e
go test -v -count=1 "$@" >"$out" 2>&1
rc=$?
set -e
cat "$out"

count() {
    # grep exits 1 on no match; that is a real zero here, not an error.
    n=$(grep -c "$1" "$out" || true)
    echo "${n:-0}"
}
ran=$(count '^=== RUN')
passed=$(count '^--- PASS')
failed=$(count '^--- FAIL')
skipped=$(count '^--- SKIP')
noresult=$(count '\[no tests to run\]')

rm -f "$out"

summary="packages=[$*] top-level: PASS=$passed FAIL=$failed SKIP=$skipped, === RUN=$ran, '[no tests to run]'=$noresult"
if [ "$rc" -ne 0 ]; then
    echo "runtests.sh: go test exited $rc - $summary" >&2
    exit "$rc"
fi
if [ "$skipped" -ne 0 ]; then
    echo "runtests.sh: $skipped test(s) SKIPPED and SKIP is not a pass (ticket 71 AC#3) - $summary" >&2
    echo "runtests.sh: either make the step's scope exclude it or run it where it can execute;" >&2
    echo "             do not relax this script." >&2
    exit 1
fi
if [ "$passed" -eq 0 ] && [ "$failed" -eq 0 ]; then
    echo "runtests.sh: go test exited 0 but reported NO top-level result at all (ticket 71 AC#3) - $summary" >&2
    echo "runtests.sh: this is what a -run pattern that matches nothing looks like. Name the tests you" >&2
    echo "             intended to run, or delete the step - a green no-op is worse than a red." >&2
    exit 1
fi
echo "runtests.sh: OK - $summary"
