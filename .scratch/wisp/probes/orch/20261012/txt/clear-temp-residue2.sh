#!/usr/bin/env bash
# Owner's word "清" (2026-10-12 16:3x) = authorization to clear the legacy residue
# directories at %TEMP% top level. This is the SECOND pass: the first one
# (orch-clear-temp-residue-20261012-164248.txt) aborted on purpose because one
# directory could not be attributed to this repo. Ticket 310 already ruled that
# directory "归属不明、永久不许动", so this pass keeps it by name and deletes the
# rest of the attributable, stale, in-scope set.
#
# Scope rules, in order:
#   1) ONLY paths matching %TEMP%\Test*  (never the wisp-test-gotmp roots, never D:/tmp)
#   2) each directory's PARENT test name must exist in this repo at HEAD
#   3) named KEEP list = anything not attributable -> stays, logged (never deleted)
#   4) skip anything modified in the last 6 hours (an in-flight leg could own it)
#   5) never takeown / icacls / chmod -- a directory that refuses just stays, logged
set -u
REPO="D:/work/workspace/projects plans/Wisp"
TEMP="/c/Users/swq/AppData/Local/Temp"
LOG="$REPO/.scratch/wisp/probes/orch/20261012/logs/orch-clear-temp2-$(date +%Y%m%d-%H%M%S).txt"
cd "$REPO" || exit 9
exec 3>"$LOG"
say() { printf '%s\n' "$*" >&3; }

say "clear-temp-residue PASS2 TS=$(date '+%Y-%m-%d %H:%M:%S %z') scope=%TEMP%/Test* only"
say "head=$(git rev-parse HEAD)"
say "count_before=$(ls -d "$TEMP"/Test* 2>/dev/null | wc -l)"
say "positive_control_all_top_entries=$(ls "$TEMP" | wc -l)"

say ""
say "=== STEP 0 pre-state roster (name + mtime, every Test* dir) ==="
for d in "$TEMP"/Test*; do
  [ -d "$d" ] || continue
  m=$(stat -c %Y "$d" 2>/dev/null || echo 0)
  say "pre dir=$(basename "$d") mtime=$(date -d "@$m" '+%m-%d_%H:%M') age_h=$(( ( $(date +%s) - m ) / 3600 ))"
done

say ""
say "=== STEP 1 live-binary check (a running test could still own a fresh dir) ==="
say "wisp_test_exe_now=$(tasklist //FI 'IMAGENAME eq wisp.test.exe' 2>/dev/null | grep -c 'wisp.test.exe')"
say "audio_test_exe_now=$(tasklist //FI 'IMAGENAME eq audio.test.exe' 2>/dev/null | grep -c 'audio.test.exe')"

say ""
say "=== STEP 2 attribution + KEEP-by-name ==="
now=$(date +%s)
cutoff=$((now-21600))
declare -a LIST=()
kept_unresolved=0; skipped_recent=0; attributed=0
for d in "$TEMP"/Test*; do
  [ -d "$d" ] || continue
  base=$(basename "$d")
  fam1=$(printf '%s' "$base" | sed 's/[0-9]*$//')
  fam2=$(printf '%s' "$fam1" | sed -E 's/[^ -~].*$//' | sed 's/[0-9_]*$//')
  hits=0; used=""
  for fam in "$fam1" "$fam2"; do
    [ -z "$fam" ] && continue
    h=$(git grep -l -- "$fam" HEAD -- cmd internal tools | wc -l)
    if [ "$h" -gt 0 ]; then hits=$h; used=$fam; break; fi
  done
  if [ "$hits" -eq 0 ]; then
    say "KEEP_UNRESOLVED=$base tried=[$fam1|${fam2}] files_naming_it=0"
    kept_unresolved=$((kept_unresolved+1))
    continue
  fi
  attributed=$((attributed+1))
  m=$(stat -c %Y "$d" 2>/dev/null || echo 0)
  if [ "$m" -gt "$cutoff" ]; then
    say "SKIP_RECENT=$base mtime=$(date -d "@$m" '+%m-%d_%H:%M')"
    skipped_recent=$((skipped_recent+1))
    continue
  fi
  say "ATTRIB=$base <- '$used' in $hits file(s) at HEAD"
  LIST+=("$d")
done
say "attributed=$attributed kept_unresolved=$kept_unresolved skipped_recent=$skipped_recent candidates=${#LIST[@]}"

say ""
say "=== STEP 3 delete, one by one, logging each ==="
ok=0; denied=0; other=0
for d in "${LIST[@]}"; do
  out=$(rm -rf "$d" 2>&1); rc=$?
  if [ $rc -eq 0 ]; then
    ok=$((ok+1)); say "OK=$(basename "$d")"
  elif printf '%s' "$out" | grep -qi 'access is denied'; then
    denied=$((denied+1)); say "DENIED=$(basename "$d") :: $(printf '%s' "$out" | head -1 | cut -c1-140)"
  else
    other=$((other+1)); say "OTHER_RC=$rc :: $(basename "$d") :: $(printf '%s' "$out" | head -1 | cut -c1-140)"
  fi
done
say "deleted_ok=$ok denied_left_in_place=$denied other_failures=$other"

say ""
say "=== STEP 4 post-state ==="
say "count_after=$(ls -d "$TEMP"/Test* 2>/dev/null | wc -l)"
say "post_roster:"
for d in "$TEMP"/Test*; do
  [ -d "$d" ] || continue
  m=$(stat -c %Y "$d" 2>/dev/null || echo 0)
  say "post dir=$(basename "$d") mtime=$(date -d "@$m" '+%m-%d_%H:%M')"
done
say "positive_control_all_top_entries_after=$(ls "$TEMP" | wc -l)"
say "untouched_default_gotmp_dirs=$(ls -d C:/Users/swq/AppData/Local/wisp-test-gotmp/*/ 2>/dev/null | wc -l)"
say "untouched_d_tmp_gotmp_dirs=$(ls -d D:/tmp/wisp-test-gotmp/*/ 2>/dev/null | wc -l)"
say "wb_exe_and_balldebug_present=$(ls build/wisp.exe 2>/dev/null | wc -l)"
say "rc=0"
