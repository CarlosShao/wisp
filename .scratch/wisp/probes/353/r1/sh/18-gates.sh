#!/bin/sh
# 353-r1 门禁四把＝★只在共享工作树、跑前 unset GOTMPDIR TMPDIR TEMP TMP（带 env 时 d22scan 里票 311
#   AC#1 shape B 那道 GUARD 会 rc=2＝射程拒绝＝⛔ 读成红，台账 A1027⑤）。
# ★退码＝重定向文件后 echo RC=$? 紧跟裸命令（⛔ 管进 tail＝取到的是 tail 的退码；⛔ PowerShell 内联 $?）。
# ★gofumpt 只列自己那几枚、⛔ -w 扫别腿文件；go vet 交叉档⛔ 与本机档成对判红绿。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 9
OUT=/d/tmp/wisp/20261020/353r1/logs/gates.txt
: > "$OUT"
{
  echo "== HEAD at gate time =="
  git log -1 --format='%H' | sed 's/^/HEAD_FULL_SHA=/'
  git log -1 --format='%s' | sed 's/^/HEAD_SUBJECT=/'
  echo
  echo "== unset GOTMPDIR TMPDIR TEMP TMP =="
  unset GOTMPDIR TMPDIR TEMP TMP
  echo "GOTMPDIR=[${GOTMPDIR-<unset>}] TMPDIR=[${TMPDIR-<unset>}] TEMP=[${TEMP-<unset>}] TMP=[${TMP-<unset>}]"
  echo
  echo "== gate 1: gofumpt -l（★只列本腿那 6 枚，⛔ -w、⛔ 扫别人文件）=="
  /d/work/base/gopath/bin/gofumpt.exe -l internal/winsec/winsec.go internal/winsec/pinnedleaf_342_windows.go \
    internal/winsec/pinnedleaf_342_unix.go internal/winsec/pinnedleaf_342_fallback.go \
    internal/winsec/nametopup_353_windows_test.go internal/winsec/ancestorwindow_342r2_windows_test.go \
    > /d/tmp/wisp/20261020/353r1/logs/gate-gofumpt.txt 2>&1
  echo "GOFUMPT_RC=$?"
  echo "LISTED=$(wc -l < /d/tmp/wisp/20261020/353r1/logs/gate-gofumpt.txt)"
  cat /d/tmp/wisp/20261020/353r1/logs/gate-gofumpt.txt
  echo
  echo "== gate 2: sh scripts/d22scan.sh =="
  sh scripts/d22scan.sh > /d/tmp/wisp/20261020/353r1/logs/gate-d22scan.txt 2>&1
  echo "D22_RC=$?"
  tail -3 /d/tmp/wisp/20261020/353r1/logs/gate-d22scan.txt
  echo
  echo "== gate 3: sh scripts/check-path-length-budget.sh（要 .git）=="
  sh scripts/check-path-length-budget.sh > /d/tmp/wisp/20261020/353r1/logs/gate-length.txt 2>&1
  echo "LEN_RC=$?"
  tail -4 /d/tmp/wisp/20261020/353r1/logs/gate-length.txt
  echo
  echo "== gate 4: go vet ./internal/winsec/（本机档＝判红绿那一枚）=="
  go vet ./internal/winsec/ > /d/tmp/wisp/20261020/353r1/logs/gate-vet-native.txt 2>&1
  echo "VET_NATIVE_RC=$?"
  wc -c < /d/tmp/wisp/20261020/353r1/logs/gate-vet-native.txt | sed 's/^/VET_NATIVE_OUTPUT_BYTES=/'
  cat /d/tmp/wisp/20261020/353r1/logs/gate-vet-native.txt
  echo
  echo "== gate 4b: 交叉档＝只登记这档能编译＋vet 无声，★⛔ 与本机档成对判红绿 =="
  GOOS=linux go vet ./internal/winsec/ > /d/tmp/wisp/20261020/353r1/logs/gate-vet-linux.txt 2>&1
  echo "VET_LINUX_RC=$?"
  wc -c < /d/tmp/wisp/20261020/353r1/logs/gate-vet-linux.txt | sed 's/^/VET_LINUX_OUTPUT_BYTES=/'
  GOOS=darwin go vet ./internal/winsec/ > /d/tmp/wisp/20261020/353r1/logs/gate-vet-darwin.txt 2>&1
  echo "VET_DARWIN_RC=$?"
  wc -c < /d/tmp/wisp/20261020/353r1/logs/gate-vet-darwin.txt | sed 's/^/VET_DARWIN_OUTPUT_BYTES=/'
  echo
  echo "== 本腿写面 blob（入库后）=="
  git rev-parse HEAD:internal/winsec/nametopup_353_windows_test.go | sed 's/^/BLOB_353_TEST=/'
  git rev-parse HEAD:internal/winsec/pinnedleaf_342_windows.go | sed 's/^/BLOB_PINNEDLEAF=/'
  git rev-parse HEAD:internal/winsec/winsec.go | sed 's/^/BLOB_WINSEC_GO=/'
} >> "$OUT" 2>&1
cat "$OUT"
