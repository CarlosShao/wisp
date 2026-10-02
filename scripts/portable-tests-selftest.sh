#!/usr/bin/env bash
# scripts/portable-tests-selftest.sh - the standing carrier for ticket 250: it
# judges the DENOMINATOR of scripts/portable-tests.sh (the line set GUARD C pins
# the named scopes against and GUARD B loops over), both directions, on a laptop.
#
# WHY IT EXISTS. The defect was: `go list` progress text on stderr was merged into
# the file whose lines get counted as packages, so a cold module cache read
# "Pinned: 25, resolved: 36" and killed the core step BEFORE `go test` ran - while
# on a warm cache (every developer laptop) stderr is empty and the same code is
# silent. AC#3 forbids shipping that as "only measurable on CI", so the shape is
# reproduced here on demand: 27 script invocations over 22 named scenarios (ticket 251
# took it to eighteen, ticket 254's four scenarios add nine more), each of which either
# has to go red or has to stay green, and the carrier exits non-zero if any of them
# does the other.
#
# TICKET 251 ADDED THE LAST EIGHT, and they judge a different hole: `winsec_pin`
# was a pin NOTHING compared - no `winsec)` branch existed, and the route CI really
# takes (scripts/winsec-tests.sh:94) hands scripts/portable-tests.sh an EXPLICIT
# package path, which landed in GUARD C's "nothing to pin" branch. So the same
# carrier now also shoots the winsec tier by name, the explicit path winsec-tests.sh
# passes, and the census tier roster, with the ticket's two seed shapes on each
# (stdout gains a package the pin does not have / the pinned package disappears).
# Case 18 is the counterfactual that makes the choice provable: the same seed on the
# same pin, but with the tier scope left in the caller's DIRECTORY form, stays GREEN,
# because a real `go list ./internal/winsec/` names the one package in that directory
# and can never see a second one created beside it. Only the GLOB form of the tier's
# scope puts a split into the resolved set, so the glob is load-bearing - and that is
# a reading in this carrier, not an argument in a comment.
# The fake go resolves its arguments the way go does when FAKEGO_SCOPE_FILTER is set
# (scripts/testdata/portable-tests/go, added by this ticket), which is what lets case
# 18 tell the two forms apart at all.
# Cases 11-18 FAIL on the pre-fix bytes by design - that is this ticket's 改前必红
# reading:  bash scripts/portable-tests-selftest.sh all /tmp/prefix.sh
#
# TICKET 254 ADDED THE LAST FOUR, and they judge a third thing: whether the audit is
# still happening ON THE ROUTE CI TAKES. ci.yml:439 never says --scope=winsec - it runs
# scripts/winsec-tests.sh, which hands scripts/portable-tests.sh an explicit package
# path, so every claim about the winsec tier has to be re-taken through that parent.
# Cases 19-20 drive the real chain (winsec-tests.sh -> portable-tests.sh) on the fake
# go and show the delegated audit both bites (a split seen through the parent, case 19)
# and is noticed when it stops (a widened explicit list silently leaves GUARD C's
# denominator, case 20 - that is winsec-tests.sh GUARD 3). Cases 21-22 pin AC#2: the
# core tier and the winsec tier now resolve internal/winsec through ONE variable, so on
# one seed they must refuse over the SAME package (21), and the directory form this
# ticket replaced is shown to hide that package outright (22, mutated copy). 20-22 go
# red on the pre-delta bytes; 19 was already true under ticket 251 and stays green
# there, which is why it is not the 改前必红 shot:
#   git show d253703a^:scripts/portable-tests.sh >/tmp/pre254-portable.sh
#   git show d253703a^:scripts/winsec-tests.sh  >/tmp/pre254-winsec.sh
#   CARRIER_WINSEC_TESTS=/tmp/pre254-winsec.sh \
#     bash scripts/portable-tests-selftest.sh all /tmp/pre254-portable.sh
#
# WHAT IS REAL AND WHAT IS NOT. The `go` process is scripts/testdata/portable-tests/go
# (see its header for why a shim is the right tool here); everything else -
# scripts/portable-tests.sh, GUARD A/B/C, the ledger staleness check,
# tools/d22scan/runtests.sh - is the committed code, executed. This carrier makes
# NO test claim about the product: no Go test binary is built or run by it, so it
# cannot contaminate the milliseconds-sensitive packages other legs are reading.
# The import paths it feeds GUARD C are extracted from the script's own core_pin
# at run time, so the carrier cannot drift from the pin it audits.
#
# The pre-fix half of the reading is not a carrier mode - it is this command, run
# against scripts/portable-tests.sh as it stood before the ticket 250 fix:
#   cold=$(mktemp -d); GOMODCACHE="$cold" GOPATH="$cold/gopath" \
#     bash scripts/portable-tests.sh --scope=core          # exit 1 in ~12 s, 303 MB
# and its verbatim log is .scratch/wisp/probes/250/r1/logs/prefix-real-coldcache.txt.
#
# Usage:
#   bash scripts/portable-tests-selftest.sh                      # all cases
#   bash scripts/portable-tests-selftest.sh coldcache-stderr     # one case, verbose
#   bash scripts/portable-tests-selftest.sh <case> <other-portable-tests.sh>
#        Runs the same cases against another copy of the script - this is how the
#        pre-fix half of AC#3's reading is taken without touching the working
#        tree:  git show <sha>:scripts/portable-tests.sh >/tmp/prefix.sh  and point
#        at /tmp/prefix.sh. A shadow root is built for it (scripts/ + the real
#        tools/d22scan/runtests.sh + go.mod/go.sum), because the script under test
#        derives its own root from its location and refuses to run without that
#        strict runner. Only `go` is faked; the guards and runtests.sh are the
#        bytes being judged.
#   CARRIER_WINSEC_TESTS=/path/to/winsec-tests.sh
#        Same idea for the PARENT script the chain cases (19-22) drive. Unset = the
#        working copy. Set it together with the second argument to judge a whole
#        commit's bytes - see the 改前必红 reading in the header.
set -u -o pipefail

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
root=$(CDPATH= cd -- "$here/.." && pwd)
cd "$root"

