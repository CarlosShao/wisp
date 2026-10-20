#!/bin/sh
cd "/d/work/workspace/projects plans/Wisp" || exit 1
cp /d/tmp/wisp/20261020/342v3/logs/writeface-hygiene.txt .scratch/wisp/probes/342/v3/logs/30-writeface-hygiene-final.txt
cp /d/tmp/wisp/20261020/342v3/logs/writeface-hygiene.txt .scratch/wisp/probes/348/v2/logs/04-writeface-hygiene.txt
echo "ROOT342=$(find .scratch/wisp/probes/342/v3 -type f | wc -l)"
echo "ROOT348=$(find .scratch/wisp/probes/348/v2 -type f | wc -l)"
echo "LOGS342=$(find .scratch/wisp/probes/342/v3/logs -type f | wc -l)"
echo "SH342=$(find .scratch/wisp/probes/342/v3/sh -type f | wc -l)"
echo "ZERO=$(find .scratch/wisp/probes/342/v3 .scratch/wisp/probes/348/v2 -type f -size 0 | wc -l)"
echo "OUT=$(find .scratch/wisp/probes/342/v3 .scratch/wisp/probes/348/v2 -type f -name '*.out' | wc -l)"
echo "-- claims about my own roots in every doc --"
grep -o '（[0-9]\{1,3\} 枚[^）]\{0,14\}' docs/evidence/s1/342-v3-verdict.md docs/evidence/s1/348-v2-verdict.md .scratch/wisp/probes/342/v3/report.md .scratch/wisp/probes/348/v2/readings-348.md /d/tmp/wisp/20261020/342v3/sh/tail-342.txt /d/tmp/wisp/20261020/342v3/sh/tail-348.txt || true

