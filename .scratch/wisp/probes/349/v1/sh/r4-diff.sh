#!/bin/sh
cd '/d/work/workspace/projects plans/Wisp' || exit 9
BASE=070de7fa1b38b5c5d52445564c7ae19d1a73d6bd
LEG=ed8db1ea27eb50d3057dd0cd7d5c955a13e978fc
echo '=== git diff -U0 BASE..LEG for internal/risk (verbatim) ==='
git diff -U0 $BASE $LEG -- internal/risk
echo '=== numstat ==='
git diff --numstat $BASE $LEG -- internal/risk
echo '=== every ^[+-] line classified: starts-with-// after leading +/- ? ==='
git diff -U0 $BASE $LEG -- internal/risk | grep -E '^[+-]' | grep -vE '^(\+\+\+|---)' > /d/tmp/wisp/20261019/349v1/logs2/diff-lines.txt
TOT=$(wc -l < /d/tmp/wisp/20261019/349v1/logs2/diff-lines.txt)
NONC=$(grep -cvE '^([+])[[:space:]]*//|^(-)[[:space:]]*//' /d/tmp/wisp/20261019/349v1/logs2/diff-lines.txt || true)
printf 'DIFFLINE_TOTAL=%s NONCOMMENTLINE=%s\n' "$TOT" "$NONC"
echo '--- the non-comment ones (should be none) ---'
grep -E '^([+])[[:space:]]*[^/]|^(-)[[:space:]]*[^/]' /d/tmp/wisp/20261019/349v1/logs2/diff-lines.txt || echo NONE
echo RC=$?
