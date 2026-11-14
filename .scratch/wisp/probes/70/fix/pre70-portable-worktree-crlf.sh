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
#   GUARD D (this ticket's AC#1, added by the r1 leg) the CENSUS refuses to exit 0
#                   while a package compiles a test file for this GOOS and no tier
#                   claims it. A/B/C audit the packages a scope NAMES; nothing
#                   audited the packages a scope FORGETS, which is why "20 of 33"
#                   could sit in a printed table that exited 0. GUARD C bites when
#                   a row is deleted from a list that still exists; GUARD D bites
#                   when the list and its pin are deleted together, and when a new
#                   package arrives with tests and no entry at all (internal/projctx,
#                   ticket 200, is the specimen).
#
# Usage:
#   bash scripts/portable-tests.sh                    # the core (ubuntu) CI scope
#   bash scripts/portable-tests.sh --scope=windows    # the test-windows scope
#   bash scripts/portable-tests.sh --scope=cli        # cmd/wisp, after third_party
#   bash scripts/portable-tests.sh --scope=winsec     # the internal/winsec tier (ticket 251)
#   bash scripts/portable-tests.sh --scope=speech     # internal/speech via the staging
#                                                     # carrier scripts/speech-tests.sh
#                                                     # (ticket 409, wave99)
#   bash scripts/portable-tests.sh ./internal/proc/   # explicit scope (A and B still
#                                                     # apply; C has nothing to pin -
#                                                     # UNLESS the list IS one tier's own
#                                                     # scope, see the ticket 251 AC#1 block)
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

# ---- TICKET 310 AC#1 (shape A-G): the test temp family gets a dedicated root.
#
# WHAT THIS BUYS, MEASURED. A test binary that dies without running its cleanup
# leaves its `t.TempDir()` directory on disk: the census leg (310-a1 20-shapes.md)
# and the orchestrator (probes/orch/20261012/logs/orch-310-recheck-b.txt) each
# measured exactly one new directory per death, twice over, for the `go test
# -timeout` panic and for an external TerminateProcess. On this host those
# directories piled up at the top of %TEMP% (21 -> 44 in one day, 55 registered at
# the time of this writing), which is why the fleet's starting gate now counts them
# (probes/orch/gate-roster-304.md section 1). Redirecting the family moves WHERE the
# residue lands; it does not collect it - ticket 310's own rule is that deleting a
# directory is irreversible and needs a per-directory roster plus a non-implementer
# review, so this step names the root and deletes nothing.
#
# WHY GOTMPDIR AND NOT TMP/TEMP. testing's makeTempDir calls
# os.MkdirTemp(os.Getenv("GOTMPDIR"), pattern) (this host's GOROOT,
# src/testing/testing.go:1613), so GOTMPDIR moves exactly the t.TempDir()/
# os.CreateTemp family. os.TempDir() on Windows goes through Win32 GetTempPath and
# does NOT read GOTMPDIR, so a test that reads os.TempDir() literally keeps its
# answer. Redirecting TMP/TEMP instead would move those answers too - that is the
# A-T variant ticket 310 ruled against (12 test files plus the live-window ones).
#
# WHY THE ROOT IS NOT UNDER %TEMP%. The census leg measured that shape and named
# the trap: a root inside %TEMP% leaves the non-recursive
# `Get-ChildItem $env:TEMP -Directory` roster reading "0 new" while the residue is
# simply out of the ruler's reach. So the default root is a SIBLING of %TEMP% -
# under %LOCALAPPDATA%, outside the OS temp dir - not a directory inside it.
#
# WHY %LOCALAPPDATA% AND NOT "THE REPO'S VOLUME", MEASURED TODAY. The first draft
# of this step put the root on the repo's volume (`D:/tmp/...`) and the winsec tier
# went 58 PASS / 0 FAIL to 56 PASS / 2 FAIL on this host:
# TestGateJudgesThePrivateSetByResolvedSID (both sub-tests) and
# TestSealNarrowsAndNamesThePrincipalItRemovedBySID judge the ACL set they seal,
# and an inherited set from another volume is not the same set. The same scope run
# against a %LOCALAPPDATA%-sibling root came back 58/0, identical to the unredirected
# run. So the location of the root is NOT a free choice, and the value below is the
# one that was measured to leave every colour alone.
#
# WHY WINDOWS ONLY. The ticket, the roster row and every measurement behind both
# are about this host's %TEMP%. The ubuntu leg is left exactly as it is because
# this leg cannot run CI to measure what a redirected temp root does to those
# steps (registered as a debt on ticket 310, not silently assumed safe).
#
# WHY THE ROOT KEEPS THE NATIVE SPELLING (ticket 317). The first cut of this step
# folded `\` to `/` while composing the default root; the folded `C:/...` value
# then reached internal/risk's literal-spelling preconditions and two real
# `mklink /J` targets on the windows leg of run 38508311376, reading as four false
# reds (`precondition: deepestExistingAncestor(...) ... want
# "C:/Users/runneradmin/..."`, `mklink ... Invalid switch - "Users"`). Measured
# both ways on a dev host: those four cases are red under a `C:/...` root and
# green under a `C:\...` root with everything else unchanged. So the default below
# is composed with the native separator, and an explicit WISP_TEST_GOTMP_ROOT is
# still taken verbatim - this step does not rewrite the operator's spelling.
if [ -z "${GOTMPDIR:-}" ]; then
    if [ "$(go env GOOS)" = windows ]; then
        if [ -n "${LOCALAPPDATA:-}" ]; then
            gotmp_base=${WISP_TEST_GOTMP_ROOT:-"$LOCALAPPDATA\\wisp-test-gotmp"}
            gotmp_run="$gotmp_base\\portable-$(date +%Y%m%d-%H%M%S)-$$"
            if mkdir -p "$gotmp_run" 2>/dev/null; then
                export GOTMPDIR="$gotmp_run"
                echo "portable-tests.sh: ticket 310 AC#1: GOTMPDIR=$GOTMPDIR - this run's" \
                    "t.TempDir() family lands here. Nothing deletes it: ticket 310 requires a" \
                    "per-directory roster plus a non-implementer review before any removal."
            else
                echo "portable-tests.sh: ticket 310 AC#1: DEGRADED - cannot create" \
                    "$gotmp_run, so GOTMPDIR stays unset and a test binary that dies will" \
                    "leave its temp dir in the OS temp root again. Not a silent fallback:" \
                    "this line is the finding." >&2
            fi
        else
            echo "portable-tests.sh: ticket 310 AC#1: DEGRADED - LOCALAPPDATA is empty, so" \
                "this step has no measured-safe root to build (a root on another volume turns" \
                "the winsec seal cases red, see the note above). GOTMPDIR stays unset." >&2
        fi
    else
        echo "portable-tests.sh: ticket 310 AC#1: GOTMPDIR left as-is on GOOS=$(go env GOOS)" \
            "(the ticket's measurement surface is the Windows host's temp root)."
    fi
fi

