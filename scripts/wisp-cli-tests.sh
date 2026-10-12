#!/usr/bin/env bash
# scripts/wisp-cli-tests.sh - the CI step that runs cmd/wisp's OWN tests on the
# windows leg (ticket 111 AC#4).
#
# AC#4 ASKED FOR A MEASURED ANSWER TO "can cmd/wisp be tested on CI at all", and
# explicitly forbade answering it by reasoning from the laptop. Ticket 98's hole
# is that `go test ./cmd/wisp/` on this host fails at LOAD time with
# STATUS_DLL_NOT_FOUND (0xc0000135): the test binary links against the sherpa-onnx
# cgo import library, so Windows refuses to start the process before Go's first
# line runs, and the failure looks like a broken package rather than a missing
# DLL. Four readings, all run this ticket, are what the wiring below is built on:
#
#   windows, no PATH help .......... the negative control, re-measured at this
#                                    HEAD: `go test -c` builds fine (rc=0), and
#                                    starting the binary gives rc=127 and
#                                    "error while loading shared libraries:
#                                    sherpa-onnx-c-api.dll: cannot open shared
#                                    object file" - 0 stdout bytes, so not one
#                                    test ever ran
#   windows, PATH=third_party/sherpa-onnx ... runtests.sh: PASS=33 FAIL=0 SKIP=0,
#                                    === RUN=67, rc=0, 51.1s  -> GREEN
#   ubuntu,   CGO_ENABLED=0 ........... cannot even resolve:
#                                    "build constraints exclude all Go files in
#                                    sherpa-onnx-go-linux" (vet rc=1)
#   ubuntu,   CGO_ENABLED=1 ........... the binary BUILDS and RUNS (no load-time
#                                    failure, so the linux .so layout is fine) but
#                                    19 of 29 top-level cases FAIL, rc=1, 127.4s
#
# So the answer is: on CI it CAN run, on exactly one of the two legs. cmd/wisp is
# wired into test-windows, where the third_party DLLs are already staged by the
# cache + build steps of that same job. The ubuntu reading is NOT silently
# dropped: it is registered as a finding on the ticket (19 reds in a package that
# has never been run on either platform before this ticket, so nobody knows yet
# how much of it is a real cross-platform bug and how much is a hosted-runner
# environment). Trimming the linux wiring until those 19 are explained would be
# the "调阈值变绿" move AC#2 forbids, so the split is by measured capability
# (cgo/shared-library layout), not by convenience.
#
# THE ADDITION OVER portable-tests.sh IS THE DIAGNOSIS, NOT THE VERDICT: every
# test claim still comes from portable-tests.sh -> tools/d22scan/runtests.sh
# (non-zero go test exit propagates, any top-level SKIP is fatal, zero PASS *and*
# zero FAIL is fatal), which adds GUARD A (the scope must have a compiled test
# file) and GUARD B (cmd/wisp must print its own anchored top-level result line).
# What this script owns is the one failure mode that is specific to it: a missing
# DLL is a setup problem, and a setup problem must say so instead of arriving at
# the CI log as an opaque exit status 0xc0000135.
set -eu -o pipefail

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
root=$(CDPATH= cd -- "$here/.." && pwd)
cd "$root"

portable="$root/scripts/portable-tests.sh"
if [ ! -f "$portable" ]; then
    echo "wisp-cli-tests.sh: $portable is missing - this script is a PATH-staging shell" \
        "around it and refuses to become the weaker instrument it replaces." >&2
    exit 2
fi

