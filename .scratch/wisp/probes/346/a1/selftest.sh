#!/usr/bin/env bash
# selftest.sh - ticket 346 AC#3: fixture-driven, offline selftest for
# ci-red-attribution.sh. No gh, no network, no repo pool, no writes outside a
# temp dir. Run from anywhere: cd "$(dirname "$0")" happens below.
#
# Usage: bash selftest.sh            (tests ./ci-red-attribution.sh)
#        TOOL=/path/to/tool.sh bash selftest.sh
#
# Every check below is a reading this ruler must be able to reproduce:
#   - the ticket-face ruler (^--- FAIL on the raw gh log) is BLIND: 0 hits
#   - stripping the job/step/timestamp prefix recovers the real roster
#   - indented subtest lines and mid-line prose never enter the denominator
#   - the same name in two jobs is one unique name, two job readings
#   - attribution is per name: criteria line / file fallback / ledger / none
#   - a clean-looking empty roster is REFUSED (exit 5), not reported green
#   - a ruler whose positive control cannot be found exits 1
#   - gh failure exits 3, no input exits 2, unreadable log exits 4
#   - the run is read-only: fixture hashes are identical before and after
set -uo pipefail
cd "$(dirname "$0")"
TOOL="${TOOL:-./ci-red-attribution.sh}"
POSCTL_NAME="TestC26RewriteAccountIsConsumedAtEverySecurityLeg"
checks=0
failures=0
err="$(mktemp -d 2>/dev/null || echo "./.selftest.$$")"
trap 'rm -rf "$err"' EXIT

ok() { checks=$((checks + 1)); printf 'ok   %s\n' "$1"; }
bad() { checks=$((checks + 1)); failures=$((failures + 1)); printf 'FAIL %s\n' "$1"; }
eq() { if [ "$2" = "$3" ]; then ok "$1"; else bad "$1 (expected [$2], got [$3])"; fi; }
has() { case "$2" in *"$3"*) ok "$1" ;; *) bad "$1 (missing [$3])" ;; esac; }
hasnt() { case "$2" in *"$3"*) bad "$1 (unexpected [$3])" ;; *) ok "$1" ;; esac; }
blocks() { printf '%s\n' "$1" | awk -v want="name=$2" '$0 == want { f = 1 } f { print } f && /^  => / { exit }'; }

# ---- 0. the ruler exists and answers --help -------------------------------
if [ ! -f "$TOOL" ]; then
  bad "0 ruler exists ($TOOL)"
  printf 'selftest: FAIL (%d of %d checks; ruler missing - this is the RED run)\n' "$failures" "$checks"
  exit 1
fi
ok "0 ruler exists ($TOOL)"
hout="$(bash "$TOOL" --help 2>&1)"; hrc=$?
eq "0 --help exit code" 0 "$hrc"
has "0 --help prints usage" "$hout" "Usage:"

# ---- 1. the ticket-face ruler is blind on a real gh-shaped log ------------
# This is the measurement that killed the ticket face's line 18 ruler: gh
# prefixes every line with "job<TAB>step<TAB>timestamp ", so nothing is at
# column 0. A ruler that counts ^--- FAIL reports a clean roster it never read.
naive="$(grep -c '^--- FAIL' fixtures/gh-prefixed.log 2>/dev/null || true)"
eq "1 ticket-face ^--- FAIL ruler is blind (0 hits)" 0 "$naive"

# ---- 2. main attribution run ----------------------------------------------
before="$(sha256sum fixtures/gh-prefixed.log fixtures/all-green.log fixtures/ledger.md fixtures/funcs.tsv fixtures/pool/*.md | sha256sum)"
out="$(bash "$TOOL" \
  --log fixtures/gh-prefixed.log \
  --log fixtures/gh-prefixed.log \
  --issues fixtures/pool \
  --ledger fixtures/ledger.md \
  --func-index fixtures/funcs.tsv \
  --posctl "$POSCTL_NAME" 2>"$err/main.err")"
rc=$?
after="$(sha256sum fixtures/gh-prefixed.log fixtures/all-green.log fixtures/ledger.md fixtures/funcs.tsv fixtures/pool/*.md | sha256sum)"

