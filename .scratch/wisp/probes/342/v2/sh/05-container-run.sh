#!/bin/sh
# 342-v2 容器驱动（POSIX 档）：在仓外导出树里跑指定 -run，落 raw log＋四数＋liveness 分子逐发点名。
# ★TMPDIR/GOTMPDIR ⛔ 绑 Windows 卷（342-v1 的定式：os.SameFile 在 bind mount 失真＝身份门被台面摘掉）。
# 用法：sh 05-container-run.sh <treeB|treeL> <label> [extra go test args]
set -e
TREE="$1"; LABEL="$2"; shift 2
ROOT=/d/tmp/wisp/20261019/342v2
LOGDIR="$ROOT/logs"; mkdir -p "$LOGDIR"
RAW="$LOGDIR/$LABEL.raw.txt"
docker run --rm \
  -v "$ROOT/$TREE:/src" \
  -v /d/work/base/gopath/pkg/mod:/gm:ro \
  -w /src \
  -e GOPROXY=off -e GOFLAGS=-mod=mod -e CGO_ENABLED=0 -e GOMODCACHE=/gm \
  golang:1.27 \
  go test ./internal/winsec/ -v -count=1 -timeout 40m "$@" > "$RAW" 2>&1
RC=$?
echo "RC=$RC LABEL=$LABEL TREE=$TREE"
RUN=$(grep -c '^=== RUN' "$RAW" || true)
PT=$(grep -c '^--- PASS' "$RAW" || true)
FT=$(grep -c '^--- FAIL' "$RAW" || true)
SK=$(grep -c '^--- SKIP' "$RAW" || true)
echo "FOUR RUN=$RUN PASS_TOP=$PT FAIL_TOP=$FT SKIP=$SK"
printf '%s\t%s\tRUN=%s PASS_TOP=%s FAIL_TOP=%s SKIP=%s RC=%s\n' "$LABEL" "$TREE" "$RUN" "$PT" "$FT" "$SK" "$RC" >> "$LOGDIR/four-counts.tsv"
echo '--- 分子逐发（truncated/deleted/pinnedParentBeforeOpen/provedLinkCycles/narrowed）---'
grep -E 'AC#1 \(ticket 342|AC#2 \(ticket 348' "$RAW" | grep -E 'attempts=' || true
echo '--- RED 句逐字 ---'
grep -E '^\s+.*AC#(1|2) RED' "$RAW" || true
echo '--- 红名册 ---'
grep -E '^--- FAIL' "$RAW" || true
