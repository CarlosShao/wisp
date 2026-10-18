#!/bin/bash
# 341-r1: roster/gate re-measure. WHY A SECOND RUN, named:
#   run12's roster ruler was `git log --grep='票 341 落地腿 341-r1'` and it swallowed the ORCHESTRATOR's own
#   commit 98dbf20a, whose subject contains the string "341-r1" (it lists all four legs). So the "越界=16"
#   reading there was MY ruler defect - exactly the leg-identity trap the dispatch's 三键 section warns about:
#   a subject substring is not a unique key. This pass uses the strict subject PREFIX of my own commits,
#   '票 341 落地腿 341-r1·' (the middle dot appears in my subjects and nowhere in the orchestrator's),
#   plus two positive controls proving the ruler is not silently empty.
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 2
out=.scratch/wisp/probes/341/r1/logs/12-roster-strict.txt
A=98dbf20a35027a9290f52649d8ab2634646cb14e
PFX='票 341 落地腿 341-r1·'
{
  echo "HEAD at measurement = $(git rev-parse --short HEAD)   my start anchor = $(git rev-parse --short $A)"
  echo "strict ruler = git log <anchor>..HEAD --format='%H' --grep=\"^$PFX\"  (--perl-regexp, anchored at subject start)"
  echo
  echo "[roster] my commits only:"
  cs=$(git log --perl-regexp --grep="^$PFX" --format='%H' "$A..HEAD" --reverse)
  echo "$cs" | while read -r c; do [ -n "$c" ] && echo "  $(git log -1 --format='%h %s' "$c")"; done
  echo "count=$(echo "$cs" | sed '/^$/d' | wc -l)"
  echo
  echo "[POSITIVE CONTROL 1] the SAME ruler with the loose substring (run12's form) matches the orchestrator's"
  echo "  dispatch commit too: count=$(git log --grep='票 341 落地腿 341-r1' --format='%H' | wc -l) vs strict $(echo "$cs" | sed '/^$/d' | wc -l)"
  git log --grep='票 341 落地腿 341-r1' --format='  loose-hit %h %s' | head -3
  echo
  echo "[POSITIVE CONTROL 2] strict ruler against another leg's prefix must find ITS commits (ruler fires):"
  echo "  322-r1 strict hits=$(git log --perl-regexp --grep='^票 322 落地腿 322-r1·' --format='%H' | wc -l)"
  echo
  echo "[per-commit name-only + numstat, strict roster]"
  for c in $cs; do
    echo "  --- $(git log -1 --format='%h' "$c") $(git log -1 --format='%s' "$c")"
    git show --name-only --format= "$c" | sed '/^$/d' | sed 's/^/        /'
    echo "      numstat:"
    git diff --numstat "$c^" "$c" | sed 's/^/        /'
    echo "      production .go paths changed = $(git diff --numstat "$c^" "$c" | grep -E '\.go$' | grep -vE '_test\.go$' | wc -l)  (want 0)"
  done
  echo
  echo "[union of my paths, strict roster] and the 越界 count:"
  union=$(for c in $cs; do git show --name-only --format= "$c"; done | sed '/^$/d' | sort -u)
  echo "$union" | sed 's/^/    /'
  echo "  outside (internal/config/**_test.go | probes/341/r1/** | issues/341-*) = $(echo "$union" | grep -vcE '^(internal/config/.*_test\.go$|\.scratch/wisp/probes/341/r1/|\.scratch/wisp/issues/341-)' || true)"
  echo
  echo "[the three forbidden production files, blob identity start-anchor vs HEAD]"
  for p in internal/config/writeguard.go internal/config/loader.go internal/config/migrate.go; do
    x=$(git rev-parse "$A:$p"); y=$(git rev-parse "HEAD:$p")
    echo "    $p same=$([ "$x" = "$y" ] && echo YES || echo NO) $x"
  done
  echo
  echo "[collision re-check: is anyone else landing on internal/config right now?]"
  echo "  files newer than 45 min:"
  find internal/config -name '*.go' -newermt '-45 minutes' -print | sed 's/^/    /'
  echo "  (my own new file writeguard_unnailed_341_test.go will appear there = mine)"
  echo "  worktree status of the face:"
  git status --porcelain -- internal/config cmd/wisp | sed 's/^/    /'
  echo "[end status]"
  echo
  echo "[ticket box ruler] unchecked/checked must stay 7/0:"
  t=.scratch/wisp/issues/341-four-hand-edit-shapes-still-unnailed-and-two-226-sentences-overclaim.md
  echo "  unchecked=$(grep -c '^[[:space:]]*- \[ \] \*\*AC#' "$t" || true) checked=$(grep -c '^[[:space:]]*- \[x\] \*\*AC#' "$t" || true) lines=$(wc -l < "$t")"
} > "$out" 2>&1
echo "written $(wc -c < "$out") bytes"
