#!/bin/sh
# gw-1 味② 撞红复量＝公共块令"撞红先安静复量 3 发再归因、⛔ 记成别腿新造红、⛔ 把摆动当闭环"。
# 首发那枚逐名换位＝container 档 internal/config：基线侧红 `TestHandDeletedSchemaVersion…` ⇄ 落地侧红
# `TestNothingWrittenPromiseOnlyHoldsWhenSchemaIsCurrent`＝两枚都＝公共块在册的 internal/config mtime 族。
# ★两侧各 3 发、串行（⛔ 并行跑改文件的脚本）；名册取两侧**并集**作差＝新增红只认并集 comm。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 9
R=/d/tmp/wisp/20261020/gw1
LOGS="$R/logs"
mkdir -p "$LOGS"
echo "== START anchor =="; git rev-parse HEAD | sed 's|^|HEAD_FULL=|'
for SHOT in r1 r2 r3; do
  for TREE in gwB gwL; do
    OUT="$LOGS/cfg-$SHOT-$TREE.raw.txt"
    MSYS_NO_PATHCONV=1 MSYS2_ARG_CONV_EXCL='*' docker run --rm \
      -v "$(cygpath -w "$R/$TREE"):/src:rw" \
      -v 'D:\work\base\gopath\pkg\mod:/gm:ro' \
      -v "$(cygpath -w "$R/gocache-col-$TREE"):/gobuild:rw" -w /src \
      -e GOPROXY=off -e GOFLAGS=-mod=mod -e CGO_ENABLED=0 -e GOMODCACHE=/gm \
      -e GOTMPDIR=/tmp -e TMPDIR=/tmp \
      golang:1.27 \
      go test -v -count=1 -timeout 30m ./internal/config/ > "$OUT" 2>&1
    RC=$?
    LINE="cfg $SHOT $TREE MODE=container RC=$RC RUN=$(grep -c '^=== RUN' "$OUT" || true) PASS_TOP=$(grep -c '^--- PASS' "$OUT" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$OUT" || true) SKIP=$(grep -c '^--- SKIP' "$OUT" || true)"
    echo "$LINE"; echo "$LINE" >> "$LOGS/four-counts.tsv"
    grep -E '^--- FAIL' "$OUT" | sed 's/^--- FAIL: //; s/ *([0-9.]*s)$//' | sort > "$LOGS/cfg-$SHOT-$TREE.rednames.txt"
    echo "  RED=$(grep -c . "$LOGS/cfg-$SHOT-$TREE.rednames.txt" || true)"
    grep . "$LOGS/cfg-$SHOT-$TREE.rednames.txt" | sed 's|^|   RED |'
  done
  echo "  -- 并集作差 $SHOT --"
  comm -13 "$LOGS/cfg-$SHOT-gwB.rednames.txt" "$LOGS/cfg-$SHOT-gwL.rednames.txt" | sed 's|^|   NEWRED |'
  comm -23 "$LOGS/cfg-$SHOT-gwB.rednames.txt" "$LOGS/cfg-$SHOT-gwL.rednames.txt" | sed 's|^|   GONERED |'
done
echo
echo "== ★三发并集（⛔ 拿单发换位当新增红）=="
cat "$LOGS"/cfg-r1-gwB.rednames.txt "$LOGS"/cfg-r2-gwB.rednames.txt "$LOGS"/cfg-r3-gwB.rednames.txt | sort -u > "$LOGS/cfg-union-B.txt"
cat "$LOGS"/cfg-r1-gwL.rednames.txt "$LOGS"/cfg-r2-gwL.rednames.txt "$LOGS"/cfg-r3-gwL.rednames.txt | sort -u > "$LOGS/cfg-union-L.txt"
echo "UNION_B=$(grep -c . "$LOGS/cfg-union-B.txt" || true)"; grep . "$LOGS/cfg-union-B.txt" | sed 's|^|  B |'
echo "UNION_L=$(grep -c . "$LOGS/cfg-union-L.txt" || true)"; grep . "$LOGS/cfg-union-L.txt" | sed 's|^|  L |'
echo "NEW_RED_UNION=$(comm -13 "$LOGS/cfg-union-B.txt" "$LOGS/cfg-union-L.txt" | grep -c . || true)"
comm -13 "$LOGS/cfg-union-B.txt" "$LOGS/cfg-union-L.txt" | sed 's|^|  NEWRED_UNION |'
echo "SAME_UNION=$(cmp -s "$LOGS/cfg-union-B.txt" "$LOGS/cfg-union-L.txt" && echo YES || echo NO)"
echo
echo "== 逐字红句（两侧各取一枚，供裁腿对判据形状）=="
for f in "$LOGS"/cfg-r1-gwB.raw.txt "$LOGS"/cfg-r1-gwL.raw.txt; do
  echo "-- $f --"
  grep -E '^\s+[a-z_]+\.go:[0-9]+: ' "$f" | grep -iE 'mtime|schema' | head -4 | sed 's|^|   REDLINE |'
done
echo "== END anchor =="; git rev-parse HEAD | sed 's|^|HEAD_FULL_END=|'
