#!/usr/bin/env bash
# phase-b.sh - 311-r1b AC#2 (b): the SENSITIVITY positive control.
# Removing ONLY the guard block(s) must bring the colour back. Done in an out-of-repo
# export tree (tstrip) so the shared working tree and the pristine tpost tree are never
# touched; every strip is bracketed by sha256 evidence (pristine -> stripped -> restored).
set -u
WT="/d/work/workspace/projects plans/Wisp"
TPOST=/d/tmp/311r1b/tpost
TSTRIP=/d/tmp/311r1b/tstrip
BK=/d/tmp/311r1b/backup
L="$WT/.scratch/wisp/probes/311/r1b/logs"
R="$WT/.scratch/wisp/probes/311/r1b/txt/run-shot.sh"
M="$WT/.scratch/wisp/probes/311/r1b/txt/mk-block.sh"
XVOL='D:/tmp/wisp-test-gotmp/311r1b-xvol'
FILES="scripts/winsec-tests.sh scripts/portable-tests.sh scripts/wisp-cli-tests.sh tools/d22scan/runtests.sh"

sha() { (cd "$1" && sha256sum $FILES); }
land() { echo "$@" >> "$L/21-strip-proof.txt"; }

if [ ! -d "$TSTRIP" ]; then cp -r "$TPOST" "$TSTRIP"; fi
mkdir -p "$BK"
for f in $FILES; do mkdir -p "$BK/$(dirname "$f")"; cp "$TPOST/$f" "$BK/$f"; done

land "=== pristine baseline (tstrp copied from tpost; backup made) ==="
land "sha_pristine_tpost:"; sha "$TPOST" | sed 's/^/  /' >> "$L/21-strip-proof.txt"
land "sha_pristine_tstrip_before_any_strip:"; sha "$TSTRIP" | sed 's/^/  /' >> "$L/21-strip-proof.txt"
land "sha_backup:"; sha "$BK" | sed 's/^/  /' >> "$L/21-strip-proof.txt"
land "block_present_tstrip_before=$(for f in $FILES; do printf '%s:%s ' "$f" "$(sh "$M" has "$TSTRIP/$f")"; done)"

restore_all() {
    for f in $FILES; do cp "$BK/$f" "$TSTRIP/$f"; done
    land "restored_from_backup; sha_tstrip_now:"
    (cd "$TSTRIP" && sha256sum $FILES) | sed 's/^/  /' >> "$L/21-strip-proof.txt"
    land "restore_equals_pristine=$( [ "$(cd "$TSTRIP" && sha256sum $FILES | cut -d' ' -f1 | tr '\n' ' ')" = "$(cd "$BK" && sha256sum $FILES | cut -d' ' -f1 | tr '\n' ' ')" ] && echo yes || echo NO)"
}

# ---- B1: all four guard blocks removed, nothing else ----
land "=== B1: strip ALL FOUR guard blocks (only the blocks) ==="
for f in $FILES; do (cd "$TSTRIP" && sh "$M" strip "$f"); done >> "$L/21-strip-proof.txt" 2>&1
land "block_present_tstrip_after_strip4=$(for f in $FILES; do printf '%s:%s ' "$f" "$(sh "$M" has "$TSTRIP/$f")"; done)"
sh "$R" "$L/22-B1-winsec-xvol-guard-removed.txt" "$TSTRIP" "WISP_TEST_GOTMP_ROOT=$XVOL" -- bash scripts/winsec-tests.sh ./internal/winsec/
sh "$R" "$L/26-B1x-winsec-default-guardremoved.txt" "$TSTRIP" -- bash scripts/winsec-tests.sh ./internal/winsec/
sh "$R" "$L/26-B1y-portable-default-guardremoved.txt" "$TSTRIP" -- bash scripts/portable-tests.sh ./internal/winsec/
restore_all

# ---- B2: only portable-tests.sh stripped (layered defence, per-entry sensitivity) ----
land "=== B2: strip ONLY scripts/portable-tests.sh ==="
(cd "$TSTRIP" && sh "$M" strip scripts/portable-tests.sh) >> "$L/21-strip-proof.txt" 2>&1
land "block_present_after_strip_portable=$(for f in $FILES; do printf '%s:%s ' "$f" "$(sh "$M" has "$TSTRIP/$f")"; done)"
sh "$R" "$L/22-B2-portable-xvol-own-guard-removed.txt" "$TSTRIP" "WISP_TEST_GOTMP_ROOT=$XVOL" -- bash scripts/portable-tests.sh ./internal/winsec/
restore_all

# ---- B3: only runtests.sh stripped (colour back at the deepest entry) ----
land "=== B3: strip ONLY tools/d22scan/runtests.sh ==="
(cd "$TSTRIP" && sh "$M" strip tools/d22scan/runtests.sh) >> "$L/21-strip-proof.txt" 2>&1
land "block_present_after_strip_runtests=$(for f in $FILES; do printf '%s:%s ' "$f" "$(sh "$M" has "$TSTRIP/$f")"; done)"
sh "$R" "$L/22-B3-runtests-xvol-guard-removed.txt" "$TSTRIP" "WISP_TEST_GOTMP_ROOT=$XVOL" -- sh tools/d22scan/runtests.sh ./internal/winsec/
sh "$R" "$L/26-B3x-runtests-default-guardremoved.txt" "$TSTRIP" -- sh tools/d22scan/runtests.sh -C tools/d22scan ./...
restore_all

land "=== final: tstrip back to pristine, tpost untouched throughout ==="
land "sha_tpost_final:"; sha "$TPOST" | sed 's/^/  /' >> "$L/21-strip-proof.txt"
land "sha_tstrip_final:"; sha "$TSTRIP" | sed 's/^/  /' >> "$L/21-strip-proof.txt"
echo "phase-b-sequence-finished"
