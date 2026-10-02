BEGIN { FS = "\t" }
{
  line = $3; for (i = 4; i <= NF; i++) line = line "\t" $i;
  sub(/^\xef\xbb\xbf/, "", line);
  sub(/^[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9:.+-]+Z? ?/, "", line);
  gsub(/\r/, "", line);
  if (line ~ /^##\[group\]Run /) { inseg = (line ~ /scope=windows/); next }
  if (!inseg) next;
  if (line ~ /^--- FAIL: Test/) { n = line; sub(/^--- FAIL: Test/, "Test", n); sub(/ .*/, "", n); print "FAIL " n }
  if (line ~ /^--- SKIP: Test/) { n = line; sub(/^--- SKIP: Test/, "Test", n); sub(/ .*/, "", n); print "SKIP " n }
}
