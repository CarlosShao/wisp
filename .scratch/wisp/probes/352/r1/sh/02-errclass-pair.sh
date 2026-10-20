#!/bin/sh
# 352-r1 AC#0/AC#1 的尺：同一枚仓外探针件打进【基线树拷贝】⇄【落地单差异树拷贝】，两平台各 ≥2 发。
# 用法：sh 02-errclass-pair.sh <host|container> <tree slot1> <tree slot2> <label>
#   slot1＝基线侧、slot2＝落地侧（AC#0 阶段两枚槽都给 treeB＝同一棵 pristine 树复量两发）。
# ★两枚槽各自独立的 mut 树与独立的 raw 件＝⛔ 后一发覆盖前一发。
# ★台面对齐（第十五波公共块 台面四味）：
#   host＝出货平台 go1.27.1 windows/amd64；GOTMPDIR/TMPDIR＝反斜杠形＋落 %LOCALAPPDATA% 同卷短根；
#     ★⛔ 用 cygpath -wp 转整条 PATH（那会把无盘符 MSYS 段指错⇒CGO 1→0⇒零枚 --- FAIL 的整包死）。
#   container＝golang:1.27 + CGO_ENABLED=0 + GOPROXY=off + 宿主 module 缓存只读；
#     ★TMPDIR/GOTMPDIR ⛔ 绑 Windows 卷（os.SameFile 失真）＝容器自 /tmp。
#     ★MSYS_NO_PATHCONV=1：不然 msys 把 /src 换成 E:/…/msys64/src⇒RC=125／RUN=0＝台损（本腿首发实测）。
# ★退码＝重定向到文件 + 取 $? 紧跟裸命令（⛔ 管进 tail）；串行，⛔ 并行跑改文件的脚本。
# ★SHAPE＝VOID 当且仅当零枚 === RUN（＝这发⛔ 到判据，台损，⛔ 计入判据）。
set -u
MODE="$1"; T1="$2"; T2="$3"; LABEL="$4"
ROOT=/d/tmp/wisp/20261020/352r1
ROOTW='D:\tmp\wisp\20261020\352r1'
LOGS="$ROOT/logs"
PROBE="$ROOT/probe/probe-352r1-errclass_test.go.txt"
[ -f "$PROBE" ] || { echo "PROBE_MISSING=$PROBE"; exit 7; }
mkdir -p "$LOGS"
i=0
for TREE in "$T1" "$T2"; do
  i=$((i+1))
  MUT="mut-$i-$TREE"
  [ -d "$ROOT/$MUT" ] || cp -a "$ROOT/$TREE" "$ROOT/$MUT"
  cp -f "$PROBE" "$ROOT/$MUT/internal/winsec/probe_352r1_err_test.go"
  OUT="$LOGS/${LABEL}-s${i}-${MUT}.raw.txt"
  PRE=$(sha256sum "$ROOT/$MUT/internal/winsec/winsec.go" | cut -d' ' -f1)
  if [ "$MODE" = container ]; then
    MSYS_NO_PATHCONV=1 MSYS2_ARG_CONV_EXCL='*' docker run --rm \
      -v "$(cygpath -w "$ROOT/$MUT"):/src:rw" \
      -v 'D:\work\base\gopath\pkg\mod:/gm:ro' \
      -v "$(cygpath -w "$ROOT/gocache-col"):/gobuild:rw" -w /src \
      -e GOPROXY=off -e GOFLAGS=-mod=mod -e CGO_ENABLED=0 -e GOMODCACHE=/gm \
      -e GOTMPDIR=/tmp -e TMPDIR=/tmp \
      golang:1.27 \
      go test ./internal/winsec/ -v -count=1 -timeout 20m -run 'TestP352R1ErrClassify' > "$OUT" 2>&1
    RC=$?
    POST_GOTMP="(container /tmp)"
  else
    TS=$(date +%Y%m%d-%H%M%S)-$$
    GOTMPW="${LOCALAPPDATA:-C:\\Users\\swq\\AppData\\Local}\\wisp-test-gotmp\\352r1-$TS"
    mkdir -p "$(cygpath -u "$GOTMPW")" || { echo "CANNOT_CREATE_GOTMP=$GOTMPW"; exit 7; }
    ( cd "$ROOTW\\$MUT" && GOTMPDIR="$GOTMPW" TMPDIR="$GOTMPW" \
        GOCACHE='D:\tmp\wisp\20261020\352r1\gocache-win' \
        go test ./internal/winsec/ -v -count=1 -timeout 20m -run 'TestP352R1ErrClassify' ) > "$OUT" 2>&1
    RC=$?
    echo "GOTMPDIR_USED=$GOTMPW"
    POST_GOTMP="$GOTMPW"
  fi
  POST=$(sha256sum "$ROOT/$MUT/internal/winsec/winsec.go" | cut -d' ' -f1)
  RESTORE=OK; [ "$PRE" = "$POST" ] || RESTORE="DRIFT $PRE -> $POST"
  RUN=$(grep -c '^=== RUN' "$OUT" || true)
  PT=$(grep -c '^--- PASS' "$OUT" || true)
  FT=$(grep -c '^--- FAIL' "$OUT" || true)
  SK=$(grep -c '^--- SKIP' "$OUT" || true)
  SHAPE=OK; [ "$RUN" = 0 ] && SHAPE=VOID
  LINE="$LABEL s$i $MUT MODE=$MODE RC=$RC RUN=$RUN PASS_TOP=$PT FAIL_TOP=$FT SKIP=$SK SHAPE=$SHAPE WINSEC_SHA256=$POST RESTORE=$RESTORE GOTMP=$POST_GOTMP"
  echo "$LINE"
  echo "$LINE" >> "$LOGS/four-counts-errclass.tsv"
  grep -E ': R[0-9]+ ' "$OUT"
  grep -E '^--- (FAIL|SKIP)' "$OUT" || true
  echo "-----"
done
