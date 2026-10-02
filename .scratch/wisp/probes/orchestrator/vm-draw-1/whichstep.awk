BEGIN { FS = "\t" }
{
  line = $3; for (i = 4; i <= NF; i++) line = line "\t" $i;
  sub(/^\xef\xbb\xbf/, "", line);
  sub(/^[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9:.+-]+Z? ?/, "", line);
  gsub(/\r/, "", line);
  if (line ~ /^##\[group\]Run /) { cur = substr(line, 14); gsub(/  +/, " ", cur); next }
  if (line ~ /TestTicket223HandEditedFsLooseningCostsAnL2Card|TestResolvePerCallBudget/)
    print FILENAME " | " substr(cur,1,60) " | " substr(line,1,120);
}
