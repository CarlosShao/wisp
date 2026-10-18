#!/usr/bin/env bash
# 336-v1 mandatory-(ii) measurement: does a THIRD world make the ":224" aUnder operand
# class-bearing on its own? Runs the read-only probe (sh/probe-third-world.txt) twice
# on one exported tree: (1) unmutated, (2) with ONLY the ":224" operand's fallback
# branch removed. If the probe's deep row reads class=B under that mutation and A
# without it, a third world gives :224 an independent red; if class stays A, the
# overlap 1b<=>1c is structural and needs a different candidate placement, not a file.
set -u
REPO="D:/work/workspace/projects plans/Wisp"
ANCHOR="${ANCHOR:?anchor}"
TREE="${TREE:?exported tree}"
LOGDIR="${LOGDIR:?log dir}"
GOTMPROOT="${GOTMPROOT:?gotmp root}"
F="$TREE/internal/risk/blacklist.go"
PRISTINE=$(git -C "$REPO" rev-parse "$ANCHOR:internal/risk/blacklist.go")
PROBE_SRC="$REPO/.scratch/wisp/probes/336/v1/sh/probe-third-world.txt"
PROBE_DST="$TREE/internal/risk/zz_probe_336v1_test.go"
export GOPATH='D:\work\base\gopath'
export PATH="$TREE/third_party/sherpa-onnx:$TREE/build:$PATH"
ROSTER='TestProbe336V1ThirdWorld|TestAnchorFallbackTeethPerAnchor336'

restore () { git -C "$REPO" show "$ANCHOR:internal/risk/blacklist.go" > "$F"; }
check_hash () { local h; h=$(git -C "$REPO" hash-object "$F")
  echo "  hash=$h expected=$PRISTINE match=$([ "$h" = "$PRISTINE" ] && echo YES || echo NO)  [$1]"; }
run () { local gt="$GOTMPROOT/$2"; mkdir -p "$gt"
  ( cd "$TREE" && GOTMPDIR="$gt" go test ./internal/risk/ -run "$1" -count=1 -timeout 30m -v ) > "$LOGDIR/$2" 2>&1
  local rc=$?
  echo "  rc=$rc RUN=$(grep -c '^=== RUN' "$LOGDIR/$2" || true) PASS_TOP=$(grep -c '^--- PASS' "$LOGDIR/$2" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$LOGDIR/$2" || true) SKIP=$(grep -c '^--- SKIP' "$LOGDIR/$2" || true)"
  grep -E 'V3WORLD|--- FAIL' "$LOGDIR/$2" | sed 's/^/    /' || true
}
{
echo "### mandatory-(ii) third-world probe: anchor=$ANCHOR tree=$TREE"
cp "$PROBE_SRC" "$PROBE_DST" || { echo "copy_rc=NONZERO"; exit 1; }
echo "  probe_blob=$(git -C "$REPO" hash-object "$PROBE_DST") lines=$(grep -c '' "$PROBE_DST")"
restore; check_hash start
echo "### leg 1 UNMUTATED:"; run "$ROSTER" m2-pristine.log
check_hash after-pristine
echo "### leg 2 MU-1b only (drop the :224 aUnder operand's fallback branch):"
export N='case aUnder(cand, anchorForms(anchorPath(appDataS, `wisp`))) &&'
export R='case anchorMatches(cand, anchorForms(anchorPath(appDataS, `wisp`)), isUnder) &&'
echo "  needle_before=$(grep -cF "$N" "$F" || true)"
perl -0pi -e 'my $q=quotemeta($ENV{N}); s/$q/$ENV{R}/' "$F"
echo "  needle_after=$(grep -cF "$N" "$F" || true) branch_removals=$(grep -cF 'anchorMatches(cand, anchorForms(' "$F" || true) (want 1)"
git -C "$REPO" hash-object "$F" | sed 's/^/  mutated_hash=/'
run "$ROSTER" m2-MU1b.log
restore; check_hash after-MU1b
echo "### leg 3 MU-1c only (the aEq half of the same case), same probe, for the pair:"
export N2='aEq(cand, anchorForms(anchorPath(appDataS, `wisp`, `config.toml`))):'
export R2='anchorMatches(cand, anchorForms(anchorPath(appDataS, `wisp`, `config.toml`)), eqPath):'
echo "  needle_before=$(grep -cF "$N2" "$F" || true)"
perl -0pi -e 'my $q=quotemeta($ENV{N2}); s/$q/$ENV{R2}/' "$F"
echo "  needle_after=$(grep -cF "$N2" "$F" || true) branch_removals=$(grep -cF 'anchorMatches(cand, anchorForms(' "$F" || true) (want 1)"
git -C "$REPO" hash-object "$F" | sed 's/^/  mutated_hash=/'
run "$ROSTER" m2-MU1c.log
restore; check_hash after-MU1c
echo "### probe stays in the exported tree (repo untouched): repo_hits=$(git -C "$REPO" ls-files -- 'internal/risk/zz_probe*' | wc -l) dirty_internal_cmd=$(git -C "$REPO" status --porcelain --untracked-files=no -- internal cmd | wc -l)"
echo "m2_rc=0"
} 2>&1 | tee "$LOGDIR/m2-third-world-runner.log"
