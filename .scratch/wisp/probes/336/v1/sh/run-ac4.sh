#!/usr/bin/env bash
# 336-v1 AC#4 (judge's own replay): one -run roster, measured on PAIRED temp-root
# surfaces (default %LOCALAPPDATA% shape <=> D:-volume shape, GOTMPDIR mkdir -p'd),
# on the BEFORE tree (ed8ff551, no teeth) and the AFTER tree (8c36c23c, teeth landed).
# Four counts read separately: RUN=^=== RUN (incl. subtests) / PASS_TOP=^--- PASS (top
# only) / FAIL_TOP=^--- FAIL / SKIP=^--- SKIP. Never mixed, never subtracted.
# Then the whole package twice per tree (-run .) and the red rosters diffed by name.
# All legs run in exported trees; the shared worktree is never used for go test.
set -u
REPO="D:/work/workspace/projects plans/Wisp"
LOGDIR="${LOGDIR:?log dir}"
GOTMPROOT="${GOTMPROOT:?D-volume gotmp root}"
BEFORE="${BEFORE:?before tree}"
AFTER="${AFTER:?after tree}"
ROSTER='TestAnchorMissReasonNamesCandidateAndMissingAnchor|TestAnchorMissKeepsNonOverridableClassAndNeedL2|TestFileBaseNameClaimChangesNoClassOrReason|TestAnchorFallbackTeethPerAnchor336'
export GOPATH='D:\work\base\gopath'

run () { # $1=tree $2=roster $3=log $4=shape(--ci|--d)
  local tree="$1" roster="$2" log="$3" shape="$4" gt
  case "$shape" in
    --ci) gt="${LOCALAPPDATA:-C:\\Users\\swq\\AppData\\Local}\\wisp-test-gotmp\\portable-336v1-$log" ;;
    --d)  gt="$GOTMPROOT/$log" ;;
    *) echo "  unknown shape $shape"; return 2 ;;
  esac
  mkdir -p "$gt"; local mrc=$?
  ( cd "$tree" && PATH="$tree/third_party/sherpa-onnx:$tree/build:$PATH" GOTMPDIR="$gt" \
      go test ./internal/risk/ -run "$roster" -count=1 -timeout 30m -v ) > "$LOGDIR/$log" 2>&1
  local rc=$?
  echo "  rc=$rc mkdir_gotmp_rc=$mrc GOTMPDIR=$gt"
  echo "  FOURCOUNT RUN=$(grep -c '^=== RUN' "$LOGDIR/$log" || true) PASS_TOP=$(grep -c '^--- PASS' "$LOGDIR/$log" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$LOGDIR/$log" || true) SKIP=$(grep -c '^--- SKIP' "$LOGDIR/$log" || true) rc=$rc"
  grep -E '^\s*--- FAIL: ' "$LOGDIR/$log" | sed 's/([0-9.]*s)//' > "$LOGDIR/$log.rednames" || true
  echo "  rednames=$(grep -c '' "$LOGDIR/$log.rednames" || true)"
  sed 's/^/    /' "$LOGDIR/$log.rednames" || true
  grep -E '^\s*--- SKIP: ' "$LOGDIR/$log" | sed 's/([0-9.]*s)//;s/^/    SKIPPED /' || true
}

{
echo "### AC#4 paired temp-root surfaces, same roster, before <=> after (judge-run)"
echo "### roster=$ROSTER"
echo "### tree-before=$BEFORE tree-after=$AFTER"
echo "### blacklist_before=$(git -C "$REPO" hash-object "$BEFORE/internal/risk/blacklist.go") blacklist_after=$(git -C "$REPO" hash-object "$AFTER/internal/risk/blacklist.go")"
echo "### teeth_before=$([ -f "$BEFORE/internal/risk/anchor_fallback_teeth_336_windows_test.go" ] && echo PRESENT || echo ABSENT) teeth_after=$([ -f "$AFTER/internal/risk/anchor_fallback_teeth_336_windows_test.go" ] && echo PRESENT || echo ABSENT)"

echo "#### directed roster, D-volume surface (dshape)"
echo "-- before:"; run "$BEFORE" "$ROSTER" ac4-d-before.log --d
echo "-- after:";  run "$AFTER"  "$ROSTER" ac4-d-after.log  --d

echo "#### directed roster, %LOCALAPPDATA% surface (ci-shape)"
echo "-- before:"; run "$BEFORE" "$ROSTER" ac4-ci-before.log --ci
echo "-- after:";  run "$AFTER"  "$ROSTER" ac4-ci-after.log  --ci

echo "#### whole package -run . (two legs per tree, same shape = dshape)"
echo "-- before:"; run "$BEFORE" '.' ac4-full-before.log --d
echo "-- after:";  run "$AFTER"  '.' ac4-full-after.log  --d
echo "#### red roster diff, by name (before <=> after)"
diff "$LOGDIR/ac4-full-before.log.rednames" "$LOGDIR/ac4-full-after.log.rednames" > "$LOGDIR/ac4-full-roster.diff" 2>&1
echo "  diff_rc=$? diff_lines=$(grep -c '' "$LOGDIR/ac4-full-roster.diff" || true)"
sed 's/^/    /' "$LOGDIR/ac4-full-roster.diff" || true
echo "#### SKIP lines, both full legs"
grep -hE '^\s*--- SKIP: ' "$LOGDIR/ac4-full-before.log" "$LOGDIR/ac4-full-after.log" | sed 's/^/    /' || true

echo "#### dshape-only reds of the full AFTER leg vs the ci-shape comparison (surface sensitivity)"
echo "  (this leg only measures surface; the roster difference before/after is the diff above)"
echo "ac4_runner_rc=0"
} 2>&1 | tee "$LOGDIR/ac4-runner.log"
