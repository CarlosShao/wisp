#!/usr/bin/env bash
# 341-v1 AC#0 + AC#1 landing-distinctness + AC#2 roster-intersection gauges.
R="/d/work/workspace/projects plans/Wisp"
W=/d/tmp/wisp/20261018/341v1
LOGROOT="$R/.scratch/wisp/probes/341/v1/logs"
cd "$R" || exit 9

echo "######## AC#0 independent census: which shapes has ANY *_test.go ever asserted? ########"
echo "ruler = grep over the tree's *_test.go for the shape's own planted-delete content anchor"
echo "        (a shape counts as NAILED only when the test both PLANTS the delete and ASSERTS a result);"
echo "        hits printed as file:line. Run twice: base tree (pre-341) and delivery tree."

census() {
  tree="$1"; label="$2"
  echo "===== $label ($tree) ====="
  for d in internal/config cmd/wisp; do
    echo "-- dir scope: $d --"
    echo "P1 plant-delete a NON-default scalar line (removes the 'size = 60' line):"
    grep -rnE 'size = 60.*\\n", "", 1' --include='*_test.go' "$tree/$d" 2>/dev/null | sed "s|$tree/||" | sed 's/^/   /'
    echo "P2 plant-delete a DEFAULT-valued scalar line (removes the 'size = 56' line):"
    grep -rnE 'size = 56.*\\n", "", 1' --include='*_test.go' "$tree/$d" 2>/dev/null | sed "s|$tree/||" | sed 's/^/   /'
    echo "P3 plant-delete a WHOLE SECTION (header [ball] plus its body in one removal):"
    grep -rnE '"\[ball\]\\n.*\\n", "", 1' --include='*_test.go' "$tree/$d" 2>/dev/null | sed "s|$tree/||" | sed 's/^/   /'
    echo "P4 plant-delete a MAP ENTRY (provider sub-table removed):"
    grep -rnE '"\\n\[llm\]\\n' --include='*_test.go' "$tree/$d" 2>/dev/null | sed "s|$tree/||" | sed 's/^/   /'
    echo "P5/P5-prime plant-delete the schema_version line:"
    grep -rnE 'schema_version = 2\\n", "", 1' --include='*_test.go' "$tree/$d" 2>/dev/null | sed "s|$tree/||" | sed 's/^/   /'
    echo "P6 plant-delete a provider base_url LINE inside a surviving table:"
    grep -rnE 'base_url = .*\\n", "", 1' --include='*_test.go' "$tree/$d" 2>/dev/null | sed "s|$tree/||" | sed 's/^/   /'
  done
}
census "$W/treeB" "BASE tree (anchor 6948eb08^, pre-341)"
census "$W/treeD" "DELIVERY tree (anchor HEAD)"

