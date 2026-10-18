#!/bin/bash
# 227-v1: final self-proof pass (after verdict edits). Rulers must be self-immune.
cd '/d/work/workspace/projects plans/Wisp' || exit 2
out='.scratch/wisp/probes/227/v1/logs/19-selfproof-final.txt'
v=docs/evidence/s1/227-v1-verdict.md
t=.scratch/wisp/issues/227-guarded-write-delete-shapes-and-migration-on-read-are-unasserted.md
pat='[T][O][D][O]|[F][I][X][M][E]|[占]位|[本]腿待补|[设]法填空'
{
  echo "== slot ruler on verdict (bracket form, self-immune) =="
  echo "  hits = $(grep -cE "$pat" "$v" || true)  (want 0)"
  echo "  posctrl HANDOVER.md = $(grep -cE "$pat" docs/reports/HANDOVER.md || true) (want >0)"
  echo "  negctrl '[Z][Z][Z]absent-token' = $(grep -cE '[Z][Z][Z]absent-token' "$v" || true) (want 0)"
  echo "  verdict lines=$(wc -l < "$v") bytes=$(wc -c < "$v") (cap 130)"
  echo "== AC box gauge (scoped form) =="
  echo "  unchecked=$(grep -c '^[[:space:]]*- \[ \] \*\*AC#' "$t" || true) checked=$(grep -c '^[[:space:]]*- \[x\] \*\*AC#' "$t" || true)"
  echo "== zero-byte census over my surface =="
  echo "  files=$(find .scratch/wisp/probes/227/v1 docs/evidence/s1/227-v1-verdict.md -type f | wc -l) zero-byte=$(find .scratch/wisp/probes/227/v1 docs/evidence/s1/227-v1-verdict.md -type f -size 0 | wc -l)"
  echo "== my roster (strict subject prefix '票 227 终裁腿 227-v1·') =="
  git log --format='%h %s' --grep='票 227 终裁腿 227-v1·' | cut -c1-90
  echo "  count=$(git log --format='%H' --grep='票 227 终裁腿 227-v1·' | wc -l) (want 4-8)"
  paths=$(for c in $(git log --format='%H' --grep='票 227 终裁腿 227-v1·'); do git show --name-only --format= "$c" | sed '/^$/d'; done | sort -u)
  echo "  union paths=$(echo "$paths" | grep -c . || true)"
  echo "  .go hits (want 0) = $(echo "$paths" | grep -cE '\.go$' || true)"
  echo "  outside allowed (want 0) = $(echo "$paths" | grep -vcE '^(docs/evidence/s1/227-v1-verdict\.md|\.scratch/wisp/probes/227/v1/|\.scratch/wisp/issues/227-)' || true)"
  echo "  [posctrl] same .go ruler over 227-r1's 8 commits union = $(for c in $(git log --format='%H' --grep='票 227 落地腿 227-r1·'); do git show --name-only --format= "$c" | sed '/^$/d'; done | sort -u | grep -cE '\.go$' || true)  (want 2 => 尺⛔ 空转)"
} > "$out" 2>&1
echo "written $(wc -c < "$out")"
