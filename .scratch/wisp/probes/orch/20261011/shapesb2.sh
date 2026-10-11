#!/usr/bin/env bash
# shape (b) RE-TAKEN: the first shot never fired the timeout (test finished in 4.1s under a 14s cap)
set -u
R="D:/work/workspace/projects plans/Wisp"
D="$R/.scratch/wisp/probes/orch/20261011"
L="$D/logs/orch-A850-shape-b2.txt"
PS1W="$(cygpath -w "$D/mockllm-procs.ps1")"
TMPX="/d/tmp/wisp-orch-20261011"
cd "$R" || exit 1
quick() { powershell -NoProfile -ExecutionPolicy Bypass -File "$PS1W" -Quick 2>&1 | tail -1; }

{ echo "=== shape (b) re-taken @ $(date '+%F %T')  HEAD=$(git rev-parse --short HEAD)"
  echo "cap set to 3s because the whole test needs ~4.1s (measured: ok/... 4.112s) => the panic must land mid-flight"
  echo "before: $(quick)"; } >> "$L"
go test ./internal/llm/openaichat/ -run '^TestMockllmGoldenByteIdenticalEvents$' -count=1 -timeout 3s -v > "$TMPX/shape-b2.txt" 2>&1
rc=$?
{ echo "rc=$rc"; echo "timeout_panic_lines=$(grep -c 'test timed out after' "$TMPX/shape-b2.txt")"; echo "panic_lines=$(grep -c '^panic:' "$TMPX/shape-b2.txt")"; } >> "$L"
grep -E 'test timed out|^panic:|^\s*--- (PASS|FAIL)|^FAIL|^ok' "$TMPX/shape-b2.txt" | head -12 >> "$L"
sleep 2
{ echo "after: $(quick)"; powershell -NoProfile -ExecutionPolicy Bypass -File "$PS1W" 2>&1; } >> "$L"
left=$(powershell -NoProfile -ExecutionPolicy Bypass -File "$PS1W" -Quick 2>&1 | grep -oE 'pids=\[[0-9,]*\]' | tr -d 'pids=[]')
{ echo "clean ids=[$left]"; } >> "$L"
[ -n "$left" ] && powershell -NoProfile -ExecutionPolicy Bypass -File "$PS1W" -CleanIds "$left" >> "$L" 2>&1
echo "final: $(quick)" >> "$L"
cat "$L"
