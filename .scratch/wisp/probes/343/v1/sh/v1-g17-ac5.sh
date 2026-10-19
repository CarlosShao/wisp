#!/bin/sh
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 90
export LC_ALL=C.UTF-8
A=/d/tmp/wisp/20261019/343v1/logs/ac5
mkdir -p $A
{
echo "############ AC#5 discipline gauges (343-v1, own runs) ############"
echo "== roster anchored two ways (want 6 both) =="
echo "-- gauge1: git log --grep anchored to subject line --"
git log --grep='^票 343 落地腿 343-r1' --format='%h %s' | wc -l
echo "-- gauge2: single %H%x09%s export + in-memory index()==1 --"
git log --format='%H%x09%s' > $A/all-subjects.tsv
perl -e 'open(F,"<","'$A'/all-subjects.tsv");$p="票 343 落地腿 343-r1";$n=0;while(<F>){chomp;($h,$s)=split(/\t/,$_,2);if(defined $s && index($s,$p)==0){$n++;print "$h $s\n"}} print "ANCHORED0_COUNT=$n\n";' | tail -3
echo "-- non-anchored grep for contrast (same-family trap) --"
git log --grep='票 343 落地腿 343-r1' --format='%h' | wc -l
echo
echo "== per-commit name-only (all 6) =="
for c in 31d4a476 2b6f8a9a 360a6af6 bf5977f4 70caae94 137044b6; do
  echo "----- $c -----"
  git show --name-only --format='%h %s' $c | sed 's/^/  /'
done
echo
echo "== .go count across the 6 commits (want exactly 1) =="
for c in 31d4a476 2b6f8a9a 360a6af6 bf5977f4 70caae94 137044b6; do git show --name-only --format= $c; done | grep -c '\.go$' || true
echo "-- which: --"
for c in 31d4a476 2b6f8a9a 360a6af6 bf5977f4 70caae94 137044b6; do git show --name-only --format= $c; done | grep '\.go$' || true
echo
echo "== any path outside the allowed write surface? (want none) =="
for c in 31d4a476 2b6f8a9a 360a6af6 bf5977f4 70caae94 137044b6; do git show --name-only --format= $c; done | grep -vE '^(\.scratch/wisp/probes/343/r1/|internal/risk/blacklist\.go|\.scratch/wisp/issues/343-)' | sed 's/^/OUTSIDE /' || true
echo
echo "== ticket boxes: unchecked / checked =="
T='.scratch/wisp/issues/343-blacklist-comment-names-aEquiv-a-symbol-that-does-not-exist-in-the-package.md'
printf 'uncheckbox_AC='; grep -c '^- \[ \] \*\*AC#' "$T"
printf 'checked_any='; { grep -c '^- \[x\]' "$T" || true; }
printf 'subject_line='; sed -n '1p' "$T" | cut -c1-40
printf 'status_line='; sed -n '3p' "$T" | cut -c1-40
echo "== ticket filename (no -done) =="
git ls-files -- .scratch/wisp/issues | grep 'issues/343-' | sed 's/^/  /'
echo "== tail append-only: diff of the ticket across the 6 commits =="
git diff --numstat 31d4a476^ 137044b6 -- "$T"
git log --format='%h %s' -- "$T" | head -8
echo "== ticket file history line counts (append => only additions after its creating commit) =="
for c in bf5977f4 137044b6; do printf '%s numstat=' "$c"; git show --numstat --format= $c -- "$T" | tr -s ' '; done
echo
echo "== 343-v1 roster (my own, anchored; want 0 at start of this run => first on disk) =="
git log --grep='^票 343 终裁腿 343-v1' --format='%h %s' | sed 's/^/  /'
echo AC5_DONE
} > $A/ac5-gauges.log 2>&1
echo DONE
cat $A/ac5-gauges.log
