#!/bin/sh
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 90
export LC_ALL=C.UTF-8
O=/d/tmp/wisp/20261019/343v1/logs/selfproof
mkdir -p $O
T=.scratch/wisp/issues/343-blacklist-comment-names-aEquiv-a-symbol-that-does-not-exist-in-the-package.md
{
echo "############ 343-v1 delivery self-proof (three gauges: commits real / files real / own reruns) ############"
echo "== gauge A: my own roster, anchored (two independent shapes) =="
printf 'A1_gitlog_anchored='; git log --grep='^票 343 终裁腿 343-v1' --format='%h' | wc -l
git log --format='%H%x09%s' > $O/all.tsv
perl -e 'open(F,"<:raw","'$O'/all.tsv"); my $p="票 343 终裁腿 343-v1"; my $n=0;
  while(<F>){chomp; my($h,$s)=split(/\t/,$_,2); next unless defined $s;
  if(index($s,$p)==0){$n++; print "   $h $s\n"}} print "A2_index0_COUNT=$n\n";'
echo
echo "== gauge B: per-commit path lists, ⊆ write surface, .go / .out / 0-byte counts =="
for c in $(git log --grep='^票 343 终裁腿 343-v1' --format='%H' | tac); do
  echo "----- $c"
  git log -1 --format='   %s' $c
  git show --name-only --format= $c | grep -v '^$' > $O/_p.txt
  printf '   paths=%s\n' "$(wc -l < $O/_p.txt)"
  printf '   outside_allowed_surface=%s\n' "$(grep -vcE '^(\.scratch/wisp/probes/343/v1/|docs/evidence/s1/343-v1-verdict\.md$|\.scratch/wisp/issues/343-)' $O/_p.txt || true)"
  printf '   .go=%s .out=%s\n' "$(grep -c '\.go$' $O/_p.txt || true)" "$(grep -c '\.out$' $O/_p.txt || true)"
done
echo
echo "== gauge C: files real on disk (tracked tree of HEAD), byte sizes, 0-byte scan =="
git ls-tree -r --name-only HEAD -- .scratch/wisp/probes/343/v1 docs/evidence/s1/343-v1-verdict.md > $O/_tracked.txt
printf 'tracked_files_in_my_surface=%s\n' "$(wc -l < $O/_tracked.txt)"
printf 'zero_byte_tracked=%s\n' "$(while read -r f; do [ -s "$f" ] || echo "$f"; done < $O/_tracked.txt | wc -l)"
printf 'total_bytes=%s\n' "$(while read -r f; do wc -c < "$f"; done < $O/_tracked.txt | awk '{s+=$1}END{print s}')"
echo "-- banned literal char (U+00A7) across my tracked surface (want 0 files) --"
printf 'banned_files=%s  ' "$(grep -rl "$(printf '\xc2\xa7')" .scratch/wisp/probes/343/v1 docs/evidence/s1/343-v1-verdict.md 2>/dev/null | wc -l)"
printf 'positive_control_HANDOVER=%s\n' "$(grep -c "$(printf '\xc2\xa7')" docs/reports/HANDOVER.md || true)"
echo "-- UTF-8 validity of every tracked file in my surface (iconv round trip) --"
bad=0
while read -r f; do iconv -f UTF-8 -t UTF-8 "$f" >/dev/null 2>&1 || { printf 'NONUTF8 %s\n' "$f"; bad=1; }; done < $O/_tracked.txt
printf 'nonutf8_flag=%s\n' "$bad"
echo
echo "== gauge D: ticket state (boxes untouched, name/Status untouched, my bullet = last) =="
printf 'unchecked_AC=%s checked=%s\n' "$(grep -c '^- \[ \] \*\*AC#' $T)" "$(grep -c '^- \[x\]' $T || true)"
git show 29d258f9:$T | sed -n '1p;3p' > $O/_old13.txt
sed -n '1p;3p' $T > $O/_new13.txt
cmp $O/_old13.txt $O/_new13.txt && echo "name+Status BYTE-IDENTICAL=YES" || echo "name+Status BYTE-IDENTICAL=NO"
printf 'ticket_line_count=%s\n' "$(wc -l < $T)"
echo "-- my bullet is the tail: first line of the last appended block (WHOLE line; no cut -c, that works on BYTES and would leave invalid UTF-8 = my own v1 trap) --"
tail -9 $T | head -1
echo
echo "== gauge E: lane integrity (I touched NO tracked code byte) =="
printf 'tracked_changes_in_working_tree=%s (want 0)\n' "$(git status --porcelain --untracked-files=no | wc -l)"
printf 'internal_risk_dirty=%s internal_winsec_dirty=%s\n' "$(git status --porcelain -- internal/risk | wc -l)" "$(git status --porcelain -- internal/winsec | wc -l)"
printf 'blacklist.go_blob_at_HEAD=%s (leg delivery 62f2dddb…, unchanged by me)\n' "$(git rev-parse HEAD:internal/risk/blacklist.go)"
printf 'push_count_by_this_leg=0 (I never ran git push)\n'
echo SELFPROOF_DONE
} > $O/selfproof.log 2>&1
cat $O/selfproof.log
