#!/usr/bin/env bash
# 336-v1 AC#0 own replay: the ticket's four starting measurements (M-V2a/b/c + M-V3),
# re-run by the JUDGE on its own exported tree, own GOTMPDIR, own scripts.
# Mutations hit ONLY the physical bytes of the exported tree (no tracked file edits,
# no -overlay). Content anchors (a token inside the case expression), NOT line numbers.
# Pristine hash verified before and after every leg, in ONE shell, serially.
# Needle/replacement travel through single-quoted env vars: the operand text contains
# backticks, which a double-quoted bash string would run as command substitution.
set -u
REPO="D:/work/workspace/projects plans/Wisp"
ANCHOR="${ANCHOR:?named commit anchor}"
TREE="${TREE:?exported tree root}"
LOGDIR="${LOGDIR:?log dir}"
GOTMPROOT="${GOTMPROOT:?D-volume gotmp root}"
F="$TREE/internal/risk/blacklist.go"
PRISTINE=$(git -C "$REPO" rev-parse "$ANCHOR:internal/risk/blacklist.go")
export GOPATH='D:\work\base\gopath'
export PATH="$TREE/third_party/sherpa-onnx:$TREE/build:$PATH"
# The three instruments that were on the books when 334-v1 read 10/3/0/0 (the roster
# ticket 336 AC#0 names). This leg does not touch them.
ROSTER_LEG='TestAnchorMissReasonNamesCandidateAndMissingAnchor|TestAnchorMissKeepsNonOverridableClassAndNeedL2|TestFileBaseNameClaimChangesNoClassOrReason'

restore () { git -C "$REPO" show "$ANCHOR:internal/risk/blacklist.go" > "$F"; }
check_hash () {
  local h; h=$(git -C "$REPO" hash-object "$F")
  echo "  hash=$h expected=$PRISTINE match=$([ "$h" = "$PRISTINE" ] && echo YES || echo NO)  [$1]"
}

run_leg () { # $1 = log name
  local gt="$GOTMPROOT/$1"; mkdir -p "$gt" || { echo "  mkdir_gotmp_rc=NONZERO"; return 1; }
  ( cd "$TREE" && GOTMPDIR="$gt" go test ./internal/risk/ -run "$ROSTER_LEG" \
      -count=1 -timeout 30m -v ) > "$LOGDIR/$1" 2>&1
  local rc=$?
  echo "  go_test_rc=$rc RUN=$(grep -c '^=== RUN' "$LOGDIR/$1" || true) PASS_TOP=$(grep -c '^--- PASS' "$LOGDIR/$1" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$LOGDIR/$1" || true) SKIP=$(grep -c '^--- SKIP' "$LOGDIR/$1" || true)"
  grep -E '^\s*--- FAIL: ' "$LOGDIR/$1" | sed 's/([0-9.]*s)//;s/^/  RED /' || true
  grep -E '^\s*--- PASS: ' "$LOGDIR/$1" | sed 's/([0-9.]*s)//;s/^/  GREEN /' || true
}

mutate () { # $1 = needle (fixed string, content anchor) $2 = replacement
  restore
  export N="$1" R="$2"
  local before after
  before=$(grep -c -F "$1" "$F" || true)
  perl -0pi -e 'my $q=quotemeta($ENV{N}); s/$q/$ENV{R}/' "$F"
  after=$(grep -c -F "$1" "$F" || true)
  echo "  落点命中枚数 before=$before after=$after  (want 1 -> 0)"
  echo "  changed_lines=$(diff <(git -C "$REPO" show "$ANCHOR:internal/risk/blacklist.go") "$F" | grep -c '^[<>]' || true)"
  echo "  diff vs named commit blob:"
  diff <(git -C "$REPO" show "$ANCHOR:internal/risk/blacklist.go") "$F" | sed 's/^/    /'
  export N='' R=''
}

{
echo "### AC#0 replay start: anchor=$ANCHOR pristine=$PRISTINE tree=$TREE"
echo "### start_head=$(git -C "$REPO" log -1 --format='%H %ct')"
restore; check_hash "restored-from-commit"
run_leg ac0-00-baseline.log

echo "### M-V2a: drop ONLY the fallback branch of the LOCALAPPDATA Microsoft Credentials case"
mutate 'localAppDataS != "" && aUnder(cand, anchorForms(anchorPath(localAppDataS, `microsoft`, `credentials`)))' 'localAppDataS != "" && anchorMatches(cand, anchorForms(anchorPath(localAppDataS, `microsoft`, `credentials`)), isUnder)'
run_leg ac0-01-MV2a.log
restore; check_hash after-M-V2a

echo "### M-V2b: drop ONLY the fallback branch of the ~/.ssh/** case"
mutate 'homeS != "" && aUnder(cand, anchorForms(anchorPath(homeS, `.ssh`)))' 'homeS != "" && anchorMatches(cand, anchorForms(anchorPath(homeS, `.ssh`)), isUnder)'
run_leg ac0-02-MV2b.log
restore; check_hash after-M-V2b

echo "### M-V2c: drop ONLY the fallback branch of the ~/.kube/config case"
mutate 'homeS != "" && aEq(cand, anchorForms(anchorPath(homeS, `.kube`, `config`)))' 'homeS != "" && anchorMatches(cand, anchorForms(anchorPath(homeS, `.kube`, `config`)), eqPath)'
run_leg ac0-03-MV2c.log
restore; check_hash after-M-V2c

echo "### M-V3: opposite shape (STRICTER) - remove the ancestry bound of uncertainAnchorMiss"
mutate 'return isUnder(cand.root, anchor.root) || isUnder(anchor.root, cand.root)' 'return true'
run_leg ac0-04-MV3.log
restore; check_hash after-M-V3

echo "### final state"
check_hash final-blacklist
echo "  shared_worktree_dirty_internal_cmd=$(git -C "$REPO" status --porcelain --untracked-files=no -- internal cmd | wc -l)"
echo "ac0_runner_rc=0"
} 2>&1 | tee "$LOGDIR/ac0-runner.log"
