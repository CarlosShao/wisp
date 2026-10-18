#!/bin/bash
# 341-r1: final state after a HUNG commit attempt, verified rather than assumed.
# What happened, named: my run24 script's `git commit -F` never returned (240s timeout -> moved to background,
# then killed). Meanwhile the orchestrator landed two commits of its own, so MY files sat STAGED in the shared
# index - the exact reverse of the A1027 trap (a later commit by another leg with a pathspec derived from the
# index would have swept them in). I retried with GIT_EDITOR=true and literal pathspecs; that is 46c920c8.
# Gauges below answer: did the hung attempt commit anything? is my roster exactly five? is the index clean of mine?
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 2
export GIT_PAGER=cat
P=.scratch/wisp/probes/341/r1
out="$P/logs/22-final-state-after-hang.txt"
{
  echo "[1] roster: strict subject-prefix ruler (⛔ substring - that one ate the orchestrator's dispatch commit)"
  git log --perl-regexp --grep='^票 341 落地腿 341-r1·' --format='%h %ad %s' --date=format:'%H:%M' | cut -c1-96
  echo "count=$(git log --perl-regexp --grep='^票 341 落地腿 341-r1·' --format='%H' | wc -l)"
  echo
  echo "[2] duplicate-commit check (the hung attempt might have committed once silently):"
  echo "  commits touching $P/logs/20-selfproof-final.txt = $(git log --format='%h' -- "$P/logs/20-selfproof-final.txt" | wc -l)  (want 1)"
  echo "  subjects: $(git log --format='%h %s' -- "$P/logs/20-selfproof-final.txt" | cut -c1-60 | tr '\n' '|')"
  echo
  echo "[3] is anything of mine still staged in the shared index? (want: nothing)"
  git diff --cached --name-only | sed 's/^/    STAGED /'
  echo "[end staged list]"
  echo
  echo "[4] per-commit out-of-scope count and production-go count, all five commits"
  for c in $(git log --perl-regexp --grep='^票 341 落地腿 341-r1·' --format='%H'); do
    oos=$(git show --name-only --format= "$c" | sed '/^$/d' | grep -vcE '^(internal/config/[^/]*_test\.go|\.scratch/wisp/probes/341/r1/|\.scratch/wisp/issues/341-)' || true)
    prod=$(git diff --numstat "$c^" "$c" | grep -E '\.go$' | grep -vE '_test\.go$' | wc -l)
    echo "    $(git rev-parse --short "$c")  out-of-scope=$oos  prod-go=$prod  files=$(git show --name-only --format= "$c" | sed '/^$/d' | wc -l)"
  done
  echo
  echo "[5] the three forbidden production blobs: anchor 98dbf20a vs HEAD"
  for p in internal/config/writeguard.go internal/config/loader.go internal/config/migrate.go; do
    echo "    $p  $(git rev-parse 98dbf20a:$p | cut -c1-12) vs $(git rev-parse HEAD:$p | cut -c1-12)"
  done
  echo
  echo "[6] ticket discipline: 44 lines / one bullet / 7 unchecked 0 checked / zero deleted lines"
  T=.scratch/wisp/issues/341-four-hand-edit-shapes-still-unnailed-and-two-226-sentences-overclaim.md
  echo "    lines=$(wc -l < "$T") bullet=$(grep -c '落地腿 341-r1 交件' "$T" || true) unchecked=$(grep -c '^[[:space:]]*- \[ \] \*\*AC#' "$T" || true) checked=$(grep -c '^[[:space:]]*- \[x\] \*\*AC#' "$T" || true)"
  echo "    deletions in the ticket across my range = $(git diff 98dbf20a..HEAD --numstat -- "$T" | awk '{print $2}')"
  echo
  echo "[7] my new test file, final: blob/lines/CR/funcs + gate"
  f=internal/config/writeguard_unnailed_341_test.go
  echo "    blob=$(git rev-parse HEAD:$f | cut -c1-12) lines=$(wc -l < "$f") CR=$(tr -dc '\r' < "$f" | wc -c) funcs=$(grep -c '^func Test' "$f")"
  unset GOTMPDIR TMPDIR
  /d/work/base/gopath/bin/gofumpt.exe -l "$f"; echo "    gofumpt rc=$?"
  go vet ./internal/config/ 2>&1 | tail -2; echo "    vet rc=${PIPESTATUS[0]}"
  echo
  echo "[8] push discipline: this leg ran git push 0 times; branch backlog ahead of origin/dev = $(git rev-list --count origin/dev..HEAD) (⛔ mine; owner is holding pushes)"
  echo "[9] bench trees under D:/tmp/wisp/20261018/341r1 (built, never deleted):"
  ls -1 /d/tmp/wisp/20261018/341r1 | sed 's/^/    /'
} > "$out" 2>&1
echo "written $(wc -c < "$out") bytes"
