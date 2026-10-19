#!/usr/bin/env bash
# 342-v1 host (Windows) colour faces: pristine baseline vs the two landing shapes, on
# repo bytes ONLY (no probe files in these trees, so the four numbers are the package's
# own roster and nothing else). treeH = git archive HEAD export, untouched.
set -u
S=/d/tmp/wisp/20261018/342v1
L=$S/logs
FORM="$1"
SRC="$S/treeH"
DST="$S/col-win-$FORM"
[ -d "$DST" ] || cp -r "$SRC" "$DST"
WG="$DST/internal/winsec/winsec.go"
echo "== pre-cut =="
OLD_N=$(grep -c 'f, err := os\.OpenFile(resolved, flags|oNoFollowPlatform, perm)' "$WG" || true)
echo "old-anchor hits = $OLD_N (want 1)"
[ "$OLD_N" = "1" ] || { echo "ABORT"; exit 2; }
case "$FORM" in
  jia) cp "$S/probe/v342v1-landing-jia-win.go.txt" "$DST/internal/winsec/v342v1landing_windows.go" ;;
  yi)  cp "$S/probe/v342v1-landing-yi-win.go.txt"  "$DST/internal/winsec/v342v1landing_windows.go" ;;
  *) echo "unknown form"; exit 2 ;;
esac
sed -i 's@f, err := os\.OpenFile(resolved, flags|oNoFollowPlatform, perm)@f, err := openPrivateAt(resolved, flags, perm)@' "$WG"
echo "== post-cut =="
OLD2=$(grep -c 'f, err := os\.OpenFile(resolved, flags|oNoFollowPlatform, perm)' "$WG" || true)
NEW2=$(grep -c 'f, err := openPrivateAt(resolved, flags, perm)' "$WG" || true)
echo "old = $OLD2 (want 0) ; new = $NEW2 (want 1)"
[ "$OLD2" = "0" ] && [ "$NEW2" = "1" ] || { echo "ABORT"; exit 2; }
echo "landing lines = $(wc -l < "$DST/internal/winsec/v342v1landing_windows.go")"
echo "winsec.go sha = $(sha256sum "$WG" | cut -d' ' -f1)"
echo "COLWIN-BUILD-OK $FORM"
