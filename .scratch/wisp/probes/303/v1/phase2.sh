#!/usr/bin/env bash
# 303-v1 phase 2: (a) gates re-run by the NON-implementer, (b) the pre-regression nail
# re-check that ticket 305's existence hinges on, (c) the ⓐ attack rig.
# Waits for the package pair (logs/runner.status = ALL_DONE) so nothing real-window-heavy
# runs concurrently with it.
set -u
MOTHER="/d/work/workspace/projects plans/Wisp"
OUT="$MOTHER/.scratch/wisp/probes/303/v1"
LOG="$OUT/logs"
CLONE="$HOME/wisp-303-v1"
export PATH="$MOTHER/third_party/sherpa-onnx:$MOTHER/build:$PATH"

while ! grep -q ALL_DONE "$LOG/runner.status" 2>/dev/null; do sleep 20; done

# ---------- (a) gates, run by this leg, in the MOTHER repo (read-only for files) -------
{ echo "### go vet ./cmd/wisp/  (run in $MOTHER, cwd mother)"; cd "$MOTHER" || exit 1;
  sh -c "go vet ./cmd/wisp/" 2>&1; echo "rc=$?"; } > "$LOG/gate-vet-plain.txt" 2>&1
{ echo "### go vet -tags winlive ./cmd/wisp/ ./internal/ball/"; cd "$MOTHER" || exit 1;
  sh -c "go vet -tags winlive ./cmd/wisp/ ./internal/ball/" 2>&1; echo "rc=$?"; } > "$LOG/gate-vet-winlive.txt" 2>&1
{ echo "### sh scripts/d22scan.sh"; cd "$MOTHER" || exit 1;
  sh scripts/d22scan.sh 2>&1; echo "rc=$?"; } > "$LOG/gate-d22scan.txt" 2>&1
{ echo "### gofmt -l cmd/wisp  (WORKING TREE ruler, mother repo)"; cd "$MOTHER" || exit 1;
  gofmt -l cmd/wisp 2>&1; echo "rc=$?"; } > "$LOG/gate-gofmt-worktree.txt" 2>&1
{ echo "### gofmt -l cmd/wisp  (HEAD BLOB ruler: git archive HEAD | tar -x into \$HOME)";
  T="$HOME/wisp-303-v1-blob"; rm -rf "$T"; mkdir -p "$T"; (cd "$MOTHER" && git archive HEAD | tar -x -C "$T");
  cd "$T" || exit 1; gofmt -l cmd/wisp 2>&1; echo "rc=$?"; } > "$LOG/gate-gofmt-headblob.txt" 2>&1
{ echo "### whole-repo gofmt residual (context only, NOT this leg's range)"; cd "$MOTHER" || exit 1;
  gofmt -l . 2>&1 | head -12; echo "total=$(gofmt -l . 2>/dev/null | wc -l)"; echo "rc=0"; } > "$LOG/gate-gofmt-repocontext.txt" 2>&1

# ---------- (b) the pre-regression table check: three nails at f718e9b6 (parent of the culprit) --
git -C "$CLONE" checkout --detach f718e9b6 > "$LOG/checkout-prereg.txt" 2>&1
echo "checkout_rc=$? detached=$(git -C "$CLONE" rev-parse --short HEAD)" >> "$LOG/checkout-prereg.txt"
{ echo "### three nails at f718e9b6 (PARENT of fb2fb802, i.e. BEFORE the regression), clone table";
  go test -C "$CLONE" ./cmd/wisp/ -count=1 -timeout 420s -v \
    -run 'TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe|TestAC14AwaitedBindingReplyReachesThePage|TestAC14GoSideEvalPushReachesThePage' 2>&1
  echo "rc=$?"; } > "$LOG/nails-at-f718e9b6.txt" 2>&1

# ---------- (c) ⓐ attack rig -------------------------------------------------------------
bash "$OUT/attack-rig.sh"
