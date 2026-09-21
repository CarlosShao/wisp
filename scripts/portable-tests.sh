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
# TICKET 111 - THE THREE GUARDS BELOW THE SCOPE, AND WHY THIS SCRIPT OWES THEM.
# "A package name in a list" and "that package was tested" were never connected
# here. AC#1's census over CI history found 20 tested packages out of `go list
# ./...`=33, and the reason is not only the missing five: it is also that an entry
# with NOTHING behind it cannot fail. internal/session and internal/watchdog are
# in the list below and have never had a test file - doc.go only, both marked
# "DEFERRED: implemented by ticket 28 / 42, this ticket only freezes the package
# boundary" - and a third case (internal/agent/scheduler) was reachable only
# because the list said ./internal/agent/... . Measured, not argued:
# `runtests.sh ./internal/session/` alone exits 1 (no top-level result at all),
# but INSIDE a 23-package invocation the strict runner's "PASS>0" test is
# satisfied by the other packages, so the empty one contributes nothing and stays
# invisible forever. That is why the scope is now resolved to import paths and
# audited per package by three guards, none of which can be talked around:
#   GUARD A (AC#3)  a declared package that compiles ZERO test files on this
#                   platform is a red step - the entry must gain a denominator
#                   or leave the list out loud;
#   GUARD B (AC#2)  every declared package must print its OWN top-level result
#                   line in this run, matched anchored and literal, so "the word
#                   appeared in the log" can never count as "the package ran";
#   GUARD C (AC#3)  the named scopes pin their resolved import paths, so
#                   deleting an entry (or a glob quietly narrowing) fails the
#                   step instead of shortening the run.
#
# Usage:
#   bash scripts/portable-tests.sh                    # the core (ubuntu) CI scope
#   bash scripts/portable-tests.sh --scope=windows    # the test-windows scope
#   bash scripts/portable-tests.sh --scope=cli        # cmd/wisp, after third_party
#   bash scripts/portable-tests.sh ./internal/proc/   # explicit scope (A and B still
#                                                     # apply; C has nothing to pin)
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

# The scopes CI runs, stated once, here - ci.yml names a scope, it does not carry
# a package list. That move is ticket 111 AC#2's: a workflow line naming a package
# is not evidence the package was tested, and the guard that makes deleting one
# fail loudly (GUARD C) only works if the list and its pin live together.
#
#   core    test-core (ubuntu). Every portable package with a test denominator on
#           both platforms. ./internal/session/ and ./internal/watchdog/ USED to
#           be here and are gone: they are doc.go-only boundary stubs
#           (DEFERRED: ticket 28 / ticket 42), so they could never go red and
#           never proved anything. ./internal/agent/... is spelled out as agent +
#           agent/approval for the same reason - the glob was quietly carrying
#           agent/scheduler, a third doc.go-only stub. When those packages get
#           code, put them back: GUARD A rejects them until they have a test file,
#           so re-entry cannot be lazy.
#   windows test-windows' portable step.
#   cli     cmd/wisp, whose test binary needs the sherpa DLLs (ticket 98's load
#           hole) - run by scripts/wisp-cli-tests.sh, which stages them first.
#
# ./internal/winsec/ joined the core list for ticket 111 AC#9: the ubuntu leg ran
# 16 globs that never named it, so the whole `!windows` half of the package that
# decides whether a path really gets an ACL had ZERO CI regression protection -
# ticket 113's linking leg (winsec_other.go) existed only in a laptop Docker run.
# Measured on the clean 8fe5c7c snapshot inside a real ubuntu container:
# === RUN=35 / PASS=20 / FAIL=0 / SKIP=0, rc=0, and winsec printed its own
# top-level result line. That is a POSIX denominator of 20 assertions, not a
# compile-only claim - `GOOS=linux go vet` cannot produce either number.
# The windows half stays where ticket 110 put it (its own step, own guards); the
# two legs now cover different files of the same package, which is the point.
mode=core
case "${1-}" in
    --scope=*) mode=${1#--scope=}; shift ;;
esac

# Each named scope carries the set its globs resolve to, on one line per import
# path. GUARD C compares against this, and --scope=census reads it, so the list a
# step claims and the audit of that claim cannot be maintained apart. Refresh with:
#   bash scripts/portable-tests.sh --scope=census
core_pin='
github.com/CarlosShao/wisp/cmd/llmrecord
github.com/CarlosShao/wisp/internal/agent
github.com/CarlosShao/wisp/internal/agent/approval
github.com/CarlosShao/wisp/internal/audio
github.com/CarlosShao/wisp/internal/ball
github.com/CarlosShao/wisp/internal/buildinfo
github.com/CarlosShao/wisp/internal/config
github.com/CarlosShao/wisp/internal/llm
github.com/CarlosShao/wisp/internal/llm/adaptertest
github.com/CarlosShao/wisp/internal/llm/anthropic
github.com/CarlosShao/wisp/internal/llm/golden
github.com/CarlosShao/wisp/internal/llm/openaichat
github.com/CarlosShao/wisp/internal/llm/openairesponses
github.com/CarlosShao/wisp/internal/memory
github.com/CarlosShao/wisp/internal/models
github.com/CarlosShao/wisp/internal/observe
github.com/CarlosShao/wisp/internal/panel
github.com/CarlosShao/wisp/internal/perm
github.com/CarlosShao/wisp/internal/plugin
github.com/CarlosShao/wisp/internal/proc
github.com/CarlosShao/wisp/internal/risk
github.com/CarlosShao/wisp/internal/secret
github.com/CarlosShao/wisp/internal/statemachine
github.com/CarlosShao/wisp/internal/tools
github.com/CarlosShao/wisp/internal/winsec
'
win_pin='
github.com/CarlosShao/wisp/cmd/llmrecord
github.com/CarlosShao/wisp/internal/ball
github.com/CarlosShao/wisp/internal/config
github.com/CarlosShao/wisp/internal/perm
github.com/CarlosShao/wisp/internal/plugin
github.com/CarlosShao/wisp/internal/proc
github.com/CarlosShao/wisp/internal/risk
github.com/CarlosShao/wisp/internal/secret
'
cli_pin='
github.com/CarlosShao/wisp/cmd/wisp
'
winsec_pin='
github.com/CarlosShao/wisp/internal/winsec
'

scope=()
pinned=''
case $mode in
core)
    scope=(
        ./internal/agent/ ./internal/agent/approval/
        ./internal/llm/... ./internal/config/...
        ./internal/memory/... ./internal/observe/... ./internal/secret/...
        ./internal/risk/... ./internal/statemachine/...
        ./internal/tools/... ./internal/models/...
        ./internal/buildinfo/... ./internal/audio/... ./internal/proc/...
        ./internal/panel/... ./internal/ball/ ./internal/perm/
        ./internal/plugin/ ./cmd/llmrecord/ ./internal/winsec/
    )
    pinned=$core_pin
    ;;
