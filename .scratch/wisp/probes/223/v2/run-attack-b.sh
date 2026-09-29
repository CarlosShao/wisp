#!/bin/sh
# 223-v2 任务②-B 攻牙：产品那句仍含 needle，但"原因／涉及"半段被整段删掉。
# 读的是"这发到底绿不绿"——绿＝awaitStdout+既有断言没有把 stdout 那半句钉住。
cd "D:/work/workspace/projects plans/Wisp" || exit 9
OUT=.scratch/wisp/probes/223/v2/attack-b-stdout-content.txt
export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"
{
  echo "=== START date ==="; date
  echo "=== anchor ==="; git rev-parse --short HEAD
  echo "=== 突变 B：stdout 那句只剩 \"wisp run: 这些段的改动本次运行不会生效。\"（原因/涉及/两件事都没发生 三段全删，needle 仍在）"
  echo "=== 尺：go test ./cmd/wisp -count=1 -v -run 'TestTicket223RestartTierSaysItWillNotApply'"
} > "$OUT" 2>&1
go test ./cmd/wisp -count=1 -v -timeout 5m -run 'TestTicket223RestartTierSaysItWillNotApply' >> "$OUT" 2>&1
echo "EXIT=$?" >> "$OUT"
{ echo "=== END date ==="; date; } >> "$OUT" 2>&1
wc -l "$OUT"
