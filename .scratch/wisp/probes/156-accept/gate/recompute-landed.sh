#!/usr/bin/env bash
# accept-r1 AC#7: recompute the FOUR landed gate logs with one ruler, then roster comms.
# Anchored '^--- FAIL' only (cmd/wisp prints the string FAIL from inside PASSING cases).
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 1
P=.scratch/wisp/probes/156
OUT=.scratch/wisp/probes/156-accept/gate
mkdir -p "$OUT"

for d in gate-pre gate-post gate-r3-pre gate-r4-post; do
  for n in cmdwisp internaltools; do
    f="$P/$d/$n.log"
    if [ ! -f "$f" ]; then echo "$d/$n MISSING"; continue; fi
    run=$(grep -a -c '^=== RUN' "$f")
    pass=$(grep -a -c '^--- PASS' "$f")
    fail=$(grep -a -c '^--- FAIL' "$f")
    skip=$(grep -a -c '^--- SKIP' "$f")
    loadfail=$(grep -a -c '0xc0000135' "$f")
    plainfailstr=$(grep -a -c 'FAIL' "$f")
    verdict=$(grep -aE '^(ok|FAIL)[ \t]' "$f" | tail -1 | tr -s ' \t' ' ')
    printf '%-12s %-14s RUN=%-4s PASS=%-4s FAIL=%-4s SKIP=%-4s 0xc0000135=%-2s printedFAILstr=%-3s %s\n' \
        "$d" "$n" "$run" "$pass" "$fail" "$skip" "$loadfail" "$plainfailstr" "$verdict"
    grep -a '^=== RUN' "$f" | sed -E 's/^=== RUN +(\S+).*/\1/' | sort > "$OUT/$d-$n-runs.txt"
    grep -aE '^--- (PASS|FAIL|SKIP)' "$f" | sed -E 's/^--- ([A-Z]+): (\S+).*/\1 \2/' | sort > "$OUT/$d-$n-verdicts.txt"
  done
done

echo
echo "== roster two-way comm: r1 pre vs r1 post"
for n in cmdwisp internaltools; do
  for k in runs verdicts; do
    a="$OUT/gate-pre-$n-$k.txt"; b="$OUT/gate-post-$n-$k.txt"
    echo "$n $k pre-only=$(comm -23 "$a" "$b" | wc -l) post-only=$(comm -13 "$a" "$b" | wc -l) sizes=$(wc -l < "$a")/$(wc -l < "$b")"
  done
done
echo "-- names only-in-post (cmdwisp runs):"
comm -13 "$OUT/gate-pre-cmdwisp-runs.txt" "$OUT/gate-post-cmdwisp-runs.txt" | sed 's/^/   /'

echo
echo "== roster two-way comm: r3 pre vs r4 post"
for n in cmdwisp internaltools; do
  for k in runs verdicts; do
    a="$OUT/gate-r3-pre-$n-$k.txt"; b="$OUT/gate-r4-post-$n-$k.txt"
    echo "$n $k pre-only=$(comm -23 "$a" "$b" | wc -l) post-only=$(comm -13 "$a" "$b" | wc -l) sizes=$(wc -l < "$a")/$(wc -l < "$b")"
  done
done

echo
echo "== cross-check my extractions against the LANDED roster files (r3/r4 shipped -runs.txt)"
for n in cmdwisp internaltools; do
  for k in runs verdicts; do
    d=gate-r3-pre
    if diff -q "$OUT/$d-$n-$k.txt" "$P/$d/$n-$k.txt" >/dev/null 2>&1; then echo "$d $n $k IDENTICAL-to-landed"; else echo "$d $n $k DIFF: $(diff "$OUT/$d-$n-$k.txt" "$P/$d/$n-$k.txt" | grep -ac '^[<>]') line(s)"; fi
    d=gate-r4-post
    if diff -q "$OUT/$d-$n-$k.txt" "$P/$d/$n-$k.txt" >/dev/null 2>&1; then echo "$d $n $k IDENTICAL-to-landed"; else echo "$d $n $k DIFF: $(diff "$OUT/$d-$n-$k.txt" "$P/$d/$n-$k.txt" | grep -ac '^[<>]') line(s)"; fi
  done
done

echo
echo "== the out-of-scope always-red case must NOT appear in any of these four logs"
for d in gate-pre gate-post gate-r3-pre gate-r4-post; do for n in cmdwisp internaltools; do
  printf '  %s/%s internal/panel refs=%s TestC21DesignTokens refs=%s\n' "$d" "$n" \
    "$(grep -a -c 'internal/panel' "$P/$d/$n.log")" "$(grep -a -c 'TestC21DesignTokensFourWayAgree' "$P/$d/$n.log")"
done; done
