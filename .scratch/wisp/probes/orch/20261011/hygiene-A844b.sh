#!/bin/sh
# hygiene ruler v2 (orchestrator self-check, read-only). v1 returned distinct_paths=0 because my
# ledger sections A840..A844 cite paths in SHORT form (`probes/309/a1/**`), not `.scratch/wisp/probes/...`.
# v1 stays on disk unedited; this file records the corrected scope.
cd "D:/work/workspace/projects plans/Wisp" || exit 2
out=.scratch/wisp/probes/orch/20261011/logs/hygiene-A844b.txt
mkdir -p .scratch/wisp/probes/orch/20261011/logs
{
  echo "date=$(date '+%Y-%m-%d %H:%M:%S %z')"
  echo "HEAD=$(git rev-parse HEAD)"
  echo "ruler=sed -n '/^## A840/,\$p' docs/reports/pending-and-issues.md | grep -oE 'probes/[A-Za-z0-9_./*-]+'  (short form, resolve to .scratch/wisp/<path>)"
  echo "v1_self_report=v1 used pattern '.scratch/wisp/probes/...' -> distinct_paths=0 = FALSE CLEAN, my ledger cites short form"
  sed -n '/^## A840/,$p' docs/reports/pending-and-issues.md \
    | grep -oE 'probes/[A-Za-z0-9_./*-]+' \
    | sed 's/[.,;)]$//' | sort -u > /tmp/hyg2-raw.txt
  echo "raw_matches=$(grep -c . /tmp/hyg2-raw.txt)"
  # split: exact file paths vs directory/glob forms (those ending in ** or a dir)
  miss_disk=0
  miss_head=0
  while read -r p; do
    full=".scratch/wisp/$p"
    case "$full" in
      *'**'*) kind="glob" ;;
      *) kind="path" ;;
    esac
    if [ "$kind" = "glob" ]; then
      base=$(dirname "$full")
      if [ ! -d "$base" ]; then echo "MISS_DISK $full (glob, base dir absent)"; miss_disk=$((miss_disk+1)); fi
      if ! git ls-tree -r --name-only HEAD -- "$base" | grep -q .; then echo "MISS_HEAD $full (glob, base untracked)"; miss_head=$((miss_head+1)); fi
    else
      if [ ! -e "$full" ]; then echo "MISS_DISK $full"; miss_disk=$((miss_disk+1)); fi
      if ! git cat-file -e "HEAD:$full" 2>/dev/null; then
        if [ -d "$full" ]; then
          if ! git ls-tree -r --name-only HEAD -- "$full" | grep -q .; then echo "MISS_HEAD $full (dir)"; miss_head=$((miss_head+1)); fi
        else
          echo "MISS_HEAD $full"; miss_head=$((miss_head+1))
        fi
      fi
    fi
  done < /tmp/hyg2-raw.txt
  echo "miss_disk=$miss_disk"
  echo "miss_head=$miss_head"
  echo "--- the matches, verbatim ---"
  cat /tmp/hyg2-raw.txt
} > "$out" 2>&1
echo "wrote $out"
wc -c "$out"
