#!/usr/bin/env bash
# 242-r3 third harness: M-B (the mint side writes a digest that is NOT the item's
# own) - ticket 242 v2 measured "all six 242 cases green" under this shape.
# This run checks whether the rewritten case + its positive control now bite.
set -u
REPO="D:/work/workspace/projects plans/Wisp"
MUT="/c/Users/swq/242r3-mut"
MUTW="C:/Users/swq/242r3-mut"
OUT="$REPO/.scratch/wisp/probes/242/r3/logs"
QUEUE_REL="internal/agent/approval/queue.go"
TEST_REL="internal/agent/approval/ticket242_binding_test.go"
: > "$OUT/mut-summary-r3c.txt"

cp "$REPO/$QUEUE_REL" "$MUT/queue-MB.go"
python - "$MUT/queue-MB.go" <<'PY'
import sys
p = sys.argv[1]
s = open(p, encoding='utf-8').read()
old = "\tit.grants.issue(nonce, it.bind)"
new = "\tit.grants.issue(nonce, it.bind + \"m\") // MUTATION M-B: mint side stores a foreign digest"
assert s.count(old) == 1, s.count(old)
open(p, 'w', encoding='utf-8', newline='').write(s.replace(old, new))
PY
grep -n "MUTATION M-B" "$MUT/queue-MB.go" >> "$OUT/mut-summary-r3c.txt"

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
  ( cd "$REPO" && go test -count=1 -timeout 120s -overlay="$json" -v -run "$filt" ./internal/agent/approval/ > "$OUT/mut-$name.txt" 2>&1 )
  local rc=$?
  echo "== MUTATION $name (run=$filt) go_test_rc=$rc" >> "$OUT/mut-summary-r3c.txt"
  grep "^--- FAIL" "$OUT/mut-$name.txt" >> "$OUT/mut-summary-r3c.txt"
  grep -c "^--- FAIL" "$OUT/mut-$name.txt" > /dev/null
  echo "grep_fail_rc=$?" >> "$OUT/mut-summary-r3c.txt"
  grep "AC#1 RED" "$OUT/mut-$name.txt" >> "$OUT/mut-summary-r3c.txt"
  grep "^--- PASS" "$OUT/mut-$name.txt" >> "$OUT/mut-summary-r3c.txt"
  for f in "$QUEUE_REL" "$TEST_REL"; do
    local h b s
    h=$( cd "$REPO" && git hash-object "$f" )
    b=$( cd "$REPO" && git rev-parse "HEAD:$f" )
    s=DIFF; [ "$h" = "$b" ] && s=SAME
    echo "blobcheck $name $f worktree=$h head=$b $s" >> "$OUT/mut-summary-r3c.txt"
  done
}

run_one "OLD2-MB" "TestTicket242" "$REPO/$QUEUE_REL::$MUTW/queue-MB.go" "$REPO/$TEST_REL::$MUTW/test-OLD2.go"
run_one "NEW-MB" "TestTicket242" "$REPO/$QUEUE_REL::$MUTW/queue-MB.go"
echo "harness r3c done" >> "$OUT/mut-summary-r3c.txt"
