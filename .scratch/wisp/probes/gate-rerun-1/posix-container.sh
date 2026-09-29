set -u
# Runs inside golang:1.27 with the throwaway snapshot mounted at /work.
# /work is a `git archive` copy OUTSIDE the repository, so editing files here is
# the shape tickets 113/104 prescribe ("变异只在 /tmp 仓外快照做、还原证 diff -q 干净").
echo "=== false-green guard: the mount must actually carry the tree ==="
test -f /work/go.mod || { echo "ABORT-EMPTY-MOUNT no go.mod"; exit 90; }
echo "go_files_under_internal=$(find /work/internal -name '*.go' | wc -l)"
echo "posix_test_files=$(ls /work/internal/winsec/*_other_test.go 2>/dev/null | wc -l)"
echo "winsec_go=$(ls /work/internal/winsec/*.go | wc -l)"
go version
echo "GOOS=$(go env GOOS) GOFLAGS=$(go env GOFLAGS)"
export GOPROXY=off
cd /work

P=internal/winsec/winsec_other.go
cp -p "$P" /tmp/winsec_other.pristine.go
echo -n "pristine_bytes=" && wc -c < "$P"
grep -n 'if ancestorIsLink(prefix) {' "$P"
grep -n 'for _, prefix := range pathPieces(path) {' "$P"

echo
echo "=== G-6 baseline: three packages -count=2 -v on linux ==="
go test ./internal/winsec/ ./internal/memory/ ./internal/risk/ -count=2 -v -p 1 > /tmp/g6-baseline.txt 2>&1
echo "g6_rc=$?"
grep -P '^(ok|FAIL)\t' /tmp/g6-baseline.txt
echo "four numbers counted WITH subtests:"
echo "RUN=$(grep -c '^=== RUN' /tmp/g6-baseline.txt) PASS=$(grep -c '^ *--- PASS' /tmp/g6-baseline.txt) FAIL=$(grep -c '^ *--- FAIL' /tmp/g6-baseline.txt) SKIP=$(grep -c '^ *--- SKIP' /tmp/g6-baseline.txt)"
echo "same log counted top-level only (the naive form): PASS=$(grep -c '^--- PASS' /tmp/g6-baseline.txt) FAIL=$(grep -c '^--- FAIL' /tmp/g6-baseline.txt) SKIP=$(grep -c '^--- SKIP' /tmp/g6-baseline.txt)"
grep '^ *--- FAIL' /tmp/g6-baseline.txt | sed 's/.*FAIL: //;s/ .*//' | sort -u > /tmp/g6-red-names.txt
echo "red_name_lines=$(wc -l < /tmp/g6-red-names.txt)"; cat /tmp/g6-red-names.txt
echo "SKIP names:"; grep '^ *--- SKIP' /tmp/g6-baseline.txt | sed 's/.*SKIP: //;s/ .*//' | sort -u

echo
echo "=== M-4a: turn the new POSIX link leg OFF ==="
sed -i 's/if ancestorIsLink(prefix) {/if false \&\& ancestorIsLink(prefix) {/' "$P"
echo -n "landing_hits=" && grep -c 'if false && ancestorIsLink(prefix) {' "$P"
go build ./internal/winsec/; echo "build_rc=$?"
go test ./internal/winsec/ -count=1 -v -p 1 > /tmp/m4a-red.txt 2>&1; echo "m4a_test_rc=$?"
grep -P '^(ok|FAIL)\t' /tmp/m4a-red.txt
grep '^ *--- FAIL' /tmp/m4a-red.txt | sed 's/.*FAIL: //;s/ .*//' | grep -v '/' | sort -u > /tmp/m4a-red-names.txt
echo "m4a_top_level_reds=$(wc -l < /tmp/m4a-red-names.txt)"; cat /tmp/m4a-red-names.txt
cp -p /tmp/winsec_other.pristine.go "$P"; diff -q "$P" /tmp/winsec_other.pristine.go && echo "restored=identical"

echo
echo "=== M-4b: half-fix, the walk consults ONE level (pathPieces(path)[:1]) ==="
sed -i 's/for _, prefix := range pathPieces(path) {/for _, prefix := range pathPieces(path)[:1] {/' "$P"
echo -n "landing_hits=" && grep -c 'range pathPieces(path)\[:1\]' "$P"
go build ./internal/winsec/; echo "build_rc=$?"
go test ./internal/winsec/ -count=1 -v -p 1 > /tmp/m4b-red.txt 2>&1; echo "m4b_test_rc=$?"
grep -P '^(ok|FAIL)\t' /tmp/m4b-red.txt
grep '^ *--- FAIL' /tmp/m4b-red.txt | sed 's/.*FAIL: //;s/ .*//' | grep -v '/' | sort -u > /tmp/m4b-red-names.txt
echo "m4b_top_level_reds=$(wc -l < /tmp/m4b-red-names.txt)"; cat /tmp/m4b-red-names.txt
cp -p /tmp/winsec_other.pristine.go "$P"; diff -q "$P" /tmp/winsec_other.pristine.go && echo "restored=identical"

echo
echo "=== G-6 gates (per package) ==="
for pkg in ./internal/winsec/ ./internal/memory/ ./internal/risk/; do
  go vet "$pkg"; echo "vet_rc[$pkg]=$?"
done
GOOS=linux go vet ./internal/winsec/ ./internal/memory/ ./internal/risk/; echo "GOOS_linux_vet_rc=$?"
echo -n "gofmt_l_count=" && gofmt -l internal/ 2>/dev/null | wc -l
gofmt -l internal/ 2>/dev/null | head -5
sh scripts/d22scan.sh > /tmp/g6-d22scan.txt 2>&1; echo "d22scan_rc=$?"
tail -30 /tmp/g6-d22scan.txt
echo "=== restore proof: the snapshot file must equal what we mounted ==="
diff -q "$P" /tmp/winsec_other.pristine.go && echo "winsec_other_identical_to_pristine=yes"
echo "=== DONE ==="
