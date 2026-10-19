#!/bin/sh
# 四数尺：只在 ^=== RUN / ^--- PASS / ^--- FAIL / ^--- SKIP 上取（⛔ 混用⛔ 相减）。
# 用法：sh 03-four-numbers.sh <log>
L="$1"
echo "FILE=$L"
echo "BYTES=$(wc -c < "$L")"
echo "RUN=$(grep -c '^=== RUN' "$L" || true)"
echo "PASS_TOP=$(grep -c '^--- PASS' "$L" || true)"
echo "FAIL_TOP=$(grep -c '^--- FAIL' "$L" || true)"
echo "SKIP=$(grep -c '^--- SKIP' "$L" || true)"
echo "PKGFAIL=$(grep -c '^FAIL\s\+github' "$L" || true)"
echo "Panic=$(grep -c '^panic' "$L" || true)"
echo "RCline=$(grep '^RC=' "$L" || true)"
echo '--- red roster (top-level only) ---'
grep '^--- FAIL' "$L" | sed 's/^--- FAIL: //; s/ .*//' | sed 's/^/RED /'
