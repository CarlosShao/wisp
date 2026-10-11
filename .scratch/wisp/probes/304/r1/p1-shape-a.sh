#!/usr/bin/env bash
# 304-r1 p1: shape (a) (normal end) for BOTH spawn points, landed tree.
# Purpose: (i) the (a) reading itself, (ii) measure elapsed so the (b) cap is set
# with evidence (A850 s2#3: a cap above the run time measures a world that never happened).
set -u
R="D:/work/workspace/projects plans/Wisp"
D="$R/.scratch/wisp/probes/304/r1"
L="$D/logs/p1-shape-a.txt"
PS="$(cygpath -w "$R/.scratch/wisp/probes/orch/20261011/mockllm-procs.ps1")"
TMPX="/d/tmp/wisp-304-r1"
mkdir -p "$TMPX" "$D/logs"
cd "$R" || exit 1
q() { powershell -NoProfile -ExecutionPolicy Bypass -File "$PS" -Quick 2>&1 | tail -1; }

{ echo "=== p1 shape(a) @ $(date '+%F %T %z') HEAD=$(git rev-parse --short HEAD)"
  echo "ruler: mockllm-procs.ps1 -Quick (Get-CimInstance Win32_Process Name='mockllm.exe', attribution guard = cmdline matches *Wisp\\internal\\llm\\testdata\\golden*)"
  echo "measured against: the WORKTREE (my landed bytes), not a HEAD blob"
  echo "before: $(q)"; } > "$L"

for spec in "SP2-openaichat:./internal/llm/openaichat/:^TestMockllmGoldenByteIdenticalEvents\$" \
            "SP1-adaptertest:./internal/llm/openaichat/:^TestMockllmSharedFaultSuite\$"; do
  name="${spec%%:*}"; rest="${spec#*:}"; pkg="${rest%%:*}"; run="${rest#*:}"
  { echo "--- $name pkg=$pkg run=$run"
    t0=$(date +%s)
    GOFLAGS= go test "$pkg" -run "$run" -count=1 -timeout 300s -v > "$TMPX/a-$name.txt" 2>&1
    rc=$?
    echo "rc=$rc elapsed=$(( $(date +%s) - t0 ))s"
    echo "markers: RUN=$(grep -cE '^=== RUN' "$TMPX/a-$name.txt") PASS=$(grep -cE '^\s*--- PASS' "$TMPX/a-$name.txt") FAIL=$(grep -cE '^\s*--- FAIL' "$TMPX/a-$name.txt") panic=$(grep -c 'panic' "$TMPX/a-$name.txt") timedout=$(grep -c 'test timed out after' "$TMPX/a-$name.txt")"
    tail -3 "$TMPX/a-$name.txt"
    echo "after-$(echo $name | cut -c1-3): $(q)"; } >> "$L" 2>&1
done
echo "=== end $(date '+%T')" >> "$L"
cat "$L"
