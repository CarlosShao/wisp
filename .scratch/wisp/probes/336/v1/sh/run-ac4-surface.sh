#!/usr/bin/env bash
# 336-v1 AC#4 surface-sensitivity leg (judge's own): the WHOLE ./internal/risk package,
# -run . , on the AFTER tree, twice — once on the %LOCALAPPDATA% temp-root surface
# (scripts/portable-tests.sh:152-159 shape) and once on the D-volume surface.
# Question measured: are the four TestSync*116 / junction reds a property of the
# ticket-336 bytes, or of the surface? Compare red rosters by name between the two.
set -u
REPO="D:/work/workspace/projects plans/Wisp"
LOGDIR="${LOGDIR:?log dir}"
GOTMPROOT="${GOTMPROOT:?D-volume gotmp root}"
AFTER="${AFTER:?after tree}"
BEFORE="${BEFORE:?before tree}"
export GOPATH='D:\work\base\gopath'

run () { # $1=tree $2=log $3=shape
  local tree="$1" log="$2" shape="$3" gt mrc
  case "$shape" in
    --ci) gt="${LOCALAPPDATA:-C:\\Users\\swq\\AppData\\Local}\\wisp-test-gotmp\\portable-336v1-$log}" ;;
    --d)  gt="$GOTMPROOT/$log" ;;
    *) echo "  unknown shape"; return 2 ;;
  esac
  mkdir -p "$gt"; mrc=$?
  ( cd "$tree" && PATH="$tree/third_party/sherpa-onnx:$tree/build:$PATH" GOTMPDIR="$gt" \
      go test ./internal/risk/ -run '.' -count=1 -timeout 30m -v ) > "$LOGDIR/$log" 2>&1
  local rc=$?
  echo "  rc=$rc mkdir_gotmp_rc=$mrc GOTMPDIR=$gt"
  echo "  FOURCOUNT RUN=$(grep -c '^=== RUN' "$LOGDIR/$log" || true) PASS_TOP=$(grep -c '^--- PASS' "$LOGDIR/$log" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$LOGDIR/$log" || true) SKIP=$(grep -c '^--- SKIP' "$LOGDIR/$log" || true) rc=$rc"
  grep -E '^\s*--- FAIL: ' "$LOGDIR/$log" | sed 's/([0-9.]*s)//' > "$LOGDIR/$log.rednames" || true
  printf 'REDNAME_COUNT=%s (ruler=grep -E "^\\s*--- FAIL: "; empty means zero reds)\n' "$(grep -c '' "$LOGDIR/$log.rednames" || true)" >> "$LOGDIR/$log.rednames"
  sed 's/^/    /' "$LOGDIR/$log.rednames" || true
}

{
echo "### AC#4 surface sensitivity: whole package, ci-shape <=> d-shape, same tree"
echo "-- AFTER tree, ci-shape (%LOCALAPPDATA% temp root):"
run "$AFTER" ac4full-ci-after.log --ci
echo "-- AFTER tree, d-shape (D: temp root):"
run "$AFTER" ac4full-d-after.log --d
echo "-- BEFORE tree, ci-shape:"
run "$BEFORE" ac4full-ci-before.log --ci
echo "### red roster: AFTER ci-shape <=> AFTER d-shape (by name)"
diff "$LOGDIR/ac4full-ci-after.log.rednames" "$LOGDIR/ac4full-d-after.log.rednames" > "$LOGDIR/ac4full-surface.diff" 2>&1
printf 'surface_diff_rc=%s surface_diff_lines=%s (ruler=diff of the two red-name rosters)\n' "$?" "$(grep -c '' "$LOGDIR/ac4full-surface.diff" || true)" >> "$LOGDIR/ac4full-surface.diff"
sed 's/^/    /' "$LOGDIR/ac4full-surface.diff"
echo "### do the ticket-336 teeth files appear in ANY red line of either surface? (want 0)"
for f in ac4full-ci-after ac4full-d-after; do
  printf "%s teeth_red_hits=%s positive_control_anchor_miss_334_red_hits=%s\n" "$f" \
    "$(grep -c 'anchor_fallback_teeth_336' "$LOGDIR/$f.log" || true)" \
    "$(grep -cE '^\s*--- FAIL: TestAnchor' "$LOGDIR/$f.log" || true)"
done
echo "### blacklists unchanged across all four trees/surfaces"
echo "  blacklist_hashes before=$(git -C "$REPO" hash-object "$BEFORE/internal/risk/blacklist.go") after=$(git -C "$REPO" hash-object "$AFTER/internal/risk/blacklist.go")"
echo "### shared worktree internal/cmd dirty (want 0 = I never ran anything there): $(git -C "$REPO" status --porcelain --untracked-files=no -- internal cmd | wc -l)"
echo "ac4_surface_rc=0"
} 2>&1 | tee "$LOGDIR/ac4-surface-runner.log"
