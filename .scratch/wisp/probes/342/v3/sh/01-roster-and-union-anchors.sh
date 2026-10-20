#!/bin/sh
# 342-v3 roster + cross-leg rulers. Two INDEPENDENT shapes, then the union; every count names its ruler.
# R1: anchored subject-only export of the whole ref + in-memory index (A1045 shape, landing form).
# R1b: same ruler on my dispatch's "continuation" form (expected 0 = A1045 that flavour).
# R2: --grep anchored at ^%s (second independent ruler).
# R2b: non-anchored --grep, printed to SHOW the difference, never used as a count.
# R3: whole-leg git diff --name-status <leg start anchor>..<tip> with new/modify classification.
# R4: per-commit production .go count + winsec.go blob.
# R5: the :118-140 self-statement sha per leg commit.
cd "/d/work/workspace/projects plans/Wisp" || exit 1
OUT=/d/tmp/wisp/20261020/342v3/logs
git log --format='%H%x09%s' > /d/tmp/wisp/20261020/342v3/all-subjects.tsv
{
echo "== R1 subject-only anchored index (perl over a single export) =="
perl -ne 'chomp; @f=split(/\t/,$_,2); print "$f[0]\t$f[1]\n" if defined $f[1] && $f[1]=~/^票 342＋票 348 落地腿 342-r2/;' /d/tmp/wisp/20261020/342v3/all-subjects.tsv > "$OUT/roster-r1-landshape.txt"
wc -l < "$OUT/roster-r1-landshape.txt"
echo "== R1b same ruler, my dispatch's continuation form (expect 0) =="
perl -ne 'chomp; @f=split(/\t/,$_,2); print "$f[0]\t$f[1]\n" if defined $f[1] && $f[1]=~/342-v3/;' /d/tmp/wisp/20261020/342v3/all-subjects.tsv > "$OUT/roster-r1b-v3shape.txt"
wc -l < "$OUT/roster-r1b-v3shape.txt"
echo "== R2 --grep anchored at start of subject (independent ruler) =="
git log --grep='^票 342＋票 348 落地腿 342-r2' --format='%h %ad %s' --date=format:'%m-%d %H:%M' > "$OUT/roster-r2-grep.txt"
wc -l < "$OUT/roster-r2-grep.txt"
echo "== R2b NON-anchored --grep=342-r2 (forbidden shape; shown only to prove the anchors matter) =="
git log --grep='342-r2' --format='%h' | wc -l
echo "== R2c union of the two anchors (R1 set + R2 set, by hash) =="
cut -f1 "$OUT/roster-r1-landshape.txt" | cut -c1-8 | sort -u > /d/tmp/wisp/20261020/342v3/u1.txt
git log --grep='^票 342＋票 348 落地腿 342-r2' --format='%h' | sort -u > /d/tmp/wisp/20261020/342v3/u2.txt
cat /d/tmp/wisp/20261020/342v3/u1.txt /d/tmp/wisp/20261020/342v3/u2.txt | sort -u | tee /d/tmp/wisp/20261020/342v3/union.txt | wc -l
echo "== R3 whole-leg diff --name-status 3abbb639..bfcac9ed (the leg's own start anchor .. its tip) =="
git diff --name-status 3abbb639..bfcac9ed > "$OUT/legdiff-namestatus.txt" 2>&1
wc -l < "$OUT/legdiff-namestatus.txt"
echo "-- status letter classification --"
awk '{print $1}' "$OUT/legdiff-namestatus.txt" | sort | uniq -c
echo "-- non _test.go .go touched across the whole leg (expect 0) --"
grep -E '\.go$' "$OUT/legdiff-namestatus.txt" | grep -v '_test\.go$' > "$OUT/legdiff-prodgo.txt" || true
wc -l < "$OUT/legdiff-prodgo.txt"
echo "-- _test.go .go touched (names) --"
grep -E '_test\.go$' "$OUT/legdiff-namestatus.txt" || true
echo "-- winsec.go touched by the leg (expect no line) --"
grep -c 'internal/winsec/winsec\.go' "$OUT/legdiff-namestatus.txt" || true
echo "== R4 per-commit production .go + winsec.go blob (10 commits) =="
cut -f1 "$OUT/roster-r1-landshape.txt" > /d/tmp/wisp/20261020/342v3/leg-commits.txt
: > "$OUT/percommit-prodgo.txt"
while read -r c; do
  n=$(git show --name-only --format= "$c" | grep -E '\.go$' | grep -v '_test\.go$' | wc -l)
  t=$(git show --name-only --format= "$c" | grep -E '\.go$' | wc -l)
  all=$(git show --name-only --format= "$c" | grep -c . || true)
  b=$(git rev-parse "$c:internal/winsec/winsec.go" 2>/dev/null | cut -c1-8)
  echo "$c files=$all anygo=$t prodgo=$n winsec.go=$b" >> "$OUT/percommit-prodgo.txt"
done < /d/tmp/wisp/20261020/342v3/leg-commits.txt
cat "$OUT/percommit-prodgo.txt"
echo "== R5 :118-140 sha16 per leg commit (self-statement never edited) =="
while read -r c; do
  s=$(git cat-file blob "$c:internal/winsec/winsec.go" | sed -n '118,140p' | sha1sum | cut -c1-16)
  echo "$c :118-140 sha16=$s"
done < /d/tmp/wisp/20261020/342v3/leg-commits.txt > "$OUT/leg-118sha.txt"
cat "$OUT/leg-118sha.txt"
echo "== R6 HEAD blob of winsec.go and the two instrument files =="
git rev-parse HEAD:internal/winsec/winsec.go
git cat-file blob HEAD:internal/winsec/winsec.go | sed -n '118,140p' | sha1sum | cut -c1-16
} > "$OUT/roster-all.txt" 2>&1
echo "RC=$?"