echo
echo "######## AC#1: the four de-gum landings, pairwise distinct lines (现验) ########"
for spec in "writeguard.go:174|175|p2 INFO-branch delete-shape clause" "writeguard.go:141|141|p3 base = deepCopyConfig(disk)" "loader.go:257|257|p5 cp.SchemaVersion = SchemaVersionCurrent" "loader.go:213|213|p6 p.BaseURL = pre.BaseURL"; do
  f=${spec%%:*}; rest=${spec#*:}; a=${rest%%|*}; b=$(echo "$rest" | cut -d'|' -f2 | cut -d' ' -f1); nm=$(echo "$rest" | cut -d'|' -f2-)
  echo "--- $nm :: internal/config/$f lines $a-$b ---"
  awk -v s="$a" -v e="$b" 'NR>=s && NR<=e {printf "   %d: %s\n", NR, $0}' "$W/treeD/internal/config/$f"
done
echo "-- distinctness check (byte offsets of each anchor in its file, all four different) --"
perl -0777 -ne 'while (/are not preserved; a key line the\s*"\+\s*"operator deleted by hand comes back at its schema default value/gs) { printf "p2 offset=%d\n", $-[0] }' "$W/treeD/internal/config/writeguard.go"
perl -0777 -ne 'while (/base = deepCopyConfig\(disk\)/gs) { printf "p3 offset=%d\n", $-[0] }' "$W/treeD/internal/config/writeguard.go"
perl -0777 -ne 'while (/cp\.SchemaVersion = SchemaVersionCurrent/gs) { printf "p5 offset=%d\n", $-[0] }' "$W/treeD/internal/config/loader.go"
perl -0777 -ne 'while (/p\.BaseURL = pre\.BaseURL/gs) { printf "p6 offset=%d\n", $-[0] }' "$W/treeD/internal/config/loader.go"
echo "-- line numbers (from grep -n) --"
grep -n 'base = deepCopyConfig(disk)' "$W/treeD/internal/config/writeguard.go"
grep -n 'operator deleted by hand comes back at its schema default value' "$W/treeD/internal/config/writeguard.go"
grep -n 'cp\.SchemaVersion = SchemaVersionCurrent' "$W/treeD/internal/config/loader.go"
grep -n 'p\.BaseURL = pre\.BaseURL' "$W/treeD/internal/config/loader.go"

echo
echo "######## AC#2: red rosters of the four shots, named, and the P2 <-> P3 intersection ########"
L="$LOGROOT"
for x in p2 p3 p5 p6; do echo "--- $x roster ---"; sort "$W/log-$x/rednames.txt"; done
echo "-- intersection P2 ∩ P3 (want: does P2 appear in P3's roster? want NO) --"
grep -c "TestHandDeletedDefaultValuedLineIsRestoredInvisibly" "$W/log-p3/rednames.txt" || true
echo "-- intersection as a set difference --"
comm -12 <(sort "$W/log-p2/rednames.txt") <(sort "$W/log-p3/rednames.txt") | sed 's/^/   COMMON: /'
echo "-- P2 shot roster vs P3 shot roster (symdiff) --"
comm -3 <(sort "$W/log-p2/rednames.txt") <(sort "$W/log-p3/rednames.txt") | sed 's/^/   /'

echo
echo "######## AC#4: base<->delivery red roster difference (want zero new red) ########"
echo "base FAIL_TOP=$(grep -c '^--- FAIL' $W/run-base-A/test.log || true) delivery FAIL_TOP=$(grep -c '^--- FAIL' $W/run-deliv-A/test.log || true)"
grep '^--- FAIL' "$W/run-base-A/test.log" | sed 's/^--- FAIL: //' | cut -d' ' -f1 | sort > "$W/base-red.txt"
grep '^--- FAIL' "$W/run-deliv-A/test.log" | sed 's/^--- FAIL: //' | cut -d' ' -f1 | sort > "$W/deliv-red.txt"
echo "-- new red in delivery (want empty) --"; comm -13 "$W/base-red.txt" "$W/deliv-red.txt"
echo "-- five new case names present in delivery RUN list --"
for t in TestHandDeletedDefaultValuedLineIsRestoredInvisibly TestHandDeletedWholeSectionComesBackFromDefaultsNotMemory TestHandDeletedSchemaVersionOnTheWriteBranchIsForcedBackAndUnclaimed TestHandDeletedProviderBaseUrlComesBackAsThePresetValue TestHandDeletedBaseUrlOfAnUnknownProviderStaysEmptyInTheFile; do
  printf "%s = %s\n" "$t" "$(grep -cE "^--- RUN $t\$|^=== RUN   $t\$" "$W/run-deliv-A/test.log" || true)"
done
echo "-- PASS lines for the five on the unmutated delivery tree --"
for t in TestHandDeletedDefaultValuedLineIsRestoredInvisibly TestHandDeletedWholeSectionComesBackFromDefaultsNotMemory TestHandDeletedSchemaVersionOnTheWriteBranchIsForcedBackAndUnclaimed TestHandDeletedProviderBaseUrlComesBackAsThePresetValue TestHandDeletedBaseUrlOfAnUnknownProviderStaysEmptyInTheFile; do
  printf "%s = %s\n" "$t" "$(grep -E "^--- PASS: $t( |\()" "$W/run-deliv-A/test.log" | head -1)"
done
