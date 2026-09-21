#!/usr/bin/env bash
# scripts/portable-tests.sh - the runner behind CI's "Portable package tests"
# step (ticket 93). It exists because that step used to be a bare
#
#	go test <16 packages> -count=1
#
# and Go books a SKIP as `ok`: internal/risk's TestSyncRegistryProbeLive skipped
# on Windows AND on ubuntu, printed `ok  github.com/CarlosShao/wisp/internal/risk`,
# and the step exited 0. A case that could never evaluate anything was booked as
# a pass - on the platform where its subject (the registry hive) does not exist,
# and on the platform where the machine holds no record to read. The step's
# output was byte-for-byte what "everything concluded" looks like.
#
# THE INSTRUMENT IS NOT NEW (ticket 93's brief: reuse its shape, do not build a
# weaker one). Every verdict comes from tools/d22scan/runtests.sh, written by
# ticket 71 AC#3 and stricter than anything this script could add:
#   1. a non-zero `go test` exit code propagates unchanged;
#   2. any top-level `--- SKIP` in what actually ran is fatal;
#   3. zero top-level PASS *and* zero FAIL is fatal (the green no-op);
# and it forces -v -count=1, so the four numbers printed below can never be
# vacuous. This script adds exactly one thing on top: a per-platform ledger of
# cases declared NOT IN SCOPE for a stated reason, plus the check that stops
# that declaration from rotting. It subtracts nothing.
#
# WHY THIS IS NOT allowlist.txt (AC#2 forbids growing that file, and it is the
# orchestrator's surface - it is untouched here). An allowlist row is
# one-directional and silent. Every entry below is an exact top-level test name
# + the package that declares it + a platform + a class + a reason, and on every
# run each active entry is re-verified against the compiled test binary for THIS
# platform via `go test -list`. Rename the test, delete it, or bury it behind a
# build tag, and the entry goes stale and this step goes red - which is exactly
# the failure mode a plain -skip list would have smuggled in. An unaccounted SKIP
# also still goes red, so a new skip cannot be parked here by accident: adding an
# entry is a deliberate, reviewed, reason-carrying edit that shows up in the CI
# log, not just in a diff.
#
# Two entries are "the platform has no such API" cases whose honest cure is the
# platform-layer move AC#2 prescribes rather than a ledger row. They are recorded
# with the reason that keeps them here (one file matches this ticket's forbidden
# pattern internal/risk/pathresolver*.go; the other needs a second volume that no
# CI host has) and named as next= work on the ticket. The registry case that
# COULD be moved was moved, in the same commit, to syncdirs_windows_test.go.
#
# Usage:
#   bash scripts/portable-tests.sh                      # the 16-package CI portable scope
#   bash scripts/portable-tests.sh ./internal/proc/ ... # a narrower scope (windows job)
set -eu -o pipefail

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
root=$(CDPATH= cd -- "$here/.." && pwd)
cd "$root"

strict="$root/tools/d22scan/runtests.sh"
if [ ! -f "$strict" ]; then
    echo "portable-tests.sh: $strict is missing - this script is a shell around it and" \
        "refuses to become the weaker instrument it was meant to replace." >&2
    exit 2
fi

