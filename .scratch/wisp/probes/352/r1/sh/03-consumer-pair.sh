#!/bin/sh
# 352-r1 AC#2 尺：五包消费者面 ⇄ internal/winsec 全包，成对打进【基线树拷贝】⇄【落地树拷贝】，逐名作差。
# 用法：sh 03-consumer-pair.sh <host|container> <label> <pkg...>
# ★台面对齐：host＝go1.27.1 windows/amd64＋GOTMPDIR 反斜杠形落 %LOCALAPPDATA% 同卷短根；
#   container＝golang:1.27＋CGO_ENABLED=0＋GOPROXY=off＋MSYS_NO_PATHCONV=1（首发实测：不转＝RC=125 台损）
#   ＋★TMPDIR/GOTMPDIR ⛔ 绑 Windows 卷（容器自 /tmp）。
# ★⛔ cmd/wisp（本波明令＝归编排者那窗，且会撞 351-r2 写面）；★⛔ ./internal/agent/**（＝230-r1 写面那棵子包）
#   ⇒ 本腿只打 ./internal/agent/ 这一枚包自身，且打在仓外导出树上（⛔ 共享工作树＝⛔ 卷走别腿在飞字节）。
# ★四数只在 ^=== RUN／^--- PASS／^--- FAIL／^--- SKIP 取、⛔ 混用⛔ 相减；逐名册另落 .rednames.txt。
# ★SHAPE=VOID 当且仅当零枚 === RUN＝台损、⛔ 计入判据。串行。
set -u
MODE="$1"; LABEL="$2"; shift 2
ROOT=/d/tmp/wisp/20261020/352r1
ROOTW='D:\tmp\wisp\20261020\352r1'
LOGS="$ROOT/logs"
mkdir -p "$LOGS"
i=0
for TREE in treeB treeL; do
  i=$((i+1))
  MUT="mut-c$i-$TREE"
  [ -d "$ROOT/$MUT" ] || cp -a "$ROOT/$TREE" "$ROOT/$MUT"
  OUT="$LOGS/${LABEL}-s$i-$TREE.raw.txt"
  PRE=$(sha256sum "$ROOT/$MUT/internal/winsec/winsec.go" | cut -d' ' -f1)
  if [ "$MODE" = container ]; then
    MSYS_NO_PATHCONV=1 MSYS2_ARG_CONV_EXCL='*' docker run --rm \
      -v "$(cygpath -w "$ROOT/$MUT"):/src:rw" \
      -v 'D:\work\base\gopath\pkg\mod:/gm:ro' \
      -v "$(cygpath -w "$ROOT/gocache-col"):/gobuild:rw" -w /src \
      -e GOPROXY=off -e GOFLAGS=-mod=mod -e CGO_ENABLED=0 -e GOMODCACHE=/gm \
      -e GOTMPDIR=/tmp -e TMPDIR=/tmp \
      golang:1.27 \
      go test -v -count=1 -timeout 40m "$@" > "$OUT" 2>&1
    RC=$?
    GT="(container /tmp)"
  else
    TS=$(date +%Y%m%d-%H%M%S)-$$
    GOTMPW="${LOCALAPPDATA:-C:\\Users\\swq\\AppData\\Local}\\wisp-test-gotmp\\352r1-$TS"
    mkdir -p "$(cygpath -u "$GOTMPW")" || { echo "CANNOT_CREATE_GOTMP=$GOTMPW"; exit 7; }
    # ★仓外导出树跑 Windows .exe 取 DLL＝PATH 用 cygpath -w 的 Windows 形＋分号分隔（342-v2 sh/12 同式）；
    #   ★⛔ cygpath -wp 转整条 PATH（那会把无盘符 MSYS 段 /mingw64/bin 指错⇒CGO 1→0＝整包死的假红）。
    DLLP=$(cygpath -w "$ROOTW\\$MUT\\build"); SHP=$(cygpath -w "$ROOTW\\$MUT\\third_party\\sherpa-onnx")
    ( cd "$ROOTW\\$MUT" && PATH="$SHP;$DLLP;$PATH" GOTMPDIR="$GOTMPW" TMPDIR="$GOTMPW" \
        GOCACHE='D:\tmp\wisp\20261020\352r1\gocache-win' \
        go test -v -count=1 -timeout 40m "$@" ) > "$OUT" 2>&1
    RC=$?
    GT="$GOTMPW"
  fi
  POST=$(sha256sum "$ROOT/$MUT/internal/winsec/winsec.go" | cut -d' ' -f1)
  RESTORE=OK; [ "$PRE" = "$POST" ] || RESTORE="DRIFT $PRE -> $POST"
  RUN=$(grep -c '^=== RUN' "$OUT" || true)
  PT=$(grep -c '^--- PASS' "$OUT" || true)
  FT=$(grep -c '^--- FAIL' "$OUT" || true)
  SK=$(grep -c '^--- SKIP' "$OUT" || true)
  SHAPE=OK; [ "$RUN" = 0 ] && SHAPE=VOID
  LINE="$LABEL s$i $TREE MODE=$MODE RC=$RC RUN=$RUN PASS_TOP=$PT FAIL_TOP=$FT SKIP=$SK SHAPE=$SHAPE WINSEC_SHA256=$POST RESTORE=$RESTORE GOTMP=$GT"
  echo "$LINE"
  echo "$LINE" >> "$LOGS/four-counts-consumers.tsv"
  grep -E '^--- FAIL' "$OUT" | sed 's/^--- FAIL: //; s/ *([0-9.]*s)//' | sort > "$LOGS/${LABEL}-s$i-$TREE.rednames.body"
  { echo "# RED_NAMES=$(grep -c . "$LOGS/${LABEL}-s$i-$TREE.rednames.body" || true)（★⛔ 0 字节件：零枚红时名册仍写这行说明；尺＝^--- FAIL 计数，同发四数见 four-counts-consumers.tsv）"; cat "$LOGS/${LABEL}-s$i-$TREE.rednames.body"; } > "$LOGS/${LABEL}-s$i-$TREE.rednames.txt"
  rm -f "$LOGS/${LABEL}-s$i-$TREE.rednames.body" 2>/dev/null || true
  echo "  RED_NAMES=$(wc -l < "$LOGS/${LABEL}-s$i-$TREE.rednames.txt")"
done
