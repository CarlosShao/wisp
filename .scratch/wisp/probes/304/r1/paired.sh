#!/usr/bin/env bash
# 304-r1 paired.sh — AC#2 paired shots. $1 = fwd (landed tree) | rev (job step removed).
# Two spawn points x two bad-world shapes: (b) go test -timeout panic / (c) external TerminateProcess.
# Every shot must land the evidence that the bad world REALLY existed in that shot,
# attributed by parent PID and by the test name inside the child's t.TempDir() path
# (A850 s2#3 trap: a "0" measured over a world that never happened is not a reading;
# and on a shared machine a bare count can be another leg's in-flight child).
set -u
DIR="${1:?fwd|rev}"
R="D:/work/workspace/projects plans/Wisp"
D="$R/.scratch/wisp/probes/304/r1"
L="$D/logs/paired-$DIR.txt"
PS="$(cygpath -w "$R/.scratch/wisp/probes/orch/20261011/mockllm-procs.ps1")"
WK="$(cygpath -w "$D/watch-kill.ps1")"
CD="$(cygpath -w "$D/cmdline-dump.ps1")"
TMPX="/d/tmp/wisp-304-r1"; EXE="$(cygpath -w "$TMPX/paired.test.exe")"
PKG="./internal/llm/openaichat/"
mkdir -p "$TMPX"
cd "$R" || exit 1
q() { powershell -NoProfile -ExecutionPolicy Bypass -File "$PS" -Quick 2>&1 | tail -1; }
roster() { powershell -NoProfile -ExecutionPolicy Bypass -File "$PS" 2>&1; }
tpid() { powershell -NoProfile -ExecutionPolicy Bypass -File "$PS" -Names "paired.test.exe" 2>&1 | tail -1; }
alive() { powershell -NoProfile -ExecutionPolicy Bypass -Command "if (Get-Process -Id $1 -ErrorAction SilentlyContinue) {'YES'} else {'NO'}" 2>&1 | tail -1; }
gcount() { printf '%s' "$1" | grep -oE 'count_repo_golden=[0-9]+' | cut -d= -f2; }
clean_left() { # $1 = roster snapshot file ; guards inside mockllm-procs.ps1 -CleanIds: name=mockllm.exe AND repo-golden cmdline AND parent gone
  local ids; ids=$(grep -oE 'repo_golden_pids=\[[0-9,]*\]' "$1" | tr -d 'repo_golden_pids=[]' | sed 's/,$//')
  if [ -n "$ids" ]; then powershell -NoProfile -ExecutionPolicy Bypass -File "$PS" -CleanIds "$ids" 2>&1; else echo "nothing-to-clean"; fi; }
attribute() { # $1 = tag : whose residue is it? the t.TempDir() first segment carries the TEST NAME (issue 304 lines 11-14)
  local f="$D/logs/cmdline-$DIR-$1.txt"
  powershell -NoProfile -ExecutionPolicy Bypass -File "$CD" > "$f" 2>&1
  sed -E 's/^/  /' "$f" | cut -c1-260 >> "$L"; }

{ echo "=== paired dir=$DIR @ $(date '+%F %T %z') HEAD=$(git rev-parse --short HEAD)"
  echo "tree measured = WORKTREE of $R (dir=$DIR); landed bytes = HEAD tree + this leg's edits"
  echo "rulers = .scratch/wisp/probes/orch/20261011/mockllm-procs.ps1 (Get-CimInstance Win32_Process Name='mockllm.exe', guard = cmdline like *Wisp\\internal\\llm\\testdata\\golden) + this leg's watch-kill.ps1 (observe-and-terminate inside one PowerShell process) + cmdline-dump.ps1 (attribution by the test name in the child's TempDir path)"
  echo "spawn points: SP2=internal/llm/openaichat/mockllm_integ_test.go:73 via -run ^TestMockllmGoldenByteIdenticalEvents$ ; SP1=internal/llm/adaptertest/mockllm.go:72 (the helper, reached from an openaichat test) via -run ^TestMockllmSharedFaultSuite$"
  echo "caps: -timeout 3s sized by p1 (SP2 4.924s test / 7s wall, SP1 3.798s test / 5s wall) => the panic has to land mid-flight"
  echo "start: $(q)"; } > "$L"

GOFLAGS= go test -c -o "$EXE" "$PKG" > "$TMPX/compile.txt" 2>&1
echo "compile_test_exe_rc=$? exe_bytes=$(wc -c < "$TMPX/paired.test.exe" 2>/dev/null)" >> "$L"