# The scope CI's test-core job runs, stated once; ci.yml points here so the list
# cannot drift between a workflow edit and a local re-run.
scope=()
if [ $# -eq 0 ]; then
    scope=(
        ./internal/agent/... ./internal/llm/... ./internal/config/...
        ./internal/memory/... ./internal/observe/... ./internal/secret/...
        ./internal/risk/... ./internal/statemachine/... ./internal/session/...
        ./internal/watchdog/... ./internal/tools/... ./internal/models/...
        ./internal/buildinfo/... ./internal/audio/... ./internal/proc/...
        ./internal/panel/...
    )
else
    scope=("$@")
fi

goos=$(go env GOOS)
if [ -z "$goos" ]; then
    echo "portable-tests.sh: go env GOOS came back empty" >&2
    exit 2
fi

# name|package|platform|class|reason
#   class: fixture = the machine or OS cannot supply the subject at all
#          opt-in  = deliberately slow or real-network, gated by an env var
#          reexec  = not a test, the child-process entry point of another test
ledger=(
    "TestDefaultDeadlineWallClockMeasurement|./internal/agent/approval/|any|opt-in|300s wall-clock measurement gated by WISP_84_MEASURE=1 (ticket84_no_owner_test.go:224, docs/evidence/s1/84-ac1-bounded-wait.md); it measures a deadline instead of asserting a property, so a CI run of it produces a duration nobody checks"
    "TestSubprocessCrashWriter|./internal/memory/|any|reexec|child-process entry point of TestCrashRecoveryKillMidWrite (concurrent_test.go:186); not a test of its own, it skips whenever it is not the re-exec'd child"
    "TestHelperProcess|./internal/proc/|windows|reexec|child-process entry point of the job-object cases (jobscope_windows_test.go:87, //go:build windows); skips unless re-exec'd with the helper env var set"
    "TestLiveWasapiSmoke|./internal/audio/|windows|fixture|live WASAPI capture needs WISP_LIVE_MIC=1 and a physical microphone (hotplug_test.go:527, //go:build windows); a hosted runner has no audio endpoint, so the case has no subject there. Real-hardware smoke is ticket 16's"
    "TestRealDownloadVadThroughPipeline|./internal/models/|any|opt-in|real-network spot check gated by WISP_IT_REAL_MIRROR=1 (manifest_real_test.go:120); CI has no business pulling model archives from a third-party mirror"
    "TestRealDownloadPuncArchiveThroughPipeline|./internal/models/|any|opt-in|real-network spot check gated by WISP_IT_REAL_MIRROR=1 (manifest_real_test.go:147); same reason as the VAD row"
    "TestSyncRegistryProbeLive|./internal/risk/|windows|fixture|P12 live registry evidence; on ubuntu it is now out of scope by platform instead of by skip (ticket 93 AC#2 moved it verbatim into syncdirs_windows_test.go), and on a Windows host whose HKCU Accounts key carries no UserFolder the hive has nothing to say. Measured on a windows host at HEAD 84e43af: --- SKIP at syncdirs_windows_test.go:133. As of ticket 110 AC#4 the windows leg carries ./internal/risk/ in its scope, so this row is re-verified against the compiled WINDOWS test binary (go test -list) and printed with its reason in that step's log - which is the step-level answer R-93-4 asked for. Remedy for a real run: a host with a sync record, or fold the shape check into the fixture-driven case. next= on ticket 93"
    "TestC26RewrittenSyncRootDoesNotDisarmSuspectNet|./internal/risk/|linux|fixture|needs a handle-resolved form of an existing directory, which POSIX has no object for (pathresolver_rewrite_account_test.go:138). AC#2's remedy is a platform-layer move, impossible from this ticket: the file matches this ticket's forbidden pattern internal/risk/pathresolver*.go. The case RUNS AND PASSES on windows, so the assertion is not lost, it is only un-evaluable here. next= an owner for that file"
    "TestD34WriteMatrix|./internal/tools/|linux|fixture|the table needs a second volume from otherVolumeDir (fs_write_test.go:140/151), and the volume-identity call it uses is a Windows API - on POSIX there is no object to ask. In a CI ubuntu container /tmp and / are one filesystem anyway. It RUNS AND PASSES on a two-volume windows host, so if a windows host ever has a single volume this entry does not cover it and the step goes red on purpose"
    "TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop|./internal/tools/|linux|fixture|same otherVolumeDir requirement as the write matrix row above (fs_write_test.go:677); on POSIX there is no volume identity to compare, so the case has no subject. Runs and passes on a two-volume windows host"
)

# Names and packages, and the -list universe for this platform: built once, from
# the scope PLUS every ledger package, so an entry whose package is outside a
# narrowed scope is still checked rather than silently spared.
ledger_pkgs=()
for entry in "${ledger[@]}"; do
    IFS='|' read -r _name pkg _platform _class _reason <<<"$entry"
    case " ${ledger_pkgs[*]-} " in
    *" $pkg "*) ;;
    *) ledger_pkgs+=("$pkg") ;;
    esac
done

universe_file=$(mktemp 2>/dev/null || echo "$root/.portable-list.$$.log")
if ! go test -list '.*' "${scope[@]}" "${ledger_pkgs[@]}" >"$universe_file" 2>&1; then
    cat "$universe_file"
    echo "portable-tests.sh: go test -list failed for [${scope[*]} ${ledger_pkgs[*]}] - the ledger" \
        "cannot be verified against a test binary that does not build." >&2
    rm -f "$universe_file"
    exit 1
fi