script="$root/scripts/portable-tests.sh"
shim="$root/scripts/testdata/portable-tests/go"
for f in "$script" "$shim"; do
    if [ ! -f "$f" ]; then
        echo "portable-tests-selftest.sh: $f is missing - the carrier refuses to judge a copy." >&2
        exit 2
    fi
done

# Optional second argument: a portable-tests.sh to judge INSTEAD of the working one.
shadow=''
if [ -n "${2-}" ]; then
    if [ ! -f "$2" ]; then
        echo "portable-tests-selftest.sh: no such script to judge: $2" >&2
        exit 2
    fi
    shadow=$(mktemp -d 2>/dev/null || echo "$root/.portable-shadow.$$")
    mkdir -p "$shadow/scripts" "$shadow/tools/d22scan"
    cp "$2" "$shadow/scripts/portable-tests.sh"
    cp "$root/tools/d22scan/runtests.sh" "$shadow/tools/d22scan/runtests.sh"
    cp "$root/go.mod" "$shadow/go.mod"
    cp "$root/go.sum" "$shadow/go.sum"
    script="$shadow/scripts/portable-tests.sh"
    echo "portable-tests-selftest.sh: judging $2 in a shadow root at $shadow"
fi

work=$(mktemp -d 2>/dev/null || echo "$root/.portable-selftest.$$")
mkdir -p "$work/bin"
cp "$shim" "$work/bin/go"
chmod +x "$work/bin/go" 2>/dev/null || true

# One truth: the pin the carrier feeds is the pin the script carries.
pin="$work/core-pin.txt"
awk -v q="'" 'index($0, "core_pin=" q) == 1 { f = 1; next }
    f { if ($0 == q) { f = 0; next } print }' "$script" >"$pin"
pinc=$(grep -c . "$pin" || true)
echo "portable-tests-selftest.sh: core_pin extracted from $script = $pinc import paths"
if [ "$pinc" -eq 0 ]; then
    echo "portable-tests-selftest.sh: could not extract core_pin - the carrier would be vacuous." >&2
    exit 2
fi

# The ten stderr lines CI's run 36889094435 printed (verbatim, ticket 250's table).
cat >"$work/progress-stderr.txt" <<'EOF'
go: downloading github.com/dustin/go-humanize v1.0.1
go: downloading github.com/google/uuid v1.6.0
go: downloading github.com/pelletier/go-toml/v2 v2.2.4
go: downloading github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec
go: downloading golang.org/x/crypto v0.57.0
go: downloading golang.org/x/sys v0.48.0
go: downloading modernc.org/libc v1.75.7
go: downloading modernc.org/mathutil v1.7.1
go: downloading modernc.org/memory v1.12.1
go: downloading modernc.org/sqlite v1.59.0
EOF

failed=0
ran=0
last_rc=0
last_log=''
only=${1-}
[ "$only" = all ] && only=''

# run <case-name> [ENV=VAL assignment words...]
run() {
    local name=$1
    shift
    last_log="$work/$name.log"
    ( PATH="$work/bin:$PATH" FAKEGO_PIN_FILE="$pin" env "$@" \
        bash "$script" --scope=core ) >"$last_log" 2>&1
    last_rc=$?
    ran=$((ran + 1))
    printf '\n== case %s: rc=%d log=%s\n' "$name" "$last_rc" "$last_log"
}

want_rc() { # want_rc <expected>
    if [ "$last_rc" -eq "$1" ]; then
        printf '   ok   exit %s, as expected\n' "$1"
    else
        printf '   FAIL exit %d, expected %s\n' "$last_rc" "$1"
        failed=$((failed + 1))
    fi
}

has() { # has <regex> [count]  - the line set the case must print
    local n
    n=$(grep -Ec "$1" "$last_log" || true)
    if [ "$n" -ge "${2:-1}" ]; then
        printf '   ok   /%s/ appears %s time(s)\n' "$1" "$n"
    else
        printf '   FAIL /%s/ appeared %s time(s), wanted at least %s\n' "$1" "$n" "${2:-1}"
        failed=$((failed + 1))
    fi
}

hasnt() { # hasnt <regex>
    local n
    n=$(grep -Ec "$1" "$last_log" || true)
    if [ "$n" -eq 0 ]; then
        printf '   ok   /%s/ absent\n' "$1"
    else
        printf '   FAIL /%s/ appeared %s time(s):\n' "$1" "$n"
        grep -En "$1" "$last_log" | sed 's/^/        /'
        failed=$((failed + 1))
    fi
}

