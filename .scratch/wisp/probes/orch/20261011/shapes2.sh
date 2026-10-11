#!/usr/bin/env bash
# AC#0 cell (2): three death shapes MEASURED, not read. Orchestrator run (Q-304-a1-A).
set -u
R="D:/work/workspace/projects plans/Wisp"
D="$R/.scratch/wisp/probes/orch/20261011"
L="$D/logs/orch-A850-shapes.txt"
PS1W="$(cygpath -w "$D/mockllm-procs.ps1")"
TMPX="/d/tmp/wisp-orch-20261011"
: > "$L"
cd "$R" || exit 1

snap() { # $1 = tag
  powershell -NoProfile -ExecutionPolicy Bypass -File "$PS1W" > "$D/logs/orch-A850-snap-$1.txt" 2>&1
  { echo "--- snap $1 @ $(date '+%T')"; cat "$D/logs/orch-A850-snap-$1.txt"; } >> "$L"
}
quick() { powershell -NoProfile -ExecutionPolicy Bypass -File "$PS1W" -Quick 2>&1 | tail -1; }
names() { powershell -NoProfile -ExecutionPolicy Bypass -File "$PS1W" -Names "$1" 2>&1 | tail -1; }
newpids() { # $1 earlier tag, $2 later tag  -> pids present later but not earlier
  comm -13 \
    <(grep -oE '^pid=[0-9]+' "$D/logs/orch-A850-snap-$1.txt" | sed 's/^pid=//' | sort -u) \
    <(grep -oE '^pid=[0-9]+' "$D/logs/orch-A850-snap-$2.txt" | sed 's/^pid=//' | sort -u) | tr '\n' ','
}

{ echo "=== RUN $(date '+%F %T %z')  HEAD=$(git rev-parse --short HEAD)  staged=$(git diff --cached --name-only | grep -c '' || true)"
  echo "ruler: mockllm.exe roster via Get-CimInstance; attribution guard = command line contains this repo's golden dir"
  echo "target test: ./internal/llm/openaichat/ -run ^TestMockllmGoldenByteIdenticalEvents\$ (spawns mockllm.exe at mockllm_integ_test.go:73, reaps at :106)"
} >> "$L"
snap base

# ---------- shape (a): test binary reaches its own end ----------
{ echo "=== SHAPE (a) normal end"; t0=$(date +%s)
  go test ./internal/llm/openaichat/ -run '^TestMockllmGoldenByteIdenticalEvents$' -count=1 -timeout 300s -v > "$TMPX/shape-a.txt" 2>&1
  rc=$?; echo "rc=$rc elapsed=$(( $(date +%s) - t0 ))s"
  echo "markers: $(grep -cE '^=== RUN' "$TMPX/shape-a.txt") RUN / $(grep -cE '^--- PASS' "$TMPX/shape-a.txt") PASS / $(grep -cE '^--- FAIL' "$TMPX/shape-a.txt") FAIL / panic=$(grep -c 'panic' "$TMPX/shape-a.txt")"
  tail -6 "$TMPX/shape-a.txt"
} >> "$L"
snap a-after
{ echo "shape_a_new_pids=[$(newpids base a-after)]"; echo "shape_a_after_$(quick)"; } >> "$L"

# ---------- shape (b): go test -timeout fires while a child is live ----------
echo "=== SHAPE (b) -timeout panic" >> "$L"
go test ./internal/llm/openaichat/ -run '^TestMockllmGoldenByteIdenticalEvents$' -count=1 -timeout 14s -v > "$TMPX/shape-b.txt" 2>&1 &
bg=$!
seen=0; seen_at=""
for i in $(seq 1 12); do
  sleep 1
  q=$(quick); n=$(printf '%s' "$q" | grep -oE 'count_repo_golden=[0-9]+' | cut -d= -f2)
  if [ "${n:-0}" -gt 0 ] && [ "$seen" = 0 ]; then seen=1; seen_at="$i"; echo "mid-run: world HAS a live child at t=${i}s :: $q" >> "$L"; fi
done
wait $bg; rc=$?
echo "rc=$rc (124/2 expected if the timeout fired)" >> "$L"
echo "markers: timeout_panic=$(grep -c 'test timed out after' "$TMPX/shape-b.txt") RUN=$(grep -cE '^=== RUN' "$TMPX/shape-b.txt") PASS=$(grep -cE '^--- PASS' "$TMPX/shape-b.txt") FAIL=$(grep -cE '^--- FAIL' "$TMPX/shape-b.txt")" >> "$L"
tail -4 "$TMPX/shape-b.txt" >> "$L"
echo "world_constructed_before_panic=$seen (t=${seen_at}no)" >> "$L"
snap b-after
{ echo "shape_b_new_pids=[$(newpids a-after b-after)]"; echo "shape_b_after_$(quick)"; } >> "$L"

# ---------- shape (c): external TerminateProcess of the test binary ----------
echo "=== SHAPE (c) external TerminateProcess" >> "$L"
go test -c -o "$(cygpath -w "$TMPX/oc.test.exe")" ./internal/llm/openaichat/ > "$TMPX/build-c.txt" 2>&1
echo "build_rc=$? exe_bytes=$(wc -c < "$TMPX/oc.test.exe" 2>/dev/null)" >> "$L"
"$(cygpath -w "$TMPX/oc.test.exe")" -test.run '^TestMockllmGoldenByteIdenticalEvents$' -test.timeout 300s > "$TMPX/shape-c.txt" 2>&1 &
sleep 2
nm=$(names "oc.test.exe")
echo "test-exe-roster: $nm" >> "$L"
tp=$(printf '%s' "$nm" | grep -oE 'pids=\[[0-9,]*\]' | tr -d 'pids=[]' | cut -d, -f1)
seen=0
for i in $(seq 1 12); do
  q=$(quick); n=$(printf '%s' "$q" | grep -oE 'count_repo_golden=[0-9]+' | cut -d= -f2)
  if [ "${n:-0}" -gt 0 ]; then seen=1; echo "mid-run: live child at t=${i}s :: $q" >> "$L"; break; fi
  sleep 1
done
if [ -n "${tp:-}" ]; then
  powershell -NoProfile -ExecutionPolicy Bypass -Command "Stop-Process -Id $tp -Force" >> "$L" 2>&1
  echo "terminated test exe pid=$tp (only the pid this script launched)" >> "$L"
else
  echo "NO test-exe pid found -> shape (c) not executed" >> "$L"
fi
sleep 2
snap c-after
{ echo "shape_c_new_pids=[$(newpids b-after c-after)]"; echo "shape_c_after_$(quick)"; echo "still-running test exes: $(names 'oc.test.exe')"; echo "markers-c: $(grep -cE '^--- PASS' "$TMPX/shape-c.txt") PASS lines in its own log"; } >> "$L"

# ---------- attribution-guarded cleanup of what MY runs left ----------
left=$(grep -oE '^pid=[0-9]+' "$D/logs/orch-A850-snap-c-after.txt" | sed 's/^pid=//' | tr '\n' ',' | sed 's/,$//')
echo "=== CLEAN (guards: name=mockllm.exe AND cmdline has repo golden dir AND parent PID gone) ids=[$left]" >> "$L"
[ -n "$left" ] && powershell -NoProfile -ExecutionPolicy Bypass -File "$PS1W" -CleanIds "$left" >> "$L" 2>&1
sleep 1
snap final
echo "=== end $(date '+%T')" >> "$L"
cat "$L"
