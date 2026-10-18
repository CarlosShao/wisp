#!/bin/bash
# 227-v1 delivery-side self-proof gauges (fixed quoting; earlier 00-start-gauges.txt had two
# mis-escaped rulers -- the box gauge there printed 0 by MY escape bug, the correct readings
# are recorded here; the PowerShell-run originals at first contact printed 6/0, see verdict L3.)
cd '/d/work/workspace/projects plans/Wisp' || exit 2
out='.scratch/wisp/probes/227/v1/logs/18-selfproof-fixed.txt'
v=docs/evidence/s1/227-v1-verdict.md
t=.scratch/wisp/issues/227-guarded-write-delete-shapes-and-migration-on-read-are-unasserted.md
pat='[T][O][D][O]|[F][I][X][M][E]|[然]后填|[然]后补|待判|未判|占位'
{
  echo "== [slot ruler] on my verdict file (bracket form => this script's own bytes cannot self-match the pattern meaningfully)"
  echo "  hits = $(grep -cE "$pat" "$v" || true)  (want 0)"
  echo "== [positive control 1] same ruler over docs/reports/HANDOVER.md (want >0 => 尺活) = $(grep -cE "$pat" docs/reports/HANDOVER.md || true)"
  echo "== [positive control 2] same ruler over .scratch/wisp/issues/README.md (want >0) = $(grep -cE "$pat" .scratch/wisp/issues/README.md || true)"
  echo "== [negative control] known-absent token over verdict (want 0) = $(grep -cE 'ZZZnotatoken' "$v" || true)"
  echo
  echo "== [AC-box gauge] ticket 227, criteria-region scoped form (翻框归编排者; 六格必须仍未勾)"
  echo "  unchecked('- [ ] **AC#') = $(grep -c '^[[:space:]]*- \[ \] \*\*AC#' "$t" || true)"
  echo "  checked('- [x] **AC#')   = $(grep -c '^[[:space:]]*- \[x\] \*\*AC#' "$t" || true)"
  echo "  line numbers of all box lines:"
  grep -n '^[[:space:]]*- \[' "$t" | cut -c1-24
  echo
  echo "== [my own roster] commits with subject prefix 票 227 终裁腿 227-v1· (want 4-8)"
  git log --format='%h %ad %s' --date=format:'%H:%M' --grep='票 227 终裁腿 227-v1·' | cut -c1-90
  echo "  count = $(git log --format='%H' --grep='票 227 终裁腿 227-v1·' | wc -l)"
  paths=$(for c in $(git log --format='%H' --grep='票 227 终裁腿 227-v1·'); do git show --name-only --format= "$c"; done | sed '/^$/d' | sort -u)
  echo "  union paths:"; echo "$paths"
  echo "  .go hits (want 0) = $(echo "$paths" | grep -c '\.go$' || true)"
  echo "  outside my three write surfaces (want 0) = $(echo "$paths" | grep -vcE '^(docs/evidence/s1/227-v1-verdict\.md|\.scratch/wisp/probes/227/v1/|\.scratch/wisp/issues/227-)' || true)"
  echo "  [positive control] same .go ruler on the leg's b72d7391 = $(git show --name-only --format= b72d7391 | grep -c '\.go$' || true)  (want 2 => 尺⛔ 空转)"
  echo
  echo "== [zero-byte census] my write surface"
  echo "  probes/227/v1 files = $(find .scratch/wisp/probes/227/v1 -type f | wc -l); zero-byte = $(find .scratch/wisp/probes/227/v1 -type f -size 0 | wc -l) (want 0)"
  echo "  NOTE kept OUT of commits (too large, and it mirrors another party's file): logs/tree-ledger.txt = $(wc -c < .scratch/wisp/probes/227/v1/logs/tree-ledger.txt) bytes; the section I cite is excerpted to logs/14-a1033-extract.txt = $(wc -c < .scratch/wisp/probes/227/v1/logs/14-a1033-extract.txt) bytes, taken from HEAD bytes."
  echo "== [verdict size] lines=$(wc -l < "$v") bytes=$(wc -c < "$v") (cap 130)"
} > "$out" 2>&1
echo "written $(wc -c < "$out")"
