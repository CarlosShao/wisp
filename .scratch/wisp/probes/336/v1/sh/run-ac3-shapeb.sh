#!/usr/bin/env bash
# 336-v1 AC#3 shape B (乙形) decisive shot, run by the JUDGE on its own exported tree.
# The instrument = the 42-line-shaped read-only probe ruler 334-v1 shipped as a
# non-tracked artefact (.scratch/wisp/probes/334/v1/logs/probe-334v1_test.go.txt,
# blob 69a70fa454886029aea9f8f716cb372908467339). It is copied into the EXPORTED tree
# only (never the repo, never tracked) under its probe name.
# Two legs, serial, same tree, same roster, content anchor, pristine hash before/after:
#   (1) unmutated   -> what does it report, and is it green?
#   (2) MU-1h (drop ONLY the fallback branch of the %LOCALAPPDATA%\Microsoft\Credentials
#       operand, the same needle the AC#2 sweep used) -> does rc still equal 0?
# That second leg is the "judgment criterion swapped to the opposite shape still all
# green = insensitive" test the dispatch makes mandatory.
# Assertion census of the probe (my own ruler, not the leg's):
set -u
REPO="D:/work/workspace/projects plans/Wisp"
ANCHOR="${ANCHOR:?named commit anchor}"
TREE="${TREE:?exported tree root}"
LOGDIR="${LOGDIR:?log dir}"
GOTMPROOT="${GOTMPROOT:?D-volume gotmp root}"
PROBE_SRC="$REPO/.scratch/wisp/probes/334/v1/logs/probe-334v1_test.go.txt"
PROBE_DST="$TREE/internal/risk/zz_probe_334v1_test.go"
F="$TREE/internal/risk/blacklist.go"
PRISTINE=$(git -C "$REPO" rev-parse "$ANCHOR:internal/risk/blacklist.go")
export GOPATH='D:\work\base\gopath'
export PATH="$TREE/third_party/sherpa-onnx:$TREE/build:$PATH"
ROSTER_PROBE='TestProbe334V1'
ROSTER_TEETH='TestAnchorFallbackTeethPerAnchor336'

restore () { git -C "$REPO" show "$ANCHOR:internal/risk/blacklist.go" > "$F"; }
check_hash () {
  local h; h=$(git -C "$REPO" hash-object "$F")
  echo "  hash=$h expected=$PRISTINE match=$([ "$h" = "$PRISTINE" ] && echo YES || echo NO)  [$1]"
}
run () { # $1=roster $2=log $3=tag
  local gt="$GOTMPROOT/$2"; mkdir -p "$gt" || { echo "  mkdir_gotmp_rc=NONZERO"; return 1; }
  local s e
  s=$(date +%s.%N)
  ( cd "$TREE" && GOTMPDIR="$gt" go test ./internal/risk/ -run "$1" -count=1 -timeout 30m -v ) > "$LOGDIR/$2" 2>&1
  local rc=$?
  e=$(date +%s.%N)
  echo "  rc=$rc wall=$(awk -v a="$s" -v b="$e" 'BEGIN{printf "%.2f", b-a}')s RUN=$(grep -c '^=== RUN' "$LOGDIR/$2" || true) PASS_TOP=$(grep -c '^--- PASS' "$LOGDIR/$2" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$LOGDIR/$2" || true) SKIP=$(grep -c '^--- SKIP' "$LOGDIR/$2" || true)  [$3]"
}

