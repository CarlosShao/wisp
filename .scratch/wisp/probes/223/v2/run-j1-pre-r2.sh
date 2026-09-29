#!/bin/sh
# 223-v2 任务③补：把 loader.go 换回 223-r2 之前的形状（blob c774f8da），
# 再跑同一枚 J1 台件，读 (b)(c) 两形改前到底报哪一句——用来判定"新注是不是 r2 改出来的回归"。
cd "D:/work/workspace/projects plans/Wisp" || exit 9
OUT=.scratch/wisp/probes/223/v2/j1-shapes-pre-r2-loader.txt
export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"
{
  echo "=== START date ==="; date
  echo "=== anchor ==="; git rev-parse --short HEAD
  echo "=== loader.go 基线（改法后）哈希 ==="; certutil -hashfile internal/config/loader.go SHA256
  echo "=== 换回 c774f8da 形状后哈希 ==="
  git cat-file blob c774f8da:internal/config/loader.go > internal/config/loader.go
  certutil -hashfile internal/config/loader.go SHA256
  echo "=== 尺：同 ③ 节那条 overlay 台件，只换 loader.go 的形状 ==="
} > "$OUT" 2>&1
go test -overlay .scratch/wisp/probes/223/v2/overlay/overlay.json ./cmd/wisp -count=1 -v -timeout 5m -run TestV2ProbeJ1Shapes >> "$OUT" 2>&1
echo "EXIT=$?" >> "$OUT"
{
  echo "=== 还原为 HEAD blob（42220c9c…）＋哈希 ==="
  git cat-file blob 42220c9c4dc4c090b3749ae2bc471d3be86bdfef > internal/config/loader.go
  certutil -hashfile internal/config/loader.go SHA256
  echo "=== 收工 git status（应为空） ==="; git status --porcelain -- cmd/wisp internal/config
  echo "=== END date ==="; date
} >> "$OUT" 2>&1
wc -l "$OUT"
