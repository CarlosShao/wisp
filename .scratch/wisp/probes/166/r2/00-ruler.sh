#!/bin/sh
# 166-r2 · the ruler for ticket 166 AC#6 (coin-flip instrument -> deterministic).
# One invocation = one round of the SAME command on the SAME tree. Every number it
# prints is a count of a NAMED pattern in the raw log it names, so any figure can
# be re-derived from that log instead of from memory.
#
# usage: sh 00-ruler.sh <label> <tree-dir> <count> <run-regex>
#   label      -> raw log name: /d/tmp/wisp-166r2/logs/<label>.log
#   tree-dir   -> absolute path of the tree being measured (the exported tree)
#   count      -> -count for go test (20 = the judgement rounds)
#   run-regex  -> -run pattern ('TestSessionList' = 2 named top-level cases)
#
# UNITS (named, so a number is never ambiguous):
#   RUN                    lines '^=== RUN'  = cases STARTED. RUN=0 = the round
#                          never ran = SURFACE FAILURE, not a reading.
#   PASS_TOP/FAIL_TOP/SKIP lines '^--- PASS' / '^--- FAIL' / '^--- SKIP'
#                          = TOP-LEVEL cases only.
#   SUBTEST_LINES          lines '^    --- ' = subtests. This package has none,
#                          so the two rulers coincide here (checked, not assumed).
#   PASS_<case>/FAIL_<case> TOP-LEVEL lines for that ONE named case.
#   *_fail_bodies          assertion BODY lines, matched by TEXT pattern, never
#                          by line number: this leg moves the assertion lines, and
#                          a 'manage_test.go:78' counter silently reads 0 on a tree
#                          that is really red (v1 of this ruler did exactly that --
#                          named in 03-teeth.txt; superseded here).
# surface_markers          walk / build failed / no tests to run / 0xc0000135 /
#                          connection refused / creating work dir  (any hit =
#                          classify as surface failure, do not count as a reading)
label=$1
tree=$2
n=$3
runre=$4
log="/d/tmp/wisp-166r2/logs/${label}.log"
# 5th arg = GOTMPDIR for this round. Default stays the D: root this leg started on,
# but ticket 311's guard (runtests.sh, called by scripts/d22scan.sh) says a test temp
# root on ANOTHER VOLUME than %LOCALAPPDATA% re-colours gates, so the volume-C pair in
# 08-headtree-c-volume.txt passes this argument explicitly. The ruler prints which root
# it used:台面名 must ride every round or two rounds are not comparable.
gotmp=${5:-D:/tmp/wisp-166r2/gotmp}

echo "label=${label}"
echo "tree=${tree}"
echo "count=${n}"
echo "run=${runre}"
echo "gotmpdir=${gotmp}"

if [ ! -d "$tree" ]; then
  echo "rc_cd=1 tree_missing=${tree}"
  exit 1
fi
mkdir -p /d/tmp/wisp-166r2/logs
mkdir -p "$gotmp"
echo "rc_mkdir_gotmp=$?"

cd "$tree" || { echo "rc_cd=1"; exit 1; }
echo "rc_cd=0 pwd=$(pwd)"
echo "go_files_in_tree=$(find . -name '*.go' | wc -l)"
echo "test_file_sha256_tree=$(sha256sum ./internal/session/manage_test.go | cut -c1-16)"

GOTMPDIR="$gotmp" go test ./internal/session/ -count="$n" -run "$runre" -v >"$log" 2>&1
rc=$?
echo "rc=${rc}"
echo "log=${log}"
echo "RUN=$(grep -c '^=== RUN' "$log")"
echo "PASS_TOP=$(grep -c '^--- PASS' "$log")"
echo "FAIL_TOP=$(grep -c '^--- FAIL' "$log")"
echo "SKIP=$(grep -c '^--- SKIP' "$log")"
echo "SUBTEST_LINES=$(grep -c '^    --- ' "$log")"
echo "PASS_ByWorkspace=$(grep -c '^--- PASS: TestSessionListByWorkspace' "$log")"
echo "FAIL_ByWorkspace=$(grep -c '^--- FAIL: TestSessionListByWorkspace' "$log")"
echo "PASS_ByTime=$(grep -c '^--- PASS: TestSessionListByTime' "$log")"
echo "FAIL_ByTime=$(grep -c '^--- FAIL: TestSessionListByTime' "$log")"
echo "PASS_SearchTitleAndBody=$(grep -c '^--- PASS: TestSessionSearchTitleAndBody' "$log")"
echo "FAIL_SearchTitleAndBody=$(grep -c '^--- FAIL: TestSessionSearchTitleAndBody' "$log")"
echo "pkg_fail_line=$(grep -cE '^FAIL[[:space:]]+github' "$log")"
echo "surface_markers=$(grep -cE 'walk|build failed|no tests to run|0xc0000135|connection refused|creating work dir' "$log")"
echo "ws_fail_bodies=$(grep -c 'workspace filter =' "$log")"
echo "ws_fail_bodies_at_lines=$(grep -oE 'manage_test\.go:[0-9]+: workspace filter' "$log" | sed -E 's/manage_test\.go:([0-9]+).*/\1/' | sort -n | uniq -c | tr '\n' ' ')"
echo "bytime_order_fail_bodies=$(grep -cE 'List\[[0-9]+\] = .*want .* \(order ' "$log")"
echo "first_fail_body=$(grep -m1 -E 'manage_test\.go:[0-9]+: (workspace filter|List\[)' "$log")"
