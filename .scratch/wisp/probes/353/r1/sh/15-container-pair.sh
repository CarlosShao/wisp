#!/bin/sh
# 353-r1 AC#3 另一平台＝容器 leg（golang:1.27 linux/amd64，CGO_ENABLED=0、GOPROXY=off）。
# 用法：sh 15-container-pair.sh <label> <treeA> <treeB>
# ★TMPDIR/GOTMPDIR ⛔ 绑 Windows 卷（容器自 /tmp；公共块：os.SameFile 失真）；★MSYS_NO_PATHCONV=1（不转＝RC=125 台损）。
# ★⛔ ./cmd/wisp；★⛔ ./... ；★逐名册另落 .rednames.txt；SHAPE=VOID 当零枚 ^=== RUN＝台损、⛔ 计入判据。
set -u
LABEL="$1"; TA="$2"; TB="$3"
ROOT=/d/tmp/wisp/20261020/353r1
LOGS="$ROOT/logs"
mkdir -p "$LOGS"
i=0
for TREE in "$TA" "$TB"; do
  i=$((i+1))
  OUT="$LOGS/${LABEL}-s$i-$TREE.raw.txt"
  { echo "=== 353r1 ${LABEL} s$i tree=$TREE MODE=container @ $(date '+%F %T %z')"; } > "$OUT"
  MSYS_NO_PATHCONV=1 MSYS2_ARG_CONV_EXCL='*' docker run --rm \
    -v "$(cygpath -w "$ROOT/$TREE"):/src:rw" \
    -v 'D:\work\base\gopath\pkg\mod:/gm:ro' \
    -v "$(cygpath -w "$ROOT/gocache-col-$TREE"):/gobuild:rw" -w /src \
    -e GOPROXY=off -e GOFLAGS=-mod=mod -e CGO_ENABLED=0 -e GOMODCACHE=/gm \
    -e GOTMPDIR=/tmp -e TMPDIR=/tmp \
    golang:1.27 \
    go test -v -count=1 -timeout 60m ./internal/winsec/ >> "$OUT" 2>&1
  RC=$?
  RUN=$(grep -c '^=== RUN' "$OUT" || true)
  PT=$(grep -c '^--- PASS' "$OUT" || true)
  FT=$(grep -c '^--- FAIL' "$OUT" || true)
  SK=$(grep -c '^--- SKIP' "$OUT" || true)
  SHAPE=OK; [ "$RUN" = 0 ] && SHAPE=VOID
  LINE="$LABEL s$i $TREE MODE=container RC=$RC RUN=$RUN PASS_TOP=$PT FAIL_TOP=$FT SKIP=$SK SHAPE=$SHAPE"
  echo "$LINE"; echo "$LINE" >> "$LOGS/four-counts.tsv"
  grep -E '^--- FAIL' "$OUT" | sed 's/^--- FAIL: //; s/ *([0-9.]*s)$//' | sort > "$LOGS/${LABEL}-s$i-$TREE.rednames.txt.body"
  { echo "# 353r1 ${LABEL} s$i $TREE container RED_NAMES=$(grep -c . "$LOGS/${LABEL}-s$i-$TREE.rednames.txt.body" || true) 尺=^--- FAIL"
    cat "$LOGS/${LABEL}-s$i-$TREE.rednames.txt.body"; } > "$LOGS/${LABEL}-s$i-$TREE.rednames.txt"
  rm -f "$LOGS/${LABEL}-s$i-$TREE.rednames.txt.body"
  echo "  RED_NAMES=$(wc -l < "$LOGS/${LABEL}-s$i-$TREE.rednames.txt")"
  tail -4 "$OUT" | sed 's/^/   TAIL| /'
done
echo "== 逐名作差 =="
comm -23 "$LOGS/${LABEL}-s1-$TA.rednames.txt" "$LOGS/${LABEL}-s2-$TB.rednames.txt" | grep -v '^#' | sed 's/^/  ONLY_IN_BASELINE: /' || true
comm -13 "$LOGS/${LABEL}-s1-$TA.rednames.txt" "$LOGS/${LABEL}-s2-$TB.rednames.txt" | grep -v '^#' | sed 's/^/  ONLY_IN_LANDED(新增红): /' || true
echo "  SYMDIFF=$(comm -3 "$LOGS/${LABEL}-s1-$TA.rednames.txt" "$LOGS/${LABEL}-s2-$TB.rednames.txt" | grep -vc '^#' || true)"
