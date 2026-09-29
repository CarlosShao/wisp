#!/bin/sh
# 223-v2 任务一：带 -v 的终态整包读数（全量落盘，绝不接 head/tail）
cd "D:/work/workspace/projects plans/Wisp" || exit 9
OUT=.scratch/wisp/probes/223/v2/full-v.txt
{
  echo "=== START date ==="
  date
  echo "=== anchor ==="
  git rev-parse --short HEAD
  echo "=== PATH export (third_party/sherpa-onnx + build 前置) ==="
} > "$OUT" 2>&1
export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"
go test ./cmd/wisp ./internal/... -count=1 -v -timeout 30m >> "$OUT" 2>&1
echo "GATE_EXIT=$?" >> "$OUT"
{
  echo "=== END date ==="
  date
} >> "$OUT" 2>&1
echo "done, lines:"
wc -l "$OUT"
