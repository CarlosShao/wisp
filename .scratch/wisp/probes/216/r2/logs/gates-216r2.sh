#!/usr/bin/env bash
# 票 216-r2 门禁四数（只在共享工作树上跑；⛔ 部分导出树）
# 逐枚落 rc=，并把 ban #8 的分母两枚（改前/改后）另跑在完整导出树上。
set -u
REPO="D:/work/workspace/projects plans/Wisp"
OUT="$(cd "$(dirname "$0")" && pwd)"
cd "$REPO" || exit 1
export GOPATH='D:\work\base\gopath'
export GOTMPDIR='D:\tmp\leg216r2gotmp'
MINE="internal/risk/allow_substring_pin_216_test.go internal/risk/name_claim_two_arm_216_test.go"

echo "########## 1 gofumpt.exe -l（真身 /d/work/base/gopath/bin/gofumpt.exe）"
/d/work/base/gopath/bin/gofumpt.exe -l $MINE; echo "rc=$? (本腿两枚文件：空列表＝干净)"
echo "--- 同尺扫整包（一枚既有件在册，非本腿所改）---"
/d/work/base/gopath/bin/gofumpt.exe -l internal/risk/; echo "rc=$?"
echo "--- 那枚既有件归因：锚上的 blob 用同一把尺＝干净（⇒ 工作树 CRLF 检出形，不是源码缺陷）"
PROBE=/d/tmp/leg216r2gotmp/provenance-at-anchor.go
git show 72b65c61:internal/risk/provenance.go > "$PROBE"
/d/work/base/gopath/bin/gofumpt.exe -l "$PROBE"; echo "rc=$? (空＝锚上那版干净)"
file internal/risk/provenance.go "$PROBE"
git status --porcelain -- internal/risk/provenance.go | sed 's/^/worktree-status: /'; echo "rc=$?"

echo "########## 2 sh scripts/d22scan.sh（含正控 runtests.sh）"
sh scripts/d22scan.sh; echo "rc=$?"

echo "########## 3 sh scripts/check-path-length-budget.sh"
sh scripts/check-path-length-budget.sh; echo "rc=$?"

echo "########## 4 go vet ./internal/risk/"
go vet ./internal/risk/; echo "rc=$?"

echo "########## done"
