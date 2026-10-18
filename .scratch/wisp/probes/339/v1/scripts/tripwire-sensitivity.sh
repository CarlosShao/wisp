#!/bin/bash
# 339-v1 台件③：绊线敏感尺（必答 ii／AC#2 POSIX 半区）。
# 落点＝仓外导出树的物理字节（⛔ 跟踪件、⛔ -overlay）；内容锚＝函数体内那枚 token（⛔ 行号）。
# 定式＝落刀前量命中枚数＝1、落刀后新锚＝1 且旧锚＝0；跑前跑后各验一次 pristine 哈希；两发串行。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
ANCHOR=7437ef882fac5b427bf17c4bf590819c58d0b375
TREE="/d/tmp/wisp/20261018/339v1/trip/tree"
PR="$TREE/internal/risk/pathresolver_other.go"
LOGDIR="$REPO/.scratch/wisp/probes/339/v1/logs"
mkdir -p "$LOGDIR"
OUT="$LOGDIR/v1-tripwire-sensitivity.log"
PRISTINE=$(git -C "$REPO" rev-parse "$ANCHOR:internal/risk/pathresolver_other.go")
h() { git -C "$REPO" hash-object "$1"; }
OLD='func resolveHandle(p string) (string, bool) { return "", false }'
NEW='func resolveHandle(p string) (string, bool) { return p, true } // MUT339V1'
{
  echo "# 339-v1 自打：把 pathresolver_other.go 的 resolveHandle stub 换成可解析 ⇒ POSIX 绊线必须当场点名红"
  echo "anchor=$ANCHOR tree=$TREE"
  echo "P0=$(h "$PR") expect=$PRISTINE"
  [ "$(h "$PR")" = "$PRISTINE" ] || { echo "GUARD-ABORT pristine"; exit 2; }
  echo "old_anchor_hits_before=$(grep -cF "$OLD" "$PR") (must be 1)"
  perl -0pi -e 's/func resolveHandle\(p string\) \(string, bool\) \{ return "", false \}/func resolveHandle(p string) (string, bool) { return p, true } \/\/ MUT339V1/' "$PR"
  echo "MUT=$(h "$PR")"
  echo "old_anchor_hits_after=$(grep -cF "$OLD" "$PR" || true) (must be 0)"
  echo "new_anchor_hits_after=$(grep -cF 'return p, true } \/\/ MUT339V1' "$PR" || grep -cF 'MUT339V1' "$PR") (must be 1)"
  MSYS_NO_PATHCONV=1 docker run --rm \
    -v 'D:/tmp/wisp/20261018/339v1/trip/tree':/src \
    -v 'D:/work/base/gopath/pkg/mod':/go/pkg/mod:ro -w /src \
    -e GOFLAGS=-mod=mod -e GOPROXY=off -e CGO_ENABLED=0 -e HOME=/src -e GOMODCACHE=/go/pkg/mod \
    golang:1.27 bash -c "mkdir -p /tmp/gt; export TMPDIR=/tmp/gt GOTMPDIR=/tmp/gt; go test ./internal/risk/ -run 'TestAnchorMissHasNoWorldOnThisPlatform339' -count=1 -timeout 30m -v 2>&1; echo RC=\$?"
  git -C "$REPO" show "$ANCHOR:internal/risk/pathresolver_other.go" > "$PR"
  echo "R=$(h "$PR") match=$([ "$(h "$PR")" = "$PRISTINE" ] && echo YES || echo NO)"
} > "$OUT" 2>&1
cat "$OUT"
