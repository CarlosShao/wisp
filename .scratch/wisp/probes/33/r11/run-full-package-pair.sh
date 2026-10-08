#!/bin/sh
# 33-r11: the full-package pair - RUN A is the pristine HEAD state reached only
# through -overlay (no tracked file is touched, nothing to revert), RUN B is the real
# worktree after this leg's rename plus guard. Both run to completion; the rc of each
# go test lands on its own line in its own log.
cd "D:/work/workspace/projects plans/Wisp" || exit 9
export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"
L=.scratch/wisp/probes/33/r11/logs
{
  echo "### RUN A - pristine HEAD via -overlay .scratch/wisp/probes/33/r11/mutations/overlay-baseline.json (guard test removed from the build)"
  echo "### command: go test -count=1 -timeout 500s -overlay=... ./cmd/wisp"
  go test -count=1 -timeout 500s -overlay=.scratch/wisp/probes/33/r11/mutations/overlay-baseline.json ./cmd/wisp
  echo "baseline_rc=$?"
} > "$L/30-baseline-full-package.txt" 2>&1
{
  echo "### RUN B - post-change, real worktree"
  echo "### command: go test -count=1 -timeout 500s ./cmd/wisp"
  go test -count=1 -timeout 500s ./cmd/wisp
  echo "postchange_rc=$?"
} > "$L/31-postchange-full-package.txt" 2>&1
echo "pair script finished at $(date '+%Y-%m-%d %H:%M:%S %z')" > "$L/32-pair-script-done.txt"
