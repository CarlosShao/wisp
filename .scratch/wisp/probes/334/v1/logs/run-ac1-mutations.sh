#!/usr/bin/env bash
# 334-v1·AC#1 那两枚成对突变的独立复跑（腿的 M-A／M-B 形，我自己伸手打）
# 加一发 M-V6：只把 reason 尾巴换成一枚假话（不提候选）⇒ 若 ⓐⓑ 不红，说明 reason 尺对"换一种假话"失明。
set -u
REPO="D:/work/workspace/projects plans/Wisp"
TREE="D:/W334V1TreeMut"; F="$TREE/internal/risk/blacklist.go"
PRISTINE=c1799cd24afce19e6accef23db23a7d7ee746dec
LOGDIR="$REPO/.scratch/wisp/probes/334/v1/logs"
ROSTER='TestAnchorMissReasonNamesCandidateAndMissingAnchor|TestAnchorMissKeepsNonOverridableClassAndNeedL2|TestFileBaseNameClaimChangesNoClassOrReason'
restore () { git -C "$REPO" show "3a118226:internal/risk/blacklist.go" > "$F"; }
check () { local h; h=$(git -C "$REPO" hash-object "$F"); echo "  hash=$h match=$([ "$h" = "$PRISTINE" ] && echo YES || echo NO) [$1]"; }
runr () { local gt; gt="$(dirname "$TREE")/W334V1Gotmp/$2"; mkdir -p "$gt"
  ( cd "$TREE" && GOPATH='D:\work\base\gopath' GOTMPDIR="$gt" go test ./internal/risk/ -run "$ROSTER" -count=1 -timeout 30m -v ) > "$LOGDIR/$2" 2>&1
  local rc=$?; echo "  go_test_rc=$rc RUN=$(grep -c '^=== RUN' "$LOGDIR/$2" || true) PASS_TOP=$(grep -c '^--- PASS' "$LOGDIR/$2" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$LOGDIR/$2" || true) SKIP=$(grep -c '^--- SKIP' "$LOGDIR/$2" || true)"
  grep -E '^\s*--- FAIL: ' "$LOGDIR/$2" | sed 's/([0-9.]*s)//;s/^/  RED /' || true; }

echo "### M-V4＝腿的 M-A 形：reason 拼装那一处改成永不产出说实话那句"
restore
perl -0pi -e 's/\tif anchorMatches\(cand, anchor, rel\) \|\| !uncertainAnchorMiss\(cand, anchor\) \{\n\t\treturn literal\n/\tif true {\n\t\treturn literal\n/' "$F"
diff <(git -C "$REPO" show "3a118226:internal/risk/blacklist.go") "$F" | sed 's/^/    /'
runr "$LOGDIR/84-mv4.log" 84-mv4.log; check after-M-V4

echo "### M-V5＝腿的 M-B 形：锚展开那一行换成丢掉 root 的裸形"
restore
perl -0pi -e 's/\tanchor := anchorForms\(anchorSpelling\)/\tanchor := pathForms{raw: normPath(anchorSpelling)}/' "$F"
diff <(git -C "$REPO" show "3a118226:internal/risk/blacklist.go") "$F" | sed 's/^/    /'
runr "$LOGDIR/85-mv5.log" 85-mv5.log; check after-M-V5

echo "### M-V6＝只把那句假话换成另一种假话（不提候选、不提缺位锚，仍带 fail-closed 字样）"
restore
perl -0pi -e 's/\treturn fmt\.Sprintf\("%s \[fail-closed, NOT a match on this rule\]: the candidate %q is not proven handle-expanded \(its own walk stopped at %q\) and the anchor %q is not installed here \(its expansion stopped at %q\), so R17 protects the ancestry chain the two share instead of proving a miss",\n\t\tliteral, cand\.raw, cand\.root, anchorSpelling, anchor\.root\)/\treturn fmt.Sprintf("%s [fail-closed, NOT a match on this rule]: this anchor is not installed here, so the shared ancestry chain is protected",\n\t\tliteral)/' "$F"
diff <(git -C "$REPO" show "3a118226:internal/risk/blacklist.go") "$F" | sed 's/^/    /'
runr "$LOGDIR/86-mv6.log" 86-mv6.log; check after-M-V6
restore; echo "### 末态"; for f in blacklist.go anchor_miss_reason_334_test.go name_claim_two_arm_216_test.go; do echo "  $f tree=$(git -C "$REPO" hash-object "$TREE/internal/risk/$f") commit=$(git -C "$REPO" rev-parse "3a118226:internal/risk/$f")"; done
echo "-- M-V1 那发的红句逐字（判级支单独动、reason 未动）"; grep -A4 'FAIL: TestAnchorMissKeepsNonOverridableClassAndNeedL2' "$LOGDIR/61-v2-MV1-leg.log" | head -14
