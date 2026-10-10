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
github.com/CarlosShao/wisp/internal/projctx
github.com/CarlosShao/wisp/internal/risk
github.com/CarlosShao/wisp/internal/secret
github.com/CarlosShao/wisp/internal/session
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

# TICKET 251 AC#4 - ONE ROSTER, READ BY EVERYONE. The tier names used to be written
# TWICE in this file: as the `case $mode in` branches, and again as
# `for m in core windows cli winsec` inside the census loop. Two lists that can
# disagree, and at the ticket's own HEAD they DID (census printed `CLAIMED BY ...
# winsec` for a tier no branch implemented) while census exited 0. tiers= is now the
# only place the names are written; the census step re-reads THIS FILE's own
# `case $mode in` block and refuses to report anything unless the two sets are equal
# in both directions, so the roster cannot be maintained apart from the branches.
tiers='core windows cli winsec census'
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
        ./internal/agent/ ./internal/agent/approval/
        ./internal/llm/... ./internal/config/...
        ./internal/memory/... ./internal/observe/... ./internal/secret/...
        ./internal/risk/... ./internal/statemachine/...
        ./internal/tools/... ./internal/models/...
        ./internal/buildinfo/... ./internal/audio/... ./internal/proc/...
        ./internal/panel/... ./internal/ball/ ./internal/perm/
        ./internal/plugin/ ./cmd/llmrecord/ "$winsec_scope"
        ./internal/session/ ./internal/projctx/
    )
    pinned=$core_pin
    ;;
windows)
    scope=(
        ./internal/proc/ ./internal/secret/ ./internal/config/ ./internal/risk/
        ./internal/ball/ ./internal/perm/ ./internal/plugin/ ./cmd/llmrecord/
        ./internal/session/ ./internal/projctx/ ./internal/audio/
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
winsec)
    # Ticket 251 AC#1: winsec_pin existed with no branch to resolve it, so the pin
    # was never compared to anything this script actually ran. Glob scope + that
    # pin, both read from the declarations above; GUARD C then bites on a split,
    # a rename or a package that stopped resolving, in BOTH directions.
    scope=("$winsec_scope")
    pinned=$winsec_pin
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
            winsec) pin=$winsec_pin ;;
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
