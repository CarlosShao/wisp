#!/usr/bin/env bash
# 336-r1 AC#4 颜色：./internal/risk 成对两把临时根台面（--ci-shape＝默认 C 卷 %LOCALAPPDATA% 组法
# ⇄ --dshape＝D:\tmp\wisp\20261018\336r1\gotmp\<名>，机主 10-18 新铁律：⛔ 盘符根），
# 同一条 -run 名册，改前（ed8ff551＝无本腿字节）⇄ 改后（8c36c23c＝含本腿牙）同一台面各一发。
# 四数尺：RUN=^=== RUN（含子测试）／PASS_TOP=^--- PASS（只顶层）／FAIL_TOP=^--- FAIL／SKIP=^--- SKIP，
# ⛔ 混用⛔ 相减；一律 -count=1 -timeout 30m；⛔ ./cmd/wisp 整包、⛔ go test ./...。
# 必答(ii) 的跑时长差＝同一条整包名册在两棵树上各两发（go 自报的包级秒数＋本腿牙的逐枚秒数）。
set -u
REPO="D:/work/workspace/projects plans/Wisp"
LOGDIR="$REPO/.scratch/wisp/probes/336/r1/logs"
BEFORE="D:/tmp/wisp/20261018/336r1/tree-before"
AFTER="D:/tmp/wisp/20261018/336r1/tree-after"
GOTMPROOT="D:/tmp/wisp/20261018/336r1/gotmp"
export GOPATH='D:\work\base\gopath'
ROSTER='TestAnchorFallbackTeethPerAnchor336|TestAnchorMissReasonNamesCandidateAndMissingAnchor|TestAnchorMissKeepsNonOverridableClassAndNeedL2|TestFileBaseNameClaimChangesNoClassOrReason'

shot () { # $1=树 $2=名册 $3=日志名 $4=shape
  local tree="$1" roster="$2" log="$3" shape="${4:---dshape}"
  local gt="$GOTMPROOT/$log"
  if [ "$shape" = "--ci-shape" ]; then
    gt="${LOCALAPPDATA:-C:\\Users\\swq\\AppData\\Local}\\wisp-test-gotmp\\portable-336r1-$(date +%Y%m%d-%H%M%S)-$$"
  fi
  mkdir -p "$gt" && mkdir_rc=$? || mkdir_rc=$?
  ( cd "$tree" && PATH="$tree/third_party/sherpa-onnx:$tree/build:$PATH" GOTMPDIR="$gt" \
      go test ./internal/risk/ -run "$roster" -count=1 -timeout 30m -v ) > "$LOGDIR/$log" 2>&1
  local rc=$?
  echo "### $log shape=$shape mkdir_gotmp_rc=$mkdir_rc GOTMPDIR=$gt"
  echo "  tree=$tree blob_blacklist=$(git -C "$REPO" hash-object "$tree/internal/risk/blacklist.go")"
  echo "  teeth_on_tree=$([ -f "$tree/internal/risk/anchor_fallback_teeth_336_windows_test.go" ] && echo present || echo ABSENT)"
  echo "  rc=$rc RUN=$(grep -c '^=== RUN' "$LOGDIR/$log" || true) PASS_TOP=$(grep -c '^--- PASS' "$LOGDIR/$log" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$LOGDIR/$log" || true) SKIP=$(grep -c '^--- SKIP' "$LOGDIR/$log" || true)"
  grep -E '^ *--- FAIL: ' "$LOGDIR/$log" | sed 's/([0-9.]*s)//;s/^/    RED /' || true
  grep -E '^(ok|FAIL|---)\s+github|^FAIL$|^ok$' "$LOGDIR/$log" | sed 's/^/  PKGTIME /' || true
  grep -oE 'TEETH336 operand=\S+ class=\S+ matched=\S+ clause=\S+' "$LOGDIR/$log" | sed 's/^/  /' || true
}

echo "### AC#4-1 成对两把临时根台面（同一条名册，改前⇄改后）"
shot "$BEFORE" "$ROSTER" ac4-01-before-cishape.log --ci-shape
shot "$AFTER"  "$ROSTER" ac4-02-after-cishape.log  --ci-shape
shot "$BEFORE" "$ROSTER" ac4-03-before-dshape.log  --dshape
shot "$AFTER"  "$ROSTER" ac4-04-after-dshape.log   --dshape

echo "### AC#4-2 必答(ii) 跑时长差＝整包 ./internal/risk（-run . ，同台面 --dshape，各两发）"
shot "$BEFORE" '.' ac4-05-full-before-dshape-a.log --dshape
shot "$BEFORE" '.' ac4-05-full-before-dshape-b.log --dshape
shot "$AFTER"  '.' ac4-06-full-after-dshape-a.log  --dshape
shot "$AFTER"  '.' ac4-06-full-after-dshape-b.log  --dshape

echo "### AC#4-3 红名册逐名作差（before ⇄ after，整包）"
for tag in before-dshape-a before-dshape-b after-dshape-a after-dshape-b; do
  grep -E '^ *--- FAIL: ' "$LOGDIR/ac4-0*-full-$tag.log" 2>/dev/null | sed 's/([0-9.]*s)//' | sed 's/^ *//' | sort > "$LOGDIR/ac4-full-$tag.rednames" || true
  echo "  $tag rednames=$(wc -l < "$LOGDIR/ac4-full-$tag.rednames")"
done
echo "  diff before-a ⇄ after-a:"; diff "$LOGDIR/ac4-full-before-dshape-a.rednames" "$LOGDIR/ac4-full-after-dshape-a.rednames" | sed 's/^/    /' || true
echo "  diff before-b ⇄ after-b:"; diff "$LOGDIR/ac4-full-before-dshape-b.rednames" "$LOGDIR/ac4-full-after-dshape-b.rednames" | sed 's/^/    /' || true
echo "  SKIP 逐名（两棵树各一发）:"
for f in ac4-05-full-before-dshape-a.log ac4-06-full-after-dshape-a.log; do
  echo "    $f:"; grep -E '^ *--- SKIP: ' "$LOGDIR/$f" | sed 's/^/      /' || true
done
echo "ac4_rc=0"
