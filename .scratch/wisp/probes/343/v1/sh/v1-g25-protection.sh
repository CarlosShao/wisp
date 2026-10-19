#!/bin/sh
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 90
A=/d/tmp/wisp/20261019/343v1/logs/ac2-protection
mkdir -p $A
{
echo "== protected instruments: blob equality START(29d258f9) vs HEAD =="
for f in internal/risk/anchor_fallback_teeth_336_windows_test.go \
         internal/risk/anchor_fallback_teeth_336_other_test.go \
         internal/risk/anchor_miss_reason_339_other_test.go \
         internal/risk/anchor_miss_reason_334_test.go \
         internal/risk/pathresolver.go internal/risk/pathresolver_other.go \
         internal/risk/assessor.go internal/risk/syncdirs.go internal/risk/provenance.go \
         internal/panel/tokens_fourway_test.go internal/panel/l2_grant_boundary_test.go \
         internal/perm/ticket90_persist_test.go internal/observe/thresholds.go; do
  b1=$(git rev-parse 29d258f9:$f 2>/dev/null || echo ABSENT-AT-START)
  b2=$(git rev-parse HEAD:$f 2>/dev/null || echo ABSENT-AT-HEAD)
  if [ "$b1" = "$b2" ]; then s=SAME; else s=DIFFERENT; fi
  s8=$(printf '%s' "$b1" | cut -c1-8); h8=$(printf '%s' "$b2" | cut -c1-8)
  printf '   %s start=%s head=%s  %s\n' "$s" "$s8" "$h8" "$f"
done
echo "== blacklist.go itself: the ONLY .go that changed across the leg's 6 commits =="
git diff --name-only 29d258f9 137044b6 -- internal cmd tools | sed 's/^/  CHANGED /'
echo "== what changed repo-wide in internal/** between leg start and leg delivery =="
git diff --numstat 29d258f9 137044b6 -- internal cmd tools docs .scratch/wisp/issues
echo "== the 4 protected-instrument blobs (short, as leg reported) =="
for f in anchor_fallback_teeth_336_windows_test.go anchor_fallback_teeth_336_other_test.go anchor_miss_reason_339_other_test.go anchor_miss_reason_334_test.go; do
  printf '%s start=%s head=%s\n' "$f" "$(git rev-parse --short 29d258f9:internal/risk/$f)" "$(git rev-parse --short HEAD:internal/risk/$f)"
done
echo DONE2
} > $A/protection.log 2>&1
echo DONE
cat $A/protection.log
