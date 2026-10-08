#!/usr/bin/env bash
# 242-r3 (implementer leg) self-mutation harness.
# Every mutated copy lives OUTSIDE the repo (C:/Users/swq/242r3-mut) and is
# redirected in with `go test -overlay`, so the repo is never written - the
# blob check at the end of every run proves it.
# Every mutation BREAKS PRODUCTION. None of them touches an assertion, so this
# is not "摘尺", it is "种错法".
set -u
REPO="D:/work/workspace/projects plans/Wisp"
MUT="/c/Users/swq/242r3-mut"
MUTW="C:/Users/swq/242r3-mut"
OUT="$REPO/.scratch/wisp/probes/242/r3/logs"
APPR_REL="internal/agent/approval/approval.go"
QUEUE_REL="internal/agent/approval/queue.go"
TEST_REL="internal/agent/approval/ticket242_binding_test.go"
mkdir -p "$MUT"
: > "$OUT/mut-summary.txt"

run_one () {
  local name="$1" filt="$2"; shift 2
  local json="$MUT/overlay-$name.json"
  python - "$json" "$@" <<'PY'
import json, sys
json_path = sys.argv[1]
d = {"Replace": {}}
for p in sys.argv[2:]:
    rel, loc = p.split("::")
    d["Replace"][rel] = loc
open(json_path, "w").write(json.dumps(d))
PY
  if [ "$filt" = "ALL" ]; then
    ( cd "$REPO" && go test -count=1 -overlay="$json" -v ./internal/agent/approval/ > "$OUT/mut-$name.txt" 2>&1 )
  else
    ( cd "$REPO" && go test -count=1 -overlay="$json" -v -run "$filt" ./internal/agent/approval/ > "$OUT/mut-$name.txt" 2>&1 )
  fi
  local rc=$?
  echo "== MUTATION $name (run=$filt) go_test_rc=$rc" >> "$OUT/mut-summary.txt"
  echo "-- $name red cases:" >> "$OUT/mut-summary.txt"
  grep "^--- FAIL" "$OUT/mut-$name.txt" >> "$OUT/mut-summary.txt"
  grep -c "^--- FAIL" "$OUT/mut-$name.txt" > /dev/null
  echo "grep_fail_rc=$?" >> "$OUT/mut-summary.txt"
  echo "-- $name FAIL reason lines:" >> "$OUT/mut-summary.txt"
  grep -A1 "^    ticket242_binding_test.go" "$OUT/mut-$name.txt" | grep "AC#1 RED" >> "$OUT/mut-summary.txt"
  echo "-- $name pass cases:" >> "$OUT/mut-summary.txt"
  grep "^--- PASS" "$OUT/mut-$name.txt" >> "$OUT/mut-summary.txt"
  for f in "$APPR_REL" "$QUEUE_REL" "$TEST_REL"; do
    local h b
    h=$( cd "$REPO" && git hash-object "$f" )
    b=$( cd "$REPO" && git rev-parse "HEAD:$f" )
    local s=DIFF
    [ "$h" = "$b" ] && s=SAME
    echo "blobcheck $name $f worktree=$h head=$b $s" >> "$OUT/mut-summary.txt"
  done
}

# ---------- the pre-change (toothless) test file, straight out of git ----------
( cd "$REPO" && git show "HEAD~1:$TEST_REL" ) > "$MUT/test-OLD.go"
wc -l < "$MUT/test-OLD.go" >> "$OUT/mut-summary.txt"

# ---------- M-A: the binding comparison is emptied (misbind accepted) ----------
cp "$REPO/$APPR_REL" "$MUT/approval-MA.go"
python - "$MUT/approval-MA.go" <<'PY'
import sys
p = sys.argv[1]
s = open(p, encoding='utf-8').read()
old = """		delete(s.values, v) // consumed whether or not the binding matched
		if !equalSecret(stored, bind) {
			return denialMisbound
		}"""
new = """		delete(s.values, v) // consumed whether or not the binding matched
		bindingEnforced := false // MUTATION M-A: the binding layer is gone
		if bindingEnforced && !equalSecret(stored, bind) {
			return denialMisbound
		}"""
