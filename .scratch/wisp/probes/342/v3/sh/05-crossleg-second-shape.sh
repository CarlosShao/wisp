#!/bin/sh
# 342-v3 second independent cross-leg ruler: per-commit new/modify classification for the leg's OWN 10
# commits (the whole-range 3abbb639..bfcac9ed count is NOT the leg: it also carries 296-r2/296-v2/351-a1
# and the orchestrator's commits = 24 commits, proven below). Also: anchored dispatch-shape roster.
cd "/d/work/workspace/projects plans/Wisp" || exit 1
OUT=/d/tmp/wisp/20261020/342v3/logs
{
echo "== how many commits actually sit inside the leg's anchor range (NOT the leg) =="
git rev-list --count 3abbb639..bfcac9ed
echo "== the leg's own 10 commits, per-commit A/M split + non-test .go =="
while read -r c; do
  a=$(git show --name-only --format= --diff-filter=A "$c" | grep -c . || true)
  m=$(git show --name-only --format= --diff-filter=M "$c" | grep -c . || true)
  d=$(git show --name-only --format= --diff-filter=D "$c" | grep -c . || true)
  prod=$(git show --name-only --format= "$c" | grep -E '\.go$' | grep -v '_test\.go$' | grep -c . || true)
  inst=$(git show --name-only --format= "$c" | grep -E '_test\.go$' | grep -c . || true)
  echo "$(echo $c | cut -c1-8) A=$a M=$m D=$d prodgo=$prod testgo=$inst"
done < /d/tmp/wisp/20261020/342v3/leg-commits.txt
echo "== union of A/M paths over the leg's 10 commits (by git log --name-only over the explicit set) =="
: > /d/tmp/wisp/20261020/342v3/leg-paths.txt
while read -r c; do git show --name-only --format= "$c" >> /d/tmp/wisp/20261020/342v3/leg-paths.txt; done < /d/tmp/wisp/20261020/342v3/leg-commits.txt
sort -u /d/tmp/wisp/20261020/342v3/leg-paths.txt | grep -c .
echo "-- of those, non-_test .go (expect 0) --"
sort -u /d/tmp/wisp/20261020/342v3/leg-paths.txt | grep -E '\.go$' | grep -v '_test\.go$' | tee "$OUT/leg-paths-prodgo.txt" | grep -c . || true
echo "-- of those, .go at all --"
sort -u /d/tmp/wisp/20261020/342v3/leg-paths.txt | grep -E '\.go$'
echo "== write-face membership: anything outside the two instrument files + two tickets + probes/342/r2 =="
sort -u /d/tmp/wisp/20261020/342v3/leg-paths.txt | grep -vE '^internal/winsec/(ancestorwindow|leafrace)_342r2_windows_test\.go$|^\.scratch/wisp/probes/342/r2/|^\.scratch/wisp/issues/(342|348)-' | tee "$OUT/leg-paths-outside.txt" | grep -c . || true
echo "== anchored roster, dispatch's continuation shape (expect 0 = A1045) =="
git log --format='%H%x09%s' | perl -ne 'chomp; @f=split(/\t/,$_,2); print "$f[0]\t$f[1]\n" if defined $f[1] && $f[1]=~/^(票 342 )?续体终裁腿 342-v3|^342-v3/;' | tee "$OUT/roster-dispatch-anchored.txt" | grep -c . || true
echo "== who owns the three unanchored 342-v3 mentions (they are the orchestrator's own commits) =="
perl -ne 'chomp; @f=split(/\t/,$_,2); print "$f[0]\n" if defined $f[1] && $f[1]=~/342-v3/;' /d/tmp/wisp/20261020/342v3/all-subjects.tsv
} > "$OUT/crossleg-second-shape.txt" 2>&1
echo RC=$?
