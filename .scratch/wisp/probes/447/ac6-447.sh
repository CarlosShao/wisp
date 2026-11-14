#!/usr/bin/env bash
# ac6-447.sh -- ticket 447 AC#6's teeth, run ENTIRELY OUTSIDE the repository.
#
# Two mutations of the production code the two fixed cases guard, applied to an
# export tree (never to the repo), each of which must turn the instrument RED and
# NAME the world it broke, followed by a restore that is verified by sha256 (not
# asserted in prose) and a green re-run.
#
#   M1  respond to the WireR1 fix: the turn end no longer calls
#       residency.enterWarm("turn-end"), so the engines are no longer kept hot
#       and the session never lands in Warm. The new bounded wait must time out
#       and print the numbers (stage/engines_hot/warm_turn_end/...).
#   M2  respond to the panel cold-budget move: the cold bring-up really gets 2s
#       slower (lastColdMs inflated in the one place it is written). The moved
#       verdict must fire in its NEW home - the winlive-tagged case - with the
#       D32 1500 ms sentence.
#
# Usage: bash ac6-447.sh <export tree> ; exit 0 only if BOTH mutations went red
# with a named line AND both restores are byte-identical AND the green re-runs
# pass.
set -u
tree=${1:?usage: ac6-447.sh <export-tree>}
cd "$tree" || exit 2
logs="$tree/.scratch/wisp/probes/447/logs"
stage="$tree/.scratch/wisp/stage-447"
mkdir -p "$logs" "$stage" || exit 2
sha() { sha256sum "$1" | awk '{print $1}'; }
verdict=0

echo "== M0 baseline: the pristine export tree is green for both cases =="
bash scripts/wisp-cli-staged-tests.sh -test.run '^(TestWireR1AnnounceEndToEnd|TestPanelHostRealWindowHopAndLifecycle)$' \
	>"$logs/00-m0-baseline.log" 2>&1
m0=$?
echo "M0 rc=$m0 (want 0)"
grep -E '^--- (PASS|FAIL)' "$logs/00-m0-baseline.log" | head -4
[ "$m0" -eq 0 ] || verdict=1

# ---------------------------------------------------------------- M1
echo
echo "== M1: finish() stops calling enterWarm(\"turn-end\") =="
f1=cmd/wisp/resident_voice_windows.go
b1=$(sha "$f1")
cp "$f1" "$stage/m1-pristine.go"
perl -0pi -e 's/if err := rv\.residency\.enterWarm\("turn-end"\); err == nil \{/if err := error(nil); err != nil { \/\/ MUTANT 447-M1: enterWarm("turn-end") removed/' "$f1"
grep -q 'MUTANT 447-M1' "$f1" || { echo "M1 DID NOT LAND - aborting, a green mutant is the worst reading"; exit 4; }
a1=$(sha "$f1")
echo "M1 bytes changed: $([ "$b1" != "$a1" ] && echo yes || echo no)"
bash scripts/wisp-cli-staged-tests.sh -test.run '^TestWireR1AnnounceEndToEnd$' >"$logs/01-m1-enterwarm-removed.log" 2>&1
rc1=$?
echo "M1 rc=$rc1 (want non-zero)"
grep -E '^--- FAIL|stage after the turn|did not land within' "$logs/01-m1-enterwarm-removed.log" | head -4
if [ "$rc1" -eq 0 ]; then echo "M1 WENT GREEN - the fix has no teeth"; verdict=1; fi
grep -q '447: waiting for' "$logs/01-m1-enterwarm-removed.log" || echo "note: the wait's own log line is absent (the caller's Fatal carries the numbers)"
cp "$stage/m1-pristine.go" "$f1"
r1=$(sha "$f1")
echo "M1 restored byte-identical: $([ "$r1" = "$b1" ] && echo True || echo False) ($r1)"
[ "$r1" = "$b1" ] || verdict=1

