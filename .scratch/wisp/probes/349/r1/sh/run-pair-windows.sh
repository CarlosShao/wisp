#!/bin/sh
# 349-r1 AC#4 成对四数（Windows 侧）：起手锚仓外树 treeBase ⇄ 交件 treeDel，同一台面、
# `-count=1 -timeout 30m`、四数 RUN/PASS_TOP/FAIL_TOP/SKIP ⛔ 混用⛔ 相减、红名册逐名作差。
# ★两发串行（⛔ 并行＝同一枚 cache/同一包争用会洗脏读数）。
# ★PATH 用 cygpath -w 转 Windows 形＋分号分隔（公共块：POSIX 形传不动 Windows 二进制的 DLL 搜索路径
#   ⇒ 造出"零枚 --- FAIL 的整包死"＝台损、⛔ 计入判据）。GOTMPDIR 先 mkdir -p。
set -u
REPO='/d/work/workspace/projects plans/Wisp'
cd "$REPO" || exit 90
ROOT=/d/tmp/wisp/20261019/349r1
L="$ROOT/logs/pair"
mkdir -p "$L" || exit 91
RUNPAIR="$REPO/.scratch/wisp/probes/349/r1/logs/pair"
mkdir -p "$RUNPAIR" || exit 92
export GOPATH='D:\work\base\gopath'
one () { # $1=tree-dir-name $2=label
  local T="$ROOT/$1" gt="$ROOT/gt/$2"
  mkdir -p "$gt" || { echo "mkdir_gotmp_rc=NONZERO"; return 1; }
  local PW; PW=$(cygpath -w "$T/third_party/sherpa-onnx"); local PB; PB=$(cygpath -w "$T/build")
  ( cd "$T" && PATH="$PW;$PB;$PATH" GOTMPDIR="$gt" TMPDIR= TEMP= TMP= \
      /d/work/base/go/bin/go.exe test ./internal/risk/ -count=1 -timeout 30m -v ) \
      > "$L/$2.log" 2>&1
  local rc=$?
  echo "### label=$2 tree=$T rc=$rc"
  echo "RUN=$(grep -c '^=== RUN' "$L/$2.log" || true)"
  echo "PASS_TOP=$(grep -c '^--- PASS' "$L/$2.log" || true)"
  echo "FAIL_TOP=$(grep -c '^--- FAIL' "$L/$2.log" || true)"
  echo "SKIP=$(grep -c '^--- SKIP' "$L/$2.log" || true)"
  echo "-- 台面故障名册（红句含这些词⇒具名、⛔ 计入判据）--"
  grep -nE 'walk|build failed|no test files|setup failed|creating work dir|0xc0000135' "$L/$2.log" || echo "none"
  echo "-- 红名册逐名 --"
  grep -E '^--- FAIL' "$L/$2.log" || echo "(zero top-level reds)"
  cp -f "$L/$2.log" "$RUNPAIR/win-$2.log"
  echo "-- bytes of log --"; wc -c "$RUNPAIR/win-$2.log"
}
{
echo "### PAIR start (treeBase＝起手锚 7b633fa9) then delivery (treeDel＝交件 db2f9eab)，串行"
one treeBase base
one treeDel  del
echo '### 红名册逐名作差（base ⇄ del；⛔ 拿四数差当"改了 N 枚测试"＝公共块明令）'
grep -E '^--- FAIL' "$L/base.log" | sed 's/([0-9.]*s)//;s/^--- FAIL: //;s/ *$//' | sort > "$L/red-base.txt"
grep -E '^--- FAIL' "$L/del.log"  | sed 's/([0-9.]*s)//;s/^--- FAIL: //;s/ *$//' | sort > "$L/red-del.txt"
wc -l "$L/red-base.txt" "$L/red-del.txt"
echo "-- only in base (fixed by me / disappeared) --"
comm -23 "$L/red-base.txt" "$L/red-del.txt" || true
echo "-- only in del (NEW reds = must be zero) --"
comm -13 "$L/red-base.txt" "$L/red-del.txt" || true
echo "-- 名册逐名（两侧全量）--"
cat "$L/red-base.txt"; echo "  ~~ del ~~"; cat "$L/red-del.txt"
echo "### SKIP 名册（两侧）"
grep -E '^--- SKIP' "$L/base.log" | sed 's/([0-9.]*s)//' || echo "base: none"
grep -E '^--- SKIP' "$L/del.log"  | sed 's/([0-9.]*s)//' || echo "del: none"
} > "$RUNPAIR/win-pair-runner.log" 2>&1
echo "wrote=$RUNPAIR/win-pair-runner.log"
cat "$RUNPAIR/win-pair-runner.log"
