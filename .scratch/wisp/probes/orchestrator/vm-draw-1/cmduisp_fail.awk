BEGIN { FS = "\t"; inseg = 0 }
{
  line = $3; for (i = 4; i <= NF; i++) line = line "\t" $i;
  sub(/^\xef\xbb\xbf/, "", line);
  sub(/^[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9:.+-]+Z? ?/, "", line);
  gsub(/\r/, "", line);
  if (line ~ /^##\[group\]Run /) { inseg = (line ~ /wisp-cli-tests/); next }
  if (!inseg) next;
  if (line ~ /^--- (FAIL|SKIP): Test/) { t = line; sub(/^--- /, "", t); sub(/:.*/, "", t); n = line; sub(/^--- (FAIL|SKIP): Test/, "Test", n); sub(/ .*/, "", n); sub(/\(.*/, "", n); lvl = n; print t " " lvl }
}