windows)
    scope=(
        ./internal/proc/ ./internal/secret/ ./internal/config/ ./internal/risk/
        ./internal/ball/ ./internal/perm/ ./internal/plugin/ ./cmd/llmrecord/
    )
    pinned=$win_pin
    ;;
cli)
    # cmd/wisp's own tests. Only scripts/wisp-cli-tests.sh may call this, because
    # the test binary dies at LOAD time without the sherpa DLLs on PATH (ticket
    # 98), and that script stages them and refuses to run if they are not there.
    scope=(./cmd/wisp/)
    pinned=$cli_pin
    ;;
census)
    scope=()
    pinned=''
    ;;
*)
    echo "portable-tests.sh: unknown --scope=$mode (known: core, windows, cli, census)" >&2
    exit 2
    ;;
esac
goos=$(go env GOOS)
if [ -z "$goos" ]; then
    echo "portable-tests.sh: go env GOOS came back empty" >&2
    exit 2
fi

if [ "$mode" = census ]; then
    # GUARD A's numbers for EVERY package in the module, plus which named scope
    # claims it. `t=`/`x=` are the test files that COMPILE INTO THIS PLATFORM's
    # test binary, so a build tag that empties a package on one platform shows up
    # as 0 there and non-zero on the other. A NO-SCOPE row is zero coverage said in
    # the same voice as a covered one, which is what AC#1 asked for: the ticket's
    # "CI 测了 20/33" claim is only meaningful if the other 13 are on the page.
    all=$(go list ./... 2>/dev/null | sort -u)
    n=$(printf '%s\n' "$all" | grep -c . || true)
    echo "portable-tests.sh: census GOOS=$goos - go list ./... = $n packages (one row each)"
    printf 'portable-tests.sh: %-46s %-11s %s\n' PACKAGE 'TESTS(t/x)' 'CLAIMED BY'
    empty=0
    noscope=0
    while IFS= read -r p; do
        [ -n "$p" ] || continue
        counts=$(go list -f '{{len .TestGoFiles}}/{{len .XTestGoFiles}}' "$p" 2>/dev/null || echo '?/?')
        where=''
        for m in core windows cli winsec; do
            case $m in
            core) pin=$core_pin ;; windows) pin=$win_pin ;; cli) pin=$cli_pin ;;
            winsec) pin=$winsec_pin ;;
            esac
            if printf '%s\n' "$pin" | grep -qxF "$p"; then where="$where$m"; fi
        done
        if [ -z "$where" ]; then where=' NO-SCOPE'; noscope=$((noscope + 1)); fi
        case $counts in
        0/0) empty=$((empty + 1)); where="$where <-NO-TESTS" ;;
        esac
        printf 'portable-tests.sh: %-46s %-11s %s\n' "$p" "$counts" "$where"
    done <<<"$all"
    echo "portable-tests.sh: census totals: packages=$n with-zero-compiled-tests=$empty claimed-by-no-scope=$noscope"
    exit 0
