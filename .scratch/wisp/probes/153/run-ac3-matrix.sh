#!/usr/bin/env bash
# Ticket 153 AC#3: the load-bearing matrix, run entirely in an OUT-OF-REPO
# snapshot of the delivered code (commit given as $1).
set -u
REPO="/d/work/workspace/projects plans/Wisp"
HELP="/c/Users/swq/AppData/Local/Temp/wisp153"
SNAP="/c/Users/swq/AppData/Local/Temp/wisp153-snap-3"
PRIS="$HELP/pristine-3"
LOG="$HELP/logs"
SHA="${1:-HEAD}"
mkdir -p "$LOG" "$PRIS/internal/agent"

mkdir -p "$SNAP"
git -C "$REPO" -c core.autocrlf=false -c core.eol=lf archive --format=tar "$SHA" | tar -x -C "$SNAP"
for f in compress.go loop.go compress_trace_test.go; do
  cp "$SNAP/internal/agent/$f" "$PRIS/internal/agent/$f"
done
echo "== matrix runs on commit $SHA; md5 of the three files under test:"
md5sum "$PRIS"/internal/agent/*.go

restore() { python "$HELP/mut153.py" restore "$SNAP" "$PRIS" >/dev/null; }

# run <label> <op>...   -- restores first, applies ops, builds, runs the package
run() {
  local label="$1"; shift
  restore
  for op in "$@"; do
    python "$HELP/mut153.py" "$op" "$SNAP" "$PRIS" || { echo "!! $label: op $op refused"; return 2; }
  done
  ( cd "$SNAP" && go build ./internal/agent/ >/dev/null 2>&1 ) || { echo "!! $label: build FAILED"; return 2; }
  ( cd "$SNAP" && go test ./internal/agent/ -count=1 -v > "$LOG/ac3-$label.txt" 2>&1 )
  local rc=$?
  local reds
  reds=$(grep -E '^ *--- FAIL' "$LOG/ac3-$label.txt" | awk '{print $3}' | sort -u | tr '\n' ' ')
  printf '%-34s rc=%s  RUN=%s  FAILs=%s  reds: %s\n' "$label" "$rc" \
    "$(grep -c '^=== RUN' "$LOG/ac3-$label.txt")" "$(grep -cE '^ *--- FAIL' "$LOG/ac3-$label.txt")" \
    "${reds:-NONE}"
}

echo
echo "== X0 delivered code, no mutation (matrix baseline)"
run X0-delivered
echo "== X1 M5 only: does the AC#1 case bite now?"
run X1-m5 m5
echo "== X2 M5 + remove the AC#1 case  (摘掉 N1)"
run X2-m5-noAC1 m5 drop-test:TestCompressionTraceSilentWhenNothingFoldableOverThreshold
echo "== X3 loop stops tagging (摘掉 N2)"
run X3-unloop unloop
echo "== X4 摘掉 N2 + 摘掉 its witness test"
run X4-unloop-no-taskid-test unloop drop-test:TestCompressionTraceCarriesTheOwningTaskID
echo "== X5 compressor stops booking the attribute (摘掉 N3)"
run X5-drop-attr drop-attr
echo "== X6 attribute booked even when nothing was tagged (placeholder laundering)"
run X6-always-attr always-attr
echo "== X7 traceTaskID fabricates a constant id"
run X7-fake-id fake-id
echo "== X8 摘掉 N3 + 摘掉 both AC#2 tests"
run X8-dropattr-no-ac2-tests drop-attr \
    drop-test:TestCompressionTraceCarriesTheOwningTaskID \
    drop-test:TestCompressionTraceNeverInventsATaskID

echo
echo "== restore proof: snapshot files vs commit $SHA"
restore
for f in compress.go loop.go compress_trace_test.go; do
  cmp -s "$SNAP/internal/agent/$f" <(git -C "$REPO" show "$SHA:internal/agent/$f") \
    && echo "RESTORED IDENTICAL internal/agent/$f" || echo "MISMATCH internal/agent/$f"
done
