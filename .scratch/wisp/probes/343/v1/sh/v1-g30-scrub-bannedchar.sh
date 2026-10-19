#!/bin/sh
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 90
T=.scratch/wisp/probes/343/v1
F=$T/logs/census/ac3-strict.txt
SEC="$(printf '\xc2\xa7')"
echo "== occurrences BEFORE scrub =="
grep -c "$SEC" $F || true
# scrub the banned literal char out of MY OWN evidence file, replacing it with its
# codepoint name; source-side bytes are untouched (this file is a copy of grep output).
perl -i -pe 'BEGIN{$s=$ENV{SEC}} $n += s/\Q$s\E/<U+00A7>/g; END{print STDERR "scrubbed=$n\n"}' $F
echo "== occurrences AFTER scrub (want 0) =="
grep -c "$SEC" $F || true
echo "== whole write-surface rescan =="
hits=$(grep -rl "$SEC" $T 2>/dev/null | wc -l)
echo "BANNEDFILES=$hits  (positive control: HANDOVER=$(grep -c "$SEC" docs/reports/HANDOVER.md || true))"
echo "== 0-byte scan over write surface =="
find $T -type f -size 0 -print | sed 's/^/ZERO /'
echo "zero_count=$(find $T -type f -size 0 | wc -l)"
echo "== .out scan =="
echo "out_count=$(find $T -name '*.out' | wc -l)"
echo "== .go scan under my probes (would join the module build if tracked) =="
echo "go_count=$(find $T -name '*.go' | wc -l)"
echo "== roster size =="
find $T -type f | wc -l
echo DONE
