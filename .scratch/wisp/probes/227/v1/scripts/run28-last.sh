#!/bin/bash
# 227-v1: last verification pass over the leg's own delivery (7 commits now).
cd '/d/work/workspace/projects plans/Wisp' || exit 2
out='.scratch/wisp/probes/227/v1/logs/28-last-verification.txt'
v=docs/evidence/s1/227-v1-verdict.md
t=.scratch/wisp/issues/227-guarded-write-delete-shapes-and-migration-on-read-are-unasserted.md
listed="173796c4 a13c87fa 88338d6a bc4a3fbb 65558bb1 761681de a1c6eeda"
derived=$(git log --format='%H' -- \
  "$v" \
  .scratch/wisp/probes/227/v1/logs/00-start-gauges.txt \
  .scratch/wisp/probes/227/v1/logs/03-mu1.txt \
  .scratch/wisp/probes/227/v1/logs/12-ac5-internal-all.txt \
  .scratch/wisp/probes/227/v1/logs/17-ac5-contention-remeasure.txt \
  .scratch/wisp/probes/227/v1/logs/21-unattributed-file.log \
  .scratch/wisp/probes/227/v1/logs/14-a1033-extract.txt \
  .scratch/wisp/probes/227/v1/scripts/run24-foreign-roster.sh \
  .scratch/wisp/probes/227/v1/logs/27-final-selfproof-corrected.txt)
{
  echo "anchor (git log -1, own reading) = $(git log -1 --format='%h %ad %s' --date=format:'%H:%M' | cut -c1-90)"
  echo
  echo "== roster cross-check (listed vs derived) =="
  echo "  listed=$(echo $listed | wc -w)  derived=$(echo "$derived" | grep -c .)"
  miss=0
  for h in $listed; do echo "$derived" | grep -q "^$h" || { echo "  listed-not-derived: $h"; miss=1; }; done
  for d in $derived; do ok=0; for h in $listed; do case "${d:0:8}" in "$h") ok=1;; esac; done; [ $ok -eq 0 ] && { echo "  derived-not-listed: ${d:0:8}"; miss=1; }; done
  echo "  two-way difference = $miss (want 0)"
  echo
  echo "== per-commit out-of-scope ruler (corrected form) =="
  for c in $listed; do
    n=$(git show --name-only --format= "$c" | sed '/^$/d' | grep -c . || true)
    o=$(git show --name-only --format= "$c" | sed '/^$/d' | grep -vcE '^(docs/evidence/s1/227-v1-verdict\.md|\.scratch/wisp/probes/227/v1/|\.scratch/wisp/issues/227-)' || true)
    g=$(git show --name-only --format= "$c" | sed '/^$/d' | grep -c '\.go$' || true)
    printf '  %s  paths=%-3s outside=%s .go=%s\n' "$c" "$n" "$o" "$g"
  done
  paths=$(for c in $listed; do git show --name-only --format= "$c" | sed '/^$/d'; done | sort -u)
  echo "  union paths=$(echo "$paths" | grep -c .)  .go=$(echo "$paths" | grep -c '\.go$' || true)  outside=$(echo "$paths" | grep -vcE '^(docs/evidence/s1/227-v1-verdict\.md|\.scratch/wisp/probes/227/v1/|\.scratch/wisp/issues/227-)' || true)"
  echo
  echo "== subject-prefix ruler (polluted; reported for the orchestrator, ⛔ this leg's key) =="
  echo "  total=$(git log --format='%H' --grep='票 227 终裁腿 227-v1·' | grep -c .)  mine=7  foreign=8"
  echo
  echo "== ticket + verdict gauges =="
  echo "  AC-unchecked=$(grep -c '^[[:space:]]*- \[ \] \*\*AC#' "$t" || true) AC-checked=$(grep -c '^[[:space:]]*- \[x\] \*\*AC#' "$t" || true) my-v1-bullets=$(grep -c '^- \*\*2026-10-18 非实现者终裁腿 227-v1' "$t" || true)"
  echo "  verdict lines=$(wc -l < "$v") (cap 130) bytes=$(wc -c < "$v")"
  echo "  slot ruler verdict=$(grep -cE '[T][O][D][O]|[F][I][X][M][E]|[占]位|[本]腿待补' "$v" || true) (want 0) | posctrl=$(grep -cE '[T][O][D][O]|[F][I][X][M][E]|[占]位|[本]腿待补' docs/reports/HANDOVER.md || true) (want >0)"
  echo "  zero-byte over my 7 commits' files = $(for c in $listed; do git show --name-only --format= "$c" | sed '/^$/d'; done | sort -u | while read -r p; do [ -f "$p" ] && [ ! -s "$p" ] && echo x; done | grep -c . || true) (want 0)"
  echo
  echo "== absolute paths of deliverables =="
  echo "  $(cygpath -w "$PWD")\\$(echo "$v" | tr '/' '\\')"
  echo "  $(cygpath -w "$PWD")\\.scratch\\wisp\\probes\\227\\v1  (logs\\ + scripts\\)"
} > "$out" 2>&1
echo "written $(wc -c < "$out")"
