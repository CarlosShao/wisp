#!/usr/bin/env bash
# Owner 2026-10-13 08:0x: "C 盘里有没有测试遗留、没用的主动清理一下"。
# Scope = EXPLICIT allowlist only. Every path here is an out-of-repo scratch clone/mutant
# from ticket 303 (bisect platform, 10-10) or ticket 84 (snapshot, 09-21); both verdicts are
# delivered and in-repo, and every one of these is recreatable from git in minutes.
# NOT in scope (kept): wisp-test-gotmp roots (43 MB, some are cited evidence), wisp-33r* in
# %TEMP% (a leg is using them right now), the unattributable TestProbeStepByStep dir, all
# AppData wisp/wisp-dev app data.
set -u
LOG="D:/work/workspace/projects plans/Wisp/.scratch/wisp/probes/orch/20261013/logs/c-clean-wisp303-$(date +%Y%m%d-%H%M%S).txt"
exec 3>"$LOG"
say() { printf '%s\n' "$*" >&3; }

say "c-clean TS=$(date '+%Y-%m-%d %H:%M:%S %z') scope=explicit allowlist (ticket 303 / 84 out-of-repo scratch)"
say "free_before=$(df -m /c | awk 'NR==2{print $4}') MB free on C:"

ALLOW=(
  "/c/Users/swq/wisp-303-bisect"
  "/c/Users/swq/wisp-303-v1"
  "/c/Users/swq/wisp-303-v1-blob"
  "/c/Users/swq/wisp-303-v1-mut-base"
  "/c/Users/swq/wisp-303-v1-mut-M1"
  "/c/Users/swq/wisp-303-v1-mut-M2"
  "/c/Users/swq/wisp-303-v1-mut-M3"
  "/c/Users/swq/wisp-303-v1-mut-M4"
  "/c/Users/swq/wisp84-snap"
)

say ""
say "=== STEP 1 guard: every path must match the allowlist shape, else abort ==="
bad=0
for d in "${ALLOW[@]}"; do
  case "$d" in
    /c/Users/swq/wisp-303-*|/c/Users/swq/wisp84-snap) : ;;
    *) say "REFUSE=$d (outside the allowlist shape)"; bad=$((bad+1)) ;;
  esac
done
say "guard_refusals=$bad"
if [ "$bad" -ne 0 ]; then say "VERDICT=ABORT (a path left the allowlist shape; nothing deleted)"; exit 7; fi

say ""
say "=== STEP 2 in-use check (no process may hold these) ==="
say "wisp_test_exe=$(tasklist //FI 'IMAGENAME eq wisp.test.exe' 2>/dev/null | grep -ci 'wisp.test')"
say "go_exe=$(tasklist //FI 'IMAGENAME eq go.exe' 2>/dev/null | grep -ci 'go.exe')"

say ""
say "=== STEP 3 delete one by one ==="
ok=0; fail=0
for d in "${ALLOW[@]}"; do
  if [ ! -e "$d" ]; then say "ABSENT=$d"; continue; fi
  sz=$(du -sm "$d" 2>/dev/null | awk '{print $1}')
  out=$(rm -rf "$d" 2>&1); rc=$?
  if [ $rc -eq 0 ]; then
    ok=$((ok+1)); say "OK=$d size_before=${sz}MB"
  else
    fail=$((fail+1)); say "FAIL rc=$rc :: $d size_before=${sz}MB :: $(printf '%s' "$out" | head -1 | cut -c1-160)"
  fi
done
say "deleted_ok=$ok failed=$fail"

say ""
say "=== STEP 4 after ==="
for d in "${ALLOW[@]}"; do [ -e "$d" ] && say "STILL_THERE=$d"; done
say "free_after=$(df -m /c | awk 'NR==2{print $4}') MB free on C:"
say "leftovers_kept_check:"
say "  gotmp_root_total=$(ls -d /c/Users/swq/AppData/Local/wisp-test-gotmp/*/ 2>/dev/null | wc -l) roots"
say "  temp_wisp33=$(ls -d /c/Users/swq/AppData/Local/Temp/wisp-33* 2>/dev/null | wc -l)"
say "  unattributable=$(ls -d /c/Users/swq/AppData/Local/Temp/TestProbeStepByStep* 2>/dev/null | wc -l)"
say "rc=0"
