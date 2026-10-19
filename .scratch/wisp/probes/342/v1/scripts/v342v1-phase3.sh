#!/usr/bin/env bash
# 342-v1 phase 3: resolve the container-config question, then take the colour faces and
# the landed-production rows on the SAME config the census leg used (⛔ GOTMPDIR bind).
# Also: the VPROD magnitude attribution (mandatory (i)) as a racer-pacing sweep, and the
# yinofb mis-shape positive control (mandatory (iii)) with its red sentences kept verbatim.
set -u
S=/d/tmp/wisp/20261018/342v1
P=$S/probe
L=$S/logs
D() { bash "$S/scripts/v342v1-dockerB.sh" "$@"; }

echo "### (0) probe bytes staged (hash after the pacing edit):"
cp "$P/v342v1-posix-race.go.txt" "$S/treeV/internal/winsec/v342v1race_test.go"
for form in jia yi; do
  cp "$P/v342v1-posix-race.go.txt" "$S/mut-$form/internal/winsec/v342v1race_test.go"
  cp "$P/v342v1-landed-prod.go.txt" "$S/mut-$form/internal/winsec/v342v1landed_test.go"
done
sha256sum "$S/treeV/internal/winsec/v342v1race_test.go" "$S/mut-jia/internal/winsec/v342v1landed_test.go"
for form in jia yi; do echo "winsec.go sha mut-$form = $(sha256sum "$S/mut-$form/internal/winsec/winsec.go" | cut -d' ' -f1)"; done

echo "### (1) pristine baseline colour on B config, quiet x3 (ticket 346 flakiness rule)"
for r in 1 2 3; do
  D treeC "go test -count=1 -timeout 40m -v ./internal/winsec/" "14-colorB-treeC-r$r.log"
  grep -E '^--- FAIL' "$L/14-colorB-treeC-r$r.log" | sed 's/ ([0-9.]*s)//' > "$L/14-colorB-treeC-r$r-failnames.txt" || true
  echo "r$r failnames: $(wc -l < "$L/14-colorB-treeC-r$r-failnames.txt")"; cat "$L/14-colorB-treeC-r$r-failnames.txt" || true
done

echo "### (2) VPROD attribution sweep: racer pacing, 1500 attempts, treeV"
for us in 0 100 300 1000; do
  D treeV "V342V1_RACER_US=$us V342V1_ATTEMPTS=1500 go test -count=1 -timeout 40m -v -run 'TestV342v1PosixAncestorRaceRows/VPROD' ./internal/winsec/" "15-vprod-pace-$us.log"
  grep -E 'V342V1 VPROD: attempts' "$L/15-vprod-pace-$us.log" || true
  grep -E 'V342V1 VPROD leg=' "$L/15-vprod-pace-$us.log" || true
done

echo "### (3) landed production rows on the two POSIX shapes (B config), 1500 x2"
for form in jia yi; do
  for r in 1 2; do
    D "mut-$form" "V342V1_ATTEMPTS=1500 go test -count=1 -timeout 40m -v -run TestV342v1PosixLandedProductionRow ./internal/winsec/" "16-landed-$form-r$r.log"
    grep -E 'V342V1 VPROD-LANDED: attempts' "$L/16-landed-$form-r$r.log" || true
    grep -E 'V342V1 VPROD-LANDED leg=' "$L/16-landed-$form-r$r.log" | head -8 || true
  done
done

echo "### (4) colour faces on B config: baseline vs the three landing faces"
for tree in treeC col-jia col-yi col-yinofb; do
  D "$tree" "go test -count=1 -timeout 40m -v ./internal/winsec/" "17-colorB-$tree.log"
  grep -E '^--- FAIL' "$L/17-colorB-$tree.log" | sed 's/ ([0-9.]*s)//' > "$L/17-colorB-$tree-failnames.txt" || true
  echo "failnames $tree ($(wc -l < "$L/17-colorB-$tree-failnames.txt")):"; cat "$L/17-colorB-$tree-failnames.txt" || true
done
echo "### (5) yinofb red sentences verbatim:"
for n in openleaf_322 openleaf_overwrite sealhandle_privateroute; do
  echo "-- $n"; grep -A 6 "^--- FAIL: .*${n}" "$L/17-colorB-col-yinofb.log" | head -28 || true
done
echo "PHASE3-DONE"
