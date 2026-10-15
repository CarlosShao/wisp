#!/bin/sh
# 120-r2 gate round 2: the AC#3 roster on its own口径, the uid-1000 permission
# faces (pre-fix vs post-fix on the same probe), and the control that neuters the
# identity gate. Every run is the same container image (golang:1.27) and the same
# fixture family. Logs land in .scratch/wisp/probes/120/r2/ and are kept.
set -e
REPO="D:/work/workspace/projects plans/Wisp"
OUT="$REPO/.scratch/wisp/probes/120/r2"
S=/d/tmp/probe120/r2
cd "$REPO"

echo "== building the fixed snapshot (HEAD bytes + this leg's two files) =="
rm -rf "$S/base-fixed"
mkdir -p "$S/base-fixed"
git archive 58de465a | tar -x -C "$S/base-fixed"
cp internal/winsec/winsec_other.go "$S/base-fixed/internal/winsec/winsec_other.go"
cp internal/winsec/posix_seal_window_120_other_test.go "$S/base-fixed/internal/winsec/posix_seal_window_120_other_test.go"
mkdir -p "$S/base-fixed/permface" "$S/base-head/permface"
cp "$S/permface/main.go" "$S/base-fixed/permface/main.go"
cp "$S/permface/main.go" "$S/base-head/permface/main.go"
grep -c 'pinPlacement' "$S/base-fixed/internal/winsec/winsec_other.go"
grep -c 'os.Chmod(path, mode)' "$S/base-head/internal/winsec/winsec_other.go"

echo "== AC#3 roster, ticket 113's own -run filter, on the fixed working tree =="
MSYS_NO_PATHCONV=1 docker run --rm -e CGO_ENABLED=0 -v "$REPO":/work -v probe120build:/root/.cache/go-build -w /work golang:1.27 \
  sh -c 'go test -count=1 -v -run "TestAC1POSIX|TestAC2POSIX|TestAC3POSIX|TestAC4POSIX" ./internal/winsec/ > /tmp/roster.log 2>&1; echo "roster rc=$?"; grep -c "^=== RUN" /tmp/roster.log; grep -c -- "--- PASS" /tmp/roster.log; grep -c -- "--- FAIL" /tmp/roster.log; grep -c -- "--- SKIP" /tmp/roster.log; sed -n "s/^\( *\)--- \(PASS\|FAIL\|SKIP\): \(Test[A-Za-z0-9_\/]*\).*/\2 \3/p" /tmp/roster.log | sort | uniq -c; cp /tmp/roster.log /work/.scratch/wisp/probes/120/r2/120-ac3-roster-113-names.log'

echo "== permission faces at uid 0, pre-fix code =="
RG="MSYS_NO_PATHCONV=1 docker run --rm -e CGO_ENABLED=0 -v $S/base-head:/work -v probe120build:/root/.cache/go-build -w /work golang:1.27"
$RG sh -c 'go run ./permface' > "$OUT/120-uid0-permface-prefix.txt" 2>&1 || true
cat "$OUT/120-uid0-permface-prefix.txt"

echo "== permission faces at uid 0, post-fix code =="
RG2="MSYS_NO_PATHCONV=1 docker run --rm -e CGO_ENABLED=0 -v $S/base-fixed:/work -v probe120build:/root/.cache/go-build -w /work golang:1.27"
$RG2 sh -c 'go run ./permface' > "$OUT/120-uid0-permface-postfix.txt" 2>&1 || true
cat "$OUT/120-uid0-permface-postfix.txt"

echo "== permission faces at uid 1000, pre-fix code =="
RG3="MSYS_NO_PATHCONV=1 docker run --rm --user 1000:1000 -e CGO_ENABLED=0 -e GOCACHE=/tmp/gocache -v $S/base-head:/work -w /work golang:1.27"
$RG3 sh -c 'go run ./permface' > "$OUT/120-uid1000-permface-prefix.txt" 2>&1 || true
cat "$OUT/120-uid1000-permface-prefix.txt"

echo "== permission faces at uid 1000, post-fix code =="
RG4="MSYS_NO_PATHCONV=1 docker run --rm --user 1000:1000 -e CGO_ENABLED=0 -e GOCACHE=/tmp/gocache -v $S/base-fixed:/work -w /work golang:1.27"
$RG4 sh -c 'go run ./permface' > "$OUT/120-uid1000-permface-postfix.txt" 2>&1 || true
cat "$OUT/120-uid1000-permface-postfix.txt"

echo "== whole package at uid 1000, pre-fix code =="
$RG3 sh -c 'id -u; go test -count=1 ./internal/winsec/' > "$OUT/120-uid1000-package-prefix.log" 2>&1 || true
tail -3 "$OUT/120-uid1000-package-prefix.log"

echo "== whole package at uid 1000, post-fix code =="
$RG4 sh -c 'id -u; go test -count=1 ./internal/winsec/' > "$OUT/120-uid1000-package-postfix.log" 2>&1 || true
tail -3 "$OUT/120-uid1000-package-postfix.log"
