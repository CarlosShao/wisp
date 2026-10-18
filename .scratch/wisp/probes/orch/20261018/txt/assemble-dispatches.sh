#!/bin/sh
set -e
S="/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/orch/20261018/txt"
O17="/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/orch/20261017/txt"
for L in 336r1 166r1 322r2 337a1 228r3c; do
  cat "$S/$L-head.txt" "$S/common-20261018.txt" "$S/$L-specific.txt" > "$S/$L-dispatch.txt"
  ln=$(wc -l < "$S/$L-dispatch.txt")
  bn=$(wc -c < "$S/$L-dispatch.txt")
  echo "$L-dispatch.txt: $ln lines / $bn bytes"
done
cp "$O17/328v1-dispatch.txt" "$S/328v1-dispatch.txt"
ln=$(wc -l < "$S/328v1-dispatch.txt")
bn=$(wc -c < "$S/328v1-dispatch.txt")
echo "328v1-dispatch.txt (replay copy): $ln lines / $bn bytes"
