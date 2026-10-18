#!/bin/bash
# 227-v1 starting three gauges, re-run and recorded (dispatch: 起手三把尺先入库)
cd '/d/work/workspace/projects plans/Wisp' || exit 2
out='.scratch/wisp/probes/227/v1/logs/00-start-gauges.txt'
t=.scratch/wisp/issues/227-guarded-write-delete-shapes-and-migration-on-read-are-unasserted.md
{
  echo "[gauge 1] anchor = git log -1 (own reading, HEAD drifts between legs)"
  git log -1 --format='%H %ad %s' --date=iso-strict | cut -c1-120
  echo
  echo "[gauge 2] AC box count, criteria-region scoped form ONLY (dispatch: 全档计数会被票尾 bullet 污染)"
  echo "  ruler: grep -c '^[[:space:]]*- \\[ \\] \\*\\*AC#'  -> $t"
  echo "unchecked=$(grep -c '^[[:space:]]*- \\[ \\] \\*\\*AC#' "$t" || true)"
  echo "checked=$(grep -c '^[[:space:]]*- \\[x\\] \\*\\*AC#' "$t" || true)"
  echo "  line numbers of the six boxes:"
  grep -n '^[[:space:]]*- \[' "$t" | cut -c1-60
  echo "  ticket file lines=$(wc -l < "$t")  (票面 48 行＋腿的票尾 bullet)"
  echo
  echo "[gauge 3] 227-r1 roster: STRICT subject-prefix ruler (the leg's first sub-string ruler was self-invalidated)"
  echo "  ruler: git log --grep='票 227 落地腿 227-r1' (prefix of the subject)"
  git log --format='%h %ad %s' --date=format:'%H:%M' --grep='票 227 落地腿 227-r1' | cut -c1-100
  echo "count=$(git log --format='%H' --grep='票 227 落地腿 227-r1' | wc -l)"
  echo "  SAME ruler in the leg's own range feac8f77..HEAD = $(git log --format='%H' feac8f77..HEAD --grep='票 227 落地腿 227-r1' | wc -l)"
  echo "  the loose sub-string ruler (--grep='227-r1') over the same range = $(git log --format='%H' feac8f77..HEAD --grep='227-r1' | wc -l)  <- swallows OTHER legs' commits (必答 v)"
  echo
  echo "  per-commit name-only union + out-of-scope ruler (corrected form, 越界尺):"
  cs=$(git log --format='%H' --grep='票 227 落地腿 227-r1')
  union=$(for c in $cs; do git show --name-only --format= "$c"; done | sed '/^$/d' | sort -u)
  echo "$union"
  echo "union paths=$(echo "$union" | wc -l)  .go=$(echo "$union" | grep -c '\\.go$' || true)"
  oos=$(echo "$union" | grep -vcE '^(internal/config/|\.scratch/wisp/probes/227/r1/|\.scratch/wisp/issues/227-)' || true)
  echo "OUTSIDE the three allowed write surfaces (for the leg's own面) = $oos (want 0)"
  echo
  echo "  dispatch/head claim '14 笔' does NOT reproduce: strict ruler = 8; 14 = file count under probes/227/r1 = $(git ls-tree -r --name-only HEAD -- .scratch/wisp/probes/227/r1 | wc -l)"
} > "$out" 2>&1
echo "written $(wc -c < "$out") bytes"
