#!/usr/bin/env bash
# 334-v1 突变台件 v2（v1 的还原式写坏了：perl 还原锚用了 "\t}" 而文件里是列 0 的 "}"
# ⇒ M-V1 残留在突变树上，M-V2／M-V3 那两读是在脏态上取的，作废。v2 的还原＝逐发
# 「从具名提交重新取回字节覆盖盘上件 ＋ 断言哈希＝交件 blob」，⛔ 再依赖 perl 反写。
# 只打仓外导出树 D:/W334V1TreeMut 的盘上字节；⛔ 跟踪件、⛔ -overlay。
set -u
REPO="D:/work/workspace/projects plans/Wisp"
TREE="D:/W334V1TreeMut"
F="$TREE/internal/risk/blacklist.go"
PRISTINE=c1799cd24afce19e6accef23db23a7d7ee746dec   # git show 3a118226:internal/risk/blacklist.go
LOGDIR="$REPO/.scratch/wisp/probes/334/v1/logs"
ROSTER_LEG='TestAnchorMissReasonNamesCandidateAndMissingAnchor|TestAnchorMissKeepsNonOverridableClassAndNeedL2|TestFileBaseNameClaimChangesNoClassOrReason'
ROSTER_CHAIN='TestProbe334V1PerAnchorChain'

restore () { git -C "$REPO" show "3a118226:internal/risk/blacklist.go" > "$F"; }
check_hash () { local h; h=$(git -C "$REPO" hash-object "$F")
  echo "  hash=$h expected=$PRISTINE match=$([ "$h" = "$PRISTINE" ] && echo YES || echo NO)  [$1]"; }

run_roster () { # $1=名册 $2=日志
  local gt; gt="$(dirname "$TREE")/W334V1Gotmp/$2"; mkdir -p "$gt" || { echo "  mkdir_gotmp_rc=NONZERO"; return 1; }
  ( cd "$TREE" && GOPATH='D:\work\base\gopath' GOTMPDIR="$gt" \
      go test ./internal/risk/ -run "$1" -count=1 -timeout 30m -v ) > "$LOGDIR/$2" 2>&1
  local rc=$?
  echo "  go_test_rc=$rc RUN=$(grep -c '^=== RUN' "$LOGDIR/$2" || true) PASS_TOP=$(grep -c '^--- PASS' "$LOGDIR/$2" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$LOGDIR/$2" || true) SKIP=$(grep -c '^--- SKIP' "$LOGDIR/$2" || true)"
  grep -E '^\s*--- FAIL: ' "$LOGDIR/$2" | sed 's/([0-9.]*s)//;s/^/  RED /' || true
}

chain_shape () { grep -hoE 'V1SHAPE (chain|offchain):[^ ]* class=[A-Za-z]+ allow=[a-z]+ needL2=[a-z]+ rule=[^ ]*' "$1" || true; }

mutate () { # $1=perl -0pi 表达式  $2=发红句标签
  restore
  perl -0pi -e "$1" "$F"
  echo "  落点逐字（⇄ 交件 blob 的行差）:"
  diff <(git -C "$REPO" show "3a118226:internal/risk/blacklist.go") "$F" | sed 's/^/    /'
}

echo "### 0 起手：先把 v1 残留还原干净"
restore; check_hash "restored-from-commit"
run_roster "$ROSTER_LEG" 60-v2-baseline-leg.log
run_roster "$ROSTER_CHAIN" 60-v2-baseline-chain.log
chain_shape "$LOGDIR/60-v2-baseline-chain.log" > "$LOGDIR/60-v2-baseline-chain.shape"
echo "  baseline_chain_lines=$(wc -l < "$LOGDIR/60-v2-baseline-chain.shape")"

echo "### M-V1 只动判级支（aRelates 尾行＝回落支 -> false），reason 生成一字未动"
mutate 's/\n\treturn uncertainAnchorMiss\(cand, anchor\)\n/\n\treturn false\n/'
run_roster "$ROSTER_LEG" 61-v2-MV1-leg.log
run_roster "$ROSTER_CHAIN" 61-v2-MV1-chain.log
restore; check_hash after-M-V1

echo "### M-V2a 只摘 %LOCALAPPDATA%\Microsoft\Credentials 那一枚 case 的回落支"
mutate 's/case localAppDataS != "" && aUnder\(cand, anchorForms\(anchorPath\(localAppDataS, `microsoft`, `credentials`\)\)\):/case localAppDataS != "" \&\& anchorMatches(cand, anchorForms(anchorPath(localAppDataS, `microsoft`, `credentials`)), isUnder):/'
run_roster "$ROSTER_LEG" 62-v2-MV2a-leg.log
run_roster "$ROSTER_CHAIN" 62-v2-MV2a-chain.log
chain_shape "$LOGDIR/62-v2-MV2a-chain.log" > "$LOGDIR/62-v2-MV2a-chain.shape"
echo "  我这把尺看到的行差（baseline ⇄ M-V2a）:"
diff "$LOGDIR/60-v2-baseline-chain.shape" "$LOGDIR/62-v2-MV2a-chain.shape" | sed 's/^/    /' || true
restore; check_hash after-M-V2a

