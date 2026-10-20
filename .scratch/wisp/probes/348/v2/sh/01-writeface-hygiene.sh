#!/bin/sh
# 342-v3 write-face hygiene ruler: line endings (LF), banned codepoint U+00A7 (and the look-alike
# interpunct U+00B7 counted separately so the two are never confused), zero-byte files, .out files,
# UTF-8 decodability (the fourth ruler that caught 342-r2's cut -c damage), and per-file sizes.
cd "$(git rev-parse --show-toplevel)" || exit 1
OUT=/d/tmp/wisp/20261020/342v3/logs
{
for root in .scratch/wisp/probes/342/v3 .scratch/wisp/probes/348/v2 docs/evidence/s1/342-v3-verdict.md docs/evidence/s1/348-v2-verdict.md; do
  echo "===== ROOT $root ====="
  if [ -d "$root" ]; then
    files=$(find "$root" -type f | sort)
  else
    files="$root"
  fi
  n=0; z=0; o=0; cr=0; a7=0; b7=0; bad=0
  for f in $files; do
    n=$((n+1))
    [ -s "$f" ] || { z=$((z+1)); echo "ZERO-BYTE $f"; }
    case "$f" in *.out) o=$((o+1)); echo "DOT-OUT $f";; esac
    c=$(tr -dc '\r' < "$f" | wc -c); [ "$c" -gt 0 ] && { cr=$((cr+c)); echo "CR=$c $f"; }
    a=$(perl -0777 -ne 'BEGIN{$t=0} while(/(\x{00A7})/g){$t++} END{print $t+0}' <(iconv -f UTF-8 -t UTF-8 "$f" 2>/dev/null) 2>/dev/null)
    iconv -f UTF-8 -t UTF-8 "$f" > /dev/null 2>&1 || { bad=$((bad+1)); echo "NOT-UTF8 $f"; }
  done
  echo "files=$n zero=$z dotout=$o cr_bytes=$cr not-utf8=$bad"
  echo "-- codepoint ruler (perl, :encoding(UTF-8) so it counts CODEPOINTS, not bytes) --"
  for f in $files; do
    perl -CSD -0777 -ne 'my $a = () = /(\N{U+00A7})/g; my $b = () = /(\N{U+00B7})/g; printf("  A7=%d B7=%d %s\n", $a, $b, $ARGV) if $a || $b;' "$f"
  done
  echo "-- total codepoint hits over the root (A7 must be 0) --"
  for f in $files; do cat "$f"; done | perl -CSD -0777 -ne 'my $a = () = /(\N{U+00A7})/g; my $b = () = /(\N{U+00B7})/g; print "A7=$a B7=$b\n";'
done
echo "===== verdict line/byte budget (dispatch: <=130 lines each, <=8 lines per cell) ====="
for f in docs/evidence/s1/342-v3-verdict.md docs/evidence/s1/348-v2-verdict.md; do
  [ -f "$f" ] && echo "$f bytes=$(wc -c < "$f") lines=$(wc -l < "$f")"
done
} > "$OUT/writeface-hygiene.txt" 2>&1
echo RC=$?
