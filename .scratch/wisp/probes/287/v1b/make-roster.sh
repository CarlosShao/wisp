#!/usr/bin/env bash
# 287-v1b：从逐枚 meta 重建突变名册（bench-run.sh 每批会清空 roster，故名册一律由本件现生成）
set -u
OUT="$(cd "$(dirname "$0")" && pwd)"
R="$OUT/mutation-roster.txt"
: >"$R"
for s in V1 V2 V3 V4 V5 V6 V7 V8; do
  m="$OUT/logs/v1b-$s.meta.txt"
  [ -f "$m" ] || { echo "===== $s ===== META MISSING" >>"$R"; continue; }
  {
    echo "===== $s ====="
    grep -E '^(SHOT|PRE-SHA|anchor_hits|WROTE|mutate_rc|POST-SHA|hunks|changed_lines|test_rc|fault-tokens|full_rc|full_TOP_PASS|RESTORED-EQUAL-OK|RESTORE-MISMATCH|MOVED)' "$m"
    sed -n '/--- case-colors ---/,/--- red-points ---/p' "$m"
    sed -n '/--- red-sentences/,$p' "$m"
  } >>"$R"
done
wc -c "$R"
