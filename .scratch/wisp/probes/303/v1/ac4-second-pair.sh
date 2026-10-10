#!/usr/bin/env bash
# 303-v1 AC#4 second-pair rig: cmd/wisp full package, TWO runs before / TWO runs after,
# ALL FOUR on the SAME table = repo-external clean clone ~/wisp-303-v1, detached checkouts.
# Ruler copied verbatim from probes/303/orch/r5-r6-package-pair.sh (go test -C <clone> ./cmd/wisp/
# -count=1 -timeout 900s -v, PATH = mother repo's sherpa-onnx + build, shell form /d/...).
set -u
CLONE="$HOME/wisp-303-v1"
OUT="/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/303/v1"
MOTHER="/d/work/workspace/projects plans/Wisp"
export PATH="$MOTHER/third_party/sherpa-onnx:$MOTHER/build:$PATH"

run_one() { # $1=sha $2=slot label
  local sha="$1" label="$2"
  git -C "$CLONE" checkout --detach "$sha" > "$OUT/logs/checkout-$label.txt" 2>&1
  echo "checkout_rc=$? detached=$(git -C "$CLONE" rev-parse --short HEAD) headfile=$(git -C "$CLONE" rev-parse --short HEAD)" >> "$OUT/logs/checkout-$label.txt"
  wc -l < "$CLONE/cmd/wisp/panel_host_windows.go" >> "$OUT/logs/checkout-$label.txt"
  git -C "$CLONE" status --porcelain | wc -l >> "$OUT/logs/checkout-$label.txt"
  echo "start=$(date '+%Y-%m-%d %H:%M:%S %z')" >> "$OUT/logs/full-$label.txt"
  go test -C "$CLONE" ./cmd/wisp/ -count=1 -timeout 900s -v >> "$OUT/logs/full-$label.txt" 2>&1
  echo "rc=$?" >> "$OUT/logs/full-$label.txt"
  echo "end=$(date '+%Y-%m-%d %H:%M:%S %z')" >> "$OUT/logs/full-$label.txt"
}

run_one cf46c24a before-1
run_one cf46c24a before-2
run_one 807497c1 after-1
run_one 807497c1 after-2
echo "ALL_DONE rc=0" >> "$OUT/logs/runner.status"
