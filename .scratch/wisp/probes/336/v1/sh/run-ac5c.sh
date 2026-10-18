#!/usr/bin/env bash
# 336-v1 AC#5 item (a) CORRECTED. My first pass (run-ac5.sh #3 / run-ac5b.sh (a)) used
# WRONG PATHS (internal/tools/rules_gateway.go, internal/risk/thresholds.go,
# internal/risk/tokens_fourway_test.go ...): those paths do not exist in the tree, so
# every count below that ruler was a MISSING-vs-MISSING comparison = a broken ruler, and
# its "same=NO / distinct=19" lines must NOT be read as findings. Correct paths taken
# from `git ls-tree -r --name-only e5021339` (see item (a3) below for the ruler).
# This version measures real blobs, per-commit over the leg's own 18 commits.
set -u
REPO="D:/work/workspace/projects plans/Wisp"
cd "$REPO" || exit 1
D=".scratch/wisp/probes/336/v1/logs"
{
echo "### (a3) path ruler that produced the corrected names:"
git ls-tree -r --name-only e5021339 | grep -iE 'thresholds\.go|rules_gateway\.go|tokens_fourway_test\.go|l2_grant_boundary_test\.go|ticket90_persist_test\.go|d22scan/allowlist\.txt' | sed 's/^/    /'
echo "### (a4) corrected per-commit blob census across the leg's 18 commits + anchor equality"
for p in internal/observe/thresholds.go internal/risk/rules_gateway.go internal/panel/tokens_fourway_test.go internal/panel/l2_grant_boundary_test.go internal/perm/ticket90_persist_test.go tools/d22scan/allowlist.txt internal/risk/blacklist.go go.mod go.sum .gitignore .github/workflows/ci.yml docs/PLAN.md; do
  before=$(git rev-parse "ed8ff551:$p" 2>/dev/null || echo MISSING)
  after=$(git rev-parse "e5021339:$p" 2>/dev/null || echo MISSING)
  uniq=$(while read -r h; do git rev-parse "$h:$p" 2>/dev/null || echo MISSING; done < "$D/00-commits-336r1.txt" | sort -u | wc -l)
  if [ "$before" = "$after" ]; then s=YES; else s=NO; fi
  printf "%-42s distinct_blobs_in_18=%s anchor_same=%s blob=%s\n" "$p" "$uniq" "$s" "${before:0:10}"
done
echo "### (a5) positive controls for the rulers above:"
echo "  control-1 (a file the leg DID add => distinct must be 1 and anchor_same NO):"
p=internal/risk/anchor_fallback_teeth_336_windows_test.go
before=$(git rev-parse "ed8ff551:$p" 2>/dev/null || echo MISSING); after=$(git rev-parse "e5021339:$p" 2>/dev/null || echo MISSING)
uniq=$(while read -r h; do git rev-parse "$h:$p" 2>/dev/null || echo MISSING; done < "$D/00-commits-336r1.txt" | sort -u | wc -l)
[ "$before" = "$after" ] && s=YES || s=NO
printf "    %s distinct_blobs_in_18=%s anchor_same=%s (before=%s after=%s)\n" "$p" "$uniq" "$s" "${before:0:8}" "${after:0:8}"
echo "  control-2 (blacklist.go blob list, all 18 lines): $(while read -r h; do git rev-parse "$h:internal/risk/blacklist.go"; done < "$D/00-commits-336r1.txt" | sort -u | tr '\n' ' ')"
echo "ac5c_rc=0"
} 2>&1 | tee "$D/07-forbidden-faces-corrected.log"
