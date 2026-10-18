#!/bin/bash
# 227-v1 FINAL self-proof, corrected roster key.
# BUG in run25 (named here, its file kept: 只建不删): I keyed my roster on the TICKET file, which
# carries other parties' history (32a70e2f 227-r1 + fac60ad4 + c774f8da + 7ffa9520) => their commits
# entered "my" union => 18 bogus "outside" hits. Correct key = paths only THIS leg ever wrote,
# cross-checked against the six hashes this leg recorded (⛔ 单枚尺自证).
cd '/d/work/workspace/projects plans/Wisp' || exit 2
out='.scratch/wisp/probes/227/v1/logs/27-final-selfproof-corrected.txt'
v=docs/evidence/s1/227-v1-verdict.md
t=.scratch/wisp/issues/227-guarded-write-delete-shapes-and-migration-on-read-are-unasserted.md
listed="173796c4 a13c87fa 88338d6a bc4a3fbb 65558bb1 761681de"
derived=$(git log --format='%H' -- \
  "$v" \
  .scratch/wisp/probes/227/v1/logs/00-start-gauges.txt \
  .scratch/wisp/probes/227/v1/logs/03-mu1.txt \
  .scratch/wisp/probes/227/v1/logs/12-ac5-internal-all.txt \
  .scratch/wisp/probes/227/v1/logs/17-ac5-contention-remeasure.txt \
  .scratch/wisp/probes/227/v1/logs/21-unattributed-file.log \
  .scratch/wisp/probes/227/v1/logs/14-a1033-extract.txt \
  .scratch/wisp/probes/227/v1/scripts/run24-foreign-roster.sh)
{
  echo "== [cross-check of the roster] listed (six hashes recorded at commit time) vs derived (path key) =="
  for h in $listed; do printf '  listed  %s  appearsInDerived=%s\n' "$h" "$(echo "$derived" | grep -c "^$h" || true)"; done
  echo "  derived count = $(echo "$derived" | grep -c .)"
  echo "  derived-but-not-listed (a miss would mean my listed roster is incomplete):"
  for d in $derived; do echo "$listed" | tr ' ' '\n' | grep -q "^${d:0:8}" || echo "     ${d:0:8} $(git log -1 --format='%s' "$d" | cut -c1-60)"; done
  echo "  listed-but-not-derived (would mean a listed commit touched ⛔ my own unique paths):"
  for h in $listed; do echo "$derived" | grep -q "^$h" || echo "     $h $(git log -1 --format='%s' "$h" | cut -c1-60)"; done
  echo
  echo "== [per-commit verification: ancestry, author, subject] =="
  for h in $listed; do
    anc=$(git merge-base --is-ancestor "$h" HEAD && echo yes || echo no)
    printf '  %s  %s  anc=%s  %s\n' "$h" "$(git log -1 --format='%an %ad' --date=format:'%H:%M' "$h")" "$anc" "$(git log -1 --format='%s' "$h" | cut -c1-52)"
  done
  echo "  count = $(echo $listed | wc -w) (cap 4-8)"
  echo
  paths=$(for c in $listed; do git show --name-only --format= "$c" | sed '/^$/d'; done | sort -u)
  echo "== [union paths over MY six commits] = $(echo "$paths" | grep -c .) 枚 =="
  echo "  .go hits (want 0) = $(echo "$paths" | grep -c '\.go$' || true)"
  echo "  outside the three allowed write surfaces (want 0) = $(echo "$paths" | grep -vcE '^(docs/evidence/s1/227-v1-verdict\.md|\.scratch/wisp/probes/227/v1/|\.scratch/wisp/issues/227-)' || true)"
  echo "  [any outside hits, listed]:"
  echo "$paths" | grep -vE '^(docs/evidence/s1/227-v1-verdict\.md|\.scratch/wisp/probes/227/v1/|\.scratch/wisp/issues/227-)' || echo "   (none)"
  echo "  [per-commit outside count]:"
  for c in $listed; do
    o=$(git show --name-only --format= "$c" | sed '/^$/d' | grep -vcE '^(docs/evidence/s1/227-v1-verdict\.md|\.scratch/wisp/probes/227/v1/|\.scratch/wisp/issues/227-)' || true)
    printf '   %s outside=%s\n' "$c" "$o"
  done
  echo "  [positive control A] out-of-scope ruler with the allowed set narrowed to internal/config/ only, over 227-r1's 8 commits = $(for c in $(git log --format='%H' --grep='票 227 落地腿 227-r1·'); do git show --name-only --format= "$c" | sed '/^$/d'; done | sort -u | grep -vcE '^internal/config/' || true) (want 15 => 尺能开火)"
  echo "  [positive control B] .go ruler over 227-r1's union = $(for c in $(git log --format='%H' --grep='票 227 落地腿 227-r1·'); do git show --name-only --format= "$c" | sed '/^$/d'; done | sort -u | grep -c '\.go$' || true) (want 2)"
  echo
  echo "== [ticket gauge] AC-unchecked=$(grep -c '^[[:space:]]*- \[ \] \*\*AC#' "$t" || true) AC-checked=$(grep -c '^[[:space:]]*- \[x\] \*\*AC#' "$t" || true) my-v1-bullet=$(grep -c '^- \*\*2026-10-18 非实现者终裁腿 227-v1' "$t" || true) leg-r1-bullet=$(grep -c '^- \*\*2026-10-18 落地腿 227-r1' "$t" || true) =="
  echo "== [verdict] lines=$(wc -l < "$v") bytes=$(wc -c < "$v") (cap 130 lines) =="
  echo "== [slot ruler] verdict=$(grep -cE '[T][O][D][O]|[F][I][X][M][E]|[占]位|[本]腿待补' "$v" || true) (want 0) | posctrl HANDOVER=$(grep -cE '[T][O][D][O]|[F][I][X][M][E]|[占]位|[本]腿待补' docs/reports/HANDOVER.md || true) (want >0) | negctrl=$(grep -cE '[Z][Z][Z]absent' "$v" || true) (want 0)"
  echo "== [zero-byte census over MY six commits' files] zero=$(for c in $listed; do git show --name-only --format= "$c" | sed '/^$/d'; done | sort -u | while read -r p; do [ -f "$p" ] && [ ! -s "$p" ] && echo "$p"; done | grep -c . || true) (want 0)"
  echo "   tracked files under my probes dir = $(git ls-tree -r --name-only HEAD -- .scratch/wisp/probes/227/v1 | wc -l), of which MINE = $(echo "$paths" | grep -c '\.scratch/wisp/probes/227/v1/')"
  echo "   NOTE: that tracked set also holds the DUPLICATE leg's 00-*..07-* (it committed those itself: d1da6386..b578fa71) = ⛔ mine, ⛔ touched by me."
  echo "   kept OUT (named, 只建不删): logs/tree-ledger.txt 5,233,854 B (whole-ledger mirror of another party's file; cited section = 14-a1033-extract.txt 9,400 B)"
  echo "   kept OUT, never committed (foreign in-flight): logs/00-start-anchors.txt 01-table-and-pristine.txt 02-ac2-decisive-shot.txt 08-answer-iv-sentences.txt"
} > "$out" 2>&1
echo "written $(wc -c < "$out")"
