#!/bin/bash
# 341-r1: damage control after ONE MISTAKE by this leg, named in full:
#   I inlined a `printf "..." ` whose format string was DOUBLE-quoted and contained BACKTICKS
#   (`git diff --numstat`, `Status:`, `-done`). Inside bash that is command substitution: git status output
#   lines were then EXECUTED as commands, repeatedly. It burned CPU and flooded stderr.
#   Nothing was written by it (the commands it "ran" were path names / not-found), but the state is verified
#   here rather than assumed: HEAD, my five commits, the worktree, and the ticket file.
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 2
out=.scratch/wisp/probes/341/r1/logs/18-mistake-and-statecheck.txt
T=.scratch/wisp/issues/341-four-hand-edit-shapes-still-unnailed-and-two-226-sentences-overclaim.md
{
  echo "[what happened] leg's own shell-discipline violation (dispatch shell rule: backticks inside double quotes)."
  echo "  job pwsh-f0ad4753d8984224 killed on request. Consequence checked below, ⛔ assumed."
  echo
  echo "[1] HEAD and my roster (strict subject-prefix ruler)"
  git log -1 --format='%H %s'
  echo "  my commits:"
  git log --perl-regexp --grep='^票 341 落地腿 341-r1·' --format='    %h %s' | sed 's/\(.\{100\}\).*/\1…/'
  echo "  count=$(git log --perl-regexp --grep='^票 341 落地腿 341-r1·' --format='%H' | wc -l)"
  echo
  echo "[2] did the runaway commit anything? (expect: NO msg-4 commit; the bullet file may be staged)"
  git status --porcelain | grep -vE '^\?\? ' | head -20
  echo "[end tracked-dirty list: only my own intended paths may appear]"
  echo
  echo "[3] the ticket file: appended bullet present once? line count? box ruler?"
  echo "  lines=$(wc -l < "$T") (was 43 at start, 44 after my append)"
  echo "  my bullet occurrences=$(grep -c '落地腿 341-r1 交件' "$T" || true)"
  echo "  unchecked=$(grep -c '^[[:space:]]*- \[ \] \*\*AC#' "$T" || true) checked=$(grep -c '^[[:space:]]*- \[x\] \*\*AC#' "$T" || true)"
  echo "  numstat vs HEAD:"
  git diff --numstat -- "$T"
  git diff --cached --numstat -- "$T"
  echo
  echo "[4] production files untouched all along (blob identity vs my start anchor)"
  A=98dbf20a35027a9290f52649d8ab2634646cb14e
  for p in internal/config/writeguard.go internal/config/loader.go internal/config/migrate.go; do
    echo "    $p anchor=$(git rev-parse "$A:$p" | cut -c1-12) index=$(git ls-files -s "$p" | awk '{print $2}' | cut -c1-12) worktree=$(git hash-object "$p" | cut -c1-12)"
  done
  echo
  echo "[5] my test file: worktree vs committed blob"
  f=internal/config/writeguard_unnailed_341_test.go
  mc=$(git log --perl-regexp --grep='^票 341 落地腿 341-r1·' --format='%H' -- "$f" | head -1)
  echo "    added in $(git rev-parse --short "$mc")  blob=$(git rev-parse "$mc:$f" | cut -c1-12)  worktree=$(git hash-object "$f" | cut -c1-12)"
  echo "    lines=$(wc -l < "$f") CR=$(tr -dc '\r' < "$f" | wc -c)"
  echo
  echo "[6] anything the runaway may have created by accident? list newest files in repo root and my face]"
  find . -maxdepth 1 -newermt '-15 minutes' -printf '    %p\n' 2>/dev/null | head -10
  find .scratch/wisp/probes/341 -newermt '-15 minutes' -printf '    %p\n' 2>/dev/null | head -20
  echo
  echo "[7] gates re-run clean after the incident (shared worktree, GOTMPDIR/TMPDIR unset first)"
  unset GOTMPDIR TMPDIR
  /d/work/base/gopath/bin/gofumpt.exe -l "$f"; echo "    gofumpt rc=$?"
  go vet ./internal/config/ 2>&1 | tail -3; echo "    vet rc=${PIPESTATUS[0]}"
  sh scripts/d22scan.sh 2>&1 | tail -1; echo "    d22scan rc=${PIPESTATUS[0]}"
  sh scripts/check-path-length-budget.sh 2>&1 | tail -1; echo "    pathlen rc=${PIPESTATUS[0]}"
  echo
  echo "[8] bench tree colours are unaffected (they were measured earlier); the shared worktree never ran go test"
  echo "    (dispatch: ⛔ 在共享工作树跑 go test／go build - this leg honoured it; the runaway ran git/status only)"
} > "$out" 2>&1
echo "written $(wc -c < "$out") bytes"
