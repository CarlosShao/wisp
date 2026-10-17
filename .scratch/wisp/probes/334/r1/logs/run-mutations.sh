#!/usr/bin/env bash
# 334-r1 突变（牙）。⛔ 改跟踪件、⛔ -overlay（那形对读盘的尺失明）；只打仓外导出树的字节。
# 三发，落点两两不同，逐发当场还原并报还原哈希（trap 里再兜一次还原）：
#   M-A 落在 reason 拼装那处分支行（说实话那句永不产出）        ⇒ AC#1 两发必须红、AC#2 必须仍绿
#   M-B 落在锚展开那一行（anchorForms 换成不带 root 的裸形）    ⇒ 同上：证明那句确实读锚的展开
#   M-C 落在 uncertainAnchorMiss 首行（判级放松＝A 降回 B）      ⇒ AC#2 正控必须咬红（票面就是钉这件事）
# 落点一律内容锚（函数体内那枚 token），⛔ 绝对行号。
set -u
REPO="D:/work/workspace/projects plans/Wisp"
ROOT="${1:?导出树}"
LOGS="$REPO/.scratch/wisp/probes/334/r1/logs"
SRC_BLOB=3a118226   # 本腿交件那笔：还原一律从这里重新取字节
F="$ROOT/internal/risk/blacklist.go"
export GOPATH='D:\work\base\gopath'
ROSTER='TestAnchorMissReasonNamesCandidateAndMissingAnchor|TestAnchorMissKeepsNonOverridableClassAndNeedL2|TestFileBaseNameClaimChangesNoClassOrReason'
gotmp="$LOCALAPPDATA\\wisp-test-gotmp\\portable-334r1-mut-$(date +%H%M%S)"
mkdir -p "$(cygpath -u "${gotmp//\\//}")"; export GOTMPDIR="$gotmp"

restore() {
    git -C "$REPO" show "$SRC_BLOB:internal/risk/blacklist.go" > "$F"
    echo "restored_tree_hash=$(git -C "$REPO" hash-object "$F") committed_blob=$(git -C "$REPO" rev-parse "$SRC_BLOB:internal/risk/blacklist.go")"
}
trap 'echo "### trap 兜底还原"; restore' EXIT

run_case() {
    local tag="$1"
    ( cd "$ROOT" && go test ./internal/risk/ -count=1 -timeout 30m -v -run "$ROSTER" ) \
        >"$LOGS/mut-$tag.log" 2>"$LOGS/mut-$tag.err.log"
    local rc=$?
    {
        echo "MUT_TAG=$tag rc=$rc"
        echo "RUN=$(grep -c '^=== RUN' "$LOGS/mut-$tag.log" || true) PASS_TOP=$(grep -c '^--- PASS' "$LOGS/mut-$tag.log" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$LOGS/mut-$tag.log" || true) SKIP=$(grep -c '^--- SKIP' "$LOGS/mut-$tag.log" || true)"
        grep '^    --- FAIL\|^--- FAIL' "$LOGS/mut-$tag.log" || true
        head -4 "$LOGS/mut-$tag.err.log" || true
    } | tee "$LOGS/mut-$tag.summary.log"
}

echo "### 0 基线（未突变，应全绿）"; run_case baseline
echo "### 1 突变前哈希"; echo "pre_tree_hash=$(git -C "$REPO" hash-object "$F") committed_blob=$(git -C "$REPO" rev-parse "$SRC_BLOB:internal/risk/blacklist.go")"

echo "### 2 M-A：reason 拼装那一处（分支守卫）⇒ 永远退回今天那枚裸字面串"
sed -i 's#\tif anchorMatches(cand, anchor, rel) || !uncertainAnchorMiss(cand, anchor) {#\tif true { // MUTATION 334 M-A reason-assembly branch: always the bare literal#' "$F"
grep -n "MUTATION 334 M-A" "$F"; run_case M-A; restore

echo "### 3 M-B：锚展开那一行 ⇒ 换成不带 root 的裸形"
sed -i 's#anchor := anchorForms(anchorSpelling)#anchor := pathForms{raw: normPath(anchorSpelling)} // MUTATION 334 M-B anchor-expansion line#' "$F"
grep -n "MUTATION 334 M-B" "$F"; run_case M-B; restore

echo "### 4 M-C：判级放松（正控，必须把 AC#2 咬红）"
sed -i 's#^func uncertainAnchorMiss(cand, anchor pathForms) bool {#func uncertainAnchorMiss(cand, anchor pathForms) bool {\n\treturn false // MUTATION 334 M-C level relaxation#' "$F"
grep -n "MUTATION 334 M-C" "$F"; run_case M-C; restore

echo "### 5 收尾核对：三枚件全部等于交件字节"
for f in blacklist.go anchor_miss_reason_334_test.go name_claim_two_arm_216_test.go; do
    echo "$f tree=$(git -C "$REPO" hash-object "$ROOT/internal/risk/$f") blob=$(git -C "$REPO" rev-parse "$SRC_BLOB:internal/risk/$f")"
done
echo "### 6 共享工作树未被本腿写过（应空）"; git -C "$REPO" status --porcelain -- internal/risk
