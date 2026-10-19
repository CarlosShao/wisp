#!/usr/bin/env bash
# 342-v1 phase 4: the 丙 colour face (report-after-the-act) + the A/B config attribution
# file, all on the B container config (⛔ GOTMPDIR bind) so it is comparable with treeC.
set -u
S=/d/tmp/wisp/20261018/342v1
P=$S/probe
L=$S/logs
R="/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/342/v1/logs"
mkdir -p "$R"

# build col-bing from the clean export tree, same anchored cut + both-side censuses
DST="$S/col-bing"
[ -d "$DST" ] || cp -r "$S/treeC" "$DST"
WG="$DST/internal/winsec/winsec.go"
OLD_N=$(grep -c 'f, err := os\.OpenFile(resolved, flags|oNoFollowPlatform, perm)' "$WG" || true)
echo "pre-cut old-anchor = $OLD_N (want 1)"
[ "$OLD_N" = "1" ] || { echo "ABORT"; exit 2; }
cp "$P/v342v1-landing-bing-linux.go.txt" "$DST/internal/winsec/v342v1landing_linux.go"
sed -i 's@f, err := os\.OpenFile(resolved, flags|oNoFollowPlatform, perm)@f, err := openPrivateAt(resolved, flags, perm)@' "$WG"
OLD2=$(grep -c 'f, err := os\.OpenFile(resolved, flags|oNoFollowPlatform, perm)' "$WG" || true)
NEW2=$(grep -c 'f, err := openPrivateAt(resolved, flags, perm)' "$WG" || true)
echo "post-cut old = $OLD2 (want 0) ; new = $NEW2 (want 1)"
[ "$OLD2" = "0" ] && [ "$NEW2" = "1" ] || { echo "ABORT: cut did not land"; exit 2; }
echo "col-bing landing lines = $(wc -l < "$DST/internal/winsec/v342v1landing_linux.go")"
echo "col-bing winsec.go sha = $(sha256sum "$WG" | cut -d' ' -f1)"

bash "$S/scripts/v342v1-dockerB.sh" col-bing "go test -count=1 -timeout 40m -v ./internal/winsec/" "19-colorB-col-bing.log"
grep -E '^--- FAIL' "$L/19-colorB-col-bing.log" | sed 's/ ([0-9.]*s)//' > "$L/19-colorB-col-bing-failnames.txt" || true
echo "col-bing failnames: $(wc -l < "$L/19-colorB-col-bing-failnames.txt")"; cat "$L/19-colorB-col-bing-failnames.txt" || true

# 丙 behavioural magnitude again on the B config (A config rows for VBING already
# recorded; the shape claim needs at least one matched-config repeat x2)
for r in 1 2; do
  bash "$S/scripts/v342v1-dockerB.sh" treeV "V342V1_ATTEMPTS=1500 go test -count=1 -timeout 40m -v -run 'TestV342v1PosixAncestorRaceRows/VBING' ./internal/winsec/" "20-vbingB-r$r.log"
  grep -E 'V342V1 VBING: attempts' "$L/20-vbingB-r$r.log" || true
done

bash "$S/scripts/v342v1-extract-iii.sh"
echo "PHASE4-DONE"
