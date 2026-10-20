#!/bin/sh
S7=$(printf '\xc2\xa7')
cd '/d/work/workspace/projects plans/Wisp' || exit 9
P=.scratch/wisp/probes/353/v1/02-ticket-draft-secC.md
echo "== 草案第 6 节附近（数枚那处）=="
grep -n '枚\|11' "$P" | cut -c1-160
echo "== 件身 =="
printf 'SECT=%s CR=%s LINES=%s PH=%s LAST=%s\n' "$(grep -o "$S7" "$P" | wc -l)" "$(tr -cd '\r' < "$P" | wc -c)" "$(wc -l < "$P")" "$(grep -cE 'TODO|FIXME|XXX' "$P" || true)" "$(tail -c 1 "$P" | od -An -tx1 | tr -d ' \n')"
