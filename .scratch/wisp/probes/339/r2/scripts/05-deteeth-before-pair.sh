#!/bin/bash
# 339-r2 AC#2 对照：同一发去牙打在改前树（HEAD 原样，票 334 仪器仍带 //go:build windows）
# 上 ⇒ Linux 侧那枚仪器根本不跑（名册里没有它的名字），红只剩 336 牙与 339 前提孪生。
# 与 04 那一发（改后树）作差＝"补世界"到底补回了哪几枚名字。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
BASE="/d/tmp/wisp/20261105/339r2"
LOGD="$REPO/.scratch/wisp/probes/339/r2/logs"
RUN='TestAnchorMiss|TestAnchorFallbackTeeth'

rm -rf "$BASE/linux-before-mut"; cp -a "$BASE/linux-before" "$BASE/linux-before-mut"
BL="$BASE/linux-before-mut/internal/risk/blacklist.go"
hit=$(grep -cF "$(printf '\treturn uncertainAnchorMiss(cand, anchor)')" "$BL" || true)
echo "anchor_hits=$hit (must be 1)"
[ "$hit" -eq 1 ] || { echo GUARD-ABORT; exit 2; }
perl -0pi -e "s/\Q$(printf '\treturn uncertainAnchorMiss(cand, anchor)')\E/\treturn false \/\/ M-V1 339-r2 deteeth/" "$BL"

MSYS_NO_PATHCONV=1 docker run --rm \
  -v "$BASE/linux-before-mut":/src -v "/d/work/base/gopath/pkg/mod":/go/pkg/mod:ro -w /src \
  -e GOFLAGS=-mod=mod -e GOPROXY=off -e CGO_ENABLED=0 -e HOME=/src -e GOMODCACHE=/go/pkg/mod \
  golang:1.27 bash -c "mkdir -p /tmp/gt; export TMPDIR=/tmp/gt GOTMPDIR=/tmp/gt; go test ./internal/risk/ -run '$RUN' -count=1 -timeout 30m -v > /tmp/gt/o 2>&1; echo GO_TEST_RC=\$?; cat /tmp/gt/o" \
  > "$LOGD/linux-before-deteeth.log" 2>&1
O="$LOGD/linux-before-deteeth.log"
echo "$(grep -E '^GO_TEST_RC=' "$O") RUN=$(grep -c '^=== RUN' "$O" || true) PASS=$(grep -c '^--- PASS' "$O" || true) FAIL=$(grep -c '^--- FAIL' "$O" || true) SKIP=$(grep -c '^--- SKIP' "$O" || true)"
grep -E '^--- (FAIL|PASS)' "$O" || true

git -C "$REPO" show HEAD:internal/risk/blacklist.go > "$BL"
R=$(git -C "$REPO" hash-object "$BL"); P=$(git -C "$REPO" rev-parse HEAD:internal/risk/blacklist.go)
echo "restored=$R match=$([ "$R" = "$P" ] && echo YES || echo NO)"

echo "=== roster delta (before-mut ⇄ after-mut) ==="
echo "---before reds:"; grep -E '^--- FAIL' "$LOGD/linux-before-deteeth.log" | sed 's/ (.*//' || true
echo "---after  reds:"; grep -E '^--- FAIL' "$LOGD/linux-deteeth.log" | sed 's/ (.*//' || true
