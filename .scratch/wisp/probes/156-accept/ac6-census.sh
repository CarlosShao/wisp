#!/usr/bin/env bash
# accept-r1 AC#6 census: per-commit, path-based, git diff-tree (NOT `git show --format=''`).
# usage: bash ac6-census.sh
set -u
ANCHOR=7ca1130
ERA_START=9835d81
OUT=.scratch/wisp/probes/156-accept/ac6
mkdir -p "$OUT"

# 15 arms, verbatim from ticket face :49-:51
ARMS=(
  "docs/PLAN.md"
  "docs/specs"
  "internal/risk"
  "internal/panel"
  "internal/agent/approval"
  "internal/observe"
  "internal/observe/thresholds.go"
  "testdata/golden"
  "internal/llm/golden"
  "tools/d22scan/allowlist.txt"
  "scripts/slo-check.ps1"
  "tools/d22scan"
  "frontend"
  "design"
  "docs/reports/pending-and-issues.md"
  "docs/reports/HANDOVER.md"
)

# ticket implementation surface, by PATH not message
IMPL=(
  "cmd/wisp"
  "docs/evidence/s1/156-exited-asks-os-r1.md"
  "docs/evidence/s1/156-r4-gates-and-cleanup-r1.md"
  ".scratch/wisp/probes/156"
)

echo "== era roster (by path, not message): commits touching the ticket implementation surface"
: > "$OUT/impl-commits.txt"
for c in $(git rev-list --reverse "$ERA_START..$ANCHOR"); do
  hit=$(git diff-tree -r --no-commit-id --name-only "$c" -- "${IMPL[@]}" | wc -l)
  if [ "$hit" -gt 0 ]; then
    echo "$c $(git log -1 --format=%s $c | cut -c1-60) implfiles=$hit" >> "$OUT/impl-commits.txt"
  fi
done
cat "$OUT/impl-commits.txt"
echo "== count: $(wc -l < "$OUT/impl-commits.txt")"

echo
echo "== era size"
git rev-list --count "$ERA_START..$ANCHOR"

echo
echo "== per-commit AC6 arm hits for each implementation commit"
: > "$OUT/ac6-per-commit.txt"
while read -r c _rest; do
  [ -z "$c" ] && continue
  line="$c"
  total=0
  for i in "${!ARMS[@]}"; do
    n=$(git diff-tree -r --no-commit-id --name-only "$c" -- "${ARMS[$i]}" | wc -l)
    line="$line A$((i+1))=$n"
    total=$((total+n))
  done
  del=$(git diff-tree -r --no-commit-id --numstat "$c" | awk '{s+=$2} END{print s+0}')
  line="$line TOTALHITS=$total deletedlines=$del"
  echo "$line" >> "$OUT/ac6-per-commit.txt"
done < <(awk '{print $1}' "$OUT/impl-commits.txt")
cat "$OUT/ac6-per-commit.txt"

echo
echo "== go.mod / go.sum per implementation commit"
while read -r c _rest; do
  [ -z "$c" ] && continue
  echo "$c gomodgomsum=[$(git diff-tree -r --no-commit-id --name-only "$c" -- go.mod go.sum | tr '\n' ' ')]"
done < <(awk '{print $1}' "$OUT/impl-commits.txt")

echo
echo "== POSITIVE CONTROLS: same ruler on commits that DO touch arms (path-picked)"
for c in 1218192 191f0d6 5866c6f 5eb0f6b; do
  [ -z "$(git cat-file -t $c 2>/dev/null)" ] && { echo "$c MISSING"; continue; }
  echo "$c $(git log -1 --format=%s $c | cut -c1-40)"
  for i in "${!ARMS[@]}"; do
    n=$(git diff-tree -r --no-commit-id --name-only "$c" -- "${ARMS[$i]}" | wc -l)
    [ "$n" -gt 0 ] && echo "   ARM$((i+1))=${ARMS[$i]} hits=$n"
  done
  gm=$(git diff-tree -r --no-commit-id --name-only "$c" -- go.mod go.sum | wc -l)
  echo "   go.mod/go.sum hits=$gm"
done

echo
echo "== NEGATIVE CONTROL: same ruler on the anchor's own commit (touched only an evidence file)"
for i in "${!ARMS[@]}"; do
  n=$(git diff-tree -r --no-commit-id --name-only "$ANCHOR" -- "${ARMS[$i]}" | wc -l)
  echo "   ARM$((i+1))=${ARMS[$i]} hits=$n"
done
