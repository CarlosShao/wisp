#!/bin/sh
S="/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/orch/20261018/txt"
SIG=$(printf '\302\247')
cat "$S/341v1-head.txt" "$S/common-20261018c.txt" "$S/341v1-specific.txt" > "$S/341v1-dispatch.txt"
echo "341v1-dispatch: $(wc -l < "$S/341v1-dispatch.txt") lines / $(wc -c < "$S/341v1-dispatch.txt") bytes / sig=$(grep -c "$SIG" "$S/341v1-dispatch.txt" 2>/dev/null || true)"
for f in 341v1-head.txt 341v1-specific.txt common-20261018c.txt; do echo "$f sig=$(grep -c "$SIG" "$S/$f" 2>/dev/null || true)"; done
