#!/usr/bin/env bash
# Ticket 153 AC#1 readings, taken with the DELIVERED test file (so every line
# number quoted from a red sentence is a line number in the shipped file).
# All runs happen in an OUT-OF-REPO snapshot; the working tree is only read.
set -u
REPO="/d/work/workspace/projects plans/Wisp"
SNAP="/c/Users/swq/AppData/Local/Temp/wisp153-snap-1"
HELP="/c/Users/swq/AppData/Local/Temp/wisp153"
LOG="$HELP/logs"
mkdir -p "$LOG"
ANCHOR="${1:-86b0161}"

restore() {
  git -C "$REPO" show "$ANCHOR:internal/agent/compress.go" > "$SNAP/internal/agent/compress.go"
  git -C "$REPO" show "$ANCHOR:internal/agent/compress_trace_test.go" > "$SNAP/internal/agent/compress_trace_test.go"
  git -C "$REPO" show "$ANCHOR:internal/agent/loop.go" > "$SNAP/internal/agent/loop.go"
}

count() {
  local f="$1"
  printf 'RUN(all)=%s top-PASS=%s sub-PASS=%s FAIL=%s SKIP=%s panic=%s\n' \
    "$(grep -c '^=== RUN' "$f")" "$(grep -c '^--- PASS' "$f")" "$(grep -cE '^ +--- PASS' "$f")" \
    "$(grep -cE '^ *--- FAIL' "$f")" "$(grep -cE '^ *--- SKIP' "$f")" "$(grep -c 'panic:' "$f")"
}

pkg() { ( cd "$SNAP" && go test ./internal/agent/ -count=1 -v > "$1" 2>&1; echo "go test rc=$?" ); }

restore
for f in compress.go compress_trace_test.go loop.go; do
  cmp -s "$SNAP/internal/agent/$f" <(git -C "$REPO" show "$ANCHOR:internal/agent/$f") \
    && echo "pristine internal/agent/$f (anchor $ANCHOR)" || echo "MISMATCH $f"
done

echo
echo "== F1  anchor compress.go + DELIVERED test file, no mutation ==> expect GREEN"
cp "$REPO/internal/agent/compress_trace_test.go" "$SNAP/internal/agent/compress_trace_test.go"
grep -n 'func TestCompressionTraceSilentWhenNothingFoldableOverThreshold' "$SNAP/internal/agent/compress_trace_test.go"
pkg "$LOG/f1-delivered-unmutated.txt"; count "$LOG/f1-delivered-unmutated.txt"
grep -E '^--- (PASS|FAIL): TestCompressionTrace' "$LOG/f1-delivered-unmutated.txt"

echo
echo "== F2  same, with M5 applied to compress.go ==> expect RED (this is the 改前红句)"
python "$HELP/m5.py" apply "$SNAP"
( cd "$SNAP" && go build ./internal/agent/ && echo "go build rc=0" )
pkg "$LOG/f2-delivered-m5.txt"; count "$LOG/f2-delivered-m5.txt"
sed -n '/RUN   TestCompressionTraceSilentWhenNothingFoldableOverThreshold/,/^--- FAIL: TestCompressionTraceSilentWhenNothingFoldable/p' "$LOG/f2-delivered-m5.txt"

echo
echo "== F3  restore compress.go to anchor, re-run ==> expect GREEN again"
python "$HELP/m5.py" revert "$SNAP"
pkg "$LOG/f3-restored.txt"; count "$LOG/f3-restored.txt"
restore
for f in compress.go compress_trace_test.go loop.go; do
  cmp -s "$SNAP/internal/agent/$f" <(git -C "$REPO" show "$ANCHOR:internal/agent/$f") \
    && echo "RESTORED IDENTICAL internal/agent/$f" || echo "MISMATCH $f"
done
