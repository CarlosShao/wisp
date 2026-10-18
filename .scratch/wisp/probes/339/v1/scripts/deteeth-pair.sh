#!/bin/bash
# 339-v1 台件④：去牙反控（AC#2），只打仓外导出树物理字节（⛔ 跟踪件、⛔ -overlay）。
# ★成对两发落点不同（公共块口径）＝两枚内容锚：
#   落点A＝aRelates 尾行 `return uncertainAnchorMiss(cand, anchor)`＝回落支整支摘掉；
#   落点B＝aRuleReason 短路首行那半支 `anchorMatches(cand, anchor, rel) || !uncertainAnchorMiss(cand, anchor)`
#          改成恒真＝那句"实话"永不产出（判级不动、只让 reason 回到假话）。
# 平台：落点A 打 Linux 容器＋Windows 本机各一发（＝本格要点名的"两平台"那一维）；落点B 打 Windows 一发。
# 定式＝落刀前旧锚命中＝1、落刀后旧锚＝0 且新锚＝1；跑前跑后各验一次 pristine 哈希；三发串行。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
ANCHOR=7437ef882fac5b427bf17c4bf590819c58d0b375
LOGDIR="$REPO/.scratch/wisp/probes/339/v1/logs"
RUNPAT='TestAnchorMissReasonNamesCandidateAndMissingAnchor|TestAnchorMissKeepsNonOverridableClassAndNeedL2|TestAnchorMissHasNoWorldOnThisPlatform339|TestAnchorFallbackTeethPerAnchor336'
mkdir -p "$LOGDIR"
h() { git -C "$REPO" hash-object "$1"; }
PRISTINE=$(git -C "$REPO" rev-parse "$ANCHOR:internal/risk/blacklist.go")
OLD_A='	return uncertainAnchorMiss(cand, anchor)'
NEW_A='	return false // MUT339V1-A'
OLD_B='	if anchorMatches(cand, anchor, rel) || !uncertainAnchorMiss(cand, anchor) {'
NEW_B='	if true || anchorMatches(cand, anchor, rel) || !uncertainAnchorMiss(cand, anchor) { // MUT339V1-B'

mutate() { # $1=tree $2=A|B
  local TREE="$1" KIND="$2" BL O N
  BL="$TREE/internal/risk/blacklist.go"
  if [ "$KIND" = A ]; then O="$OLD_A"; N="$NEW_A"
    perl -0pi -e 's/\treturn uncertainAnchorMiss\(cand, anchor\)/\treturn false \/\/ MUT339V1-A/' "$BL"
  else O="$OLD_B"; N="$NEW_B"
    perl -0pi -e 's/\tif anchorMatches\(cand, anchor, rel\) \|\| !uncertainAnchorMiss\(cand, anchor\) \{/\tif true || anchorMatches(cand, anchor, rel) || !uncertainAnchorMiss(cand, anchor) { \/\/ MUT339V1-B/' "$BL"
  fi
  echo "KIND=$KIND MUT=$(h "$BL")"
  echo "old_hits_before_check=$(grep -cF "$O" "$BL" || true) (must be 0 after)"
  echo "new_hits_after=$(grep -cF "$N" "$BL" || true) (must be 1)"
}
precheck() { # $1=tree $2=A|B
  local TREE="$1" KIND="$2" BL O
  BL="$TREE/internal/risk/blacklist.go"
  if [ "$KIND" = A ]; then O="$OLD_A"; else O="$OLD_B"; fi
  echo "P0=$(h "$BL") expect=$PRISTINE"
  [ "$(h "$BL")" = "$PRISTINE" ] || { echo GUARD-ABORT-pristine; return 2; }
  echo "old_hits_before=$([ "$(grep -cF "$O" "$BL" || true)" = 1 ] && echo 1 || echo BAD) (must be 1)"
}
restore() {
  local BL="$1/internal/risk/blacklist.go"
  git -C "$REPO" show "$ANCHOR:internal/risk/blacklist.go" > "$BL"
  echo "R=$(h "$BL") match=$([ "$(h "$BL")" = "$PRISTINE" ] && echo YES || echo NO)"
}
four() { local F="$1"; echo "$F RUN=$(grep -c '^=== RUN' "$F" || true) PASS_TOP=$(grep -c '^--- PASS' "$F" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$F" || true) SKIP=$(grep -c '^--- SKIP' "$F" || true)"; grep -E '^--- FAIL' "$F" || true; }

