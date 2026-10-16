#!/bin/sh
# 166-r2 · PHYSICAL mutation of the exported tree's stepping clock (teeth / zheng kong).
# Not -overlay, not a build tag: the .go file on disk really changes, so the
# compiler and the test binary both see the mutated bytes.
#
# usage: sh 06-mutate-clock.sh <tree-dir> mut|restore
#   mut     -> deletes 'tick = tick.Add(time.Minute); ' from newTestManager, so
#              every call returns the SAME instant = a deliberately flat updated_at
#   restore -> copies the pristine .orig back and proves byte identity (cmp rc=0)
tree=$1
mode=$2
f="$tree/internal/session/manage_test.go"

if [ ! -f "$f" ]; then
  echo "rc_target=1 missing=$f"
  exit 1
fi

case "$mode" in
mut)
  if [ ! -f "$f.orig" ]; then
    cp "$f" "$f.orig"
    echo "rc_cp_orig=$? saved=$f.orig"
  else
    echo "rc_cp_orig=skipped orig_already_exists=1"
  fi
  echo "sha_before=$(sha256sum "$f" | cut -c1-16)"
  echo "advance_lines_before=$(grep -c 'tick = tick.Add(time.Minute)' "$f")"
  awk '{ if ($0 ~ /tick = tick\.Add\(time\.Minute\)/) sub(/tick = tick\.Add\(time\.Minute\); /, ""); print }' "$f.orig" > "$f.mut"
  echo "rc_awk=$?"
  echo "advance_lines_in_mut=$(grep -c 'tick = tick.Add(time.Minute)' "$f.mut")"
  echo "return_tick_lines_in_mut=$(grep -c 'func() time.Time { return tick }' "$f.mut")"
  echo "lines_in_mut=$(wc -l < "$f.mut")"
  mv "$f.mut" "$f"
  echo "rc_mv=$?"
  echo "sha_after=$(sha256sum "$f" | cut -c1-16)"
  ;;
restore)
  if [ ! -f "$f.orig" ]; then
    echo "rc_restore=1 no_orig=$f.orig"
    exit 1
  fi
  cp "$f.orig" "$f"
  echo "rc_cp_restore=$?"
  cmp "$f" "$f.orig"
  echo "rc_cmp_vs_orig=$?"
  echo "sha_restored=$(sha256sum "$f" | cut -c1-16)"
  echo "advance_lines_restored=$(grep -c 'tick = tick.Add(time.Minute)' "$f")"
  ;;
*)
  echo "usage: $0 <tree> mut|restore"
  exit 2
  ;;
esac
