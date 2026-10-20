#!/bin/sh
# 353-r1 编译面自检：三枚 GOOS 档 `go build` ＋ `go vet` ＋ gofumpt 只验自己那几枚。
# ★退码＝重定向文件后 echo RC=$? 紧跟裸命令（⛔ 管进 tail；★⛔ 在 PowerShell 内联里写 $?＝它抢先展开成 True/False）。
# ★⛔ 在共享工作树跑 go test（公共块 台面四味 ③＝判据只认仓外到达数；本件只编译/vet）。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 9
unset GOTMPDIR TMPDIR TEMP TMP
L=/d/tmp/wisp/20261020/353r1/logs
mkdir -p "$L"
{
  echo "== @ $(date '+%F %T %z') HEAD=$(git log -1 --format=%h) =="
  echo "== go build ./internal/winsec/（本机 windows/amd64）=="
  go build ./internal/winsec/ > "$L/build-native.txt" 2>&1
  echo "BUILD_NATIVE_RC=$?"
  cat "$L/build-native.txt"
  echo "== GOOS=linux go build ./internal/winsec/ =="
  GOOS=linux go build ./internal/winsec/ > "$L/build-linux.txt" 2>&1
  echo "BUILD_LINUX_RC=$?"
  cat "$L/build-linux.txt"
  echo "== GOOS=darwin go build ./internal/winsec/ =="
  GOOS=darwin go build ./internal/winsec/ > "$L/build-darwin.txt" 2>&1
  echo "BUILD_DARWIN_RC=$?"
  cat "$L/build-darwin.txt"
  echo "== GOOS=js GOARCH=wasm go build ./internal/winsec/（fallback 档＝预期⛔ 编译＝存量形状，⛔ 判红）=="
  GOOS=js GOARCH=wasm go build ./internal/winsec/ > "$L/build-js.txt" 2>&1
  echo "BUILD_JS_RC=$?"
  head -6 "$L/build-js.txt"
  echo "== go vet ./internal/winsec/（本机档＝判红绿那一枚）=="
  go vet ./internal/winsec/ > "$L/vet-native.txt" 2>&1
  echo "VET_NATIVE_RC=$?"
  cat "$L/vet-native.txt"
  echo "== GOOS=linux go vet（★交叉档⛔ 与本机档成对判红绿＝只登记能编译）=="
  GOOS=linux go vet ./internal/winsec/ > "$L/vet-linux.txt" 2>&1
  echo "VET_LINUX_RC=$?"
  cat "$L/vet-linux.txt"
  echo "== gofumpt -l（★只验自己这 5 枚，⛔ -w、⛔ 扫别腿/别人文件）=="
  /d/work/base/gopath/bin/gofumpt.exe -l internal/winsec/pinnedleaf_342_windows.go internal/winsec/pinnedleaf_342_unix.go internal/winsec/pinnedleaf_342_fallback.go internal/winsec/winsec.go internal/winsec/nametopup_353_windows_test.go > "$L/gofumpt-mine.txt" 2>&1
  echo "GOFUMPT_RC=$?"
  echo "LISTED=$(wc -l < "$L/gofumpt-mine.txt")"
  cat "$L/gofumpt-mine.txt"
} > "$L/compile-selfcheck.txt" 2>&1
cat "$L/compile-selfcheck.txt"
