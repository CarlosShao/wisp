#!/bin/sh
# Codepoint comparison across precedent verdict files (U+00A7 is the banned one; U+00B7 is the
# interpunct this house uses as a separator and appears in the accepted precedent files too).
cd "/d/work/workspace/projects plans/Wisp" || exit 1
for f in docs/evidence/s1/342-v2-verdict.md docs/evidence/s1/348-v1-verdict.md docs/evidence/s1/296-v2-verdict.md docs/evidence/s1/342-v3-verdict.md docs/evidence/s1/348-v2-verdict.md; do
  [ -f "$f" ] || { echo "ABSENT $f"; continue; }
  a7=$(perl -CSD -0777 -ne 'my $c = () = /(\N{U+00A7})/g; print $c;' "$f")
  b7=$(perl -CSD -0777 -ne 'my $c = () = /(\N{U+00B7})/g; print $c;' "$f")
  printf '%-46s A7=%s B7=%s lines=%s bytes=%s\n' "$f" "$a7" "$b7" "$(wc -l < "$f")" "$(wc -c < "$f")"
done
