#!/bin/sh
# 编排者收件复跑 227-r1 第三步（改正射程）：决定性那一发＝在 matrixBase（manager_test.go:38）删掉 schema_version 行
set -u
REPO="/d/work/workspace/projects plans/Wisp"
W="/d/tmp/wisp/20261018/orch/227-verify"
cd "$W/tree" || exit 9
M=internal/config/manager_test.go
T=internal/config/writeguard_226_test.go
export GOTMPDIR="$W/gotmp" TMPDIR="$W/gotmp" TEMP="$W/gotmp" TMP="$W/gotmp"
{
echo "=== pristine 两枚 sha256 ==="
sha256sum "$M" "$T"
echo "=== [1] 正控：未突变时那枚用例必须绿 ==="
go test ./internal/config/ -run 'TestGuardedWriteDoesNotRewriteAFileThatAlreadySaysIt' -count=1 -timeout 30m -v 2>&1 | grep -E '^(--- (PASS|FAIL)|ok  |FAIL)' | head -4
echo "PRISTINE_RC=${PIPESTATUS[0]}"
echo "=== [2] 突变：删掉 matrixBase 里那行 schema_version = 2 ==="
sed -i '38d' "$M"
echo "manager_test.go:36-40 now:"; sed -n '36,40p' "$M"
echo "mutated sha256 = $(sha256sum "$M" | cut -d' ' -f1)"
go test ./internal/config/ -run 'TestGuardedWriteDoesNotRewriteAFileThatAlreadySaysIt' -count=1 -timeout 30m -v 2>&1 | grep -E '^(=== RUN|--- (PASS|FAIL)|ok  |FAIL|    writeguard)' | head -10
echo "MUT_RC=${PIPESTATUS[0]}"
echo "=== [3] 整包四数（突变态，看还有谁红） ==="
go test ./internal/config/ -count=1 -timeout 30m -v 2>&1 > "$W/logs/mut-pkg.log"
echo "RUN=$(grep -c '^=== RUN' "$W/logs/mut-pkg.log") PASS_TOP=$(grep -c '^--- PASS' "$W/logs/mut-pkg.log") FAIL_TOP=$(grep -c '^--- FAIL' "$W/logs/mut-pkg.log") SKIP=$(grep -c '^--- SKIP' "$W/logs/mut-pkg.log")"
grep '^--- FAIL' "$W/logs/mut-pkg.log" | head -8
echo "=== [4] 还原＋复验＋复跑绿 ==="
git -C "$REPO" show "bcf90e06:$M" > "$M"
echo "restored sha256 = $(sha256sum "$M" | cut -d' ' -f1)"
go test ./internal/config/ -run 'TestGuardedWriteDoesNotRewriteAFileThatAlreadySaysIt' -count=1 -timeout 30m -v 2>&1 | grep -E '^(--- (PASS|FAIL)|ok  |FAIL)' | head -4
echo "RESTORED_RC=${PIPESTATUS[0]}"
echo "=== [5] 交件树整包四数复跑（对照腿报 190/109/0/0） ==="
go test ./internal/config/ -count=1 -timeout 30m -v 2>&1 > "$W/logs/final-pkg.log"
echo "RUN=$(grep -c '^=== RUN' "$W/logs/final-pkg.log") PASS_TOP=$(grep -c '^--- PASS' "$W/logs/final-pkg.log") FAIL_TOP=$(grep -c '^--- FAIL' "$W/logs/final-pkg.log") SKIP=$(grep -c '^--- SKIP' "$W/logs/final-pkg.log")"
} 2>&1
