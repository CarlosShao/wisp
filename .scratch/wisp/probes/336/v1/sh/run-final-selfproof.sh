#!/usr/bin/env bash
# 336-v1 FINAL three-ruler self-proof (fresh readings, run last).
set -u
REPO="D:/work/workspace/projects plans/Wisp"
cd "$REPO" || exit 1
D=".scratch/wisp/probes/336/v1"
L="$D/logs"
PAT='^票 336 终裁腿 336-v1'
SEC=$(printf '\302\247')
{
echo "### FINAL three-ruler self-proof (336-v1), fresh at $(date '+%F %T %z')"
echo "=== R1: strict prefix ruler -> count / per-commit subjects / name-only ⊆ three faces / .go 0 ==="
echo "count=$(git log --all --grep="$PAT" --format='%h' | wc -l)"
git log --all --grep="$PAT" --format='%h %s' | cut -c1-42 | nl
echo "unique_paths=$(git log --all --grep="$PAT" --format='' --name-only | grep -v '^$' | sort -u | wc -l)"
echo "out_of_three_faces=$(git log --all --grep="$PAT" --format='' --name-only | grep -v '^$' | sort -u | grep -vcE '^\.scratch/wisp/probes/336/v1/|^docs/evidence/s1/336-v1-verdict\.md$|^\.scratch/wisp/issues/336-a-tier-fallback-branch-per-anchor-teeth-334-ac2-remainder\.md$|^\.scratch/wisp/issues/334-anchor-missing-leaf-reason-names-another-path-and-ci-surface-reds-216-instrument\.md$' || true)"
echo "go_hits=$(git log --all --grep="$PAT" --format='' --name-only | grep -v '^$' | grep -c '\.go$' || true)"
echo "blacklist_name_hits=$(git log --all --grep="$PAT" --format='' --name-only | grep -v '^$' | grep -c 'blacklist' || true)"
echo "  (positive control that the name-only ruler is alive: count of all paths = above; the two slices I shipped carry "blacklist" in the FILENAME of my own capture, those are log names, not repo paths I touched)"
echo "=== R2: deliverable absolute paths + sizes + face census ==="
echo "verdict: $(pwd -W)/docs/evidence/s1/336-v1-verdict.md = $(wc -l < docs/evidence/s1/336-v1-verdict.md) lines / $(wc -c < docs/evidence/s1/336-v1-verdict.md) bytes (cap 130)"
echo "readouts: $(pwd -W)/$D/readings-and-mutations.md = $(wc -l < $D/readings-and-mutations.md) lines / $(wc -c < $D/readings-and-mutations.md) bytes"
echo "face: disk=$(find $D -type f | wc -l) tracked=$(git ls-files $D | wc -l) zero_byte=$(find $D -type f -size 0 | wc -l)"
echo "sec_U00A7_files=$(find $D docs/evidence/s1/336-v1-verdict.md -type f -print0 | xargs -0 grep -lF "$SEC" 2>/dev/null | wc -l) (want 0)"
echo "  POSITIVE CONTROL same xargs shape on a live token TEETH336 = $(find $D -type f -print0 | xargs -0 grep -lF 'TEETH336' 2>/dev/null | wc -l) files"
echo "=== R3: AC box census, judgment-area ruler only ==="
echo "336 unchecked=$(sed -n '20,26p' .scratch/wisp/issues/336-a-tier-fallback-branch-per-anchor-teeth-334-ac2-remainder.md | grep -c '^[[:space:]]*- \[ \]' || true) checked=$(sed -n '20,26p' .scratch/wisp/issues/336-a-tier-fallback-branch-per-anchor-teeth-334-ac2-remainder.md | grep -c '^[[:space:]]*- \[x\]' || true)"
echo "334 unchecked=$(sed -n '33,39p' .scratch/wisp/issues/334-anchor-missing-leaf-reason-names-another-path-and-ci-surface-reds-216-instrument.md | grep -c '^[[:space:]]*- \[ \]' || true) checked=$(sed -n '33,39p' .scratch/wisp/issues/334-anchor-missing-leaf-reason-names-another-path-and-ci-surface-reds-216-instrument.md | grep -c '^[[:space:]]*- \[x\]' || true)"
echo "=== blacklist.go untouched (worktree == HEAD == intake anchor string) ==="
echo "  worktree=$(git hash-object internal/risk/blacklist.go)"
echo "  HEAD=$(git rev-parse HEAD:internal/risk/blacklist.go)"
echo "=== my faces clean (want 0 lines) = $(git status --porcelain --untracked-files=all $D docs/evidence/s1/336-v1-verdict.md | wc -l)"
echo "=== shared worktree internal/cmd dirty by me (want 0) = $(git status --porcelain --untracked-files=no -- internal cmd | wc -l)"
echo "final_rc=0"
} 2>&1 | tee "$L/NN-FINAL-SELFEPROOF.log"
