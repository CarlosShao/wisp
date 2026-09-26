#!/usr/bin/env bash
# Ticket 153 AC#1 readings driver. Everything happens in an OUT-OF-REPO snapshot;
# the shared working tree is only ever read (git show).
set -u
REPO="/d/work/workspace/projects plans/Wisp"
SNAP="/c/Users/swq/AppData/Local/Temp/wisp153-snap-1"
HELP="/c/Users/swq/AppData/Local/Temp/wisp153"
LOG="$HELP/logs"
mkdir -p "$LOG"
ANCHOR="${1:-HEAD}"

restore() {
  git -C "$REPO" show "$ANCHOR:internal/agent/compress.go" > "$SNAP/internal/agent/compress.go"
  git -C "$REPO" show "$ANCHOR:internal/agent/compress_trace_test.go" > "$SNAP/internal/agent/compress_trace_test.go"
  git -C "$REPO" show "$ANCHOR:internal/agent/loop.go" > "$SNAP/internal/agent/loop.go"
}

count() {
  local f="$1"
  printf 'RUN(all)=%s top-PASS=%s sub-PASS=%s FAIL=%s SKIP=%s panic=%s unique=%s\n' \
    "$(grep -c '^=== RUN' "$f")" \
    "$(grep -c '^--- PASS' "$f")" \
    "$(grep -cE '^ +--- PASS' "$f")" \
    "$(grep -cE '^ *--- FAIL' "$f")" \
    "$(grep -cE '^ *--- SKIP' "$f")" \
    "$(grep -c 'panic:' "$f")" \
    "$(grep -E '^--- (PASS|FAIL|SKIP)' "$f" | awk '{print $3}' | sort -u | wc -l)"
}

pkg() { ( cd "$SNAP" && go test ./internal/agent/ -count=1 -v > "$1" 2>&1; echo "go test rc=$?" ); }

echo "== snapshot md5 before restoring (pristine check)"
restore
md5sum "$SNAP/internal/agent/compress.go" "$SNAP/internal/agent/compress_trace_test.go" \
       "$SNAP/internal/agent/loop.go"

echo
echo "== R0  anchor=$ANCHOR, snapshot == anchor, FOUR nails only"
pkg "$LOG/r0-anchor-baseline.txt"; count "$LOG/r0-anchor-baseline.txt"

echo
echo "== R1  M5 applied (guard -> if c.Need(hist)), FOUR nails only"
python "$HELP/m5.py" apply "$SNAP"
( cd "$SNAP" && go build ./internal/agent/ && echo "go build rc=0" )
pkg "$LOG/r1-m5-4nails.txt"; count "$LOG/r1-m5-4nails.txt"

echo
echo "== R2  M5 applied + the new AC#1 case  ==> expect RED"
python "$HELP/m5.py" append "$SNAP" "$HELP/ac1_snippet.txt"
( cd "$SNAP" && go vet ./internal/agent/ && echo "go vet rc=0" )
pkg "$LOG/r2-m5-plus-ac1.txt"; count "$LOG/r2-m5-plus-ac1.txt"
grep -A3 -E '^ *--- FAIL' "$LOG/r2-m5-plus-ac1.txt" | head -20

echo
echo "== R3  M5 reverted, new AC#1 case still present ==> expect GREEN"
python "$HELP/m5.py" revert "$SNAP"
pkg "$LOG/r3-unmutated-plus-ac1.txt"; count "$LOG/r3-unmutated-plus-ac1.txt"
grep -E '^--- (PASS|FAIL): TestCompressionTrace' "$LOG/r3-unmutated-plus-ac1.txt"

echo
echo "== R4  snapshot restored to anchor (byte-exact proof)"
restore
md5sum "$SNAP/internal/agent/compress.go" "$SNAP/internal/agent/compress_trace_test.go" \
       "$SNAP/internal/agent/loop.go"
for f in compress.go compress_trace_test.go loop.go; do
  if cmp -s "$SNAP/internal/agent/$f" <(git -C "$REPO" show "$ANCHOR:internal/agent/$f"); then
    echo "RESTORED IDENTICAL internal/agent/$f"
  else
    echo "MISMATCH internal/agent/$f"
  fi
done
