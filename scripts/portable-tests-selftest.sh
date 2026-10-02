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
# reproduced here on demand: seventeen cases, each of which either has to go red or
# has to stay green, and the carrier exits non-zero if any of them does the other.
#
# TICKET 251 ADDED THE LAST SEVEN, and they judge a different hole: `winsec_pin`
# was a pin NOTHING compared - no `winsec)` branch existed, and the route CI really
# takes (scripts/winsec-tests.sh:94) hands scripts/portable-tests.sh an EXPLICIT
# package path, which landed in GUARD C's "nothing to pin" branch. So the same
# carrier now also shoots the winsec tier by name, the explicit path winsec-tests.sh
# passes, and the census tier roster, with the ticket's two seed shapes on each
# (stdout gains a package the pin does not have / the pinned package disappears).
# Cases 11-17 FAIL on the pre-fix bytes by design - that is this ticket's 改前必红
# reading:  bash scripts/portable-tests-selftest.sh all /tmp/prefix.sh
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
# make_shadow_root <script> - the same shadow the second CLI argument builds: a
# script under test derives its own root from its location and refuses to run
# without tools/d22scan/runtests.sh next to it.
make_shadow_root() {
    local shadow
    shadow=$(mktemp -d 2>/dev/null || echo "$root/.portable-shadow.$$")
    mkdir -p "$shadow/scripts" "$shadow/tools/d22scan"
    cp "$1" "$shadow/scripts/portable-tests.sh"
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
if selected winsec-tier-clean; then
    run_with winsec-tier-clean "$winsec_pin_file" -- --scope=winsec
    want_rc 0
    denominator_is "$winc"
    hasnt 'GUARD [ABC] -'
    has '^runtests\.sh: OK'
fi

# 12. AC#1 + AC#2, seed '种假包名' on the NAMED tier: the tier's glob resolves a
#     second package under internal/winsec that no pin row claims.
if selected winsec-tier-split-goes-red; then
    run_with winsec-tier-split-goes-red "$winsec_split" -- --scope=winsec
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
    run_with winsec-explicit-split-goes-red "$winsec_split" -- ./internal/winsec/
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
        "FAKEGO_LIST_DROP=$winsec_path" -- ./internal/winsec/
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