# The denominator GUARD C pins and GUARD B loops over is observable in the log:
# GUARD B prints exactly one result row per line of "$resolved".
denominator_is() {
    local n
    n=$(grep -Ec '^portable-tests\.sh:   (ok|FAIL) \(own line\)' "$last_log" || true)
    if [ "$n" -eq "$1" ]; then
        printf '   ok   denominator (GUARD B rows, one per resolved line) = %s\n' "$n"
    else
        printf '   FAIL denominator = %s, expected %s\n' "$n" "$1"
        failed=$((failed + 1))
    fi
}

skip() { printf '\n== case %s: skipped (single-case run)\n' "$1"; }
selected() { [ -z "$only" ] || [ "$only" = "$1" ]; }

# ---- ticket 251: the helpers the winsec/census cases need -----------------------
#
# make_shadow_root <script> [winsec-tests.sh] - the same shadow the second CLI
# argument builds: a script under test derives its own root from its location and
# refuses to run without tools/d22scan/runtests.sh next to it. Ticket 254 adds the
# optional second argument, because that ticket's claim lives in the PARENT script
# (see the chain cases at the end of this file) and a chain shadow has to carry the
# bytes of both links, not just the child's.
make_shadow_root() {
    local shadow
    shadow=$(mktemp -d 2>/dev/null || echo "$root/.portable-shadow.$$")
    mkdir -p "$shadow/scripts" "$shadow/tools/d22scan"
    cp "$1" "$shadow/scripts/portable-tests.sh"
    if [ -n "${2-}" ]; then
        if [ ! -f "$2" ]; then
            echo "portable-tests-selftest.sh: no such winsec-tests.sh to judge: $2" >&2
            exit 2
        fi
        cp "$2" "$shadow/scripts/winsec-tests.sh"
    fi
    cp "$root/tools/d22scan/runtests.sh" "$shadow/tools/d22scan/runtests.sh"
    cp "$root/go.mod" "$shadow/go.mod"
    cp "$root/go.sum" "$shadow/go.sum"
    printf '%s' "$shadow"
}

# run_with <case-name> <pin-file> [ENV=VAL ...] -- <script args...>
# run() with the scope no longer hardwired: ticket 251's cases have to point the
# winsec tier by name AND pass the explicit path winsec-tests.sh passes.
run_with() {
    local name=$1
    local pinfile=$2
    shift 2
    local envs=()
    while [ "$1" != "--" ]; do envs+=("$1"); shift; done
    shift
    last_log="$work/$name.log"
    ( PATH="$work/bin:$PATH" FAKEGO_PIN_FILE="$pinfile" \
        env ${envs[@]+"${envs[@]}"} bash "$script" "$@" ) >"$last_log" 2>&1
    last_rc=$?
    ran=$((ran + 1))
    printf '\n== case %s: rc=%d log=%s\n' "$name" "$last_rc" "$last_log"
}

# One truth again: the winsec tier's pin comes out of the script under test, so the
# carrier cannot drift from what the script claims. Its cases go red on the pre-fix
# bytes on purpose (no branch, no explicit-path pin) - see the header.
winsec_pin_file="$work/winsec-pin.txt"
awk -v q="'" 'index($0, "winsec_pin=" q) == 1 { f = 1; next }
    f { if ($0 == q) { f = 0; next } print }' "$script" >"$winsec_pin_file"
winc=$(grep -c . "$winsec_pin_file" || true)
echo "portable-tests-selftest.sh: winsec_pin extracted from $script = $winc import path(s)"
if [ "$winc" -eq 0 ]; then
    echo "portable-tests-selftest.sh: could not extract winsec_pin - the winsec cases would be vacuous." >&2
    exit 2
fi
IFS= read -r winsec_path <"$winsec_pin_file"
# seed 251-a: the split - a second package appears UNDER the pinned import path.
winsec_split="$work/winsec-pin-split.txt"
{ cat "$winsec_pin_file"; printf '%s\n' "$winsec_path/acl"; } >"$winsec_split"
echo "portable-tests-selftest.sh: split seed = $winsec_path + ${winsec_path}/acl (2 lines)"
# seed 254: the same split seen from the CORE tier's pin - a universe in which a
# package exists beside internal/winsec that NEITHER pin row claims. Cases 21 and 22
# run every tier against this one file, so "the two tiers disagree" can only mean one
# thing: they were given the same world.
core_split="$work/core-pin-split.txt"
{ cat "$pin"; printf '%s\n' "$winsec_path/acl"; } >"$core_split"
echo "portable-tests-selftest.sh: 254 split universe = core_pin ($pinc) + ${winsec_path}/acl ($((pinc + 1)) lines)"

# TICKET 254: the chain runner. Cases 11-18 judge scripts/portable-tests.sh on its
# own; the claim this ticket files is about the route CI REALLY takes -
# .github/workflows/ci.yml:439 runs scripts/winsec-tests.sh, which hands the child an
# EXPLICIT package path, and since this ticket the parent refuses to call itself green
# unless the child says it audited that path (winsec-tests.sh GUARD 3). A guard that
# exists only in the parent is invisible to every case above, so the chain gets its
# own shadow root and its own runner here. Only `go` is faked; both scripts and
# tools/d22scan/runtests.sh are the bytes being judged, executed.
#
# CARRIER_WINSEC_TESTS names the winsec-tests.sh bytes to judge, exactly the way the
# second CLI argument names the portable-tests.sh bytes: the 改前必红 half of this
# ticket points BOTH links at one commit's bytes without touching the working tree.
# Default = the working copy.
winsec_script=${CARRIER_WINSEC_TESTS:-$root/scripts/winsec-tests.sh}

