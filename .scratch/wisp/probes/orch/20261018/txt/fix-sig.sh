#!/bin/sh
S="/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/orch/20261018/txt"
SIG=$(printf '\302\247')
for f in 342a1-specific.txt 342a1-head.txt 339v1-specific.txt 338v1-specific.txt 341r1-head.txt common-20261018b.txt; do
  p="$S/$f"
  perl -i -pe "s/\Q$SIG\E13\.1/13.1 /g; s/\Q$SIG\E1\`\/\`\Q$SIG\E3\`/H1\/H3\`/g; s/\Q$SIG\E/U+00A7/g;" "$p"
  n=$(grep -c "$SIG" "$p" || true)
  echo "$f sig=$n"
done
