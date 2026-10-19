#!/bin/sh
# 349-v1 必答(iii)：Windows 四象限＝（base ⇄ del）×（台损档 GOTMPDIR POSIX 形 ⇄ 干净档 GOTMPDIR Windows 形）
# 每象限一发，四发严格串行（⛔ 并行＝争用洗脏）。四数尺＝RUN/PASS_TOP/FAIL_TOP/SKIP，⛔ 混用⛔ 相减。
# ⛔ 复用 349-r1/343-v1 的树；两棵树都是本腿自己 git archive 导出。
set -u
REPO='/d/work/workspace/projects plans/Wisp'
ROOT=/d/tmp/wisp/20261019/349v1
L="$ROOT/logsQ"
mkdir -p "$L"
GOEXE=/d/work/base/go/bin/go.exe

one () { # $1=tree  $2=shapeLabel(dmg|clean)  $3=which(base|del)
  local T="$ROOT/$1" SHAPE="$2" WHICH="$3"
  local tag="$SHAPE-$WHICH" gt
  if [ "$SHAPE" = dmg ]; then
    mkdir -p "$ROOT/gtDmg/$WHICH"
    gt="$ROOT/gtDmg/$WHICH"                      # POSIX 形＝复现 349-r1／343-v1 pair/ 那一发
  else
    mkdir -p "$ROOT/gtClean/$WHICH"
    gt="$ROOT/gtClean/$WHICH"
    gt=$(cygpath -w "$gt")                        # Windows 形＝复现 343-v1 pair2/ 那一发
  fi
  local PW PB
  PW=$(cygpath -w "$T/third_party/sherpa-onnx"); PB=$(cygpath -w "$T/build")
  echo "### tag=$tag tree=$T GOTMPDIR=$gt"
  ( cd "$T" && PATH="$PW;$PB;$PATH" GOTMPDIR="$gt" TMPDIR= TEMP= TMP= \
      "$GOEXE" test ./internal/risk/ -count=1 -timeout 30m -v ) > "$L/$tag.log" 2>&1
  echo "rc=$?"
  echo "RUN=$(grep -c '^=== RUN' "$L/$tag.log" || true)"
  echo "PASS_TOP=$(grep -c '^--- PASS' "$L/$tag.log" || true)"
  echo "FAIL_TOP=$(grep -c '^--- FAIL' "$L/$tag.log" || true)"
  echo "SKIP=$(grep -c '^--- SKIP' "$L/$tag.log" || true)"
  echo "-- 台损关键词具名 --"
  grep -nE 'mklink|walk|build failed|no test files|setup failed|creating work dir|0xc0000135|Invalid switch' "$L/$tag.log" | head -12 || echo none
  echo "-- 红名册逐名 --"
  grep -E '^--- FAIL' "$L/$tag.log" | sed 's/([0-9.]*s)//;s/^--- FAIL: //;s/ *$//' | sort > "$L/$tag.rednames.txt"
  cat "$L/$tag.rednames.txt"
  echo "-- TestResolvePerCallBudget 行 --"
  grep -E '^(    )*--- (PASS|FAIL|SKIP): TestResolvePerCallBudget' "$L/$tag.log" || echo "(absent)"
  echo "-- bytes --"; wc -c "$L/$tag.log"
  echo ""
}
{
echo "### 349-v1 WINDOWS QUADRANT start $(date -Iseconds)"
echo "### anchors: base=070de7fa(349-r1 起手前最后一笔) del=ed8db1ea(349-r1 第 7 笔交件)"
one treeBase dmg base
one treeDel  dmg del
one treeBase clean base
one treeDel  clean del
echo "### === QUADRANT REDSET DIFFS ==="
for p in "dmg-base dmg-del" "clean-base clean-del" "dmg-base clean-base" "dmg-del clean-del"; do
  set -- $p
  echo "-- symdiff $1 ⇄ $2 --"
  comm -3 "$L/$1.rednames.txt" "$L/$2.rednames.txt" || true
done
echo "### counts matrix"
for t in dmg-base dmg-del clean-base clean-del; do
  printf '%s RUN=%s PASS=%s FAIL=%s SKIP=%s\n' "$t" \
    "$(grep -c '^=== RUN' "$L/$t.log" || true)" \
    "$(grep -c '^--- PASS' "$L/$t.log" || true)" \
    "$(grep -c '^--- FAIL' "$L/$t.log" || true)" \
    "$(grep -c '^--- SKIP' "$L/$t.log" || true)"
done
echo "### QUADRANT_DONE $(date -Iseconds)"
} > "$L/quadrant.log" 2>&1
echo "wrote=$L/quadrant.log"