# run_chain <case-name> <universe-pin-file> [ENV=VAL ...] -- <winsec-tests.sh args...>
# The universe file plays the role of "the packages that exist on disk": the fake go
# resolves the scope against it the way real go does. Every shape inside one case
# shares ONE universe file, so the only thing a case varies is the argument list the
# parent is given - "same seed, two shapes" in the ticket's words.
run_chain() {
    local name=$1
    local universe=$2
    shift 2
    local envs=()
    while [ "$1" != "--" ]; do envs+=("$1"); shift; done
    shift
    last_log="$work/$name.log"
    local shadow
    shadow=$(make_shadow_root "$script" "$winsec_script")
    ( cd "$shadow" && PATH="$work/bin:$PATH" FAKEGO_PIN_FILE="$universe" \
        FAKEGO_SCOPE_FILTER=1 FAKEGO_GOOS=windows \
        env ${envs[@]+"${envs[@]}"} bash scripts/winsec-tests.sh "$@" ) >"$last_log" 2>&1
    last_rc=$?
    ran=$((ran + 1))
    printf '\n== case %s: rc=%d log=%s\n' "$name" "$last_rc" "$last_log"
}

# 1. control: the pin resolves, nothing on stderr. The step must be GREEN, which
#    is what stops cases 2-9 from passing for the wrong reason.
if selected clean; then
    run clean
    want_rc 0
    denominator_is "$pinc"
    hasnt 'GUARD [ABC] -'
    has '^runtests\.sh: OK'
fi

# 2. AC#1 + AC#3, the ticket's own positive control: stderr carries the cold-cache
#    progress lines, stdout is exactly the pin. Denominator must stay 25 and GUARD
#    C must stay silent - and stderr must still be VISIBLE in the step log.
if selected coldcache-stderr; then
    run coldcache-stderr FAKEGO_STDERR_FILE="$work/progress-stderr.txt"
    want_rc 0
    denominator_is "$pinc"
    hasnt 'GUARD [ABC] -'
    has 'stderr \(progress, not packages, not in the denominator\)'
    has '^portable-tests\.sh:   err\| go: downloading ' 10
fi

# 3. AC#2 positive control, the ticket's literal wording: a pseudo package carrying
#    a space reaches STDOUT. The shape check must refuse the run, not filter it.
if selected seed-stdout-space; then
    printf '%s\n' 'github.com/CarlosShao/wisp/internal/fake pkg' >"$work/seed-space.txt"
    run seed-stdout-space FAKEGO_STDOUT_EXTRA_FILE="$work/seed-space.txt"
    want_rc 1
    has "GUARD C - go list's stdout carried 1 line\(s\) that are not"
    has 'github.com/CarlosShao/wisp/internal/fake pkg'
    hasnt '^runtests\.sh: OK'
fi

# 4. AC#2, the same seeded line shape but with the ORIGINAL defect's text: proves a
#    future re-merge of the streams is caught by the grammar and not by luck.
if selected seed-stdout-progress; then
    run seed-stdout-progress FAKEGO_STDOUT_EXTRA_FILE="$work/progress-stderr.txt"
    want_rc 1
    has "GUARD C - go list's stdout carried 10 line\(s\) that are not"
    hasnt '^runtests\.sh: OK'
fi

# 5. AC#2's other half: a seeded line with a VALID import-path shape but not in the
#    pin must be caught by GUARD C's set comparison, i.e. the shape check did not
#    become the only thing left standing.
if selected seed-stdout-valid-but-unpinned; then
    printf '%s\n' 'github.com/CarlosShao/wisp/internal/notpinned' >"$work/seed-valid.txt"
    run seed-stdout-valid-but-unpinned FAKEGO_STDOUT_EXTRA_FILE="$work/seed-valid.txt"
    want_rc 1
    has 'GUARD C - scope mode=core resolved to a DIFFERENT package'
    has '> github.com/CarlosShao/wisp/internal/notpinned'
fi

# 6. AC#1's reverse face: `go list` really fails. The step must still exit 1 and
#    print BOTH streams (the pre-fix branch at :262-266, kept, not softened).
if selected golist-fails; then
    printf '%s\n' 'github.com/CarlosShao/wisp/internal/partiallyresolved' >"$work/fail-stdout.txt"
    run golist-fails \
        FAKEGO_STDERR_FILE="$work/progress-stderr.txt" \
        FAKEGO_STDOUT_EXTRA_FILE="$work/fail-stdout.txt" \
        FAKEGO_LIST_RC=1
    want_rc 1
    has 'exited non-zero. Its stdout was:'
    has '^portable-tests\.sh:   out\| github.com/CarlosShao/wisp/internal/partiallyresolved'
    has '^portable-tests\.sh:   err\| go: downloading modernc.org/sqlite v1.59.0'
    has 'the scope cannot be audited'
fi

# 7. AC#4: GUARD C in the OTHER direction - a scope entry that quietly stops
#    resolving must fail the step, not shorten it.
if selected pin-drift-one-package-missing; then
    run pin-drift-one-package-missing \
        FAKEGO_LIST_DROP="github.com/CarlosShao/wisp/internal/winsec"
    want_rc 1
    has 'GUARD C - scope mode=core resolved to a DIFFERENT package'
    has '< github.com/CarlosShao/wisp/internal/winsec'
