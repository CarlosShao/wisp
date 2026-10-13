#!/usr/bin/env bash
# Pass 2 (2026-10-13 08:1x): the small wisp-named leftovers in the user profile that pass 1
# deliberately left behind. Same rules: explicit allowlist, guard by shape, per-item log,
# nothing outside the list is touched. These are all leg scratch from tickets 84/91/301/303.
set -u
LOG="D:/work/workspace/projects plans/Wisp/.scratch/wisp/probes/orch/20261013/logs/c-clean-small-$(date +%Y%m%d-%H%M%S).txt"
exec 3>"$LOG"
say() { printf '%s\n' "$*" >&3; }
say "pass2 TS=$(date '+%Y-%m-%d %H:%M:%S %z')"

ALLOW=(
  "/c/Users/swq/wisp-303-steps"
  "/c/Users/swq/wisp-301-r2-tmp"
  "/c/Users/swq/wisp35v7mut"
  "/c/Users/swq/wisp-303-r1-fixed.go"
  "/c/Users/swq/wisp-303-v1-clone.log"
  "/c/Users/swq/wisp91-probe-read-sess91a7c3.txt"
  "/c/Users/swq/wisp91-probe-write-sess91a7c3.txt"
)
bad=0
for p in "${ALLOW[@]}"; do
  case "$p" in /c/Users/swq/wisp*) : ;; *) say "REFUSE=$p"; bad=$((bad+1)) ;; esac
done
say "guard_refusals=$bad"
[ "$bad" -ne 0 ] && { say "VERDICT=ABORT"; exit 7; }

ok=0; skipped=0
for p in "${ALLOW[@]}"; do
  if [ ! -e "$p" ]; then say "ABSENT=$p"; skipped=$((skipped+1)); continue; fi
  kind=file; [ -d "$p" ] && kind=dir
  sz=$(du -sm "$p" 2>/dev/null | awk '{print $1}')
  out=$(rm -rf "$p" 2>&1); rc=$?
  if [ $rc -eq 0 ]; then ok=$((ok+1)); say "OK=$p kind=$kind size_before=${sz}MB"; else say "FAIL rc=$rc :: $p :: $(printf '%s' "$out" | head -1 | cut -c1-140)"; fi
done
say "deleted_ok=$ok absent_or_skipped=$skipped"

# what is LEFT under the home dir that still carries the wisp name (for the record)
say ""
say "=== left after pass2 (name contains wisp, depth 1) ==="
ls -d /c/Users/swq/*wisp* 2>/dev/null | while read -r x; do say "LEFT=$x"; done
say "left_count=$(ls -d /c/Users/swq/*wisp* 2>/dev/null | wc -l)"
say "rc=0"
