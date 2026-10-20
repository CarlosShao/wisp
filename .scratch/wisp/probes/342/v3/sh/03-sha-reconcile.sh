#!/bin/sh
# Reconciling the :118-140 sha16 label. 342-v2 + the dispatch both quote e8a2f33dea8cdeab; my first ruler
# gave 88f51af530a75de9. The substance (baseline == HEAD, byte-identical) is already proven; this names
# WHICH range/normalisation each number covers so no one reads two labels as a change.
cd "/d/work/workspace/projects plans/Wisp" || exit 1
OUT=/d/tmp/wisp/20261020/342v3/logs
{
echo "== candidate ranges on HEAD:internal/winsec/winsec.go =="
for r in 118,138 118,140 118,141 117,140; do
  a=$(git cat-file blob HEAD:internal/winsec/winsec.go | sed -n "${r}p" | sha1sum | cut -c1-16)
  b=$(git cat-file blob fde1df02:internal/winsec/winsec.go | sed -n "${r}p" | sha1sum | cut -c1-16)
  echo "sed -n '${r}p'  HEAD=$a  BASELINE=$b  same=$([ "$a" = "$b" ] && echo YES || echo NO)"
done
echo "== working-tree bytes vs blob bytes (CRLF?) =="
git cat-file blob HEAD:internal/winsec/winsec.go | sed -n '118,140p' | tr -d '\r' | sha1sum | cut -c1-16
sed -n '118,140p' internal/winsec/winsec.go | tr -d '\r' | sha1sum | cut -c1-16
echo "== byte counts of the range on HEAD vs baseline (identity by bytes, not labels) =="
echo "HEAD=$(git cat-file blob HEAD:internal/winsec/winsec.go | sed -n '118,140p' | wc -c) BASE=$(git cat-file blob fde1df02:internal/winsec/winsec.go | sed -n '118,140p' | wc -c)"
echo "== whole-file sha16 for reference =="
git rev-parse HEAD:internal/winsec/winsec.go | cut -c1-16
echo "== the whole comment block that mentions the ancestor legs: true extent =="
git cat-file blob HEAD:internal/winsec/winsec.go | grep -n 'only half closed\|not claimed closed\|ancestor leg is open there too' 
echo "== does the AC#5 ticket text say 118-138 while the verdict said 118-140? =="
grep -o ':118-1[0-9][0-9]' .scratch/wisp/issues/342-*.md docs/evidence/s1/342-v2-verdict.md docs/evidence/s1/348-v1-verdict.md | sort | uniq -c
} > "$OUT/sha-reconcile.txt" 2>&1
echo RC=$?
