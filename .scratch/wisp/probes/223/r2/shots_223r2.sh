#!/bin/sh
# 223-r2 shot loop: run the AC#7 test N times with the exact yard the ticket
# names (PATH first, then go test -count=1 -run), log EVERY shot's exit code
# full-fidelity (never piped through head/tail).
# usage: sh shots_223r2.sh <N>
cd "/d/work/workspace/projects plans/Wisp" || exit 99
export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"
N="$1"
i=1
fails=0
while [ "$i" -le "$N" ]; do
  echo "=== shot $i/$N start=$(date)"
  go test ./cmd/wisp -count=1 -run 'TestTicket223RestartTierSaysItWillNotApply' 2>&1
  code=$?
  echo "=== shot $i/$N EXIT=$code"
  if [ "$code" -ne 0 ]; then
    fails=$((fails+1))
  fi
  i=$((i+1))
done
echo "TOTAL=$N FAILS=$fails"
