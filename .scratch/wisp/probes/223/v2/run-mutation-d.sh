#!/bin/sh
# 223-v2 任务②-D：复现实现腿 m2（本腿没有转述它，自己重跑）——
# 把生产钩子赋值点 cmd/wisp/config_reload.go:114 那一行注释掉，
# 跑 AC#3 那两枚"放宽必带 L2 复确认"的具名用例，必须红。
cd "D:/work/workspace/projects plans/Wisp" || exit 9
OUT=.scratch/wisp/probes/223/v2/mutation-d-no-confirm-hook.txt
export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"
{
  echo "=== START date ==="; date
  echo "=== anchor ==="; git rev-parse --short HEAD
  echo "=== 突变前 config_reload.go 基线哈希 ==="; certutil -hashfile cmd/wisp/config_reload.go SHA256
  echo "=== 突变＝注释掉 :114 rt.mgr.ConfirmLocked = rt.confirmLockedLoosening（:115 的 OnRestartPending 与 :119 的 tick 保留）"
  echo "=== 尺：go test ./cmd/wisp -count=1 -v -timeout 8m -run 'TestTicket223HandEditedFsLooseningCostsAnL2Card|TestTicket223ModeLooseningChangesTheRunningModeAfterAllow'"
} > "$OUT" 2>&1
go build ./cmd/wisp/ >> "$OUT" 2>&1
echo "BUILD_EXIT=$?" >> "$OUT"
go test ./cmd/wisp -count=1 -v -timeout 8m -run 'TestTicket223HandEditedFsLooseningCostsAnL2Card|TestTicket223ModeLooseningChangesTheRunningModeAfterAllow' >> "$OUT" 2>&1
echo "EXIT=$?" >> "$OUT"
{ echo "=== END date ==="; date; } >> "$OUT" 2>&1
wc -l "$OUT"
