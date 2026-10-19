#!/bin/sh
# 349-r1 AC#4 Linux 侧成对（配方逐字抄自 343-r1 的 linux-343.sh，只换树/日志落点＝⛔ 复用别腿的树）。
# 两棵树都是我自己的 git archive 导出：treeBase＝起手锚 7b633fa9、treeDel＝交件锚 db2f9eab。
# 台损判据：红句含 walk／build failed／no test files／setup failed／creating work dir／0xc0000135
# ⇒ 具名写是哪一行、⛔ 计入判据。两发串行。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
OUTDIR="$REPO/.scratch/wisp/probes/349/r1/logs/pair"
mkdir -p "$OUTDIR"
CACHE="$(cd "$REPO" && go env GOMODCACHE)"
one () { # $1=tree-dir $2=label
  local T="$1" LABEL="$2"
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
  echo "-- SKIP 名册 --"
  grep -E '^--- SKIP' "$OUT" || echo "(none)"
  echo "-- GO_TEST_RC --"
  grep -E '^GO_TEST_RC=' "$OUT" || true
  echo "-- 台面故障名册 --"
  grep -nE 'walk|build failed|no test files|setup failed|creating work dir|0xc0000135' "$OUT" || echo "none"
}
{
echo "### 349-r1 Linux pair (串行: base 先、del 后)"
one /d/tmp/wisp/20261019/349r1/treeBase base
one /d/tmp/wisp/20261019/349r1/treeDel  del
echo "### 红名册逐名作差（lin base ⇄ lin del）"
grep -hE '^--- FAIL' "$OUTDIR/lin-base.log" | sed 's/([0-9.]*s)//;s/^--- FAIL: //;s/ *$//' | sort > "$OUTDIR/lin-red-base.txt"
grep -hE '^--- FAIL' "$OUTDIR/lin-del.log"  | sed 's/([0-9.]*s)//;s/^--- FAIL: //;s/ *$//' | sort > "$OUTDIR/lin-red-del.txt"
wc -l "$OUTDIR/lin-red-base.txt" "$OUTDIR/lin-red-del.txt"
echo "-- only in base --"; comm -23 "$OUTDIR/lin-red-base.txt" "$OUTDIR/lin-red-del.txt" || true
echo "-- only in del (want empty) --"; comm -13 "$OUTDIR/lin-red-base.txt" "$OUTDIR/lin-red-del.txt" || true
} > "$OUTDIR/lin-pair-runner.log" 2>&1
echo "wrote=$OUTDIR/lin-pair-runner.log"
cat "$OUTDIR/lin-pair-runner.log"
