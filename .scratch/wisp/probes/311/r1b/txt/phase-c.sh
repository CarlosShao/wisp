#!/usr/bin/env bash
# phase-c.sh - 311-r1b: closes the two pairings phase A/B left thin.
#  (1) shape 2 (same volume, different directory) needs a guard-ABSENT side at the
#      winsec entry too - T_pre still HAS that block (3b857e53), so the absent side
#      is the stripped tree, exactly as in phase B1.
#  (2) shape 2 at the portable entry: the guard-absent side is T_pre (that tree's
#      portable-tests.sh carries no 311 block at all), which is the same world as
#      "block removed" for that entry.
#  (3) one more offender shape at two entries: caller-provided GOTMPDIR only
#      (WISP_TEST_GOTMP_ROOT unset) - the guard must still read the volume and refuse.
set -u
WT="/d/work/workspace/projects plans/Wisp"
TPRE=/d/tmp/311r1b/tpre
TSTRIP=/d/tmp/311r1b/tstrip
BK=/d/tmp/311r1b/backup
L="$WT/.scratch/wisp/probes/311/r1b/logs"
R="$WT/.scratch/wisp/probes/311/r1b/txt/run-shot.sh"
M="$WT/.scratch/wisp/probes/311/r1b/txt/mk-block.sh"
XVOL='D:/tmp/wisp-test-gotmp/311r1b-xvol'
SAMEVOL='C:/Users/swq/AppData/Local/311r1b-samevol'
FILES="scripts/winsec-tests.sh scripts/portable-tests.sh scripts/wisp-cli-tests.sh tools/d22scan/runtests.sh"
land() { echo "$@" >> "$L/21-strip-proof.txt"; }

echo "== C1: winsec entry, shape 2 (same volume), ALL FOUR blocks removed =="
land "=== C1: strip ALL FOUR again (for the shape-2 guard-absent side) ==="
for f in $FILES; do (cd "$TSTRIP" && sh "$M" strip "$f"); done >> "$L/21-strip-proof.txt" 2>&1
land "c1_block_present_after_strip=$(for f in $FILES; do printf '%s:%s ' "$f" "$(sh "$M" has "$TSTRIP/$f")"; done)"
sh "$R" "$L/26-C1-winsec-samevol-guardremoved.txt" "$TSTRIP" "WISP_TEST_GOTMP_ROOT=$SAMEVOL" -- bash scripts/winsec-tests.sh ./internal/winsec/
for f in $FILES; do cp "$BK/$f" "$TSTRIP/$f"; done
land "c1_restored_sha_match=$( [ "$(cd "$TSTRIP" && sha256sum $FILES)" = "$(cd "$BK" && sha256sum $FILES)" ] && echo yes || echo NO)"

echo "== C2: portable entry, shape 2, guard-absent side (T_pre has no block there) =="
sh "$R" "$L/29-P7-pre-portable-samevol.txt" "$TPRE" "WISP_TEST_GOTMP_ROOT=$SAMEVOL" -- bash scripts/portable-tests.sh ./internal/winsec/

echo "== C3: caller-provided GOTMPDIR only, cross-volume (two entries) =="
sh "$R" "$L/20-C3-runtests-gotmpdir-only-xvol.txt" "/d/tmp/311r1b/tpost" "GOTMPDIR=$XVOL/gotmpdir-only" -- sh tools/d22scan/runtests.sh -C tools/d22scan ./...
sh "$R" "$L/20-C4-cli-gotmpdir-only-xvol-worktree.txt" "$WT" "GOTMPDIR=$XVOL/gotmpdir-only-cli" -- bash scripts/wisp-cli-tests.sh
echo "phase-c-finished"
