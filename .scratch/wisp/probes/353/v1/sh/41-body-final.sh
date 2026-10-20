#!/bin/sh
# 353-v1 交件侧件身尺（量的；一律对**入库 blob**；CR＝tr -cd '\r'｜分节码位＝grep -o 双字节）
cd '/d/work/workspace/projects plans/Wisp' || exit 9
S7=$(printf '\xc2\xa7')
FILES="docs/evidence/s1/353-v1-verdict.md .scratch/wisp/probes/353/v1/00-ruler-register.md .scratch/wisp/probes/353/v1/01-independent-readings.txt .scratch/wisp/probes/353/v1/02-ticket-draft-secC.md"
echo "== 工作树形（写前预检）=="
for p in $FILES; do
  [ -f "$p" ] || { echo "MISSING $p"; continue; }
  b=$(wc -c < "$p"); l=$(wc -l < "$p"); cr=$(tr -cd '\r' < "$p" | wc -c)
  s=$(grep -o "$S7" "$p" | wc -l); ph=$(grep -cE 'TODO|FIXME|XXX' "$p" || true)
  lb=$(tail -c 1 "$p" | od -An -tx1 | tr -d ' \n')
  ic=$(iconv -f UTF-8 -t UTF-8 "$p" > /dev/null 2>&1 && echo OK || echo BAD)
  printf '%-62s BYTES=%s LINES=%s CR=%s SECT=%s PH=%s LAST=%s ICONV=%s\n' "$(basename $p)" "$b" "$l" "$cr" "$s" "$ph" "$lb" "$ic"
done
echo "== 正控（同尺打一枚已知含分节码位的件＝票 353 面＝0，另一枚已知的＝census）=="
for p in .scratch/wisp/issues/353-pinned-name-swap-after-seal-deletes-an-object-the-caller-never-named.md docs/evidence/s1/353-a1-census.md; do
  git cat-file blob "HEAD:$p" > /tmp/pc.bin
  printf '  %-58s SECT=%s CR=%s LINES=%s\n' "$(basename $p)" "$(grep -o "$S7" /tmp/pc.bin | wc -l)" "$(tr -cd '\r' < /tmp/pc.bin | wc -c)" "$(wc -l < /tmp/pc.bin)"
done
echo "== sh 包（逐枚非 0 字节＋无 .out 形）=="
for p in .scratch/wisp/probes/353/v1/sh/*.sh .scratch/wisp/probes/353/v1/90-*.txt; do
  printf '  %s %s\n' "$(wc -c < "$p")" "$p"
done
echo "  .out 枚数=$(ls -1 .scratch/wisp/probes/353/v1/sh/ | grep -c '\.out$' || true)"
echo "  0 字节枚数=$(find .scratch/wisp/probes/353/v1 -type f -size 0 | wc -l)"
