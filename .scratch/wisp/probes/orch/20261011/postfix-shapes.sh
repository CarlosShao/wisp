#!/usr/bin/env bash
# 编排者独立复量：甲形落地后，形 (b)（go test -timeout panic）与形 (c)（外力 TerminateProcess）
# 这两发"坏世界"到底还留不留残留。每一发都同批落"坏世界真的被造出来"的证据（⛔ 拿没跑当绿）。
set -u
R="D:/work/workspace/projects plans/Wisp"
D="$R/.scratch/wisp/probes/orch/20261011"
L="$D/logs/orch-A853-postfix-shapes.txt"
PS1W="$(cygpath -w "$D/mockllm-procs.ps1")"
TMPX="/d/tmp/wisp-orch-20261011"
mkdir -p "$TMPX"
quick() { powershell -NoProfile -ExecutionPolicy Bypass -File "$PS1W" -Quick 2>&1 | tail -1; }

{ echo "=== postfix shape (b)+(c) @ $(date '+%F %T')  HEAD=$(git rev-parse --short HEAD)"
  echo "before: $(quick)"; } >> "$L"

# ---------- 形 (b) 第一枚：spawn point #1（internal/llm 整包，走 adaptertest.StartMockllm）
echo "--- (b)#1 go test ./internal/llm/ -count=1 -timeout 8s" >> "$L"
go test ./internal/llm/ -count=1 -timeout 8s > "$TMPX/b1.txt" 2>&1 &
BGPID=$!
maxb1=0; pidsb1=""
for i in $(seq 1 14); do
  q=$(quick); c=$(printf '%s' "$q" | grep -oE 'count_repo_golden=[0-9]+' | cut -d= -f2)
  p=$(printf '%s' "$q" | grep -oE 'pids=\[[0-9,]*\]' | sed 's/pids=\[//;s/\]//')
  [ "${c:-0}" -gt "$maxb1" ] && { maxb1=$c; pidsb1=$p; }
  sleep 1
  kill -0 $BGPID 2>/dev/null || break
done
wait $BGPID; rc_b1=$?
{ echo "(b)#1 rc=$rc_b1 max_child_seen_while_in_flight=$maxb1 pids=[$pidsb1]"
  echo "(b)#1 panic_fired=$(grep -c 'test timed out after' "$TMPX/b1.txt")  (>=1 ＝坏世界真被造出来)"
  sleep 2
  echo "(b)#1 after: $(quick)"; } >> "$L"

# ---------- 形 (b) 第二枚：spawn point #2（openaichat 本地 startMockllm，3s 帽对 4.1s 用例）
echo "--- (b)#2 go test ./internal/llm/openaichat/ -run '^TestMockllmGoldenByteIdenticalEvents\$' -timeout 3s" >> "$L"
go test ./internal/llm/openaichat/ -run '^TestMockllmGoldenByteIdenticalEvents$' -count=1 -timeout 3s -v > "$TMPX/b2.txt" 2>&1 &
BGPID=$!
maxb2=0; pidsb2=""
for i in $(seq 1 10); do
  q=$(quick); c=$(printf '%s' "$q" | grep -oE 'count_repo_golden=[0-9]+' | cut -d= -f2)
  p=$(printf '%s' "$q" | grep -oE 'pids=\[[0-9,]*\]' | sed 's/pids=\[//;s/\]//')
  [ "${c:-0}" -gt "$maxb2" ] && { maxb2=$c; pidsb2=$p; }
  sleep 1
  kill -0 $BGPID 2>/dev/null || break
done
wait $BGPID; rc_b2=$?
{ echo "(b)#2 rc=$rc_b2 max_child_seen_while_in_flight=$maxb2 pids=[$pidsb2]"
  echo "(b)#2 panic_fired=$(grep -c 'test timed out after' "$TMPX/b2.txt")"
  sleep 2
  echo "(b)#2 after: $(quick)"; } >> "$L"

# ---------- 形 (c)：编译测试二进制后真跑，等它把孩子起来，再外力 Stop-Process 那枚二进制
echo "--- (c) build test binaries" >> "$L"
go test -c -o "$TMPX/oc.test.exe" ./internal/llm/openaichat/ > "$TMPX/cbuild.txt" 2>&1; cb=$?
go test -c -o "$TMPX/llm.test.exe" ./internal/llm/ > "$TMPX/cbuild2.txt" 2>&1; cb2=$?
echo "(c) build_rc openaichat=$cb internal_llm=$cb2" >> "$L"
if [ "$cb" = "0" ]; then
  "$TMPX/oc.test.exe" -test.run '^TestMockllmGoldenByteIdenticalEvents$' -test.timeout=120s > "$TMPX/c1.txt" 2>&1 &
  TP=$!
  child=""; tries=0
  while [ $tries -lt 30 ] && [ -z "$child" ]; do
    tries=$((tries+1))
    child=$(powershell -NoProfile -ExecutionPolicy Bypass -Command "Get-CimInstance Win32_Process -Filter \"Name='mockllm.exe'\" | Where-Object { \$_.ParentProcessId -eq $TP } | ForEach-Object { \$_.ProcessId }" 2>/dev/null | tr -d '\r' | head -1)
    sleep 1
  done
  { echo "(c)#2 parent_test_pid=$TP child_seen=$child after_polls=$tries  (⛔ child ⇒ 这一发没起世界，读数作废)"
    powershell -NoProfile -ExecutionPolicy Bypass -Command "Stop-Process -Id $TP -Force" 2>&1
    sleep 3
    gone=$(powershell -NoProfile -ExecutionPolicy Bypass -Command "if (Get-Process -Id $child -ErrorAction SilentlyContinue) {'STILL-ALIVE'} else {'child-gone'}" 2>/dev/null | tr -d '\r')
    echo "(c)#2 child_state_after_external_kill=$gone"; echo "(c)#2 after: $(quick)"; } >> "$L"
fi
if [ "$cb2" = "0" ]; then
  "$TMPX/llm.test.exe" -test.timeout=120s > "$TMPX/c2.txt" 2>&1 &
  TP2=$!
  child2=""; tries=0
  while [ $tries -lt 40 ] && [ -z "$child2" ]; do
    tries=$((tries+1))
    child2=$(powershell -NoProfile -ExecutionPolicy Bypass -Command "Get-CimInstance Win32_Process -Filter \"Name='mockllm.exe'\" | Where-Object { \$_.ParentProcessId -eq $TP2 } | ForEach-Object { \$_.ProcessId }" 2>/dev/null | tr -d '\r' | head -1)
    sleep 1
  done
  { echo "(c)#1 parent_test_pid=$TP2 child_seen=$child2 after_polls=$tries"
    powershell -NoProfile -ExecutionPolicy Bypass -Command "Stop-Process -Id $TP2 -Force" 2>&1
    sleep 3
    gone2=$(powershell -NoProfile -ExecutionPolicy Bypass -Command "if (Get-Process -Id $child2 -ErrorAction SilentlyContinue) {'STILL-ALIVE'} else {'child-gone'}" 2>/dev/null | tr -d '\r')
    echo "(c)#1 child_state_after_external_kill=$gone2"; echo "(c)#1 after: $(quick)"; } >> "$L"
fi
{ echo "final: $(quick)"; powershell -NoProfile -ExecutionPolicy Bypass -File "$PS1W" 2>&1 | tail -3; } >> "$L"
echo "script_rc=0" >> "$L"
tail -28 "$L"
