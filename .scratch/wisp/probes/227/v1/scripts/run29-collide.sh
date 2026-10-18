#!/bin/bash
# 227-v1: collision check — did ANY of my seven commits carry a file the DUPLICATE leg authored?
cd '/d/work/workspace/projects plans/Wisp' || exit 2
out='.scratch/wisp/probes/227/v1/logs/29-collision-check.txt'
mine="173796c4 a13c87fa 88338d6a bc4a3fbb 65558bb1 761681de a1c6eeda"
foreign="d1da6386 daec56ac 73dac93c 69ce6600 57eb72cf 9ac4adef db6c6156 b578fa71"
{
  echo "== the duplicate leg's 8 commits, and the paths each added =="
  for c in $foreign; do
    printf '  %s  %s\n' "$c" "$(git log -1 --format='%ad %s' --date=format:'%H:%M' "$c" | cut -c1-70)"
    git show --name-only --format= "$c" | sed '/^$/d' | sed 's/^/      /'
  done
  echo
  echo "== intersection of THEIR path set with MY path set (want EMPTY) =="
  tp=$(for c in $foreign; do git show --name-only --format= "$c" | sed '/^$/d'; done | sort -u)
  mp=$(for c in $mine;   do git show --name-only --format= "$c" | sed '/^$/d'; done | sort -u)
  echo "  their paths = $(echo "$tp" | grep -c .) ; my paths = $(echo "$mp" | grep -c .)"
  inter=$(comm -12 <(echo "$tp") <(echo "$mp"))
  if [ -z "$inter" ]; then echo "  INTERSECTION = 0 枚 (want 0)"; else echo "  INTERSECTION:"; echo "$inter" | sed 's/^/    /'; fi
  echo
  echo "== same check by authorship per path (git log -1 --format=%h for each of my paths) =="
  coll=0
  for p in $mp; do
    last=$(git log -1 --format='%h' -- "$p")
    case " $mine " in *" $last "*) : ;; *) echo "    path $p last touched by $last (⛔ mine)"; coll=$((coll+1));; esac
  done
  echo "  paths whose last committer is ⛔ one of my seven = $coll (want 0)"
  echo
  echo "== and the two shared-surface risks named in the dispatch =="
  echo "  the ticket file: my append = one line-pair only; numstat vs its parent commit:"
  git diff --numstat 65558bb1^ 65558bb1 -- .scratch/wisp/issues/227-guarded-write-delete-shapes-and-migration-on-read-are-unasserted.md
  echo "  verdict file commits touching it (want only mine):"
  git log --format='    %h %s' -- docs/evidence/s1/227-v1-verdict.md | cut -c1-70
} > "$out" 2>&1
echo "written $(wc -c < "$out")"
