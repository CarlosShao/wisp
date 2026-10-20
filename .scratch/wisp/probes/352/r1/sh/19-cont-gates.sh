#!/bin/bash
# 352-r1 续跑（cont）· 门禁四数复量（只在共享工作树；跑前 unset GOTMPDIR TMPDIR TEMP TMP）
# 退码取法＝脚本内重定向到文件后 echo RC=$? 紧跟裸命令（公共块读数纪律）。
# ★续跑⛔ 动任何产码字节⇒四道门应与笔⑦入库版同绿；这里只现量登记，⛔ 据别腿红判自己失败。
cd '/d/work/workspace/projects plans/Wisp' || exit 9
unset GOTMPDIR TMPDIR TEMP TMP
echo "PWD=$(pwd)"
echo "ENV GOTMPDIR=[${GOTMPDIR-UNSET}] TMPDIR=[${TMPDIR-UNSET}] TEMP=[${TEMP-UNSET}] TMP=[${TMP-UNSET}]"
GOFUMPT=/d/work/base/gopath/bin/gofumpt.exe
echo '===== 门1 gofumpt -l internal/winsec/winsec.go ====='
"$GOFUMPT" -l internal/winsec/winsec.go > /d/tmp/wisp/20261020/352r1cont/tmp/g1.txt 2>&1
echo RC=$?
echo "  LISTED_LINES=$(grep -c '' /d/tmp/wisp/20261020/352r1cont/tmp/g1.txt || true)"
echo '===== 门2 sh scripts/d22scan.sh ====='
sh scripts/d22scan.sh > /d/tmp/wisp/20261020/352r1cont/tmp/g2.txt 2>&1
echo RC=$?
grep -oE 'clean - no D22 ban violations|internal/=603|cmd/=143' /d/tmp/wisp/20261020/352r1cont/tmp/g2.txt | head -3 | sed 's/^/  /'
echo '===== 门3 sh scripts/check-path-length-budget.sh ====='
sh scripts/check-path-length-budget.sh > /d/tmp/wisp/20261020/352r1cont/tmp/g3.txt 2>&1
echo RC=$?
grep -oE 'VERDICT GREEN|over-budget=60|roster entries=60|not in roster=0' /d/tmp/wisp/20261020/352r1cont/tmp/g3.txt | sed 's/^/  /'
echo '===== 门4 go vet ./internal/winsec/ ====='
go vet ./internal/winsec/ > /d/tmp/wisp/20261020/352r1cont/tmp/g4.txt 2>&1
echo RC=$?
echo "  VET_OUT_BYTES=$(wc -c < /d/tmp/wisp/20261020/352r1cont/tmp/g4.txt)"
echo '===== 交叉档（只登记"能编译＋vet 无声"，⛔ 与本机档成对判红绿）====='
GOOS=linux go vet ./internal/winsec/ > /d/tmp/wisp/20261020/352r1cont/tmp/g5.txt 2>&1
echo RC_linux=$?
GOOS=darwin go vet ./internal/winsec/ > /d/tmp/wisp/20261020/352r1cont/tmp/g6.txt 2>&1
echo RC_darwin=$?
echo '===== 存量具名（⛔ 我动、⛔ 顺手重排）gofumpt -l internal cmd tools ====='
"$GOFUMPT" -l internal cmd tools > /d/tmp/wisp/20261020/352r1cont/tmp/g7.txt 2>&1
echo RC=$?
echo "  STOCK_DIRTY=$(grep -c '' /d/tmp/wisp/20261020/352r1cont/tmp/g7.txt || true)"
grep -v '^$' /d/tmp/wisp/20261020/352r1cont/tmp/g7.txt | sed 's/^/    /' | cat
echo '===== 那 4 枚存量件里有没有我这趟新造的（应全住别腿/别处，⛔ 我产码）====='
echo DONE
