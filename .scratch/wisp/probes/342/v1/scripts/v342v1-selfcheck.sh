#!/usr/bin/env bash
# 342-v1 self-check rulers over its own write face, with a positive control on every
# "0 hits" claim (dispatch reading discipline: before reporting 0, run a ruler over a file
# where it must hit and look at stderr).
set -u
R="/d/work/workspace/projects plans/Wisp"
cd "$R" || exit 2
V=docs/evidence/s1/342-v1-verdict.md
B=/d/tmp/wisp/20261018/342v1/logs

echo "== (a) forbidden separator char U+00A7, matched BY BYTE at runtime =="
# FIRST CUT OF THIS RULER WAS BROKEN and is booked here: the pattern was written as escape
# TEXT, so grep read the backslashes as ordinary characters and every number it printed was
# meaningless -- it "hit 1" on this very script's own printf text. The fix builds the real
# bytes at runtime with printf octal escapes, which also keeps the char itself out of this
# file (dispatch: ⛔ write the banned char, write the code point name).
SEP=$(printf '\302\247')
printf 'x\302\247y\n' > "$B/v342v1-sep-pos.txt"
echo "FORGE CHECK (the ruler must be non-zero on a file that really carries the byte): hits = $(grep -c -- "$SEP" "$B/v342v1-sep-pos.txt" || true)"
echo "POSITIVE CTRL 1: same ruler on docs/reports/HANDOVER.md = $(grep -c -- "$SEP" docs/reports/HANDOVER.md || true)"
echo "verdict hits   = $(grep -c -- "$SEP" "$V" || true)"
echo "v1 face files carrying it = $(grep -rl -- "$SEP" .scratch/wisp/probes/342/v1/ 2>/dev/null | wc -l || true)"

echo
echo "== (b) placeholder ruler: TODO / TBD / FIXME / XXX / << / 占位 =="
echo "verdict hits = $(grep -cE 'TODO|TBD|FIXME|XXX|<<|占位' "$V" || true)"
echo "POSITIVE CTRL = same ruler on docs/reports/HANDOVER.md: $(grep -cE 'TODO|TBD|FIXME|XXX|<<|占位' docs/reports/HANDOVER.md || true)"

echo
echo "== (c) zero-byte check over the whole write face =="
echo "zero-byte count = $(find .scratch/wisp/probes/342/v1 "$V" -type f -size 0 | wc -l)"
echo "files in v1 face = $(find .scratch/wisp/probes/342/v1 -type f | wc -l)"

echo
echo "== (d) .out extension check (must be 0: .gitignore carries a global *.out) =="
echo "hits = $(find .scratch/wisp/probes/342/v1 -type f -name '*.out' | wc -l)"
echo "POSITIVE CTRL = same find-shape over .txt: $(find .scratch/wisp/probes/342/v1 -type f -name '*.txt' | wc -l)"

echo
echo "== (e) no .go bytes in my commits: .go hits per hash in MYHASHES =="
for h in $MYHASHES; do echo "$h .go-hits=$(git show --name-only --format='' "$h" | grep -c '[.]go$' || true) files=$(git show --name-only --format='' "$h" | grep -c . || true)"; done
echo "POSITIVE CTRL = a commit that DID ship .go bytes must hit non-zero (31d4a476): $(git show --name-only --format='' 31d4a476 | grep -c '[.]go$' || true)"

echo
echo "== (f) ticket 342 frame ruler (seven cells, ⛔ any flip) =="
T=.scratch/wisp/issues/342-privatefile-open-follows-a-swapped-ancestor-on-both-platforms-no-ticket-owns-it.md
echo "unchecked AC cells = $(grep -c '^- \[ ] \*\*AC#' "$T" || true)"
echo "checked cells      = $(grep -c '^- \[x]' "$T" || true)"
echo "POSITIVE CTRL = same unchecked-ruler on ticket 348: $(grep -c '^- \[ ] \*\*AC#' .scratch/wisp/issues/348-refusal-path-remove-follows-an-ancestor-link-and-deletes-a-foreign-file.md || true)"

echo
echo "== (g) no recommendation-of-shape sentence in the verdict =="
echo "hits = $(grep -cE '推荐|应当选|建议选|应该选' "$V" || true)"
echo "POSITIVE CTRL = same ruler hits the ticket's own words: $(grep -cE '推荐|应当选|建议选|应该选' .scratch/wisp/issues/342-privatefile-open-follows-a-swapped-ancestor-on-both-platforms-no-ticket-owns-it.md || true)"
