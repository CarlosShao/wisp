BEGIN { FS = "\t" }
{
  job = $1
  line = ""
  for (i = 3; i <= NF; i++) line = line $i
  g = line
  sub(/^[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9:.+-]+Z?/, "", g)
  gsub(/\r/, "", g)
  sub(/^ +/, "", g)
  if (g ~ /^##\[group\]Run /) {
    step = g; sub(/^##\[group\]Run */, "", step); gsub(/ +/, " ", step)
    cur[job] = step; next
  }
  if (g ~ /^##\[error\]/) { sub(/^##\[error\] */, "ERR|", g) }
  if (g !~ /--- FAIL:|--- SKIP:|--- PASS:|^ERR\||^FAIL|^ +FAIL|exit status [1-9]|panic:/) next
  st = (cur[job] == "") ? "(pre-step)" : cur[job]
  print job "\t" substr(st,1,52) "\t" g
}
