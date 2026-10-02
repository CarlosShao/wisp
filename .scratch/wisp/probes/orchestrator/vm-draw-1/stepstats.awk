BEGIN { FS = "\t" }
{
  line = $3;
  for (i = 4; i <= NF; i++) line = line "\t" $i;
  sub(/^\xef\xbb\xbf/, "", line);
  sub(/^[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9:.+-]+Z? ?/, "", line);
  gsub(/\r/, "", line);
  gsub(/\033\[[0-9;]*m/, "", line);
  if (line ~ /^##\[group\]Run /) {
    cur = substr(line, 14);
    gsub(/  +/, " ", cur);
    order[++ncur] = cur;
    seen[cur] = 1;
    next;
  }
  if (line ~ /^##\[group\]Post /) { cur = "POST:"; next }
  if (cur == "") next;
  if (line ~ /^=== RUN/) { R[cur]++; next }
  if (line ~ /^--- PASS:/) { P[cur]++; next }
  if (line ~ /^--- FAIL:/) { F[cur]++; next }
  if (line ~ /^--- SKIP:/) { S[cur]++; next }
}
END {
  for (i = 1; i <= ncur; i++) {
    k = order[i];
    if (R[k] + P[k] + F[k] + S[k] == 0) continue;
    printf "RUN=%-5s PASS=%-5s FAIL=%-4s SKIP=%-4s | %s\n", R[k]+0, P[k]+0, F[k]+0, S[k]+0, substr(k, 1, 90);
  }
}
