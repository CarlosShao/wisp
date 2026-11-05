#!/bin/bash
# 339-r2 门禁四数：只在共享工作树里落 rc（gofumpt 真身／d22scan／长度门／go vet ./internal/risk/）
# ⛔ ./cmd/wisp 整包、⛔ go test ./...；跑前 unset GOTMPDIR TMPDIR TEMP TMP（票 311 GUARD）。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
LOGD="$REPO/.scratch/wisp/probes/339/r2/logs"
FILES="internal/risk/anchor_miss_reason_334_test.go internal/risk/anchor_miss_reason_334_windows_test.go internal/risk/anchor_miss_reason_339_other_test.go"
cd "$REPO" || exit 9
unset GOTMPDIR TMPDIR TEMP TMP

echo "== 1) go build ./cmd/... ./internal/... ./tools/... =="
go build ./cmd/... ./internal/... ./tools/... > "$LOGD/gate-build.log" 2>&1; echo "BUILD_RC=$?"
head -5 "$LOGD/gate-build.log"

echo "== 2) go test ./internal/risk/ -count=1 -timeout 240s =="
go test ./internal/risk/ -count=1 -timeout 240s > "$LOGD/gate-test-risk.log" 2>&1; echo "TEST_RC=$?"
grep -E '^(ok|FAIL|---)' "$LOGD/gate-test-risk.log" | head -5

echo "== 3) gofumpt -l (写面三件) =="
/d/work/base/gopath/bin/gofumpt.exe -l $FILES > "$LOGD/gate-gofumpt.log" 2>&1; echo "GOFUMPT_RC=$?  LINES=$(wc -l < "$LOGD/gate-gofumpt.log")"
cat "$LOGD/gate-gofumpt.log"

echo "== 4) sh scripts/d22scan.sh =="
sh scripts/d22scan.sh > "$LOGD/gate-d22scan.log" 2>&1; echo "D22_RC=$?"
tail -3 "$LOGD/gate-d22scan.log"

echo "== 5) scripts/check-path-length-budget.sh =="
sh scripts/check-path-length-budget.sh > "$LOGD/gate-pathlen.log" 2>&1; echo "PATHLEN_RC=$?"
tail -3 "$LOGD/gate-pathlen.log"

echo "== 6) go vet ./internal/risk/ =="
go vet ./internal/risk/ > "$LOGD/gate-vet.log" 2>&1; echo "VET_RC=$?"
cat "$LOGD/gate-vet.log"