{
echo "### AC#3 shape B decisive shot: anchor=$ANCHOR pristine=$PRISTINE tree=$TREE"
echo "### probe blob (source artefact) = $(git -C "$REPO" rev-parse HEAD:.scratch/wisp/probes/334/v1/logs/probe-334v1_test.go.txt)"
cp "$PROBE_SRC" "$PROBE_DST" || { echo "probe_copy_rc=NONZERO"; exit 1; }
echo "  probe_on_disk_blob=$(git -C "$REPO" hash-object "$PROBE_DST")"
echo "  probe_lines=$(grep -c "" "$PROBE_DST")"
echo "### assertion census of the probe (own ruler, counts inside the file just placed):"
echo "  t_Log=$(grep -c -E 't\.Log(f)?\(' "$PROBE_DST" || true)  t_Fatalf=$(grep -c -E 't\.Fatalf\(' "$PROBE_DST" || true)  t_Errorf_or_Error=$(grep -c -E 't\.Errorf?\(' "$PROBE_DST" || true)  t_Skip=$(grep -c -E 't\.Skip' "$PROBE_DST" || true)"
echo "  build_tags=$(grep -c '^//go:build' "$PROBE_DST" || true)  (0 => compiles on both platforms)"
restore; check_hash start

echo "### leg 1: UNMUTATED bytes, roster=$ROSTER_PROBE"
run "$ROSTER_PROBE" ac3-b-pristine.log pristine-shapeB
echo "  chain lines (the ones a real assertion would need): V1SHAPE=$(grep -c 'V1SHAPE ' "$LOGDIR/ac3-b-pristine.log" || true) V1OVERRIDE=$(grep -c 'V1OVERRIDE ' "$LOGDIR/ac3-b-pristine.log" || true)"
grep -E 'chain:case9-local-creds|chain:case4b-ssh-solo' "$LOGDIR/ac3-b-pristine.log" | sed 's/^/    /' || true
check_hash after-leg1

echo "### leg 2: MU-1h only (same content needle the AC#2 sweep used)"
export N='localAppDataS != "" && aUnder(cand, anchorForms(anchorPath(localAppDataS, `microsoft`, `credentials`)))'
export R='localAppDataS != "" && anchorMatches(cand, anchorForms(anchorPath(localAppDataS, `microsoft`, `credentials`)), isUnder)'
echo "  落点命中枚数 before=$(grep -c -F "$N" "$F" || true)"
perl -0pi -e 'my $q=quotemeta($ENV{N}); s/$q/$ENV{R}/' "$F"
echo "  落点命中枚数 after=$(grep -c -F "$N" "$F" || true) (want 1 -> 0)"
echo "  branch_removals_on_disk=$(grep -c -F 'anchorMatches(cand, anchorForms(' "$F" || true) (want 1)"
diff <(git -C "$REPO" show "$ANCHOR:internal/risk/blacklist.go") "$F" | sed 's/^/    /'
echo "  mutated_hash=$(git -C "$REPO" hash-object "$F")"
run "$ROSTER_PROBE" ac3-b-MU1h.log MU1h-shapeB
echo "  --- shapeB under the tooth-removal (does the CLASS change show up as red?):"
grep -E 'chain:case9-local-creds' "$LOGDIR/ac3-b-MU1h.log" | sed 's/^/    /' || true
run "$ROSTER_TEETH" ac3-a-MU1h.log MU1h-shapeA
echo "  --- shapeA (the shipped teeth) under the SAME single mutation:"
grep -E '^\s*--- FAIL: ' "$LOGDIR/ac3-a-MU1h.log" | sed 's/^/    /' || true
restore; check_hash after-leg2

echo "### leg 3: shape A pristine timing on the same tree for the cost pair"
run "$ROSTER_TEETH" ac3-a-pristine.log pristine-shapeA
restore; check_hash after-leg3

echo "### shape files, pristine vs MU-1h (the readable diff the probe produces, if any)"
grep -E '^ +anchor_fallback|V1SHAPE chain:case9|V1SHAPE chain:case4b' "$LOGDIR/ac3-b-pristine.log" > "$LOGDIR/ac3-b-pristine.shape" || true
grep -E '^ +anchor_fallback|V1SHAPE chain:case9|V1SHAPE chain:case4b' "$LOGDIR/ac3-b-MU1h.log" > "$LOGDIR/ac3-b-MU1h.shape" || true
echo "  shape_diff_lines=$(diff "$LOGDIR/ac3-b-pristine.shape" "$LOGDIR/ac3-b-MU1h.shape" | grep -c '^[<>]' || true)"
diff "$LOGDIR/ac3-b-pristine.shape" "$LOGDIR/ac3-b-MU1h.shape" | sed 's/^/    /'
echo "  probe_stays_in_exported_tree=$([ -f "$PROBE_DST" ] && echo YES || echo NO) (临时件只建不删；它只在仓外树里，⛔ 进仓、⛔ 进跟踪面)"
echo "  repo_side_probe_hits=$(git -C "$REPO" ls-files -- 'internal/risk/zz_probe*' | wc -l) (want 0)"
echo "  shared_worktree_dirty_internal_cmd=$(git -C "$REPO" status --porcelain --untracked-files=no -- internal cmd | wc -l)"
echo "ac3_runner_rc=0"
} 2>&1 | tee "$LOGDIR/ac3-runner.log"
