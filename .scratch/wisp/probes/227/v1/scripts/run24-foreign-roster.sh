#!/bin/bash
cd '/d/work/workspace/projects plans/Wisp' || exit 2
for c in d1da6386 daec56ac 73dac93c 69ce6600 57eb72cf 9ac4adef db6c6156 b578fa71; do
  echo "--- $c  $(git log -1 --format='%ad %s' --date=format:'%H:%M' $c | cut -c1-70)"
  git show --name-only --format= "$c" | sed '/^$/d' | sed 's/^/   /'
done
echo
echo "== union of their paths, and collisions with MY three write surfaces =="
paths=$(for c in d1da6386 daec56ac 73dac93c 69ce6600 57eb72cf 9ac4adef db6c6156 b578fa71; do git show --name-only --format= "$c" | sed '/^$/d'; done | sort -u)
echo "$paths" | sed 's/^/  /'
echo "their-path-count=$(echo "$paths" | grep -c .)"
echo "under .scratch/wisp/probes/227/v1 = $(echo "$paths" | grep -c '^\.scratch/wisp/probes/227/v1/' || true)"
echo "theirs touching docs/evidence/s1/227-v1-verdict.md = $(echo "$paths" | grep -c '227-v1-verdict.md' || true)"
echo "theirs touching ticket 227 = $(echo "$paths" | grep -c 'issues/227-' || true)"
echo "theirs touching .go = $(echo "$paths" | grep -c '\.go$' || true)"
