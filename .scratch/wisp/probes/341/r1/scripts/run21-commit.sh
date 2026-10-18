#!/bin/bash
# 341-r1: the final commit, done the RIGHT way this time (message pre-written to a file, literal pathspecs,
# single-quoted strings only, ⛔ backticks inside double quotes = the exact rule this leg broke once).
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 2
T=.scratch/wisp/issues/341-four-hand-edit-shapes-still-unnailed-and-two-226-sentences-overclaim.md
P=.scratch/wisp/probes/341/r1
git add -- \
  "$T" \
  "$P/report.md" \
  "$P/scripts/bullet-341-append.txt" \
  "$P/scripts/msg-4-bullet.txt" \
  "$P/scripts/msg-5-bullet-final.txt" \
  "$P/scripts/run19-mistakecheck.sh" \
  "$P/scripts/run20-junkcheck.sh" \
  "$P/logs/18-mistake-and-statecheck.txt" \
  "$P/logs/19-runaway-junk-files.txt"
git commit -F "$P/scripts/msg-5-bullet-final.txt" -- \
  "$T" \
  "$P/report.md" \
  "$P/scripts/bullet-341-append.txt" \
  "$P/scripts/msg-4-bullet.txt" \
  "$P/scripts/msg-5-bullet-final.txt" \
  "$P/scripts/run19-mistakecheck.sh" \
  "$P/scripts/run20-junkcheck.sh" \
  "$P/logs/18-mistake-and-statecheck.txt" \
  "$P/logs/19-runaway-junk-files.txt"
echo "commit rc=$?"
git show --name-only --format='%h %s' HEAD
echo '---- my face status (want empty) ----'
git status --porcelain -- "$P" "$T" internal/config
echo '[end]'
