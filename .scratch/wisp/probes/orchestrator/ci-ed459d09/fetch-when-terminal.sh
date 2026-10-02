#!/usr/bin/env bash
# 取数程：等 CI run 到终态，然后把三枚 job 的原文日志逐枚落盘（不用 --log-failed，防少算）。
# 写点唯一＝.scratch/wisp/probes/orchestrator/ci-ed459d09/  只建不删。
set -u
RUN=36956193382
OUT=".scratch/wisp/probes/orchestrator/ci-ed459d09"
mkdir -p "$OUT"

echo "watch start $(date -Iseconds)" > "$OUT/watch.txt"
gh run watch "$RUN" --exit-status > "$OUT/watch.log" 2>&1
RC=$?
gh run view "$RUN" --json conclusion,status,updatedAt > "$OUT/terminal.json" 2>>"$OUT/watch.txt"
echo "watch rc=$RC end $(date -Iseconds)" >> "$OUT/watch.txt"

# 逐枚 job 取 id 与名字
gh run view "$RUN" --json jobs --jq '.jobs[] | [.databaseId, .name, .conclusion] | @tsv' > "$OUT/jobs.tsv" 2>>"$OUT/watch.txt"

while IFS=$'\t' read -r JID JNAME JCONC; do
  [ -z "${JID:-}" ] && continue
  safe=$(echo "$JNAME" | tr ' /()' '____' )
  gh run view "$RUN" --job "$JID" --log > "$OUT/job-$safe.log" 2>/dev/null
  printf '%s\t%s\t%s\tlines=%s\n' "$JID" "$JNAME" "$JCONC" "$(wc -l < "$OUT/job-$safe.log")" >> "$OUT/job-sizes.txt"
done < "$OUT/jobs.tsv"

echo "done $(date -Iseconds)" >> "$OUT/watch.txt"