eq "2 main run exit code" 0 "$rc"
if [ "$rc" -eq 0 ]; then
  has "2 two logs were read" "$out" "logs=2"
  has "2 job x name denominator (identical rows collapse across logs)" "$out" "顶格 FAIL (job x name)=6"
  has "2 unique names and UNREGISTERED count" "$out" "unique red names=5  UNREGISTERED=2"
  has "2 positive control speaks" "$out" "POSITIVE CONTROL ok"

  b="$(blocks "$out" TestFixtureKnownByFile)"
  has "2 KnownByFile: two jobs, sorted" "$b" "jobs=test-core,test-windows"
  has "2 KnownByFile: file fallback resolved" "$b" "file=fixture_owner_test.go"
  has "2 KnownByFile: file-named owner 322" "$b" "具名文件在册票=322"
  has "2 KnownByFile: REGISTERED" "$b" "=> REGISTERED"

  b="$(blocks "$out" TestFixtureKnownByLedger)"
  has "2 KnownByLedger: ledger column" "$b" "台账=A9001"
  has "2 KnownByLedger: REGISTERED" "$b" "=> REGISTERED"

  b="$(blocks "$out" "$POSCTL_NAME")"
  has "2 POSCTL: criteria column points at 93" "$b" "判据在册票=93"
  has "2 POSCTL: REGISTERED" "$b" "=> REGISTERED"

  b="$(blocks "$out" TestFixtureProseOnlyZZZ)"
  has "2 ProseOnly: prose count is informational" "$b" "散文提及=1"
  has "2 ProseOnly: prose alone is NOT an owner" "$b" "=> UNREGISTERED"

  b="$(blocks "$out" TestFixtureUnregisteredXYZ)"
  has "2 Unregistered: no file resolved" "$b" "file=未解析"
  has "2 Unregistered: UNREGISTERED" "$b" "=> UNREGISTERED"

  hasnt "2 indented subtest never enters the roster" "$out" "TestFixtureKnownByFile/subcase"
  hasnt "2 mid-line prose never enters the roster" "$out" "TestFixtureNotAtLineStart"
  eq "2 run is read-only (fixture hashes unchanged)" "$before" "$after"
else
  bad "2 main run produced no output to assert on (see $err/main.err)"
fi

# ---- 3. an empty roster is refused, never reported green ------------------
empty_err="$err/empty.err"; empty_out="$(bash "$TOOL" --log fixtures/all-green.log \
  --issues fixtures/pool --ledger fixtures/ledger.md --func-index fixtures/funcs.tsv \
  --posctl "$POSCTL_NAME" 2>"$empty_err")"; erc=$?
eq "3 empty roster exit code is 5" 5 "$erc"
has "3 empty roster is refused out loud" "$(cat "$empty_err")" "REFUSING"
hasnt "3 empty roster never prints a clean summary" "$empty_out" "unique red names=0"

# ---- 4. positive control fails before attribution ------------------------
pc_out="$(bash "$TOOL" --log fixtures/gh-prefixed.log \
  --issues fixtures/pool --ledger fixtures/ledger.md --func-index fixtures/funcs.tsv \
  --posctl TestNotInAnyPool 2>&1)"; prc=$?
eq "4 missing positive control exits before attribution" 1 "$prc"
has "4 positive control failure is named" "$pc_out" "POSITIVE CONTROL FAILED"

# ---- 5. usage / input errors ---------------------------------------------
bash "$TOOL" >/dev/null 2>"$err/noinput.err"; nrc=$?
eq "5 no input exits 2" 2 "$nrc"
bash "$TOOL" --log fixtures/does-not-exist.log --issues fixtures/pool \
  --ledger fixtures/ledger.md >/dev/null 2>"$err/nolog.err"; lrc=$?
eq "5 unreadable log exits 4" 4 "$lrc"
has "5 unreadable log is named" "$(cat "$err/nolog.err")" "unreadable log"

# ---- 6. gh failure is a reading, not a green ----------------------------
mkdir -p fixtures/bin
chmod +x fixtures/bin/gh 2>/dev/null || true
PATH="$PWD/fixtures/bin:$PATH" bash "$TOOL" 999999 --issues fixtures/pool \
  --ledger fixtures/ledger.md >/dev/null 2>"$err/gh.err"; grc=$?
eq "6 gh failure exits 3" 3 "$grc"
has "6 gh failure is named with the run id" "$(cat "$err/gh.err")" "gh run view 999999 FAILED"

# ---- 7. summary ----------------------------------------------------------
if [ "$failures" -eq 0 ]; then
  printf 'selftest: PASS (%d checks)\n' "$checks"
  exit 0
fi
printf 'selftest: FAIL (%d of %d checks)\n' "$failures" "$checks"
exit 1
