#!/usr/bin/env bash
# 342-v1 mandatory question (i): the same five POSIX shapes, re-run as separate
# container invocations so each repeat is an independent measurement. ⛔ one number.
set -u
S=/d/tmp/wisp/20261018/342v1
for r in 1 2 3; do
  bash "$S/scripts/v342v1-race.sh" treeV "r$r" 1500 > "$S/logs/10-race-r$r-summary.txt" 2>&1
  echo "== repeat $r: $(grep -E 'V342V1 V[A-Z]+: attempts' "$S/logs/10-race-r$r-summary.txt" | wc -l) rows =="
  grep -E 'V342V1 V[A-Z]+: attempts' "$S/logs/10-race-r$r-summary.txt" || true
done
echo "RACE-REPEATS-DONE"
