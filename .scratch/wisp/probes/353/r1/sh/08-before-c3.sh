#!/bin/sh
# 353-r1 随笔③前自检：gofumpt 只验自己那枚（⛔ -w 扫别人文件；本腿只对**自己那两枚**用过 -w，此处逐字登记）
# ＋ vet ＋★工作树脏枚数盘点（⛔ 把别腿脏文件记进自己笔里＝公共块 共享 index 那一节）。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 9
unset GOTMPDIR TMPDIR TEMP TMP
L=/d/tmp/wisp/20261020/353r1/logs
{
  echo "== @ $(date '+%F %T %z') HEAD=$(git log -1 --format=%h) =="
  echo "== gofumpt -l 只列本腿那两枚 =="
  /d/work/base/gopath/bin/gofumpt.exe -l internal/winsec/ancestorwindow_342r2_windows_test.go internal/winsec/nametopup_353_windows_test.go
  echo "GOFUMPT_RC=$?"
  echo "== go vet ./internal/winsec/（本机档）=="
  go vet ./internal/winsec/ > "$L/vet-now.txt" 2>&1
  echo "VET_RC=$?"
  wc -c < "$L/vet-now.txt" | sed 's/^/VET_OUTPUT_BYTES=/'
  cat "$L/vet-now.txt"
  echo "== 工作树脏枚数（★本腿只提交自己那枚；别的脏文件＝别腿在飞＝⛔ 卷走）=="
  git status --porcelain -- internal/winsec internal/models cmd/wisp
  echo "== 本腿那枚的 diff -U0 逐行 =="
  git diff -U0 -- internal/winsec/ancestorwindow_342r2_windows_test.go | grep -E '^[+-]' | grep -v '^[+-][+-]'
} 2>&1
