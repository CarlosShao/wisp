#!/usr/bin/env bash
# summarize.sh - 311-r1b: one table of per-shot counters (every count standing next to
# its denominator: the capture's byte size and its rc= line), plus the winsec name roster
# comparison pre<->post with the two rulers side by side (sub-tests included vs top-level
# only - they are NOT added and NOT interchangeable).
set -u
WT="/d/work/workspace/projects plans/Wisp"
L="$WT/.scratch/wisp/probes/311/r1b/logs"
OUT="$L/41-shot-table.txt"
TAB="$L/42-winsec-roster.txt"
N() { grep -c "$1" "$2" 2>/dev/null || true; }
NL() { grep -cE "$1" "$2" 2>/dev/null || true; }
names() { grep -E '^--- '"$1"': ' "$2" | sed -E 's/^--- [A-Z]+: ([^ ]+).*/\1/' | sort; }

: > "$OUT"
echo "# 41-shot-table.txt - rulers: grep -c on each capture (denominator = bytes + rc line on the same row)" >> "$OUT"
echo "shot_file | bytes | rc | RUN=^===RUN | PASSall | PASS_top | FAILall | FAIL_top | TempDirCleanup | AccessDenied | 310hits | 311GUARD | 311INCONC | first_GOTMPDIR_line" >> "$OUT"
for f in "$L"/1[0-9]*pre*.txt "$L"/2[0-9]-*.txt "$L"/3[0-9]-*.txt; do
    [ -f "$f" ] || continue
    b=$(wc -c < "$f")
    rc=$(grep -E '^rc=' "$f" | tr -d '\n')
    echo "$(basename "$f") | $b | $rc | RUN=$(N '^=== RUN' "$f") | PASSall=$(N '^ *--- PASS' "$f") | PASS_top=$(N '^--- PASS' "$f") | FAILall=$(N '^ *--- FAIL' "$f") | FAIL_top=$(N '^--- FAIL' "$f") | TDcleanup=$(N 'TempDir RemoveAll cleanup' "$f") | AccessDenied=$(N 'Access is denied' "$f") | 310=$(N 'ticket 310 AC#1' "$f") | 311GUARD=$(N '311 AC#1 shape B: GUARD' "$f") | 311INCONC=$(N '311 AC#1 shape B: INCONCLUSIVE' "$f") | $(grep -m1 -o 'GOTMPDIR=[^ ]*' "$f" || true)" >> "$OUT"
done
echo "--- four-numbers / census / verdict lines, verbatim, per shot that has one ---" >> "$OUT"
grep -H -E 'four numbers|census totals|runtests.sh: OK|VERDICT|gate is RED' "$L"/1[0-9]*pre*.txt "$L"/2[0-9]-*.txt "$L"/3[0-9]-*.txt 2>/dev/null | cut -c1-240 >> "$OUT"

: > "$TAB"
echo "# 42-winsec-roster.txt - name roster: two rulers side by side, never added, never interchanged" >> "$TAB"
echo "ruler A (sub-tests included) = grep -c '^ *--- PASS' ; ruler B (top-level only) = grep -c '^--- PASS'" >> "$TAB"
echo "top_names = grep -cE '^--- (PASS|FAIL|SKIP): ' (denominator for a name roster row)" >> "$TAB"
echo "self-correction: the first trip of THIS sub-ruler used \\| inside an ERE, which matched nothing and printed a fake top_names=0; fixed here, and the corrected value is printed next to its denominator." >> "$TAB"
for f in 10-winsec-pre-run1 10-winsec-pre-run2 30-G4-post-winsec-run1 30-G5-post-winsec-run2 26-D1-winsec-samevol 22-B1-winsec-xvol-guard-removed 22-B3-runtests-xvol-guard-removed; do
    p="$L/$f.txt"
    [ -f "$p" ] || { echo "$f MISSING" >> "$TAB"; continue; }
    echo "$f: A_pass=$(N '^ *--- PASS' "$p") B_pass=$(N '^--- PASS' "$p") A_fail=$(N '^ *--- FAIL' "$p") B_fail=$(N '^--- FAIL' "$p") top_names=$(NL '^--- (PASS|FAIL|SKIP): ' "$p") TDcleanup=$(N 'TempDir RemoveAll cleanup' "$p")" >> "$TAB"
    names 'FAIL' "$p" | sed 's/^/    FAILNAME /' >> "$TAB"
done
echo "--- pre(2 shots) vs post(2 shots) PASS-name sets: intersection + per-name diff ---" >> "$TAB"
names 'PASS' "$L/10-winsec-pre-run1.txt" > /tmp/pre1; names 'PASS' "$L/10-winsec-pre-run2.txt" > /tmp/pre2
names 'PASS' "$L/30-G4-post-winsec-run1.txt" > /tmp/pst1; names 'PASS' "$L/30-G5-post-winsec-run2.txt" > /tmp/pst2
comm -12 /tmp/pre1 /tmp/pre2 > /tmp/preint; comm -12 /tmp/pst1 /tmp/pst2 > /tmp/pstint
echo "pre_names_run1=$(wc -l < /tmp/pre1) pre_names_run2=$(wc -l < /tmp/pre2) pre_intersection=$(wc -l < /tmp/preint)" >> "$TAB"
echo "post_names_run1=$(wc -l < /tmp/pst1) post_names_run2=$(wc -l < /tmp/pst2) post_intersection=$(wc -l < /tmp/pstint)" >> "$TAB"
echo "names_in_pre_not_in_post=$(comm -23 /tmp/preint /tmp/pstint | wc -l)"; comm -23 /tmp/preint /tmp/pstint | sed 's/^/    LOST /' >> "$TAB"
echo "names_in_post_not_in_pre=$(comm -13 /tmp/preint /tmp/pstint | wc -l)"; comm -13 /tmp/preint /tmp/pstint | sed 's/^/    GAINED /' >> "$TAB"
echo "--- red-name roster on the cross-volume, guard-removed shots (AC#2 (c) evidence) ---" >> "$TAB"
grep -H -E '^--- FAIL|TempDir RemoveAll cleanup' "$L"/22-B1-winsec-xvol-guard-removed.txt "$L"/22-B3-runtests-xvol-guard-removed.txt 2>/dev/null | cut -c1-200 >> "$TAB"
echo "summarize-finished"
