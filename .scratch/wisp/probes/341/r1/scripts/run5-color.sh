#!/bin/bash
# 341-r1 bench: internal/config four numbers, baseline tree vs delivery tree, on BOTH temp-root surfaces.
# Surface A "C default": GOTMPDIR unset, TMP/TEMP left as inherited (C: volume).
# Surface B "D pre-created": GOTMPDIR + TMP/TEMP/TEMPDIR point into D:/tmp/wisp/20261018/341r1/tmpB,
#                              which is mkdir -p'ed first (a missing dir shows up as "creating work dir" red).
# Rulers (dispatch 颜色尺): RUN= grep -c '^=== RUN' / PASS_TOP= '^--- PASS' / FAIL_TOP= '^--- FAIL' / SKIP= '^--- SKIP'
# -count=1 -timeout 30m, the four numbers are never mixed and never subtracted.
set -u
repo='/d/work/workspace/projects plans/Wisp'
root='D:/tmp/wisp/20261018/341r1'
out="$repo/.scratch/wisp/probes/341/r1/logs"
cd "$repo" || exit 2
mkdir -p "$root/tmpB" "$root/gotmp" || exit 2
for surf in Adflt Bdtemp; do
  for t in base deliv; do
    tree="$root/$t"
    log="$out/05-color-$surf-$t.txt"
    head="$out/05-color-$surf-$t-head.txt"
    (
      if [ "$surf" = Adflt ]; then
        unset GOTMPDIR
        echo "SURFACE=A default C-volume temp (GOTMPDIR unset; TMP=$TMP TEMP=$TEMP as inherited)"
      else
        export GOTMPDIR="$root/gotmp"
        export TMP="$root/tmpB" TEMP="$root/tmpB" TMPDIR="$root/tmpB"
        echo "SURFACE=B D-precreated (GOTMPDIR=$GOTMPDIR TMP=$TMP)"
      fi
      cd "$tree" || exit 2
      echo "TREE=$t goliles=$(find . -name '*.go' | wc -l)"
      go test ./internal/config/ -count=1 -timeout 30m -v > "$log" 2>&1
      rc=$?
      RUN=$(grep -c '^=== RUN' "$log" || true)
      PASSTOP=$(grep -c '^--- PASS' "$log" || true)
      FAILTOP=$(grep -c '^--- FAIL' "$log" || true)
      SKIP=$(grep -c '^--- SKIP' "$log" || true)
      echo "rc=$rc RUN=$RUN PASS_TOP=$PASSTOP FAIL_TOP=$FAILTOP SKIP=$SKIP"
      echo "---- FAIL roster ----"
      grep '^--- FAIL' "$log" || echo "(none)"
      echo "---- my six cases, named ----"
      grep -E '^(--- |    )?(=== RUN|--- PASS|--- FAIL): Test(HandDeleted|NothingWritten|DeletedLine)' "$log" | grep -E '^(--- (PASS|FAIL)|=== RUN)' | grep -E 'Test(HandDeletedDefaultValuedLineIsRestoredInvisibly|HandDeletedWholeSectionComesBackFromDefaultsNotMemory|HandDeletedSchemaVersionOnTheWriteBranchIsForcedBackAndUnclaimed|HandDeletedProviderBaseUrlComesBackAsThePresetValue|HandDeletedBaseUrlOfAnUnknownProviderStaysEmptyInTheFile)' || echo "(my case names not found in the run output)"
      echo "$log"
    ) > "$head" 2>&1
    echo "== $surf/$t: $(grep '^rc=' "$head")"
  done
done
echo done