echo "### M-V2b 只摘 ~/.ssh 那一枚 case 的回落支"
mutate 's/case homeS != "" && aUnder\(cand, anchorForms\(anchorPath\(homeS, `\.ssh`\)\)\):/case homeS != "" \&\& anchorMatches(cand, anchorForms(anchorPath(homeS, `.ssh`)), isUnder):/'
run_roster "$ROSTER_LEG" 63-v2-MV2b-leg.log
run_roster "$ROSTER_CHAIN" 63-v2-MV2b-chain.log
chain_shape "$LOGDIR/63-v2-MV2b-chain.log" > "$LOGDIR/63-v2-MV2b-chain.shape"
echo "  我这把尺看到的行差（baseline ⇄ M-V2b）:"
diff "$LOGDIR/60-v2-baseline-chain.shape" "$LOGDIR/63-v2-MV2b-chain.shape" | sed 's/^/    /' || true
restore; check_hash after-M-V2b

echo "### M-V2c 只摘 ~/.kube/config 那一枚 case 的回落支"
mutate 's/case homeS != "" && aEq\(cand, anchorForms\(anchorPath\(homeS, `\.kube`, `config`\)\)\):/case homeS != "" \&\& anchorMatches(cand, anchorForms(anchorPath(homeS, `.kube`, `config`)), eqPath):/'
run_roster "$ROSTER_LEG" 64-v2-MV2c-leg.log
run_roster "$ROSTER_CHAIN" 64-v2-MV2c-chain.log
chain_shape "$LOGDIR/64-v2-MV2c-chain.log" > "$LOGDIR/64-v2-MV2c-chain.shape"
echo "  我这把尺看到的行差（baseline ⇄ M-V2c）:"
diff "$LOGDIR/60-v2-baseline-chain.shape" "$LOGDIR/64-v2-MV2c-chain.shape" | sed 's/^/    /' || true
restore; check_hash after-M-V2c

echo "### M-V3 反形（变严）：回落支去掉祖先关系界＝任何候选都算命中"
mutate 's/\treturn isUnder\(cand\.root, anchor\.root\) \|\| isUnder\(anchor\.root, cand\.root\)/\treturn true/'
run_roster "$ROSTER_LEG" 65-v2-MV3-leg.log
run_roster "$ROSTER_CHAIN" 65-v2-MV3-chain.log
chain_shape "$LOGDIR/65-v2-MV3-chain.log" > "$LOGDIR/65-v2-MV3-chain.shape"
echo "  我这把尺看到的行差（baseline ⇄ M-V3）:"
diff "$LOGDIR/60-v2-baseline-chain.shape" "$LOGDIR/65-v2-MV3-chain.shape" | sed 's/^/    /' || true
restore; check_hash after-M-V3

echo "### AC#3 成对反控：把 AC#3 那枚新夹具改回 t.TempDir() ⇒ 必须回到 24／5／1／0 rc=1"
restore
cp "$F" /tmp/v1-blacklist-copy.go
NF="$TREE/internal/risk/name_claim_two_arm_216_test.go"
git -C "$REPO" show "3a118226:internal/risk/name_claim_two_arm_216_test.go" > /tmp/v1-nameclaim-pristine.go
perl -0pi -e 's/dir := candidateOffAnchorChain334\(t\)/dir := t.TempDir()/' "$NF"
echo "  夹具反控落点行差:"; diff /tmp/v1-nameclaim-pristine.go "$NF" | sed 's/^/    /'
GTR="$(dirname "$TREE")/W334V1Gotmp/66-v2-ac3-revert"; mkdir -p "$GTR"
( cd "$TREE" && GOPATH='D:\work\base\gopath' GOTMPDIR="$GTR" go test ./internal/risk/ \
    -run 'TestAllowSideFunctionsCarryNoSubstringOperator|TestSubstringOperatorCounter216IsNotBlind|TestFileBaseNameClaimChangesNoClassOrReason|TestArgvZeroClaimChangesNoShellVerdict|TestHostClaimChangesNoNetworkVerdict|TestOverrideKeySpelledWithClaimDoesNotUnlock' \
    -count=1 -timeout 30m -v ) > "$LOGDIR/66-v2-ac3-revert.log" 2>&1
rc=$?
echo "  go_test_rc=$rc RUN=$(grep -c '^=== RUN' "$LOGDIR/66-v2-ac3-revert.log" || true) PASS_TOP=$(grep -c '^--- PASS' "$LOGDIR/66-v2-ac3-revert.log" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$LOGDIR/66-v2-ac3-revert.log" || true) SKIP=$(grep -c '^--- SKIP' "$LOGDIR/66-v2-ac3-revert.log" || true)"
grep -E '^\s*--- FAIL: ' "$LOGDIR/66-v2-ac3-revert.log" | sed 's/([0-9.]*s)//;s/^/  RED /' || true
cp /tmp/v1-nameclaim-pristine.go "$NF"
echo "  夹具还原 hash=$(git -C "$REPO" hash-object "$NF") 交件 blob=$(git -C "$REPO" rev-parse "3a118226:internal/risk/name_claim_two_arm_216_test.go")"
check_hash final-blacklist

echo "### 末态：突变树上三枚件哈希（必须＝交件 blob）"
for f in blacklist.go anchor_miss_reason_334_test.go name_claim_two_arm_216_test.go zz_probe_334v1_test.go; do
  echo "  $f tree=$(git -C "$REPO" hash-object "$TREE/internal/risk/$f") commit=$(git -C "$REPO" rev-parse "3a118226:internal/risk/$f" 2>/dev/null || echo '（探针非跟踪件）')"
done
echo "  共享工作树 internal/** cmd/** 是否被本腿动过：git status --porcelain --untracked-files=no → [$(git -C "$REPO" status --porcelain --untracked-files=no)]"
