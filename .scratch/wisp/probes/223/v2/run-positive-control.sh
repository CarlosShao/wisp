#!/bin/sh
# 223-v2 任务②-A 正控：产品那句 stdout Fprintf 被整段注释掉后，
# TestTicket223RestartTierSaysItWillNotApply 必须红。
cd "D:/work/workspace/projects plans/Wisp" || exit 9
OUT=.scratch/wisp/probes/223/v2/positive-control-v2.txt
export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"
{
  echo "=== START date ==="; date
  echo "=== anchor ==="; git rev-parse --short HEAD
  echo "=== 突变说明：cmd/wisp/config_reload.go 的 reportRestartPending stdout Fprintf 整段注释掉（两行 stderr 审计保留）"
  echo "=== 尺：go test ./cmd/wisp -count=1 -v -run 'TestTicket223RestartTierSaysItWillNotApply'"
} > "$OUT" 2>&1
go test ./cmd/wisp -count=1 -v -timeout 5m -run 'TestTicket223RestartTierSaysItWillNotApply' >> "$OUT" 2>&1
echo "EXIT=$?" >> "$OUT"
echo "=== END date ==="; date >> "$OUT" 2>&1
wc -l "$OUT"
