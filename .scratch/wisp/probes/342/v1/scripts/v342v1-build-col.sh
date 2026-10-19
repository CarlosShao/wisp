#!/usr/bin/env bash
# 342-v1 colour-face builder: the four-number comparisons for AC#0's "which existing pins
# go red" must NOT contain this leg's own probe files (they add their own RUN rows).
# So the colour faces are built from a CLEAN export tree: treeC (git archive HEAD, no
# probes) + landing + the same content-anchored cut, old/new anchor census on both sides.
set -u
S=/d/tmp/wisp/20261018/342v1
SRC="$S/treeC"
FORM="$1"
DST="$S/col-$FORM"
[ -d "$DST" ] || cp -r "$SRC" "$DST"
# treeC is a fresh git archive: it never received probe files, and this script does NOT
# delete anything (dispatch rule: temp artefacts are created, never removed). The stray
# census below is the check that the face really is repo bytes + landing only.
WG="$DST/internal/winsec/winsec.go"
echo "== pre-cut =="
OLD_N=$(grep -c 'f, err := os\.OpenFile(resolved, flags|oNoFollowPlatform, perm)' "$WG" || true)
STRAY=$(ls "$DST"/internal/winsec/v342v1*_test.go 2>/dev/null | wc -l)
echo "old-anchor hits = $OLD_N (want 1) ; stray probe files = $STRAY (want 0)"
[ "$OLD_N" = "1" ] && [ "$STRAY" = "0" ] || { echo "ABORT"; exit 2; }
case "$FORM" in
  jia)    cp "$S/probe/v342v1-landing-jia-linux.go.txt"     "$DST/internal/winsec/v342v1landing_linux.go" ;;
  yi)     cp "$S/probe/v342v1-landing-yi-linux.go.txt"      "$DST/internal/winsec/v342v1landing_linux.go" ;;
  yinofb) cp "$S/probe/v342v1-landing-yinofb-linux.go.txt"  "$DST/internal/winsec/v342v1landing_linux.go" ;;
  *) echo "unknown form"; exit 2 ;;
esac
sed -i 's@f, err := os\.OpenFile(resolved, flags|oNoFollowPlatform, perm)@f, err := openPrivateAt(resolved, flags, perm)@' "$WG"
echo "== post-cut =="
OLD2=$(grep -c 'f, err := os\.OpenFile(resolved, flags|oNoFollowPlatform, perm)' "$WG" || true)
NEW2=$(grep -c 'f, err := openPrivateAt(resolved, flags, perm)' "$WG" || true)
echo "old = $OLD2 (want 0) ; new = $NEW2 (want 1)"
[ "$OLD2" = "0" ] && [ "$NEW2" = "1" ] || { echo "ABORT: cut did not land"; exit 2; }
echo "landing lines added = $(wc -l < "$DST/internal/winsec/v342v1landing_linux.go")"
echo "winsec.go sha = $(sha256sum "$WG" | cut -d' ' -f1)"
echo "COL-BUILD-OK $FORM"