fi

# 8. AC#4: GUARD A, seeded empty denominator.
if selected guard-a-empty-test-package; then
    run guard-a-empty-test-package \
        FAKEGO_GUARD_A_SEED="github.com/CarlosShao/wisp/internal/session"
    want_rc 1
    has 'GUARD A - these packages are declared in scope mode=core'
    has 'github.com/CarlosShao/wisp/internal/session'
fi

# 9. AC#4: GUARD B, a resolved package that prints no top-level result line.
if selected guard-b-no-result-line; then
    run guard-b-no-result-line \
        FAKEGO_DROP_PKG="github.com/CarlosShao/wisp/internal/perm"
    want_rc 1
    has 'GUARD B - 1 of .* packages in scope mode=core printed'
    has 'github.com/CarlosShao/wisp/internal/perm'
fi

# 10. AC#4: the empty-scope hard exit. Unreachable from the CLI at this HEAD
#     (every named scope is non-empty, `census` exits earlier, the explicit-path
#     branch only takes over when $# > 0), so the carrier runs the guard's OWN
#     committed BYTES: the four-line block is sliced out verbatim (if / echo / exit
#     2 / fi), printed, and fed an empty scope array.
if selected empty-scope-slice; then
    slice="$work/empty-scope-slice.sh"
    anchor='if [ ${#scope[@]} -eq 0 ]; then'
    awk -v s="$anchor" '$0 == s { f = 1 }
        f { print } f && $0 == "fi" { exit }' "$script" >"$slice"
    slicec=$(grep -c . "$slice" || true)
    printf '\n== case empty-scope-slice: %s line(s) sliced verbatim from %s\n' "$slicec" "$script"
    ran=$((ran + 1))
    if [ "$slicec" -ne 4 ]; then
        printf '   FAIL the slice is not the 4-line guard (if/echo/exit/fi) any more (%s lines) -' "$slicec"
        printf ' re-read the carrier\n        instead of trusting this case:\n'
        sed 's/^/        /' "$slice"
        failed=$((failed + 1))
    else
        sed 's/^/   slice| /' "$slice"
        {
            echo 'set -eu -o pipefail'
            echo 'scope=()'
            echo 'mode=carrier'
            cat "$slice"
        } >"$work/empty-scope-carrier.sh"
        bash "$work/empty-scope-carrier.sh" >"$work/empty-scope-slice.log" 2>&1
        last_rc=$?
        last_log="$work/empty-scope-slice.log"
        printf '   (log: %s)\n' "$last_log"
        want_rc 2
        has 'refusing to be a green no-op'
    fi
fi

# 11. ticket 251 AC#1's control for the new tier: --scope=winsec resolves to exactly
#     winsec_pin. Without this the next three cases could pass for the wrong reason.
#     FAKEGO_SCOPE_FILTER makes the fake go resolve the args the way real go does
#     (./internal/winsec/ is ONE package; ./internal/winsec/... is the subtree), so
#     "the pin and the resolved set differ" below is a claim about resolution, not
#     about a file of names.
if selected winsec-tier-clean; then
    run_with winsec-tier-clean "$winsec_pin_file" "FAKEGO_SCOPE_FILTER=1" -- --scope=winsec
    want_rc 0
    denominator_is "$winc"
    hasnt 'GUARD [ABC] -'
    has '^runtests\.sh: OK'
fi

# 12. AC#1 + AC#2, seed '种假包名' on the NAMED tier: the tier's glob resolves a
#     second package under internal/winsec that no pin row claims.
if selected winsec-tier-split-goes-red; then
    run_with winsec-tier-split-goes-red "$winsec_split" "FAKEGO_SCOPE_FILTER=1" -- --scope=winsec
    want_rc 1
    has 'GUARD C - scope mode=winsec resolved to a DIFFERENT package'
    has 'Pinned: 1, resolved: 2'
    has '> github\.com/CarlosShao/wisp/internal/winsec/acl'
    hasnt '^runtests\.sh: OK'
fi

# 13. AC#1's other half, and the one that matters for CI: the EXPLICIT path
#     scripts/winsec-tests.sh:94 actually passes. Same seed as case 12, and the log
#     must say the explicit list was recognised as the tier's own scope - otherwise
#     this is just case 12 again, not the CI route being pinned.
if selected winsec-explicit-split-goes-red; then
    run_with winsec-explicit-split-goes-red "$winsec_split" "FAKEGO_SCOPE_FILTER=1" -- ./internal/winsec/
    want_rc 1
    has 'explicit scope .*IS the winsec tier'
    has 'GUARD C - scope mode=winsec resolved to a DIFFERENT package'
    has '> github\.com/CarlosShao/wisp/internal/winsec/acl'
    hasnt '^runtests\.sh: OK'
fi

# 14. AC#2's '删真包名' control on the same route: the pinned package stops
#     resolving. Pre-fix this was the loudest reading of the ticket - an explicit
#     scope that resolved to NOTHING still exited 0, because GUARD B audited an
#     empty denominator and GUARD C had nothing to pin.
if selected winsec-explicit-pin-gone-goes-red; then
    run_with winsec-explicit-pin-gone-goes-red "$winsec_pin_file" \
        "FAKEGO_SCOPE_FILTER=1" "FAKEGO_LIST_DROP=$winsec_path" -- ./internal/winsec/
    want_rc 1
    has 'GUARD C - scope mode=winsec resolved to a DIFFERENT package'
    has 'Pinned: 1, resolved: 0'
    has '< github\.com/CarlosShao/wisp/internal/winsec'
    hasnt '^runtests\.sh: OK'
