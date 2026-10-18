#!/bin/bash
# 227-v1 delivery-side self-proof: empty-slot ruler on my verdict file, with a REAL positive control
# (a token-shaped pattern that must exist somewhere else), zero-byte census, ticket box gauge.
cd '/d/work/workspace/projects plans/Wisp' || exit 2
out='.scratch/wisp/probes/227/v1/logs/15-selfproof-slot-and-box.txt'
v=docs/evidence/s1/227-v1-verdict.md
t=.scratch/wisp/issues/227-guarded-write-delete-shapes-and-migration-on-read-are-unasserted.md
{
  # pattern written as bracket classes so THIS file's own bytes cannot self-match it
  pat='[T][O][D][O]|[F][I][X][M][E]|待'填'|占位'
  echo "== slot ruler (character-class form; this file's bytes cannot self-match) =="
  echo "  ruler: grep -cE \"\$pat\" $v   (量的是本判语件, 射程＝单文件)"
  echo "verdict slot hits = $(grep -cE "$pat" "$v" || true)  (want 0)"
  echo "== positive control: same ruler shape over files that DO carry such tokens =="
  echo "  docs/reports/HANDOVER.md hits = $(grep -cE "$pat" docs/reports/HANDOVER.md || true)  (want >0 => 尺活)"
  echo "  .scratch/wisp/issues/README.md hits = $(grep -cE "$pat" .scratch/wisp/issues/README.md || true)  (want >0)"
  echo "== negative control: known-absent token over the verdict (want 0) =="
  echo "  ZZZnotatoken hits = $(grep -cE 'ZZZnotatoken' "$v" || true)"
  echo "== verdict size (cap: <=130 行) =="
  echo "  lines=$(wc -l < "$v") bytes=$(wc -c < "$v")"
  echo "== zero-byte census over my write surface (want 0) =="
  echo "  probes/227/v1 files=$(find .scratch/wisp/probes/227/v1 -type f | wc -l)  zero-byte=$(find .scratch/wisp/probes/227/v1 -type f -size 0 | wc -l)"
  echo "  NOTE kept-OUT-of-commit: logs/tree-ledger.txt (5.2MB mirror of the ledger, another party's file; the extract of the"
  echo "       section I cite lives at logs/14-a1033-extract.txt, 9,400 bytes, same bytes taken from HEAD)."
  echo "== ticket 227 AC-box gauge, criteria-region scoped form (six boxes must stay unchecked; 翻框归编排者) =="
  echo "  unchecked = $(grep -c '^[[:space:]]*- \[ \] \*\*AC#' "$t" || true)   checked = $(grep -c '^[[:space:]]*- \[x\] \*\*AC#' "$t" || true)"
  echo "  (the tail bullet appended by 227-r1 is line $(grep -n '^\- \*\*2026-10-18 落地腿 227-r1' "$t" | cut -d: -f1) and carries no box => scoped count == whole-file count here; 全档尺本腿⛔ 用)"
  echo "== out-of-scope ruler over MY own commits (corrected form; run at delivery time) =="
  git log --format='%h %s' --grep='票 227 终裁腿 227-v1' | cat
  echo "  my-commit count = $(git log --format='%H' --grep='票 227 终裁腿 227-v1' | wc -l)"
  paths=$(for c in $(git log --format='%H' --grep='票 227 终裁腿 227-v1'); do git show --name-only --format= "$c"; done | sed '/^$/d' | sort -u)
  echo "$paths" | cat
  echo "  paths=$(echo "$paths" | wc -l)  .go-hits=$(echo "$paths" | grep -c '\.go$' || true)  (want 0)"
  echo "  outside-my-three-write-surfaces = $(echo "$paths" | grep -vcE '^(docs/evidence/s1/227-v1-verdict\.md|\.scratch/wisp/probes/227/v1/|\.scratch/wisp/issues/227-)' || true)  (want 0)"
  echo "  positive control: same shape on 227-r1's b72d7391 .go hits = $(git show --name-only --format= b72d7391 | grep -c '\.go$' || true)  (want 2 => 尺⛔ 空转)"
} > "$out" 2>&1
echo "written $(wc -c < "$out")"