fi

if [ $# -gt 0 ]; then
    # Explicit paths are the local-debug form (and how winsec-tests.sh delegates).
    # GUARD A and GUARD B still apply to whatever is named; only GUARD C's pin is
    # skipped, because an ad-hoc list has no expectation to be faithful to.
    scope=("$@")
    pinned=''
fi
if [ ${#scope[@]} -eq 0 ]; then
    echo "portable-tests.sh: empty scope (mode=$mode) - refusing to be a green no-op" >&2
    exit 2
fi

# ---- resolve the scope once, so all three guards read the same truth ---------
resolved=$(mktemp 2>/dev/null || echo "$root/.portable-resolved.$$.txt")
if ! go list "${scope[@]}" >"$resolved" 2>&1; then
    cat "$resolved"
    echo "portable-tests.sh: go list [${scope[*]}] failed - the scope cannot be audited" \
        "against a pattern that does not resolve." >&2
    rm -f "$resolved"
    exit 1
fi
grep -v '^$' "$resolved" | sort -u >"$resolved.sorted" || true
mv "$resolved.sorted" "$resolved"
pkgcount=$(wc -l <"$resolved" | tr -d '[:space:]')

# GUARD C: the named scopes pin their own resolved set. Deleting an entry from the
# list above, or letting a glob quietly stop covering a package, is a red step -
# not a run that tests one fewer thing and prints the same green. This is the
# shape ticket 93's "条目腐坏即红" and winsec-tests.sh's guard 1 already use; it
# fires in BOTH directions, so renaming ball into a package nobody named is caught
# too.
if [ -n "$pinned" ]; then
    want=$(printf '%s\n' "$pinned" | grep -v '^[[:space:]]*$' | sort -u)
    got=$(cat "$resolved")
    if [ "$want" != "$got" ]; then
        {
            echo "portable-tests.sh: GUARD C - scope mode=$mode resolved to a DIFFERENT package"
            echo "portable-tests.sh:   set than the one pinned next to it. Pinned: $(printf '%s\n' "$want" | wc -l | tr -d '[:space:]'), resolved: $pkgcount."
            diff <(printf '%s\n' "$want") <(printf '%s\n' "$got") | sed 's/^/portable-tests.sh:   /' || true
            echo "portable-tests.sh: a deleted or newly-uncovered package must fail the step, not"
            echo "portable-tests.sh: shorten it. Either restore the scope entry or update the pin"
            echo "portable-tests.sh: in the SAME commit and say why in the CI log."
        } >&2
        rm -f "$resolved"
        exit 1
    fi
fi

# GUARD A: the loud empty denominator. `.TestGoFiles`+`.XTestGoFiles` are the
# files that COMPILE INTO THIS PLATFORM's test binary, so a build tag that
# excludes a package's tests on this GOOS is caught here as well as a missing
# file. An entry like this can never go red and never proves anything, which is
# the exact hole internal/session and internal/watchdog sat in.
denom=$(go list -f '{{if or .TestGoFiles .XTestGoFiles}}{{else}}EMPTY {{.ImportPath}}
{{end}}' "${scope[@]}" | grep '^EMPTY ' || true)
if [ -n "$denom" ]; then
    {
        echo "portable-tests.sh: GUARD A - these packages are declared in scope mode=$mode but"
        echo "portable-tests.sh:   compile NO test file at all for GOOS=$goos:"
        printf '%s\n' "$denom" | sed 's|^|portable-tests.sh:   |'
        echo "portable-tests.sh: a scope entry with no denominator cannot fail, so keeping it there"
        echo "portable-tests.sh: is a false claim of coverage (ticket 111 AC#3). Give the package a"
        echo "portable-tests.sh: test, or take the entry out and say where its coverage lives."
    } >&2
    rm -f "$resolved"
    exit 1
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
    "TestWorkspaceSwitchRefusesAJunctionToOutside|./internal/tools/|linux|fixture|ticket111 AC#10. The case is AC#3(ii) with a REAL junction: on windows it builds one with cmd /c mklink /J and asserts the switch is refused (runs and passes here). On POSIX there is nothing to refuse, because the detection it exercises is itself a Windows implementation - internal/tools/paths_workspace_test.go:198 skips with 'C26's reparse detection is a Windows implementation (risk.pathresolver_other.go reparseComponents returns nil elsewhere); nothing to deny on linux'. Registered rather than deleted because the file's own comment promises 'the skip is reported, never averaged away', and runtests.sh making every SKIP fatal is what turns that promise into a red step: an unaccounted skip on the ubuntu leg (test-core step 7, measured PASS=578 FAIL=0 SKIP=1 on run 35599458439) is exactly what this row exists to name. OWNER of the POSIX half: whoever lands reparseComponents in risk.pathresolver_other.go - that file is forbidden to this ticket (internal/risk/pathresolver*.go, same wall as the row above), so the POSIX-equivalent criterion is NOT faked here. The row is re-verified against the compiled LINUX test binary on every run: delete the skip, and the case starts running here; delete the case, and GUARD's staleness check goes red."
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

# GUARD B (ticket 111 AC#2 + AC#8): a package counts as tested only if THIS run
# printed a top-level result line for its exact import path.
#
# The anchor is the whole point. R-110-3 is a matching expression that let a bare
# substring "winsec" be read as 18 hits when the package was tested 0 times; the
# general form of that mistake is "the word appears in the log" == "the package
# ran". So: line-anchored to ^ok/FAIL (a result line Go prints once per package),
# the import path regex-ESCAPED (an unescaped `.` in github.com/... matches any
# character, so internalXwinsec would vouch for internal/winsec), and followed by
# whitespace or end-of-line, which is what keeps internal/winsecfrom counting for
# internal/winsec. A seeded positive control is on the ticket: a test that only
# PRINTS another package's name contributes nothing to this table.
escape_re() { printf '%s' "$1" | sed 's/[][\\^$.*+?(){}|]/\\&/g'; }
# Go's package result line is `ok  \t<import path>\t0.123s`, tab-delimited with a
# duration (or `[build failed]`). Demanding that tail is what keeps a test that
# writes `os.Stdout.WriteString("ok  github.com/...internal/perm\t0.01s")` from
# impersonating a package: it can fake the prefix at column 0, faking the whole
# shape plus being inside the right package's output block is a different claim.
# The residual limit is stated rather than hidden: this is a line-shape match on
# captured output, not a property of the go test protocol. The `-count=1` the
# strict runner forces is what makes `[no test files]` and `(cached)` impossible
# to confuse with a real run here.
result_tail='[[:space:]]+([0-9]+\.[0-9]+s|\[build failed\])$'

missing=0
declare -A per_pkg=()
while IFS= read -r pkg; do
    [ -n "$pkg" ] || continue
    line=$(grep -E "^(ok|FAIL)[[:space:]]+$(escape_re "$pkg")${result_tail}" "$capture" | tail -1 || true)
    if [ -z "$line" ]; then
        per_pkg["$pkg"]='<NO TOP-LEVEL RESULT LINE> missing'
        missing=$((missing + 1))
    else
        case $line in
        ok*) per_pkg["$pkg"]="ok (own line)" ;;
        *) per_pkg["$pkg"]="FAIL (own line)" ;;
        esac
    fi
    printf 'portable-tests.sh:   %-14s %s\n' "${per_pkg[$pkg]}" "$pkg"
done <"$resolved"

if [ "$missing" -ne 0 ]; then
    {
        echo "portable-tests.sh: GUARD B - $missing of $pkgcount packages in scope mode=$mode printed"
        echo "portable-tests.sh:   NO top-level result line, so this step tested them ZERO times:"
        for p in "${!per_pkg[@]}"; do
            if [ "${per_pkg[$p]}" = '<NO TOP-LEVEL RESULT LINE> missing' ]; then
                echo "portable-tests.sh:   $p"
            fi
        done
        echo "portable-tests.sh: a package that goes missing here is one that a build tag emptied, that"
        echo "portable-tests.sh: was renamed, or that the run never reached. That is a red, not a run"
        echo "portable-tests.sh: that got shorter (ticket 111 AC#2/AC#8)."
    } >&2
    if [ "$rc" -eq 0 ]; then rc=1; fi
fi

rm -f "$capture" "$resolved"
[ "$rc" -eq 0 ] || echo "portable-tests.sh: strict runner exited $rc for scope=[${scope[*]}]" >&2
exit "$rc"