fi

# 15. AC#4's control: census reports the tier roster it is about to be judged by,
#     and the winsec row is claimed by BOTH core and winsec.
if selected census-tier-roster-clean; then
    run_with census-tier-roster-clean "$pin" -- --scope=census
    want_rc 0
    has 'census totals: packages='
    has 'github\.com/CarlosShao/wisp/internal/winsec .*corewinsec'
fi

# 16. AC#4's seed: the tier roster loses a name while the `case $mode in` branch
#     stays - the exact shape this ticket was filed on (census claimed a winsec tier
#     no branch implemented, and exited 0 doing it). The seed is put on a MUTATED
#     COPY in its own shadow root; if the copy's roster line is not there to mutate,
#     the case says so instead of reporting a vacuous green.
if selected census-tier-roster-drift-goes-red; then
    newroster="tiers='core windows cli census'"
    mutant="$work/portable-tests-drift.sh"
    awk -v r="$newroster" '/^tiers=/{ print r; next } { print }' "$script" >"$mutant"
    touched=$(diff "$script" "$mutant" | grep -c '^[<>]' || true)
    printf '\n== case census-tier-roster-drift-goes-red: seed diff shows %s line(s) (2 = the one' "$touched"
    printf " tiers= line replaced, 0 = nothing to mutate)\n"
    ran=$((ran + 1))
    if [ "$touched" -eq 0 ]; then
        printf '   FAIL no tiers= roster line in the script under test - the seed cannot take, and a case that'
        printf '   cannot seed\n        is not evidence (pre-fix bytes fail here: that is ticket 251 AC#4).\n'
        failed=$((failed + 1))
    else
        mshadow=$(make_shadow_root "$mutant")
        last_log="$work/census-tier-roster-drift-goes-red.log"
        ( PATH="$work/bin:$PATH" FAKEGO_PIN_FILE="$pin" \
            bash "$mshadow/scripts/portable-tests.sh" --scope=census ) >"$last_log" 2>&1
        last_rc=$?
        printf '   (log: %s)\n' "$last_log"
        want_rc 1
        has 'census - the tier roster \(tiers=\) and the'
        has 'branches : .*winsec'
        hasnt 'census totals: packages='
    fi
fi

# 17. AC#3's regression pin: the fifth tier must not have softened the loud failure.
#     An unknown --scope is still exit 2 on stderr, naming every tier that exists.
if selected unknown-scope-still-hard; then
    run_with unknown-scope-still-hard "$pin" -- --scope=winsecfoo
    want_rc 2
    has 'unknown --scope=winsecfoo \(known: core, windows, cli, winsec, census\)'
    hasnt 'refusing to be a green no-op'
    hasnt '^runtests\.sh: OK'
fi

# 18. AC#1's 选型凭据. The ticket offers two ways to close the hole and asks for
#     evidence from whoever picks the second one, so here it is, as a reading instead
#     of an argument: a MUTATED copy of the current script that pins winsec_pin on the
#     explicit route but leaves the caller's DIRECTORY form as the scope - the literal
#     shape of option (b). Same seed, fake go resolving like real go:
#       real script (tier scope = ./internal/winsec/...)  -> RED,  Pinned: 1, resolved: 2
#       mutant    (scope = the caller's ./internal/winsec/) -> GREEN
#     because `go list ./internal/winsec/` names the one package in that directory and
#     can never see a second one created beside it. "The tier's scope is a glob" is
#     therefore load-bearing, not decoration - and the counterfactual is pinned here so
#     it cannot be quietly dropped later.
if selected winsec-explicit-glob-is-what-bites; then
    noglob="$work/portable-tests-noglob.sh"
    awk -v old='        scope=("$winsec_scope")' -v new='        scope=("$@")' \
        '$0 == old { print new; next } { print }' "$script" >"$noglob"
    touched=$(diff "$script" "$noglob" | grep -c '^<' || true)
    printf '\n== case winsec-explicit-glob-is-what-bites: mutant replaced %s line (the scope of the explicit-path branch)\n' "$touched"
    ran=$((ran + 1))
    if [ "$touched" -ne 1 ]; then
        printf '   FAIL the mutant seed took %s line(s), expected exactly 1 - this case cannot judge a' "$touched"
        printf ' shape it did not find\n        (the explicit-path substitution moved; re-read the script, not the carrier)\n'
        failed=$((failed + 1))
    else
        last_log="$work/winsec-explicit-glob-is-what-bites-real.log"
        ( PATH="$work/bin:$PATH" FAKEGO_PIN_FILE="$winsec_split" FAKEGO_SCOPE_FILTER=1 \
            bash "$script" ./internal/winsec/ ) >"$last_log" 2>&1
        last_rc=$?
        printf '   real bytes (tier scope is the glob):\n'
        want_rc 1
        has 'GUARD C - scope mode=winsec resolved to a DIFFERENT package'
        nshadow=$(make_shadow_root "$noglob")
        last_log="$work/winsec-explicit-glob-is-what-bites-mutant.log"
        ( PATH="$work/bin:$PATH" FAKEGO_PIN_FILE="$winsec_split" FAKEGO_SCOPE_FILTER=1 \
            bash "$nshadow/scripts/portable-tests.sh" ./internal/winsec/ ) >"$last_log" 2>&1
        last_rc=$?
        printf '   mutant (same pin, caller keeps the directory form) - the counterfactual:\n'
        want_rc 0
        hasnt 'GUARD C - scope mode=winsec'
        has '^runtests\.sh: OK'
        printf '   (logs: %s / %s)\n' "${last_log%-mutant.log}-real.log" "$last_log"
    fi
