#!/usr/bin/env bash
# 111-r7 g4 (ticket 111 AC#8): the positive self-proof that the COMMITTED matchers
# distinguish "the package was tested" from "the package name appears inside a string".
#
# WHAT IS REAL: both copies of the matcher (portable-tests.sh GUARD B and winsec-tests.sh
# guard 2) are extracted from the committed bytes at run time; the only replaced process
# is `go`, by the committed shim scripts/testdata/portable-tests/go - the same convention
# scripts/portable-tests-selftest.sh uses. Nothing under scripts/ is modified.
#
# Usage: bash g4-ac8-proof.sh <scratch-dir>
set -u -o pipefail
root=$(CDPATH= cd -- "$(dirname -- "$0")/../../../../../.." && pwd)
work=${1:?usage: g4-ac8-proof.sh <scratch-dir>}
mkdir -p "$work/bin"
cp "$root/scripts/testdata/portable-tests/go" "$work/bin/go"
chmod +x "$work/bin/go"

PORTABLE="$root/scripts/portable-tests.sh"
WINSEC="$root/scripts/winsec-tests.sh"
PKG='github.com/CarlosShao/wisp/internal/perm'
OTHER='github.com/CarlosShao/wisp/internal/plugin'

echo "== committed matcher lines (verbatim, with line numbers) =="
grep -n 'escape_re()\|result_tail=\|line=\$(grep -E "\^(ok|FAIL)' "$PORTABLE"
grep -n 'esc_pkgpath=\|result_line=\$(grep -E' "$WINSEC"

# Extract the matcher pieces from the committed bytes (no retyping).
escape_re="$(sed -n '874p' "$PORTABLE")"
result_tail="$(sed -n '884p' "$PORTABLE")"
esc_line="$(sed -n '234p' "$WINSEC")"
res_line="$(sed -n '235p' "$WINSEC")"
eval "$escape_re"
eval "$result_tail"

# ---- Part B: the AC#8-literal shape on a capture we control ------------------
# `false result line for plugin (the control) + two string-only mentions of PKG`.
cap="$work/capture-b.txt"
{
  printf '=== RUN   TestRealThing\n'
  printf -- '--- PASS: TestRealThing (0.01s)\n'
  printf '    fixture_test.go:1: golden fixture recorded under %s (string-only mention)\n' "$PKG"
  printf '    fixture_test.go:2: a typo in a doc string spells githubXcom/CarlosShao/wisp/internal/perm\n'
  printf 'ok  \t%s\t0.012s\n' "$OTHER"
} >"$cap"
echo
echo "== Part B: capture-b.txt =="
cat "$cap"
loose=$(grep -c "$PKG" "$cap" || true)
dotwild=$(grep -c 'github.com/CarlosShao/wisp/internal/perm' "$cap" || true)
newfmt=$(grep -Ec "^(ok|FAIL)[[:space:]]+$(escape_re "$PKG")${result_tail}" "$cap" || true)
l_other=$(grep -E "^(ok|FAIL)[[:space:]]+$(escape_re "$OTHER")${result_tail}" "$cap" | tail -1 || true)
echo "loose grep -c '<PKG>'                          = $loose   (the word is in the log)"
echo "loose grep -c 'github.com/...' (dots are wild) = $dotwild"
echo "committed matcher, tests PKG                   = $newfmt   (string-only mention grants nothing)"
echo "committed matcher, tests OTHER (control)       = [$l_other]"
# the SECOND site, winsec-tests.sh guard 2, judged on the same capture with its own
# variable names (`pkgpath`, `capture`) - both lines are the committed bytes.
pkgpath="$PKG"
capture="$cap"
eval "$esc_line"
eval "$res_line"
echo "winsec guard-2 matcher, tests PKG             = [$result_line]"

# ---- Part A: end-to-end, the REAL script + REAL runtests.sh, go faked --------
awk -v q="'" 'index($0, "core_pin=" q) == 1 { f = 1; next }
    f { if ($0 == q) { f = 0; next } print }' "$PORTABLE" >"$work/core-pin.txt"
printf 'go: downloading %s v0.0.1 (string-only mention, not a result line)\n' "$PKG" >"$work/decoy.txt"

echo
echo "== Part A control: no seeds (PKG must get its own line, rc=0) =="
PATH="$work/bin:$PATH" FAKEGO_PIN_FILE="$work/core-pin.txt" \
  bash "$PORTABLE" --scope=core >"$work/run-a0.log" 2>&1
echo "rc=$?"
grep -n "internal/perm" "$work/run-a0.log" | head -5

echo
echo "== Part A: PKG's result line withheld + its name planted in a string (rc=1, GUARD B names it) =="
PATH="$work/bin:$PATH" FAKEGO_PIN_FILE="$work/core-pin.txt" \
  FAKEGO_DROP_PKG="$PKG" FAKEGO_STDERR_FILE="$work/decoy.txt" \
  bash "$PORTABLE" --scope=core >"$work/run-a1.log" 2>&1
echo "rc=$?"
echo "loose grep -c '<PKG> over the whole log' = $(grep -c "$PKG" "$work/run-a1.log" || true)"
echo "lines naming PKG:"
grep -n "$PKG" "$work/run-a1.log" | sed 's/^/  /'
echo "execution-shaped lines naming PKG (=== RUN/--- PASS/FAIL/SKIP) = $(grep -cE '(=== RUN|--- PASS|--- FAIL|--- SKIP).*'"$PKG" "$work/run-a1.log" || true)"
echo "GUARD B verdict row for PKG:"
grep -n "GUARD B\|$PKG" "$work/run-a1.log" | grep -E 'missing|GUARD B' | sed 's/^/  /' | head -5
