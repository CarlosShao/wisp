#!/usr/bin/env bash
# phase-pre-shapes.sh - 311-r1b: the SAME shapes shot on the PRE side (export tree of
# HEAD 2567e3a9, where only winsec-tests.sh carries the block) so "the guard breaks none
# of the three shapes" is a pre<->post count comparison and not an absence of evidence.
set -u
WT="/d/work/workspace/projects plans/Wisp"
TPRE=/d/tmp/311r1b/tpre
L="$WT/.scratch/wisp/probes/311/r1b/logs"
R="$WT/.scratch/wisp/probes/311/r1b/txt/run-shot.sh"
SAMEVOL='C:/Users/swq/AppData/Local/311r1b-samevol-pre'

sh "$R" "$L/29-P1-pre-runtests-nolappdata.txt"  "$TPRE" "LOCALAPPDATA=" -- sh tools/d22scan/runtests.sh -C tools/d22scan ./...
sh "$R" "$L/29-P2-pre-winsec-nolappdata.txt"    "$TPRE" "LOCALAPPDATA=" -- bash scripts/winsec-tests.sh ./internal/winsec/
sh "$R" "$L/29-P3-pre-runtests-gooslinux.txt"   "$TPRE" "GOOS=linux"    -- sh tools/d22scan/runtests.sh -C tools/d22scan ./...
sh "$R" "$L/29-P4-pre-winsec-gooslinux.txt"     "$TPRE" "GOOS=linux"    -- bash scripts/winsec-tests.sh ./internal/winsec/
sh "$R" "$L/29-P5-pre-runtests-samevol.txt"     "$TPRE" "WISP_TEST_GOTMP_ROOT=$SAMEVOL" -- sh tools/d22scan/runtests.sh -C tools/d22scan ./...
sh "$R" "$L/29-P6-pre-portable-nolappdata.txt"  "$TPRE" "LOCALAPPDATA=" -- bash scripts/portable-tests.sh ./internal/winsec/
echo "phase-pre-shapes-finished"