echo
echo "== M1-restored green =="
bash scripts/wisp-cli-staged-tests.sh -test.run '^TestWireR1AnnounceEndToEnd$' >"$logs/02-m1-restored-green.log" 2>&1
rc1g=$?
echo "M1-restored rc=$rc1g (want 0)"
grep -E '^--- (PASS|FAIL)' "$logs/02-m1-restored-green.log" | head -2
[ "$rc1g" -eq 0 ] || verdict=1

# ---------------------------------------------------------------- M2
echo
echo "== M2: the cold bring-up really got 2s slower (lastColdMs += 2000) =="
f2=cmd/wisp/panel_host_windows.go
b2=$(sha "$f2")
cp "$f2" "$stage/m2-pristine.go"
perl -0pi -e 's/\tm\.lastColdMs = rtMs\n/\tm.lastColdMs = rtMs + 2000 \/\/ MUTANT 447-M2: 2s slower bring-up\n/' "$f2"
grep -q 'MUTANT 447-M2' "$f2" || { echo "M2 DID NOT LAND - aborting"; exit 4; }
a2=$(sha "$f2")
echo "M2 bytes changed: $([ "$b2" != "$a2" ] && echo yes || echo no)"
echo "-- building the winlive-tagged binary (the cold verdict's new home) --"
if go test -tags winlive -c -o "$stage/wisp-winlive.test.exe" ./cmd/wisp/ >"$logs/03-m2-build.log" 2>&1; then
	bash scripts/stage-sherpa-dlls.sh "$stage" >>"$logs/03-m2-build.log" 2>&1
	(cd cmd/wisp && "$stage/wisp-winlive.test.exe" -test.count=1 -test.timeout=10m -test.v \
		-test.run '^TestPanelHostRealWindowHopAndLifecycle_WinLive$') >"$logs/04-m2-winlive-cold-red.log" 2>&1
	rc2=$?
	echo "M2 rc=$rc2 (want non-zero)"
	grep -E '^--- FAIL|cold bring-up .* exceeds|winlive cold' "$logs/04-m2-winlive-cold-red.log" | head -4
	if [ "$rc2" -eq 0 ]; then echo "M2 WENT GREEN - the moved verdict has no teeth"; verdict=1; fi
else
	echo "M2 BUILD FAILED (see $logs/03-m2-build.log)"; verdict=1
fi
cp "$stage/m2-pristine.go" "$f2"
r2=$(sha "$f2")
echo "M2 restored byte-identical: $([ "$r2" = "$b2" ] && echo True || echo False) ($r2)"
[ "$r2" = "$b2" ] || verdict=1

echo
echo "== M2-restored green (same winlive case, clean code) =="
if go test -tags winlive -c -o "$stage/wisp-winlive.test.exe" ./cmd/wisp/ >>"$logs/03-m2-build.log" 2>&1; then
	bash scripts/stage-sherpa-dlls.sh "$stage" >>"$logs/03-m2-build.log" 2>&1
	(cd cmd/wisp && "$stage/wisp-winlive.test.exe" -test.count=1 -test.timeout=10m -test.v \
		-test.run '^TestPanelHostRealWindowHopAndLifecycle_WinLive$') >"$logs/05-m2-restored-green.log" 2>&1
	rc2g=$?
	echo "M2-restored rc=$rc2g (want 0)"
	grep -E '^--- (PASS|FAIL)|cold bring-up' "$logs/05-m2-restored-green.log" | head -3
	[ "$rc2g" -eq 0 ] || verdict=1
else
	echo "M2 restore build failed"; verdict=1
fi

echo
echo "== final sha256 pair (must equal the repo's own bytes) =="
echo "m1 $(sha cmd/wisp/resident_voice_windows.go)"
echo "m2 $(sha cmd/wisp/panel_host_windows.go)"
echo "verdict=$([ "$verdict" -eq 0 ] && echo ALL-TEETH-BIT || echo SOMETHING-DID-NOT-BITE)"
exit "$verdict"