# The scopes CI runs, stated once, here - ci.yml names a scope, it does not carry
# a package list. That move is ticket 111 AC#2's: a workflow line naming a package
# is not evidence the package was tested, and the guard that makes deleting one
# fail loudly (GUARD C) only works if the list and its pin live together.
#
#   core    test-core (ubuntu). Every portable package with a test denominator on
#           both platforms. ./internal/watchdog/ is NOT here and cannot be: it is a
#           doc.go-only boundary stub (DEFERRED: ticket 42), so it could never go red
#           and never proved anything; ./internal/agent/... is spelled out as agent +
#           agent/approval for the same reason - the glob was quietly carrying
#           agent/scheduler, a second doc.go-only stub. When those packages get code,
#           put them back: GUARD A rejects them until they have a test file, and
#           GUARD D names them the moment they do have one and no entry.
#           ./internal/session/ RE-ENTERED on 2026-10-05 (ticket 111 r1 leg): ticket
#           111 took it out because it was a doc.go stub, then ticket 224 put
#           grants.go + session.go (371 lines) and two test files behind it, and
#           nothing noticed - the package sat at zero CI coverage for four days while
#           the census that could see it was a command nobody ran. ./internal/projctx/
#           arrived the same way from ticket 200 (one xtest file, 442 lines, never in
#           any scope in any commit). Both compile tests on BOTH platforms, both carry
#           ZERO //go:build lines, ZERO t.Skip, ZERO os.Getenv and ZERO exec.Command
#           (measured by grep, so neither can hide a skip from runtests.sh), and every
#           package they import is already in this scope (memory / panel / risk / perm
#           / config). Their recorded green readings are on WINDOWS hosts
#           (probes/197/v1/M1a-sweep/run-whole_repository.txt lines 62 and 71, and
#           probes/200/r2/final-gate-committed.txt:27); the ubuntu reading has never
#           been taken by anyone, so this row's first ubuntu run is the reading this
#           leg owes - a red there is a finding to register on the ticket, not a
#           reason to touch an assertion, a threshold or a build tag (AC#2).
#   windows test-windows' portable step - the same denominator in its windows shape.
#   cli     cmd/wisp, whose test binary needs the sherpa DLLs (ticket 98's load
#           hole) - run by scripts/wisp-cli-tests.sh, which stages them first.
#   winsec  the internal/winsec tier (ticket 251 AC#1). This tier is NEW: the pin
#           below has existed since ticket 111 with NO reader but the census roster,
#           because `case $mode in` had no winsec branch AND the route CI actually
#           takes (scripts/winsec-tests.sh:94, which hands this script explicit
#           package paths) arrived in GUARD C's "nothing to pin" branch. So a split
#           of or a rename inside internal/winsec was, until ticket 251, a step that
#           kept passing while testing one package less. The tier's scope is written
#           as the GLOB ./internal/winsec/... on purpose: unlike the bare directory,
#           the glob puts a newly created subpackage INTO the resolved set, where
#           GUARD C's pin comparison is what makes it a red step.
#           census is not a scope - it is the audit of the four above, and it is in
#           the tiers= roster below because AC#4 (ticket 251) requires the tier names
#           census reads and the branches that give a tier its scope to be ONE list.
#   speech  the internal/speech tier (ticket 409, wave99). This tier exists because
#           the package cannot be judged by a bare `go test` invocation on Windows:
#           its cgo test binary links the sherpa-onnx import library and the loader
#           resolves the whole DLL tree at process LOAD time, before any Go code, so
#           the directory of the executable is the only search slot a runner can
#           still influence. The tier's run step therefore never reaches the strict
#           runner; it hands the package to scripts/speech-tests.sh (the staging
#           carrier: `go test -c`, stage the three pinned DLLs next to the binary,
#           run it from the package directory) and reads the verdict off the same
#           -v stream. The four numbers, the ledger staleness check and the guards
#           above the run step all still apply; only the per-package `ok <path>`
#           result-line table (GUARD B) is replaced by its speech-shaped
#           equivalent, because a raw binary run prints no such line.
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
#
# ./internal/speech/cer/ JOINED THE CORE LIST 2026-10-26 (ticket 15 AC#2, leg
# cerci-r1), and this paragraph is the reason rather than a description. SPEC-10:60-61
# and SPEC-04:125 require the two committed sets under testdata/asr-baseline/ to run
# WITH CI ("otherwise the baseline is empty talk"), and this package is the reader and
# the 6%/15% gate that scores them. It is claimable on the portable tier for measured
# reasons, not for convenience: ZERO cgo directives, ZERO //go:build lines, four test
# files that compile on both platforms, and - now that the corpus and the frozen
# 5-column sources.tsv are committed - `tools/d22scan/runtests.sh ./internal/speech/cer/`
# reports PASS=24 FAIL=0 SKIP=0, so the claim buys a real denominator instead of a
# green over nothing. Before this line the package sat UNCLAIMED-HAS-TESTS on the
# census the windows leg runs (a package that runs only on laptops is the hole ticket
# 111 AC#1 was filed for). The scope entry and the pin line move in the SAME commit
# because GUARD C compares them.
#
# ./internal/present/, ./internal/watchdog/ and ./internal/audio/aecspike/ JOINED
# THE CORE LIST 2026-11-01 (leg ci-r1, the census's own remedy for GUARD D). All
# three sat UNCLAIMED-HAS-TESTS on the windows census (run 41274267651: present
# 4/0, watchdog 3/0, aecspike 7/0) after their test files landed without a claim.
# Measured reasons, per the speech/cer bar above: present carries 4 test files and
# watchdog 4, both with ZERO //go:build lines, ZERO t.Skip, ZERO os.Getenv and ZERO
# exec.Command (grep-measured; both ran green on this windows host at the claim
# commit). aecspike carries 3 test files that compile on linux (the other 5 are
# //go:build windows rig/opt-in files gated by WISP_AEC_SPIKE / WISP_AEC392R2, so
# no portable carrier can run them - their readings stay dev-box readings per
# their own headers), and its non-test linux files import nothing that links
# sherpa, so the ubuntu test binary loads (the windows-shaped local reading dies
# at load on missing DLLs instead - 0xc0000135 - which is why its first ubuntu run
# below is the reading this claim owes). The parent's audio glob was narrowed from
# ./internal/audio/... to ./internal/audio/ + ./internal/audio/aecspike/ in the
# SAME commit because internal/audio/aec compiles ZERO test files on linux (its
# only test file is windows && cgo && webrtcap), so the glob would have pulled a
# package GUARD A refuses to declare. The scope entry and the pin line move in the
# SAME commit because GUARD C compares them. ./internal/speech/ stays unclaimed as
# documented below; its census row is reported to the owner, not resolved here.
#
# ./internal/speech/ IS DELIBERATELY NOT CLAIMED, here or in any other tier, and that
# is a reading rather than an omission: the parent package dies at LOAD time without
# the three sherpa DLLs (`go test ./internal/speech` = `exit status 0xc0000135` before
# a single test runs), it locates its models under third_party/spike-models (a fetched
# directory: `git ls-files third_party` = 0 files), and it carries 7 `t.Skip*` sites
# that fire whenever a session cannot be created. tools/d22scan/runtests.sh:234-239
# makes ANY `--- SKIP` fatal by design, so claiming it without a carrier that stages
# the DLLs and the models would be a permanent red - and silencing those skips is
# exactly what that runner exists to refuse. It stays on the census as
# UNCLAIMED-HAS-TESTS until such a carrier claims it; nothing is hidden by leaving it
# off this list.
#
# UPDATE 2026-11-08 (ticket 409, wave99): that reading is now superseded by a carrier
# that removes its premise. scripts/speech-tests.sh builds the test binary with
# `go test -c`, stages the three pinned DLLs NEXT TO IT and runs the binary from the
# package directory - the loader's FIRST search slot, ahead of System32 - so the
# shadowed OS-provided onnxruntime 1.17.1 can no longer win: the 41 shadow skips
# turned into real passes (measured 2026-11-08: PASS=83 FAIL=0 SKIP=6 rc=0, the 6
# residual skips naming a missing punctuation model, not the shadow). The speech tier
# below therefore claims github.com/CarlosShao/wisp/internal/speech and routes it
# through that carrier instead of the strict runner - a bare `go test` in a named
# scope would still die at LOAD time or book the model-missing skips against a runner
# that treats every skip as fatal, which is the exact permanent red the paragraph
# above refused to buy. The census row moves from UNCLAIMED-HAS-TESTS to CLAIMED BY
# speech in the same commit, because GUARD D compares this file's pins against the
# packages it sees.
#
# UPDATE 2026-11-08 (census leg, the last four PACKAGE-level holes GUARD D had):
# with speech closed the census still refused with unclaimed-with-tests=4 -
# cmd/balldebug, internal/agent/scheduler, internal/command and
# internal/plugin/manifestslot. All four are shape 甲 (has tests, never claimed):
# each compiles a non-zero test count for GOOS=windows (1/0, 2/0, 1/0, 7/0) and
# no named scope claimed it, so no CI step ran it and none could go red over it.
# The per-package reading lives in docs/evidence/s1/census-guard-d.md; the claims
# land here, split by the same bar this list already uses:
#
# - internal/agent/scheduler (ticket 47's acceptance legs) and internal/command
#   (ticket 213's) join CORE on the measured present/watchdog bar: zero
#   //go:build lines, zero t.Skip, zero os.Getenv, zero exec.Command
#   (grep-measured over their _test.go files), both GOOS=linux test binaries
#   compile and `go vet` clean, and both ran green on the windows host at the
#   claim commit. scheduler re-enters a scope it was explicitly parked out of:
#   ticket 111's Progress log registered it as a DEFERRED stub (zero tests, only
#   doc.go) whose implementation ticket must put the package back in scope when
#   its tests land - the tests landed, so this entry is that registered duty,
#   not a new claim.
# - cmd/balldebug (ticket 71's diff-gate tests, //go:build windows) and
#   internal/plugin/manifestslot (tickets 50/51/52/138, seven test files) join
#   the WINDOWS tier, and CORE is deliberately not their tier: balldebug's only
#   test file is windows-tagged, so the ubuntu core job would refuse the entry
#   through GUARD A the day it landed (zero compiled tests for GOOS=linux); and
#   manifestslot's tests exec tasklist and probe SystemRoot/LOCALAPPDATA - the
#   zero-exec.Command bar above is exactly what they fail, so their readings
#   belong on the leg that can supply the subject. Both ran green on the windows
#   host at the claim commit (manifestslot 17.7s, balldebug 0.4s, zero skips).
# The scope entries and the pin lines move in the SAME commit because GUARD C
# compares them. tiers=, the census arm and every guard's logic are untouched:
# this commit only spends claims; it does not loosen the instrument.
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
github.com/CarlosShao/wisp/internal/agent/scheduler
github.com/CarlosShao/wisp/internal/audio
github.com/CarlosShao/wisp/internal/audio/aecspike
github.com/CarlosShao/wisp/internal/ball
github.com/CarlosShao/wisp/internal/buildinfo
github.com/CarlosShao/wisp/internal/command
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
github.com/CarlosShao/wisp/internal/present
github.com/CarlosShao/wisp/internal/proc
github.com/CarlosShao/wisp/internal/projctx
github.com/CarlosShao/wisp/internal/risk
github.com/CarlosShao/wisp/internal/secret
github.com/CarlosShao/wisp/internal/session
github.com/CarlosShao/wisp/internal/speech/cer
github.com/CarlosShao/wisp/internal/statemachine
github.com/CarlosShao/wisp/internal/tools
github.com/CarlosShao/wisp/internal/watchdog
github.com/CarlosShao/wisp/internal/winsec
'
win_pin='
github.com/CarlosShao/wisp/cmd/balldebug
github.com/CarlosShao/wisp/cmd/llmrecord
github.com/CarlosShao/wisp/internal/audio
github.com/CarlosShao/wisp/internal/ball
github.com/CarlosShao/wisp/internal/config
github.com/CarlosShao/wisp/internal/perm
github.com/CarlosShao/wisp/internal/plugin
github.com/CarlosShao/wisp/internal/plugin/manifestslot
github.com/CarlosShao/wisp/internal/proc
github.com/CarlosShao/wisp/internal/projctx
github.com/CarlosShao/wisp/internal/risk
github.com/CarlosShao/wisp/internal/secret
github.com/CarlosShao/wisp/internal/session
'
cli_pin='
github.com/CarlosShao/wisp/cmd/wisp
'
winsec_pin='
github.com/CarlosShao/wisp/internal/winsec
'
speech_pin='
github.com/CarlosShao/wisp/internal/speech
'

