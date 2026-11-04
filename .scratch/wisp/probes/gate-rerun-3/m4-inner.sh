#!/usr/bin/env bash
# M4 (票113:54) inner script — runs INSIDE golang:1.27 container, /work = snapshot of HEAD
set -u
cd /work || { echo "cd_work_FAIL"; exit 9; }
echo "=== false-green guard ==="
[ -f go.mod ] && echo "go_mod_present=yes" || echo "go_mod_present=NO"
echo "go_files_internal=$(find /work/internal -name '*.go' | wc -l)"
echo "posix_winsec_other=$(find /work/internal/winsec -name 'winsec_other*.go' | wc -l)"
go version
echo "=== pristine positive control: go test ./internal/winsec/ -count=1 ==="
go test ./internal/winsec/ -count=1 2>&1 | tail -4
PC=$?
echo "pristine_tail_done"
echo "=== M4a: turn the new leg off (ancestorIsLink pass-through) ==="
grep -n 'ancestorIsLink(prefix)' internal/winsec/winsec_other.go | head
cp -p internal/winsec/winsec_other.go /tmp/winsec_other.pristine
sed -i 's/if ancestorIsLink(prefix) {/if false \&\& ancestorIsLink(prefix) {/' internal/winsec/winsec_other.go
echo "landing_m4a=$(grep -c 'if false && ancestorIsLink(prefix)' internal/winsec/winsec_other.go)"
go build ./internal/winsec/ 2>&1 | head -3; echo "build_m4a_rc=$?"
go test ./internal/winsec/ -count=1 -v > /tmp/m4a.log 2>&1; echo "m4a_test_rc=$?"
echo "m4a_fail_top=$(grep -cE '^--- FAIL' /tmp/m4a.log)"
grep -E '^--- FAIL' /tmp/m4a.log | sed -E 's/.*FAIL: ([^ ]+).*/\1/' | sort -u | head -25
echo "=== restore after M4a ==="
cp -p /tmp/winsec_other.pristine internal/winsec/winsec_other.go
diff -q /tmp/winsec_other.pristine internal/winsec/winsec_other.go && echo "restored_m4a_identical=yes"
echo "=== M4b: half-fix, consult only first prefix ==="
sed -i 's/for _, prefix := range pathPieces(path) {/for _, prefix := range pathPieces(path)[:1] {/' internal/winsec/winsec_other.go
echo "landing_m4b=$(grep -c 'pathPieces(path)\[:1\]' internal/winsec/winsec_other.go)"
go build ./internal/winsec/ 2>&1 | head -3; echo "build_m4b_rc=$?"
go test ./internal/winsec/ -count=1 -v > /tmp/m4b.log 2>&1; echo "m4b_test_rc=$?"
echo "m4b_fail_top=$(grep -cE '^--- FAIL' /tmp/m4b.log)"
grep -E '^--- FAIL' /tmp/m4b.log | sed -E 's/.*FAIL: ([^ ]+).*/\1/' | sort -u | head -25
echo "=== restore after M4b ==="
cp -p /tmp/winsec_other.pristine internal/winsec/winsec_other.go
diff -q /tmp/winsec_other.pristine internal/winsec/winsec_other.go && echo "restored_m4b_identical=yes"
echo "DONE_M4"
