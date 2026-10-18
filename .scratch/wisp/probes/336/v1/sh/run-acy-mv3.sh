#!/usr/bin/env bash
# 336-v1 extra shot (informs AC#2 + ticket 334 AC#2): does the SHIPPED teeth instrument
# bite on the OPPOSITE shape — M-V3, i.e. uncertainAnchorMiss made STRICTER (its ancestry
# bound removed, so every uncertain candidate counts as a miss)? 334-v1 measured that the
# then-shipped instruments only caught M-V3 by accident, and the gate cell never went red.
# Two serial legs on one exported tree: pristine control, then M-V3 only.
set -u
REPO="D:/work/workspace/projects plans/Wisp"
ANCHOR="${ANCHOR:?anchor}"
TREE="${TREE:?exported tree}"
LOGDIR="${LOGDIR:?logdir}"
GOTMPROOT="${GOTMPROOT:?gotmp}"
F="$TREE/internal/risk/blacklist.go"
PRISTINE=$(git -C "$REPO" rev-parse "$ANCHOR:internal/risk/blacklist.go")
export GOPATH='D:\work\base\gopath'
export PATH="$TREE/third_party/sherpa-onnx:$TREE/build:$PATH"
ROSTER='TestAnchorFallbackTeethPerAnchor336'
restore () { git -C "$REPO" show "$ANCHOR:internal/risk/blacklist.go" > "$F"; }
check_hash () { local h; h=$(git -C "$REPO" hash-object "$F")
  echo "  hash=$h expected=$PRISTINE match=$([ "$h" = "$PRISTINE" ] && echo YES || echo NO)  [$1]"; }
run () { local gt="$GOTMPROOT/$1"; mkdir -p "$gt"
  ( cd "$TREE" && GOTMPDIR="$gt" go test ./internal/risk/ -run "$ROSTER" -count=1 -timeout 30m -v ) > "$LOGDIR/$1" 2>&1
  local rc=$?
  echo "  rc=$rc RUN=$(grep -c '^=== RUN' "$LOGDIR/$1" || true) PASS_TOP=$(grep -c '^--- PASS' "$LOGDIR/$1" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$LOGDIR/$1" || true) SKIP=$(grep -c '^--- SKIP' "$LOGDIR/$1" || true)"
  grep -E '^[[:space:]]*--- FAIL: ' "$LOGDIR/$1" | sed 's/([0-9.]*s)//;s/^/    RED /' || true
}
{
echo "### M-V3 (stricter) on the ticket-336 teeth: anchor=$ANCHOR tree=$TREE"
restore; check_hash start
echo "### leg 1 pristine control:"; run acY-pristine-teeth.log
echo "### leg 2 M-V3 only (remove the ancestry bound):"
export N='return isUnder(cand.root, anchor.root) || isUnder(anchor.root, cand.root)'
export R='return true'
echo "  needle_before=$(grep -cF "$N" "$F" || true)"
perl -0pi -e 'my $q=quotemeta($ENV{N}); s/$q/$ENV{R}/' "$F"
echo "  needle_after=$(grep -cF "$N" "$F" || true) (want 1 -> 0)"
echo "  mutated_hash=$(git -C "$REPO" hash-object "$F")"
run acY-MV3-teeth.log
grep -E 'offchain|TEETH336 operand' "$LOGDIR/acY-MV3-teeth.log" | grep -B1 'FAIL' | head -8 | sed 's/^/    /' || true
restore; check_hash after-MV3
echo "  shared_worktree_dirty_internal_cmd=$(git -C "$REPO" status --porcelain --untracked-files=no -- internal cmd | wc -l)"
echo "acY_rc=0"
} 2>&1 | tee "$LOGDIR/acY-mv3-teeth.log"
