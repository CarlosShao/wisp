#!/usr/bin/env bash
W=/d/tmp/wisp/20261018/341v1
SQ="'"
DQ='"'
echo "=== canonical single-quoted string lines dumped in MY runs ==="
for f in "$W/log-p6/test.log" "$W/run-deliv-A/test.log" "$W/log-p3/test.log"; do
  echo "-- $f --"
  echo "   single-quote-shape hits  = $(grep -cE "= ${SQ}[^${SQ}]*${SQ}" "$f" || true)"
  echo "   double-quote-shape hits  = $(grep -cE "= ${DQ}" "$f" || true)"
  grep -m2 -nE "= ${SQ}" "$f" | sed 's/^/     /'
done
echo
echo "=== exact lines cited in the verdict (base_url / protocol) ==="
grep -m4 -nE "base_url = ${SQ}${SQ}|protocol = ${SQ}openai-chat${SQ}|model = ${SQ}" "$W/log-p6/test.log" | sed 's/^/   /'
