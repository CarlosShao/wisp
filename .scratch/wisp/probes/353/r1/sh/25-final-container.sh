#!/bin/sh
# 353-r1 交件前第二平台＝container（golang:1.27 linux/amd64，CGO_ENABLED=0、GOPROXY=off）打成对两棵终量树。
# ★TMPDIR/GOTMPDIR ⛔ 绑 Windows 卷（容器自 /tmp）；★MSYS_NO_PATHCONV=1；★⛔ ./cmd/wisp；★⛔ ./...
set -u
R=/d/tmp/wisp/20261020/353r1
LOGS="$R/logs"
i=0
for TREE in finalB finalL; do
  i=$((i+1))
  OUT="$LOGS/final-col-s$i-$TREE.raw.txt"
  { echo "=== 353r1 finalcol s$i tree=$TREE MODE=container @ $(date '+%F %T %z')"; } > "$OUT"
  MSYS_NO_PATHCONV=1 MSYS2_ARG_CONV_EXCL='*' docker run --rm \
    -v "$(cygpath -w "$R/$TREE"):/src:rw" \
    -v 'D:\work\base\gopath\pkg\mod:/gm:ro' \
    -v "$(cygpath -w "$R/gocache-col-$TREE"):/gobuild:rw" -w /src \
    -e GOPROXY=off -e GOFLAGS=-mod=mod -e CGO_ENABLED=0 -e GOMODCACHE=/gm \
    -e GOTMPDIR=/tmp -e TMPDIR=/tmp \
    golang:1.27 \
    go test -v -count=1 -timeout 60m ./internal/winsec/ >> "$OUT" 2>&1
  RC=$?
  LINE="finalcol s$i $TREE MODE=container RC=$RC RUN=$(grep -c '^=== RUN' "$OUT" || true) PASS_TOP=$(grep -c '^--- PASS' "$OUT" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$OUT" || true) SKIP=$(grep -c '^--- SKIP' "$OUT" || true) SHAPE=$([ "$(grep -c '^=== RUN' "$OUT" || true)" = 0 ] && echo VOID || echo OK)"
  echo "$LINE"; echo "$LINE" >> "$LOGS/final-fourcounts.tsv"
  grep -E '^--- FAIL' "$OUT" | sed 's/^--- FAIL: //; s/ *([0-9.]*s)$//' | sort > "$LOGS/final-col-s$i.rednames.txt"
  echo "  RED_NAMES=$(grep -c . "$LOGS/final-col-s$i.rednames.txt" || true)"
  cat "$LOGS/final-col-s$i.rednames.txt" | sed 's/^/   RED /'
  tail -2 "$OUT" | sed 's/^/   TAIL| /'
done
echo "== container 侧逐名作差（尺＝名集）=="
comm -13 "$LOGS/final-col-s1.rednames.txt" "$LOGS/final-col-s2.rednames.txt" | sed 's/^/  ONLY_IN_LANDED: /' || true
echo "  NEW_RED_COUNT=$(comm -13 "$LOGS/final-col-s1.rednames.txt" "$LOGS/final-col-s2.rednames.txt" | grep -c . || true)"