# TICKET 251 AC#4 - ONE ROSTER, READ BY EVERYONE. The tier names used to be written
# TWICE in this file: as the `case $mode in` branches, and again as
# `for m in core windows cli winsec` inside the census loop. Two lists that can
# disagree, and at the ticket's own HEAD they DID (census printed `CLAIMED BY ...
# winsec` for a tier no branch implemented) while census exited 0. tiers= is now the
# only place the names are written; the census step re-reads THIS FILE's own
# `case $mode in` block and refuses to report anything unless the two sets are equal
# in both directions, so the roster cannot be maintained apart from the branches.
tiers='core windows cli winsec speech census'
# The winsec tier's scope, named once: the tier branch below, the explicit-path audit
# further down and the core tier's winsec row ALL read it, so "what the tier covers"
# cannot drift from "what the tier's pin claims" the way the tier NAME used to, and
# the two tiers that claim internal/winsec cannot resolve it two different ways
# (ticket 254 AC#2: core used to name the DIRECTORY here while the tier names the
# GLOB, which is exactly the shape that lets one tier go green on a package split
# while the other goes red on it).
winsec_dir=./internal/winsec/
winsec_scope=./internal/winsec/...

scope=()
pinned=''
case $mode in
core)
    scope=(
        ./internal/agent/ ./internal/agent/approval/ ./internal/agent/scheduler/
        ./internal/llm/... ./internal/config/...
        ./internal/memory/... ./internal/observe/... ./internal/secret/...
        ./internal/risk/... ./internal/statemachine/...
        ./internal/tools/... ./internal/models/...
        ./internal/buildinfo/... ./internal/command/ ./internal/audio/ ./internal/audio/aecspike/
        ./internal/proc/...
        ./internal/panel/... ./internal/ball/ ./internal/perm/
        ./internal/plugin/ ./cmd/llmrecord/ "$winsec_scope"
        ./internal/session/ ./internal/projctx/ ./internal/present/
        ./internal/watchdog/
        ./internal/speech/cer/
    )
    pinned=$core_pin
    ;;
windows)
    scope=(
        ./internal/proc/ ./internal/secret/ ./internal/config/ ./internal/risk/
        ./internal/ball/ ./internal/perm/ ./internal/plugin/ ./internal/plugin/manifestslot/ ./cmd/balldebug/ ./cmd/llmrecord/
        ./internal/session/ ./internal/projctx/ ./internal/audio/
    )
    pinned=$win_pin
    ;;
cli)
    # cmd/wisp's own tests. Only scripts/wisp-cli-tests.sh may call this, because
    # the test binary dies at LOAD time without the sherpa DLLs on PATH (ticket
    # 98), and that script stages them and refuses to run if they are not there.
    #
    # TICKET 409's FINDING, EXTENDED HERE (round 87): PATH is the WRONG lever for
    # the SECOND-level dependency. Windows resolves the whole import tree at
    # process load time, searching the executable's directory first and System32
    # before PATH - so on a host whose System32 carries an OS-provided
    # onnxruntime.dll (measured here: 1.17.1), the loader binds THAT copy
    # whatever PATH says, and every voice-backed test in cmd/wisp reports the
    # shadowing verdict instead of measuring the product. The 2026-11-12
    # whole-package run read 39 such reds out of 52 on this host, none of them a
    # product defect. internal/speech solved the identical problem in ticket 409
    # by staging the pinned DLLs next to the test binary and running that binary
    # from the package directory; the run step below routes cli through the same
    # shape (scripts/wisp-cli-staged-tests.sh).
    scope=(./cmd/wisp/)
    pinned=$cli_pin
    ;;
winsec)
    # Ticket 251 AC#1: winsec_pin existed with no branch to resolve it, so the pin
    # was never compared to anything this script actually ran. Glob scope + that
    # pin, both read from the declarations above; GUARD C then bites on a split,
    # a rename or a package that stopped resolving, in BOTH directions.
    scope=("$winsec_scope")
    pinned=$winsec_pin
    ;;
speech)
    # Ticket 409 / wave99: the tier claims the package the census used to report as
    # UNCLAIMED-HAS-TESTS, and the run step routes it through the staging carrier
    # (see the mode=speech branch below the guards). Scope and pin follow the same
    # one-list rule as every other tier: GUARD C bites when the package is split,
    # renamed or stops resolving, in both directions.
    scope=(./internal/speech/)
    pinned=$speech_pin
    ;;
census)
    scope=()
    pinned=''
    ;;
