#!/usr/bin/env bash
# 181-c1 read-only git-detection census probe (ticket 181 + 186).
# READ-ONLY: it never writes to the repo, never runs any mutating git command,
# and only runs `git ... list` style readers. Output goes to THIS probe's own
# logs dir (never overwrites another ticket's readings -- see ledger A367).
set -u
cd "$(git rev-parse --show-toplevel)" || exit 1
OUT_DIR=".scratch/wisp/probes/181/c1/logs"
mkdir -p "$OUT_DIR"
LOG="$OUT_DIR/git-census.log"
: > "$LOG"
say() { printf '%s\n' "$*" >> "$LOG"; }

say "== head/time =="
date >> "$LOG"
git log -1 --format='%h %cd' >> "$LOG"

say "== ruler 1: does any Go file read git? (expect empty) =="
grep -rln 'exec\.Command("git")\|\.git/HEAD\|rev-parse' --include=*.go internal/ cmd/ | grep -v _test.go >> "$LOG"

say "== ruler 2/3: workspace switch handler + callers =="
grep -n "func RequestWorkspaceSwitch" internal/panel/workspace.go >> "$LOG"
grep -rn "RequestWorkspaceSwitch" --include=*.go internal/ cmd/ | grep -v _test.go >> "$LOG"

say "== ruler 4: write-surface gate (must be empty) =="
git status --porcelain -- internal/ cmd/ >> "$LOG"

say "== Q1: .git shape at repo root + HEAD content =="
test -d .git && say ".git IS DIR"
test -f .git && say ".git IS FILE"
say "HEAD bytes: $(cat .git/HEAD)"

say "== Q1/Q2: ground truth (read-only listing) =="
git worktree list --porcelain >> "$LOG"

say "== Q1: linked worktree .git files (file shape, gitdir: prefix) =="
git worktree list --porcelain | grep '^worktree ' | sed 's/^worktree //' | while IFS= read -r wt; do
  [ "$wt" = "$(pwd -W 2>/dev/null || pwd)" ] && continue
  win="$(printf '%s' "$wt" | sed 's#^\([A-Za-z]\):/#/\L\1/#')"
  say "path=$wt"
  test -f "$win/.git" && say "  .git IS FILE content=[$(cat "$win/.git")]"
  test -d "$win/.git" && say "  .git IS DIR"
done >> "$LOG"

say "== Q2: main repo .git/worktrees enumeration =="
for w in .git/worktrees/*/; do
  say "worktree-dir=$w HEAD=[$(cat "$w/HEAD")] gitdir=[$(cat "$w/gitdir")] commondir=[$(cat "$w/commondir")]"
  say "  entries: $(ls "$w" | tr '\n' ' ')"
done >> "$LOG"

say "== Q3: branch enumeration (recursive vs one-level) =="
say "one-level ls count = $(ls .git/refs/heads | wc -l)"
say "recursive find count = $(find .git/refs/heads -type f | wc -l)"
say "ground truth git branch --list count = $(git branch --list | wc -l)"
find .git/refs/heads -type f | sort >> "$LOG"
say "packed-refs present = $(test -f .git/packed-refs && echo YES || echo NO)"
say "refs/remotes files = $(find .git/refs/remotes -type f 2>/dev/null | wc -l)"
say "remote url lines in .git/config (count only, values never logged) = $(grep -c 'url = ' .git/config)"

say "== Q4: inbound hop consumers (expect zero production callers) =="
grep -rn "ParseComposerRequest" --include=*.go internal/ cmd/ | grep -v _test.go >> "$LOG"
grep -rn "\.HandleModeRequest(\|\.HandleWorkspaceRequest(\|HandleAttachmentAdd(" --include=*.go internal/ cmd/ | grep -v _test.go >> "$LOG"
grep -rn "MethodWorkspaceRequest" --include=*.go internal/ cmd/ | grep -v _test.go >> "$LOG"
grep -rn "WebMessage\|ReceiveMessage\|OnMessage" --include=*.go internal/ cmd/ | grep -v _test.go >> "$LOG"

say "== Q5: workspace-root ripple targets =="
grep -rn "SetWorkspaceRoot\|WorkspaceRoot()" --include=*.go internal/ | grep -v _test.go >> "$LOG"
grep -rn "NewSpiller(" --include=*.go internal/ cmd/ | grep -v _test.go >> "$LOG"
grep -rn "BudgetsFor" --include=*.go internal/ | grep -v _test.go >> "$LOG"
grep -rn "OpenScope\|MarkWithHostPath" --include=*.go internal/risk internal/tools | grep -v _test.go | head -8 >> "$LOG"

say "== done =="
