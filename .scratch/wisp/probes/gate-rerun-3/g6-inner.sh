#!/usr/bin/env bash
set -u
cd /work || { echo cd_FAIL; exit 9; }
go test -count=2 -v ./internal/winsec/ ./internal/memory/ ./internal/risk/ > /tmp/g6.log 2>&1
echo "g6_rc=$?"
echo "RUN=$(grep -cE '^=== RUN ' /tmp/g6.log) PASS=$(grep -cE '^--- PASS' /tmp/g6.log) FAIL=$(grep -cE '^--- FAIL' /tmp/g6.log) SKIP=$(grep -cE '^--- SKIP' /tmp/g6.log)"
echo "--- SKIP names ---"
grep -E '^--- SKIP' /tmp/g6.log | sed -E 's/.*SKIP: ([^ ]+).*/\1/' | sort -u
echo "--- FAIL names ---"
grep -E '^--- FAIL' /tmp/g6.log | sed -E 's/.*FAIL: ([^ ]+).*/\1/' | sort -u
echo "--- package lines ---"
grep -E '^(ok|FAIL)' /tmp/g6.log | head
