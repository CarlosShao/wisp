#!/usr/bin/env bash
# 329-v1 (verdict leg) mandated replay: named-roster gauge + three/four mutation shots.
# ALL byte edits land in OUT-OF-REPO export trees (/d/W329VA, /d/W329VB), never in tracked files.
# Every shot restores on the spot (and a trap restores at exit); every restore reports sha256.
set -u
REPO="/d/work/workspace/projects plans/Wisp"
LOGS="$REPO/.scratch/wisp/probes/329/v1/logs"
A=/d/W329VA
B=/d/W329VB
GOTMP=/d/wisp-gotmp-329v1
BK="$GOTMP/backup"
mkdir -p "$BK" "$LOGS"

DELIVER_REL="cmd/wisp/panel_snapshot_deliver_329.go"
PUMP_REL="internal/panel/pump.go"
DELIVER="$B/$DELIVER_REL"
PUMP="$B/$PUMP_REL"

cp "$DELIVER" "$BK/deliver329.orig"
cp "$PUMP" "$BK/pump.orig"

H_DELIVER=$(sha256sum "$BK/deliver329.orig" | cut -d' ' -f1)
H_PUMP=$(sha256sum "$BK/pump.orig" | cut -d' ' -f1)

restore_all() {
  cp "$BK/deliver329.orig" "$DELIVER"
  cp "$BK/pump.orig" "$PUMP"
}
trap restore_all EXIT

ROSTER='Test197NoAllowDoorIsReachableFromASubagentsAssembly|Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands|TestBookedSnapshotReachesHostBitForBit329|TestRunBooksWithASnapshotOfItsLiveQueue|TestRunPacketCarriesTheLoadedInstructionFiles|TestRunPacketCarriesTheSubagentItsRosterRowFed|TestRunPacketLightsTheRosterRowOfItsLiveL1Window|TestRunPacketMarksTheRosterRowACardIsHolding|TestRunPacketReportsTheStreamLogPastItsBound|TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound|TestSnapshotWorkspaceSectionReportsTheNarrowing|TestSubagentStreamKeyHasOneMintSite|TestThePumpIsDrivenNotJustAssembled'
TARGET='^TestBookedSnapshotReachesHostBitForBit329$'

run_case() { # $1 tree  $2 -run regex  $3 logfile  $4 label
  local tree=$1 re=$2 log=$3 label=$4
  : > "$log"
  {
    echo "### $label"
    echo "### time: $(date '+%Y-%m-%d %H:%M:%S %z')"
    echo "### cmd: cd $tree && PATH=\$PWD/third_party/sherpa-onnx:\$PWD/build:\$PATH GOTMPDIR=$GOTMP go test ./cmd/wisp/ -count=1 -timeout 30m -v -run '$re'"
  } >> "$log"
  ( cd "$tree" && PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" GOTMPDIR="$GOTMP" go test ./cmd/wisp/ -count=1 -timeout 30m -v -run "$re" ) >> "$log" 2>&1
  local rc=$?
  {
    echo "RC=$rc"
    echo "GAUGE RUN=$(grep -c '^=== RUN' "$log" || true) PASS_TOP=$(grep -c '^--- PASS' "$log" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$log" || true) SKIP=$(grep -c '^--- SKIP' "$log" || true)"
    echo "FAIL_NAMES:"; grep -E '^ *--- FAIL' "$log" || echo "(none)"
    echo "PASS_NAMES:"; grep -E '^--- PASS' "$log" || echo "(none)"
    echo "RED_LEAD_LINES:"; grep -E '^\s+[a-z_/]+\.go:[0-9]+:' "$log" | head -20 || echo "(none)"
    echo
  } >> "$log"
}

# ---------- part 1: the named-roster gauge (AC#7 alternate route, before vs after) ----------
run_case "$A" "$ROSTER" "$LOGS/roster-named-before-feab45e3.log" "AC#7 alternate route: named 329-adjacent roster, tree A (feab45e3, before)"
run_case "$B" "$ROSTER" "$LOGS/roster-named-after-1d1261e0.log"   "AC#7 alternate route: same named roster, tree B (1d1261e0, after)"

# ---------- part 2: mandated mutation shots, tree B physical bytes ----------
echo "PRISTIN_HASH deliver329=$H_DELIVER" > "$LOGS/mutation-hashes.log"
echo "PRISTIN_HASH pump=$H_PUMP" >> "$LOGS/mutation-hashes.log"

