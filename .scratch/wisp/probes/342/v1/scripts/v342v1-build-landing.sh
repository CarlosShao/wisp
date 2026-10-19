#!/usr/bin/env bash
# 342-v1 landing builder. ⛔ the shared working tree: only exported trees under
# /d/tmp/wisp/20261018/342v1. Content anchors, verified on BOTH sides of the cut
# (public-block rule: pre-census expect 1 old hit, post-census expect 0 old + 1 new;
# otherwise an empty operation would masquerade as a successful re-run).
#   old anchor (want exactly 1): f, err := os.OpenFile(resolved, flags|oNoFollowPlatform, perm)
#   new anchor (want exactly 1): f, err := openPrivateAt(resolved, flags, perm)
set -u
S='/d/tmp/wisp/20261018/342v1'
FORM="$1"
SRC="$S/treeV"
DST="$S/mut-$FORM"
[ -d "$DST" ] || cp -r "$SRC" "$DST"
WG="$DST/internal/winsec/winsec.go"

echo "== pre-cut =="
OLD_N=$(grep -c 'f, err := os\.OpenFile(resolved, flags|oNoFollowPlatform, perm)' "$WG" || true)
echo "old-anchor hits = $OLD_N (want 1)"
[ "$OLD_N" = "1" ] || { echo "ABORT: anchor census mismatch"; exit 2; }
PRISTINE=$(sha256sum "$WG" | cut -d' ' -f1)
echo "pristine-sha winsec.go = $PRISTINE"

case "$FORM" in
  jia)     cp "$S/probe/v342v1-landing-jia-linux.go.txt"  "$DST/internal/winsec/v342v1landing_linux.go" ;;
  yi)      cp "$S/probe/v342v1-landing-yi-linux.go.txt"   "$DST/internal/winsec/v342v1landing_linux.go" ;;
  yinofb)  cp "$S/probe/v342v1-landing-yinofb-linux.go.txt" "$DST/internal/winsec/v342v1landing_linux.go" ;;
  *) echo "unknown form $FORM"; exit 2 ;;
esac
cp "$S/probe/v342v1-posix-race.go.txt" "$DST/internal/winsec/v342v1race_test.go"
sed -i 's@f, err := os\.OpenFile(resolved, flags|oNoFollowPlatform, perm)@f, err := openPrivateAt(resolved, flags, perm)@' "$WG"

echo "== post-cut =="
OLD2=$(grep -c 'f, err := os\.OpenFile(resolved, flags|oNoFollowPlatform, perm)' "$WG" || true)
NEW2=$(grep -c 'f, err := openPrivateAt(resolved, flags, perm)' "$WG" || true)
echo "old-anchor after cut = $OLD2 (want 0) ; new-anchor = $NEW2 (want 1)"
[ "$OLD2" = "0" ] && [ "$NEW2" = "1" ] || { echo "ABORT: cut did not land"; exit 2; }
echo "mutated-sha winsec.go = $(sha256sum "$WG" | cut -d' ' -f1)"
sha256sum "$DST/internal/winsec/v342v1landing_linux.go" "$DST/internal/winsec/v342v1race_test.go"
echo "BUILD-OK $FORM"
