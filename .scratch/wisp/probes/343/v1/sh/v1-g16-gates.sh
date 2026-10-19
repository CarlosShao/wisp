#!/bin/sh
# AC#4 门禁四数 —— 只在共享工作树跑，跑前 unset GOTMPDIR TMPDIR TEMP TMP（票 311 shape B GUARD）。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 90
unset GOTMPDIR TMPDIR TEMP TMP
G=/d/tmp/wisp/20261019/343v1/logs/gates
mkdir -p $G
echo "=== 0. env sanity (want all empty) ==="
printf 'GOTMPDIR=[%s] TMPDIR=[%s] TEMP=[%s] TMP=[%s]\n' "${GOTMPDIR:-}" "${TMPDIR:-}" "${TEMP:-}" "${TMP:-}"

echo "=== 1. gofumpt -l on the file this ticket touched (blacklist.go) ==="
/d/work/base/gopath/bin/gofumpt.exe -l internal/risk/blacklist.go > $G/gate1-gofumpt-mine.txt 2>&1
echo "rc=$? lines=$(wc -l < $G/gate1-gofumpt-mine.txt)"
cat $G/gate1-gofumpt-mine.txt
echo "-- registry: whole-repo gofumpt -l (existing CRLF dirty set, for the named register) --"
/d/work/base/gopath/bin/gofumpt.exe -l ./internal ./cmd ./tools > $G/gate1b-gofumpt-repo.txt 2>&1
echo "rc=$? count=$(wc -l < $G/gate1b-gofumpt-repo.txt)"
cat $G/gate1b-gofumpt-repo.txt

echo "=== 2. sh scripts/d22scan.sh ==="
sh scripts/d22scan.sh > $G/gate2-d22scan.txt 2>&1
echo "rc=$?"
tail -6 $G/gate2-d22scan.txt

echo "=== 3. sh scripts/check-path-length-budget.sh ==="
sh scripts/check-path-length-budget.sh > $G/gate3-length.txt 2>&1
echo "rc=$?"
tail -8 $G/gate3-length.txt

echo "=== 4. go vet ./internal/risk/ ==="
GOTMPW="$(cygpath -w /d/tmp/wisp/20261019/343v1/gotmp/vet)"
mkdir -p "$GOTMPW"
GOPATH='D:\work\base\gopath' GOTMPDIR="$GOTMPW" /d/work/base/go/bin/go.exe vet ./internal/risk/ > $G/gate4-vet.txt 2>&1
echo "rc=$?"
cat $G/gate4-vet.txt

echo "=== positive control for the banned-literal-char gauge: HANDOVER has many, my write surface has 0 ==="
printf 'HANDOVER_hits='; grep -c $'\xc2\xa7' docs/reports/HANDOVER.md || true
printf 'ticket343_hits='; grep -c $'\xc2\xa7' '.scratch/wisp/issues/343-blacklist-comment-names-aEquiv-a-symbol-that-does-not-exist-in-the-package.md' || true
printf 'my_tree_files_scanned='; find .scratch/wisp/probes/343/v1 -type f | wc -l
printf 'my_tree_char_hits='; grep -rl $'\xc2\xa7' .scratch/wisp/probes/343/v1 2>/dev/null | wc -l
echo GATES_DONE