shot_b() { # $1 label  $2 -run filter  $3 cap seconds
  local lbl="$1" run="$2" cap="$3" seen=0 seen_at="" rc out qs n bppid
  out="$TMPX/b-$lbl.txt"
  { echo "--- SHAPE(b) $lbl : go test -timeout ${cap}s ; the same test reaches its own end at ~4-5s (see p1)"
    echo "before: $(q)"; } >> "$L"
  GOFLAGS= go test "$PKG" -run "$run" -count=1 -timeout "${cap}s" -v > "$out" 2>&1 &
  local bg=$!
  for i in $(seq 1 40); do
    kill -0 "$bg" 2>/dev/null || break
    qs=$(q); n=$(gcount "$qs")
    if [ "${n:-0}" -gt 0 ] && [ "$seen" = 0 ]; then
      seen=1; seen_at="$i"
      roster > "$D/logs/construct-$DIR-b-$lbl.txt" 2>&1
      bppid=$(head -1 "$D/logs/construct-$DIR-b-$lbl.txt" | grep -oE 'ppid=[0-9]+' | cut -d= -f2)
      { echo "world_constructed t=${i}x0.5s: go test still in flight (bash kill-0 held on wrapper pid=$bg) ; $qs"
        head -2 "$D/logs/construct-$DIR-b-$lbl.txt" | sed -E 's/^/  /'
        echo "  test-binary_pid=$bppid present_at_that_instant=$(alive "$bppid") (Windows pid => the presence check is valid; bash MSYS pids are not Windows pids)"; } >> "$L"
    fi
    sleep 0.5
  done
  wait "$bg"; rc=$?
  { echo "rc=$rc  timeout_panic_lines=$(grep -c 'test timed out after' "$out" || true)  panic_lines=$(grep -c '^panic:' "$out" || true)"
    grep -E 'test timed out after|^panic:' "$out" | head -2 | sed -E 's/^/  bad-world-evidence: /'
    echo "after: $(q)"; } >> "$L"
  roster > "$D/logs/roster-$DIR-b-$lbl.txt" 2>&1
  head -3 "$D/logs/roster-$DIR-b-$lbl.txt" | sed -E 's/^/  /' >> "$L"
  { echo "  attribution of the residue in this shot:"; attribute "b-$lbl"; } >> "$L"
  { echo "VERDICT-b-$lbl: world_constructed=$seen (t=${seen_at:-NEVER}) ; $(q)"
    echo "  clean: $(clean_left "$D/logs/roster-$DIR-b-$lbl.txt" | tr '\n' ' ')"; } >> "$L"
}

shot_c() { # $1 label  $2 -test.run filter
  local lbl="$1" run="$2" tp="" seen=0 termed=0 rc out nm w
  out="$TMPX/c-$lbl.txt"
  { echo "--- SHAPE(c) $lbl : run the compiled test exe, then watch-kill.ps1 (single PowerShell process) observes a live mockllm.exe child OF THAT PID and TerminateProcess-es exactly that pid"
    echo "before: $(q)"; } >> "$L"
  "$EXE" -test.run "$run" -test.timeout 300s > "$out" 2>&1 &
  local bg=$!
  sleep 0.5
  nm=$(tpid); echo "test-exe-roster: $nm (bash wrapper pid=$bg; Windows pids are the ones in pids=[])" >> "$L"
  tp=$(printf '%s' "$nm" | grep -oE 'pids=\[[0-9,]*\]' | tr -d 'pids=[]' | cut -d, -f1)
  if [ -n "$tp" ]; then
    w=$(powershell -NoProfile -ExecutionPolicy Bypass -File "$WK" -TestPid "$tp" -MaxSeconds 25 2>&1)
    printf '%s\n' "$w" > "$D/logs/watch-$DIR-c-$lbl.txt"
    printf '%s\n' "$w" | sed -E 's/^/  watch-kill: /' >> "$L"
    printf '%s' "$w" | grep -q 'world_constructed' && seen=1
    printf '%s' "$w" | grep -q 'terminate_issued' && termed=1
  else
    echo "  test exe pid not found -> shape (c) NOT EXECUTED for $lbl" >> "$L"
  fi
  wait "$bg" 2>/dev/null; rc=$?
  { echo "exe_wait_rc=$rc (TerminateProcess exit code is what Stop-Process -Force sets; not used as the judgement)"
    echo "after: $(q)"; } >> "$L"
  roster > "$D/logs/roster-$DIR-c-$lbl.txt" 2>&1
  head -3 "$D/logs/roster-$DIR-c-$lbl.txt" | sed -E 's/^/  /' >> "$L"
  { echo "  attribution of the residue in this shot:"; attribute "c-$lbl"; } >> "$L"
  { echo "VERDICT-c-$lbl: world_constructed=$seen terminate_issued=$termed target_pid=$tp ; $(q)"
    echo "  clean: $(clean_left "$D/logs/roster-$DIR-c-$lbl.txt" | tr '\n' ' ')"; } >> "$L"
}

shot_b SP2 '^TestMockllmGoldenByteIdenticalEvents$' 3
shot_c SP2 '^TestMockllmGoldenByteIdenticalEvents$'
shot_b SP1 '^TestMockllmSharedFaultSuite$' 3
shot_c SP1 '^TestMockllmSharedFaultSuite$'
{ echo "=== end dir=$DIR $(date '+%T')  final: $(q)"; } >> "$L"
grep -cE "" "$L" | sed 's/^/log_lines=/'
grep -E "VERDICT|world_constructed|terminate_|bad-world-evidence|NOT EXECUTED|compile_test|test-binary_pid|parent_present|tempdir_testname|count_repo_golden" "$L" | cut -c1-175