active=()
stale=()
for entry in "${ledger[@]}"; do
    IFS='|' read -r name pkg platform class reason <<<"$entry"
    if [ -z "$name" ] || [ -z "$pkg" ] || [ -z "$platform" ] || [ -z "$class" ] || [ -z "$reason" ]; then
        echo "portable-tests.sh: malformed ledger entry, need name|pkg|platform|class|reason: $entry" >&2
        rm -f "$universe_file"
        exit 2
    fi
    if ! printf '%s\n' "$name" | grep -Eq '^Test[A-Za-z0-9_]+$'; then
        echo "portable-tests.sh: ledger name $name is not a bare top-level Test identifier," \
            "so it cannot be turned into an anchored -skip pattern." >&2
        rm -f "$universe_file"
        exit 2
    fi
    case $platform in
    any | "$goos") ;;
    *) continue ;; # accounted on the other platform's job, not on this one
    esac
    active+=("$name")
    # The name must exist in the compiled test binary for the package it is
    # accounted against. Two traps here, both paid for: `grep -qx` in a pipeline
    # exits on the first match, which SIGPIPEs `go test` and - with pipefail -
    # reads as "not matched" on linux (the check would then be red on every
    # platform for every entry, i.e. an instrument nobody trusts); and an exact
    # whole-line match is required, because a substring test would let
    # TestFooBar vouch for an entry named TestFoo. Here-string, no pipeline.
    listed=$(go test -list "^${name}\$" "$pkg" 2>/dev/null || true)
    if ! grep -qx "$name" <<<"$listed"; then
        stale+=("$name  [$pkg  platform=$platform]")
    fi
done

echo "portable-tests.sh: platform=$goos scope=[${scope[*]}]"
echo "portable-tests.sh: ${#ledger[@]} ledger entries, ${#active[@]} accounted on this platform:"
for entry in "${ledger[@]}"; do
    IFS='|' read -r name _pkg platform class reason <<<"$entry"
    case $platform in
    any | "$goos") printf 'portable-tests.sh:   %-52s %s\n' "$name" "[$class] $reason" ;;
    esac
done

if [ "${#stale[@]}" -ne 0 ]; then
    echo "portable-tests.sh: these ledger entries name NO TEST in the compiled test binary for $goos:" >&2
    printf '    %s\n' "${stale[@]}" >&2
    echo "portable-tests.sh: a -skip pattern that matches nothing is the green no-op ticket 71 was" \
        "written about. Either the test was renamed or deleted (fix or delete the entry), or it" \
        "moved behind a build tag that excludes this platform (then say so in the entry)." >&2
    rm -f "$universe_file"
    exit 1
fi
rm -f "$universe_file"

if [ "${#active[@]}" -eq 0 ]; then
    # Nothing applies on this platform. Use a pattern that cannot match any test
    # name, so the flag is always present and its effect is always exercised.
    skip_pattern='^portable_tests_ledger_is_empty_on_this_platform$'
else
    alt=""
    for name in "${active[@]}"; do
        if [ -z "$alt" ]; then alt=$name; else alt="$alt|$name"; fi
    done
    skip_pattern="^($alt)\$"
fi
echo "portable-tests.sh: -skip pattern built from the ledger: $skip_pattern"

capture=$(mktemp 2>/dev/null || echo "$root/.portable-tests.$$.log")
rc=0
# tee costs nothing and keeps the strict runner's verdict intact while giving
# this script the bytes for the four numbers and for naming an offender.
sh "$strict" "${scope[@]}" -count=1 -skip "$skip_pattern" 2>&1 | tee "$capture" || rc=$?

count() { grep -c "$1" "$capture" || true; }
ran=$(count '^=== RUN')
passed=$(count '^--- PASS')
failed=$(count '^--- FAIL')
skipped=$(count '^--- SKIP')

echo "portable-tests.sh: four numbers (all from -v output): === RUN=$ran  --- PASS=$passed  --- FAIL=$failed  --- SKIP=$skipped"
if [ "$skipped" -ne 0 ]; then
    echo "portable-tests.sh: unaccounted SKIP lines, each with the file:line and reason it printed:" >&2
    grep -B1 -- '^--- SKIP' "$capture" >&2 || true
fi

rm -f "$capture"
[ "$rc" -eq 0 ] || echo "portable-tests.sh: strict runner exited $rc for scope=[${scope[*]}]" >&2
exit "$rc"
