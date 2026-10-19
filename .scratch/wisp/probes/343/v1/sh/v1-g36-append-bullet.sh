#!/bin/sh
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 90
T=.scratch/wisp/issues/343-blacklist-comment-names-aEquiv-a-symbol-that-does-not-exist-in-the-package.md
B=/d/tmp/wisp/20261019/343v1/msg/ticket-bullet.md
BEFORE=$(wc -l < "$T")
SEC="$(printf '\xc2\xa7')"
echo "== bullet file hygiene before appending: banned char / 0 bytes / utf8 =="
printf 'bullet_banned_char_hits='; grep -c "$SEC" "$B" || true
printf 'bullet_bytes='; wc -c < "$B"
iconv -f UTF-8 -t UTF-8 "$B" >/dev/null 2>&1 && echo bullet_utf8=CLEAN || echo bullet_utf8=BAD
echo "== append (>> only; ⛔ 覆盖、⛔ 改任何既有行) =="
# ensure the ticket ends with exactly one newline, then append the bullet block
tail -c 1 "$T" | od -c | head -1
cat "$B" >> "$T"
AFTER=$(wc -l < "$T")
echo "lines before=$BEFORE after=$AFTER delta=$((AFTER-BEFORE))"
echo "== append-only proof: the old content must be a byte-prefix of the new content =="
git show HEAD:"$T" > /d/tmp/wisp/20261019/343v1/logs/ticket-at-HEAD.md
cmp -n $(wc -c < /d/tmp/wisp/20261019/343v1/logs/ticket-at-HEAD.md) /d/tmp/wisp/20261019/343v1/logs/ticket-at-HEAD.md "$T" && echo PREFIX_IDENTICAL=YES || echo PREFIX_IDENTICAL=NO
echo "== diff numstat of the working tree vs HEAD for the ticket (want 0 deletions) =="
git diff --numstat -- "$T"
git diff -- "$T" | grep -c '^-[^-]' || true
echo "== box gauges after the append (want 6 / 0) =="
printf 'unchecked='; grep -c '^- \[ \] \*\*AC#' "$T"
printf 'checked='; grep -c '^- \[x\]' "$T" || true
echo "== banned char over the ticket after append =="
printf 'ticket_banned_char='; grep -c "$SEC" "$T" || true
echo "== status line / title untouched? =="
git diff -- "$T" | grep -E '^[+-]' | grep -vE '^(\+\+\+|---)' | cut -c1-60 | head -5
echo DONE
