#!/bin/sh
# 353-r1 AC#3 逐名作差尺：★同一对测试字节、只差产码字节＝【treeB2（产码＝基线）】⇄【treeL（产码＝乙已落）】。
# 用法：sh 11-pername-pair.sh <label> <treeA> <treeB>
# ★⛔ ./cmd/wisp（321-r1 写面＋本波明令）；★⛔ ./... ；★⛔ 在共享工作树跑。串行、⛔ 并行改文件。
# 四数只在 ^=== RUN／^--- PASS／^--- FAIL／^--- SKIP 取；逐名册另落 .rednames.txt，作差另算。
set -u
LABEL="$1"; TA="$2"; TB="$3"
ROOT=/d/tmp/wisp/20261020/353r1
ROOTW='D:\tmp\wisp\20261020\353r1'
LOGS="$ROOT/logs"
mkdir -p "$LOGS"
i=0
for TREE in "$TA" "$TB"; do
  i=$((i+1))
  MUT="pn$i-$TREE-$LABEL"
  [ -d "$ROOT/$MUT" ] || cp -a "$ROOT/$TREE" "$ROOT/$MUT"
  OUT="$LOGS/${LABEL}-s$i-$TREE.raw.txt"
  { echo "=== 353r1 ${LABEL} s$i tree=$TREE pername-pair @ $(date '+%F %T %z')"; } > "$OUT"
  TS=$(date +%Y%m%d-%H%M%S)-$$
  GOTMPW="${LOCALAPPDATA:-C:\\Users\\swq\\AppData\\Local}\\wisp-test-gotmp\\353r1-$TS"
  mkdir -p "$(cygpath -u "$GOTMPW")" || { echo "CANNOT_CREATE_GOTMP=$GOTMPW" >> "$OUT"; exit 7; }
  DLLP=$(cygpath -w "$ROOTW\\$MUT\\build"); SHP=$(cygpath -w "$ROOTW\\$MUT\\third_party\\sherpa-onnx")
  ( cd "$ROOTW\\$MUT" && PATH="$SHP;$DLLP;$PATH" GOTMPDIR="$GOTMPW" TMPDIR="$GOTMPW" \
      GOCACHE="$ROOTW\\gocache-win" go test -v -count=1 -timeout 60m ./internal/winsec/ ) >> "$OUT" 2>&1
  RC=$?
  RUN=$(grep -c '^=== RUN' "$OUT" || true)
  PT=$(grep -c '^--- PASS' "$OUT" || true)
  FT=$(grep -c '^--- FAIL' "$OUT" || true)
  SK=$(grep -c '^--- SKIP' "$OUT" || true)
  SHAPE=OK; [ "$RUN" = 0 ] && SHAPE=VOID
  LINE="$LABEL s$i $TREE RC=$RC RUN=$RUN PASS_TOP=$PT FAIL_TOP=$FT SKIP=$SK SHAPE=$SHAPE GOTMP=$GOTMPW"
  echo "$LINE"; echo "$LINE" >> "$LOGS/four-counts.tsv"
  grep -E '^--- FAIL' "$OUT" | sed 's/^--- FAIL: //; s/ *([0-9.]*s)$//' | sort > "$LOGS/${LABEL}-s$i-$TREE.rednames.txt.body"
  { echo "# 353r1 ${LABEL} s$i $TREE RED_NAMES=$(grep -c . "$LOGS/${LABEL}-s$i-$TREE.rednames.txt.body" || true) 尺=^--- FAIL"
    cat "$LOGS/${LABEL}-s$i-$TREE.rednames.txt.body"; } > "$LOGS/${LABEL}-s$i-$TREE.rednames.txt"
  rm -f "$LOGS/${LABEL}-s$i-$TREE.rednames.txt.body"
  echo "  RED_NAMES=$(wc -l < "$LOGS/${LABEL}-s$i-$TREE.rednames.txt")"
done
echo "== 逐名作差（★⛔ 拿四数差当"我的字节改了 N 枚测试"）=="
comm -23 "$LOGS/${LABEL}-s1-$TA.rednames.txt" "$LOGS/${LABEL}-s2-$TB.rednames.txt" | grep -v '^#' | sed 's/^/  ONLY_IN_BASELINE(消掉的红): /' || true
comm -13 "$LOGS/${LABEL}-s1-$TA.rednames.txt" "$LOGS/${LABEL}-s2-$TB.rednames.txt" | grep -v '^#' | sed 's/^/  ONLY_IN_LANDED(新增红): /' || true
echo "  SYMDIFF_COUNT=$(comm -3 "$LOGS/${LABEL}-s1-$TA.rednames.txt" "$LOGS/${LABEL}-s2-$TB.rednames.txt" | grep -vc '^#' || true)"
