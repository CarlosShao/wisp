#!/usr/bin/env bash
# 242-r3 second harness: re-runs the two "before the change" measurements with the
# genuinely toothless copy, and fixes the M-D3 mutation (run 1 produced a BUILD
# FAILURE because the replacement text put a // comment inside a one-line composite
# literal -> "syntax error: unexpected newline in composite literal",
# logs/mut-NEW-MD3.txt). A build failure is not a red, so run 1's MD3 numbers are
# void and this file re-measures them.
set -u
REPO="D:/work/workspace/projects plans/Wisp"
MUT="/c/Users/swq/242r3-mut"
MUTW="C:/Users/swq/242r3-mut"
OUT="$REPO/.scratch/wisp/probes/242/r3/logs"
APPR_REL="internal/agent/approval/approval.go"
QUEUE_REL="internal/agent/approval/queue.go"
TEST_REL="internal/agent/approval/ticket242_binding_test.go"
: > "$OUT/mut-summary-r3b.txt"

# the toothless copy, pinned to the commit BEFORE this leg landed (df1b962e^),
# NOT to HEAD~1 - other legs keep moving HEAD under a running leg.
( cd "$REPO" && git show "df1b962e^:$TEST_REL" ) > "$MUT/test-OLD2.go"
echo "OLD copy lines: $(wc -l < "$MUT/test-OLD2.go") (the new file is $(wc -l < "$REPO/$TEST_REL") lines)" >> "$OUT/mut-summary-r3b.txt"

run_one () {
  local name="$1" filt="$2"; shift 2
  local json="$MUT/overlay-$name.json"
  python - "$json" "$@" <<'PY'
import json, sys
d = {"Replace": {}}
for p in sys.argv[2:]:
    rel, loc = p.split("::")
    d["Replace"][rel] = loc
open(sys.argv[1], "w").write(json.dumps(d))
PY
  if [ "$filt" = "ALL" ]; then
    ( cd "$REPO" && go test -count=1 -overlay="$json" -v ./internal/agent/approval/ > "$OUT/mut-$name.txt" 2>&1 )
  else
    ( cd "$REPO" && go test -count=1 -overlay="$json" -v -run "$filt" ./internal/agent/approval/ > "$OUT/mut-$name.txt" 2>&1 )
  fi
  local rc=$?
  echo "== MUTATION $name (run=$filt) go_test_rc=$rc" >> "$OUT/mut-summary-r3b.txt"
  echo "-- $name red cases:" >> "$OUT/mut-summary-r3b.txt"
  grep "^--- FAIL" "$OUT/mut-$name.txt" >> "$OUT/mut-summary-r3b.txt"
  grep -c "^--- FAIL" "$OUT/mut-$name.txt" > /dev/null
  echo "grep_fail_rc=$?" >> "$OUT/mut-summary-r3b.txt"
  echo "-- $name FAIL reason lines:" >> "$OUT/mut-summary-r3b.txt"
  grep "AC#1 RED" "$OUT/mut-$name.txt" >> "$OUT/mut-summary-r3b.txt"
  echo "-- $name pass cases:" >> "$OUT/mut-summary-r3b.txt"
  grep "^--- PASS" "$OUT/mut-$name.txt" >> "$OUT/mut-summary-r3b.txt"
  echo "-- $name overlay sanity: the seq ruler must be ABSENT from an OLD-file run and PRESENT otherwise" >> "$OUT/mut-summary-r3b.txt"
  grep -c "BindDigestSeparatesItemsBySequenceNumber" "$OUT/mut-$name.txt" >> "$OUT/mut-summary-r3b.txt"
  for f in "$APPR_REL" "$QUEUE_REL" "$TEST_REL"; do
    local h b s
    h=$( cd "$REPO" && git hash-object "$f" )
    b=$( cd "$REPO" && git rev-parse "HEAD:$f" )
    s=DIFF; [ "$h" = "$b" ] && s=SAME
    echo "blobcheck $name $f worktree=$h head=$b $s" >> "$OUT/mut-summary-r3b.txt"
  done
}

# ---------- M-D3, done so it still compiles ----------
cp "$REPO/$APPR_REL" "$MUT/approval-MD3b.go"
python - "$MUT/approval-MD3b.go" <<'PY'
import sys
p = sys.argv[1]
s = open(p, encoding='utf-8').read()
old = "strconv.FormatUint(seq, 10)"
new = "strconv.FormatUint(0, 10)"
assert s.count(old) == 1, s.count(old)
open(p, 'w', encoding='utf-8', newline='').write(s.replace(old, new))
PY
grep -n "strconv.FormatUint(0, 10)" "$MUT/approval-MD3b.go" >> "$OUT/mut-summary-r3b.txt"

# 1. the toothless case as it stood BEFORE this leg, unmutated: green
run_one "OLD2-plain" "TestTicket242" "$REPO/$TEST_REL::$MUTW/test-OLD2.go"
# 2. the toothless case + M-A: still green on the cross-card case = no teeth
run_one "OLD2-MA" "TestTicket242" "$REPO/$APPR_REL::$MUTW/approval-MA.go" "$REPO/$TEST_REL::$MUTW/test-OLD2.go"
# 3. the rewritten case + seq folded out: the new seq ruler must be the red one
run_one "NEW-MD3b" "TestTicket242" "$REPO/$APPR_REL::$MUTW/approval-MD3b.go"
# 4. same, whole package: v2 measured "M-D3 => whole package green"; this is the after
run_one "NEW-MD3b-FULL" "ALL" "$REPO/$APPR_REL::$MUTW/approval-MD3b.go"
echo "harness r3b done" >> "$OUT/mut-summary-r3b.txt"
