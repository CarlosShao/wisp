#!/bin/sh
# 编排者收件复跑（227-r1 决定性那一发）：① 交件树五枚用例必须绿 ② 把 :393 那枚的 fixture 里 schema_version 删掉 ⇒ 它必须红
set -u
REPO="/d/work/workspace/projects plans/Wisp"
W="/d/tmp/wisp/20261018/orch/227-verify"
ANCHOR=bcf90e06
mkdir -p "$W/tree" "$W/gotmp" "$W/logs"
cd "$REPO" || exit 9
{
echo "### 编排者复跑 227-r1 anchor=$ANCHOR"
git archive "$ANCHOR" go.mod go.sum internal cmd tools scripts frontend models docs | tar -x -C "$W/tree"
echo "go files = $(find "$W/tree" -name '*.go' | wc -l)"
cd "$W/tree" || exit 9
export GOTMPDIR="$W/gotmp" TMPDIR="$W/gotmp" TEMP="$W/gotmp" TMP="$W/gotmp"
ROST='TestHandDeletedScalarLineComesBackAtSchemaDefault|TestHandDeletedMapEntryIsNotRestored|TestNothingWrittenPromiseOnlyHoldsWhenSchemaIsCurrent|TestNothingWrittenLeavesFileUntouchedWhenSchemaIsCurrent|TestDeletedLineRestorationIsNotClaimedAsAWrittenKey'
echo "=== [1] 五枚新用例在未改动交件树上必须全绿（证⛔ 恒红） ==="
go test ./internal/config/ -run "$ROST" -count=1 -timeout 30m -v 2>&1 | grep -E '^(=== RUN|--- (PASS|FAIL)|ok  |FAIL)' | head -14
echo "RC=${PIPESTATUS[0]}"
echo "=== [2] 整包成对四数（交件树） ==="
go test ./internal/config/ -count=1 -timeout 30m -v 2>&1 > "$W/logs/pkg.log"
echo "RUN=$(grep -c '^=== RUN' "$W/logs/pkg.log") PASS_TOP=$(grep -c '^--- PASS' "$W/logs/pkg.log") FAIL_TOP=$(grep -c '^--- FAIL' "$W/logs/pkg.log") SKIP=$(grep -c '^--- SKIP' "$W/logs/pkg.log")"
echo "=== [3] 决定性那一发：:393 那枚用例的 fixture 里删掉 schema_version 行 ⇒ 它必须红 ==="
grep -n 'func TestGuardedWriteDoesNotRewriteAFileThatAlreadySaysIt' internal/config/writeguard_226_test.go
sha256sum internal/config/writeguard_226_test.go | cut -d' ' -f1
awk 'NR>=393 && NR<=410' internal/config/writeguard_226_test.go | grep -n schema_version
} > "$W/logs/orch-227-recheck.log" 2>&1
echo "LOG=$W/logs/orch-227-recheck.log"
cat "$W/logs/orch-227-recheck.log"
