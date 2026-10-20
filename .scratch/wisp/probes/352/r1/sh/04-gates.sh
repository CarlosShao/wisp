#!/bin/sh
# 352-r1 门禁四数＝★只在共享工作树、跑前 unset GOTMPDIR TMPDIR TEMP TMP（带 env 时 d22scan 里票 311
#   AC#1 shape B 那道 GUARD 会 rc=2＝射程拒绝＝⛔ 读成红，台账 A1027⑤）。
# ★退码＝重定向文件 + echo RC=$? 紧跟裸命令（⛔ 管进 tail＝那取到的是 tail 的退码）。
# ★gofumpt 只验自己那枚、⛔ -w、⛔ 扫别腿文件；go vet ⛔ 与交叉档成对判红绿（交叉档只登记"能编译"）。
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 9
OUT=/d/tmp/wisp/20261020/352r1/logs/gates.txt
: > "$OUT"
{
  echo "== HEAD at gate time（⛔ cut -c：按字节截中文＝非法 UTF-8）=="
  git log -1 --format='%H' | sed 's/^/HEAD_FULL_SHA=/'
  git log -1 --format='%s' | sed 's/^/HEAD_SUBJECT=/'
  echo
  echo "== unset GOTMPDIR TMPDIR TEMP TMP =="
  unset GOTMPDIR TMPDIR TEMP TMP
  echo "GOTMPDIR=[${GOTMPDIR-<unset>}] TMPDIR=[${TMPDIR-<unset>}] TEMP=[${TEMP-<unset>}] TMP=[${TMP-<unset>}]"
  echo
  echo "== gate 1: gofumpt -l internal/winsec/winsec.go（只验自己这枚；⛔ -w）=="
  /d/work/base/gopath/bin/gofumpt.exe -l internal/winsec/winsec.go
  echo "GOFUMPT_RC=$?"
  echo "(空列出＝干净；本机工作树有既有件因 CRLF 被标脏＝存量，本腿⛔ 顺手重排别腿文件)"
  echo
  echo "== gate 1b: 既有 CRLF 脏枚数登记（同一条尺、只数不改）=="
  /d/work/base/gopath/bin/gofumpt.exe -l ./... > /d/tmp/wisp/20261020/352r1/logs/gofumpt-all.txt 2>&1
  echo "GOFUMPT_ALL_RC=$?  LISTED=$(wc -l < /d/tmp/wisp/20261020/352r1/logs/gofumpt-all.txt)"
  grep -c 'winsec/winsec.go' /d/tmp/wisp/20261020/352r1/logs/gofumpt-all.txt | sed 's/^/MY_FILE_LISTED=/' || true
  echo
  echo "== gate 2: sh scripts/d22scan.sh =="
  sh scripts/d22scan.sh > /d/tmp/wisp/20261020/352r1/logs/gate-d22scan.txt 2>&1
  echo "D22_RC=$?"
  tail -2 /d/tmp/wisp/20261020/352r1/logs/gate-d22scan.txt
  echo
  echo "== gate 3: sh scripts/check-path-length-budget.sh（要 .git）=="
  sh scripts/check-path-length-budget.sh > /d/tmp/wisp/20261020/352r1/logs/gate-length.txt 2>&1
  echo "LEN_RC=$?"
  tail -3 /d/tmp/wisp/20261020/352r1/logs/gate-length.txt
  echo
  echo "== gate 4: go vet ./internal/winsec/（本机档＝判红绿的那一枚）=="
  go vet ./internal/winsec/ > /d/tmp/wisp/20261020/352r1/logs/gate-vet-native.txt 2>&1
  echo "VET_NATIVE_RC=$?"
  wc -c < /d/tmp/wisp/20261020/352r1/logs/gate-vet-native.txt | sed 's/^/VET_NATIVE_OUTPUT_BYTES=/'
  echo
  echo "== gate 4b: 交叉档＝只登记"这档能编译＋vet 无声"，★⛔ 与本机档成对判红绿 =="
  GOOS=linux go vet ./internal/winsec/ > /d/tmp/wisp/20261020/352r1/logs/gate-vet-linux.txt 2>&1
  echo "VET_LINUX_RC=$?"
  wc -c < /d/tmp/wisp/20261020/352r1/logs/gate-vet-linux.txt | sed 's/^/VET_LINUX_OUTPUT_BYTES=/'
  GOOS=darwin go vet ./internal/winsec/ > /d/tmp/wisp/20261020/352r1/logs/gate-vet-darwin.txt 2>&1
  echo "VET_DARWIN_RC=$?"
  wc -c < /d/tmp/wisp/20261020/352r1/logs/gate-vet-darwin.txt | sed 's/^/VET_DARWIN_OUTPUT_BYTES=/'
  echo
  echo "== winsec.go blob 各笔串（供编排者判 351-r2 读数可比性）=="
  git rev-parse HEAD:internal/winsec/winsec.go | sed 's/^/HEAD_BLOB=/'
  git hash-object internal/winsec/winsec.go | sed 's/^/WORKTREE_BLOB=/'
  for c in 0fb399c3 ecffb93b aacde632 c16ff47e; do
    git rev-parse "$c:internal/winsec/winsec.go" 2>/dev/null | sed "s|^|BLOB_AT_$c=|"
  done
} >> "$OUT" 2>&1
echo "WROTE $OUT"
wc -l "$OUT"
