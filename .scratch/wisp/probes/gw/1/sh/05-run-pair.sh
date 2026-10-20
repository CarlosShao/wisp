#!/bin/sh
# gw-1 味②⇄③ 成对跑色驱动（★每枚包单独一发＝逐枚 rc／四数无歧义；⛔ 整包 ./cmd/wisp、⛔ ./...）
# 用法：sh 05-run-pair.sh <host|container> <label> <treeB> <treeL> <pkg...>
# 台面（缺任一味＝假红）：
#  host＝go1.27.1 windows/amd64；GOTMPDIR⇄TMPDIR＝cygpath -w 反斜杠形、落 %LOCALAPPDATA% 同卷、先 mkdir -p
#       （★⛔ D:\tmp\wisp 作台根＝109-r1 量到那味：那里带机器范围 BUILTIN\Users 两条继承，
#        会 icacls /remove 再判"继承还在不在"的用例读到的是台根策略⇒假红；读数见 probes/109/r1 第 121-127 行）
#  ★PATH＝cygpath -w 的两枚 DLL 段＋分号（⛔ cygpath -wp 转整条 PATH＝/mingw64/bin 指错⇒CGO 1→0⇒整包死的假红）
#  container＝golang:1.27＋CGO_ENABLED=0＋GOPROXY=off＋MSYS_NO_PATHCONV=1＋GOMODCACHE 只读挂宿主缓存
#       ＋★TMPDIR/GOTMPDIR ⛔ 绑 Windows 卷（容器自 /tmp）
#  ★四数只在 ^=== RUN／^--- PASS／^--- FAIL／^--- SKIP 取、⛔ 混用⛔ 相减；逐名册另落 .rednames.txt
#  ★SHAPE=VOID 当且仅当零枚 ^=== RUN＝台损、⛔ 计入判据；红句含 walk／build failed／setup failed／0xc0000135 同判
set -u
MODE="$1"; LABEL="$2"; TB="$3"; TL="$4"; shift 4
ROOT=/d/tmp/wisp/20261020/gw1
LOGS="$ROOT/logs"
mkdir -p "$LOGS"
PKGS="$*"
i=0
for TREE in "$TB" "$TL"; do
  i=$((i+1))
  case "$TREE" in /*) T="$TREE" ;; *) T="$ROOT/$TREE" ;; esac
  for P in $PKGS; do
    SHORT=$(printf '%s' "$P" | sed 's|^./internal/||; s|^./cmd/||; s|/$||')
    OUT="$LOGS/${LABEL}-s$i-$(basename "$T")-${SHORT}.raw.txt"
    { echo "=== gw-1 $LABEL s$i tree=$(basename "$T") pkg=$P MODE=$MODE @ $(date '+%F %T %z')"
      } > "$OUT"
    if [ "$MODE" = container ]; then
      MSYS_NO_PATHCONV=1 MSYS2_ARG_CONV_EXCL='*' docker run --rm \
        -v "$(cygpath -w "$T"):/src:rw" \
        -v 'D:\work\base\gopath\pkg\mod:/gm:ro' \
        -v "$(cygpath -w "$ROOT/gocache-col-$(basename "$T")"):/gobuild:rw" -w /src \
        -e GOPROXY=off -e GOFLAGS=-mod=mod -e CGO_ENABLED=0 -e GOMODCACHE=/gm \
        -e GOTMPDIR=/tmp -e TMPDIR=/tmp \
        golang:1.27 \
        go test -v -count=1 -timeout 30m "$P" >> "$OUT" 2>&1
      RC=$?
      GT='(container /tmp)'
    else
      TS=$(date +%Y%m%d-%H%M%S)-$$
      GOTMPW="${LOCALAPPDATA:-C:\\Users\\swq\\AppData\\Local}\\wisp-test-gotmp\\gw1-$TS"
      mkdir -p "$(cygpath -u "$GOTMPW")" || { echo "CANNOT_CREATE_GOTMP=$GOTMPW" >> "$OUT"; exit 7; }
      TW=$(cygpath -w "$T")
      DLLP=$(cygpath -w "$T/build"); SHP=$(cygpath -w "$T/third_party/sherpa-onnx")
      ( cd "$TW" && PATH="$SHP;$DLLP;$PATH" GOTMPDIR="$GOTMPW" TMPDIR="$GOTMPW" \
          GOCACHE="D:\\tmp\\wisp\\20261020\\gw1\\gocache-win" \
          go test -v -count=1 -timeout 30m "$P" ) >> "$OUT" 2>&1
      RC=$?
      GT="$GOTMPW"
    fi
    RUN=$(grep -c '^=== RUN' "$OUT" || true)
    PT=$(grep -c '^--- PASS' "$OUT" || true)
    FT=$(grep -c '^--- FAIL' "$OUT" || true)
    SK=$(grep -c '^--- SKIP' "$OUT" || true)
    SHAPE=OK; [ "$RUN" = 0 ] && SHAPE=VOID
    VOIDKIND=none
    if [ "$SHAPE" = VOID ]; then
      grep -qE '0xc0000135' "$OUT" && VOIDKIND=0xc0000135
      grep -qE 'build failed|\[build failed\]' "$OUT" && VOIDKIND="${VOIDKIND}build-failed"
      grep -qiE 'walk|setup failed|no test files' "$OUT" && VOIDKIND="${VOIDKIND}walk-or-setup"
      [ "$VOIDKIND" = none ] && VOIDKIND=unknown
    fi
    LINE="$LABEL s$i $(basename "$T") $P MODE=$MODE RC=$RC RUN=$RUN PASS_TOP=$PT FAIL_TOP=$FT SKIP=$SK SHAPE=$SHAPE VOIDKIND=$VOIDKIND GOTMP=$GT"
    echo "$LINE"
    echo "$LINE" >> "$LOGS/four-counts.tsv"
    grep -E '^--- FAIL' "$OUT" | sed 's/^--- FAIL: //; s/ *([0-9.]*s)$//' | sort > "$LOGS/${LABEL}-s$i-$(basename "$T")-${SHORT}.rednames.body"
    { echo "# gw-1 RED_NAMES=$(grep -c . "$LOGS/${LABEL}-s$i-$(basename "$T")-${SHORT}.rednames.body" || true) 尺=^--- FAIL 计数（同发四数见 four-counts.tsv；零枚时本行即说明＝⛔ 0 字节件）"
      cat "$LOGS/${LABEL}-s$i-$(basename "$T")-${SHORT}.rednames.body"; } > "$LOGS/${LABEL}-s$i-$(basename "$T")-${SHORT}.rednames.txt"
    rm -f "$LOGS/${LABEL}-s$i-$(basename "$T")-${SHORT}.rednames.body" 2>/dev/null || true
    echo "  RED_NAMES=$(wc -l < "$LOGS/${LABEL}-s$i-$(basename "$T")-${SHORT}.rednames.txt")"
    tail -2 "$OUT" | sed 's/^/   TAIL| /'
  done
done
echo "== 逐名作差（每枚包分别 comm，⛔ 混包）=="
for P in $PKGS; do
  SHORT=$(printf '%s' "$P" | sed 's|^./internal/||; s|^./cmd/||; s|/$||')
  B="$LOGS/${LABEL}-s1-$(basename "$TB")-${SHORT}.rednames.txt"
  L="$LOGS/${LABEL}-s2-$(basename "$TL")-${SHORT}.rednames.txt"
  grep -v '^#' "$B" | sort > "$B.b.only"; grep -v '^#' "$L" | sort > "$L.l.only"
  echo "-- $P NEW_RED(only in landed)=$(comm -13 "$B.b.only" "$L.l.only" | grep -c . || true) GONE_RED(only in base)=$(comm -23 "$B.b.only" "$L.l.only" | grep -c . || true)"
  comm -13 "$B.b.only" "$L.l.only" | sed 's/^/   NEWRED /'
  comm -23 "$B.b.only" "$L.l.only" | sed 's/^/   GONERED /'
  rm -f "$B.b.only" "$L.l.only"
done