*)
    echo "portable-tests.sh: unknown --scope=$mode (known: ${tiers// /, })" >&2
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
    #
    # AC#4 (ticket 251) speaks BEFORE any of that is measured: this census used to
    # read a tier roster written HERE (`for m in core windows cli winsec`) while the
    # tiers that can actually be RUN are written AGAIN as the `case $mode in` branches
    # above. At the ticket's own HEAD the two sets already disagreed - census printed
    # `CLAIMED BY ... winsec` for a tier NO branch implemented, and exited 0 doing it.
    # Census now refuses to run at all unless the roster and this file's own branches
    # are equal, read back out of the bytes it is standing in.
    own="$here/portable-tests.sh"
    if [ ! -r "$own" ]; then
        echo "portable-tests.sh: census - cannot read this script's own bytes at $own, so the tier" \
            "roster cannot be compared with the case branches. Refused (ticket 251 AC#4)." >&2
        exit 1
    fi
    branches=$(awk '
        /^case \$mode in$/ { f = 1; next }
        f && /^esac$/ { exit }
        f && /^\*\)$/ { print "*"; next }
        f && /^[a-z][a-z0-9]*\)$/ { print substr($0, 1, length($0) - 1) }' "$own" | sort -u | tr '\n' ' ')
    roster=$(printf '%s\n' $tiers '*' | sort -u | tr '\n' ' ')
    if [ "$branches" != "$roster" ]; then
        {
            echo "portable-tests.sh: census - the tier roster (tiers=) and the 'case \$mode in'"
            echo "portable-tests.sh:   branches are NOT the same set, so the CLAIMED BY column below"
            echo "portable-tests.sh:   would name tiers that cannot be run, or hide tiers that can:"
            printf 'portable-tests.sh:   tiers=   : %s\n' "$roster"
            printf 'portable-tests.sh:   branches : %s\n' "$branches"
            echo "portable-tests.sh: refused rather than reported - a census whose roster drifted from"
            echo "portable-tests.sh: the runnable tiers is '配置里有一行 = 它跑过了' wearing a table"
            echo "portable-tests.sh: (ticket 251 AC#4)."
        } >&2
        exit 1
    fi
    # THE go list CALL, AND WHY IT USED TO DIE WITHOUT A WORD. This block ran
    # `all=$(go list ./... 2>/dev/null | sort -u)` under `set -eu -o pipefail`, so a
    # non-zero `go list` aborted the script with rc=1 and ZERO bytes of output -
    # measured on this host 2026-10-05 14:07 with GOOS=linux (rc=1, stdout 34 lines,
    # stderr 3 lines, log 0 lines). A red with no words is the "step green can equal
    # nothing tested" disease wearing the other colour, so the streams are kept
    # separate and a failed resolve REFUSES out loud, naming its own rc.
    census_out=$(mktemp 2>/dev/null || echo "$root/.portable-census.$$.txt")
    census_err=$(mktemp 2>/dev/null || echo "$root/.portable-census-err.$$.txt")
    go_rc=0
    go list ./... >"$census_out" 2>"$census_err" || go_rc=$?
    if [ "$go_rc" -ne 0 ]; then
        {
            echo "portable-tests.sh: census - go list ./... exited $go_rc. stdout named" \
                "$(grep -c . "$census_out" || true) package(s); its stderr said:"
            sed 's/^/portable-tests.sh:   err| /' "$census_err"
            echo "portable-tests.sh: refused rather than reported - a roster built on a partial package set"
            echo "portable-tests.sh: can only UNDER-REPORT holes, which is the exact claim ticket 111 was"
            echo "portable-tests.sh: filed against. internal/models/assembly_reachability_121_test.go:74-92"
            echo "portable-tests.sh: records that the non--e form of this query does not succeed for"
            echo "portable-tests.sh: GOOS=linux (sherpa-onnx-go-linux excludes all its files there), so the"
            echo "portable-tests.sh: census is wired on the leg where the plain form DOES succeed."
        } >&2
        rm -f "$census_out" "$census_err"
        exit 1
    fi
    rm -f "$census_err"
    all=$(grep -v '^[[:space:]]*$' "$census_out" | sort -u || true)
    rm -f "$census_out"
    n=$(printf '%s\n' "$all" | grep -c . || true)
    echo "portable-tests.sh: census GOOS=$goos - go list ./... = $n packages (one row each)"
    printf 'portable-tests.sh: %-46s %-11s %s\n' PACKAGE 'TESTS(t/x)' 'CLAIMED BY'
    empty=0
    noscope=0
    guardd=0
    unclaimed=''
    while IFS= read -r p; do
        [ -n "$p" ] || continue
        counts=$(go list -f '{{len .TestGoFiles}}/{{len .XTestGoFiles}}' "$p" 2>/dev/null || echo '?/?')
        where=''
        for m in $tiers; do
            if [ "$m" = census ]; then continue; fi # the auditor owns no pin (ticket 251 AC#4)
            case $m in
            core) pin=$core_pin ;; windows) pin=$win_pin ;; cli) pin=$cli_pin ;;
            winsec) pin=$winsec_pin ;; speech) pin=$speech_pin ;;
            *)
                {
                    echo "portable-tests.sh: census - tier '$m' is in tiers= but this loop has no pin for"
                    echo "portable-tests.sh:   it, so census could only ever report it as covering nothing."
                    echo "portable-tests.sh: Name the pin (ticket 251 AC#4), do not add a tier silently."
                } >&2
                exit 1
                ;;
            esac
            if printf '%s\n' "$pin" | grep -qxF "$p"; then where="$where$m"; fi
        done
        if [ -z "$where" ]; then
            where=' NO-SCOPE'
            noscope=$((noscope + 1))
            # GUARD D's split of the NO-SCOPE rows: `0/0` is a package with no
            # denominator on THIS platform (GUARD A's business, and a scope entry
            # for it would be the false claim ticket 111 AC#3 refuses), while
            # anything else - a non-zero count, or `?/?`/empty because the count
            # could not be read - is a package with tests that no step will ever
            # run. Default-deny: an unreadable count is counted as a hole, never
            # as covered.
            case $counts in
            0/0) ;;
            *)
                unclaimed="$unclaimed$p  tests-compiled-for-$goos=${counts:-UNKNOWN}"$'\n'
                guardd=$((guardd + 1))
                where="$where <-UNCLAIMED-HAS-TESTS"
                ;;
            esac
        fi
        case $counts in
        0/0) empty=$((empty + 1)); where="$where <-NO-TESTS" ;;
        esac
        printf 'portable-tests.sh: %-46s %-11s %s\n' "$p" "$counts" "$where"
    done <<<"$all"
    echo "portable-tests.sh: census totals: packages=$n with-zero-compiled-tests=$empty claimed-by-no-scope=$noscope unclaimed-with-tests=$guardd"
    if [ "$guardd" -ne 0 ]; then
        {
            echo "portable-tests.sh: GUARD D - $guardd package(s) compile a test file for GOOS=$goos and"
            echo "portable-tests.sh:   NO named scope claims them, so no CI step runs them and no CI step"
            echo "portable-tests.sh:   can ever go red over them:"
            printf '%s' "$unclaimed" | sed 's/^/portable-tests.sh:   /'
            echo "portable-tests.sh: this row is the ticket 111 field itself ('CI 只测 33 个包里的 20 个')."
            echo "portable-tests.sh: The census used to PRINT these rows and exit 0, so the hole was visible"
            echo "portable-tests.sh: only to whoever thought to run it by hand - which is how internal/session"
            echo "portable-tests.sh: (2 files, ticket 224) and internal/projctx (1 file, ticket 200) sat at"
            echo "portable-tests.sh: zero coverage for four days after their tests landed. Pull the package"
            echo "portable-tests.sh: into a named scope and update that tier's pin in the SAME commit, or"
            echo "portable-tests.sh: state on the ticket where its coverage lives. What this guard forbids"
            echo "portable-tests.sh: is the third option: a build tag, a t.Skip or a deleted test."
        } >&2
        exit 1
    fi
    exit 0
fi

