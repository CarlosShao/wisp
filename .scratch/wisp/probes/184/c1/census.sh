#!/usr/bin/env bash
# 184-c1 read-only census. Every internal/** line number is taken from the HEAD
# tree (git show HEAD:.. / git grep .. HEAD) so it is immune to a concurrent
# 183-v1 mutation of internal/risk/** in the working tree. Zero production code
# is written; this only prints measurements.
set +e
cd "$(dirname "$0")/../../../.." || exit 1   # logdir-relative to repo root, NOT runtime.Caller
echo "184-c1 census  date=$(date '+%Y-%m-%d %H:%M:%S')  HEAD=$(git rev-parse HEAD)  branch=$(git rev-parse --abbrev-ref HEAD)"
echo "git status --porcelain -- internal/ cmd/ :"; git status --porcelain -- internal/ cmd/; echo "[end status]"

r(){ echo; echo "### CMD: $*"; eval "$@"; }

echo "==================== 8 RULERS (re-measured vs ticket) ===================="
r "git show HEAD:internal/agent/journal.go | grep -n -B1 -A8 'Decision values'"
r "git show HEAD:internal/memory/models.go | grep -n -A9 'toolCallDecisions'"
r "git show HEAD:internal/memory/models.go | grep -n -A2 'invalid tool_call.decision'"
r "git show HEAD:internal/memory/dao_toolcall.go | grep -n -A4 'func (s \*Store) DecideToolCall'"
r "git show HEAD:docs/specs/SPEC-02-data-storage.md | nl -ba | sed -n '78,80p'"
r "git show HEAD:docs/PLAN.md | grep -n '取证记录'"
r "git grep -n -I 'GrantID:' HEAD -- internal/ cmd/ ; echo rc=\$?"
r "git grep -n -I 'DecideToolCall(' HEAD -- internal/ cmd/"
r "git show HEAD:internal/memory/privacy.go | grep -n -A2 'if r.Decision != \"\"'"
r "git show HEAD:internal/tools/bridge.go | grep -n 'dec.DecisionColumn = agent.DecisionAllow'"
r "git show HEAD:internal/agent/loop.go | grep -n 'j.decide(ctx, rowID, DecisionAllow)'"
r "git show HEAD:internal/agent/journal.go | grep -n -A7 'func riskColumn'"
r "git show HEAD:internal/agent/journal.go | grep -n 'riskColumn('"
r "git grep -n 'PassThroughUnclassifiedRisk' HEAD -- internal/ cmd/ | grep -v '_test.go'"
r "git grep -n 'PassThroughUnclassifiedRisk: true' HEAD"

echo "==================== AC#1 five booking sites + neighbours ===================="
r "git show HEAD:internal/tools/bridge.go | nl -ba | sed -n '371,416p'"
r "git show HEAD:internal/agent/loop.go | nl -ba | sed -n '783,806p'"
r "git show HEAD:internal/tools/bridge.go | nl -ba | sed -n '943,965p'"

echo "==================== AC#2 decision-COLUMN readers (not the tools.Decision struct) ===================="
r "git grep -nE 'if r.Decision|\\+ r.Decision|r\\.Decision !=|got\\.Decision|\\.Decision = decision' HEAD -- internal/ cmd/"
r "git grep -nE 'tools\\.Decision|risk\\.Decision' HEAD -- cmd/ internal/panel/ internal/tools/bridge.go internal/agent/approval/ | head"

echo "==================== AC#4 cost-table inputs ===================="
r "git show HEAD:internal/memory/schema.go | grep -n -A3 'Enum-shaped columns'"
r "git show HEAD:internal/memory/schema.go | grep -n 'decision '"
r "git show HEAD:internal/memory/schema.go | grep -nE 'migrationChain|{from:|SchemaVersionTarget'"
r "git grep -n 'SchemaVersionTarget' HEAD -- internal/memory/schema.go"
r "git grep -nE 'len\\(.*[Dd]ecision\\)|== 5|five decision' HEAD -- '*_test.go' ; echo rc=\$?"
r "git ls-files HEAD 2>/dev/null | grep -iE 'migration' ; echo rc=\$? ; echo \"[migration dir search done]\""

echo "==================== AC#5b D45 approval_grant PRODUCERS (non-test) ===================="
r "git grep -n 'InsertGrant(' HEAD -- internal/ cmd/ | grep -v '_test.go' ; echo rc=\$?"

echo "==================== AC#6 decision VALUES passed to j.decide (out-of-map risk) ===================="
r "git grep -n '\\.decide(ctx\\|\\.decide(wctx' HEAD -- internal/agent/loop.go"

echo "[census done]"