fi

if selected chain-explicit-split-audit-bites; then
    # positive: the SAME seed case 12 shoots the named tier with, handed over the route
    # ci.yml:439 takes - and it takes it with ZERO arguments, so that is how the chain
    # is called here (winsec-tests.sh then builds its own one-path scope). The audit has
    # to bite THROUGH the parent, not merely be announced by it: "my scope was
    # recognised" that nobody could falsify is the same claim as "配置里有一行 = 它跑过了".
    run_chain chain-explicit-split-audit-bites "$winsec_split" --
    want_rc 1
    has 'explicit scope .*IS the winsec tier'
    has 'GUARD C - scope mode=winsec resolved to a DIFFERENT package'
    has '> github\.com/CarlosShao/wisp/internal/winsec/acl'
    hasnt '^runtests\.sh: OK'
    # guard 3 declares itself only on an OTHERWISE GREEN step: a run the child already
    # refused must not get a second, misattributed message. Pinned here because the
    # only way to lose it is to notice it was never checked.
    hasnt 'winsec-tests\.sh: GUARD 3'
    # reverse face: same chain, same fake go, nothing split - green, and the parent
    # stays quiet BECAUSE the child said it audited. Without this the case above could
    # go red for the wrong reason (a chain that never audits anything).
    run_chain chain-explicit-clean-stays-green "$winsec_pin_file" --
    want_rc 0
    has 'explicit scope .*IS the winsec tier'
    has '^runtests\.sh: OK'
    hasnt 'winsec-tests\.sh: GUARD'
    hasnt 'GUARD [ABC] -'
fi

# 20. AC#1's ⓑ teeth: the hole GUARD 3 exists for. This script's own header advertises
#     "extra packages may be appended by hand while debugging", and the moment it does,
#     the child drops out of its one-path tier branch into "nothing to pin": the step
#     stays GREEN (before this ticket's guard existed it measured exactly that - rc=0,
#     zero audit lines, .scratch/wisp/probes/254/r1/logs/before-widened.txt). One
#     universe, three shapes, only the argument list varies.
if selected chain-widened-scope-loses-audit; then
    run_chain chain-widened-scope-loses-audit "$pin" -- ./internal/winsec/ ./internal/agent/
    want_rc 2
    has 'winsec-tests\.sh: GUARD 3'
    has 'the delegated audit therefore did not happen'
    hasnt 'IS the winsec tier'
    # the child WAS green - that is the whole point: rc=2 here comes from the audit
    # going missing, not from a test failing.
    has '^runtests\.sh: OK'
    # reverse face, and the literal CI call (zero arguments, the parent's own one-path
    # default): same universe, same fake go - green, because the child says it audited.
    run_chain chain-ci-shape-is-audited "$pin" --
    want_rc 0
    has 'IS the winsec tier'
    hasnt 'winsec-tests\.sh: GUARD 3'
    # and the OTHER branch AC#1 offered (ⓐ, point the parent at --scope=winsec) is not
    # reachable by passing arguments either: guard 1 refuses a list that does not name
    # the package. So the ⓐ/ⓑ choice cannot be made silently from a call site - it
    # needs an edit to this script, which is where the guards live.
    run_chain chain-scope-flag-not-accepted "$pin" -- --scope=winsec
    want_rc 2
    has 'winsec-tests\.sh: GUARD 1 - the package list for this step does not name'
fi

