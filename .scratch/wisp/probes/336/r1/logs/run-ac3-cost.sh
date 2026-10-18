#!/usr/bin/env bash
# 336-r1 AC#3 乙形代价现量（⛔ 推理代替读数）：把 334-v1 那把 42 行三值对撞尺
# （.scratch/wisp/probes/334/v1/logs/probe-334v1_test.go.txt，哈希 69a70fa4…）原样装进
# 仓外导出树，量三件事：① 跑时长与四数（未突变 ⇄ 摘一枚回落支）② 它能不能在 CI 跑到
# （＝该尺无 build tag ⇒ 两平台都编译；Linux 档发由容器台件另取）③ 维护面＝要写几枚 want。
# 尺＝同一枚文件、同一台树、同一条 -run 名册；突变只打仓外树物理字节，逐发 pristine 哈希。
set -u
REPO="D:/work/workspace/projects plans/Wisp"
ANCHOR="${ANCHOR:?具名提交锚}"
TREE="${TREE:?仓外导出台根}"
LOGDIR="$REPO/.scratch/wisp/probes/336/r1/logs"
GOTMPROOT="${GOTMPROOT:-/d/tmp/wisp/20261018/336r1/gotmp}"
F="$TREE/internal/risk/blacklist.go"
PROBE_SRC="$REPO/.scratch/wisp/probes/334/v1/logs/probe-334v1_test.go.txt"
PRISTINE=$(git -C "$REPO" rev-parse "$ANCHOR:internal/risk/blacklist.go")
export GOPATH='D:\work\base\gopath'
export PATH="$TREE/third_party/sherpa-onnx:$TREE/build:$PATH"
ROSTER='TestProbe334V1PerAnchorChain|TestProbe334V1RealMachine|TestProbe334V1SinkReach|TestProbe334V1FixturePrefix'

echo "### 装尺 probe=$PROBE_SRC"
cp "$PROBE_SRC" "$TREE/internal/risk/zz_probe_336y_test.go"
echo "  probe_blob=$(git -C "$REPO" hash-object "$TREE/internal/risk/zz_probe_336y_test.go") 334-v1 原件哈希=69a70fa454886029aea9f8f716cb372908467339"
echo "  probe_lines=$(wc -l < "$TREE/internal/risk/zz_probe_336y_test.go")  打印行数（无断言＝t.Log）=$(grep -c 't\.Log' "$TREE/internal/risk/zz_probe_336y_test.go" || true)  断言行数=$(grep -cE 't\.(Errorf|Fatalf)' "$TREE/internal/risk/zz_probe_336y_test.go" || true)"

run () { # $1=日志名 $2=额外名册
  local gt="$GOTMPROOT/$1"; mkdir -p "$gt" || { echo "  mkdir_gotmp_rc=NONZERO"; return 1; }
  local t0 t1
  t0=$(date +%s.%N)
  ( cd "$TREE" && GOTMPDIR="$gt" go test ./internal/risk/ -run "$2" -count=1 -timeout 30m -v ) > "$LOGDIR/$1" 2>&1
  local rc=$?; t1=$(date +%s.%N)
  echo "  rc=$rc wall=$(awk -v a="$t0" -v b="$t1" 'BEGIN{printf "%.2fs", b-a}') RUN=$(grep -c '^=== RUN' "$LOGDIR/$1" || true) PASS_TOP=$(grep -c '^--- PASS' "$LOGDIR/$1" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$LOGDIR/$1" || true) SKIP=$(grep -c '^--- SKIP' "$LOGDIR/$1" || true)"
  grep -E '^(ok|FAIL)[ \t]+github' "$LOGDIR/$1" | sed 's/^/    PKGTIME /'
  grep -E '^ *--- FAIL: ' "$LOGDIR/$1" | sed 's/([0-9.]*s)//;s/^/    RED /' || true
}

echo "### 乙-1 未突变（成对正控＝这一尺今天必绿，因为它⛔ 断言）"
git -C "$REPO" show "$ANCHOR:internal/risk/blacklist.go" > "$F"
echo "  hash=$(git -C "$REPO" hash-object "$F") expected=$PRISTINE"
run ac3y-00-pristine-chain.log "$ROSTER"
grep -hoE 'V1SHAPE (chain|offchain):[^ ]* class=[A-Za-z]+ allow=[a-z]+ needL2=[a-z]+ rule=[^ ]*' "$LOGDIR/ac3y-00-pristine-chain.log" > "$LOGDIR/ac3y-00-pristine-chain.shape" || true
echo "  shape_lines=$(wc -l < "$LOGDIR/ac3y-00-pristine-chain.shape")"

echo "### 乙-2 只摘 local-credentials 一枚的回落支（＝AC#2 的 MU-1h 同一落点）"
git -C "$REPO" show "$ANCHOR:internal/risk/blacklist.go" > "$F"
perl -0pi -e 's/case localAppDataS != "" && aUnder\(cand, anchorForms\(anchorPath\(localAppDataS, `microsoft`, `credentials`\)\)\):/case localAppDataS != "" \&\& anchorMatches(cand, anchorForms(anchorPath(localAppDataS, `microsoft`, `credentials`)), isUnder):/' "$F"
echo "  落点命中=$(grep -c 'anchorMatches(cand, anchorForms' "$F")"
run ac3y-01-MU1h-chain.log "$ROSTER"
grep -hoE 'V1SHAPE (chain|offchain):[^ ]* class=[A-Za-z]+ allow=[a-z]+ needL2=[a-z]+ rule=[^ ]*' "$LOGDIR/ac3y-01-MU1h-chain.log" > "$LOGDIR/ac3y-01-MU1h-chain.shape" || true
echo "  尺看到的行差（pristine ⇄ MU-1h，逐字）:"
diff "$LOGDIR/ac3y-00-pristine-chain.shape" "$LOGDIR/ac3y-01-MU1h-chain.shape" | sed 's/^/    /' || true
echo "  ★注意＝这一尺⛔ 断言，所以 rc 仍 0：搬成仓内断言必须逐枚补 want（维护面见 report）"

echo "### 乙-3 甲齿对照（同一发突变下我的牙的 rc／四数）"
run ac3y-02-MU1h-teeth.log 'TestAnchorFallbackTeethPerAnchor336'
git -C "$REPO" show "$ANCHOR:internal/risk/blacklist.go" > "$F"
echo "  final_hash=$(git -C "$REPO" hash-object "$F") expected=$PRISTINE"
echo "  两尺件数与行数＝甲：1 枚 windows 件 $(wc -l < "$TREE/internal/risk/anchor_fallback_teeth_336_windows_test.go") 行＋1 枚 other 件 $(wc -l < "$TREE/internal/risk/anchor_fallback_teeth_336_other_test.go" 2>/dev/null || echo 0) 行；乙：1 枚 $(wc -l < "$TREE/internal/risk/zz_probe_336y_test.go") 行（无断言）"
echo "ac3y_rc=0"
