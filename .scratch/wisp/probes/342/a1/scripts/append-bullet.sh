#!/usr/bin/env bash
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 2
T=".scratch/wisp/issues/342-privatefile-open-follows-a-swapped-ancestor-on-both-platforms-no-ticket-owns-it.md"
before=$(wc -l < "$T")
before_crc=$(tr -cd "\r" < "$T" | wc -c)
cat /d/tmp/wisp/20261018/342a1/bullet-342.txt >> "$T"
after=$(wc -l < "$T")
after_crc=$(tr -cd "\r" < "$T" | wc -c)
echo "lines before=$before after=$after ; CR-bytes before=$before_crc after=$after_crc (want both 0)"
git diff --numstat -- "$T"
echo -n "deleted-lines in diff = "; git diff -- "$T" | grep -c "^-[^-]" || true
echo "== ticket Status/name/boxes untouched check =="
git diff -- "$T" | grep -E "^[+-]" | grep -cE "Status:|^- \[|## " || true
tail -1 "$T" | head -c 200; echo