# 21. AC#2: the two tiers that both claim internal/winsec must AGREE about a split of
#     it. One universe with a package created beside internal/winsec; both tiers run
#     against it; both must refuse, and the package they refuse over must be the SAME
#     line. Two tiers, two spellings, two verdicts is the shape the ticket forbids, and
#     the textual face of the rule is checked in the same case: the core branch may not
#     carry its own spelling of that package, it has to name the shared variable.
if selected core-and-winsec-see-the-same-split; then
    run_with core-tier-split-goes-red "$core_split" "FAKEGO_SCOPE_FILTER=1" -- --scope=core
    want_rc 1
    has 'GUARD C - scope mode=core resolved to a DIFFERENT package'
    has '> github\.com/CarlosShao/wisp/internal/winsec/acl'
    run_with winsec-tier-split-same-seed "$core_split" "FAKEGO_SCOPE_FILTER=1" -- --scope=winsec
    want_rc 1
    has 'GUARD C - scope mode=winsec resolved to a DIFFERENT package'
    has '> github\.com/CarlosShao/wisp/internal/winsec/acl'
    a=$(sed -n 's/^portable-tests\.sh:   > //p' "$work/core-tier-split-goes-red.log" | sort -u)
    b=$(sed -n 's/^portable-tests\.sh:   > //p' "$work/winsec-tier-split-same-seed.log" | sort -u)
    if [ -n "$a" ] && [ "$a" = "$b" ]; then
        printf '   ok   both tiers refuse the same seed over the SAME package: %s\n' "$a"
    else
        printf '   FAIL the two tiers disagree about what bit - core says [%s], the winsec tier' "$a"
        printf ' says [%s].\n        One package, two answers: ticket 254 AC#2.\n' "$b"
        failed=$((failed + 1))
    fi
    core_branch=$(awk '/^core\)$/{f=1;next} f&&/^    ;;$/{exit} f' "$script")
    hard=$(printf '%s\n' "$core_branch" | grep -cF './internal/winsec/' || true)
    shared=$(printf '%s\n' "$core_branch" | grep -cF '"$winsec_scope"' || true)
    if [ "$hard" -eq 0 ] && [ "$shared" -eq 1 ]; then
        printf '   ok   the core branch names the winsec package through the ONE shared variable'
        printf ' (%s reference, %s private spelling)\n' "$shared" "$hard"
    else
        printf '   FAIL core branch carries %s private spelling(s) of ./internal/winsec/ and %s' "$hard" "$shared"
        printf ' shared\n        reference(s) - two spellings means the two tiers can resolve one package\n'
        printf '        two ways (ticket 254 AC#2). Pre-d253703a bytes fail here.\n'
        failed=$((failed + 1))
    fi
    # control for all of the above: on an unsplit universe BOTH tiers are green, i.e.
    # the case above went red because of the seed and not because the tiers are broken.
    run_with core-tier-clean-globform "$pin" "FAKEGO_SCOPE_FILTER=1" -- --scope=core
    want_rc 0
    denominator_is "$pinc"
    hasnt 'GUARD [ABC] -'
fi

# 22. AC#2's counterfactual, the same shape as case 18 and for the same reason: "the
#     core row had to change to the glob" is load-bearing only if the directory form
#     demonstrably hides the split. MUTATED COPY, same seed, same fake go - core goes
#     GREEN (it cannot see a package created beside the directory) while the winsec
#     branch of the very same script goes RED over it. That asymmetry inside one file
#     is the defect this ticket closes; pinning it is what stops a later edit from
#     quietly restoring the directory form and calling the two tiers consistent.
if selected core-dir-form-hides-the-split; then
    dirform="$work/portable-tests-core-dirform.sh"
    awk -v old='        ./internal/plugin/ ./cmd/llmrecord/ "$winsec_scope"' \
        -v new='        ./internal/plugin/ ./cmd/llmrecord/ ./internal/winsec/' \
        '$0 == old { print new; next } { print }' "$script" >"$dirform"
    touched=$(diff "$script" "$dirform" | grep -c '^<' || true)
    printf '\n== case core-dir-form-hides-the-split: mutant replaced %s line (core'"'"'s winsec row;' "$touched"
    printf ' 0 = nothing to mutate, and a seed that cannot take is not evidence)\n'
    ran=$((ran + 1))
    if [ "$touched" -ne 1 ]; then
        printf '   FAIL the mutant seed took %s line(s), expected exactly 1 - the core tier'"'"'s winsec' "$touched"
        printf ' row moved\n        out from under this case; re-read the script instead of trusting it\n'
        failed=$((failed + 1))
    else
        dshadow=$(make_shadow_root "$dirform")
        last_log="$work/core-dir-form-hides-the-split.log"
        ( cd "$dshadow" && PATH="$work/bin:$PATH" FAKEGO_PIN_FILE="$core_split" \
            FAKEGO_SCOPE_FILTER=1 bash scripts/portable-tests.sh --scope=core ) >"$last_log" 2>&1
        last_rc=$?
        printf '   mutant (core names the DIRECTORY, split package beside it invisible):\n'
        printf '   (log: %s)\n' "$last_log"
        want_rc 0
        hasnt 'GUARD [ABC] -'
        has '^runtests\.sh: OK'
        denominator_is "$pinc"
        # the seed is live in the same bytes: the tier branch, run by the mutant, does
        # see it. Core green + tier red on ONE universe is the double standard.
        last_log="$work/core-dir-form-hides-the-split-tier.log"
        ( cd "$dshadow" && PATH="$work/bin:$PATH" FAKEGO_PIN_FILE="$core_split" \
            FAKEGO_SCOPE_FILTER=1 bash scripts/portable-tests.sh --scope=winsec ) >"$last_log" 2>&1
        last_rc=$?
        printf '   same mutant, --scope=winsec, same universe - the tier it does not hide:\n'
        want_rc 1
        has 'GUARD C - scope mode=winsec resolved to a DIFFERENT package'
        has '> github\.com/CarlosShao/wisp/internal/winsec/acl'
        printf '   (logs: %s / %s)\n' "$work/core-dir-form-hides-the-split.log" "$last_log"
    fi
fi

if [ "$only" != "" ] && [ "$ran" -eq 0 ]; then
    echo "portable-tests-selftest.sh: no such case: $only" >&2
    exit 2
fi

printf '\nportable-tests-selftest.sh: %d case(s) ran, %d assertion(s) failed\n' "$ran" "$failed"
printf 'portable-tests-selftest.sh: carrier scratch kept at %s (rules: temp files are created, not deleted)\n' "$work"
if [ "$failed" -ne 0 ]; then
    echo "portable-tests-selftest.sh: RED - a guard stopped biting, or the carrier went stale." >&2
    exit 1
fi
echo "portable-tests-selftest.sh: GREEN - every seeded anomaly was refused, and the clean scope passed."
exit 0
