#!/usr/bin/env bash
# 342-a1 sim-landing builder (⛔ runs on exported trees ONLY; the shared working tree and
# the repository's internal/** bytes are never touched by this leg).
# Usage: bash build-mut.sh <tree-name: jia|yi> <platform: linux|win>
# Content anchors (落刀前验命中＝期望值，落刀后再验一次＝新锚命中+旧锚归零；公共块定式):
#   old anchor (expect exactly 1): f, err := os.OpenFile(resolved, flags|oNoFollowPlatform, perm)
#   new anchor (expect exactly 1): f, err := openPrivateAt(resolved, flags, perm)
set -u
S='/d/tmp/wisp/20261018/342a1'
FORM="$1"; PLAT="$2"
SRC="$S/treeB"
DST="$S/mut-$FORM-$PLAT"
[ -d "$DST" ] || cp -r "$SRC" "$DST"
WG="$DST/internal/winsec/winsec.go"

echo "== pre-cut anchor census =="
OLD_N=$(grep -c 'f, err := os\.OpenFile(resolved, flags|oNoFollowPlatform, perm)' "$WG" || true)
echo "old-anchor hits = $OLD_N (want 1)"
[ "$OLD_N" = "1" ] || { echo "ABORT: anchor census mismatch, ⛔ cutting"; exit 2; }

case "$FORM-$PLAT" in
  jia-linux)
    cp "$S/probe/sim_jia_posix_linux.go.txt" "$DST/internal/winsec/simlanding_openat_linux.go"
    ;;
  yi-linux)
    cp "$S/probe/sim_yi_posix_linux.go.txt" "$DST/internal/winsec/simlanding_pin_linux.go"
    sed -i 's/func openPrivatePinned(/func openPrivateAt(/' "$DST/internal/winsec/simlanding_pin_linux.go"
    NEWFN=$(grep -c 'func openPrivateAt(' "$DST/internal/winsec/simlanding_pin_linux.go" || true)
    echo "renamed landing fn hits = $NEWFN (want 1)"
    ;;
  jia-win)
    cp "$S/probe/sim_jia_win_windows.go.txt" "$DST/internal/winsec/simlanding_jia_windows.go"
    ;;
  yi-win)
    cp "$S/probe/sim_yi_win_windows.go.txt" "$DST/internal/winsec/simlanding_yi_windows.go"
    ;;
  *) echo "unknown form/platform"; exit 2;;
esac

sed -i 's@f, err := os\.OpenFile(resolved, flags|oNoFollowPlatform, perm)@f, err := openPrivateAt(resolved, flags, perm)@' "$WG"

echo "== post-cut anchor census =="
OLD2=$(grep -c 'f, err := os\.OpenFile(resolved, flags|oNoFollowPlatform, perm)' "$WG" || true)
NEW2=$(grep -c 'f, err := openPrivateAt(resolved, flags, perm)' "$WG" || true)
echo "old-anchor after cut = $OLD2 (want 0) ; new-anchor = $NEW2 (want 1)"
[ "$OLD2" = "0" ] && [ "$NEW2" = "1" ] || { echo "ABORT: cut did not land as expected"; exit 2; }
sha256sum "$WG" "$DST/internal/winsec/simlanding_"*.go 2>/dev/null
echo "BUILD-OK $FORM $PLAT"
