#!/usr/bin/env bash
# 304-r1 paired.sh — AC#2 paired shots. $1 = fwd (landed tree) | rev (job step removed).
# Two spawn points x two bad-world shapes (b) timeout-panic / (c) external TerminateProcess.
# Every shot lands the evidence line that the bad world REALLY existed in that shot
# (A850 s2#3 trap: "0" measured over a world that never happened is not a reading).
set -u
DIR="${1:?fwd|rev}"
R="D:/work/workspace/projects plans/Wisp"
D="$R/.scratch/wisp/probes/304/r1"
L="$D/logs/paired-$DIR.txt"
PS="$(cygpath -w "$R/.scratch/wisp/probes/orch/20261011/mockllm-procs.ps1")"
TMPX="/d/tmp/wisp-304-r1"; EXE="$(cygpath -w "$TMPX/paired.test.exe")"
PKG="./internal/llm/openaichat/"
mkdir -p "$TMPX"
cd "$R" || exit 1
q() { powershell -NoProfile -ExecutionPolicy Bypass -File "$PS" -Quick 2>&1 | tail -1; }
roster() { powershell -NoProfile -ExecutionPolicy Bypass -File "$PS" 2>&1; }
nproc() { powershell -NoProfile -ExecutionPolicy Bypass -File "$PS" -Names "paired.test.exe" 2>&1 | tail -1; }
alive() { powershell -NoProfile -ExecutionPolicy Bypass -Command "if (Get-Process -Id $1 -ErrorAction SilentlyContinue) {'YES'} else {'NO'}" 2>&1 | tail -1; }
golden_count() { printf '%s' "$1" | grep -oE 'count_repo_golden=[0-9]+' | cut -d= -f2; }
clean_left() { # $1 = roster snapshot text file
  local ids; ids=$(grep -oE 'repo_golden_pids=\[[0-9,]*\]' "$1" | tr -d 'repo_golden_pids=[]' | sed 's/,$//')
  if [ -n "$ids" ]; then powershell -NoProfile -ExecutionPolicy Bypass -File "$PS" -CleanIds "$ids" 2>&1; else echo "nothing-to-clean"; fi; }

{ echo "=== paired dir=$DIR @ $(date '+%F %T %z') HEAD=$(git rev-parse --short HEAD)"
  echo "tree measured = WORKTREE of $R (dir=$DIR); landed bytes = $(git rev-parse --short HEAD) + my uncommitted edits"
  echo "ruler = .scratch/wisp/probes/orch/20261011/mockllm-procs.ps1 (Get-CimInstance Win32_Process Name='mockllm.exe', guard = cmdline like *Wisp\\internal\\llm\\testdata\\golden*)"
  echo "targets: SP2=mockllm_integ_test.go:73 via -run ^TestMockllmGoldenByteIdenticalEvents$ ; SP1=adaptertest/mockllm.go:72 via -run ^TestMockllmSharedFaultSuite$"
  echo "start: $(q)"; } > "$L"

GOFLAGS= go test -c -o "$EXE" "$PKG" >> "$L" 2>&1
echo "compile_test_exe_rc=$? exe_bytes=$(wc -c < "$TMPX/paired.test.exe" 2>/dev/null)" >> "$L"

