#!/usr/bin/env bash
# 336-v1 AC#3/AC#1 supporting rulers: line-kind census of the shape-B probe output
# (how many readable lines a ported assertion would have to carry named wants for),
# and whether the leg's report still carries the claim it retracted.
set -u
REPO="D:/work/workspace/projects plans/Wisp"
cd "$REPO" || exit 1
D=".scratch/wisp/probes/336/v1/logs"
{
echo "### line-kind census of the shape-B pristine leg (40-ac3-b-pristine.log)"
for k in "V1SHAPE chain:" "V1SHAPE chain-clean:" "V1SHAPE offchain:" "V1SHAPE real-temp" "V1SHAPE real-appdata" "V1OVERRIDE chain:" "V1OVERRIDE real-appdata" "V1PREFIX anchor_root" "V1LSTAT"; do
  printf "  %-26s %s\n" "$k" "$(grep -cF "$k" "$D/40-ac3-b-pristine.log" || true)"
done
echo "  V1SHAPE_total=$(grep -c 'V1SHAPE ' "$D/40-ac3-b-pristine.log" || true) V1OVERRIDE_total=$(grep -c 'V1OVERRIDE ' "$D/40-ac3-b-pristine.log" || true)"
echo "  chain-family lines (chain + chain-clean + offchain) = $(grep -cE '^.*V1SHAPE (chain:|chain-clean:|offchain:)' "$D/40-ac3-b-pristine.log" || true)"
echo "### does any under-row world make the aEq half a PROVEN match (i.e. the row is not a branch-alone tooth)? measured by the teeth log itself:"
grep -E 'TEETH336 operand=' "$D/30-ac2-00-pristine-teeth.log" | sed 's/^/    /'
echo "### the leg report at e5021339:"
echo "  lines=$(git show e5021339:.scratch/wisp/probes/336/r1/report.md | grep -c '')"
echo "  retracted_claim_still_there=$(git show e5021339:.scratch/wisp/probes/336/r1/report.md | grep -c '闸门格在红名册里' || true) (want 0)"
echo "  honest_line_present=$(git show e5021339:.scratch/wisp/probes/336/r1/report.md | grep -cE '恒不红' || true)"
echo "support_rc=0"
} 2>&1 | tee "$D/08-supporting-rulers.log"
