#!/bin/sh
# 353-r1 AC#3 定名尺＝整包 ./internal/winsec/ 在【treeB2＝产码基线＋同一对测试字节】⇄【treeL＝乙已落】各连打 N 发，
#   然后按**逐名**出矩阵（★⛔ 拿四数差当"我的字节改了 N 枚测试"＝公共块 颜色尺 那条；★⛔ 单发名册当定论）。
# 用法：sh 14-matrix.sh <N>
# ★竞态／非确定性读数⛔ 要求逐字相同＝这里复现的是"哪一枚在两侧都响／哪一枚只在某一侧常响"的形状。
set -u
N="${1:-2}"
ROOT=/d/tmp/wisp/20261020/353r1
ROOTW='D:\tmp\wisp\20261020\353r1'
LOGS="$ROOT/logs"
mkdir -p "$LOGS"
for TREE in treeB2 treeL; do
  for shot in $(seq 1 "$N"); do
    MUT="mx-$TREE-$shot"
    [ -d "$ROOT/$MUT" ] || cp -a "$ROOT/$TREE" "$ROOT/$MUT"
    OUT="$LOGS/mx-$TREE-$shot.raw.txt"
    TS=$(date +%Y%m%d-%H%M%S)-$$
    GOTMPW="${LOCALAPPDATA:-C:\\Users\\swq\\AppData\\Local}\\wisp-test-gotmp\\353r1-mx-$TS"
    mkdir -p "$(cygpath -u "$GOTMPW")" || { echo "CANNOT_CREATE_GOTMP=$GOTMPW"; exit 7; }
    DLLP=$(cygpath -w "$ROOTW\\$MUT\\build"); SHP=$(cygpath -w "$ROOTW\\$MUT\\third_party\\sherpa-onnx")
    ( cd "$ROOTW\\$MUT" && PATH="$SHP;$DLLP;$PATH" GOTMPDIR="$GOTMPW" TMPDIR="$GOTMPW" \
        GOCACHE="$ROOTW\\gocache-win" go test -v -count=1 -timeout 60m ./internal/winsec/ ) > "$OUT" 2>&1
    RC=$?
    LINE="MATRIX $TREE shot=$shot RC=$RC RUN=$(grep -c '^=== RUN' "$OUT" || true) PASS_TOP=$(grep -c '^--- PASS' "$OUT" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$OUT" || true) SKIP=$(grep -c '^--- SKIP' "$OUT" || true)"
    echo "$LINE"; echo "$LINE" >> "$LOGS/four-counts.tsv"
    grep -E '^--- FAIL' "$OUT" | sed 's/^--- FAIL: //; s/ *([0-9.]*s)$//' | sort > "$LOGS/mx-$TREE-$shot.rednames.txt"
  done
done
{
  echo "== 逐名矩阵（★⛔ 四数差＝这枚尺只数名）=="
  cat "$LOGS"/mx-treeB2-*.rednames.txt | sort -u > "$LOGS/mx-allnames.txt"
  printf 'NAME'; for t in treeB2 treeL; do for s in $(seq 1 "$N"); do printf ' |%s-%s' "$t" "$s"; done; done; printf '\n'
  while read -r n; do
    [ -n "$n" ] || continue
    printf '%s' "$n"
    for t in treeB2 treeL; do
      for s in $(seq 1 "$N"); do
        if grep -qxF "$n" "$LOGS/mx-$t-$s.rednames.txt"; then printf ' |RED'; else printf ' |green'; fi
      done
    done
    printf '\n'
  done < "$LOGS/mx-allnames.txt"
  echo "== 从未红过的名（＝两侧全 green）⛔ 出现在此表＝表里只有⛔ 至少响过一发的名；全量名数见下 =="
  echo "DISTINCT_NAMES_EVER_RED=$(wc -l < "$LOGS/mx-allnames.txt")"
  echo "TOTAL_TEST_FUNCS=$(grep -rhcE '^func Test' "$ROOT"/treeL/internal/winsec/*_test.go | paste -sd+ | bc)"
} 2>&1
