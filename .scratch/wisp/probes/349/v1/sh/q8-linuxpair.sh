#!/bin/sh
# 349-v1 AC#4 Linux 侧成对（配方＝公共块台面配方＋343/349 之先例：docker golang:1.27，
# 只挂我自己的树；⛔ 复用 349-r1/343-v1 的树）。两发串行。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
OUTDIR=/d/tmp/wisp/20261019/349v1/logsLin
mkdir -p "$OUTDIR"
CACHE="$(cd "$REPO" && /d/work/base/go/bin/go.exe env GOMODCACHE | tr -d '\r')"
one () { # $1=tree-dir $2=label
  local T="/d/tmp/wisp/20261019/349v1/$1" LABEL="$2"
  local OUT="$OUTDIR/lin-$LABEL.log"
  echo "tree=$T label=$LABEL modcache=$CACHE"
  MSYS_NO_PATHCONV=1 docker run --rm \
    -v "$T":/src \
    -v "$CACHE":/go/pkg/mod:ro \
    -w /src \
    -e GOFLAGS=-mod=mod -e GOPROXY=off -e CGO_ENABLED=0 -e HOME=/src -e GOMODCACHE=/go/pkg/mod \
    golang:1.27 bash -c "mkdir -p /tmp/gt; export TMPDIR=/tmp/gt GOTMPDIR=/tmp/gt; go test ./internal/risk/ -count=1 -timeout 30m -v > /tmp/gt/out.txt 2>&1; echo GO_TEST_RC=\$?; cat /tmp/gt/out.txt" > "$OUT" 2>&1
  echo "== label=$LABEL bytes=$(wc -c < "$OUT") =="
  echo "RUN=$(grep -c '^=== RUN' "$OUT" || true)"
  echo "PASS_TOP=$(grep -c '^--- PASS' "$OUT" || true)"
  echo "FAIL_TOP=$(grep -c '^--- FAIL' "$OUT" || true)"
  echo "SKIP=$(grep -c '^--- SKIP' "$OUT" || true)"
  echo "-- 红名册 --"
  grep -E '^--- FAIL' "$OUT" || echo "(zero top-level reds)"
  echo "-- 台面故障名册 --"
  grep -nE 'build failed|no test files|setup failed|creating work dir|0xc0000135' "$OUT" || echo "none"
}
{
echo "### 349-v1 Linux pair (串行: base 先、del 后) $(date -Iseconds)"
one treeBase base
one treeDel  del
echo "### 红名册逐名作差"
grep -hE '^--- FAIL' "$OUTDIR/lin-base.log" | sed 's/([0-9.]*s)//;s/^--- FAIL: //;s/ *$//' | sort > "$OUTDIR/lin-red-base.txt"
grep -hE '^--- FAIL' "$OUTDIR/lin-del.log"  | sed 's/([0-9.]*s)//;s/^--- FAIL: //;s/ *$//' | sort > "$OUTDIR/lin-red-del.txt"
wc -l "$OUTDIR/lin-red-base.txt" "$OUTDIR/lin-red-del.txt"
echo "-- only in base --"; comm -23 "$OUTDIR/lin-red-base.txt" "$OUTDIR/lin-red-del.txt" || true
echo "-- only in del (want empty) --"; comm -13 "$OUTDIR/lin-red-base.txt" "$OUTDIR/lin-red-del.txt" || true
echo "### LIN_DONE $(date -Iseconds)"
} > "$OUTDIR/lin-pair.log" 2>&1
echo "wrote=$OUTDIR/lin-pair.log"
