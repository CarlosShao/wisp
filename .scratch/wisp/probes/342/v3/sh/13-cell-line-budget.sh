#!/bin/sh
# Per-cell line budget ruler for the two verdict files. A "cell" = the contiguous non-blank block that
# starts at a `- **AC#` bullet (markdown wraps each judgement into one long line, so counting physical
# lines from the bullet to the next blank line is the honest measure of the <=8-lines-per-cell cap).
cd "/d/work/workspace/projects plans/Wisp" || exit 1
for f in docs/evidence/s1/342-v3-verdict.md docs/evidence/s1/348-v2-verdict.md; do
  echo "===== $f  (lines=$(wc -l < "$f") bytes=$(wc -c < "$f") cap 130) ====="
  awk '
    /^- \*\*AC#/ { inblock=1; start=NR; n=0; label=substr($0,1,60); next }
    inblock {
      if ($0 ~ /^[[:space:]]*$/) { printf "  cell at line %d : %d lines  | %s\n", start, n, label; inblock=0; next }
      n++
    }
    END { if (inblock) printf "  cell at line %d : %d lines (EOF)  | %s\n", start, n, label }
  ' "$f"
  echo "  -- cells whose first line counts as line 1 --"
  awk 'BEGIN{cap=0} /^- \*\*AC#/ {c++; getline?0:0} END{print "  AC cells found = " c}' "$f"
done
echo "===== header/section blocks (## ...) are NOT cells; their sizes for the record ====="
for f in docs/evidence/s1/342-v3-verdict.md docs/evidence/s1/348-v2-verdict.md; do
  echo "--- $f"
  awk '/^## /{if(s){printf "  %d..%d = %d lines | %s\n", s, NR-1, NR-1-s+1, t}; s=NR; t=substr($0,1,50)} END{if(s){n=NR; printf "  %d..%d = %d lines | %s\n", s, n, n-s+1, t}}' "$f"
done
