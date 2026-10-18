#!/bin/bash
# 339-r1 绊线敏感尺：证明 anchor_miss_reason_339_other_test.go 的绊线不是装饰——
# 把 pathresolver_other.go 的 resolveHandle stub 换成"能解析"（导出树物理字节，⛔ 跟踪件），
# 我的 POSIX 绊线必须当场点名红（这正是牙齿要搬来的时刻）。发后当场还原＋哈希自证。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
ANCHOR=7437ef882fac5b427bf17c4bf590819c58d0b375   # 我交码笔＝交付字节
TREE="/d/tmp/wisp/20261018/339r1/linux/tripwire/tree"
PR="$TREE/internal/risk/pathresolver_other.go"
LOGD="$REPO/.scratch/wisp/probes/339/r1/logs"
mkdir -p "$TREE" "$LOGD"
git -C "$REPO" archive "$ANCHOR" go.mod go.sum internal cmd tools scripts frontend models docs | tar -xf - -C "$TREE"
PRISTINE=$(git -C "$REPO" rev-parse "$ANCHOR:internal/risk/pathresolver_other.go")
gh() { git -C "$REPO" hash-object "$1"; }
echo "P0=$(gh "$PR") expect=$PRISTINE"
hit=$(grep -cF 'func resolveHandle(p string) (string, bool) { return "", false }' "$PR")
echo "anchor_hits=$hit (must be 1)"
[ "$hit" -eq 1 ] || { echo GUARD-ABORT > "$LOGD/tripwire-sensitivity.log"; exit 2; }
perl -0pi -e 's/func resolveHandle\(p string\) \(string, bool\) \{ return "", false \}/func resolveHandle(p string) (string, bool) { return p, true } \/\/ MUT: pretend resolver/' "$PR"
echo "MUT=$(gh "$PR")"
MSYS_NO_PATHCONV=1 docker run --rm -v "$TREE":/src -v "$(cd "$REPO" && go env GOMODCACHE)":/go/pkg/mod:ro -w /src \
  -e GOFLAGS=-mod=mod -e GOPROXY=off -e CGO_ENABLED=0 -e HOME=/src -e GOMODCACHE=/go/pkg/mod \
  golang:1.27 bash -c "mkdir -p /tmp/gt; export TMPDIR=/tmp/gt GOTMPDIR=/tmp/gt; go test ./internal/risk/ -run 'TestAnchorMissHasNoWorldOnThisPlatform339' -count=1 -timeout 30m -v 2>&1; echo RC=\$?" > "$LOGD/tripwire-sensitivity.log" 2>&1
git -C "$REPO" show "$ANCHOR:internal/risk/pathresolver_other.go" > "$PR"
echo "R=$(gh "$PR") match=$( [ "$(gh "$PR")" = "$PRISTINE" ] && echo YES || echo NO )" >> "$LOGD/tripwire-sensitivity.log"
cat "$LOGD/tripwire-sensitivity.log"
