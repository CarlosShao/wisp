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

# ---- TICKET 310 AC#1 (shape A-G): the test temp family gets a dedicated root.
# This is the fifth entry that really starts `go test`: scripts/d22scan.sh:51
# calls `sh tools/d22scan/runtests.sh -C tools/d22scan ./...`, CI's "D22 scanner
# positive control" step calls the same runner, and the `go test` below is the
# real process start. The block is copied verbatim from scripts/winsec-tests.sh
# (only the entry name and the echo prefix change) so all four entries that start
# tests present ONE shape, not four dialects. The guard keeps the delegation
# contract: when scripts/portable-tests.sh has already exported GOTMPDIR and then
# calls this runner, we MUST NOT open a second root - one root per run, named by
# the outermost entry; the delegated run inherits it. Why GOTMPDIR and not
# TMP/TEMP, why the root is a SIBLING of %TEMP% (never inside it, where a
# non-recursive roster would read "0 new" while the residue is simply out of the
# ruler's reach), why it is not on the repo's volume (a cross-volume root turns
# the winsec seal cases red - measured, ticket 311), why nothing is deleted here,
# and why the ubuntu leg is left alone: the long note in scripts/portable-tests.sh.
#
# The GOOS=windows + empty-LOCALAPPDATA branch prints one honest line instead of
# going silent (ticket 310 AC#1, debt Q-310-orch-Q): without it a caller with no
# %LOCALAPPDATA% falls through to the OS temp root and NOTHING says so, which is
# the shape 310-v1 registered (Q-310-v1-5). It is a finding printed to stderr,
# NOT a fail-closed exit and NOT a refusal to run the tests (orchestrator ruling:
# keep the loud-degradation shape; revocation phrase "310 改 fail-closed").
if [ -z "${GOTMPDIR:-}" ] && [ "$(go env GOOS)" = windows ]; then
    if [ -n "${LOCALAPPDATA:-}" ]; then
        gotmp_run="${WISP_TEST_GOTMP_ROOT:-"${LOCALAPPDATA//\\//}/wisp-test-gotmp"}/d22scan-$(date +%Y%m%d-%H%M%S)-$$"
        if mkdir -p "$gotmp_run" 2>/dev/null; then
            export GOTMPDIR="$gotmp_run"
            echo "runtests.sh: ticket 310 AC#1: GOTMPDIR=$GOTMPDIR - this run's" \
                "t.TempDir() family lands here; nothing deletes it."
        else
            echo "runtests.sh: ticket 310 AC#1: DEGRADED - cannot create $gotmp_run," \
                "so GOTMPDIR stays unset and a test binary that dies leaves its temp dir in" \
                "the OS temp root again. Not a silent fallback: this line is the finding." >&2
        fi
    else
        echo "runtests.sh: ticket 310 AC#1: DEGRADED - LOCALAPPDATA is empty, so this step has" \
            "no measured-safe root to build (a root on another volume turns the winsec seal" \
            "cases red, see the note in scripts/portable-tests.sh). GOTMPDIR stays unset." >&2
    fi
fi

# ---- TICKET 311 AC#1 (shape B, the orchestrator's ruling "B + C"; shape A was
# refused): the VOLUME this run's test temp root sits on is part of this gate's
# verdict, so a root that names a volume other than %LOCALAPPDATA%'s is REFUSED HERE,
# on the exit code, instead of reaching the log as a gate that quietly changed colour.
# Measured, not theorized (ticket 311 census, reproduced by leg 311-a1 and by the
# orchestrator the same day): ONE unchanged internal/winsec tree reads
# PASS=58 FAIL=0 rc=0 from a %LOCALAPPDATA% (C:) root and PASS=56 FAIL=2 rc=1 from a
# D:/tmp root, and those two reds are Go's own t.TempDir cleanup failing
# (testing.go:1617 "TempDir RemoveAll cleanup: ... Access is denied"), not an
# assertion - which is why this is a refusal here rather than a fixture change there.
# This block adds NO assertion, touches no threshold, golden or file under
# internal/winsec, and none may be touched to make it green (ticket 311 forbidden
# list; D22). Shape A is refused by that same reading: what moved between the two
# shots was the cleanup, not the gate's judgement, so no fixture changes here.
#
# PLACEMENT IS LOAD-BEARING, AND MEASURED: the ticket-310 root block of this entry
# runs BEFORE its guard chain (that is where the DEGRADED line comes from, measured
# by the orchestrator at A862 and re-measured per entry by leg 311-r1b), so this guard
# sits AFTER the whole guard chain above and BEFORE the run below, and one trip reads,
# in order: [310 DEGRADED/GOTMPDIR line] -> [311 GUARD line] -> [no test process at
# all]. Ordering readings: .scratch/wisp/probes/311/r1b/logs/23-order-degraded-then-guard.txt
# and .scratch/wisp/probes/311/r1b/logs/24-order-guard-after-guard1.txt - the two files
# this comment used to name sat under probes/311/r1/ and never landed (leg 311-r1 died).
#
# THREE SHAPES THIS MUST NOT BREAK, each with its own reading on the ticket:
# WISP_TEST_GOTMP_ROOT unset (the factory default), a root on the SAME volume under a
# different directory, and anything where `go env GOOS` is not windows or
# %LOCALAPPDATA% is empty - that last branch is the one ticket 310 made LOUD, and this
# guard stays silent there too (no line, no exit code). Everything else about this
# branch is deliberately NOT that precedent: it fires on the exit code on purpose
# (ruling B; revocation phrase "311 撤乙"). A value's volume is its drive letter, read
# in both the Windows form (C:/x and C:\x) and the form this shell prints (/c/x); a
# value neither form can judge (UNC, bare relative) prints an INCONCLUSIVE line and
# leaves the exit code alone - undecided shape, registered as debt Q-311-r1-1. What
# this guard cannot reach at all is a caller that skips these entries and runs
# `go test` directly (leg 311-a1 did exactly that), which is why shape C stayed in
# force beside it (.scratch/wisp/probes/orch/gate-roster-304.md section 1).
if [ "$(go env GOOS)" = windows ] && [ -n "${LOCALAPPDATA:-}" ]; then
    wisp_311_volume_of() {
        wisp_311_arg=${1//\\//}
        case $wisp_311_arg in
            [A-Za-z]:*) printf '%.1s' "$wisp_311_arg" | tr 'a-z' 'A-Z' ;;
            /[A-Za-z]/*) printf '%.1s' "${wisp_311_arg#/}" | tr 'a-z' 'A-Z' ;;
            *) printf '' ;;
        esac
    }
    wisp_311_ref_path=${LOCALAPPDATA//\\//}
    wisp_311_ref_vol=$(wisp_311_volume_of "$wisp_311_ref_path")
    wisp_311_cross=''
    for wisp_311_probe in "GOTMPDIR=${GOTMPDIR:-}" "WISP_TEST_GOTMP_ROOT=${WISP_TEST_GOTMP_ROOT:-}"; do
        wisp_311_name=${wisp_311_probe%%=*}
        wisp_311_value=${wisp_311_probe#*=}
        if [ -z "$wisp_311_value" ]; then
            continue
        fi
        wisp_311_vol=$(wisp_311_volume_of "$wisp_311_value")
        if [ -z "$wisp_311_vol" ] || [ -z "$wisp_311_ref_vol" ]; then
            echo "runtests.sh: ticket 311 AC#1 shape B: INCONCLUSIVE - no volume can be read from $wisp_311_name=$wisp_311_value (reference %LOCALAPPDATA%=$wisp_311_ref_path), so nothing is judged here and this step runs on the root as given. This is NOT the check passing: registered as debt Q-311-r1-1 on ticket 311." >&2
            continue
        fi
        if [ "$wisp_311_vol" != "$wisp_311_ref_vol" ]; then
            wisp_311_cross="$wisp_311_cross $wisp_311_name=$wisp_311_value(volume $wisp_311_vol)"
        fi
    done
    if [ -n "$wisp_311_cross" ]; then
        echo "runtests.sh: ticket 311 AC#1 shape B: GUARD - the test temp root names ANOTHER VOLUME than %LOCALAPPDATA%: reference volume $wisp_311_ref_vol (path $wisp_311_ref_path), offender:$wisp_311_cross" >&2
        echo "runtests.sh: ticket 311 AC#1 shape B: this gate's COLOUR is a function of that volume - the same unchanged internal/winsec code reads PASS=58 FAIL=0 on a volume-$wisp_311_ref_vol root and PASS=56 FAIL=2 with Go's own TempDir cleanup failing on a D:/tmp root (ticket 311). Refusing to run is the fix; reporting a re-coloured gate as a verdict is not." >&2
        echo "runtests.sh: ticket 311 AC#1 shape B: point WISP_TEST_GOTMP_ROOT (or GOTMPDIR) at a directory on volume $wisp_311_ref_vol, or unset both, then run this step again. Do NOT relax an assertion, threshold, golden or SLO to make it green." >&2
        exit 2
    fi
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