shot() { # $1 file  $2 sed expr  $3 marker-token  $4 logfile  $5 label
  local file=$1 expr=$2 token=$3 log=$4 label=$5
  restore_all
  sed -i "$expr" "$file"
  local applied
  applied=$(grep -c "$token" "$file" || true)
  {
    echo "### $label"
    echo "### time: $(date '+%Y-%m-%d %H:%M:%S %z')"
    echo "### mutated file: $file   marker count after mutation = $applied"
  } >> "$log"
  ( cd "$B" && PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" GOTMPDIR="$GOTMP" go test ./cmd/wisp/ -count=1 -timeout 30m -v -run "$TARGET" ) >> "$log" 2>&1
  local rc=$?
  # restore on the spot
  restore_all
  local hd hp
  hd=$(sha256sum "$DELIVER" | cut -d' ' -f1)
  hp=$(sha256sum "$PUMP" | cut -d' ' -f1)
  {
    echo "RC=$rc"
    echo "FAIL_NAMES:"; grep -E '^ *--- FAIL' "$log" || echo "(none)"
    echo "RED_LEAD_LINES:"; grep -E '^\s+[a-z_/]+\.go:[0-9]+:' "$log" | head -20 || echo "(none)"
    echo "RESTORE deliver329=$hd  match_pristine=$([ "$hd" = "$H_DELIVER" ] && echo YES || echo NO)"
    echo "RESTORE pump=$hp  match_pristine=$([ "$hp" = "$H_PUMP" ] && echo YES || echo NO)"
    echo "MUTANT_MARKER_LEFTOVER=$(grep -rc 'MUT329V1' "$DELIVER" "$PUMP" | tr '\n' ' ')"
    echo
  } >> "$log"
  echo "### $label" >> "$LOGS/mutation-hashes.log"
  echo "    RC=$rc  RESTORE deliver329=$hd (match=$([ "$hd" = "$H_DELIVER" ] && echo YES || echo NO))  pump=$hp (match=$([ "$hp" = "$H_PUMP" ] && echo YES || echo NO))" >> "$LOGS/mutation-hashes.log"
}

# M1 de-fang: push an EMPTY payload instead of the booked bytes.
shot "$DELIVER" 's|h\.evalPanelJS(panelSnapshotEvalJS(data))|h.evalPanelJS(panelSnapshotEvalJS(nil)) // MUT329V1-defang-empty|' 'MUT329V1-defang-empty' \
  "$LOGS/mutation-v1-defang-empty.log" "MUTATION 1 (de-fang, empty payload): deliver_329.go push -> panelSnapshotEvalJS(nil)"

# M1x de-fang, byte-sensitive: push the booked bytes minus their LAST byte (same shape, same call, one byte short).
shot "$DELIVER" 's|h\.evalPanelJS(panelSnapshotEvalJS(data))|h.evalPanelJS(panelSnapshotEvalJS(data[:len(data)-1])) // MUT329V1-defang-lastbyte|' 'MUT329V1-defang-lastbyte' \
  "$LOGS/mutation-v1-defang-lastbyte.log" "MUTATION 1x (de-fang, drop last byte): deliver_329.go push -> panelSnapshotEvalJS(data[:len(data)-1])"

# M2 counter-control: marshal side loses one pending card (bytes no longer equal marshal of the booked snapshot).
shot "$PUMP" 's|\tdata, err := json\.Marshal(snap)|\tmutSnap := snap; mutSnap.Pending = nil; data, err := json.Marshal(mutSnap) // MUT329V1-counter-marshal-oneshort|' 'MUT329V1-counter-marshal-oneshort' \
  "$LOGS/mutation-v2-counter-marshal.log" "MUTATION 2 (counter-control, marshal side one pending card short): pump.go Publish marshals a trimmed copy"

# M3: remove the production reader at deliver_329.go:95 (the lastPanelSnapshot read).
shot "$DELIVER" 's|_, data, seen := rt\.lastPanelSnapshot()|var data []byte; seen := false // MUT329V1-reader-removed-at-95|' 'MUT329V1-reader-removed-at-95' \
  "$LOGS/mutation-v3-reader-removal.log" "MUTATION 3 (production reader removed): deliver_329.go line 95 lastPanelSnapshot read neutered"

# ---------- part 3: post-restore control, tree B green again ----------
run_case "$B" "$TARGET" "$LOGS/targeted-green-after-all-mutations.log" "control after all restores: targeted 329 ruler, tree B"
sha256sum "$DELIVER" "$PUMP" "$B/cmd/wisp/panel_pump.go" >> "$LOGS/mutation-hashes.log"
echo "ALL_DONE $(date '+%Y-%m-%d %H:%M:%S %z')" >> "$LOGS/mutation-hashes.log"
