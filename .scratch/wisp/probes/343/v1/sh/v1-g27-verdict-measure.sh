#!/bin/sh
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 90
F=docs/evidence/s1/343-v1-verdict.md
echo "LINES=$(wc -l < $F) BYTES=$(wc -c < $F)"
echo "--- per-section line counts (cap 8 per AC cell) ---"
awk '/^## /{if(name!="")print n"  "name; name=$0; n=0; next}{if(name!="")n++}END{if(name!="")print n"  "name}' $F
echo "--- banned literal char (U+00A7) hits ---"
grep -c "$(printf '\xc2\xa7')" $F || true
echo "--- CR bytes ---"
tr -cd '\r' < $F | wc -c
echo "--- positive control same gauge on HANDOVER ---"
grep -c "$(printf '\xc2\xa7')" docs/reports/HANDOVER.md || true
