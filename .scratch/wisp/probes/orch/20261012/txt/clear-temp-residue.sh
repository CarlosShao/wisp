#!/usr/bin/env bash
# Owner's single word "清" (16:3x) = authorization to clear the 59 legacy residue
# directories at %TEMP% top level. Scope rules, in order:
#   1) ONLY paths matching %TEMP%\Test*  (never the wisp-test-gotmp roots, never D:/tmp)
#   2) each name family must be a real test name in THIS repo (ruler below)
#   3) skip anything modified in the last 6 hours (an in-flight leg could own it)
#   4) never takeown / icacls / chmod -- a directory that refuses just stays, logged
set -u
REPO="D:/work/workspace/projects plans/Wisp"
TEMP="/c/Users/swq/AppData/Local/Temp"
LOG="$REPO/.scratch/wisp/probes/orch/20261012/logs/orch-clear-temp-residue-$(date +%Y%m%d-%H%M%S).txt"
cd "$REPO" || exit 9
exec 3>"$LOG"
say() { printf '%s\n' "$*" >&3; }
say "clear-temp-residue TS=$(date '+%Y-%m-%d %H:%M:%S %z') scope=%TEMP%/Test* only"
say "count_before=$(ls -d "$TEMP"/Test* 2>/dev/null | wc -l)"
say "positive_control_all_top_entries=$(ls "$TEMP" | wc -l)"

say ""
say "=== STEP 1 attribution: every directory's PARENT test name must exist in this repo ==="
# Go names a temp dir TestX<digits> for a parent test, and TestX/sub becomes
# TestX_sub... -- so the attributable unit is the longest leading identifier that
# appears in the source. Try the whole name with trailing digits stripped, then cut
# at the first non-ASCII byte (subtest suffixes with CJK) and strip again.
unresolved=0
for d in "$TEMP"/Test*; do
  [ -d "$d" ] || continue
  base=$(basename "$d")
  fam1=$(printf '%s' "$base" | sed 's/[0-9]*$//')
  fam2=$(printf '%s' "$fam1" | sed -E 's/[^ -~].*$//' | sed 's/[0-9_]*$//')
  hits=0
  for fam in "$fam1" "$fam2"; do
    [ -z "$fam" ] && continue
    h=$(git grep -l -- "$fam" HEAD -- cmd internal tools | wc -l)
    [ "$h" -gt 0 ] && hits=$h && break
  done
  say "dir=$base tried=[$fam1|${fam2}] files_naming_it=$hits"
  [ "$hits" -eq 0 ] && unresolved=$((unresolved+1))
done
say "dirs_with_zero_hits=$unresolved"

say ""
say "=== STEP 2 freshness guard: skip anything newer than 6 h ==="
now=$(date +%s)
cutoff=$((now-21600))
skipped_recent=0
LIST=()
for d in "$TEMP"/Test*; do
  [ -d "$d" ] || continue
  m=$(stat -c %Y "$d" 2>/dev/null || echo 0)
  if [ "$m" -gt "$cutoff" ]; then
    say "SKIP_RECENT=$d mtime=$(date -d "@$m" '+%m-%d_%H:%M')"
    skipped_recent=$((skipped_recent+1))
  else
    LIST+=("$d")
  fi
done
say "candidates=${#LIST[@]} skipped_recent=$skipped_recent"

if [ "$unresolved" -ne 0 ]; then
  say "VERDICT=ABORT (a directory is not attributable to this repo; nothing deleted)"
  exit 7
fi

say ""
say "=== STEP 3 delete, one by one, logging each ==="
ok=0; denied=0; other=0
for d in "${LIST[@]}"; do
  out=$(rm -rf "$d" 2>&1); rc=$?
  if [ $rc -eq 0 ]; then
    ok=$((ok+1)); say "OK=$d"
  elif printf '%s' "$out" | grep -qi 'access is denied'; then
    denied=$((denied+1)); say "DENIED=$d :: $(printf '%s' "$out" | head -1 | cut -c1-140)"
  else
    other=$((other+1)); say "OTHER_RC=$rc :: $d :: $(printf '%s' "$out" | head -1 | cut -c1-140)"
  fi
done
say "deleted_ok=$ok denied_left_in_place=$denied other_failures=$other"
say "count_after=$(ls -d "$TEMP"/Test* 2>/dev/null | wc -l)"
say "positive_control_all_top_entries_after=$(ls "$TEMP" | wc -l)"
say "untouched_roots_check_default_gotmp_dirs=$(ls -d C:/Users/swq/AppData/Local/wisp-test-gotmp/*/ 2>/dev/null | wc -l)"
say "untouched_roots_check_d_tmp_dirs=$(ls -d D:/tmp/wisp-test-gotmp/*/ 2>/dev/null | wc -l)"
say "rc=0"
