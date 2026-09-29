#!/bin/sh
# 223-v2：本腿自己复量那枚曾间歇红的用例 10 发（同一把尺：单包、无 -v、只看 exit 码）。
cd "D:/work/workspace/projects plans/Wisp" || exit 9
OUT=.scratch/wisp/probes/223/v2/flake-recheck-10.txt
export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"
{ echo "=== START ==="; date; echo "尺：go test ./cmd/wisp -count=1 -run 'TestTicket223RestartTierSaysItWillNotApply' × 10"; } > "$OUT" 2>&1
FAILS=0
for i in $(seq 1 10); do
  go test ./cmd/wisp -count=1 -timeout 5m -run 'TestTicket223RestartTierSaysItWillNotApply' >> "$OUT" 2>&1
  E=$?
  echo "shot $i/10 EXIT=$E" >> "$OUT"
  [ "$E" -ne 0 ] && FAILS=$((FAILS+1))
done
{ echo "TOTAL=10 FAILS=$FAILS"; echo "=== END ==="; date; } >> "$OUT" 2>&1
tail -3 "$OUT"
