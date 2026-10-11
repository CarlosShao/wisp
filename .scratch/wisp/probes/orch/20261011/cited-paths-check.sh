#!/bin/sh
cd "D:/work/workspace/projects plans/Wisp" || exit 99
OUT=.scratch/wisp/probes/orch/20261011/logs/orch-cited-paths-check.txt
git show HEAD:docs/reports/pending-and-issues.md \
  | grep -oE 'probes/[A-Za-z0-9_./-]+\.(md|txt|sh)' \
  | sed 's|^|.scratch/wisp/|' | sort -u > /tmp/cited3.txt
total=$(wc -l < /tmp/cited3.txt)
at_head=0; on_disk_only=0; nowhere=0
: > "$OUT.roster"
while IFS= read -r p; do
  if git cat-file -e "HEAD:$p" 2>/dev/null; then
    at_head=$((at_head+1))
  elif [ -e "$p" ]; then
    on_disk_only=$((on_disk_only+1)); printf 'ON-DISK-NOT-HEAD %s B %s\n' "$(stat -c %s "$p")" "$p" >> "$OUT.roster"
  else
    nowhere=$((nowhere+1)); printf 'NOT-ANYWHERE %s\n' "$p" >> "$OUT.roster"
  fi
done < /tmp/cited3.txt
{
  echo "cited_total=$total"
  echo "at_head=$at_head"
  echo "on_disk_not_head=$on_disk_only"
  echo "nowhere=$nowhere"
  echo "--- roster ---"
  cat "$OUT.roster"
} > "$OUT"
cat "$OUT"
