#!/bin/sh
# 197-r4: single-point mutation / positive-control runs for AC#4's blockedOnApproval
# judgement. One mutation at a time, each restored from HEAD before the next.
# Readings land next to this script.
#
# These are -run filtered runs: they are INTERMEDIATE readings for the mutation
# table. The package-level "green" claims live in final-v-*.txt (whole package, -v).
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9
export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"
P=.scratch/wisp/probes/197/r4

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
  grep -E 'subagent_blocked_197_test\.go:[0-9]+:' "$P/$1.txt" | head -6
}

BLOCKED='TestRunPacketMarksTheRosterRowACardIsHolding'
CARR='TestRunPacketCarriesTheSubagentItsRosterRowFed|TestRunPacketReportsTheStreamLogPastItsBound|TestSubagentStreamKeyHasOneMintSite'
PANEL5='TestTheRosterReaderPutsSubagentsOnTheWire|TestAPumpWithoutARosterReaderSendsFourKeys|TestAStatusOutsideD43NeverReachesTheWire|TestTruncationFactsRideThePacket|TestDroppedStreamsAreNamedOnTheWire'

echo "=========== BASELINE (committed tree + this case, no mutation) ==========="
run m0-baseline-cmdwisp ./cmd/wisp/ "$BLOCKED|$CARR"
run m0-baseline-panel ./internal/panel/ "$PANEL5"

echo "=========== M1 unplug the join's reader: the pump stops handing the cards over ==========="
# internal/panel/pump.go: the roster section is still built, still carries rows, and
# is now built over an empty card list -> blockedOnApproval can never be true.
sed -i 's/snap.Tasks = TaskRosterSectionFrom(p.src.Tasks(), cards)/snap.Tasks = TaskRosterSectionFrom(p.src.Tasks(), nil) \/\/ MUTATION-M1/' internal/panel/pump.go
grep -n 'MUTATION-M1' internal/panel/pump.go
run m1-red-cmdwisp ./cmd/wisp/ "$BLOCKED"
run m1-red-panel ./internal/panel/ "$PANEL5"
restore internal/panel/pump.go
run m1-restored-cmdwisp ./cmd/wisp/ "$BLOCKED"

echo "=========== M2 the flag hard-off (constant false) ==========="
sed -i 's/BlockedOnApproval: row.TaskID != "" \&\& waiting\[row.TaskID\],/BlockedOnApproval: false, \/\/ MUTATION-M2/' internal/panel/subagent_roster_197.go
grep -n 'MUTATION-M2' internal/panel/subagent_roster_197.go
run m2-red-cmdwisp ./cmd/wisp/ "$BLOCKED"
run m2-red-panel ./internal/panel/ "$PANEL5"
restore internal/panel/subagent_roster_197.go
run m2-restored-cmdwisp ./cmd/wisp/ "$BLOCKED"

echo "=========== M3 the flag hard-on (constant true) ==========="
# This is the 空转 shape ticket 181 AC#7 was filed over: a field that prints one
# word whatever happens. The case's SECOND reading (after the card left the queue)
# and the root-row reading are what catch it.
sed -i 's/BlockedOnApproval: row.TaskID != "" \&\& waiting\[row.TaskID\],/BlockedOnApproval: true, \/\/ MUTATION-M3/' internal/panel/subagent_roster_197.go
grep -n 'MUTATION-M3' internal/panel/subagent_roster_197.go
run m3-red-cmdwisp ./cmd/wisp/ "$BLOCKED"
restore internal/panel/subagent_roster_197.go
run m3-restored-cmdwisp ./cmd/wisp/ "$BLOCKED"

echo "=========== M4 unplug the whole carrier at the assembly root ==========="
sed -i 's/Tasks: rt.taskRosterState,/Tasks: nil, \/\/ MUTATION-M4/' cmd/wisp/run.go
grep -n 'MUTATION-M4' cmd/wisp/run.go
run m4-red-cmdwisp ./cmd/wisp/ "$BLOCKED"
restore cmd/wisp/run.go
run m4-restored-cmdwisp ./cmd/wisp/ "$BLOCKED"

echo "=========== M5 the card names a different key (the row must NOT flip) ==========="
# Test-side: the queue still holds exactly one live card, but its correlation id is
# not this row's task id. If blockedOnApproval were "any card anywhere", this would
# stay green; it must go red.
sed -i 's/TaskID: id, CorrelationID: id, Tool: blocked197Tool,/TaskID: id, CorrelationID: id + "-elsewhere", Tool: blocked197Tool, \/\/ MUTATION-M5/' cmd/wisp/subagent_blocked_197_test.go
grep -n 'MUTATION-M5' cmd/wisp/subagent_blocked_197_test.go
run m5-red-cmdwisp ./cmd/wisp/ "$BLOCKED"
restore cmd/wisp/subagent_blocked_197_test.go
run m5-restored-cmdwisp ./cmd/wisp/ "$BLOCKED"

echo "=========== FINAL TREE STATE (tracked production + this case) ==========="
git status --porcelain -- cmd internal
echo "FINAL_DONE"
