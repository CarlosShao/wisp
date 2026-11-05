#!/bin/bash
# 339-r2 AC#2 去牙反控（两平台成对）：把 aRelates 的回落支整支摘掉（同 334-v1 的 M-V1 形）
# ⇒ 修好后的仪器必须在两平台都点名红；还原后必须回绿。
# 突变只打在仓外导出树的物理字节上（⛔ 共享工作树），逐发 git show 还原 + hash 自证。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
BASE="/d/tmp/wisp/20261105/339r2"
LOGD="$REPO/.scratch/wisp/probes/339/r2/logs"
RUN='TestAnchorMiss|TestAnchorFallbackTeeth'
PRISTINE=$(git -C "$REPO" rev-parse HEAD:internal/risk/blacklist.go)
echo "PRISTINE_BLACKLIST=$PRISTINE"

mutate() { # $1 = tree
  local bl="$1/internal/risk/blacklist.go"
  local hit
  hit=$(grep -cF "$(printf '\treturn uncertainAnchorMiss(cand, anchor)')" "$bl" || true)
  echo "  anchor_hits=$hit (must be 1)"
  [ "$hit" -eq 1 ] || { echo "  GUARD-ABORT"; return 2; }
  perl -0pi -e "s/\Q$(printf '\treturn uncertainAnchorMiss(cand, anchor)')\E/\treturn false \/\/ M-V1 339-r2 deteeth/" "$bl"
  echo "  mutated_blob=$(git -C "$REPO" hash-object "$bl")"
}
restore() { local bl="$1/internal/risk/blacklist.go"; git -C "$REPO" show HEAD:internal/risk/blacklist.go > "$bl"
  local r; r=$(git -C "$REPO" hash-object "$bl"); echo "  restored_blob=$r match=$([ "$r" = "$PRISTINE" ] && echo YES || echo NO)"; }

echo "=== LINUX ==="
rm -rf "$BASE/linux-mut"; cp -a "$BASE/linux-after" "$BASE/linux-mut"
mutate "$BASE/linux-mut"
MSYS_NO_PATHCONV=1 docker run --rm \
  -v "$BASE/linux-mut":/src -v "/d/work/base/gopath/pkg/mod":/go/pkg/mod:ro -w /src \
  -e GOFLAGS=-mod=mod -e GOPROXY=off -e CGO_ENABLED=0 -e HOME=/src -e GOMODCACHE=/go/pkg/mod \
  golang:1.27 bash -c "mkdir -p /tmp/gt; export TMPDIR=/tmp/gt GOTMPDIR=/tmp/gt; go test ./internal/risk/ -run '$RUN' -count=1 -timeout 30m -v > /tmp/gt/o 2>&1; echo GO_TEST_RC=\$?; cat /tmp/gt/o" \
  > "$LOGD/linux-deteeth.log" 2>&1
O="$LOGD/linux-deteeth.log"
echo "  $(grep -E '^GO_TEST_RC=' "$O") RUN=$(grep -c '^=== RUN' "$O" || true) PASS=$(grep -c '^--- PASS' "$O" || true) FAIL=$(grep -c '^--- FAIL' "$O" || true) SKIP=$(grep -c '^--- SKIP' "$O" || true)"
grep -E '^--- FAIL' "$O" || true
grep -nE 'fixture guard|class: got|Gate: got|unlocked by an override|off-chain arm' "$O" | head -20 || true
restore "$BASE/linux-mut"

echo "=== WINDOWS ==="
rm -rf "$BASE/win-mut"; cp -a "$BASE/win-after" "$BASE/win-mut"
mutate "$BASE/win-mut"
export TEMP='D:\tmp\wisp\20261105\339r2\gt' TMP="$TEMP" TMPDIR="$TEMP" GOTMPDIR="$TEMP"
export GOFLAGS=-mod=mod GOPROXY=off CGO_ENABLED=0 GOCACHE='D:\tmp\wisp\20261105\339r2\gocache'
( cd "$BASE/win-mut" && go test ./internal/risk/ -run "$RUN" -count=1 -timeout 30m -v > "$LOGD/win-deteeth.log" 2>&1; echo "  WIN_RC=$?" )
O="$LOGD/win-deteeth.log"
echo "  RUN=$(grep -c '^=== RUN' "$O" || true) PASS=$(grep -c '^--- PASS' "$O" || true) FAIL=$(grep -c '^--- FAIL' "$O" || true) SKIP=$(grep -c '^--- SKIP' "$O" || true)"
grep -E '^--- FAIL' "$O" || true
grep -nE 'fixture guard|class: got|Gate: got|unlocked by an override|off-chain arm' "$O" | head -20 || true
restore "$BASE/win-mut"