# ---- appended: final-selfproof2.sh (the ruler whose output is above) ----
#!/bin/sh
# 342-v3 FINAL self-proof (v2, no /dev/stderr: that capture is the documented sandbox boundary on this
# host and it corrupted the previous run's output).
cd "/d/work/workspace/projects plans/Wisp" || exit 1
{
echo "== HEAD =="
git log -1 --format='%H %ad %s' --date=format:'%m-%d %H:%M' | cut -c1-90
echo "== roster shape A (anchored subject) =="
git log --format='%H%x09%s' > /d/tmp/wisp/20261020/342v3/all2.tsv
perl -ne 'chomp; @f=split(/\t/,$_,2); print "$f[0]\n" if defined $f[1] && $f[1]=~/^票 342＋票 348 续体终裁腿 342-v3·/' /d/tmp/wisp/20261020/342v3/all2.tsv > /d/tmp/wisp/20261020/342v3/rosterA.txt
echo "A=$(wc -l < /d/tmp/wisp/20261020/342v3/rosterA.txt)"
echo "== roster shape B (--grep anchored) =="
echo "B=$(git log --grep='^票 342＋票 348 续体终裁腿 342-v3·' --format='%h' | wc -l)"
echo "== roster shape C (the ASCII-subject correction commit) =="
git log --grep='^342-v3 commit-' --format='%h %s' | cut -c1-60
echo "== union by hash (A+B+C) =="
{ cat /d/tmp/wisp/20261020/342v3/rosterA.txt; git log --grep='^票 342＋票 348 续体终裁腿 342-v3·' --format='%H'; git log --grep='^342-v3 commit-' --format='%H'; } | sort -u > /d/tmp/wisp/20261020/342v3/my-union.txt
echo "UNION=$(wc -l < /d/tmp/wisp/20261020/342v3/my-union.txt)"
echo "== per-commit membership over the union =="
while read -r c; do
  names=$(git show --name-only --format= "$c" | grep -c . || true)
  outside=$(git show --name-only --format= "$c" | grep -vE '^(docs/evidence/s1/(342-v3|348-v2)-verdict\.md|\.scratch/wisp/probes/34[28]/v[23]/|\.scratch/wisp/issues/(342|348)-)' | grep -c . || true)
  prodgo=$(git show --name-only --format= "$c" | grep -E '\.go$' | grep -v '_test\.go$' | grep -c . || true)
  anc=$(git merge-base --is-ancestor "$c" HEAD && echo Y || echo N)
  echo "  $(echo "$c" | cut -c1-8) names=$names outside=$outside prodgo=$prodgo ancestor=$anc"
done < /d/tmp/wisp/20261020/342v3/my-union.txt > /d/tmp/wisp/20261020/342v3/percommit.txt
cat /d/tmp/wisp/20261020/342v3/percommit.txt
echo "== non-_test .go touched by the leg over 677d4837..HEAD (expect 0) =="
git diff --name-only 677d4837..HEAD -- '*.go' | grep -v '_test\.go$' > /d/tmp/wisp/20261020/342v3/leg-prodgo.txt || true
echo "prodgo=$(grep -c . /d/tmp/wisp/20261020/342v3/leg-prodgo.txt || true)"
echo "== _test .go touched by the leg (expect empty: the leg wrote no .go at all) =="
git diff --name-only 677d4837..HEAD -- '*.go' > /d/tmp/wisp/20261020/342v3/leg-anygo.txt || true
echo "anygo=$(grep -c . /d/tmp/wisp/20261020/342v3/leg-anygo.txt || true)"
echo "== verdict files =="
for p in docs/evidence/s1/342-v3-verdict.md docs/evidence/s1/348-v2-verdict.md; do
  echo "  $(pwd)/$p bytes=$(wc -c < "$p") lines=$(wc -l < "$p") CR=$(tr -dc '\r' < "$p" | wc -c) A7=$(perl -CSD -0777 -ne 'my $c = () = /(\N{U+00A7})/g; print $c;' "$p")"
done
echo "== write-face counts at HEAD (ls-tree) =="
echo "  probes/342/v3=$(git ls-tree -r --name-only HEAD -- .scratch/wisp/probes/342/v3 | grep -c .) probes/348/v2=$(git ls-tree -r --name-only HEAD -- .scratch/wisp/probes/348/v2 | grep -c .)"
echo "  zero=$(git ls-tree -r HEAD -- .scratch/wisp/probes/342/v3 .scratch/wisp/probes/348/v2 | awk '$3=="0"' | grep -c . || true) dotout=$(git ls-tree -r --name-only HEAD -- .scratch/wisp/probes/342/v3 .scratch/wisp/probes/348/v2 | grep -c '\.out$' || true)"
echo "  A7 over roots=$(find .scratch/wisp/probes/342/v3 .scratch/wisp/probes/348/v2 -type f -exec cat {} + | perl -CSD -0777 -ne 'my $c = () = /(\N{U+00A7})/g; print $c;')"
echo "== winsec.go HEAD blob =="
git rev-parse HEAD:internal/winsec/winsec.go
echo "== box rulers + tails =="
for n in 342 348; do f=$(ls .scratch/wisp/issues/${n}-*.md); echo "  $(basename "$f") [x]=$(grep -c '^- \[x\] \*\*AC' "$f") [ ]=$(grep -c '^- \[ \] \*\*AC' "$f") mybullets=$(grep -c 'agent=`342-v3`' "$f")"; done
git diff --numstat 677d4837..HEAD -- .scratch/wisp/issues/342-*.md .scratch/wisp/issues/348-*.md
echo "== dirty state of my write faces (expect 0) =="
echo "  dirty=$(git status --porcelain -- .scratch/wisp/probes/342/v3 .scratch/wisp/probes/348/v2 docs/evidence/s1/342-v3-verdict.md docs/evidence/s1/348-v2-verdict.md .scratch/wisp/issues/342-privatefile-open-follows-a-swapped-ancestor-on-both-platforms-no-ticket-owns-it.md .scratch/wisp/issues/348-refusal-path-remove-follows-an-ancestor-link-and-deletes-a-foreign-file.md | grep -c . || true)"
echo "== unpushed count (for the record; the leg never pushes) =="
echo "  ahead=$(git rev-list --count origin/dev..HEAD 2>/dev/null || echo '?')"
} > /d/tmp/wisp/20261020/342v3/logs/final-selfproof.txt 2>&1
echo RC=$?

# ---- appended: restore-verify.sh ----
#!/bin/sh
# 342-v3 restore-equals-HEAD verification: the teeth shots restored from the leg's own pristine copies;
# prove that restored bytes == git cat-file blob HEAD:<path> for every mutated file (the public block's
# restore rule is "restore from the HEAD blob"; hash equality is the same fact measured differently).
cd "/d/work/workspace/projects plans/Wisp" || exit 1
{
for t in treeL treeT; do
  p="internal/winsec/pinnedleaf_342_windows.go"
  h1=$(git cat-file blob "HEAD:$p" | sha256sum | cut -d' ' -f1)
  h2=$(sha256sum < "/d/tmp/wisp/20261020/342v3/$t/$p" | cut -d' ' -f1)
  echo "$t: HEAD-blob sha256=$h1"
  echo "$t: restored  sha256=$h2  equal=$([ "$h1" = "$h2" ] && echo YES || echo NO)"
done
for t in treeL treeT; do
  p="internal/winsec/winsec.go"
  h1=$(git cat-file blob "HEAD:$p" | sha256sum | cut -d' ' -f1)
  h2=$(sha256sum < "/d/tmp/wisp/20261020/342v3/$t/$p" | cut -d' ' -f1)
  echo "$t winsec.go equal-to-HEAD=$([ "$h1" = "$h2" ] && echo YES || echo NO) blob=$(git rev-parse HEAD:$p | cut -c1-8)"
done
} > /d/tmp/wisp/20261020/342v3/logs/restore-verify.txt 2>&1
cat /d/tmp/wisp/20261020/342v3/logs/restore-verify.txt
