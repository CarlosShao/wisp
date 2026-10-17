#!/usr/bin/env bash
# 334-v1 突变台件：只打在【仓外导出台树】的盘上字节（⛔ 跟踪件、⛔ -overlay）。
# 每一发串行：先断言 pristine 哈希 → 打 → 跑名册 → 当场还原 → 再断言哈希回到交件 blob。
# 落点全用内容锚（函数体内那枚 token），⛔ 绝对行号。
set -u
REPO="D:/work/workspace/projects plans/Wisp"
TREE="D:/W334V1TreeMut"
F="$TREE/internal/risk/blacklist.go"
PRISTINE=c1799cd24afce19e6accef23db23a7d7ee746dec   # ＝ git show 3a118226:internal/risk/blacklist.go 的 blob
LOGDIR="$REPO/.scratch/wisp/probes/334/v1/logs"
ROSTER_LEG='TestAnchorMissReasonNamesCandidateAndMissingAnchor|TestAnchorMissKeepsNonOverridableClassAndNeedL2|TestFileBaseNameClaimChangesNoClassOrReason'

run_roster () { # $1=名册正则 $2=日志名
  local gt; gt="$(dirname "$TREE")/W334V1Gotmp/$2"
  mkdir -p "$gt" || { echo "  mkdir_gotmp_rc=NONZERO $gt"; return 1; }
  local t0; t0=$(date '+%F %T %z')
  ( cd "$TREE" && GOPATH='D:\work\base\gopath' \
      GOTMPDIR="$gt" \
      go test ./internal/risk/ -run "$1" -count=1 -timeout 30m -v ) > "$LOGDIR/$2" 2>&1
  local rc=$?
  echo "  [$t0 → $(date '+%T')] go_test_rc=$rc RUN=$(grep -c '^=== RUN' "$LOGDIR/$2" || true) PASS_TOP=$(grep -c '^--- PASS' "$LOGDIR/$2" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$LOGDIR/$2" || true) SKIP=$(grep -c '^--- SKIP' "$LOGDIR/$2" || true)"
  grep -E '^\s+--- FAIL: |^--- FAIL: ' "$LOGDIR/$2" | sed 's/([0-9.]*s)//' | sed 's/^/  RED /' || true
}

check_hash () { # $1=期望标签
  local h; h=$(git -C "$REPO" hash-object "$F")
  echo "  hash_now=$h expected=$PRISTINE match=$([ "$h" = "$PRISTINE" ] && echo YES || echo NO) ($1)"
}

echo "### 0 pristine"
check_hash before-any-mutation
run_roster "$ROSTER_LEG" 30-mut-baseline.log
run_roster 'TestProbe334V1PerAnchorChain' 30-probe-baseline.log
grep -hoE 'V1SHAPE chain:.*' "$LOGDIR/30-probe-baseline.log" > "$LOGDIR/30-probe-baseline.shape" || true

echo "### M-V1 只动判级支（aRelates 尾行回落支 -> false），reason 生成一字未动"
grep -c '^\treturn uncertainAnchorMiss(cand, anchor)$' "$F"
perl -0pi -e 's/\n\treturn uncertainAnchorMiss\(cand, anchor\)\n/\n\treturn false\n/' "$F"
grep -c '^\treturn false$' "$F"
run_roster "$ROSTER_LEG" 31-mutMV1-legroster.log
run_roster 'TestProbe334V1PerAnchorChain' 31-mutMV1-chain.log
perl -0pi -e 's/\n\treturn false\n\t\}\n\n\/\/ anchorMatches/\n\treturn uncertainAnchorMiss(cand, anchor)\n\t}\n\n\/\/ anchorMatches/' "$F"
check_hash after-M-V1-restore

echo "### M-V2 只把 ~/.ssh 那一枚 case 的判级换成不带回落支的 anchorMatches"
perl -0pi -e 's/case homeS != "" && aUnder\(cand, anchorForms\(anchorPath\(homeS, `\.ssh`\)\)\):/case homeS != "" \&\& anchorMatches(cand, anchorForms(anchorPath(homeS, `.ssh`)), isUnder):/' "$F"
grep -c 'anchorMatches(cand, anchorForms(anchorPath(homeS, `\.ssh`)), isUnder)' "$F"
run_roster "$ROSTER_LEG" 32-mutMV2-legroster.log
run_roster 'TestProbe334V1' 32-mutMV2-chain.log
grep -hoE 'V1SHAPE chain:.*|V1SHAPE offchain:.*' "$LOGDIR/32-mutMV2-chain.log" > "$LOGDIR/32-mutMV2-chain.shape" || true
diff "$LOGDIR/30-probe-baseline.shape" <(grep -E 'chain:' "$LOGDIR/32-mutMV2-chain.shape") > "$LOGDIR/32-mutMV2-chain.diff" || true
echo "  chain_diff_lines=$(wc -l < "$LOGDIR/32-mutMV2-chain.diff")"
cat "$LOGDIR/32-mutMV2-chain.diff"
perl -0pi -e 's/case homeS != "" && anchorMatches\(cand, anchorForms\(anchorPath\(homeS, `\.ssh`\)\), isUnder\):/case homeS != "" \&\& aUnder(cand, anchorForms(anchorPath(homeS, `.ssh`))):/' "$F"
check_hash after-M-V2-restore

echo "### M-V3 反向（变严）：回落支去掉祖先关系界＝任何候选都算命中"
perl -0pi -e 's/\treturn isUnder\(cand\.root, anchor\.root\) \|\| isUnder\(anchor\.root, cand\.root\)/\treturn true/' "$F"
run_roster "$ROSTER_LEG" 33-mutMV3-legroster.log
run_roster 'TestProbe334V1PerAnchorChain' 33-mutMV3-chain.log
grep -hoE 'V1SHAPE chain:.*|V1SHAPE offchain:.*' "$LOGDIR/33-mutMV3-chain.log" > "$LOGDIR/33-mutMV3-chain.shape" || true
diff "$LOGDIR/30-probe-baseline.shape" "$LOGDIR/33-mutMV3-chain.shape" > "$LOGDIR/33-mutMV3-chain.diff" || true
echo "  chain_diff_lines=$(wc -l < "$LOGDIR/33-mutMV3-chain.diff")"
grep -E '^[<>]' "$LOGDIR/33-mutMV3-chain.diff" | head -20
perl -0pi -e 's/\treturn true\n\}\n\n\/\/ hasGitConfigSegment/\treturn isUnder(cand.root, anchor.root) || isUnder(anchor.root, cand.root)\n}\n\n\/\/ hasGitConfigSegment/' "$F"
check_hash after-M-V3-restore-final

echo "### 4 收尾：交件树三枚件哈希（⛔ 任何一发残留在突变树上）"
for f in blacklist.go anchor_miss_reason_334_test.go name_claim_two_arm_216_test.go; do
  echo "  $f tree=$(git -C "$REPO" hash-object "$TREE/internal/risk/$f") commit=$(git -C "$REPO" rev-parse "3a118226:internal/risk/$f")"
done
