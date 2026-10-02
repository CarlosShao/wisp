BEGIN { FS = "\t" }
{
  line = $3; for (i = 4; i <= NF; i++) line = line "\t" $i;
  sub(/^\xef\xbb\xbf/, "", line);
  sub(/^[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9:.+-]+Z? ?/, "", line);
  gsub(/\r/, "", line);
  if (line ~ /^##\[group\]Run /) { inseg = (line ~ /wisp-cli-tests/); next }
  if (!inseg) next;
  if (line ~ /^=== RUN   TestTicket223HandEditedFsLooseningCostsAnL2Card$/) { found = 1; exit }
  if (line ~ /^=== RUN   Test[^\/]*$/) print line;
}
