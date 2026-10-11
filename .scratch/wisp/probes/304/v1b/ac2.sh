#!/bin/bash
# usage: ac2.sh TREE LOGPREFIX SHAPE PKG TEST
# SHAPE b = run compiled test binary with -test.timeout 3s (panic world, in-flight watch)
# SHAPE c = run compiled test binary with -test.timeout 600s, kill IT by named path once the child is seen
# World-construction evidence is landed in the same log; a run where watch_seen=0 is VOID (not a reading).
TREE="$1"; LOG="$2"; SHAPE="$3"; PKG="$4"; TEST="$5"
WT=$(cygpath -w "$TREE")
SCRIPTDIR=$(cd "$(dirname "$0")" && pwd)
RO=$(cygpath -w "$SCRIPTDIR/ro.ps1")
[ -f "$SCRIPTDIR/ro.ps1" ] || { echo "RO_MISSING=$SCRIPTDIR/ro.ps1"; exit 9; }
BIN="$TREE/ac2t.exe"
cd "$TREE" || exit 9
{
  echo "== ac2.sh shape=$SHAPE pkg=$PKG test=$TEST tree=$TREE =="
  echo "== compile: go test -c =="
  GOFLAGS= go test -c -o "$BIN" "$PKG" > "$LOG-compile.txt" 2>&1
  echo "compile_rc=$?"
  if [ ! -f "$BIN" ]; then echo "BIN_MISSING=1"; exit 8; fi
  echo "== pre roster =="
  powershell -NoProfile -ExecutionPolicy Bypass -File "$RO" -Tree "$WT" 2>&1
  echo "pre_rc=$?"
} > "$LOG-head.txt" 2>&1
(
  seen=0
  for i in $(seq 1 300); do
    out=$(powershell -NoProfile -ExecutionPolicy Bypass -File "$RO" -Tree "$WT" 2>&1)
    if echo "$out" | grep -q '^hit pid='; then
      echo "$out" | grep '^hit pid=' | head -3
      echo "watch_child_seen=1 t_poll=$i"
      seen=1
      break
    fi
    sleep 0.2
  done
  echo "watch_child_seen=$seen"
) > "$LOG-watch.txt" 2>&1 &
WPID=$!
if [ "$SHAPE" = "b" ]; then
  "$BIN" -test.run "^$TEST$" -test.v -test.timeout 3s > "$LOG-body.txt" 2>&1
  rc=$?
elif [ "$SHAPE" = "c" ]; then
  "$BIN" -test.run "^$TEST$" -test.v -test.timeout 600s > "$LOG-body.txt" 2>&1 &
  TPID=$!
  i=0
  while [ $i -lt 200 ]; do
    if grep -q 'watch_child_seen=1' "$LOG-watch.txt" 2>/dev/null; then break; fi
    sleep 0.2; i=$((i+1))
  done
  WT="$WT" powershell -NoProfile -Command '$p = @(Get-CimInstance Win32_Process | Where-Object { $_.Name -eq "ac2t.exe" -and $_.CommandLine -and $_.CommandLine.ToLower().Contains($env:WT.ToLower()) }); foreach($x in $p){ Write-Output ("kill_testbinary pid=" + $x.ProcessId + " ppid=" + $x.ParentProcessId + " parent_alive_at_kill=YES"); Stop-Process -Id $x.ProcessId -Force }; Write-Output ("kill_targets=" + $p.Count)' 2>&1
  wait $TPID; rc=$?
  sleep 1
else rc=7; fi
wait $WPID
{
  echo "go_test_rc=$rc"
  echo "== watch (world-construction evidence) =="; cat "$LOG-watch.txt"
  echo "== panic lines count =="; grep -c 'panic: test timed out after 3s' "$LOG-body.txt"
  echo "== panic/timeout context (first 6) =="; grep -n 'timed out after\|^panic' "$LOG-body.txt" | head -6
  echo "== post roster (own tree only) =="
  powershell -NoProfile -ExecutionPolicy Bypass -File "$RO" -Tree "$WT" 2>&1
  echo "post_rc=$?"
} >> "$LOG-head.txt" 2>&1
