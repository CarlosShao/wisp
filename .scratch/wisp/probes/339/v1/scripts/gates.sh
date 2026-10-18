#!/bin/bash
# 339-v1 台件⑪：★门禁四数（只在共享工作树、跑前 unset GOTMPDIR TMPDIR TEMP TMP＝A1027 第⑤节定式）。
# 逐枚落 rc=；红句含 walk/build failed/no test files/setup failed/0xc0000135 ⇒ 判台面故障具名⛔ 计入判据。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
cd "$REPO" || exit 9
unset GOTMPDIR TMPDIR TEMP TMP
LOG=".scratch/wisp/probes/339/v1/logs/v1-gates-shared-tree.txt"
{
  echo "unset 后 GOTMPDIR=[${GOTMPDIR:-UNSET}] TMPDIR=[${TMPDIR:-UNSET}]（★跑前清 env＝票 311 AC#1 shape B GUARD 定式）"
  echo "--- 门① gofumpt -l（真身 /d/work/base/gopath/bin/gofumpt.exe，射程＝本腿写面的 .go 件；本腿零 .go ⇒ 尺打目录＝空输出）---"
  echo "（★具名一枚坏尺：第一发把这把尺打到了 .md 判语件上 ⇒ rc=2 illegal character U+0023＝尺用错对象、⛔ 读成红、⛔ 计入判据；gofumpt 只吃 Go 源）"
  GOFILES=$(find .scratch/wisp/probes/339/v1 docs/evidence/s1 -maxdepth 2 -name '*.go' -type f | wc -l)
  echo "my_go_files=$GOFILES (本腿写面 .go 枚数＝0＝正控：尺对空集的正确输出＝空)"
  OUT1=$(/d/work/base/gopath/bin/gofumpt.exe -l .scratch/wisp/probes/339/v1 docs/evidence/s1 2>&1); RC1=$?
  echo "gofumpt_rc=$RC1 output_lines=$(printf '%s' "$OUT1" | grep -c . || true)"
  printf '%s\n' "$OUT1"
  echo "--- 门①b 全仓 gofumpt -l（存量对照，逐枚点名⛔ 我写面的件；⛔ 顺手重排别人件）---"
  /d/work/base/gopath/bin/gofumpt.exe -l . 2>/dev/null | grep -E '339/v1|evidence/s1' || echo "repo_dirty_listed_not_mine=$( /d/work/base/gopath/bin/gofumpt.exe -l . 2>/dev/null | grep -c . || true)"
  echo "--- 门② d22scan ---"
  OUT2=$(sh scripts/d22scan.sh 2>&1); RC2=$?
  echo "d22_rc=$RC2"; printf '%s\n' "$OUT2" | tail -3
  echo "--- 门③ 长度门（要 .git）---"
  OUT3=$(sh scripts/check-path-length-budget.sh 2>&1); RC3=$?
  echo "len_rc=$RC3"; printf '%s\n' "$OUT3" | grep -E 'VERDICT|denominator|longest'
  echo "--- 门④ go vet ./internal/risk/ ---"
  OUT4=$(go vet ./internal/risk/ 2>&1); RC4=$?
  echo "vet_rc=$RC4"; printf '%s\n' "$OUT4"
  echo "--- 门④b 两平台 vet 交叉（容器侧，交件树）---"
  echo "（另见 logs/v1-*-fullpkg：包内编译已证；容器 vet 读数＝实现腿 linux-final-vet.log、本腿⛔ 抄：见 v1-jia-naive 那份 vet 尺同源）"
  echo "RC_TABLE rc1=$RC1 rc2=$RC2 rc3=$RC3 rc4=$RC4"
} > "$LOG" 2>&1
cat "$LOG"
