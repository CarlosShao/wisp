#!/bin/sh
# 333-v1 AC#2 ruler: does the category-predicate commit leave any literal codepoint roster
# in the ADDED lines of internal/panel/display_scrub_216.go ? positive control = same ruler
# aimed at the DELETED side must hit. Run from the repo, diffing two NAMED commits
# (ae710c49 -> 859026d2), never the shared working tree.
R="D:/work/workspace/projects plans/Wisp"
cd "$R"
F=internal/panel/display_scrub_216.go
D=$(git diff ae710c49 859026d2 -- "$F")
printf '%s\n' "$D" | grep '^+' | grep -v '^+++' > /tmp/ac2-added.txt
printf '%s\n' "$D" | grep '^-' | grep -v '^---' > /tmp/ac2-del.txt
echo "== anchor full shas =="
git rev-parse ae710c49 859026d2
echo "== added line count ="; wc -l < /tmp/ac2-added.txt
echo "== deleted line count ="; wc -l < /tmp/ac2-del.txt
echo "== RULER on ADDED side (must be 0) =="
for pat in 'case 0x' '0x' 'U+' '\\x'; do
  n=$(grep -c -- "$pat" /tmp/ac2-added.txt || true)
  echo "added  hits for [$pat] = $n"
done
echo "== POSITIVE CONTROL: same ruler on DELETED side (must be > 0) =="
for pat in 'case 0x' '0x' 'U+' '\\x'; do
  n=$(grep -c -- "$pat" /tmp/ac2-del.txt || true)
  echo "deleted hits for [$pat] = $n"
done
echo "== the changed function body verbatim (hunk around isDisplayWashCf) =="
git diff -U3 ae710c49 859026d2 -- "$F" | grep -A16 'isDisplayWashCf' | grep -E '^[+-]' | head -40
echo "== per-codepoint ruler on ADDED side: any hex literal at all? =="
grep -o -E '0x[0-9A-Fa-f]{2,6}' /tmp/ac2-added.txt | sort -u | wc -l
echo "== second-source hunt inside isDisplayWashCf body (delivery tree) =="
sed -n '150,170p' /d/W333V1-SuqCO8/internal/panel/display_scrub_216.go
echo "== package-level identifiers referenced by the predicate =="
awk '/^func isDisplayWashCf/,/^}/' /d/W333V1-SuqCO8/internal/panel/display_scrub_216.go
echo "== grep: any hand roster left anywhere in the file (hex literals) =="
grep -n -o -E '0x[0-9A-Fa-f]{2,6}' /d/W333V1-SuqCO8/internal/panel/display_scrub_216.go | sort -u
