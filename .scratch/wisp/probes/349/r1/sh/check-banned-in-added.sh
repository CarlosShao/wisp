#!/bin/sh
# 证明：我新加的那四枚注释行里 0 枚被禁字面字符（码位 U+00A7）；
# 并且 .go 起手 blob ⇄ 交件 blob 的该码位计数相等 ⇒ 我没往任何 .go 字节里带过它。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 90
A=7b633fa9713b0514806c8b808ed4db91b1485cdc
echo "== .go blob counts: START vs DELIVERY (equal => my added lines contain none) =="
for f in blacklist provenance syncdirs; do
  b=$(git cat-file blob "$A:internal/risk/$f.go" | perl -CSD -0ne 'print scalar(() = /\x{00A7}/g)')
  d=$(git cat-file blob "HEAD:internal/risk/$f.go" | perl -CSD -0ne 'print scalar(() = /\x{00A7}/g)')
  printf '%-16s start=%s delivery=%s delta=%s\n' "$f.go" "${b:-0}" "${d:-0}" "$(( ${d:-0} - ${b:-0} ))"
done
echo "== the + side of my diff (my added bytes), scanned for the codepoint =="
git diff "$A" HEAD -- internal/risk | grep '^+' | grep -v '^+++' > /tmp/plus-349.txt 2>/dev/null || \
  git diff "$A" HEAD -- internal/risk | grep '^+' | grep -v '^+++' > .scratch/wisp/probes/349/r1/tmp-plus.txt
P=.scratch/wisp/probes/349/r1/tmp-plus.txt
[ -f /tmp/plus-349.txt ] && cp -f /tmp/plus-349.txt "$P"
echo "  added_lines=$(grep -c . "$P")"
echo "  U+00A7 hits in added lines = $(perl -CSD -0ne 'print scalar(() = /\x{00A7}/g)' "$P")"
echo "== positive control: same gauge over the whole delivery blob of provenance.go (want nonzero = 存量) =="
perl -CSD -0ne 'print scalar(() = /\x{00A7}/g)' internal/risk/provenance.go
echo ""
