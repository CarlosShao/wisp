#!/bin/bash
# 227-v1 delivery-side final self-proof (path-based roster: the subject-prefix ruler is
# polluted by the duplicate same-name leg, see logs/24-*.log).
cd '/d/work/workspace/projects plans/Wisp' || exit 2
out='.scratch/wisp/probes/227/v1/logs/26-final-selfproof.txt'
v=docs/evidence/s1/227-v1-verdict.md
t=.scratch/wisp/issues/227-guarded-write-delete-shapes-and-migration-on-read-are-unasserted.md
mine=$(git log --format='%H' -- "$v" "$t" .scratch/wisp/probes/227/v1/logs/00-start-gauges.txt .scratch/wisp/probes/227/v1/logs/03-mu1.txt .scratch/wisp/probes/227/v1/logs/12-ac5-internal-all.txt .scratch/wisp/probes/227/v1/logs/17-ac5-contention-remeasure.txt)
{
  echo "== [roster, path-based] commits touching MY unique files (⛔ subject-prefix ruler: 12 hits = 4 mine + 8 the duplicate leg)"
  git log --format='%h %ad %s' --date=format:'%H:%M' -- "$v" "$t" .scratch/wisp/probes/227/v1/logs/00-start-gauges.txt .scratch/wisp/probes/227/v1/logs/03-mu1.txt .scratch/wisp/probes/227/v1/logs/12-ac5-internal-all.txt .scratch/wisp/probes/227/v1/logs/17-ac5-contention-remeasure.txt | cut -c1-100
  echo "  count = $(echo "$mine" | wc -l)"
  echo
  echo "== [out-of-scope ruler, corrected form] union of paths over my commits =="
  paths=$(for c in $mine; do git show --name-only --format= "$c" | sed '/^$/d'; done | sort -u)
  echo "$paths" | sed 's/^/   /'
  echo "  paths = $(echo "$paths" | grep -c .)"
  echo "  .go hits (want 0) = $(echo "$paths" | grep -c '\.go$' || true)"
  echo "  outside three write surfaces (want 0) = $(echo "$paths" | grep -vcE '^(docs/evidence/s1/227-v1-verdict\.md|\.scratch/wisp/probes/227/v1/|\.scratch/wisp/issues/227-)' || true)"
  echo "  [positive control] same .go ruler over 227-r1's prefix roster = $(for c in $(git log --format='%H' --grep='票 227 落地腿 227-r1·'); do git show --name-only --format= "$c" | sed '/^$/d'; done | sort -u | grep -c '\.go$' || true) (want 2 => 尺⛔ 空转)"
  echo
  echo "== [verdict] $(wc -l < "$v") lines / $(wc -c < "$v") bytes (cap 130) =="
  echo "== [slot ruler, bracket form] verdict hits = $(grep -cE '[T][O][D][O]|[F][I][X][M][E]|[占]位|[本]腿待补' "$v" || true) (want 0)"
  echo "   [posctrl] HANDOVER.md = $(grep -cE '[T][O][D][O]|[F][I][X][M][E]|[占]位|[本]腿待补' docs/reports/HANDOVER.md || true) (want >0)"
  echo "   [negctrl] '[Z][Z][Z]absent' = $(grep -cE '[Z][Z][Z]absent' "$v" || true) (want 0)"
  echo
  echo "== [ticket gauge] AC-unchecked=$(grep -c '^[[:space:]]*- \[ \] \*\*AC#' "$t" || true) AC-checked=$(grep -c '^[[:space:]]*- \[x\] \*\*AC#' "$t" || true) ; my-v1-bullets=$(grep -c '^- \*\*2026-10-18 非实现者终裁腿 227-v1' "$t" || true) (want 1); leg-r1-bullets=$(grep -c '^- \*\*2026-10-18 落地腿 227-r1' "$t" || true) (want 1) =="
  echo
  echo "== [zero-byte census over my write surface] files=$(find .scratch/wisp/probes/227/v1 "$v" -type f | wc -l) zero-byte=$(find .scratch/wisp/probes/227/v1 "$v" -type f -size 0 | wc -l) (want 0) =="
  echo "   kept OUT (named in logs + verdict): logs/tree-ledger.txt = $(wc -c < .scratch/wisp/probes/227/v1/logs/tree-ledger.txt) bytes (whole-ledger mirror of another party's file; the cited section is 14-a1033-extract.txt)"
  echo "   kept OUT (never committed, foreign in-flight bytes): logs/00-start-anchors.txt 01-table-and-pristine.txt 02-ac2-decisive-shot.txt 08-answer-iv-sentences.txt  <- those four are the DUPLICATE leg's (it committed its own copies itself: d1da6386..b578fa71)"
  echo
  echo "== [absolute paths] verdict = $(pwd)/$v"
  echo "   probes dir = $(pwd)/.scratch/wisp/probes/227/v1"
} > "$out" 2>&1
echo "written $(wc -c < "$out")"
