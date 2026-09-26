#!/usr/bin/env bash
# gate shape = CI-同形 for cmd/wisp: dll dir on PATH (MSYS form), full package, no -run,
# no -cover, no -overlay.  One run per sha, sequentially (they must not share CPU).
set -u
base=/d/work/tmp/wisp152-accept-r1
for k in post pre; do
  cd "$base/snap-$k"
  export PATH="$base/snap-$k/third_party/sherpa-onnx:$PATH"
  echo "### $k HEAD=$(git -C '/d/work/workspace/projects plans/Wisp' rev-parse --short 97cfc6e) snap=$k" > /dev/null
  go test -count=1 -v ./cmd/wisp/ > "$base/logs/gate-$k.log" 2>&1
  echo "gate-$k rc=$?"
done
