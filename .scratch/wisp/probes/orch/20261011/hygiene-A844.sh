#!/bin/sh
# writeface: fleet hygiene ruler (orchestrator self-check, read-only)
# Question: do files referenced by my most recent ledger sections A840..A844 exist on disk AND in HEAD?
cd "D:/work/workspace/projects plans/Wisp" || exit 2
out=.scratch/wisp/probes/orch/20261011/logs/hygiene-A844.txt
mkdir -p .scratch/wisp/probes/orch/20261011/logs
{
  echo "date=$(date '+%Y-%m-%d %H:%M:%S %z')"
  echo "HEAD=$(git rev-parse HEAD)"
  echo "ruler=ledger sections A840..A844 -> paths matching .scratch/wisp/probes/... ; checked two ways: test -f (disk) and git cat-file -e HEAD:<path> (committed)"
  sed -n '/^## A840/,$p' docs/reports/pending-and-issues.md \
    | grep -oE '\.scratch/wisp/probes/[A-Za-z0-9_./-]+' \
    | sed 's/[.,;)]$//' | sort -u > /tmp/hyg-paths.txt
  echo "distinct_paths=$(grep -c . /tmp/hyg-paths.txt)"
  miss_disk=0
  miss_head=0
  while read -r p; do
    if [ ! -e "$p" ]; then
      echo "MISS_DISK $p"
      miss_disk=$((miss_disk+1))
    fi
    if ! git cat-file -e "HEAD:$p" 2>/dev/null; then
      echo "MISS_HEAD $p"
      miss_head=$((miss_head+1))
    fi
  done < /tmp/hyg-paths.txt
  echo "miss_disk=$miss_disk"
  echo "miss_head=$miss_head"
  echo "note: MISS_HEAD on a directory path is expected only if git does not track dirs by that name"
} > "$out" 2>&1
echo "wrote $out"
wc -c "$out"