# ---- TICKET 310 AC#1 (shape A-G): claim the dedicated GOTMPDIR root at THIS
# entry, before the PATH-staging guard below can short-circuit, because cmd/wisp is
# the package whose test binary is the heaviest temp-dir user this host has. The
# delegated portable-tests.sh builds a root only while GOTMPDIR is still empty, so
# it reuses this one instead of opening a second root. Rationale (GOTMPDIR vs
# TMP/TEMP, why the root is outside the OS temp dir, why nothing is deleted here,
# why non-Windows is untouched) is the long note in scripts/portable-tests.sh.
#
# The GOOS=windows + empty-LOCALAPPDATA branch is now LOUD (ticket 310 AC#1, debt
# Q-310-orch-Q). It used to be the conjunctive form `... && [ -n LOCALAPPDATA ]`,
# so when %LOCALAPPDATA% was empty the whole block fell through in silence - no
# GOTMPDIR line, no DEGRADED line, nothing said the redirection did not happen. It
# is now a finding on stderr, still NOT a fail-closed exit and NOT a refusal to
# run (orchestrator ruling; revocation phrase "310 改 fail-closed").
if [ -z "${GOTMPDIR:-}" ] && [ "$(go env GOOS)" = windows ]; then
    if [ -n "${LOCALAPPDATA:-}" ]; then
        gotmp_run="${WISP_TEST_GOTMP_ROOT:-"${LOCALAPPDATA//\\//}/wisp-test-gotmp"}/cli-$(date +%Y%m%d-%H%M%S)-$$"
        if mkdir -p "$gotmp_run" 2>/dev/null; then
            export GOTMPDIR="$gotmp_run"
            echo "wisp-cli-tests.sh: ticket 310 AC#1: GOTMPDIR=$GOTMPDIR - this run's" \
                "t.TempDir() family lands here; nothing deletes it."
        else
            echo "wisp-cli-tests.sh: ticket 310 AC#1: DEGRADED - cannot create $gotmp_run," \
                "so GOTMPDIR stays unset and a test binary that dies leaves its temp dir in" \
                "the OS temp root again. Not a silent fallback: this line is the finding." >&2
        fi
    else
        echo "wisp-cli-tests.sh: ticket 310 AC#1: DEGRADED - LOCALAPPDATA is empty, so this step has no" \
            "measured-safe root to build (a root on another volume turns the winsec seal cases" \
            "red, see the note in scripts/portable-tests.sh). GOTMPDIR stays unset." >&2
    fi
fi

if [ "$(go env GOOS)" != windows ]; then
    # Not a skip and not a pass: this step is wired into the job whose runner has
    # the shared DLLs. On linux cmd/wisp builds and runs but 19 cases are red
    # (measured, in the header), which is a finding for the package owner, not a
    # thing to be quietly absorbed here.
    echo "wisp-cli-tests.sh: GUARD - this is the windows leg of the cmd/wisp gate," \
        "GOOS=$(go env GOOS). See the header for the ubuntu reading; do not relax this line." >&2
    exit 2
fi

# The pinned native dependency list is deps.toml, not a copy of file names here:
# the sections [sherpa-onnx.dll."<name>"] are what fetch-deps.ps1 extracts into
# third_party/sherpa-onnx/, so a pin added there is required here automatically.
dll_dir="$root/third_party/sherpa-onnx"
if [ ! -d "$dll_dir" ]; then
    echo "wisp-cli-tests.sh: GUARD - $dll_dir does not exist. The windows job stages it" \
        "(actions/cache on the deps.toml key, then scripts/build.ps1 -> fetch-deps.ps1)." \
        "Without it the step below would die at process load with 0xc0000135 and look like" \
        "a broken package. This is ticket 98's hole, named." >&2
    exit 1
fi

need=$(sed -n 's/^\[sherpa-onnx\.dll\."\(.*\)"\]$/\1/p' "$root/deps.toml")
if [ -z "$need" ]; then
    echo "wisp-cli-tests.sh: GUARD - no [sherpa-onnx.dll.*] sections found in deps.toml, so" \
        "this script has nothing to check. Empty instrument: red, not green." >&2
    exit 1
fi
missing=""
while IFS= read -r dll; do
    [ -n "$dll" ] || continue
    if [ ! -f "$dll_dir/$dll" ]; then missing="$missing $dll"; fi
done <<<"$need"
if [ -n "$missing" ]; then
    echo "wisp-cli-tests.sh: GUARD - pinned DLL(s) absent from $dll_dir:$missing" >&2
    echo "wisp-cli-tests.sh: cmd/wisp's test binary links the sherpa cgo import library, so a" >&2
    echo "wisp-cli-tests.sh: missing DLL fails at LOAD time (0xc0000135) before any test runs." >&2
    echo "wisp-cli-tests.sh: this is a setup failure of the step ahead of it, not a package bug." >&2
    exit 1
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
            echo "wisp-cli-tests.sh: ticket 311 AC#1 shape B: INCONCLUSIVE - no volume can be read from $wisp_311_name=$wisp_311_value (reference %LOCALAPPDATA%=$wisp_311_ref_path), so nothing is judged here and this step runs on the root as given. This is NOT the check passing: registered as debt Q-311-r1-1 on ticket 311." >&2
            continue
        fi
        if [ "$wisp_311_vol" != "$wisp_311_ref_vol" ]; then
            wisp_311_cross="$wisp_311_cross $wisp_311_name=$wisp_311_value(volume $wisp_311_vol)"
        fi
    done
    if [ -n "$wisp_311_cross" ]; then
        echo "wisp-cli-tests.sh: ticket 311 AC#1 shape B: GUARD - the test temp root names ANOTHER VOLUME than %LOCALAPPDATA%: reference volume $wisp_311_ref_vol (path $wisp_311_ref_path), offender:$wisp_311_cross" >&2
        echo "wisp-cli-tests.sh: ticket 311 AC#1 shape B: this gate's COLOUR is a function of that volume - the same unchanged internal/winsec code reads PASS=58 FAIL=0 on a volume-$wisp_311_ref_vol root and PASS=56 FAIL=2 with Go's own TempDir cleanup failing on a D:/tmp root (ticket 311). Refusing to run is the fix; reporting a re-coloured gate as a verdict is not." >&2
        echo "wisp-cli-tests.sh: ticket 311 AC#1 shape B: point WISP_TEST_GOTMP_ROOT (or GOTMPDIR) at a directory on volume $wisp_311_ref_vol, or unset both, then run this step again. Do NOT relax an assertion, threshold, golden or SLO to make it green." >&2
        exit 2
    fi
fi

# PATH is handed over in the form the loader accepts, which is NOT the form you
# would guess. Both were measured on this host at this HEAD:
#   `pwd -W` form (D:/work/.../third_party/sherpa-onnx) -> exit status 0xc0000135,
#       the ticket 98 symptom, reproduced by this script's own first draft;
#   the shell's own path form (/d/work/... with $dll_dir) -> PASS=33 FAIL=0 SKIP=0.
# MSYS2 rewrites PATH for native child processes, and the drive-letter form with a
# space in it survives that rewrite in a shape the DLL search does not accept. So
# the plain value is the correct one, and `pwd -W` is a trap here.
export PATH="$dll_dir:$PATH"

echo "wisp-cli-tests.sh: scope=./cmd/wisp/ dll_dir=$dll_dir (pinned: $(printf '%s\n' $need | tr '\n' ' '))"
echo "wisp-cli-tests.sh: handing the verdict to portable-tests.sh --scope=cli"
bash "$portable" --scope=cli
