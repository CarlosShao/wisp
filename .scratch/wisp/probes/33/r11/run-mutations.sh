#!/bin/sh
# 33-r11 mutation runs. Each mutated copy reaches the compiler only through
# -overlay, so no tracked file is touched; the denominator of every run is the
# -run list written out in full, and the rc of each go test lands on its own line.
cd "D:/work/workspace/projects plans/Wisp" || exit 9
export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"
L=.scratch/wisp/probes/33/r11/logs
R=.scratch/wisp/probes/33/r11/mutations
NAMES='TestLockedSuffixNeverTakesTheLockItself33r11|TestRosteredLockedMethodReturnsWhileCallerHoldsTheLock33r11|TestLockedNamingCensusBitesItsOwnFixture33r11|TestLockedNamingCensusRejectsTheOldName33r11|TestCompletedWithinReportsAParkedCall33r11|TestAC13ColdStartPageOverEndsOnEntryContentNotTheProbe'

{
  echo "### CONTROL - overlay lands check: the declaration is renamed inside the overlay copy while the call site still asks for the old name, so a tree that really used the copy cannot compile. rc=1 here is the pass condition; it is what makes every GREEN below a green of the mutated tree rather than of the untouched one."
  echo "### -run denominator: none (build only, go vet)"
  go vet -overlay=$R/overlay-control.json ./cmd/wisp/
  echo "control_rc=$?"
} > "$L/40-mutation-control-undefined-symbol.txt" 2>&1

for tag in m1 m2 m3; do
  {
    case $tag in
      m1) echo "### M1 - the pre-rename name put back on the self-locking round trip (decl + its call site, so it compiles). Expect: unrostered AND self-locking, two red lines from the census case.";;
      m2) echo "### M2 - a brand new Locked-suffixed method that takes no lock at all. Expect: ONLY the roster complaint, which is the exemption list's own tooth.";;
      m3) echo "### M3 - the rostered honest method made to lock itself. Expect: census red AND the behavioural case red because the call parks under the held m.mu.";;
    esac
    echo "### -run denominator: 6 named cases (5 of this leg's own plus 33-r10's page-over case, which must stay green)"
    go test -count=1 -v -timeout 240s -overlay=$R/overlay-$tag.json -run "$NAMES" ./cmd/wisp
    echo "mutation_${tag}_rc=$?"
  } > "$L/41-mut-$tag.txt" 2>&1
done
echo "mutation script finished at $(date '+%Y-%m-%d %H:%M:%S %z')" > "$L/49-mutation-script-done.txt"
