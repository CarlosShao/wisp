#!/bin/sh
S="/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/orch/20261018/txt"
SIG=$(printf '\302\247')
for f in 343r1-head.txt 343r1-specific.txt common-20261018c.txt; do
  n=$(grep -c "$SIG" "$S/$f" 2>/dev/null || true)
  echo "$f: sig=$n lines=$(wc -l < "$S/$f") bytes=$(wc -c < "$S/$f")"
done
cat "$S/343r1-head.txt" "$S/common-20261018c.txt" "$S/343r1-specific.txt" > "$S/343r1-dispatch.txt"
echo "343r1-dispatch: $(wc -l < "$S/343r1-dispatch.txt") lines / $(wc -c < "$S/343r1-dispatch.txt") bytes / sig=$(grep -c "$SIG" "$S/343r1-dispatch.txt" 2>/dev/null || true)"
