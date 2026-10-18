#!/bin/bash
# 227-v1 AC#5 contention re-measure (ticket: internal/risk 的 TestResolvePerCallBudget 是争用型假红，安静复量为准)
set -u
root='D:/tmp/wisp/20261018/227v1'
out='/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/227/v1/logs'
export GOTMPDIR="$root/gotmp"; unset TMPDIR
{
  echo "quiet re-measure = targeted -run on the one differing name, serial, 3 rounds per tree"
  echo "CPU now: $(powershell.exe -NoProfile -Command 'Get-CimInstance Win32_Processor | Measure-Object -Property LoadPercentage -Average | Select -Expand Average' 2>/dev/null || echo unknown)"
  for tree in deliv base; do
    for i in 1 2 3; do
      cd "$root/$tree" || exit 2
      go test ./internal/risk/ -run 'TestResolvePerCallBudget' -count=1 -timeout 30m -v > "$out/17-quiet-$tree-$i.txt" 2>&1
      rc=$?
      echo "$tree round$i rc=$rc  $(grep -E '^--- (PASS|FAIL): TestResolvePerCallBudget' "$out/17-quiet-$tree-$i.txt" | head -1)"
    done
  done
  echo "== and the full-package context for the two trees (risk only, -run empty = whole pkg) =="
  for tree in deliv base; do
    cd "$root/$tree" || exit 2
    go test ./internal/risk/ -count=1 -timeout 30m > "$out/17-pkg-$tree.txt" 2>&1
    echo "$tree whole-internal/risk rc=$? -> $(tail -1 "$out/17-pkg-$tree.txt")"
  done
} > "$out/17-ac5-contention-remeasure.txt" 2>&1
echo "written $(wc -c < "$out/17-ac5-contention-remeasure.txt")"
