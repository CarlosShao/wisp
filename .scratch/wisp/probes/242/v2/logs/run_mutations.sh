#!/usr/bin/env bash
# 242-v2 acceptance leg: independent mutations against the grant binding layer.
# All mutated copies live OUTSIDE the repo; `go test -overlay` redirects them in.
# Nothing here relaxes an assertion - every mutation BREAKS production.
set -u
REPO="D:/work/workspace/projects plans/Wisp"
MUT="/c/Users/swq/242v2-mut"
MUTW="C:/Users/swq/242v2-mut"
OUT="$REPO/.scratch/wisp/probes/242/v2/logs"
APPR_REL="internal/agent/approval/approval.go"
QUEUE_REL="internal/agent/approval/queue.go"
mkdir -p "$MUT"

run_one () {
  local name="$1"; shift
  local json="$MUT/overlay-$name.json"
  python - "$json" "$@" <<'PY'
import json,sys
json_path=sys.argv[1]; pairs=sys.argv[2:]
d={"Replace":{}}
for p in pairs:
    rel,loc=p.split("::")
    d["Replace"][rel]=loc
open(json_path,"w").write(json.dumps(d))
PY
  ( cd "$REPO" && go test -count=1 -overlay="$json" -v -run 'TestTicket242' ./internal/agent/approval/ \
      > "$OUT/mut-$name-test.txt" 2>&1 )
  echo "MUTATION $name go_test_rc=$?" | tee -a "$OUT/mut-summary.txt"
  echo "-- $name red cases:" | tee -a "$OUT/mut-summary.txt"
  grep "^--- FAIL" "$OUT/mut-$name-test.txt" | tee -a "$OUT/mut-summary.txt"
  grep -c "^--- FAIL" "$OUT/mut-$name-test.txt" > /dev/null; echo "grep_fail_rc=$?" | tee -a "$OUT/mut-summary.txt"
  echo "-- $name pass cases:" | tee -a "$OUT/mut-summary.txt"
  grep "^--- PASS" "$OUT/mut-$name-test.txt" | tee -a "$OUT/mut-summary.txt"
}

# baseline roster for this run
: > "$OUT/mut-summary.txt"

# ---------- M-A: binding comparison removed (misbind accepted) ----------
cp "$REPO/$APPR_REL" "$MUT/approval-MA.go"
python - "$MUT/approval-MA.go" <<'PY'
import sys
p=sys.argv[1]; s=open(p,encoding='utf-8').read()
old="""		delete(s.values, v) // consumed whether or not the binding matched
		if !equalSecret(stored, bind) {
			return denialMisbound
		}"""
new="""		delete(s.values, v) // consumed whether or not the binding matched
		bindingEnforced := false // MUTATION M-A: the binding layer is gone
		if bindingEnforced && !equalSecret(stored, bind) {
			return denialMisbound
		}"""
assert s.count(old)==1, s.count(old)
open(p,'w',encoding='utf-8',newline='').write(s.replace(old,new))
PY
run_one "MA" "$REPO/$APPR_REL::$MUTW/approval-MA.go"

# ---------- M-B: the row is issued under a digest that is NOT the item's own ----------
cp "$REPO/$QUEUE_REL" "$MUT/queue-MB.go"
python - "$MUT/queue-MB.go" <<'PY'
import sys
p=sys.argv[1]; s=open(p,encoding='utf-8').read()
old="	it.grants.issue(nonce, it.bind)"
new="	it.grants.issue(nonce, it.bind+\"m\") // MUTATION M-B: mint writes a foreign digest"
assert s.count(old)==1, s.count(old)
open(p,'w',encoding='utf-8',newline='').write(s.replace(old,new))
PY
run_one "MB" "$REPO/$QUEUE_REL::$MUTW/queue-MB.go"

# ---------- M-C: burn-order - a rejected spend does NOT consume the nonce ----------
cp "$REPO/$APPR_REL" "$MUT/approval-MC.go"
python - "$MUT/approval-MC.go" <<'PY'
import sys
p=sys.argv[1]; s=open(p,encoding='utf-8').read()
old="""		delete(s.values, v) // consumed whether or not the binding matched
		if !equalSecret(stored, bind) {
			return denialMisbound
		}"""
new="""		if !equalSecret(stored, bind) {
			return denialMisbound // MUTATION M-C: reject no longer burns the nonce
		}
		delete(s.values, v)"""
assert s.count(old)==1, s.count(old)
open(p,'w',encoding='utf-8',newline='').write(s.replace(old,new))
PY
run_one "MC" "$REPO/$APPR_REL::$MUTW/approval-MC.go"

# ---------- M-D: bindDigest stops folding the sequence number ----------
cp "$REPO/$APPR_REL" "$MUT/approval-MD.go"
python - "$MUT/approval-MD.go" <<'PY'
import sys
p=sys.argv[1]; s=open(p,encoding='utf-8').read()
old="	for _, s := range []string{corr, taskID, tool, level, strconv.FormatUint(seq, 10)} {"
new="	_ = seq // MUTATION M-D"
old2="	for _, s := range []string{corr, taskID, tool, level, strconv.FormatUint(seq, 10)} {\n\t\t_, _ = sum.Write([]byte{byte(len(s))})\n\t\t_, _ = sum.Write([]byte(s))\n\t}"
new2="	for _, s := range []string{corr, taskID, tool, level, \"0\"} { // MUTATION M-D: seq not folded in\n\t\t_, _ = sum.Write([]byte{byte(len(s))})\n\t\t_, _ = sum.Write([]byte(s))\n\t}\n\t_ = seq"
assert s.count(old2)==1, s.count(old2)
open(p,'w',encoding='utf-8',newline='').write(s.replace(old2,new2))
PY
run_one "MD" "$REPO/$APPR_REL::$MUTW/approval-MD.go"

# ---------- M-E: a second, unrelated card's store is consulted too ----------
# (cross-item reachability: spend falls through to denialNone when this item's
#  store has no row but the digest matches some other value)
cp "$REPO/$APPR_REL" "$MUT/approval-ME.go"
python - "$MUT/approval-ME.go" <<'PY'
import sys
p=sys.argv[1]; s=open(p,encoding='utf-8').read()
old="	return denialSpentNonce\n}\n\n// revoke drops every live nonce"
new="	return denialNone // MUTATION M-E: unknown nonce opens the card\n}\n\n// revoke drops every live nonce"
assert s.count(old)==1, s.count(old)
open(p,'w',encoding='utf-8',newline='').write(s.replace(old,new))
PY
run_one "ME" "$REPO/$APPR_REL::$MUTW/approval-ME.go"

echo "DONE rc=$?"