shot_b() { # $1 label  $2 -run filter  $3 cap seconds
  local lbl="$1" run="$2" cap="$3" out="$TMPX/b-$lbl.txt" seen=0 seen_at="" rc
  { echo "--- SHAPE(b) $lbl : go test -timeout ${cap}s (test itself takes ~4-5s when it finishes, see p1)"
    echo "before: $(q)"; } >> "$L"
  GOFLAGS= go test "$PKG" -run "$run" -count=1 -timeout "${cap}s" -v > "$out" 2>&1 &
  local bg=$!
  for i in $(seq 1 40); do
    if ! kill -0 "$bg" 2>/dev/null; then break; fi
    local n; n=$(golden_count "$(q)")
    if [ "${n:-0}" -gt 0 ] && [ "$seen" = 0 ]; then seen=1; seen_at="$i"; echo "world_constructed: live mockllm.exe child seen at t=${i}.5s while the test process was STILL RUNNING ($n pid(s)) :: $(q)" >> "$L"; fi
    sleep 0.5
  done
  wait "$bg"; rc=$?
  { echo "rc=$rc  timeout_panic_lines=$(grep -c 'test timed out after' "$out" || true)  panic_lines=$(grep -c '^panic:' "$out" || true)"
    grep -E 'test timed out after|^panic:' "$out" | head -2
    echo "after: $(q)"; } >> "$L"
  roster > "$D/logs/roster-$DIR-b-$lbl.txt" 2>&1
  sed -E 's/^(pid=.*)/  \1/' "$D/logs/roster-$DIR-b-$lbl.txt" | head -8 >> "$L"
  { echo "VERDICT-b-$lbl: world_constructed=$seen (t=${seen_at:-NO}) ; $(q)"
    echo "clean: $(clean_left "$D/logs/roster-$DIR-b-$lbl.txt" | tr '\n' ' ')"; } >> "$L"
}

shot_c() { # $1 label  $2 -test.run filter
  local lbl="$1" run="$2" out="$TMPX/c-$lbl.txt" tp="" seen=0 rc
  { echo "--- SHAPE(c) $lbl : launch compiled test exe, wait for a LIVE child, then TerminateProcess the exe (Stop-Process -Force) only the pid this script launched"
    echo "before: $(q)"; } >> "$L"
  "$EXE" -test.run "$run" -test.timeout 300s > "$out" 2>&1 &
  local bg=$!
  sleep 1
  local nm; nm=$(nproc); echo "test-exe-roster: $nm" >> "$L"
  tp=$(printf '%s' "$nm" | grep -oE 'pids=\[[0-9,]*\]' | tr -d 'pids=[]' | cut -d, -f1)
  for i in $(seq 1 30); do
    local n; n=$(golden_count "$(q)")
    if [ "${n:-0}" -gt 0 ]; then seen=1; echo "world_constructed: live mockllm.exe child seen at t=${i}s, test exe pid=$tp alive=$(alive "$tp") :: $(q)" >> "$L"; break; fi
    sleep 1
  done
  if [ -n "$tp" ]; then
    powershell -NoProfile -ExecutionPolicy Bypass -Command "Stop-Process -Id $tp -Force" >> "$L" 2>&1
    echo "terminated: I launched pid=$tp and stopped exactly that pid; parent_confirmed_gone=$(alive "$tp")" >> "$L"
  else
    echo "NO test-exe pid found -> shape (c) NOT executed for $lbl" >> "$L"
  fi
  wait "$bg" 2>/dev/null; rc=$?
  { echo "exe_wait_rc=$rc (killed -> non-zero expected)"; echo "after: $(q)"; } >> "$L"
  roster > "$D/logs/roster-$DIR-c-$lbl.txt" 2>&1
  sed -E 's/^(pid=.*)/  \1/' "$D/logs/roster-$DIR-c-$lbl.txt" | head -8 >> "$L"
  { echo "VERDICT-c-$lbl: world_constructed=$seen ; $(q)"
    echo "clean: $(clean_left "$D/logs/roster-$DIR-c-$lbl.txt" | tr '\n' ' ')"; } >> "$L"
}

shot_b SP2 '^TestMockllmGoldenByteIdenticalEvents$' 3
shot_c SP2 '^TestMockllmGoldenByteIdenticalEvents$'
shot_b SP1 '^TestMockllmSharedFaultSuite$' 3
shot_c SP1 '^TestMockllmSharedFaultSuite$'
{ echo "=== end dir=$DIR $(date '+%T')  final: $(q)"; } >> "$L"
grep -cE "" "$L" | sed 's/^/log_lines=/'
grep -E "VERDICT|world_constructed|terminated|count_repo_golden" "$L" | cut -c1-170
