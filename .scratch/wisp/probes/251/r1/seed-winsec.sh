#!/usr/bin/env bash
# .scratch/wisp/probes/251/r1/seed-winsec.sh - ticket 251's 改前必红 credential.
#
# Same seeds, same carrier parts, twice: once against the bytes at this leg's
# starting anchor (改前) and once against the working tree (改后), so "the guard
# started biting because of this ticket" is a reading and not an assertion. The `go`
# PROCESS is the committed shim scripts/testdata/portable-tests/go; the script under
# test, its three guards and tools/d22scan/runtests.sh run for real. NO real
# go test / go build is started by this probe - that is what the shim is for, and
# cmd/wisp + internal/config + internal/ball are owned by other legs right now.
#
# Seeds, all on the route CI ACTUALLY takes for winsec (scripts/winsec-tests.sh:94 ->
# `bash scripts/portable-tests.sh ./internal/winsec/`, i.e. an explicit path):
#   control  the tier resolves to exactly its pin             -> GREEN both ways
#   split    a second package appears under internal/winsec   -> 改前 GREEN / 改后 RED
#   gone     the pinned package stops resolving               -> 改前 GREEN / 改后 RED
#   scope    --scope=winsec                                    -> 改前 rc=2 / 改后 GREEN
#   unknown  --scope=winsecfoo                                 -> rc=2 both ways (AC#3)
#   core     winsec missing from the core scope                -> RED both ways (untouched guard)
#   roster   tiers= mutated to drop winsec, case branch kept   -> 改前 n/a / 改后 census RED
#
# Usage:
#   bash .scratch/wisp/probes/251/r1/seed-winsec.sh
#   bash .scratch/wisp/probes/251/r1/seed-winsec.sh /tmp/251-prefix-portable-tests.sh
set -u -o pipefail

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
root=$(CDPATH= cd -- "$here/../../../../.." && pwd)
cd "$root"

target=${1-$root/scripts/portable-tests.sh}
shim="$root/scripts/testdata/portable-tests/go"
winsec_import=github.com/CarlosShao/wisp/internal/winsec

work=$(mktemp -d 2>/dev/null || echo "$root/.portable-251-seed.$$")
mkdir -p "$work/bin" "$work/logs"
cp "$shim" "$work/bin/go"
chmod +x "$work/bin/go" 2>/dev/null || true

make_shadow() { # make_shadow <script-to-judge> -> echoes the shadow root
    local shadow
    shadow=$(mktemp -d 2>/dev/null || echo "$work/shadow.$RANDOM")
    mkdir -p "$shadow/scripts" "$shadow/tools/d22scan"
    cp "$1" "$shadow/scripts/portable-tests.sh"
    cp "$root/tools/d22scan/runtests.sh" "$shadow/tools/d22scan/runtests.sh"
    cp "$root/go.mod" "$shadow/go.mod"
    cp "$root/go.sum" "$shadow/go.sum"
    printf '%s' "$shadow"
}

if [ "$(CDPATH= cd -- "$(dirname -- "$target")" && pwd)/$(basename -- "$target")" = "$root/scripts/portable-tests.sh" ]; then
    script="$root/scripts/portable-tests.sh"
    echo "seed-winsec.sh: judging the WORKING TREE bytes"
else
    shadow=$(make_shadow "$target")
    script="$shadow/scripts/portable-tests.sh"
    echo "seed-winsec.sh: judging $target in a shadow root at $shadow"
fi

# The tier's own pin, read out of the script under test (one truth; same awk shape
# the carrier uses for core_pin). A script without a winsec_pin falls back to the
# import path, which is what a split-proof needs, not what it pretends to audit.
pin="$work/winsec-pin.txt"
awk -v q="'" 'index($0, "winsec_pin=" q) == 1 { f = 1; next }
    f { if ($0 == q) { f = 0; next } print }' "$script" >"$pin"
[ -s "$pin" ] || printf '%s\n' "$winsec_import" >"$pin"
echo "seed-winsec.sh: winsec_pin read from the script under test = $(grep -c . "$pin") line(s)"

seed="$work/winsec-split.txt"
{ cat "$pin"; printf '%s\n' "$winsec_import/acl"; } >"$seed"

run_shot() { # run_shot <name> <script> <ENV=V>... -- <script args...>
    local name=$1
    local which=$2
    shift 2
    local envs=()
    while [ "$1" != "--" ]; do envs+=("$1"); shift; done
    shift
    local log="$work/logs/$name.txt"
    ( PATH="$work/bin:$PATH" env "${envs[@]}" bash "$which" "$@" ) >"$log" 2>&1
    local rc=$?
    printf '\n== %s: rc=%d  [args: %s]\n' "$name" "$rc" "$*"
    grep -nE 'portable-tests\.sh: GUARD [ABC] -|Pinned: [0-9]+, resolved|four numbers|unknown --scope|refusing to be a green no-op|strict runner exited|runtests\.sh: OK|explicit scope|census -|tiers=|branches :' "$log" |
        sed 's/^/   | /' || true
    printf '   (log: %s, %s non-empty line(s))\n' "$log" "$(grep -c . "$log")"
}