if [ $# -gt 0 ]; then
    # Explicit paths are the local-debug form (and, until ticket 251, the way
    # scripts/winsec-tests.sh delegated). GUARD A and GUARD B have always applied to
    # whatever is named. GUARD C used to be skipped wholesale here, "because an
    # ad-hoc list has no expectation to be faithful to" - true of an ad-hoc list,
    # false of the one call CI really makes: winsec-tests.sh:94 hands over exactly
    # ./internal/winsec/, which IS a tier's own scope and does have an expectation to
    # be faithful to. So that list, and only that list, is run as its tier and pinned
    # (AC#1's option (b), on the route the gate actually takes - the tier branch above
    # alone would have left CI's real path unpinned). Any other explicit list keeps
    # today's meaning: nothing to pin, A and B still apply.
    if [ $# -eq 1 ] && { [ "$1" = "$winsec_dir" ] || [ "$1" = "$winsec_scope" ]; }; then
        mode=winsec
        scope=("$winsec_scope")
        pinned=$winsec_pin
        echo "portable-tests.sh: explicit scope [$1] IS the winsec tier's own scope - running it" \
            "as mode=$mode against winsec_pin (ticket 251 AC#1)"
    else
        scope=("$@")
        pinned=''
    fi
fi
if [ ${#scope[@]} -eq 0 ]; then
    echo "portable-tests.sh: empty scope (mode=$mode) - refusing to be a green no-op" >&2
    exit 2
fi

# ---- resolve the scope once, so all three guards read the same truth ---------
#
# TICKET 250: THE DENOMINATOR IS STDOUT ALONE. This block used to read
#
#     if ! go list "${scope[@]}" >"$resolved" 2>&1; then
#
# and then counted EVERY line of that file - `pkgcount` below, and GUARD B's
# per-package loop at the bottom of the script, both read "$resolved". `go list`
# writes package paths to stdout and progress/diagnostics to stderr, so on a cold
# module cache the `go: downloading ...` lines were counted as packages while the
# package set itself was untouched. Measured on this host 2026-10-02 08:5x with
# GOMODCACHE pointed at an empty directory and nothing else changed:
#
#     portable-tests.sh:   set than the one pinned next to it. Pinned: 25, resolved: 36.
#
# - eleven `go: downloading` lines, exit 1, twelve seconds in, BEFORE `go test`
# ever started, so the whole core scope's readings went missing. A warm cache prints
# nothing to stderr (measured: stdout 25 lines / stderr 0 bytes at this same HEAD),
# which is why it is silent on every developer laptop and red only on CI.
#
# Keeping stderr out of the denominator is not the same as throwing it away: a
# successful `go list` echoes its stderr into the step log verbatim, and a FAILED
# one still exits 1 with BOTH streams printed - the pre-fix failure branch (it used
# to `cat` the one merged file) is preserved, not softened (ticket 250 AC#1's
# reverse face).
#
# The standing carrier that judges this block both ways on a laptop, ten cases, no
# real test binary involved, is scripts/portable-tests-selftest.sh (ticket 250
# AC#3). Run it before editing anything between here and GUARD C below.
resolved=$(mktemp 2>/dev/null || echo "$root/.portable-resolved.$$.txt")
resolved_err=$(mktemp 2>/dev/null || echo "$root/.portable-resolved-err.$$.txt")
if ! go list "${scope[@]}" >"$resolved" 2>"$resolved_err"; then
    {
        echo "portable-tests.sh: go list [${scope[*]}] exited non-zero. Its stdout was:"
        sed 's/^/portable-tests.sh:   out| /' "$resolved"
        echo "portable-tests.sh: its stderr was:"
        sed 's/^/portable-tests.sh:   err| /' "$resolved_err"
        echo "portable-tests.sh: go list [${scope[*]}] failed - the scope cannot be audited" \
            "against a pattern that does not resolve."
    } >&2
    rm -f "$resolved" "$resolved_err"
    exit 1
fi
if [ -s "$resolved_err" ]; then
    echo "portable-tests.sh: go list wrote $(grep -c . "$resolved_err" || true) line(s) to stderr" \
        "(progress, not packages, not in the denominator):"
    sed 's/^/portable-tests.sh:   err| /' "$resolved_err"
fi
rm -f "$resolved_err"
grep -v '^$' "$resolved" | sort -u >"$resolved.sorted" || true
mv "$resolved.sorted" "$resolved"
pkgcount=$(wc -l <"$resolved" | tr -d '[:space:]')

# Ticket 250 AC#2: the shape check, in capability form, NOT a word list.
#
# Every line that made it into the denominator has to BE an import path. Asserting
# "this is not the text `go: downloading`" would be the wrong instrument - Go rewords
# its progress output and the guard goes blind again - so what is asserted is the
# grammar a path has (go/src/internal/module check.go: an element is a non-empty run
# of letters, digits and - . _ + ~, starting alphanumeric, with ! for the case
# escape; elements joined by /). No whitespace anywhere in the line is the part that
# progress text, diagnostics and "matched no packages" prose all fail and no import
# path can.
#
# A line that fails is REFUSED, not filtered. Dropping the unrecognised lines and
# carrying on would make GUARD B audit a smaller scope than the pin claims while
# printing the same green - the exact shape the comment below stands guard over, so
# the whole step goes red and says which lines it choked on.
import_path_re='^[A-Za-z0-9][A-Za-z0-9._+~!-]*(/[A-Za-z0-9][A-Za-z0-9._+~!-]*)*$'
notpaths=''
notpathcount=0
while IFS= read -r line || [ -n "$line" ]; do
    [[ $line =~ $import_path_re ]] && continue
    notpaths="$notpaths$line"$'\n'
    notpathcount=$((notpathcount + 1))
done <"$resolved"
if [ "$notpathcount" -ne 0 ]; then
    {
        echo "portable-tests.sh: GUARD C - go list's stdout carried $notpathcount line(s) that are not"
        echo "portable-tests.sh:   import paths, so this $pkgcount-line denominator cannot be trusted:"
        printf '%s' "$notpaths" | sed 's/^/portable-tests.sh:   /'
        echo "portable-tests.sh: refused rather than filtered: a guard that drops the lines it does not"
        echo "portable-tests.sh: recognise measures less than it claims to (ticket 250 AC#2)."
    } >&2
    rm -f "$resolved"
    exit 1
fi

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
    "TestAppendSessionTurnCrashChild|./internal/memory/|any|reexec|child-process entry point of TestAppendSessionTurnCrashAtomicAcrossKill (dao_session_append_326_test.go:44, ticket 326 AC#3); not a test of its own - it skips unless re-exec'd with WISP_326_APPEND_CRASH_CHILD=1, and the parent that drives it (exec of os.Args[0] + TerminateProcess mid-write) runs and passed on BOTH platforms' steps, so the assertion itself is not lost, only the child role is un-runnable as a top-level case. Same class as the TestSubprocessCrashWriter row above"
    "TestHelperProcess|./internal/proc/|windows|reexec|child-process entry point of the job-object cases (jobscope_windows_test.go:87, //go:build windows); skips unless re-exec'd with the helper env var set"
    "TestLiveWasapiSmoke|./internal/audio/|windows|fixture|live WASAPI capture needs WISP_LIVE_MIC=1 and a physical microphone (hotplug_test.go:527, //go:build windows); a hosted runner has no audio endpoint, so the case has no subject there. Real-hardware smoke is ticket 16's"
    "TestAC297LivePacketByteFace|./internal/audio/|windows|fixture|ticket 297 AC#0's real-device half: one in-process run against THIS machine's default capture endpoint that prints the mix-format quad, the first raw packet's samples read as int16 and again as float32, and the window levels (livepacket_297_windows_test.go:41, //go:build windows), gated behind WISP_LIVE_MIC=1 exactly like TestLiveWasapiSmoke above; a hosted runner has no microphone, so the case has no subject there, and one unaccounted top-level SKIP is exactly what turned the windows step red when ticket 297-l2 landed the file (roster grep 297 was 0). The reading the ticket asked for was taken on the orchestrator's machine and lives under .scratch/wisp/probes/297/l1/; real-hardware reruns stay local, same rule as ticket 16's smoke row above"
    "TestRealDownloadVadThroughPipeline|./internal/models/|any|opt-in|real-network spot check gated by WISP_IT_REAL_MIRROR=1 (manifest_real_test.go:120); CI has no business pulling model archives from a third-party mirror"
    "TestRealDownloadPuncArchiveThroughPipeline|./internal/models/|any|opt-in|real-network spot check gated by WISP_IT_REAL_MIRROR=1 (manifest_real_test.go:147); same reason as the VAD row"
    "TestSyncRegistryProbeLive|./internal/risk/|windows|fixture|P12 live registry evidence; on ubuntu it is now out of scope by platform instead of by skip (ticket 93 AC#2 moved it verbatim into syncdirs_windows_test.go), and on a Windows host whose HKCU Accounts key carries no UserFolder the hive has nothing to say. Measured on a windows host at HEAD 84e43af: --- SKIP at syncdirs_windows_test.go:133. As of ticket 110 AC#4 the windows leg carries ./internal/risk/ in its scope, so this row is re-verified against the compiled WINDOWS test binary (go test -list) and printed with its reason in that step's log - which is the step-level answer R-93-4 asked for. Remedy for a real run: a host with a sync record, or fold the shape check into the fixture-driven case. next= on ticket 93"
    "TestWorkspaceSwitchRefusesAJunctionToOutside|./internal/tools/|linux|fixture|ticket111 AC#10, re-based by ticket 347 after a MEASURED correction. The case is AC#3(ii) with a REAL junction: on windows it builds one with cmd /c mklink /J and asserts the switch is refused (runs and passes here). What this row said until 347 - 'on POSIX there is nothing to refuse, because the detection is a Windows implementation (reparseComponents returns nil elsewhere)' - is now half false, and the half that is false is the one a reader would act on: (1) the DETECTION is no longer Windows-only, because ticket 347 landed risk.pathresolver_other.go's reparseComponents as the lstat symlink walk, so a symlink is denied on POSIX exactly as a junction is denied on Windows; (2) the case nevertheless STILL SKIPS on linux at paths_workspace_test.go:198, re-measured after 347, because that skip is a platform guard (runtime.GOOS != windows) and not a capability probe - no test starts running here, and the row accounts for exactly what it always did. What is stale is the sentence inside the skip, not the skip. OWNER of the POSIX half: ticket 347 landed the detection a POSIX equivalent needs, but the conversion - build a SYMLINK instead of a junction, and rewrite that message - lives in internal/tools/paths_workspace_test.go, a file ticket 347's write surface forbids, so the POSIX-equivalent criterion is NOT faked here and the conversion is dispatched separately. Re-verified against the compiled LINUX test binary every run: delete the skip and the case starts running here; delete the case and GUARD's staleness check goes red."
    "TestD34WriteMatrix|./internal/tools/|linux|fixture|the table needs a second volume from otherVolumeDir (fs_write_test.go:140/151), and the volume-identity call it uses is a Windows API - on POSIX there is no object to ask. In a CI ubuntu container /tmp and / are one filesystem anyway. It RUNS AND PASSES on a two-volume windows host, so if a windows host ever has a single volume this entry does not cover it and the step goes red on purpose"
    "TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop|./internal/tools/|linux|fixture|same otherVolumeDir requirement as the write matrix row above (fs_write_test.go:677); on POSIX there is no volume identity to compare, so the case has no subject. Runs and passes on a two-volume windows host"
    "TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath174r3|./internal/tools/|linux|fixture|ticket174 AC#6 shape 甲 (the ledger row, not the POSIX-equivalent shape 乙). The case is AC#2d's junction arm: the filed path must be a REAL NTFS junction (cmd /c mklink /J) whose link spelling C26 refuses - per the fixture comment (task_output_pointer_notice_test.go:376-380) a non-exempt reparse traversal is the ONLY reachable input that gives this branch a path literal at all. On POSIX nothing can be refused, because the detection it exercises is a Windows implementation (risk/pathresolver_other.go reparseComponents returns nil elsewhere - the same wall the workspace-junction row above names), so this machine has no subject: task_output_pointer_notice_test.go:408 skips with 前置条件缺失：造不出真 NTFS junction（mklink /J %s %s: %v: %s）, and the %v the ubuntu leg filled in was exec: \"cmd\": executable file not found in \$PATH: (archived run-37166458550-failed.log, the file then at :428 - the line moves as other criteria land, so the row quotes the skip rather than the line). It RUNS AND PASSES on windows (real mklink /J, no skip), so the assertion is not lost, it is only un-evaluable here. Registered rather than deleted because runtests.sh makes every SKIP fatal (ticket 71 AC#3) and one unaccounted skip is exactly what turned test-core step7 red in the archived run (four numbers there: PASS=1010 FAIL=4 SKIP=1 - the FAILs have their own owners, the SKIP=1 is what this row accounts for): -skip filters an accounted name silently, so this is the repo's chosen mechanism for declaring the skip and re-checking it every run instead of letting it land as an unnamed red. OWNER of the POSIX half: whoever lands AC#6 shape 乙 on ticket 174, a POSIX-runnable equivalent of the same criterion (the notice may not relay the filed path, no junction needed); this row does not fake that half. Re-verified against the compiled LINUX test binary every run: delete the skip and the case starts running here; delete the case and the staleness check goes red. next= on ticket 174"
    "TestReadFailureDoesNotAdoptStatTheNextTickSeesTheWrite|./internal/config/|windows|fixture|ticket 70 AC#2 triage: the case's mechanism contract (CreateFile dwShareMode=0 conflicts with a concurrent reader) is not delivered on the hosted windows runner session (red with 'must surface an error, got report &{Hot:[ball] ...}' on runs 43362781437 and 43243081268, unchanged since ci-r1's 2026-11-01 reading); the PROBE-UNDO gate in manager_readfail_42r1_envprobe_windows_test.go skips it there with the undo condition named, and it RUNS AND PASSES on every host where the probe passes (measured 21x green locally, incl. GOTOOLCHAIN=go1.27.4). Re-enable on a runner image where share-0 conflicts again: delete the probe gate in TestReadFailureDoesNotAdoptStatTheNextTickSeesTheWrite"
    "TestTicket223PermissionDeniedSitsInItsOwnSentence|./cmd/wisp/|windows|fixture|ticket 223 AC#4's permission sentence needs a REAL read denial (icacls /deny Everyone:(R)) to plant its premise, and the hosted windows runner session (runneradmin, first red run 43413718734 test-windows step 7) holds an administrator token whose Effective Access ignores a deny ACE on the Everyone group - icacls lands as a no-op, os.ReadFile still succeeds, and the case's own positive control proves the premise absent. Until 2026-11-13 that shape was t.Fatalf('attempt %d: icacls denied nothing...'), i.e. a missing SUBJECT was booked as a broken assertion; the wave-8b leg turned it into a named skip (config_reload_perm_223_windows_test.go) with the reason printed, and it RUNS AND PASSES - all assertions intact - on any host where the deny actually binds (normal-user sessions, the dev box). Re-enable everywhere: a hosted runner image running the step as a non-admin, or a UAC-split-token shape the case could detect; delete the skip branch to return the Fatalf. Accounted by ticket 111 wave 8b (ledger row requested by the dispatch; ticket 71 AC#3 rule that a SKIP must carry a ledger row)"
    "Test16r1VoiceDemoWakeAsrAgentSpeakSettleRealHardware|./cmd/wisp/|windows|fixture|ticket 16 AC#1's demo speaks its reply through the PRODUCTION announce path (real waveOut render, voice_out_windows.go:217); a hosted runner reports zero waveOut output devices, so the D43#15/#26 Speaking walk has no subject there and the announce line degrades to the D10 notification (resident_voice_engines_windows.go:1171) - CI red on run 43425302222 with 'no waveOut output device' + 'timed out waiting for the wav to drain (4.wav)'. The wave-10 leg added the precondition probe demo16AnnounceNeedsAnOutputDevice (voice_e2e_demo_16r1_windows_test.go, openWaveOutSpeaker(22050) per the 26p rig precedent) which SKIPS with that reason only when no device answers; on any machine with an output device the probe closes immediately and every assertion (the whole D43 walk, speak counters, settle trace) runs unchanged - measured PASS on the dev box 2026-11-13. Re-enable everywhere: a runner with any audio endpoint (real or virtual) - delete the probe call to return the unconditional run. Accounted by ticket 111 wave 10 (ledger row requested by the dispatch; ticket 71 AC#3)"
    "Test41a2r1WarmGateFeedsKwsNotTheAsr|./cmd/wisp/|windows|fixture|the Warm-gate case measures its turns between KNOWN drain positions of the C8 seam (voice41a2r1Feed), and the seam is a 6-frame/192ms bounded channel whose meter DROPS on overflow and counts it (internal/audio/audio.go:115-145, the 'bounded channel full (slow consumer)' WARN). On a loaded host the ASR hand-off spike (CI run 43425302222: 24.03s, dropped_frames_total=1, kws_hits=1 asr_fed=13312 - the behaviour landed) overruns that window, frames are gone, and the drained counter can never reach the asserted position: the drain timeout is constructively unmeetable, not a behaviour red. The wave-10 leg turned the drain timeout into a verdict helper (voice41a2r1DrainTimeout, resident_voice_41a2r1_windows_test.go): injector Stats() frames_dropped==0 stays a real red (consumer stalled without a seam drop), frames_dropped>0 skips NAMED with the sent/dropped/drained readings. Preconditions in place: it RUNS AND PASSES with all assertions intact (measured on the dev box before the move 2026-11-13). Re-enable everywhere: a host fast enough that the hand-off never overruns 192ms of buffered audio (or an off-thread seam redesign, owner internal/audio). Accounted by ticket 111 wave 10 (ledger row requested by the dispatch; ticket 71 AC#3)"
    "TestPanelHostRealWindowHopAndLifecycle|./cmd/wisp/|windows|fixture|D32's panel latency VERDICT (hot re-show <= 200ms, PLAN.md:2345/:2358 frozen - not relaxed, not retried) moved to the winlive tier (panel_host_windows_budget_winlive_test.go) by ticket 111 wave 10, following the 33-r5 precedent (the 2s exit clause): the hosted runner's desktop composition session read hot re-show 220.7ms (CI run 43425302222) while the SAME binary on the dev box reads 42.1ms and cold sits inside the 1500ms budget on BOTH hosts - a same-code 5x divergence at one hop is the runner's session, not the product (D32 16.3.1's own metering doctrine: 先定口径再定数字). This row accounts for the shape that remains visible in CI: the DEFAULT-tier case still measures and logs hot/cold every run and feeds the aggregate, but the over-budget t.Errorf lives only where it is decidable - on a desktop box (winlive has NO CI job; the cost is named in panel_host_windows_budget_winlive_test.go and docs/evidence/s1/111-cmdwisp-last-mile.md). Every one-sample-decidable dimension stays asserted in CI: same HWND, single-window host-pid identity, no listen sockets, Destroy leaves no window. Re-enable in the default tier: a CI runner whose interactive desktop session measures this hop like real hardware (e.g. a self-hosted windows runner with a real session) - move the t.Errorf back from the winlive file. Accounted by ticket 111 wave 10 (ticket 71 AC#3: a skipped VERDICT must be accounted; the case itself still RUNS in CI)"
    "TestPanelHostLatencyPercentilesAC2|./cmd/wisp/|windows|fixture|the P50/P95 aggregate of the SAME D32 budgets (cold 1500 / hot 200) moved to the winlive tier with the lifecycle verdict (panel_host_windows_budget_winlive_test.go), ticket 111 wave 10 - same reason, same precedent, budgets unchanged (see the row above; CI run 43425302222 red on the runner's 220.7ms hot sample). The default-tier case still prints the full percentile table every run (the reading) and keeps its empty-aggregate named skip; the over-budget t.Errorf now lives only on a desktop box. Re-enable in the default tier with the lifecycle row. Accounted by ticket 111 wave 10 (ticket 71 AC#3)"
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
            echo "portable-tests.sh: ticket 311 AC#1 shape B: INCONCLUSIVE - no volume can be read from $wisp_311_name=$wisp_311_value (reference %LOCALAPPDATA%=$wisp_311_ref_path), so nothing is judged here and this step runs on the root as given. This is NOT the check passing: registered as debt Q-311-r1-1 on ticket 311." >&2
            continue
        fi
        if [ "$wisp_311_vol" != "$wisp_311_ref_vol" ]; then
            wisp_311_cross="$wisp_311_cross $wisp_311_name=$wisp_311_value(volume $wisp_311_vol)"
        fi
    done
    if [ -n "$wisp_311_cross" ]; then
        echo "portable-tests.sh: ticket 311 AC#1 shape B: GUARD - the test temp root names ANOTHER VOLUME than %LOCALAPPDATA%: reference volume $wisp_311_ref_vol (path $wisp_311_ref_path), offender:$wisp_311_cross" >&2
        echo "portable-tests.sh: ticket 311 AC#1 shape B: this gate's COLOUR is a function of that volume - the same unchanged internal/winsec code reads PASS=58 FAIL=0 on a volume-$wisp_311_ref_vol root and PASS=56 FAIL=2 with Go's own TempDir cleanup failing on a D:/tmp root (ticket 311). Refusing to run is the fix; reporting a re-coloured gate as a verdict is not." >&2
        echo "portable-tests.sh: ticket 311 AC#1 shape B: point WISP_TEST_GOTMP_ROOT (or GOTMPDIR) at a directory on volume $wisp_311_ref_vol, or unset both, then run this step again. Do NOT relax an assertion, threshold, golden or SLO to make it green." >&2
        exit 2
    fi
fi

# ---- TICKET 314 AC#3 / ledger A976: the cli tier's END ANCHOR, computed BEFORE the run.
#
# WHY THIS EXISTS. cmd/wisp's whole-package run is the one run in this repository
# that can stop before its last case and still leave a plausible-looking verdict:
# an external kill, a -test.timeout panic, or a hang the operator ends by hand all
# leave the four numbers present and the package result line present, so every
# pre-existing judgement here (GUARD B's `ok <pkg>` line, the four numbers, the
# carrier's own ok/FAIL line) reads "a run happened". Measured on 2026-11-13: a
# full run reached === RUN=380 of a 619-name roster, ended on
# `panic: test timed out after 10m0s`, and booked FAIL=2 - two real reds - while
# ~239 cases, including the last enumerated one, never ran. Under A976 those four
# numbers are NOT a coverage denominator, and nothing in this script could say so.
#
# WHAT IS ASSERTED. The roster of top-level cases is enumerated HERE, from the
# compiled test binary for THIS platform, with the same `go test -list` shape the
# ledger check above already relies on - so this adds no environment dependency:
# if the binary cannot be listed, the ledger check has already refused. The LAST
# name in that enumeration is the run's end anchor. Go registers top-level cases
# in file order and runs them in that same order, so the last name printed by
# -list is the last case the run can reach. The run must reach it AND print its
# verdict line (`--- PASS|FAIL|SKIP: <name>`). Seeing the `=== RUN` line alone is
# not enough: a process that dies inside the last case prints that too, and telling
# "stopped early" from "finished" is the entire point.
#
# WHAT IT DOES NOT DO. It does not assert ran == roster: subtests legitimately add
# `=== RUN` lines, the carrier's named fixture fallback (-test.skip, ticket 111
# wave 8b) legitimately removes cases, and accounted skips legitimately remove
# others. Both numbers are REPORTED, never equated. The anchor is the binary
# verdict. This touches no test, no threshold, no golden, no ledger row and no rc
# that was not already non-zero.
cli_anchor_name=''
cli_roster_count=0
if [ "$mode" = cli ]; then
    cli_roster=$(mktemp 2>/dev/null || echo "$root/.portable-cli-roster.$$.txt")
    cli_roster_err=$(mktemp 2>/dev/null || echo "$root/.portable-cli-roster-err.$$.txt")
    if ! go test -list '.*' "${scope[@]}" >"$cli_roster" 2>"$cli_roster_err"; then
        {
            echo "portable-tests.sh: END ANCHOR (ticket 314) - \`go test -list '.*' ${scope[*]}\` exited non-zero,"
            echo "portable-tests.sh:   so the cli tier's roster cannot be enumerated and there is no end anchor to"
            echo "portable-tests.sh:   assert against. Refusing to run: a guard that cannot name its endpoint is"
            echo "portable-tests.sh:   the green no-op this script exists to prevent. Its stderr was:"
            sed 's/^/portable-tests.sh:   err| /' "$cli_roster_err"
        } >&2
        rm -f "$cli_roster" "$cli_roster_err"
        exit 1
    fi
    # `go test -list` writes the case names AND then Go's own `PASS` /
    # `ok  <pkg> <dur>` trailer to stdout; keep only bare Test identifiers, so
    # the anchor is always a case name and never the package result line.
    grep -E '^Test[A-Za-z0-9_]+$' "$cli_roster" >"${cli_roster}.names" || true
    cli_roster_count=$(grep -c . "${cli_roster}.names" || true)
    cli_anchor_name=$(grep . "${cli_roster}.names" | tail -n 1 || true)
    rm -f "$cli_roster" "$cli_roster_err" "${cli_roster}.names"
    if [ "$cli_roster_count" -eq 0 ] || [ -z "$cli_anchor_name" ]; then
        echo "portable-tests.sh: END ANCHOR (ticket 314) - the roster came back empty ($cli_roster_count names)," \
            "so there is no last case to anchor on. Refusing to run." >&2
        exit 2
    fi
    case $cli_anchor_name in
    Test[A-Za-z0-9_]*) ;;
    *)
        echo "portable-tests.sh: END ANCHOR (ticket 314) - the last enumerated name '$cli_anchor_name' is not a" \
            "bare top-level Test identifier, so it cannot be anchored on. Refusing to run." >&2
        exit 2
        ;;
    esac
    echo "portable-tests.sh: END ANCHOR (ticket 314): roster=$cli_roster_count top-level case(s) from" \
        "\`go test -list '.*' ${scope[*]}\`; the run must reach the last one and print its verdict:" \
        "$cli_anchor_name"
fi

capture=$(mktemp 2>/dev/null || echo "$root/.portable-tests.$$.log")
rc=0
# tee costs nothing and keeps the strict runner's verdict intact while giving
# this script the bytes for the four numbers and for naming an offender.
#
# TICKET 409 / wave99: the speech tier never reaches the strict runner. A bare
# `go test ./internal/speech/` cannot be judged here - the cgo test binary's DLL
# tree is resolved by the Windows loader at process load time, before any Go code,
# so on a host whose System32 carries an older onnxruntime the run is 47 skips
# (which the strict runner books as fatal), and without the DLLs on the first
# search slot it dies with 0xc0000135 before a single test runs. The tier's scope
# is still resolved, pinned (GUARD C) and audited (GUARD A) above; only the run
# step is delegated, to the carrier that stages the pinned DLLs next to the test
# binary and runs it from the package directory. Its stdout is the same -v stream
# the strict runner would have produced, so the four numbers below are read off
# it unchanged.
if [ "$mode" = speech ]; then
    carrier="$here/speech-tests.sh"
    if [ ! -f "$carrier" ]; then
        echo "portable-tests.sh: $carrier is missing - the speech tier routes through it and" \
            "refuses to invent a weaker run in its place." >&2
        rm -f "$capture"
        exit 2
    fi
    echo "portable-tests.sh: mode=speech routes [${scope[*]}] through the staging carrier" \
        "scripts/speech-tests.sh (ticket 409) - the strict runner cannot judge this package:"
    echo "portable-tests.sh: the loader resolves the sherpa DLL tree at process load time," \
        "so the carrier stages the pinned DLLs next to the test binary first."
    bash "$carrier" 2>&1 | tee "$capture" || rc=$?
elif [ "$mode" = cli ]; then
    # Same reasoning as the speech branch above: cmd/wisp's test binary links the
    # sherpa import library too, so its second-level onnxruntime dependency is
    # resolved by the loader BEFORE any Go code runs. PATH cannot influence that
    # ordering (System32 is searched before PATH), which is exactly what ticket
    # 409 established for internal/speech. The carrier stages the pinned DLLs
    # next to the binary and runs it from the package directory; its stdout is
    # the same -v stream, so the four numbers below are read off it unchanged.
    # This does NOT relax the ticket-98 GUARD, the SKIP rule or any assertion -
    # only HOW the binary is launched changes.
    carrier="$here/wisp-cli-staged-tests.sh"
    if [ ! -f "$carrier" ]; then
        echo "portable-tests.sh: $carrier is missing - the cli tier routes through it and" \
            "refuses to invent a weaker run in its place." >&2
        rm -f "$capture"
        exit 2
    fi
    echo "portable-tests.sh: mode=cli routes [${scope[*]}] through the staging carrier" \
        "scripts/wisp-cli-staged-tests.sh - the strict runner cannot judge this package:"
    echo "portable-tests.sh: the loader resolves the sherpa DLL tree at process load time and" \
        "System32 is searched before PATH, so the carrier stages the pinned DLLs next to the" \
        "test binary first (ticket 409's shape, extended to cmd/wisp)."
    bash "$carrier" 2>&1 | tee "$capture" || rc=$?
else
    sh "$strict" "${scope[@]}" -count=1 -skip "$skip_pattern" 2>&1 | tee "$capture" || rc=$?
fi

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

if [ "$mode" = speech ]; then
    # GUARD B (speech shape, ticket 409 / wave99): the carrier runs the raw test
    # binary, so Go's per-package `ok <import path> <duration>` result line does
    # not exist in the output and the anchored per-package table below cannot
    # apply. The speech-shaped equivalent of "the package was tested" is the -v
    # denominator itself: a capture with zero top-level PASS and zero top-level
    # FAIL lines ran nothing at all - the green no-op the strict runner refuses
    # as its third check - so it is refused here with the same voice. The
    # remaining skips are NOT made fatal here, deliberately and on the record:
    # the carrier's own exit code already carries every assertion verdict, and
    # the 6 residual skips name a missing punctuation model (a fixture-class
    # environment gap, printed with their reasons above), not a softened
    # assertion - the strict runner's skip rule stays untouched for every other
    # tier.
    if [ "$passed" -eq 0 ] && [ "$failed" -eq 0 ]; then
        {
            echo "portable-tests.sh: GUARD B (speech) - the carrier's -v output carries NO top-level"
            echo "portable-tests.sh:   result line, so this step tested [${scope[*]}] ZERO times:"
            echo "portable-tests.sh:   === RUN=$ran  --- PASS=$passed  --- FAIL=$failed  --- SKIP=$skipped"
            echo "portable-tests.sh: that is the green no-op shape ticket 71 AC#3 was written against."
        } >&2
        if [ "$rc" -eq 0 ]; then rc=1; fi
    fi
else
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
fi

# ---- TICKET 314 AC#3 / ledger A976: the cli tier's END ANCHOR verdict ----------------
#
# Placed AFTER GUARD B on purpose: GUARD B asks "did the package print its own
# result line", which an early-stopped run answers YES (the carrier prints `ok` or
# `FAIL` from the binary's real exit code and measured duration either way). This
# block asks the only question that separates a completed suite from a cut one:
# did the LAST enumerated case actually reach its verdict? It reports both counts,
# the missing anchor and the last line of output, and it separates "killed by an
# external clock/signal" from "an assertion failed" - those two have different
# meanings and must not be filed under the same word.
if [ "$mode" = cli ]; then
    cli_anchor_re=$(printf '%s' "$cli_anchor_name" | sed 's/[][\\^$.*+?(){}|]/\\&/g')
    if ! grep -Eq "^--- (PASS|FAIL|SKIP): ${cli_anchor_re}( |\$)" "$capture"; then
        cli_last_line=$(tail -n 1 "$capture" || true)
        # Three different things must not share one word: an external clock
        # (timeout(1) / CI), a signal or native crash that took the whole binary
        # down (bash spells a signal death 128+n; on this host an unhandled
        # 0xC0000005 reaches MSYS2 as rc=139, which both
        # cmd/wisp/panel_host_windows.go:581-586 and
        # docs/evidence/s1/33-panel-host-c27-r7.md:99 name "the whole process
        # down"), and Go's own -test.timeout panic. None of them is a test
        # failure, and the first two mean the run was cut, not red.
        cli_killed=''
        case $rc in
        124) cli_killed='rc=124 (external timeout(1))' ;;
        137) cli_killed='rc=137 (SIGKILL)' ;;
        139) cli_killed='rc=139 (SIGSEGV / 0xC0000005 - whole process down)' ;;
        *)
            if [ "$rc" -ge 129 ] && [ "$rc" -le 143 ]; then
                cli_killed="rc=$rc (128+signal $((rc - 128)) - killed by a signal)"
            fi
            ;;
        esac
        if [ -z "$cli_killed" ] && grep -q 'panic: test timed out' "$capture"; then
            cli_killed='panic: test timed out (Go -test.timeout)'
        fi
        {
            echo "portable-tests.sh: END ANCHOR (ticket 314) - the cli run did NOT reach its last enumerated case."
            echo "portable-tests.sh:   anchor    : a '--- PASS|FAIL|SKIP: $cli_anchor_name' verdict line"
            echo "portable-tests.sh:   missing   : the anchor above never appeared in the -v output"
            echo "portable-tests.sh:   ran       : === RUN=$ran  --- PASS=$passed  --- FAIL=$failed  --- SKIP=$skipped"
            echo "portable-tests.sh:   roster    : $cli_roster_count top-level case(s) enumerated before the run"
            echo "portable-tests.sh:   last line : ${cli_last_line:-<capture is empty>}"
            if [ -n "$cli_killed" ]; then
                echo "portable-tests.sh:   verdict   : KILLED, NOT FAILED - the run was cut by an external clock or signal"
                echo "portable-tests.sh:               ($cli_killed). This is NOT a test failure and must not be reported as one;"
                echo "portable-tests.sh:               it is the early-stop/hang shape ticket 314 is about."
            else
                echo "portable-tests.sh:   verdict   : EARLY STOP / HANG - the run ended before its last case and no external"
                echo "portable-tests.sh:               kill signature was found. The four numbers above are NOT a coverage"
                echo "portable-tests.sh:               denominator (ticket 314 / ledger A976)."
            fi
            echo "portable-tests.sh:   Do NOT read this run's four numbers as a denominator, and do NOT report a killed run"
            echo "portable-tests.sh:   as a test failure. Re-run under the same conditions; if it stops again, the stopping"
            echo "portable-tests.sh:   case is the last '=== RUN' name printed above."
        } >&2
        if [ "$rc" -eq 0 ]; then rc=1; fi
    fi
fi

rm -f "$capture" "$resolved"
[ "$rc" -eq 0 ] || echo "portable-tests.sh: strict runner exited $rc for scope=[${scope[*]}]" >&2
exit "$rc"
