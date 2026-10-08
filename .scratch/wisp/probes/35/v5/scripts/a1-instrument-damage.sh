#!/bin/bash
# 35-v5: does the (ii)/(iii) split damage any REUSABLE instrument?
# The v4 driver mutates the fixture BY LINE NUMBER (mkmut = awk NR==n). Its anchors are
# recorded in comments in that script; if the 35-r5 edit shifted those lines, re-running the
# v4 driver at HEAD silently produces a mutant that does not test what its name claims.
set -u
REPO="/d/work/workspace/projects plans/Wisp"
LOGS="$REPO/.scratch/wisp/probes/35/v5/logs"
cd "$REPO" || exit 1
FIX=cmd/wisp/panel_transport_35r2_test.go
BEFORE="$LOGS/before-fixture.copy"

{ echo "=== (a) v4 fixture-mutant line anchors: claimed vs actual, before vs after ==="
  echo "ruler = sed -n '<n>p' on each blob; LANDING = the printed today-line still is the target the script's comment names"
  while IFS='|' read -r n claim target; do
    b="$(sed -n "${n}p" "$BEFORE" | sed 's/^[[:space:]]*//')"
    a="$(sed -n "${n}p" "$FIX"  | sed 's/^[[:space:]]*//')"
    case "$a" in
      "$target"*) land=LANDS;;
      //*)        land="DEAD (hits a comment line)";;
      "")         land="DEAD (hits a blank line)";;
      *)          land="MISLANDS";;
    esac
    echo "line $n [$claim] wants \"$target\""
    echo "   before: $b"
    echo "   after : $a"
    echo "   verdict: $land"
    echo "   actual target now lives at line(s): $(grep -n -F "$target" "$FIX" | cut -d: -f1 | tr '\n' ' ')"
  done <<'ANCHORS'
164|v-set-gone|if _, exists := o.props[name]; exists && o.unwritable[name] {
1288|v-detector-gone|if recv != cw {
1290|v-no-panic|panic(&jsPanic{"TypeError: Illegal invocation:
1947|v-ma-world-swap|bf, spy, thrown := runShippedHookInOneWorld35r4(t, false)
1977|v-mb-world-swap|bf, spy, thrown := runShippedHookInOneWorld35r4(t, true)
1954|v-label-ma|t.Fatalf("M-A RED (receiver lost at the native exit)
2048|v-label-mb|t.Fatalf("M-B RED (guard's typeof half gone)
ANCHORS
  echo "rc=0"; } > "$LOGS/a08-v4-anchor-damage.txt" 2>&1

{ echo "=== (b) prod-hook mutant anchors (panel_host_windows.go) - prod unchanged, so expect ALL LAND ==="
  for n in 675 680 682; do echo "line $n : $(sed -n "${n}p" cmd/wisp/panel_host_windows.go | sed 's/^[[:space:]]*//')"; done
  echo "rc=0"; } > "$LOGS/a09-prod-anchors.txt" 2>&1

{ echo "=== (c) the repo-wide literal census the orchestrator restated as '19' ==="
  PAT="M-B"" RED"
  echo "A) tracked, ALL files, LINES            = $(git grep -h "$PAT" | wc -l)"
  echo "B) tracked, non-.go, LINES              = $(git grep -h "$PAT" -- . ':(exclude)*.go' | wc -l)"
  echo "C) tracked, non-.go, FILES              = $(git grep -l "$PAT" -- . ':(exclude)*.go' | wc -l)"
  echo "D) tracked, *.md, LINES                 = $(git grep -h "$PAT" -- '*.md' | wc -l)"
  echo "E) tracked, *.md, FILES                 = $(git grep -l "$PAT" -- '*.md' | wc -l)"
  echo "F) tracked, non-.go minus probes logs/  = $(git grep -h "$PAT" -- . ':(exclude)*.go' ':(exclude).scratch/wisp/probes/35/*/logs/**' | wc -l)"
  echo "G) .go files (expect 0)                 = $(git grep -h "$PAT" -- '*.go' | wc -l)"
  echo "--- per-file breakdown, non-.go ---"
  git grep -c "$PAT" -- . ':(exclude)*.go'
  echo "rc=0"; } > "$LOGS/a10-mb-red-census.txt" 2>&1

for f in a08-v4-anchor-damage a09-prod-anchors a10-mb-red-census; do echo "### $f"; cat "$LOGS/$f.txt"; echo; done
