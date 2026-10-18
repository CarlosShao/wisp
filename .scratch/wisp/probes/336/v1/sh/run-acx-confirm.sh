#!/usr/bin/env bash
# 336-v1 intake-anchor confirmation, re-run by the JUDGE (dispatch line 7 says the
# orchestrator already did this triple; I treat that as UNVERIFIED and redo it myself).
# Three serial legs, ONE tree, ONE combined roster (the four instruments on the books:
# ticket-334's two + ticket-216's one + ticket-336's teeth):
#   1) pristine bytes  -> expect rc=0 and every top-level green
#   2) MU-1h only (drop the :237 operand's fallback branch, content anchor) -> expect
#      rc=1 and the red named by the roster name, not just by count
#   3) restore + rerun -> expect rc=0 again
# Hash verified before and after every leg, in one shell (no cross-shell trap).
set -u
REPO="D:/work/workspace/projects plans/Wisp"
ANCHOR="${ANCHOR:?anchor}"
TREE="${TREE:?exported tree}"
LOGDIR="${LOGDIR:?log dir}"
GOTMPROOT="${GOTMPROOT:?gotmp root}"
F="$TREE/internal/risk/blacklist.go"
PRISTINE=$(git -C "$REPO" rev-parse "$ANCHOR:internal/risk/blacklist.go")
export GOPATH='D:\work\base\gopath'
export PATH="$TREE/third_party/sherpa-onnx:$TREE/build:$PATH"
ROSTER='TestAnchorMissReasonNamesCandidateAndMissingAnchor|TestAnchorMissKeepsNonOverridableClassAndNeedL2|TestFileBaseNameClaimChangesNoClassOrReason|TestAnchorFallbackTeethPerAnchor336'

restore () { git -C "$REPO" show "$ANCHOR:internal/risk/blacklist.go" > "$F"; }
check_hash () { local h; h=$(git -C "$REPO" hash-object "$F")
  echo "  hash=$h expected=$PRISTINE match=$([ "$h" = "$PRISTINE" ] && echo YES || echo NO)  [$1]"; }
run () { local gt="$GOTMPROOT/$1"; mkdir -p "$gt"
  ( cd "$TREE" && GOTMPDIR="$gt" go test ./internal/risk/ -run "$ROSTER" -count=1 -timeout 30m -v ) > "$LOGDIR/$1" 2>&1
  local rc=$?
  echo "  rc=$rc RUN=$(grep -c '^=== RUN' "$LOGDIR/$1" || true) PASS_TOP=$(grep -c '^--- PASS' "$LOGDIR/$1" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$LOGDIR/$1" || true) SKIP=$(grep -c '^--- SKIP' "$LOGDIR/$1" || true)"
  grep -E '^[[:space:]]*--- (FAIL|PASS): Test' "$LOGDIR/$1" | sed 's/([0-9.]*s)//' > "$LOGDIR/$1.tops" || true
  grep -E '^[[:space:]]*--- FAIL: ' "$LOGDIR/$1" | sed 's/([0-9.]*s)//;s/^/    RED /' || true
  printf '    GREEN_TOP=%s RED_TOP=%s (positive control: the same two greps on the paired pristine log)\n' \
    "$(grep -cE '^[[:space:]]*--- PASS: Test' "$LOGDIR/$1" || true)" "$(grep -cE '^[[:space:]]*--- FAIL: Test' "$LOGDIR/$1" || true)"
}
{
echo "### intake-anchor confirmation triple: anchor=$ANCHOR pristine=$PRISTINE tree=$TREE roster=4 instruments"
restore; check_hash start
echo "### leg 1 pristine:"; run acX-1-pristine.log
echo "### leg 2 MU-1h only:"
export N='localAppDataS != "" && aUnder(cand, anchorForms(anchorPath(localAppDataS, `microsoft`, `credentials`)))'
export R='localAppDataS != "" && anchorMatches(cand, anchorForms(anchorPath(localAppDataS, `microsoft`, `credentials`)), isUnder)'
echo "  needle_before=$(grep -cF "$N" "$F" || true)"
perl -0pi -e 'my $q=quotemeta($ENV{N}); s/$q/$ENV{R}/' "$F"
echo "  needle_after=$(grep -cF "$N" "$F" || true) branch_removals=$(grep -cF 'anchorMatches(cand, anchorForms(' "$F" || true) (want 0 / 1)"
echo "  mutated_hash=$(git -C "$REPO" hash-object "$F")"
run acX-2-MU1h.log
echo "  verbatim red sentence:"; grep -A12 -- '--- FAIL: TestAnchorFallbackTeethPerAnchor336/local-credentials' "$LOGDIR/acX-2-MU1h.log" | sed 's/^/    /' || true
echo "### leg 3 restore + rerun:"; restore; check_hash after-MU1h; run acX-3-restored.log
check_hash final
echo "  shared_worktree_dirty_internal_cmd=$(git -C "$REPO" status --porcelain --untracked-files=no -- internal cmd | wc -l)"
echo "acX_rc=0"
} 2>&1 | tee "$LOGDIR/acX-confirm-anchor.log"
