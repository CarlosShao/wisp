#!/usr/bin/env bash
# ticket 158 acceptance r1 - AC#1 load-bearing rig.
#
# Everything runs through `go test -overlay`: the shared worktree is never
# written, and r1's committed guard test is never deleted on disk - it is only
# replaced for the compiler. Never combined with -cover* (dispatch 5: the
# overlay is silently ignored when both are used).
#
# Mutants (built from `git cat-file blob c7f638c:internal/tools/bridge.go`, not
# from the worktree, so the input is pinned to the audited revision):
#   M1 = mark(): delete b.OpenTask(dec.TaskID)          (the reproduced P-2)
#   M2 = mark(): widen the guard so EVERY task opens     (r1's own degenerate shape)
#   M3 = OpenTask(): delete b.prov.OpenScope(taskID)      (bridge ledger yes, risk layer no)
# Guard variants of internal/tools/bridge_scope_open_ticket158_test.go:
#   stub  = package with the test absent (simulates "not yet written", AC#1 发①)
#   no1 / no2 / no3 = one 判据 removed each
set -eu
cd "$(git rev-parse --show-toplevel)" || exit 9
A=c7f638c
P="$PWD/.scratch/wisp/probes/158/accept-r1"
RW='D:/work/workspace/projects plans/Wisp'
BRIDGE_KEY="$RW/internal/tools/bridge.go"
GUARD_KEY="$RW/internal/tools/bridge_scope_open_ticket158_test.go"

mkdir -p "$P/mut" "$P/ovl"
git cat-file blob "$A:internal/tools/bridge.go" > "$P/mut/bridge.as-is.go"
git cat-file blob "$A:internal/tools/bridge_scope_open_ticket158_test.go" > "$P/mut/guard.as-is.go"

# --- mutants ---------------------------------------------------------------
cp "$P/mut/bridge.as-is.go" "$P/mut/bridge.M1.go"; sed -i '559d' "$P/mut/bridge.M1.go"
cp "$P/mut/bridge.as-is.go" "$P/mut/bridge.M2.go"; sed -i '552s#.*#\tif b.prov == nil || dec.TaskID == "" {#' "$P/mut/bridge.M2.go"
cp "$P/mut/bridge.as-is.go" "$P/mut/bridge.M3.go"; sed -i '642s/b\.prov\.OpenScope(taskID)/_ = open/' "$P/mut/bridge.M3.go"

# --- guard variants --------------------------------------------------------
printf 'package tools\n' > "$P/mut/guard.stub.go"
cp "$P/mut/guard.as-is.go" "$P/mut/guard.no1.go"; sed -i '110,115d' "$P/mut/guard.no1.go"
cp "$P/mut/guard.as-is.go" "$P/mut/guard.no2.go"; sed -i '117,135d' "$P/mut/guard.no2.go"
cp "$P/mut/guard.as-is.go" "$P/mut/guard.no3.go"; sed -i '85,97d' "$P/mut/guard.no3.go"

echo "## mutation diffs vs the pinned anchor copy (each must be one hunk)"
for m in M1 M2 M3; do echo "--- bridge $m"; diff "$P/mut/bridge.as-is.go" "$P/mut/bridge.$m.go" || true; done
for g in no1 no2 no3; do echo "--- guard $g: deleted line count = $(diff "$P/mut/guard.as-is.go" "$P/mut/guard.$g.go" | grep -c '^<')"; done

# --- overlays (python writes the JSON; no hand-rolled comma logic) ---------
ovl() { # $1 name $2 bridge-suffix (NONE allowed) $3 guard-suffix (NONE allowed)
	python - "$P/ovl/$1.json" "$BRIDGE_KEY" "$GUARD_KEY" "$RW" "$2" "$3" <<'PY'
import json, sys
out, bk, gk, rw, b, g = sys.argv[1:7]
rep = {}
if b != "NONE":
    rep[bk] = "%s/.scratch/wisp/probes/158/accept-r1/mut/bridge.%s.go" % (rw, b)
if g != "NONE":
    rep[gk] = "%s/.scratch/wisp/probes/158/accept-r1/mut/guard.%s.go" % (rw, g)
json.dump({"Replace": rep}, open(out, "w"), indent=2)
print(out)
PY
}
ovl r1-stubmut  M1  stub
ovl r2-stub     NONE stub
ovl r3-m2full   M2  as-is
ovl r4-m2no3    M2  no3
ovl r5-m1no1    M1  no1
ovl r6-m1no2    M1  no2
ovl r7-m3full   M3  as-is
echo "## overlay contents"
for f in "$P"/ovl/*.json; do echo "--- $f"; cat "$f"; done
