#!/bin/bash
# 341-r1: name and measure the FOUR junk files my runaway inline command created in the repo ROOT.
# ⛔ deleted (dispatch: temporary pieces are built, never deleted; and the repo root is NOT my write face, so
# cleaning them is the orchestrator's call, not mine). They are UNTRACKED => they are in no commit of mine;
# the registry below is the hand-off.
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 2
out=.scratch/wisp/probes/341/r1/logs/19-runaway-junk-files.txt
{
  echo "[how they happened] my inline printf used a bash DOUBLE-QUOTED format string containing BACKTICKS."
  echo "  Backticks there = command substitution; the substituted text (git status / git diff output lines, and"
  echo "  CJK prose fragments) was then EXECUTED by bash repeatedly. Lines containing '>' became shell"
  echo "  redirections, which created files in the CURRENT directory = the repo root. Rule violated:"
  echo "  dispatch shell block, quoted verbatim in spirit: double-quoted strings must not contain backticks or \$."
  echo
  echo "[1] every untracked entry in the repo ROOT that is not a directory, with size and mtime"
  for p in ./*; do
    [ -f "$p" ] || continue
    git ls-files --error-unmatch "$p" >/dev/null 2>&1 && { echo "    TRACKED, not mine to judge: $p"; continue; }
    printf '    UNTRACKED %s bytes=%s mtime=%s\n' "$p" "$(wc -c < "$p")" "$(stat -c '%y' "$p" | cut -c1-19)"
  done
  echo
  echo "[2] git status --porcelain for the root level only (want: exactly these four, nothing of mine tracked)"
  git status --porcelain -- . ':!:*/' | head -20
  echo
  echo "[3] content head of each (they are shell-output fragments, not data; recorded so the orchestrator can"
  echo "    confirm they contain nothing that looks like an instruction being followed)"
  for p in './⛔' './主料＝（26' './前一腿' './本件＝票'; do
    [ -f "$p" ] || { echo "    $p ABSENT now (something else moved it - report, ⛔ assume)"; continue; }
    echo "    ---- $p ($(wc -c < "$p") bytes) ----"
    head -c 400 "$p" | cat -v | sed 's/^/        /'
    echo
  done
  echo "[4] my commits do NOT contain them (strict roster, per-commit name-only)"
  for c in $(git log --perl-regexp --grep='^票 341 落地腿 341-r1·' --format='%H'); do
    echo "    $(git rev-parse --short "$c"): $(git show --name-only --format= "$c" | sed '/^$/d' | tr '\n' ' ')" | cut -c1-150
  done
  echo "    (any path in those lines outside internal/config/**_test.go, .scratch/wisp/probes/341/r1/,"
  echo "     .scratch/wisp/issues/341-* = 越界; the junk files must NOT appear)"
} > "$out" 2>&1
echo "written $(wc -c < "$out") bytes"
