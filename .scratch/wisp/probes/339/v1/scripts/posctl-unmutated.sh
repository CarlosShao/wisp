#!/bin/bash
# 339-v1 台件⑧：去牙两发的★成对未突变正控（同一枚 -run 名册、同一棵树、pristine 哈希两侧各验一次）。
# ⛔ 空操作伪装成功＝落刀前后命中数都要量；这里量的是"未突变态"＝旧锚命中必须＝1、MUT 锚必须＝0。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
S="$REPO/.scratch/wisp/probes/339/v1"
ANCHOR=7437ef882fac5b427bf17c4bf590819c58d0b375
PRISTINE=$(git -C "$REPO" rev-parse "$ANCHOR:internal/risk/blacklist.go")
RUNPAT='TestAnchorMissReasonNamesCandidateAndMissingAnchor|TestAnchorMissKeepsNonOverridableClassAndNeedL2|TestAnchorMissHasNoWorldOnThisPlatform339|TestAnchorFallbackTeethPerAnchor336'
OLD='	return uncertainAnchorMiss(cand, anchor)'
h() { git -C "$REPO" hash-object "$1"; }
numbers() {
  local F="$1"
  echo "RUN=$(grep -c '^=== RUN' "$F" || true) PASS_TOP=$(grep -c '^--- PASS' "$F" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$F" || true) SKIP=$(grep -c '^--- SKIP' "$F" || true)"
  grep -E '^--- FAIL' "$F" || true
  grep -E '^RC=' "$F" || true
}

{
  echo "### 正控①LINUX 未突变 tree=dtl run=$RUNPAT"
  BL=/d/tmp/wisp/20261018/339v1/dtl/tree/internal/risk/blacklist.go
  echo "P0=$(h "$BL") expect=$PRISTINE old_anchor_hits=$(grep -cF "$OLD" "$BL" || true) mut_hits=$(grep -cF 'MUT339V1-A' "$BL" || true)"
  MSYS_NO_PATHCONV=1 docker run --rm -v 'D:/tmp/wisp/20261018/339v1/dtl/tree':/src \
    -v 'D:/work/base/gopath/pkg/mod':/go/pkg/mod:ro -w /src \
    -e GOFLAGS=-mod=mod -e GOPROXY=off -e CGO_ENABLED=0 -e HOME=/src -e GOMODCACHE=/go/pkg/mod \
    golang:1.27 bash -c "mkdir -p /tmp/gt; export TMPDIR=/tmp/gt GOTMPDIR=/tmp/gt; go test ./internal/risk/ -run '$RUNPAT' -count=1 -timeout 30m -v 2>&1; echo RC=\$?"
} > "$S/logs/v1-posctl-A-linux-unmutated.log" 2>&1
echo "== LINUX 正控 =="; numbers "$S/logs/v1-posctl-A-linux-unmutated.log"

{
  echo "### 正控②WINDOWS 未突变 tree=dtb run=$RUNPAT"
  BL=/d/tmp/wisp/20261018/339v1/dtb/tree/internal/risk/blacklist.go
  echo "P0=$(h "$BL") expect=$PRISTINE old_anchor_hits=$(grep -cF "$OLD" "$BL" || true) mutA_hits=$(grep -cF 'MUT339V1-A' "$BL" || true) mutB_hits=$(grep -cF 'MUT339V1-B' "$BL" || true)"
  cd /d/tmp/wisp/20261018/339v1/dtb/tree || exit 9
  mkdir -p /d/tmp/wisp/20261018/339v1/gt-win
  TMPDIR=/d/tmp/wisp/20261018/339v1/gt-win GOTMPDIR=/d/tmp/wisp/20261018/339v1/gt-win \
  GOFLAGS=-mod=mod GOPROXY=off CGO_ENABLED=0 \
  go test ./internal/risk/ -run "$RUNPAT" -count=1 -timeout 30m -v 2>&1; echo "RC=$?"
} > "$S/logs/v1-posctl-B-windows-unmutated.log" 2>&1
echo "== WINDOWS 正控（落点B那棵树，未突变态）=="
numbers "$S/logs/v1-posctl-B-windows-unmutated.log"
