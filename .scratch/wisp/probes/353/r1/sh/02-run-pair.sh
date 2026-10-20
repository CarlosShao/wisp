#!/bin/sh
# 353-r1 颜色尺：成对把同一枚定向 `-run` 打进【基线树 treeB】⇄【落地树 treeL】的仓外物理字节。
# 用法：sh 02-run-pair.sh <host|container> <label> <go-test-regex-or-ALL> [extra .go to sync...]
# ★台面对齐（公共块 台面四味）：
#   host＝go1.27.1 windows/amd64＋GOTMPDIR/TMPDIR＝反斜杠短形、落 %LOCALAPPDATA% 同卷、先 mkdir -p；
#   container＝golang:1.27＋CGO_ENABLED=0＋GOPROXY=off＋MSYS_NO_PATHCONV=1＋★TMPDIR/GOTMPDIR ⛔ 绑 Windows 卷。
# ★仓外跑 Windows .exe 取 DLL＝cygpath -w 的 Windows 形 PATH；★⛔ cygpath -wp 转整条 PATH（gcc 丢⇒CGO 1→0＝整包死）。
# ★四数只在 ^=== RUN／^--- PASS／^--- FAIL／^--- SKIP 取，⛔ 混用⛔ 相减；SHAPE=VOID 当且仅当零枚 === RUN＝台损。
# ★逐名红册另落 .rednames.txt（⛔ 0 字节件）。★串行、⛔ 并行跑改文件的脚本。
set -u
MODE="$1"; LABEL="$2"; RUNPAT="$3"; shift 3
ROOT=/d/tmp/wisp/20261020/353r1
ROOTW='D:\tmp\wisp\20261020\353r1'
LOGS="$ROOT/logs"
mkdir -p "$LOGS"
i=0
for TREE in treeB treeL; do
  i=$((i+1))
  MUT="pair$i-$TREE-$LABEL"
  [ -d "$ROOT/$MUT" ] || cp -a "$ROOT/$TREE" "$ROOT/$MUT"
  # ★把调用方指定的件同步进这棵树的物理字节（两侧同一对字节＝只差本票那枚本体改动）
  # 形状＝incoming/<tree>/<f>（本腿不拿一目录形喂两边）
  for f in "$@"; do
    src="$ROOT/incoming/$TREE/$f"
    if [ -f "$src" ]; then
      cp -f "$src" "$ROOT/$MUT/internal/winsec/$f"
      echo "SYNCED $src -> $ROOT/$MUT/internal/winsec/$f"
    else
      echo "NO_SYNC_FOR_$TREE_$f (keeps the tree's own bytes)"
    fi
  done
  OUT="$LOGS/${LABEL}-s$i-$TREE.raw.txt"
  PKG=./internal/winsec/
  { echo "=== 353r1 ${LABEL} s$i tree=$TREE MODE=$MODE RUNPAT=${RUNPAT} @ $(date '+%F %T %z')"; } > "$OUT"
  if [ "$MODE" = container ]; then
    MSYS_NO_PATHCONV=1 MSYS2_ARG_CONV_EXCL='*' docker run --rm \
      -v "$(cygpath -w "$ROOT/$MUT"):/src:rw" \
      -v 'D:\work\base\gopath\pkg\mod:/gm:ro' \
      -v "$(cygpath -w "$ROOT/gocache-col"):/gobuild:rw" -w /src \
      -e GOPROXY=off -e GOFLAGS=-mod=mod -e CGO_ENABLED=0 -e GOMODCACHE=/gm \
      -e GOTMPDIR=/tmp -e TMPDIR=/tmp \
      golang:1.27 \
      go test -v -count=1 -timeout 40m -run "$RUNPAT" "$PKG" >> "$OUT" 2>&1
    RC=$?
    GT="(container /tmp)"
  else
    TS=$(date +%Y%m%d-%H%M%S)-$$
    GOTMPW="${LOCALAPPDATA:-C:\\Users\\swq\\AppData\\Local}\\wisp-test-gotmp\\353r1-$TS"
    mkdir -p "$(cygpath -u "$GOTMPW")" || { echo "CANNOT_CREATE_GOTMP=$GOTMPW" >> "$OUT"; exit 7; }
    DLLP=$(cygpath -w "$ROOTW\\$MUT\\build")
    SHP=$(cygpath -w "$ROOTW\\$MUT\\third_party\\sherpa-onnx")
    ( cd "$ROOTW\\$MUT" && PATH="$SHP;$DLLP;$PATH" GOTMPDIR="$GOTMPW" TMPDIR="$GOTMPW" \
        GOCACHE="$ROOTW\\gocache-win" go test -v -count=1 -timeout 40m -run "$RUNPAT" "$PKG" ) \
      >> "$OUT" 2>&1
    RC=$?
    GT="$GOTMPW"
  fi
  PINSHA=$(sha256sum "$ROOT/$MUT/internal/winsec/pinnedleaf_342_windows.go" | cut -d' ' -f1)
  RUN=$(grep -c '^=== RUN' "$OUT" || true)
  PT=$(grep -c '^--- PASS' "$OUT" || true)
  FT=$(grep -c '^--- FAIL' "$OUT" || true)
  SK=$(grep -c '^--- SKIP' "$OUT" || true)
  SHAPE=OK; [ "$RUN" = 0 ] && SHAPE=VOID
  LINE="$LABEL s$i $TREE MODE=$MODE RC=$RC RUN=$RUN PASS_TOP=$PT FAIL_TOP=$FT SKIP=$SK SHAPE=$SHAPE GOTMP=$GT"
  echo "$LINE"
  echo "$LINE" >> "$LOGS/four-counts.tsv"
  echo "$LINE PINNEDLEAF_SHA256=$PINSHA"
  grep -E '^--- FAIL' "$OUT" | sed 's/^--- FAIL: //; s/ *([0-9.]*s)$//' | sort > "$LOGS/${LABEL}-s$i-$TREE.rednames.body"
  { echo "# 353r1 ${LABEL} s$i $TREE MODE=$MODE RED_NAMES=$(grep -c . "$LOGS/${LABEL}-s$i-$TREE.rednames.body" || true) 尺=^--- FAIL 计数（同发四数见 four-counts.tsv）";
    cat "$LOGS/${LABEL}-s$i-$TREE.rednames.body"; } > "$LOGS/${LABEL}-s$i-$TREE.rednames.txt"
  rm -f "$LOGS/${LABEL}-s$i-$TREE.rednames.body" 2>/dev/null || true
  echo "  RED_NAMES=$(wc -l < "$LOGS/${LABEL}-s$i-$TREE.rednames.txt")"
  grep -E '353r1 AC#2 over|353r1 AC#2 shot .*LIVENESS|353r1 AC#2 CONTROL' "$OUT" | sed 's/^/  NUM| /'
done
