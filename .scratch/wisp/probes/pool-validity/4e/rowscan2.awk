#!/usr/bin/awk -f
# pool-validity 4e: per-batch census of numeric table rows.
# ROSTER = every ticket id that appears as a table row (any section of the file).
# JUDGED = rows whose tier field begins with a tier word; 未判/尚未判/待填 stay roster-only.
function sorted(k) {
  m = 0; delete arr;
  for (x in k) arr[++m] = x + 0;
  for (i = 1; i <= m; i++) for (j = i + 1; j <= m; j++) if (arr[j] < arr[i]) { t = arr[i]; arr[i] = arr[j]; arr[j] = t }
  s = ""; for (i = 1; i <= m; i++) s = s " " arr[i];
  return s;
}
BEGIN { FS = "|" }
BEGINFILE { delete roster; delete judged; file = FILENAME }
/^\| *[0-9*]+ *\|/ {
  id = $2; gsub(/[ *]/, "", id);
  if (id !~ /^[0-9]+$/) next;
  roster[id] = 1;
  for (i = 3; i <= NF; i++) {
    v = $i; gsub(/^ +| +$/, "", v); gsub(/^\*+|\*+$/, "", v);
    if (v ~ /^(仍成立|已失效|差翻勾|量不到)/) { judged[id] = 1; break }
    if (v ~ /^(未判|尚未判|待填)/) { break }
  }
}
ENDFILE { printf "#F %s\n#ROSTER%s\n#JUDGED%s\n", file, sorted(roster), sorted(judged) }
