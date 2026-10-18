#!/bin/bash
# 339-r1 AC#2/AC(iv) 去牙反控 Linux 侧（甲形把牙齿搬到 windows 后，POSIX 侧只剩绊线）。
# 目的＝量死：去牙（摘 aRelates 回落支）在 Linux 能不能咬红我的仪器。⛔ 推理代替读数。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
TREE="/d/tmp/wisp/20261018/339r1/linux/after"
BL="$TREE/internal/risk/blacklist.go"
ANCHOR=e6fd0f5051812a23c72db444d12d76b4d27285a7
LOGD="$REPO/.scratch/wisp/probes/339/r1/logs"
mkdir -p "$LOGD"
PRISTINE=$(git -C "$REPO" rev-parse "$ANCHOR:internal/risk/blacklist.go")
gh() { git -C "$REPO" hash-object "$1"; }

# 去牙：在导出树物理字节上摘回落支（⛔ 跟踪文件），跑容器，然后还原。
hit=$(grep -cF "$(printf '\treturn uncertainAnchorMiss(cand, anchor)')" "$BL")
echo "anchor_hits=$hit (must be 1); pristine=$(gh "$BL") expect=$PRISTINE"
[ "$hit" -eq 1 ] || { echo GUARD-ABORT; exit 2; }
perl -0pi -e "s/\Q$(printf '\treturn uncertainAnchorMiss(cand, anchor)')\E/\treturn false \/\/ M-V1 ticket339 deteeth linux/" "$BL"
echo "mutated_hash=$(gh "$BL")"

MSYS_NO_PATHCONV=1 docker run --rm \
  -v "$TREE":/src -v "$(cd "$REPO" && go env GOMODCACHE)":/go/pkg/mod:ro -w /src \
  -e GOFLAGS=-mod=mod -e GOPROXY=off -e CGO_ENABLED=0 -e HOME=/src -e GOMODCACHE=/go/pkg/mod \
  golang:1.27 bash -c "mkdir -p /tmp/gt; export TMPDIR=/tmp/gt GOTMPDIR=/tmp/gt; go test ./internal/risk/ -run 'TestAnchorMissReason|TestAnchorMissKeeps|TestAnchorMissHasNoWorldOnThisPlatform339|TestAnchorFallbackTeeth' -count=1 -timeout 30m -v > /tmp/gt/o 2>&1; echo GO_TEST_RC=\$?; cat /tmp/gt/o" > "$LOGD/linux-ac2-deteeth.log" 2>&1

O="$LOGD/linux-ac2-deteeth.log"
echo "linux-ac2-deteeth: RUN=$(grep -c '^=== RUN' "$O"||true) PASS=$(grep -c '^--- PASS' "$O"||true) FAIL=$(grep -c '^--- FAIL' "$O"||true) SKIP=$(grep -c '^--- SKIP' "$O"||true)"
grep -E '^--- (FAIL|PASS)' "$O" || true
grep -E 'fault|no test files|build failed' "$O" || true

git -C "$REPO" show "$ANCHOR:internal/risk/blacklist.go" > "$BL"
R=$(gh "$BL"); echo "restored_hash=$R match=$( [ "$R" = "$PRISTINE" ] && echo YES || echo NO )"
