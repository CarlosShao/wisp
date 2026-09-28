#!/bin/sh
# 197-r3b: single-point mutation / negative-control runs, one at a time, each restored
# from HEAD before the next. Readings land next to this script.
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9
export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"
P=.scratch/wisp/probes/197/r3b

restore() {
  git cat-file blob "HEAD:$1" > "$1"
  if [ -z "$(git status --porcelain -- "$1")" ]; then
    echo "RESTORED_CLEAN $1"
  else
    echo "RESTORED_DIRTY $1"
  fi
}

run() { # label  pkg  -run-regexp
  go test -count=1 -run "$3" "$2" > "$P/$1.txt" 2>&1
  rc=$?
  echo "--- $1 rc=$rc"
  grep -E '^(ok|FAIL|--- FAIL|    --- FAIL|\?)' "$P/$1.txt" | head -8
  grep -E '(subagent_carrier_197_test|subagent_roster_197_test)\.go:[0-9]+:' "$P/$1.txt" | head -4
}

CARR='TestRunPacketCarriesTheSubagentItsRosterRowFed|TestRunPacketReportsTheStreamLogPastItsBound|TestSubagentStreamKeyHasOneMintSite'
PANEL5='TestTheRosterReaderPutsSubagentsOnTheWire|TestAPumpWithoutARosterReaderSendsFourKeys|TestAStatusOutsideD43NeverReachesTheWire|TestTruncationFactsRideThePacket|TestDroppedStreamsAreNamedOnTheWire'

echo "=========== BASELINE (committed tree, no mutation) ==========="
run m0-baseline-cmdwisp ./cmd/wisp/ "$CARR"
run m0-baseline-panel ./internal/panel/ "$PANEL5"

echo "=========== M1 unplug the pump reader (cmd/wisp/run.go) ==========="
sed -i 's/Tasks: rt.taskRosterState,/Tasks: nil, \/\/ MUTATION-M1 unplug the roster reader/' cmd/wisp/run.go
grep -n 'MUTATION-M1' cmd/wisp/run.go
run m1-red-cmdwisp ./cmd/wisp/ "$CARR"
run m1-panel-blind ./internal/panel/ "$PANEL5"
restore cmd/wisp/run.go
run m1-restored-cmdwisp ./cmd/wisp/ "$CARR"

echo "=========== M2 keep a refused status on the wire (internal/panel) ==========="
sed -i 's/if !view.StatusKnown {/if false { \/\/ MUTATION-M2 keep the refused name/' internal/panel/subagent_roster_197.go
grep -n 'MUTATION-M2' internal/panel/subagent_roster_197.go
run m2-red-panel ./internal/panel/ TestAStatusOutsideD43NeverReachesTheWire
restore internal/panel/subagent_roster_197.go
run m2-restored-panel ./internal/panel/ TestAStatusOutsideD43NeverReachesTheWire

echo "=========== M3 child keeps no sink (internal/tools) ==========="
sed -i 's/opt.Sink = subagentTextSink{d: t.d}/opt.Sink = nil \/\/ MUTATION-M3 no child sink/' internal/tools/subagent_197.go
grep -n 'MUTATION-M3' internal/tools/subagent_197.go
run m3-red-cmdwisp ./cmd/wisp/ TestRunPacketCarriesTheSubagentItsRosterRowFed
restore internal/tools/subagent_197.go
run m3-restored-cmdwisp ./cmd/wisp/ TestRunPacketCarriesTheSubagentItsRosterRowFed

echo "=========== M4 re-fork the stream key literal (internal/panel/pump.go) ==========="
sed -i 's/SubagentStreamKeyPrefix = streamkey.SubagentPrefix/SubagentStreamKeyPrefix = "subagent:" \/\/ MUTATION-M4 refork/' internal/panel/pump.go
grep -n 'MUTATION-M4' internal/panel/pump.go
run m4-red-cmdwisp ./cmd/wisp/ TestSubagentStreamKeyHasOneMintSite
restore internal/panel/pump.go
run m4-restored-cmdwisp ./cmd/wisp/ TestSubagentStreamKeyHasOneMintSite

echo "=========== FINAL TREE STATE ==========="
git status --porcelain -- cmd internal
echo "FINAL_DONE"