# ── 发①落点A ⇄ LINUX 容器 ──
{
  echo "### SHOT1 landingA(LINUX) run=$RUNPAT"
  precheck /d/tmp/wisp/20261018/339v1/dtl/tree A && {
    mutate /d/tmp/wisp/20261018/339v1/dtl/tree A
    MSYS_NO_PATHCONV=1 docker run --rm -v 'D:/tmp/wisp/20261018/339v1/dtl/tree':/src \
      -v 'D:/work/base/gopath/pkg/mod':/go/pkg/mod:ro -w /src \
      -e GOFLAGS=-mod=mod -e GOPROXY=off -e CGO_ENABLED=0 -e HOME=/src -e GOMODCACHE=/go/pkg/mod \
      golang:1.27 bash -c "mkdir -p /tmp/gt; export TMPDIR=/tmp/gt GOTMPDIR=/tmp/gt; go test ./internal/risk/ -run '$RUNPAT' -count=1 -timeout 30m -v 2>&1; echo RC=\$?"
    restore /d/tmp/wisp/20261018/339v1/dtl/tree
  }
} > "$LOGDIR/v1-deteeth-A-linux.log" 2>&1
cat "$LOGDIR/v1-deteeth-A-linux.log"; echo "----"

# ── 发②落点A ⇄ WINDOWS 本机 ──
{
  echo "### SHOT2 landingA(WINDOWS) run=$RUNPAT"
  precheck /d/tmp/wisp/20261018/339v1/dtw/tree A && {
    mutate /d/tmp/wisp/20261018/339v1/dtw/tree A
    ( cd /d/tmp/wisp/20261018/339v1/dtw/tree || exit 9
      mkdir -p /d/tmp/wisp/20261018/339v1/gt-win
      TMPDIR=/d/tmp/wisp/20261018/339v1/gt-win GOTMPDIR=/d/tmp/wisp/20261018/339v1/gt-win \
      GOFLAGS=-mod=mod GOPROXY=off CGO_ENABLED=0 \
      go test ./internal/risk/ -run "$RUNPAT" -count=1 -timeout 30m -v 2>&1; echo "RC=$?" )
    restore /d/tmp/wisp/20261018/339v1/dtw/tree
  }
} > "$LOGDIR/v1-deteeth-A-windows.log" 2>&1
cat "$LOGDIR/v1-deteeth-A-windows.log"; echo "----"

# ── 发③落点B ⇄ WINDOWS 本机（另一枚内容锚＝成对两发落点不同）──
{
  echo "### SHOT3 landingB(WINDOWS) run=$RUNPAT"
  precheck /d/tmp/wisp/20261018/339v1/dtb/tree B && {
    mutate /d/tmp/wisp/20261018/339v1/dtb/tree B
    ( cd /d/tmp/wisp/20261018/339v1/dtb/tree || exit 9
      mkdir -p /d/tmp/wisp/20261018/339v1/gt-win
      TMPDIR=/d/tmp/wisp/20261018/339v1/gt-win GOTMPDIR=/d/tmp/wisp/20261018/339v1/gt-win \
      GOFLAGS=-mod=mod GOPROXY=off CGO_ENABLED=0 \
      go test ./internal/risk/ -run "$RUNPAT" -count=1 -timeout 30m -v 2>&1; echo "RC=$?" )
    restore /d/tmp/wisp/20261018/339v1/dtb/tree
  }
} > "$LOGDIR/v1-deteeth-B-windows.log" 2>&1
cat "$LOGDIR/v1-deteeth-B-windows.log"; echo "----"

echo "=== 四数逐枚点名 ==="
for L in v1-deteeth-A-linux v1-deteeth-A-windows v1-deteeth-B-windows; do four "$LOGDIR/$L.log"; done
