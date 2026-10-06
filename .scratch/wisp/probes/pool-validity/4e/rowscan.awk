#!/usr/bin/awk -f
# pool-validity 4e: roster-vs-judged extractor for the six batch files.
# A numeric table row counts toward ROSTER always; it counts toward JUDGED
# only when one of its fields begins with a tier word (or an unjudged marker).
BEGIN { FS = "|" }
/^\| *[0-9]+ *\|/ || /^\| *\*\*[0-9]+\*\* *\|/ {
  id = $2; gsub(/[ *]/, "", id);
  if (id !~ /^[0-9]+$/) next;
  tier = "";
  for (i = 3; i <= NF; i++) {
    v = $i; gsub(/^ +| +$/, "", v); gsub(/^\*+|\*+$/, "", v);
    if (v ~ /^(仍成立|已失效|差翻勾|量不到|未判|尚未判|待填)/) { tier = v; break }
  }
  print id "\t" (tier == "" ? "LISTONLY" : (tier ~ /^(未判|尚未判|待填)/ ? "UNJUDGED" : "JUDGED")) "\t" substr(tier, 1, 12)
}
