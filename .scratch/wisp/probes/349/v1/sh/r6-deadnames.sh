#!/bin/sh
cd '/d/work/workspace/projects plans/Wisp' || exit 9
BASE=070de7fa1b38b5c5d52445564c7ae19d1a73d6bd
LEG=ed8db1ea27eb50d3057dd0cd7d5c955a13e978fc
echo '######## PRECISION SHAPES ON DELIVERED SIDE (blob tree, git grep -E word-boundary) ########'
echo '--- 1) bare Complete() call shape: exact token with word boundaries ---'
for R in $BASE $LEG; do
  printf 'REV=%s CompleteCallExact=%s\n' "$R" "$(git grep -c -E '(^|[^A-Za-z0-9_])Complete\(\)' $R -- internal/risk 2>/dev/null | wc -l)"
done
echo '--- 1b) positive control: SyncDetectionComplete present (should be >0) ---'
printf 'SYNCDETECTIONCOMPLETE_del=%s\n' "$(git grep -c 'SyncDetectionComplete' $LEG -- internal/risk 2>/dev/null | wc -l)"
printf 'SYNCDETECTIONCOMPLETE_lines_del=%s\n' "$(git grep -n 'SyncDetectionComplete' $LEG -- internal/risk 2>/dev/null | wc -l)"
echo '--- 1c) the BAD ruler (bare substring Complete) for contrast: how many hits ---'
printf 'BARESUBSTR_base=%s del=%s\n' "$(git grep -c 'Complete' $BASE -- internal/risk 2>/dev/null | awk -F: '{s+=$NF} END{print s+0}')" "$(git grep -c 'Complete' $LEG -- internal/risk 2>/dev/null | awk -F: '{s+=$NF} END{print s+0}')"
echo '--- 2) fires() exact call shape ---'
for R in $BASE $LEG; do
  printf 'REV=%s firesCallExact=%s\n' "$R" "$(git grep -n -E '(^|[^A-Za-z0-9_.])fires\(' $R -- internal/risk 2>/dev/null | wc -l)"
done
echo '--- 2b) positive control: bare substring fires (nonzero) ---'
printf 'FIRES_BARESUBSTR_del=%s\n' "$(git grep -n 'fires' $LEG -- internal/risk 2>/dev/null | wc -l)"
git grep -n 'fires' $LEG -- internal/risk 2>/dev/null | head -5
echo '--- 3) taskScopeID exact word ---'
for R in $BASE $LEG; do
  printf 'REV=%s taskScopeID=%s\n' "$R" "$(git grep -n -E '(^|[^A-Za-z0-9_])taskScopeID([^A-Za-z0-9_]|$)' $R -- internal/risk 2>/dev/null | wc -l)"
done
echo '--- 3b) positive control scopeID present ---'
printf 'SCOPEID_del_lines=%s\n' "$(git grep -n 'scopeID' $LEG -- internal/risk 2>/dev/null | wc -l)"
echo '--- 4) disposalScope exact word ---'
for R in $BASE $LEG; do
  printf 'REV=%s disposalScope=%s\n' "$R" "$(git grep -n -E '(^|[^A-Za-z0-9_])disposalScope([^A-Za-z0-9_]|$)' $R -- internal/risk 2>/dev/null | wc -l)"
done
echo '--- 4b) positive control DisposalScope (real type, other pkg) whole-repo ---'
printf 'DISPOSALSCOPE_repo_lines=%s\n' "$(git grep -n 'DisposalScope' $LEG -- internal 2>/dev/null | wc -l)"
echo RC=$?
