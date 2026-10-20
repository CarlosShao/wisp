#!/bin/sh
# 353-r1 AC#3 尺：整包 ./internal/winsec/ 成对打进【基线树 treeB（乙未落＋本腿仪器件）】⇄【落地树 treeL（乙已落）】。
# 用法：sh 10-full-pair.sh <host|container> <label>
# ★⛔ ./cmd/wisp（321-r1 写面＋本波明令）；★⛔ ./... ；★⛔ 在共享工作树跑（半保存态假红＝公共块）。
# 四数只在 ^=== RUN／^--- PASS／^--- FAIL／^--- SKIP 取、⛔ 混用⛔ 相减；逐名红册另落 .rednames.txt。
# SHAPE=VOID 当且仅当零枚 ^=== RUN＝台面故障、⛔ 计入判据（红句含 walk／build failed／0xc0000135 同判）。
set -u
MODE="$1"; LABEL="$2"
ROOT=/d/tmp/wisp/20261020/353r1
ROOTW='D:\tmp\wisp\20261020\353r1'
LOGS="$ROOT/logs"
mkdir -p "$LOGS"
i=0
for TREE in treeB treeL; do
  i=$((i+1))
  MUT="full$i-$TREE-$LABEL"
  [ -d "$ROOT/$MUT" ] || cp -a "$ROOT/$TREE" "$ROOT/$MUT"
  OUT="$LOGS/${LABEL}-s$i-$TREE.raw.txt"
  { echo "=== 353r1 ${LABEL} s$i tree=$TREE MODE=$MODE FULL-PACKAGE ./internal/winsec/ @ $(date '+%F %T %z')"; } > "$OUT"
  if [ "$MODE" = container ]; then
    MSYS_NO_PATHCONV=1 MSYS2_ARG_CONV_EXCL='*' docker run --rm \
      -v "$(cygpath -w "$ROOT/$MUT"):/src:rw" \
      -v 'D:\work\base\gopath\pkg\mod:/gm:ro' \
      -v "$(cygpath -w "$ROOT/gocache-col-$TREE"):/gobuild:rw" -w /src \
      -e GOPROXY=off -e GOFLAGS=-mod=mod -e CGO_ENABLED=0 -e GOMODCACHE=/gm \
      -e GOTMPDIR=/tmp -e TMPDIR=/tmp \
      golang:1.27 \
      go test -v -count=1 -timeout 60m ./internal/winsec/ >> "$OUT" 2>&1
    RC=$?
    GT="(container /tmp)"
  else
    TS=$(date +%Y%m%d-%H%M%S)-$$
    GOTMPW="${LOCALAPPDATA:-C:\\Users\\swq\\AppData\\Local}\\wisp-test-gotmp\\353r1-$TS"
    mkdir -p "$(cygpath -u "$GOTMPW")" || { echo "CANNOT_CREATE_GOTMP=$GOTMPW" >> "$OUT"; exit 7; }
    DLLP=$(cygpath -w "$ROOTW\\$MUT\\build"); SHP=$(cygpath -w "$ROOTW\\$MUT\\third_party\\sherpa-onnx")
    ( cd "$ROOTW\\$MUT" && PATH="$SHP;$DLLP;$PATH" GOTMPDIR="$GOTMPW" TMPDIR="$GOTMPW" \
        GOCACHE="$ROOTW\\gocache-win" go test -v -count=1 -timeout 60m ./internal/winsec/ ) >> "$OUT" 2>&1
    RC=$?
    GT="$GOTMPW"
  fi
  RUN=$(grep -c '^=== RUN' "$OUT" || true)
  PT=$(grep -c '^--- PASS' "$OUT" || true)
  FT=$(grep -c '^--- FAIL' "$OUT" || true)
  SK=$(grep -c '^--- SKIP' "$OUT" || true)
  SHAPE=OK; [ "$RUN" = 0 ] && SHAPE=VOID
  LINE="$LABEL s$i $TREE MODE=$MODE RC=$RC RUN=$RUN PASS_TOP=$PT FAIL_TOP=$FT SKIP=$SK SHAPE=$SHAPE GOTMP=$GT"
  echo "$LINE"
  echo "$LINE" >> "$LOGS/four-counts.tsv"
  grep -E '^--- FAIL' "$OUT" | sed 's/^--- FAIL: //; s/ *([0-9.]*s)$//' | sort > "$LOGS/${LABEL}-s$i-$TREE.rednames.body"
  { echo "# 353r1 ${LABEL} s$i $TREE MODE=$MODE RED_NAMES=$(grep -c . "$LOGS/${LABEL}-s$i-$TREE.rednames.body" || true) 尺=^--- FAIL（同发四数见 four-counts.tsv）"
    cat "$LOGS/${LABEL}-s$i-$TREE.rednames.body"; } > "$LOGS/${LABEL}-s$i-$TREE.rednames.txt"
  rm -f "$LOGS/${LABEL}-s$i-$TREE.rednames.body" 2>/dev/null || true
  echo "  RED_NAMES=$(wc -l < "$LOGS/${LABEL}-s$i-$TREE.rednames.txt")"
  grep -E '^\s+--- (FAIL|SKIP)' "$OUT" | sed 's/^ *//' | sort -u > "$LOGS/${LABEL}-s$i-$TREE.subnames.txt"
  echo "  SUB_FAILSKIP=$(wc -l < "$LOGS/${LABEL}-s$i-$TREE.subnames.txt")"
  echo "  TAIL:"; tail -3 "$OUT" | sed 's/^/   /'
done
