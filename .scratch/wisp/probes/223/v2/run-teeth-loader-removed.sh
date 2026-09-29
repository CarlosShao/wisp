#!/bin/sh
# 223-v2 任务②-C：把 internal/config/loader.go 换回 223-r2 之前（起手锚点 c774f8da）的形状，
# 跑 TestTicket223R2FailureSentenceRouting，逐名读红册——实现腿自述"恰红 5 形（C1/E1/G1/L1/J1）、3 对照形不动"。
cd "D:/work/workspace/projects plans/Wisp" || exit 9
OUT=.scratch/wisp/probes/223/v2/teeth-loader-removed-v2.txt
export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"
{
  echo "=== START date ==="; date
  echo "=== anchor ==="; git rev-parse --short HEAD
  echo "=== loader.go 突变前基线哈希 ==="
  certutil -hashfile internal/config/loader.go SHA256
  echo "=== 突变命令：git cat-file blob c774f8da:internal/config/loader.go > internal/config/loader.go"
} > "$OUT" 2>&1
git cat-file blob c774f8da:internal/config/loader.go > internal/config/loader.go
{
  echo "=== 换回后哈希（应与 HEAD 不同、与 c774f8da 相同） ==="
  certutil -hashfile internal/config/loader.go SHA256
  echo "=== 尺：go test ./cmd/wisp -count=1 -v -run 'TestTicket223R2FailureSentenceRouting'"
} >> "$OUT" 2>&1
go test ./cmd/wisp -count=1 -v -timeout 5m -run 'TestTicket223R2FailureSentenceRouting' >> "$OUT" 2>&1
echo "EXIT=$?" >> "$OUT"
{
  echo "=== 还原命令：git cat-file blob 42220c9c4dc4c090b3749ae2bc471d3be86bdfef > internal/config/loader.go ==="
} >> "$OUT" 2>&1
echo "END_STEP1" >> "$OUT"
wc -l "$OUT"
