#!/usr/bin/env bash
# 334-v1·AC#3 成对反控重跑（v2 台件那一发我把 GOTMPDIR 放在了 D: 卷＝台面缺陷，不是读数）：
# 反控必须与腿的改前红同一把台面——CI 形 GOTMPDIR＝$LOCALAPPDATA\\wisp-test-gotmp\\portable-…，
# 因为那一发红本身就活在"临时根落在真主目录链内"这件事上，换根＝换判据。
# 随后：把 v2 探针（加了 solo 世界）在 T_BEFORE ⇄ T_AFTER 两棵树上重跑一次三值对撞。
set -u
REPO="D:/work/workspace/projects plans/Wisp"
LOGDIR="$REPO/.scratch/wisp/probes/334/v1/logs"
TREE_MUT="D:\\W334V1TreeMut"
NF="/d/W334V1TreeMut/internal/risk/name_claim_two_arm_216_test.go"
ROSTER6='TestAllowSideFunctionsCarryNoSubstringOperator|TestSubstringOperatorCounter216IsNotBlind|TestFileBaseNameClaimChangesNoClassOrReason|TestArgvZeroClaimChangesNoShellVerdict|TestHostClaimChangesNoNetworkVerdict|TestOverrideKeySpelledWithClaimDoesNotUnlock'

echo "### 1 AC#3 反控（CI 形台面）：交件字节＋夹具改回 t.TempDir()"
git -C "$REPO" show "3a118226:internal/risk/name_claim_two_arm_216_test.go" > "$NF"
perl -0pi -e 's/dir := candidateOffAnchorChain334\(t\)/dir := t.TempDir()/' "$NF"
echo "  落点行差（⇄ 交件 blob）:"
diff <(git -C "$REPO" show "3a118226:internal/risk/name_claim_two_arm_216_test.go") "$NF" | sed 's/^/    /'
sh "$REPO/.scratch/wisp/probes/334/v1/logs/run-risk.sh" "$TREE_MUT" "$ROSTER6" 70-ac3-revert-cishape.log --ci-shape > "$LOGDIR/70-ac3-revert-cishape.runner.log" 2>&1
grep -hE "FOURCOUNT|go_test_rc|blacklist_blob|redlist=" "$LOGDIR/70-ac3-revert-cishape.runner.log"
sed 's/^/  RED /' "$LOGDIR/70-ac3-revert-cishape.log.redlist"
# 还原夹具并报哈希
git -C "$REPO" show "3a118226:internal/risk/name_claim_two_arm_216_test.go" > "$NF"
echo "  夹具还原 hash=$(git -C "$REPO" hash-object "$NF") 交件 blob=$(git -C "$REPO" rev-parse "3a118226:internal/risk/name_claim_two_arm_216_test.go")"

echo "### 2 交件字节在同一把 CI 形台面上的正控（夹具未动 ⇒ 必须 24/6/0/0 rc=0）"
sh "$REPO/.scratch/wisp/probes/334/v1/logs/run-risk.sh" "$TREE_MUT" "$ROSTER6" 71-ac3-pristine-cishape.log --ci-shape > "$LOGDIR/71-ac3-pristine-cishape.runner.log" 2>&1
grep -hE "FOURCOUNT|go_test_rc|redlist=" "$LOGDIR/71-ac3-pristine-cishape.runner.log"

echo "### 3 v2 探针的三值对撞重跑（T_BEFORE ⇄ T_AFTER，solo 世界已含）"
for pair in before:before after:after; do :; done
sh "$REPO/.scratch/wisp/probes/334/v1/logs/run-risk.sh" 'D:\W334V1TreeBefore' 'TestProbe334V1' 72-probe-v2-before.log --ci-shape > "$LOGDIR/72-probe-v2-before.runner.log" 2>&1
sh "$REPO/.scratch/wisp/probes/334/v1/logs/run-risk.sh" 'D:\W334V1TreeAfter'  'TestProbe334V1' 73-probe-v2-after.log  --ci-shape > "$LOGDIR/73-probe-v2-after.runner.log" 2>&1
grep -hE "FOURCOUNT|blacklist_blob" "$LOGDIR/72-probe-v2-before.runner.log" "$LOGDIR/73-probe-v2-after.runner.log"
for t in before after; do
  grep -hoE 'V1(SHAPE|OVERRIDE) .*' "$LOGDIR/7$([ $t = before ] && echo 2 || echo 3)-probe-v2-$t.log" | sed -E 's/(portable-334v1-)[0-9-]+/[TS]/' > "$LOGDIR/v2-$t.shape"
done
wc -l "$LOGDIR/v2-before.shape" "$LOGDIR/v2-after.shape"
diff "$LOGDIR/v2-before.shape" "$LOGDIR/v2-after.shape" > "$LOGDIR/v2-shape.diff" 2>&1
echo "  THREE_VALUE_DIFF_lines=$(wc -l < "$LOGDIR/v2-shape.diff") rc=$?  (0＝逐字相同)"
cat "$LOGDIR/v2-shape.diff"
