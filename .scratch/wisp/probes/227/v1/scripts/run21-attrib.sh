#!/bin/bash
# 227-v1: compact attribution of the 32 reds in the deliv-side ./internal/... run, with the
# FIRST red sentence per case (so the raw ~1.08MB full logs need not be committed to be auditable).
cd '/d/work/workspace/projects plans/Wisp' || exit 2
f=.scratch/wisp/probes/227/v1/logs/12-ac5-internal-all-full.txt
out=.scratch/wisp/probes/227/v1/logs/20-ac5-red-attribution.txt
{
  echo "scope = deliv(HEAD export tree) ./internal/... -count=1 -timeout 30m ; four numbers 2077/1365/32/8 (rc=1)"
  echo "per red: name <- owning test file (last .go: token printed above the FAIL line) <- first 判定 sentence"
  echo "classification key: [台面]=export-tree gap (design//.git//.scratch absent, mklink/J blocked by long GOTMPDIR, SID/ACL seeded shapes) ; [争用]=ticket-named contention family"
  awk '
    /^[ ]+[A-Za-z0-9_.-]+\.go:[0-9]+:/ { last=$0 }
    /^=== RUN   / { runline=$2 }
    /^--- FAIL:/ {
      name=$3
      printf "%s <- %s\n", name, last
      if (last != "") last=""
    }
  ' "$f" > /tmp/redlist.txt
  n=0
  while IFS= read -r l; do
    n=$((n+1))
    name=$(echo "$l" | sed 's/ <- .*//')
    src=$(echo "$l" | sed 's/^.* <- //')
    first=$(grep -A3 "    $src" "$f" | grep -m1 -E "^ *$src" || true)
    echo "[$n] $l"
  done < /tmp/redlist.txt
  echo
  echo "== first red sentence per case (verbatim, truncated to 220 cols) =="
  grep -E '^\s+[a-zA-Z0-9_.-]+\.go:[0-9]+:' "$f" | head -70
  echo
  echo "== package-level FAIL lines =="
  grep -E '^FAIL\s+github' "$f"
} > "$out" 2>&1
echo "written $(wc -c < "$out") lines=$(wc -l < "$out")"