echo
run_shot control-explicit-winsec "$script" "FAKEGO_PIN_FILE=$pin" -- ./internal/winsec/
run_shot split-explicit-winsec "$script" "FAKEGO_PIN_FILE=$seed" -- ./internal/winsec/
run_shot gone-explicit-winsec "$script" "FAKEGO_PIN_FILE=$pin" "FAKEGO_LIST_DROP=$winsec_import" -- ./internal/winsec/
run_shot scope-winsec "$script" "FAKEGO_PIN_FILE=$pin" -- --scope=winsec
run_shot unknown-scope "$script" "FAKEGO_PIN_FILE=$pin" -- --scope=winsecfoo
run_shot winsec-missing-from-core "$script" "FAKEGO_PIN_FILE=$pin" "FAKEGO_LIST_DROP=$winsec_import" -- --scope=core

# The same three seeds with the fake go resolving its ARGUMENTS the way real go does
# (./internal/winsec/ = the one package in that directory; ./internal/winsec/... = the
# subtree). Without this the split shot would be an artifact of a shim that prints
# every name it knows for any argument; with it, the shots mean what the ticket's own
# risk sentence says ("拆出第二包").
echo
echo "### (filtered) the same three seeds, resolution modelled like real go"
run_shot split-explicit-filtered "$script" "FAKEGO_PIN_FILE=$seed" "FAKEGO_SCOPE_FILTER=1" -- ./internal/winsec/
run_shot scope-winsec-filtered "$script" "FAKEGO_PIN_FILE=$seed" "FAKEGO_SCOPE_FILTER=1" -- --scope=winsec
run_shot gone-explicit-filtered "$script" "FAKEGO_PIN_FILE=$pin" "FAKEGO_SCOPE_FILTER=1" "FAKEGO_LIST_DROP=$winsec_import" -- ./internal/winsec/

# The counterfactual (only meaningful against the working-tree bytes): option (b) as
# the ticket literally words it - the explicit route gets the pin, but the scope stays
# the caller's DIRECTORY form. Same split seed, real-like resolution: GREEN. That is
# why the tier's own scope here is the glob, and it is a reading, not an argument.
echo
echo "### counterfactual: pinned but still the caller's directory form"
noglob="$work/portable-tests-noglob.sh"
if grep -qx '        scope=("$winsec_scope")' "$script"; then
    awk -v old='        scope=("$winsec_scope")' -v new='        scope=("$@")' \
        '$0 == old { print new; next } { print }' "$script" >"$noglob"
    printf '   mutant replaced %s line (the explicit branch scope)\n' \
        "$(diff "$script" "$noglob" | grep -c '^<' || true)"
    nshadow=$(make_shadow "$noglob")
    run_shot counterfactual-noglob "$nshadow/scripts/portable-tests.sh" \
        "FAKEGO_PIN_FILE=$seed" "FAKEGO_SCOPE_FILTER=1" -- ./internal/winsec/
    run_shot real-bytes-glob-bites "$script" \
        "FAKEGO_PIN_FILE=$seed" "FAKEGO_SCOPE_FILTER=1" -- ./internal/winsec/
else
    echo "   (no explicit-path substitution in the script under test - 改前 bytes: the"
    echo "    counterfactual has nothing to mutate; the 改前 split shot above IS the"
    echo "    reading, and it came back rc=0)"
fi

# AC#4's seed: mutate ONLY the roster line, keep the case branch, run --scope=census.
echo
mutant="$work/mutated-portable-tests.sh"
if grep -q "^tiers='.*winsec.*'$" "$script"; then
    sed "s/^tiers='core windows cli winsec census'\$/tiers='core windows cli census'/" "$script" >"$mutant"
    touched=$(diff "$script" "$mutant" | grep -c '^[<>].*tiers=' || true)
    echo "== roster-mutation-seed: $touched line(s) changed in tiers= (0 would make this seed vacuous)"
    if [ "$touched" -eq 0 ]; then
        echo "   !! seed did not take - refusing to call this a reading"
    else
        mshadow=$(make_shadow "$mutant")
        run_shot census-roster-drift "$mshadow/scripts/portable-tests.sh" "FAKEGO_PIN_FILE=$pin" -- --scope=census
    fi
else
    echo "== census-roster-drift: this script has no tiers= roster (改前 bytes) - the drift is"
    echo "   its NATURAL state there: census prints 'CLAIMED BY ... winsec' for a tier no branch"
    echo "   implements and exits 0. That reading is logs/census-prefix.txt (real go, not a seed)."
fi
