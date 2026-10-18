#!/bin/bash
# 341-r1: named corrections to two readings in 20-selfproof-final.txt. Both are RULER defects of mine,
# ⛔ findings about the delivery, and both are corrected here rather than left ambiguous for the收件 side.
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 2
P=.scratch/wisp/probes/341/r1
out="$P/logs/21-selfproof-corrections.txt"
{
  echo "[correction 1] ruler 3's case-name regex dropped the fifth case."
  echo "  I wrote 'HandDeletedBaseUrlOfUnknown'; the case is TestHandDeletedBaseUrlOfAnUnknownProviderStaysEmptyInTheFile"
  echo "  ('An' between 'of' and 'Unknown'), so grep printed 4 of 5 names. The named reading, re-run with the exact name:"
  for s in Adflt Bdtemp; do
    printf '    %s deliv: %s\n' "$s" "$(grep -E "^--- (PASS|FAIL): TestHandDeletedBaseUrlOfAnUnknownProviderStaysEmptyInTheFile" "$P/logs/05-color-$s-deliv.txt")"
    printf '    %s base : %s\n' "$s" "$(grep -cE 'TestHandDeletedBaseUrlOfAnUnknownProviderStaysEmptyInTheFile' "$P/logs/05-color-$s-base.txt")  (0 = the case did not exist at the baseline = the census reading)"
  done
  echo "  and all five together, from 15-carrier-table.txt (which used the full names): 5 PASS x 2 surfaces."
  echo
  echo "[correction 2] 'placeholder hits = 6 files' counted my OWN rulers, which contain the pattern text."
  echo "  Re-measured: the pattern search runs over files that themselves grep for TODO|FIXME|XXX."
  printf '    files hit by the placeholder pattern = '
  grep -rlE 'TODO|FIXME|XXX' "$P" 2>/dev/null | wc -l
  echo "    which ones, and what they are:"
  for f in $(grep -rlE 'TODO|FIXME|XXX' "$P" 2>/dev/null); do
    echo "      $f  -> $(grep -nE 'TODO|FIXME|XXX' "$f" | head -2 | cut -c1-120 | tr '\n' '|')"
  done
  echo "  delivery-side reading (⛔ the ruler text itself): my TEST FILE's placeholder hits = "
  printf '    %s\n' "$(grep -cE 'TODO|FIXME|XXX' internal/config/writeguard_unnailed_341_test.go || true)"
  echo
  echo "[correction 3] ruler 1's per-commit list is now five commits, not four: the delivery commit that"
  echo "  contains these two scripts + logs/20 + logs/21 is written AFTER 20-selfproof-final.txt was measured,"
  echo "  so that file cannot list itself (self-reference). Rerun on the收件 side:"
  echo '    git log --perl-regexp --grep="^票 341 落地腿 341-r1·" --format="%h %s"'
  echo "  and the 越界尺 per commit:"
  echo '    for c in $(git log --perl-regexp --grep="^票 341 落地腿 341-r1·" --format=%H); do'
  echo '      git show --name-only --format= $c | sed "/^$/d" | grep -vcE "^(internal/config/[^/]*_test.go|.scratch/wisp/probes/341/r1/|.scratch/wisp/issues/341-)"'
  echo '    done   (every line must read 0)'
  echo
  echo "[the push reading, restated with its ruler] git rev-list --count origin/dev..HEAD = $(git rev-list --count origin/dev..HEAD)"
  echo "  that number is the WHOLE branch's unpushed backlog (the machine owner is holding pushes), ⛔ mine:"
  echo "  of those, my commits = $(git log --perl-regexp --grep='^票 341 落地腿 341-r1·' --format='%H' | wc -l); this leg never ran git push."
} > "$out" 2>&1
echo "written $(wc -c < "$out") bytes"
