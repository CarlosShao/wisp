#!/bin/bash
# 227-v1 gauge A/B: internal/config four numbers on delivery tree (HEAD) and baseline tree (feac8f77)
root='D:/tmp/wisp/20261018/227v1'
out='/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/227/v1/logs'
mkdir -p "$root/gotmp"
export GOTMPDIR='D:/tmp/wisp/20261018/227v1/gotmp'
unset TMPDIR
for t in deliv base; do
  cd "$root/$t" || exit 2
  echo "TREE=$t anchor=$1 gofiles=$(find . -name '*.go' | wc -l)" > "$out/01-$t-config-head.txt"
  go test ./internal/config/ -count=1 -timeout 30m -v > "$out/01-$t-config.txt" 2>&1
  rc=$?
  RUN=$(grep -c '^=== RUN' "$out/01-$t-config.txt" || true)
  PASSTOP=$(grep -c '^--- PASS' "$out/01-$t-config.txt" || true)
  FAILTOP=$(grep -c '^--- FAIL' "$out/01-$t-config.txt" || true)
  SKIP=$(grep -c '^--- SKIP' "$out/01-$t-config.txt" || true)
  echo "tree=$t rc=$rc RUN=$RUN PASS_TOP=$PASSTOP FAIL_TOP=$FAILTOP SKIP=$SKIP" >> "$out/01-$t-config-head.txt"
  echo "---- named five (227 cases) + the 226 :393 case ----" >> "$out/01-$t-config-head.txt"
  grep -E '^--- (PASS|FAIL): Test(HandDeleted|NothingWritten|DeletedLine|GuardedWriteDoesNotRewrite)' "$out/01-$t-config.txt" >> "$out/01-$t-config-head.txt" || echo "NONE MATCHED" >> "$out/01-$t-config-head.txt"
  echo "$t DONE rc=$rc"
done
