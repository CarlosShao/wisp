#!/usr/bin/env bash
# 212-v3 second battery: which NODE SET feeds ban #9 (f.Comments only?),
# and which TREES are in range (tools/**, _test.go).
TREE=/d/tmp/212v3tree
LOGS="/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/212/v3/logs"
cd "/d/work/workspace/projects plans/Wisp/tools/d22scan" || exit 9

run_case() {
  local name="$1"
  go run . -root "D:/tmp/212v3tree" > "$LOGS/scope-$name.txt" 2>&1
  local rc=$?
  echo "rc=$rc" >> "$LOGS/scope-$name.txt"
  local hits toks
  hits=$(grep -c "phantom-citation" "$LOGS/scope-$name.txt")
  toks=$(grep -o 'comment cites repo path "[^"]*"' "$LOGS/scope-$name.txt" | sed 's/^/    /')
  echo "CASE $name rc=$rc phantom_lines=$hits" >> "$LOGS/scope-summary.txt"
  if [ -n "$toks" ]; then echo "$toks" >> "$LOGS/scope-summary.txt"; fi
}

: > "$LOGS/scope-summary.txt"
echo "second battery at $(date '+%Y-%m-%d %H:%M:%S%z')" >> "$LOGS/scope-summary.txt"

# S1 a phantom path only inside a STRING LITERAL of a production file
printf 'package probe\n\nconst beacon = "read docs/evidence/s1/212-v3-string-only.md first"\n\nfunc s1() {}\n' > "$TREE/internal/probe/zs1.go"
run_case s1-string-literal

# S2 a phantom path inside a /* block */ comment
printf 'package probe\n\n/* block form: docs/evidence/s1/212-v3-blockcomment.md is absent */\nfunc s2() {}\n' > "$TREE/internal/probe/zs2.go"
run_case s2-block-comment

# S3 a phantom path in a _test.go comment (declared OUT of range)
printf 'package probe\n\n// see docs/evidence/s1/212-v3-testonly.md for readings\nfunc TestS3(t *testing.T) {}\n' > "$TREE/internal/probe/zs3_test.go"
printf 'package probe\n\nimport "testing"\n' >> "$TREE/internal/probe/zs3_test.go"
run_case s3-underscore-test-go

# S4 a phantom path in tools/** production go comment (declared OUT of range)
mkdir -p "$TREE/tools/probe"
printf 'package probe\n\n// see docs/evidence/s1/212-v3-toolsonly.md for readings\nfunc s4() {}\n' > "$TREE/tools/probe/s4.go"
run_case s4-tools-tree

# S5 all four carriers removed -> back to baseline green
rm -f "$TREE/internal/probe/zs1.go" "$TREE/internal/probe/zs2.go" "$TREE/internal/probe/zs3_test.go" "$TREE/tools/probe/s4.go"
rmdir "$TREE/tools/probe" 2>/dev/null
run_case s5-carriers-removed

echo "leftover probe dir: $(ls "$TREE/internal/probe" | tr '\n' ' ')" >> "$LOGS/scope-summary.txt"
