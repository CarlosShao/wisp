#!/bin/sh
S="/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/orch/20261018/txt"
for L in 227r1 340a1; do
  cat "$S/$L-head.txt" "$S/common-20261018.txt" "$S/$L-specific.txt" > "$S/$L-dispatch.txt"
  echo "$L-dispatch.txt: $(wc -l < "$S/$L-dispatch.txt") lines / $(wc -c < "$S/$L-dispatch.txt") bytes"
done
