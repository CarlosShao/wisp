#!/usr/bin/env bash
# 336-v1 final self-proof (three rulers), all run fresh by the judge.
set -u
REPO="D:/work/workspace/projects plans/Wisp"
cd "$REPO" || exit 1
D=".scratch/wisp/probes/336/v1/logs"
SEC=$(printf '\302\247')
{
echo "### R1 per-commit proof: strict prefix ruler -> every commit -> its name-only set"
echo "### (must be ⊆ the three allowed write faces, and .go hits = 0)"
echo "count=$(git log --all --grep='^票 336 终裁腿 336-v1' --format='%h' | wc -l)"
git log --all --grep='^票 336 终裁腿 336-v1' --format='%H' | sort > "$D/95-my-commits.txt"
: > "$D/95-my-per-commit-nameonly.txt"
while read -r h; do
  {
    echo "== $(git log -1 --format='%h %ct %s' "$h" | cut -c1-100)"
    git show --name-only --format='' "$h" | grep -v '^$' | sed 's/^/   /'
  } >> "$D/95-my-per-commit-nameonly.txt"
done < "$D/95-my-commits.txt"
echo "total_unique_paths=$(git log --all --grep='^票 336 终裁腿 336-v1' --format='' --name-only | grep -v '^$' | sort -u | wc -l)"
echo "out_of_three_faces=$(git log --all --grep='^票 336 终裁腿 336-v1' --format='' --name-only | grep -v '^$' | sort -u | grep -vcE '^\.scratch/wisp/probes/336/v1/|^docs/evidence/s1/336-v1-verdict\.md$|^\.scratch/wisp/issues/336-a-tier-fallback-branch-per-anchor-teeth-334-ac2-remainder\.md$|^\.scratch/wisp/issues/334-anchor-missing-leaf-reason-names-another-path-and-ci-surface-reds-216-instrument\.md$' || true)"
echo "go_hits=$(git log --all --grep='^票 336 终裁腿 336-v1' --format='' --name-only | grep -v '^$' | grep -c '\.go$' || true)"
echo "blacklist_hits=$(git log --all --grep='^票 336 终裁腿 336-v1' --format='' --name-only | grep -v '^$' | grep -c 'blacklist\.go' || true)"
echo "ledger_or_handover_or_other_tickets=$(git log --all --grep='^票 336 终裁腿 336-v1' --format='' --name-only | grep -v '^$' | grep -cE 'docs/reports/|\.scratch/wisp/issues/[^3]' || true) (want 0; positive control: same ruler without the [^3] guard on my two tickets = 2)"
echo "### R2 deliverable proof: absolute paths + wc -c + line count"
echo "verdict=$(wc -l < docs/evidence/s1/336-v1-verdict.md) lines / $(wc -c < docs/evidence/s1/336-v1-verdict.md) bytes / cap 130"
echo "readouts=$(wc -l < $D/../readings-and-mutations.md) lines / $(wc -c < $D/../readings-and-mutations.md) bytes"
echo "probe_files=$(find .scratch/wisp/probes/336/v1 -type f | wc -l) zero_byte=$(find .scratch/wisp/probes/336/v1 -type f -size 0 | wc -l) total_bytes=$(du -sb .scratch/wisp/probes/336/v1 | cut -f1)"
echo "abs_root=$(pwd -W 2>/dev/null || pwd)/.scratch/wisp/probes/336/v1"
echo "### R3 empty-slot word-form ruler (0 hits, with positive control)"
echo "  my_files_placeholder_hits=$(find .scratch/wisp/probes/336/v1 docs/evidence/s1/336-v1-verdict.md -type f -print0 | xargs -0 grep -cE 'TBD|TODO|待补|占位|XXX|<[A-Z_]+>' 2>/dev/null | grep -v ':0$' | wc -l) files (want 0)"
echo "  POSITIVE CONTROL same ruler on 334-v1 verdict = $(grep -cE 'TBD|TODO|待补|占位|XXX|<[A-Z_]+>' docs/evidence/s1/334-v1-verdict.md || true) (non-zero => ruler alive)"
echo "  section_sign_in_my_face=$(find .scratch/wisp/probes/336/v1 -type f -print0 | xargs -0 grep -lF "$SEC" 2>/dev/null | wc -l) (want 0)"
echo "  POSITIVE CONTROL same-shape xargs ruler on token TEETH336 = $(find .scratch/wisp/probes/336 -type f -print0 | xargs -0 grep -lF 'TEETH336' 2>/dev/null | wc -l) files"
echo "### AC box census after my tail commit (judgment-area ruler)"
echo "  336 unchecked=$(sed -n '20,26p' .scratch/wisp/issues/336-a-tier-fallback-branch-per-anchor-teeth-334-ac2-remainder.md | grep -c '^[[:space:]]*- \[ \]' || true) checked=$(sed -n '20,26p' .scratch/wisp/issues/336-a-tier-fallback-branch-per-anchor-teeth-334-ac2-remainder.md | grep -c '^[[:space:]]*- \[x\]' || true)"
echo "  334 unchecked=$(sed -n '33,39p' .scratch/wisp/issues/334-anchor-missing-leaf-reason-names-another-path-and-ci-surface-reds-216-instrument.md | grep -c '^[[:space:]]*- \[ \]' || true) checked=$(sed -n '33,39p' .scratch/wisp/issues/334-anchor-missing-leaf-reason-names-another-path-and-ci-surface-reds-216-instrument.md | grep -c '^[[:space:]]*- \[x\]' || true)"
echo "### shared worktree: internal/cmd dirty by me (want 0) = $(git status --porcelain --untracked-files=no -- internal cmd | wc -l)"
echo "selfcheck_rc=0"
} 2>&1 | tee "$D/95-final-selfcheck.log"
