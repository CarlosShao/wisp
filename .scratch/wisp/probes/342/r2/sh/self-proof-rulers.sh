#!/bin/sh
# 342-r2 self-proof rulers, SCOPED TO THIS LEG'S OWN COMMITS (the anchor..HEAD range contains other
# legs' commits: 227-v2 / 346-a1 / 296-r2 / orchestrator records, so an unscoped ruler would credit
# their bytes to this leg - A1045's family, this time on the write ruler).
cd "/d/work/workspace/projects plans/Wisp" || exit 1
ANCHOR=3abbb6391a8a497c0332e5363f24c8a250a7eb78   # this leg's first git log -1, taken at 21:0x
MYCOMMITS="4f4f8419 a782196e 2626119a 4ff76a0c 0e487f06"

echo "== R0: every commit between the anchor and HEAD, attributed by subject (none of them is mine"
echo "   unless it starts with 票 342＋票 348 落地腿 342-r2) =="
git log --format='%h %s' "$ANCHOR"..HEAD | cut -c1-72

echo
echo "== R1: production (non-_test) paths touched by MY commits =="
for c in $MYCOMMITS; do
  git show --name-only --format= "$c" | grep -v '_test\.go$' | sed "s|^|PROD-Touched by $c: |"
done
echo "R1_PROD_PATHS_TOTAL=$(for c in $MYCOMMITS; do git show --name-only --format= $c; done | grep -vc '_test\.go$')"

echo "== R1b: winsec.go (wrapPath :430 + the :118-140 statement) inside my own diff =="
echo "R1B_WINSEC_GO_LINES=$(for c in $MYCOMMITS; do git show --name-only --format= $c; done | grep -c 'winsec\.go$')"

echo "== R2: unique .go paths this leg wrote, and the declared write face =="
for c in $MYCOMMITS; do git show --name-only --format= $c; done | sort -u | sed 's/^/  /'

echo "== R3: capital-initial top-level additions in MY diff =="
for c in $MYCOMMITS; do git show --format= --unified=0 "$c"; done | grep -E '^\+\s*(func|var|type|const) [A-Z]' | sed 's/^/  HIT: /'
echo "R3_CAPITAL=$(for c in $MYCOMMITS; do git show --format= --unified=0 $c; done | grep -cE '^\+\s*(func|var|type|const) [A-Z]')"
echo "R3B_STRICT_EXCLUDING_TEST_ENTRYPOINTS=$(for c in $MYCOMMITS; do git show --format= --unified=0 $c; done | grep -E '^\+\s*(func|var|type|const) [A-Z]' | grep -vcE '\+func Test')"

echo "== R4: function-shape changes in production files (numstat over non-_test paths) =="
SHAPE=0
for c in $MYCOMMITS; do
  n=$(git show --numstat --format= "$c" | grep -v '_test\.go$' | wc -l)
  SHAPE=$((SHAPE + n))
done
echo "R4_PROD_NUMSTAT_LINES=$SHAPE"

echo "== R5: write-face mutex - did ANY other commit in the range touch internal/winsec? =="
git log --format='%h %s' "$ANCHOR"..HEAD -- internal/winsec | cut -c1-72 | sed 's/^/  WINSEC-COMMITS: /'

echo "== R6: line counts of my two instrument files (照实报, the 160-line cap governs REPORT files only) =="
wc -l internal/winsec/ancestorwindow_342r2_windows_test.go internal/winsec/leafrace_342r2_windows_test.go

echo "== R7: byte-level placeholders / zero-byte guards in my probe root =="
find .scratch/wisp/probes/342/r2 -type f -size 0 | sed 's/^/  ZEROBYTE: /'
find .scratch/wisp/probes/342/r2 -name '*.out' | sed 's/^/  OUT: /'
echo "R7_FILES=$(find .scratch/wisp/probes/342/r2 -type f | wc -l)"
