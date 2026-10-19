#!/bin/sh
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 90
export LC_ALL=C.UTF-8
A=/d/tmp/wisp/20261019/343v1/logs/ac5
T=.scratch/wisp/probes/343/v1
mkdir -p $A
# 343-v1 AC#5 gauges, RE-ISSUED (v2). v1 had `cut -c1-40` on two multibyte lines: cut works on
# BYTES, so it truncated mid-character and left one invalid UTF-8 sequence in my own evidence
# file. Caught by my iconv round-trip hygiene scan (g31) BEFORE commit. v2 quotes those lines
# whole; every other reading below is byte-identical to v1 (same commands, same trees).
{
echo "############ AC#5 discipline gauges (343-v1, own runs, v2 = no byte-truncation) ############"
echo "== roster anchored, two independent gauges (want 6 both) =="
printf 'gauge1_anchored_gitlog='; git log --grep='^票 343 落地腿 343-r1' --format='%h' | wc -l
git log --format='%H%x09%s' > $A/all-subjects.tsv
perl -e '
  open(F,"<:raw","'$A'/all-subjects.tsv") or die;
  my $p = "票 343 落地腿 343-r1";
  my $n = 0;
  while (<F>) { chomp; my ($h,$s) = split(/\t/,$_,2);
    next unless defined $s;
    if (index($s,$p)==0) { $n++; print "$h\n"; } }
  print "gauge2_index0_COUNT=$n\n";' > $A/gauge2.txt
tail -1 $A/gauge2.txt
printf 'contrast_NONANCHORED_gitlog='; git log --grep='票 343 落地腿 343-r1' --format='%h' | wc -l
echo "-- the extra hit the unanchored form swallows (subject+body match): --"
git log --grep='票 343 落地腿 343-r1' --format='%h' > $A/nonanch.txt
sed -n 's/^\([0-9a-f]\{7\}\).*/\1/p' $A/gauge2.txt > $A/anch7.txt
comm -23 <(sort $A/nonanch.txt) <(sort $A/anch7.txt) | while read x; do git log -1 --format='  SWALLOWED %h %s' "$x" | cut -c1-200; done
echo
echo "== per-commit name-only for the 6 leg commits =="
for c in 31d4a476 2b6f8a9a 360a6af6 bf5977f4 70caae94 137044b6; do
  echo "----- $c -----"
  git show --name-only --format='%h %s' $c | sed 's/^/  /'
done
echo
echo "== union of the 6 commits' paths, and the write-surface test =="
for c in 31d4a476 2b6f8a9a 360a6af6 bf5977f4 70caae94 137044b6; do git show --name-only --format= $c; done | grep -v '^$' | sort -u > $A/union-paths.txt
printf 'union_paths='; wc -l < $A/union-paths.txt
printf 'go_files='; grep -c '\.go$' $A/union-paths.txt || true
grep '\.go$' $A/union-paths.txt | sed 's/^/  GO /'
printf 'outside_allowed_surface='; grep -vcE '^(\.scratch/wisp/probes/343/r1/|internal/risk/blacklist\.go$|\.scratch/wisp/issues/343-)' $A/union-paths.txt || true
echo
echo "== ticket boxes (want 6 unchecked / 0 checked) =="
T343=.scratch/wisp/issues/343-blacklist-comment-names-aEquiv-a-symbol-that-does-not-exist-in-the-package.md
printf 'unchecked_AC='; grep -c '^- \[ \] \*\*AC#' "$T343"
printf 'checked_any='; grep -c '^- \[x\]' "$T343" || true
echo "== ticket subject line and Status line, WHOLE, current HEAD =="
sed -n '1p' "$T343"
sed -n '3p' "$T343"
echo "== the same two lines at the leg's start anchor 29d258f9 =="
git show 29d258f9:"$T343" | sed -n '1p'
git show 29d258f9:"$T343" | sed -n '3p'
echo "== unchanged? (byte compare of lines 1 and 3) =="
a=$(git show 29d258f9:"$T343" | sed -n '1,3p')
b=$(sed -n '1,3p' "$T343")
[ "$a" = "$b" ] && echo LINES_1_3_IDENTICAL=YES || echo LINES_1_3_IDENTICAL=NO
echo "== filename unchanged (no -done), uniqueness gauged on the INDEX of commits (tracked tree) =="
git ls-files -- .scratch/wisp/issues | grep 'issues/343-' | sed 's/^/  /'
echo "== tail append-only: numstat 29d258f9 -> HEAD for the ticket, plus deletions =="
git diff --numstat 29d258f9 HEAD -- "$T343"
printf 'deleted_lines_in_ticket_diff='; git diff 29d258f9 HEAD -- "$T343" | grep -c '^-[^-]' || true
echo "== which commit(s) touched the ticket since the leg started =="
git log --format='%h %s' 29d258f9..HEAD -- "$T343" | sed 's/^/  /'
echo
echo "== my own roster (343-v1, anchored) at the time of this run =="
git log --grep='^票 343 终裁腿 343-v1' --format='  %h %s'
echo AC5_V2_DONE
} > $A/ac5-gauges-v2.log 2>&1
cp -f $A/ac5-gauges-v2.log $T/logs/ac5/ac5-gauges.log
iconv -f UTF-8 -t UTF-8 $T/logs/ac5/ac5-gauges.log >/dev/null 2>&1 && echo UTF8=CLEAN || echo UTF8=STILL_BAD
printf 'banned_char_files_in_surface='; grep -rl "$(printf '\xc2\xa7')" $T 2>/dev/null | wc -l
printf 'zero_byte_files='; find $T -type f -size 0 | wc -l
# 临时件只建不删: the earlier partial file stays on disk, marked superseded instead of removed.
if [ -f $T/logs/ac5/ac5-ticket-lines-no-truncation.log ]; then
  { echo "NOTE: this file was created by an aborted fixup attempt (g32) before the v2 gauges above";
    echo "re-ran everything. It is SUPERSEDED by ac5-gauges.log (v2) and is left in place untracked."; } >> $T/logs/ac5/ac5-ticket-lines-no-truncation.log
fi
echo V2_DONE
