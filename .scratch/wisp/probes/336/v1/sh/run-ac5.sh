#!/usr/bin/env bash
# 336-v1 AC#5 forbidden-face rulers (judge-run on named commits; each count carries a
# positive control). Read-only: git show / git grep / git diff on named trees.
set -u
REPO="D:/work/workspace/projects plans/Wisp"
cd "$REPO" || exit 1
D=".scratch/wisp/probes/336/v1/logs"
{
echo "### 1. internal/risk numstat, before(ed8ff551) -> teeth commit(8c36c23c)"
git diff --numstat ed8ff551 8c36c23c -- internal/risk/
echo "### 2. whole leg span (ed8ff551 -> e5021339) numstat for risk / tools / panel / cmd-wisp"
git diff --numstat ed8ff551 e5021339 -- internal/risk/ internal/tools/ internal/panel/ cmd/wisp/
echo "### 3. frozen + gate-semantic files: blob equality across the span"
for p in internal/risk/thresholds.go internal/tools/rules_gateway.go internal/risk/tokens_fourway_test.go internal/risk/l2_grant_boundary_test.go internal/risk/ticket90_persist_test.go go.mod go.sum .gitignore .github/workflows/ci.yml docs/PLAN.md; do
  b=$(git rev-parse "ed8ff551:$p" 2>/dev/null || echo MISSING)
  a=$(git rev-parse "e5021339:$p" 2>/dev/null || echo MISSING)
  if [ "$b" = "$a" ]; then s=YES; else s=NO; fi
  printf "%-46s same=%s before=%s\n" "$p" "$s" "${b:0:10}"
done
echo "### 4. allowlist / golden / threshold files touched by the 18 commits (want 0)"
grep -icE 'allowlist|golden|threshold' "$D/00-paths-by-legs.txt" || true
echo "  positive control (the same file greps non-zero for a token that IS there): $(grep -icE 'internal/risk' "$D/00-paths-by-legs.txt" || true)"
echo "### 5. new exported names in the two test files (want: Test* only)"
for p in internal/risk/anchor_fallback_teeth_336_windows_test.go internal/risk/anchor_fallback_teeth_336_other_test.go; do
  echo "-- $p"
  git show "e5021339:$p" | grep -E '^(func|var|type|const) [A-Z]' | sed 's/^/    /'
done
echo "### 6. t.Skip in the two test files (want 0) + positive control that grep works here"
for p in internal/risk/anchor_fallback_teeth_336_windows_test.go internal/risk/anchor_fallback_teeth_336_other_test.go; do
  printf "%s t_Skip=%s t_Log=%s build_tag=%s\n" "$p" \
    "$(git show "e5021339:$p" | grep -cE 't\.Skip' || true)" \
    "$(git show "e5021339:$p" | grep -cE 't\.Log' || true)" \
    "$(git show "e5021339:$p" | head -1)"
done
echo "### 7. stale-comment reading: blacklist.go:193 names aEquiv — is there such a function? (production bytes, NOT the leg's)"
git grep -c 'aEquiv' e5021339 -- internal/risk/ || true
echo "  definition hits: $(git grep -c 'func aEquiv' e5021339 -- internal/risk/ || true) (want 0 => the comment names an absent helper)"
echo "  positive control: func aEq definitions = $(git grep -c 'func aEq' e5021339 -- internal/risk/ || true)"
echo "### 8. the two A-tier cases that do NOT go through aEq/aUnder: who pins them (AC#1 completeness)"
echo "-- browser credential store:"
git grep -n 'browser credential store' e5021339 -- internal/risk/ internal/tools/ | sed 's|e5021339:||' | sed 's/^/    /'
echo "-- .git/config A rule:"
git grep -nE 'A_tier_|git/config' e5021339 -- 'internal/risk/*_test.go' | sed 's|e5021339:||' | sed 's/^/    /' | head -10
echo "### 9. A-tier case count vs operand count (my own ruler, scope = the switch at blacklist.go)"
git show e5021339:internal/risk/blacklist.go | sed -n '219,239p' | grep -cE '^\tcase ' || true
echo "ac5_rc=0"
} 2>&1 | tee "$D/05-forbidden-faces.log"
