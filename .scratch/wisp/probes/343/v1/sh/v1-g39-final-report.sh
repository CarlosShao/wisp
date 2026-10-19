#!/bin/sh
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 90
export LC_ALL=C.UTF-8
O=/d/tmp/wisp/20261019/343v1/logs/final
mkdir -p $O
{
echo "== my roster, anchored, two independent gauges =="
printf 'A1_gitlog_anchored='; git log --grep='^票 343 终裁腿 343-v1' --format='%h' | wc -l
git log --format='%H%x09%s' > $O/all.tsv
perl -e 'open(F,"<:raw","'$O'/all.tsv"); my $p="票 343 终裁腿 343-v1"; my $n=0;
 while(<F>){chomp; my($h,$s)=split(/\t/,$_,2); next unless defined $s; if(index($s,$p)==0){$n++}}
 print "A2_index0_COUNT=$n\n";'
echo "== per commit: hash, subject, path count, .go count, and the key file byte sizes =="
for c in $(git log --grep='^票 343 终裁腿 343-v1' --format='%H' | tac); do
  printf '%s | %s | paths=%s | .go=%s\n' "$(git log -1 --format='%h' $c)" "$(git log -1 --format='%s' $c | cut -c1-60)" \
    "$(git show --name-only --format= $c | grep -vc '^$')" "$(git show --name-only --format= $c | grep -c '\.go$' || true)"
done
echo "== deliverable absolute paths + bytes + lines =="
printf 'verdict: %s bytes=%s lines=%s\n' "$(pwd)/docs/evidence/s1/343-v1-verdict.md" "$(wc -c < docs/evidence/s1/343-v1-verdict.md)" "$(wc -l < docs/evidence/s1/343-v1-verdict.md)"
printf 'probes root: %s files=%s bytes=%s\n' "$(pwd)/.scratch/wisp/probes/343/v1" \
  "$(git ls-files -- .scratch/wisp/probes/343/v1 | wc -l)" \
  "$(git ls-files -- .scratch/wisp/probes/343/v1 | while read -r f; do wc -c < "$f"; done | awk '{s+=$1}END{print s}')"
echo "== the two AC#2 hashes, recomputed one last time from the object store =="
git cat-file blob c1799cd24afce19e6accef23db23a7d7ee746dec | sed 's://.*::' | tr -d ' \t\r' | sed '/^$/d' | sha256sum | sed 's/^/start-side /'
git cat-file blob 62f2dddbde7c72eb12cbeb4cf1ec1edc2708254b | sed 's://.*::' | tr -d ' \t\r' | sed '/^$/d' | sha256sum | sed 's/^/head-side  /'
echo "== ticket final gauges =="
T=.scratch/wisp/issues/343-blacklist-comment-names-aEquiv-a-symbol-that-does-not-exist-in-the-package.md
printf 'unchecked=%s checked=%s lines=%s banned=%s\n' "$(grep -c '^- \[ \] \*\*AC#' $T)" "$(grep -c '^- \[x\]' $T || true)" "$(wc -l < $T)" "$(grep -c "$(printf '\xc2\xa7')" $T || true)"
printf 'working_tree_tracked_changes=%s\n' "$(git status --porcelain --untracked-files=no | wc -l)"
printf 'HEAD=%s\n' "$(git log -1 --format='%H %s' | cut -c1-90)"
} > $O/final-report.log 2>&1
cat $O/final-report.log
