#!/bin/bash
# 编排者亲跑：票 298 AC#4（改前 48e7b0a4 ⇄ 改后 18d43f8e，同一台面＝两棵仓外导出树，各 2 发）
# 规矩：改前那发⛔ 在共享工作树上跑；每把门禁件自落 rc=；0 旁边必须站着分母。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
OUT="$REPO/.scratch/wisp/probes/orch/20261012/logs/298ac4"
mkdir -p "$OUT"
SHERPA="$REPO/third_party/sherpa-onnx"
BUILD="$REPO/build"
RUNTAG=$(date +%Y%m%d-%H%M%S)
ROOT="$LOCALAPPDATA/wisp-test-gotmp/orch-298ac4-$RUNTAG"
mkdir -p "$ROOT" 2>/dev/null || true
ROOT="${ROOT//\\//}"

echo "runtag=$RUNTAG"
echo "gotmp_root=$ROOT"

before=$(mktemp -d)
after=$(mktemp -d)
echo "before_tree=$before"
echo "after_tree=$after"
( cd "$REPO" && git archive --format=tar 48e7b0a4 ) | tar -x -C "$before"
echo "archive_before_rc=$?"
( cd "$REPO" && git archive --format=tar 18d43f8e ) | tar -x -C "$after"
echo "archive_after_rc=$?"

for t in "$before" "$after"; do
  n=$(find "$t" -name '*.go' -type f | wc -l)
  echo "go_files_in_tree=$n tree=$t"
done

# 正控前置：先各自编译一次，编译失败＝台面故障，⛔ 计入牙
for side in before after; do
  t=$( [ $side = before ] && echo "$before" || echo "$after" )
  ( cd "$t" && GOFLAGS= go build ./cmd/wisp/ ) >"$OUT/$side-build.txt" 2>&1
  echo "$side build rc=$?" | tee -a "$OUT/summary.txt"
done

for side in before after; do
  t=$( [ $side = before ] && echo "$before" || echo "$after" )
  for shot in 1 2; do
    live=$(tasklist //FI "IMAGENAME eq wisp.exe" 2>/dev/null | grep -c 'wisp.exe' || true)
    echo "shot=$side-$shot wispexe_before=$live"
    ( cd "$t" && GOTMPDIR="$ROOT" PATH="$SHERPA:$BUILD:$PATH" GOFLAGS= go test ./cmd/wisp/ -count=1 ) \
      >"$OUT/$side-$shot.txt" 2>&1
    rc=$?
    ok=$(grep -c '^ok' "$OUT/$side-$shot.txt" || true)
    fail=$(grep -c '^FAIL' "$OUT/$side-$shot.txt" || true)
    dfail=$(grep -c -- '--- FAIL' "$OUT/$side-$shot.txt" || true)
    size=$(wc -c < "$OUT/$side-$shot.txt")
    echo "shot=$side-$shot rc=$rc ok=$ok fail=$fail dashfail=$dfail bytes=$size" | tee -a "$OUT/summary.txt"
  done
done

# 两把尺并排（工作树 ⇄ blob），只报枚数、⛔ 解读
for side in before after; do
  t=$( [ $side = before ] && echo "$before" || echo "$after" )
  gf=$(cd "$t" && gofmt -l cmd/wisp 2>/dev/null | tr '\n' ' ')
  echo "gofmt_l_cmd_wisp tree=$side gfn=$([ -n "$gf" ] && echo "$gf" | wc -w || echo 0) names=[$gf]" | tee -a "$OUT/summary.txt"
done
GOFUMPT="$(go env GOPATH | sed 's|\\|/|g')/bin/gofumpt.exe"
for side in before after; do
  t=$( [ $side = before ] && echo "$before" || echo "$after" )
  if [ -x "$GOFUMPT" ]; then
    gf=$( (cd "$t" && "$GOFUMPT" -l cmd/wisp 2>/dev/null) | tr '\n' ' ')
    echo "gofumpt_l_cmd_wisp tree=$side names=[$gf]" | tee -a "$OUT/summary.txt"
  else
    echo "gofumpt NOT_FOUND at $GOFUMPT" | tee -a "$OUT/summary.txt"
  fi
done
echo "done at $(date)" | tee -a "$OUT/summary.txt"
