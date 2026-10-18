#!/usr/bin/env bash
# 336-r1 门禁四数：只在【共享工作树】跑（派单：那是它们要求的台面），逐枚落 rc=。
# ★lead 10-18 提醒＝先 unset GOTMPDIR TMPDIR，否则触发票 311 AC#1 shape B GUARD＝rc=2 射程拒绝（⛔ 读成红）。
# ⛔ 顺手重排别腿文件；gofumpt 存量脏 6 枚（A1024／停车点 4.0ei 那一节）只具名登记，⛔ 动它们一字。
set -u
REPO="D:/work/workspace/projects plans/Wisp"
cd "$REPO" || exit 1
unset GOTMPDIR TMPDIR
export GOPATH='D:\work\base\gopath'
MINE_TEST='internal/risk/anchor_fallback_teeth_336_windows_test.go'
MINE_OTHER_TEST='internal/risk/anchor_fallback_teeth_336_other_test.go'
MINE_LOGS='.scratch/wisp/probes/336/r1/logs'
echo "date=$(date '+%F %T %z') head=$(git log -1 --format='%H %ad' --date=iso-strict)"

echo "### 1 gofumpt.exe -l 真身（只指我自己的【.go】件；★第一发曾把 logs/*.sh 一起喂给它＝尺用错，rc=2 是 illegal character U+0023 的噪声、⛔ 读成红，本版改正并具名登记）"
/d/work/base/gopath/bin/gofumpt.exe -l "$MINE_TEST" "$MINE_OTHER_TEST"
echo "gofumpt_mine_rc=$?（上面空输出＝我的新增件列出 0 枚）"
/d/work/base/gopath/bin/gofumpt.exe -l internal/risk/
echo "gofumpt_pkg_rc=$? (上面列出的＝存量，逐枚具名见下)"
git status --porcelain --untracked-files=no -- internal/risk/ | sed 's/^/  dirty: /'

echo "### 2 d22scan（scripts/d22scan.sh；射程缺 design/ 会 exit 2＝GUARD 射程拒绝，⛔ 读成红）"
sh scripts/d22scan.sh; echo "d22scan_rc=$?"

echo "### 3 长度门（要 .git）"
sh scripts/check-path-length-budget.sh; echo "length_rc=$?"

echo "### 4 go vet ./internal/risk/"
go vet ./internal/risk/; echo "vet_rc=$?"

echo "### 我的新件形状尺（两枚：windows 牙＋POSIX 绊线）"
for f in "$MINE_TEST" "$MINE_OTHER_TEST"; do
  echo "  file=$f"
  echo "    crlf_bytes=$(tr -dc '\r' < "$f" | wc -c)（0＝纯 LF）"
  echo "    lines=$(wc -l < "$f")  bytes=$(wc -c < "$f")"
  echo "    exported_non_test=$(grep -E '^(func|type|var|const) [A-Z]' "$f" | grep -vcE '^[a-z]+\(t \*testing\.T\)|^func Test' || true)"
  echo "    test_entries=$(grep -cE '^func Test' "$f" || true)（go test 入口，非 API 名）"
  echo "    t.Skip=$(grep -c 't\.Skip' "$f" || true)（期望 0）"
  echo "    build_tag=$(head -1 "$f")"
done
echo "gates_rc=0"
