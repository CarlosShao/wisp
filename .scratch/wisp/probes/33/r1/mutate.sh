#!/usr/bin/env bash
# 33-r1 mutation vehicle: proves each judgement nail in composer_dispatch_test.go
# can go RED, and that the file is byte-identical afterwards (hash compared).
#
# Zero deletes: the pre-mutation copy is restored with `cp`, never removed, and
# the pristine snapshot stays in this probe directory as the evidence.
set -u
# The repo root is the worktree root, not a relative walk: this script sits five
# components deep (.scratch/wisp/probes/33/r1), and an earlier version of this
# line walked four and pointed every path at .scratch/wisp/internal/panel - which
# made every mutation below a no-op and every reading meaningless.
ROOT="$(git rev-parse --show-toplevel)"
[ -n "$ROOT" ] || { echo "cannot resolve repo root"; exit 1; }
TARGET="$ROOT/internal/panel/composer_dispatch.go"
[ -f "$TARGET" ] || { echo "missing target $TARGET"; exit 1; }
PROBE="$ROOT/.scratch/wisp/probes/33/r1"
LOG="$PROBE/logs"
mkdir -p "$LOG"
cd "$ROOT" || exit 1

cp "$TARGET" "$PROBE/composer_dispatch.go.pristine"
BEFORE="$(sha256sum "$TARGET" | cut -d' ' -f1)"
# A empty hash would make every comparison below vacuously true, so the vehicle
# refuses to run rather than report a clean restore it never checked.
[ "${#BEFORE}" -ge 32 ] || { echo "pristine hash empty/unexpected: $BEFORE"; exit 1; }
echo "== pristine sha256 $BEFORE"

echo "== GREEN run (all 33-r1 nails, unmutated) =="
go test -count=1 -run 'TestRawEnvelope|TestAccepted|TestRoutingDoesNot|TestUnlisted|TestRosterMismatch|TestEveryListed|TestAttachedHandler|TestTamperedSource|TestAttribution|TestDispatcherSpells|TestSliceAAttaches|TestTheInboundHop' ./internal/panel/ >"$LOG/green.txt" 2>&1
echo "rc=$?"
tail -3 "$LOG/green.txt"

# name | sed expression | -run pattern expected to go red
mutate() {
  local tag="$1" expr="$2" pattern="$3"
  echo
  echo "== MUTANT $tag =="
  echo "sed: $expr"
  sed -i "$expr" "$TARGET"
  if ! sha256sum -c <(echo "$BEFORE  $TARGET") >/dev/null 2>&1; then
    echo "mutation applied (hash changed): $(sha256sum "$TARGET" | cut -c1-16)"
  else
    echo "!! MUTANT $tag DID NOT CHANGE THE FILE - vehicle is broken, reading is meaningless"
    return 1
  fi
  go test -count=1 -run "$pattern" ./internal/panel/ >"$LOG/red-$tag.txt" 2>&1
  local rc=$?
  echo "rc=$rc"
  grep -E '^(--- FAIL|FAIL|ok|--- PASS)' "$LOG/red-$tag.txt" | head -20
  echo "FAIL lines: $(grep -c '^\s*--- FAIL' "$LOG/red-$tag.txt")"
  cp "$PROBE/composer_dispatch.go.pristine" "$TARGET"
}

mutate M1 's|return d.Mode.HandleModeRequest(ctx, req)|return nil|' 'TestRawEnvelopeReachesTheAttachedModeHandlerExactlyOnce'
mutate M2 's|return d.rosterMismatch(req)|return nil|' 'TestRosterMismatchBackstopRefusesInsteadOfAccepting'
mutate M3 's|return d.unattached(req)|return nil|g' 'TestEveryListedMethodWithoutAHandlerIsRefusedNotDropped'
# M4 rewrites the correlation key, not the source. That is the honest target:
# ParseComposerRequest pins Source to exactly "panel-composer" before routing, so
# a router that "normalised" Source could not be observed through that field at
# all - parse would have refused the request first. RequestID is the attribution
# field that survives parse and reaches a handler verbatim, so laundering it is
# the shape this nail can actually catch.
mutate M4 's|\tswitch req.Method {|\treq.RequestID = "rid-mutated-by-the-router" // MUTANT M4\n\tswitch req.Method {|' 'TestAttributionFieldsReachTheHandlerUnchanged'

echo
echo "== RESTORE CHECK =="
AFTER="$(sha256sum "$TARGET" | cut -d' ' -f1)"
echo "before=$BEFORE"
echo "after =$AFTER"
if [ "$BEFORE" = "$AFTER" ]; then echo "IDENTICAL: file restored byte for byte"; else echo "!! NOT RESTORED"; fi

echo
echo "== GREEN run again after restore =="
go test -count=1 -run 'TestRawEnvelope|TestAccepted|TestRoutingDoesNot|TestUnlisted|TestRosterMismatch|TestEveryListed|TestAttachedHandler|TestTamperedSource|TestAttribution|TestDispatcherSpells|TestSliceAAttaches|TestTheInboundHop' ./internal/panel/ >"$LOG/green-after.txt" 2>&1
echo "rc=$?"
tail -3 "$LOG/green-after.txt"
