#!/usr/bin/env bash
# phase-a.sh - 311-r1b AC#2 teeth, the "guard PRESENT" side.
# Every shot goes through txt/run-shot.sh, so each capture carries its own command
# line, its env assignments, its台面 (workdir) and a trailing rc= line.
#   TPOST  = out-of-repo export tree of this leg's committed state (git archive HEAD)
#   WT     = the shared working tree (used ONLY where a shot cannot pass without the
#            untracked third_party DLLs, and only for shots that refuse BEFORE any
#            test starts - cmd/wisp belongs to leg 295-r1 this wave)
set -u
WT="/d/work/workspace/projects plans/Wisp"
TPOST=/d/tmp/311r1b/tpost
L="$WT/.scratch/wisp/probes/311/r1b/logs"
R="$WT/.scratch/wisp/probes/311/r1b/txt/run-shot.sh"
XVOL='D:/tmp/wisp-test-gotmp/311r1b-xvol'
BSLASH='D:\tmp\wisp-test-gotmp\311r1b-bslash'
SAMEVOL='C:/Users/swq/AppData/Local/311r1b-samevol'
NOTADIR='/d/tmp/311r1b-notadir'
RELBARE='311r1b-relbare'

# the DEGRADED-then-GUARD trip needs a value that cannot be mkdir'd but whose volume
# is still readable: a regular FILE standing where the root base should be.
[ -e "$NOTADIR" ] || printf '311-r1b: not a directory, planted on purpose for the DEGRADED branch\n' > "$NOTADIR"

echo "== A: cross-volume refusal, one per entry (AC#2 (a)) =="
sh "$R" "$L/20-A1-winsec-xvol.txt"        "$TPOST" "WISP_TEST_GOTMP_ROOT=$XVOL"  -- bash scripts/winsec-tests.sh ./internal/winsec/
sh "$R" "$L/20-A2-portable-xvol.txt"      "$TPOST" "WISP_TEST_GOTMP_ROOT=$XVOL"  -- bash scripts/portable-tests.sh ./internal/winsec/
sh "$R" "$L/20-A3-runtests-xvol.txt"      "$TPOST" "WISP_TEST_GOTMP_ROOT=$XVOL"  -- sh tools/d22scan/runtests.sh -C tools/d22scan ./...
sh "$R" "$L/20-A4-cli-xvol-worktree.txt"  "$WT"    "WISP_TEST_GOTMP_ROOT=$XVOL"  -- bash scripts/wisp-cli-tests.sh

echo "== B: line order inside one trip (AC#2 (d), the two readings 3b857e53 named but never ran) =="
sh "$R" "$L/23-order-degraded-then-guard.txt" "$TPOST" "WISP_TEST_GOTMP_ROOT=$NOTADIR" -- bash scripts/winsec-tests.sh ./internal/winsec/
sh "$R" "$L/24-order-guard-after-guard1.txt"  "$TPOST" "WISP_TEST_GOTMP_ROOT=$XVOL"     -- bash scripts/winsec-tests.sh --scope=winsec

echo "== C: spelling coverage of the volume reader (no new judgement, no exit-code change) =="
sh "$R" "$L/25-C1-runtests-xvol-backslash.txt"    "$TPOST" "WISP_TEST_GOTMP_ROOT=$BSLASH" -- sh tools/d22scan/runtests.sh -C tools/d22scan ./...
sh "$R" "$L/25-C2-runtests-inconclusive-rel.txt"  "$TPOST" "WISP_TEST_GOTMP_ROOT=$RELBARE" -- sh tools/d22scan/runtests.sh -C tools/d22scan ./...

echo "== D: the three shapes the guard must NOT break, guard-present side =="
sh "$R" "$L/26-D1-winsec-samevol.txt"   "$TPOST" "WISP_TEST_GOTMP_ROOT=$SAMEVOL" -- bash scripts/winsec-tests.sh ./internal/winsec/
sh "$R" "$L/26-D2-portable-samevol.txt" "$TPOST" "WISP_TEST_GOTMP_ROOT=$SAMEVOL" -- bash scripts/portable-tests.sh ./internal/winsec/
sh "$R" "$L/26-D3-runtests-samevol.txt" "$TPOST" "WISP_TEST_GOTMP_ROOT=$SAMEVOL" -- sh tools/d22scan/runtests.sh -C tools/d22scan ./...
sh "$R" "$L/27-E1-winsec-nolappdata.txt"  "$TPOST" "LOCALAPPDATA=" -- bash scripts/winsec-tests.sh ./internal/winsec/
sh "$R" "$L/27-E2-runtests-nolappdata.txt" "$TPOST" "LOCALAPPDATA=" -- sh tools/d22scan/runtests.sh -C tools/d22scan ./...
sh "$R" "$L/27-E3-portable-nolappdata.txt" "$TPOST" "LOCALAPPDATA=" -- bash scripts/portable-tests.sh ./internal/winsec/
sh "$R" "$L/28-F1-winsec-gooslinux.txt" "$TPOST" "GOOS=linux" -- bash scripts/winsec-tests.sh ./internal/winsec/
sh "$R" "$L/28-F2-cli-gooslinux.txt"    "$TPOST" "GOOS=linux" -- bash scripts/wisp-cli-tests.sh
sh "$R" "$L/28-F3-runtests-gooslinux.txt" "$TPOST" "GOOS=linux" -- sh tools/d22scan/runtests.sh -C tools/d22scan ./...

echo "== E: AC#4 gate readings on this leg's committed state =="
sh "$R" "$L/30-G1-post-build.txt"   "$TPOST" -- env GOFLAGS= go build ./...
sh "$R" "$L/30-G2-post-census.txt"  "$TPOST" -- bash scripts/portable-tests.sh --scope=census
sh "$R" "$L/30-G3-post-d22scan.txt" "$TPOST" -- sh scripts/d22scan.sh
sh "$R" "$L/30-G4-post-winsec-run1.txt" "$TPOST" -- bash scripts/winsec-tests.sh ./internal/winsec/
sh "$R" "$L/30-G5-post-winsec-run2.txt" "$TPOST" -- bash scripts/winsec-tests.sh ./internal/winsec/
sh "$R" "$L/30-G6-post-census-worktree.txt" "$WT" -- bash scripts/portable-tests.sh --scope=census
echo "phase-a-sequence-finished"
