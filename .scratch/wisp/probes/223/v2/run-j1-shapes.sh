#!/bin/sh
# 223-v2 任务③：J1 三形（＋对照形）——台件只读数不判等，走 overlay，零源码改动。
cd "D:/work/workspace/projects plans/Wisp" || exit 9
OUT=.scratch/wisp/probes/223/v2/j1-three-shapes.txt
export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"
{
  echo "=== START date ==="; date
  echo "=== anchor ==="; git rev-parse --short HEAD
  echo "=== 台件＝overlay 虚拟件 cmd/wisp/zz_v2probe_cmdwisp_test.go（物理件 .scratch/wisp/probes/223/v2/overlay/），只 t.Logf 不判等"
  echo "=== 尺：go test -overlay .scratch/wisp/probes/223/v2/overlay/overlay.json ./cmd/wisp -count=1 -v -run TestV2ProbeJ1Shapes"
} > "$OUT" 2>&1
go test -overlay .scratch/wisp/probes/223/v2/overlay/overlay.json ./cmd/wisp -count=1 -v -timeout 5m -run TestV2ProbeJ1Shapes >> "$OUT" 2>&1
echo "EXIT=$?" >> "$OUT"
{ echo "=== END date ==="; date; } >> "$OUT" 2>&1
wc -l "$OUT"
