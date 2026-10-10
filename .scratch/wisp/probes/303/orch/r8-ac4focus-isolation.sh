#!/usr/bin/env bash
# Orchestrator: is TestAC4FocusReturnToPriorWindowGap33r5's red in the after-run caused by
# the fix, or is it desktop-state dependent? Same surface (external clone), targeted, interleaved.
set -u
CL="$HOME/wisp-303-bisect"
REPO="/d/work/workspace/projects plans/Wisp"
OUT="$REPO/.scratch/wisp/probes/303/orch"
NAME=TestAC4FocusReturnToPriorWindowGap33r5
run() { # $1=sha $2=label $3=count
  git -C "$CL" checkout -q --detach "$1"
  local f="$OUT/r8-$2-count$3.txt"
  echo "sha=$(git -C "$CL" log -1 --format='%h') label=$2 count=$3 start=$(date '+%T')" > "$f"
  PATH="$REPO/third_party/sherpa-onnx:$REPO/build:$PATH" go test -C "$CL" ./cmd/wisp/ -count="$3" -timeout 300s -v -run "$NAME" >> "$f" 2>&1
  echo "rc=$?" >> "$f"
  echo "  colors: $(grep -cE "^--- PASS: $NAME" "$f") pass / $(grep -cE "^--- FAIL: $NAME" "$f") fail" >> "$f"
}
run 807497c1 head-fix 5
run cf46c24a parent-nofix 5
run 807497c1 head-fix-again 3
git -C "$CL" checkout -q --detach "$(git -C "$REPO" rev-parse origin/dev)"
echo "restored=$(git -C "$CL" log -1 --format='%h')"
echo "end=$(date '+%T')"
