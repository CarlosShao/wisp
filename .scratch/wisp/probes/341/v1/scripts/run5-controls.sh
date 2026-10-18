#!/usr/bin/env bash
# 341-v1 AC#1 controls: per-name verdicts for each de-gum shot + the unmutated
# positive-control PASS lines on the delivery tree.
W=/d/tmp/wisp/20261018/341v1
L="/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/341/v1/logs"

CASES="TestHandDeletedDefaultValuedLineIsRestoredInvisibly
TestHandDeletedWholeSectionComesBackFromDefaultsNotMemory
TestHandDeletedSchemaVersionOnTheWriteBranchIsForcedBackAndUnclaimed
TestHandDeletedProviderBaseUrlComesBackAsThePresetValue
TestHandDeletedBaseUrlOfAnUnknownProviderStaysEmptyInTheFile
TestHandDeletedScalarLineComesBackAtSchemaDefault
TestHandDeletedMapEntryIsNotRestored
TestNothingWrittenPromiseOnlyHoldsWhenSchemaIsCurrent
TestNothingWrittenLeavesFileUntouchedWhenSchemaIsCurrent
TestDeletedLineRestorationIsNotClaimedAsAWrittenKey
TestLoadFilePresetInheritance
TestTicket257R1AC1CleanMachineRefusesAllSevenFields"

echo "######## AC#1 controls: per-name verdicts in each shot (exact grep ruler) ########"
echo "Ruler: grep -E '^ *--- (PASS|FAIL): <name>( |\\()' on the -v log; scope = ./internal/config/"
echo "on the export tree, mutated, -count=1 -timeout 30m."
echo
echo "★NAMED RULER DEFECT of my own (registered, not used as a reading): the awk pairing"
echo "'failing case :: first red line' inside 10-mu-*.txt mis-associates names, because"
echo "go test -v prints a case's failure bodies BEFORE its '^--- FAIL:' line. The authoritative"
echo "per-case mapping is the block below plus the named red roster in each 10-mu-*.txt."
echo
for x in p2 p3 p5 p6; do
  out="$W/log-$x/test.log"
  echo "===== SHOT $x ====="
  echo "RUN=$(grep -c '^=== RUN' "$out" || true) PASS_TOP=$(grep -c '^--- PASS' "$out" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$out" || true) SKIP=$(grep -c '^--- SKIP' "$out" || true)"
  while IFS= read -r t; do
    [ -z "$t" ] && continue
    v=$(grep -E "^ *--- (PASS|FAIL): $t( |\()" "$out" | sed 's/^ *//' | awk '{print $2}' | head -1)
    printf '  %-64s %s\n' "$t" "${v:-ABSENT}"
  done < <(printf '%s\n' "$CASES")
  echo "  -- red roster (top-level, named) --"
  sed 's/^/     /' "$W/log-$x/rednames.txt"
  echo
done
echo "-- unmutated delivery tree: PASS lines for the five new cases (positive controls, surface A) --"
grep -E '^--- PASS: TestHandDeleted|^--- PASS: TestNothingWritten|^--- PASS: TestDeletedLine' "$W/run-deliv-A/test.log" | sed 's/^/  /'
echo "-- baseline surface: the five must be ABSENT there (that is the AC#0 reading) --"
for t in TestHandDeletedDefaultValuedLineIsRestoredInvisibly TestHandDeletedWholeSectionComesBackFromDefaultsNotMemory TestHandDeletedSchemaVersionOnTheWriteBranchIsForcedBackAndUnclaimed TestHandDeletedProviderBaseUrlComesBackAsThePresetValue TestHandDeletedBaseUrlOfAnUnknownProviderStaysEmptyInTheFile; do
  printf '  %-64s hits=%s\n' "$t" "$(grep -cE "^=== RUN   $t\$" "$W/run-base-A/test.log" || true)"
done
