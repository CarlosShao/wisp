#!/usr/bin/env bash
# 334-v1 幻影凭据普查（AC#3-c 那格，本仓已入账七次的高危形）。
# 射程＝【已入库字节】：具名提交树 3a118226 的 internal/risk/** 注释 ＋ 腿交回的四枚 .md。
# 判据＝凡"指名一枚测试／件当凭据"的写法，逐枚验：①名字真在册（go test 名册尺）②件真在树里。
set -u
REPO="D:/work/workspace/projects plans/Wisp"
ANCHOR=3a118226
LOG="$REPO/.scratch/wisp/probes/334/v1/logs/50-phantom-census.log"
: > "$LOG"
exec >> "$LOG" 2>&1
echo "date=$(date '+%F %T %z') anchor=$ANCHOR"

# 名册尺①：具名树里所有 Test* 函数名（含子测试名另算）
git -C "$REPO" grep -h -oE '^func (Test[A-Za-z0-9_]+)' "$ANCHOR" -- '*_test.go' | sed 's/^func //' | sort -u > /tmp/v1-testnames.txt
echo "roster_testfuncs=$(wc -l < /tmp/v1-testnames.txt) (尺＝git grep '^func Test' 射程＝$ANCHOR 全仓 *_test.go)"
# 名册尺②：具名树里所有跟踪文件路径
git -C "$REPO" ls-tree -r --name-only "$ANCHOR" > /tmp/v1-files.txt
echo "roster_files=$(wc -l < /tmp/v1-files.txt)"

echo; echo "## 正控（先证明这把尺看得见真东西）"
for nm in TestFileBaseNameClaimChangesNoClassOrReason TestAnchorMissKeepsNonOverridableClassAndNeedL2; do
  echo "  正控 $nm in 名册? $(grep -cx "$nm" /tmp/v1-testnames.txt || true)"
done
echo "  正控 internal/risk/blacklist.go in 树? $(grep -cxF 'internal/risk/blacklist.go' /tmp/v1-files.txt || true)"

echo; echo "## 待验对象＝注释与读数件里指名的测试名／件名"
# 从三枚产码件（取注释行）＋四枚读数件里抓 Test[A-Z]… 形与 xxx_test.go／internal/**.go 形
{ git -C "$REPO" show "$ANCHOR:internal/risk/blacklist.go";
  git -C "$REPO" show "$ANCHOR:internal/risk/anchor_miss_reason_334_test.go";
  git -C "$REPO" show "$ANCHOR:internal/risk/name_claim_two_arm_216_test.go";
  for f in ac0-ac4-readings.md ac1-ac2-ac3-ac5-readings.md selfcheck-answers.md incident-swallowed-index-entry.md; do
    git -C "$REPO" show "$ANCHOR:.scratch/wisp/probes/334/r1/$f" 2>/dev/null || cat "$REPO/.scratch/wisp/probes/334/r1/$f"
  done
} > /tmp/v1-scan.txt
echo "scan_bytes=$(wc -c < /tmp/v1-scan.txt)"
grep -oE 'Test[A-Za-z0-9_]+|[A-Za-z0-9_./-]+\.go' /tmp/v1-scan.txt | sort -u > /tmp/v1-refs.txt
echo "refs_found=$(wc -l < /tmp/v1-refs.txt)"

echo; echo "### A. 指名的测试名，逐枚验在册（NO＝幻影）"
grep -E '^Test' /tmp/v1-refs.txt | while read -r nm; do
  hit=$(grep -cx "$nm" /tmp/v1-testnames.txt || true)
  echo "  test:$nm 在册=$hit"
done

echo; echo "### B. 指名的 .go 件，逐枚验在树（NO＝幻影）"
grep -E '\.go$' /tmp/v1-refs.txt | while read -r f; do
  # 只按 basename 找（读数里常省包名），报命中路径
  hit=$(grep -E "(^|/)$f\$" /tmp/v1-files.txt | head -3 | tr '\n' ' ' || true)
  echo "  file:$f 命中=${hit:-NONE}"
done

echo; echo "### C. 编排者点名的那枚中途报告引用（blacklist_test.go ＋ 那个用例名）"
echo "  internal/risk/blacklist_test.go 在 $ANCHOR 树里? $(git -C "$REPO" cat-file -e "$ANCHOR:internal/risk/blacklist_test.go" 2>/dev/null && echo YES || echo NO)"
echo "  用例名 TestApprovalSensitivePathIsNotOverridable 全仓命中枚数=$(git -C "$REPO" grep -c 'TestApprovalSensitivePathIsNotOverridable' "$ANCHOR" -- . 2>/dev/null | wc -l || true)"
echo "  同一枚名在 HEAD 树里命中=$(git -C "$REPO" log --all --oneline -S 'TestApprovalSensitivePathIsNotOverridable' | wc -l || true) 笔提交动过该串"
echo "  四枚读数件里是否也写了这枚名? $(grep -c 'TestApprovalSensitivePathIsNotOverridable' /tmp/v1-scan.txt || true) 处"
echo "  'non-overridable' 这个词在 $ANCHOR 的命中文件: $(git -C "$REPO" grep -l 'non-overridable' "$ANCHOR" -- 'internal/**' | tr '\n' ' ')"
echo "census_rc=0"
