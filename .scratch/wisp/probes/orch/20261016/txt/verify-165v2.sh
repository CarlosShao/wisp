#!/usr/bin/env bash
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9
ANCHOR=$(git rev-parse --short HEAD)
TREE=/d/tmp/orch-165v2-verify-a
mkdir -p "$TREE" "$TREE/logs" || exit 9
echo "anchor=$ANCHOR"
git archive "$ANCHOR" go.mod go.sum internal cmd tools scripts frontend models | tar -x -C "$TREE" || exit 8
mkdir -p "$TREE/build" "$TREE/third_party/sherpa-onnx" "$TREE/frontend/dist"
cp build/*.dll "$TREE/build/" 2>/dev/null; cp third_party/sherpa-onnx/*.dll "$TREE/third_party/sherpa-onnx/" 2>/dev/null
cp -r frontend/dist/. "$TREE/frontend/dist/" 2>/dev/null
echo "go_files=$(find "$TREE" -name '*.go' | wc -l)"
mkdir -p /d/tmp/orch-gotmp || exit 7
cd "$TREE" || exit 9
export GOTMPDIR='D:\tmp\orch-gotmp'
{
  echo "=== forward targeted run"
  go test ./internal/agent/ -count=1 -timeout 30m -v -run '^TestPlanPresent' 2>&1 | grep -E '^(=== RUN|--- (PASS|FAIL|SKIP)|ok|FAIL|panic)' | head -20
  echo "rc_forward=${PIPESTATUS[0]}"
  echo "=== whole package -v denominators"
  go test ./internal/agent/ -count=1 -timeout 30m -v > $TREE/logs/full-v.txt 2>&1
  echo "rc_full=$?"
} 2>&1 | tee "$TREE/logs/verify.txt"
echo "RUN=$(grep -c '^=== RUN' $TREE/logs/full-v.txt)"
echo "PASS_TOP=$(grep -c '^--- PASS' $TREE/logs/full-v.txt)"
echo "FAIL_TOP=$(grep -c '^--- FAIL' $TREE/logs/full-v.txt)"
echo "SKIP=$(grep -c '^--- SKIP' $TREE/logs/full-v.txt)"
echo "PANIC=$(grep -c '^panic' $TREE/logs/full-v.txt)"
sha256sum internal/agent/testdata/golden/plan-present.sse internal/agent/plan_present_doorbell_165r3_test.go