assert s.count(old) == 1, s.count(old)
open(p, 'w', encoding='utf-8', newline='').write(s.replace(old, new))
PY

# ---------- M-D3: seq stops being folded into the digest ----------
cp "$REPO/$APPR_REL" "$MUT/approval-MD3.go"
python - "$MUT/approval-MD3.go" <<'PY'
import sys
p = sys.argv[1]
s = open(p, encoding='utf-8').read()
old = "strconv.FormatUint(seq, 10)"
new = 'strconv.FormatUint(0, 10) // MUTATION M-D3: seq folded out ("seq" kept referenced so the copy still builds)'
assert s.count(old) == 1, s.count(old)
open(p, 'w', encoding='utf-8', newline='').write(s.replace(old, new))
PY

# ---------- M-SAME: bindDigest stops separating cards at all ----------
cp "$REPO/$APPR_REL" "$MUT/approval-SAME.go"
python - "$MUT/approval-SAME.go" <<'PY'
import sys
p = sys.argv[1]
s = open(p, encoding='utf-8').read()
old1 = 'for _, s := range []string{corr, taskID, tool, level, strconv.FormatUint(seq, 10)} {'
new1 = 'for _, s := range []string{"same", strconv.FormatUint(0, 10)} { // MUTATION M-SAME: identity fields dropped'
old2 = "\t_, _ = sum.Write(args)\n\treturn hex.EncodeToString(sum.Sum(nil))"
new2 = '\t_, _ = sum.Write([]byte("same-args")) // MUTATION M-SAME: args dropped too -> one digest for every card\n\treturn hex.EncodeToString(sum.Sum(nil))'
assert s.count(old1) == 1 and s.count(old2) == 1, (s.count(old1), s.count(old2))
s = s.replace(old1, new1).replace(old2, new2)
open(p, 'w', encoding='utf-8', newline='').write(s)
PY

# ---------- M-E: the membership scan stops refusing unknown nonces ----------
cp "$REPO/$APPR_REL" "$MUT/approval-ME.go"
python - "$MUT/approval-ME.go" <<'PY'
import sys
p = sys.argv[1]
s = open(p, encoding='utf-8').read()
old = "\t}\n\treturn denialSpentNonce\n}"
new = "\t}\n\treturn denialNone // MUTATION M-E: an unrecognised nonce is accepted\n}"
assert s.count(old) == 1, s.count(old)
open(p, 'w', encoding='utf-8', newline='').write(s.replace(old, new))
PY

# ================= runs =================
# 1. the OLD case as delivered: green, which is the defect being fixed
run_one "OLDplain" "TestTicket242" "$REPO/$TEST_REL::$MUTW/test-OLD.go"
# 2. OLD case + M-A: still green on the cross-card case => it had no teeth
run_one "OLD-MA" "TestTicket242" "$REPO/$APPR_REL::$MUTW/approval-MA.go" "$REPO/$TEST_REL::$MUTW/test-OLD.go"
# 3. NEW case + M-A: must turn red (this leg's acceptance core)
run_one "NEW-MA" "TestTicket242" "$REPO/$APPR_REL::$MUTW/approval-MA.go"
# 4. same mutation, whole package (collateral read)
run_one "NEW-MA-FULL" "ALL" "$REPO/$APPR_REL::$MUTW/approval-MA.go"
# 5. NEW + seq folded out: the new seq ruler must be the one that goes red
run_one "NEW-MD3" "TestTicket242" "$REPO/$APPR_REL::$MUTW/approval-MD3.go"
run_one "NEW-MD3-FULL" "ALL" "$REPO/$APPR_REL::$MUTW/approval-MD3.go"
# 6. NEW + every card shares one digest
run_one "NEW-SAME" "TestTicket242" "$REPO/$APPR_REL::$MUTW/approval-SAME.go"
# 7. NEW + unknown nonce accepted: the discriminator assertion must go red
run_one "NEW-ME" "TestTicket242" "$REPO/$APPR_REL::$MUTW/approval-ME.go"

echo "harness done" >> "$OUT/mut-summary.txt"
