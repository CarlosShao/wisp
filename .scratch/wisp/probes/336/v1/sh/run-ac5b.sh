#!/usr/bin/env bash
# 336-v1 AC#5 follow-ups: (a) the span in item #3 crosses OTHER legs' commits, so the
# blob-equality claim must be measured per-commit over the leg's own 18, not over a
# range; (b) aEquiv stale-comment provenance; (c) case count with a literal tab.
set -u
REPO="D:/work/workspace/projects plans/Wisp"
cd "$REPO" || exit 1
D=".scratch/wisp/probes/336/v1/logs"
{
echo "### (a) per-commit blob equality for the frozen / gate-semantic files ACROSS THE LEG'S OWN 18 COMMITS"
for p in internal/risk/thresholds.go internal/tools/rules_gateway.go internal/risk/tokens_fourway_test.go internal/risk/l2_grant_boundary_test.go internal/risk/ticket90_persist_test.go internal/risk/blacklist.go; do
  uniq=$(while read -r h; do git rev-parse "$h:$p" 2>/dev/null || echo MISSING; done < "$D/00-commits-336r1.txt" | sort -u | wc -l)
  printf "%-46s distinct_blobs_in_18_commits=%s want=1 same_as_anchor=%s\n" "$p" "$uniq" \
    "$([ "$(git rev-parse "e5021339:$p")" = "$(git rev-parse "ed8ff551:$p")" ] && echo YES || echo 'NO (another leg moved it in the range, not this one — see (a2))')"
done
echo "### (a2) who actually moved those files in ed8ff551..e5021339 (subject prefixes, so attribution is by leg, not by me guessing)"
for p in internal/risk/thresholds.go internal/tools/rules_gateway.go internal/risk/tokens_fourway_test.go; do
  echo "-- $p"
  git log --format='    %h %s' ed8ff551..e5021339 -- "$p" | cut -c1-110
done
echo "### (b) aEquiv: the only hit in the package is a COMMENT naming a helper that has no definition. Provenance:"
git log -S 'aEquiv' --format='  %h %ad %s' --date=short -- internal/risk/blacklist.go | cut -c1-140
echo "### (c) A-tier case count in classifyForms (literal tab ruler):"
TAB=$(printf '\t')
git show e5021339:internal/risk/blacklist.go | sed -n '219,239p' | grep -c "^${TAB}case " || true
echo "  and of those, how many carry an aEq/aUnder operand (the branch-capable ones):"
git show e5021339:internal/risk/blacklist.go | sed -n '219,239p' | grep -cE "aEq|aUnder" || true
echo "ac5b_rc=0"
} 2>&1 | tee "$D/06-forbidden-faces-followup.log"
