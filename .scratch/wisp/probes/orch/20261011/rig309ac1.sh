#!/usr/bin/env bash
# orchestrator's OWN counter-shape rig for ticket 309 AC#1 (⛔ adopting the leg's table untested)
# two directions the cell demands + a pristine control; every shot restored and byte-proved.
set -u
R="D:/work/workspace/projects plans/Wisp"
D="$R/.scratch/wisp/probes/orch/20261011"
L="$D/logs/orch-A851-309ac1-rig.txt"
F="internal/tools/bridge.go"
BK="/d/tmp/wisp-orch-20261011/bridge.go.pre-309ac1"
: > "$L"
cd "$R" || exit 1
TEST='Test309AssessorForEmptyTaskIDDropsTaintDetector'

cp "$F" "$BK"
SHA_PRE=$(sha256sum "$BK" | cut -c1-16)
{ echo "=== pre: sha256[:16]=$SHA_PRE bytes=$(wc -c <"$F") HEAD=$(git rev-parse --short HEAD)"
  echo "guard_line=$(grep -n 'b.injected || b.prov == nil' "$F" | cut -c1-90)"
  echo "attach_line=$(grep -n 'WithTaintDetector(b.prov.Detector(taskID))' "$F" | cut -c1-90)"; } >> "$L"

shot() { # $1 name, $2 old, $3 new
  local name="$1" old="$2" new="$3"
  { echo "--- $name"; echo "old=[$old]"; echo "new=[$new]"; } >> "$L"
  if ! grep -qF "$old" "$F"; then echo "ABORT: old string not found (mutation would be a no-op = fake green)" >> "$L"; return 1; fi
  awk -v o="$old" -v n="$new" '{ i=index($0,o); if(i>0){ $0=substr($0,1,i-1) n substr($0,i+length(o)) } print }' "$F" > "$F.tmp" && mv "$F.tmp" "$F"
  echo "mutated_lines=$(git diff --numstat -- "$F" | tr '\n' ' ')" >> "$L"
  GOFLAGS= go test ./internal/tools/ -run "$TEST" -count=1 > "$D/logs/orch-A851-$name.out" 2>&1
  local rc=$?
  { echo "rc=$rc"; grep -E '^(--- FAIL|--- PASS|ok|FAIL|PASS)' "$D/logs/orch-A851-$name.out" | head -3
    grep -m1 -E '309|detector|taint' "$D/logs/orch-A851-$name.out" | cut -c1-170; } >> "$L"
  cp "$BK" "$F"
  if cmp -s "$BK" "$F"; then echo "restored=IDENTICAL sha256[:16]=$(sha256sum "$F" | cut -c1-16)" >> "$L"; else echo "RESTORE FAILED !!!" >> "$L"; fi
  echo "diff_after_restore_lines=$(git diff --numstat -- "$F" | grep -c '' )" >> "$L"
}

# control: pristine tree must be green (a rig without this control proves nothing)
GOFLAGS= go test ./internal/tools/ -run "$TEST" -count=1 > "$D/logs/orch-A851-control.out" 2>&1
{ echo "=== CONTROL (pristine) rc=$?"; grep -E '^(--- PASS|--- FAIL|ok|FAIL|PASS)' "$D/logs/orch-A851-control.out" | head -2; } >> "$L"

shot SHOT1-drop-empty-limb 'if b.injected || b.prov == nil || taskID == "" {' 'if b.injected || b.prov == nil {'
shot SHOT2-detector-not-attached 'WithTaintDetector(b.prov.Detector(taskID))' 'WithTaintDetector(nil)'

{ echo "=== final proof"
  echo "cmp_with_backup=$(cmp -s "$BK" "$F" && echo IDENTICAL || echo DIFF)"
  echo "sha256_now=$(sha256sum "$F" | cut -c1-16) (pre=$SHA_PRE)"
  echo "worktree_diff_bridge=$(git diff --numstat -- "$F" | tr '\n' ' ')"
  git show "HEAD:$F" > /d/tmp/wisp-orch-20261011/bridge.go.head.blob
  echo "cmp_with_HEAD_blob=$(cmp -s /d/tmp/wisp-orch-20261011/bridge.go.head.blob "$F" && echo IDENTICAL || echo DIFF)"
  echo "dirty_internal_cmd_tools=$(git status --porcelain -- internal cmd tools | wc -l | tr -d ' ')"
  echo "staged=$(git diff --cached --name-only | grep -c '' || true)"; } >> "$L"
cat "$L"
