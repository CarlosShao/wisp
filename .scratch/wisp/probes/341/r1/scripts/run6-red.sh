#!/bin/bash
# 341-r1: pull the verbatim red sentences out of the delivery-tree run (the two FAIL_TOP the color gauge reported).
set -u
out='/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/341/r1/logs'
log="$out/05-color-Adflt-deliv.txt"
{
  echo "[which cases are red]"
  grep '^--- FAIL' "$log"
  echo
  echo "[my five case names, with their verdict]"
  for f in TestHandDeletedDefaultValuedLineIsRestoredInvisibly \
           TestHandDeletedWholeSectionComesBackFromDefaultsNotMemory \
           TestHandDeletedSchemaVersionOnTheWriteBranchIsForcedBackAndUnclaimed \
           TestHandDeletedProviderBaseUrlComesBackAsThePresetValue \
           TestHandDeletedBaseUrlOfAnUnknownProviderStaysEmptyInTheFile; do
    v=$(grep -E "^--- (PASS|FAIL): $f" "$log" || echo "NOT-RUN")
    echo "    $v"
  done
  echo
  echo "[verbatim failure bodies for the red ones]"
  awk '/^=== RUN/{name=$0} /^--- FAIL:/{show=1} show && /^ *writeguard_unnailed_341_test.go:/{print FILENAME":"NR": "$0} /^--- FAIL:/{print "HEADER: "$0} show && /^ *--- /&&!/^--- FAIL/{show=0}' "$log" | head -40
  echo
  echo "[context: the 20 lines after each FAIL header]"
  grep -A 22 '^--- FAIL' "$log" | head -60
} > "$out/06-red-analysis.txt" 2>&1
echo "written $(wc -c < "$out/06-red-analysis.txt") bytes"
