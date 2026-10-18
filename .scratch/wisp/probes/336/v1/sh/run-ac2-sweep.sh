#!/usr/bin/env bash
# 336-v1 AC#2 own sweep: for EACH of the eight A-tier branch operands (aEq/aUnder),
# remove ONLY that one operand's fail-closed branch (swap it for anchorMatches with
# the same relation and the same anchor spelling), run the new teeth roster AND the
# three pre-existing instruments, then restore. Judge-run; own tree, own GOTMPDIR.
# Content anchors (fixed-string needles containing backticks -> passed through env,
# never a double-quoted perl expression). 落点命中枚数 printed per leg. Pristine hash
# checked before and after every leg, in one shell, serially.
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
ROSTER_TEETH='TestAnchorFallbackTeethPerAnchor336'
ROSTER_LEG='TestAnchorMissReasonNamesCandidateAndMissingAnchor|TestAnchorMissKeepsNonOverridableClassAndNeedL2|TestFileBaseNameClaimChangesNoClassOrReason'
BACKMARK='anchorMatches(cand, anchorForms('

restore () { git -C "$REPO" show "$ANCHOR:internal/risk/blacklist.go" > "$F"; }
check_hash () {
  local h; h=$(git -C "$REPO" hash-object "$F")
  echo "  hash=$h expected=$PRISTINE match=$([ "$h" = "$PRISTINE" ] && echo YES || echo NO)  [$1]"
}
run () { # $1=roster $2=log name
  local gt="$GOTMPROOT/$2"; mkdir -p "$gt" || { echo "  mkdir_gotmp_rc=NONZERO"; return 1; }
  ( cd "$TREE" && GOTMPDIR="$gt" go test ./internal/risk/ -run "$1" -count=1 -timeout 30m -v ) > "$LOGDIR/$2" 2>&1
  local rc=$?
  echo "  rc=$rc RUN=$(grep -c '^=== RUN' "$LOGDIR/$2" || true) PASS_TOP=$(grep -c '^--- PASS' "$LOGDIR/$2" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$LOGDIR/$2" || true) SKIP=$(grep -c '^--- SKIP' "$LOGDIR/$2" || true)"
  grep -E '^\s*--- FAIL: ' "$LOGDIR/$2" | sed 's/([0-9.]*s)//;s/^/    RED /' || true
}
mutate () { # $1=needle $2=replacement
  restore
  export N="$1" R="$2"
  local before after lines
  before=$(grep -c -F "$1" "$F" || true)
  perl -0pi -e 'my $q=quotemeta($ENV{N}); s/$q/$ENV{R}/' "$F"
  after=$(grep -c -F "$1" "$F" || true)
  lines=$(diff <(git -C "$REPO" show "$ANCHOR:internal/risk/blacklist.go") "$F" | grep -c '^[<>]' || true)
  echo "  落点命中枚数 needle_before=$before needle_after=$after (want 1 -> 0)"
  echo "  branch_removals_on_disk=$(grep -c -F "$BACKMARK" "$F" || true) (want 1: exactly one operand lost its branch)"
  echo "  changed_lines=$lines"
  diff <(git -C "$REPO" show "$ANCHOR:internal/risk/blacklist.go") "$F" | sed 's/^/    /'
  export N='' R=''
}

{
echo "### AC#2 sweep: anchor=$ANCHOR pristine=$PRISTINE tree=$TREE"
restore; check_hash start
echo "### paired positive control (UNMUTATED bytes)"
echo "  TEETH:"; run "$ROSTER_TEETH" ac2-00-pristine-teeth.log
grep -E 'TEETH336 operand=' "$LOGDIR/ac2-00-pristine-teeth.log" | sed 's/^/    /' || true
echo "  LEG:"; run "$ROSTER_LEG" ac2-00-pristine-leg.log
check_hash after-pristine

while IFS='|' read -r tag needle repl; do
  [ -z "$tag" ] && continue
  echo "### $tag"
  mutate "$needle" "$repl"
  echo "  TEETH:"; run "$ROSTER_TEETH" "ac2-${tag}-teeth.log"
  grep -E 'TEETH336 operand=' "$LOGDIR/ac2-${tag}-teeth.log" | sed 's/^/    /' || true
  echo "  LEG:"; run "$ROSTER_LEG" "ac2-${tag}-leg.log"
  restore; check_hash "after-$tag"
done <<'EOF'
MU-1a|homeS != "" && aEq(cand, anchorForms(anchorPath(homeS, `.git-credentials`)))|homeS != "" && anchorMatches(cand, anchorForms(anchorPath(homeS, `.git-credentials`)), eqPath)
MU-1b|case aUnder(cand, anchorForms(anchorPath(appDataS, `wisp`))) &&|case anchorMatches(cand, anchorForms(anchorPath(appDataS, `wisp`)), isUnder) &&
MU-1c|aEq(cand, anchorForms(anchorPath(appDataS, `wisp`, `config.toml`))):|anchorMatches(cand, anchorForms(anchorPath(appDataS, `wisp`, `config.toml`)), eqPath):
MU-1d|homeS != "" && aUnder(cand, anchorForms(anchorPath(homeS, `.ssh`)))|homeS != "" && anchorMatches(cand, anchorForms(anchorPath(homeS, `.ssh`)), isUnder)
MU-1e|homeS != "" && aEq(cand, anchorForms(anchorPath(homeS, `.aws`, `credentials`)))|homeS != "" && anchorMatches(cand, anchorForms(anchorPath(homeS, `.aws`, `credentials`)), eqPath)
MU-1f|homeS != "" && aEq(cand, anchorForms(anchorPath(homeS, `.kube`, `config`)))|homeS != "" && anchorMatches(cand, anchorForms(anchorPath(homeS, `.kube`, `config`)), eqPath)
MU-1g|appDataS != "" && aUnder(cand, anchorForms(anchorPath(appDataS, `microsoft`, `protect`)))|appDataS != "" && anchorMatches(cand, anchorForms(anchorPath(appDataS, `microsoft`, `protect`)), isUnder)
MU-1h|localAppDataS != "" && aUnder(cand, anchorForms(anchorPath(localAppDataS, `microsoft`, `credentials`)))|localAppDataS != "" && anchorMatches(cand, anchorForms(anchorPath(localAppDataS, `microsoft`, `credentials`)), isUnder)
EOF

echo "### final state"
check_hash final
echo "  shared_worktree_dirty_internal_cmd=$(git -C "$REPO" status --porcelain --untracked-files=no -- internal cmd | wc -l)"
echo "### red roster, teeth legs, per mutation (names only)"
for tag in MU-1a MU-1b MU-1c MU-1d MU-1e MU-1f MU-1g MU-1h; do
  printf "%s teeth_red=%s leg_red=%s\n" "$tag" \
    "$(grep -cE '^\s*--- FAIL: ' "$LOGDIR/ac2-${tag}-teeth.log" || true)" \
    "$(grep -cE '^\s*--- FAIL: ' "$LOGDIR/ac2-${tag}-leg.log" || true)"
done
echo "ac2_sweep_rc=0"
} 2>&1 | tee "$LOGDIR/ac2-sweep-runner.log"
