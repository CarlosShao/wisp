#!/bin/sh
# 353-v1 R57 = 交件侧件身尺（四枚交付件＋全部 commit 正文文件逐枚量；工作树形⇄入库 blob 形两取）
S7=$(printf '\xc2\xa7')
cd '/d/work/workspace/projects plans/Wisp' || exit 9
echo "=== 工作树形 ==="
for f in /d/tmp/wisp/20261020/353v1/msg-0*.txt \
         docs/evidence/s1/353-v1-verdict.md \
         .scratch/wisp/probes/353/v1/00-ruler-register.md \
         .scratch/wisp/probes/353/v1/01-independent-readings.txt \
         .scratch/wisp/probes/353/v1/02-ticket-draft-secC.txt \
         .scratch/wisp/probes/353/v1/02-ticket-draft-secC.md \
         .scratch/wisp/probes/353/v1/90-truncated-nonzero-rows.txt ; do
  [ -f "$f" ] || continue
  printf '%-56s SECT=%-3s CR=%-3s LINES=%-4s PH=%-3s LAST=%s\n' "$(basename "$f")" \
    "$(grep -o "$S7" "$f" | wc -l)" "$(tr -cd '\r' < "$f" | wc -c)" "$(wc -l < "$f")" \
    "$(grep -cE 'TODO|FIXME|XXX' "$f" || true)" "$(tail -c 1 "$f" | od -An -tx1 | tr -d ' \n')"
done
echo "=== 入库 blob 形（四枚已入库交付件）==="
for p in docs/evidence/s1/353-a1-census.md ; do :; done
for p in .scratch/wisp/probes/353/v1/00-ruler-register.md .scratch/wisp/probes/353/v1/01-independent-readings.txt .scratch/wisp/probes/353/v1/02-ticket-draft-secC.md .scratch/wisp/probes/353/v1/90-truncated-nonzero-rows.txt docs/evidence/s1/353-v1-verdict.md ; do
  if git cat-file -e "HEAD:$p" 2>/dev/null; then
    git cat-file blob "HEAD:$p" > /d/tmp/wisp/20261020/353v1/_b.bin
    printf '%-56s SECT=%-3s CR=%-3s LINES=%-4s BYTES=%s\n' "$(basename "$p")" \
      "$(grep -o "$S7" /d/tmp/wisp/20261020/353v1/_b.bin | wc -l)" \
      "$(tr -cd '\r' < /d/tmp/wisp/20261020/353v1/_b.bin | wc -c)" \
      "$(wc -l < /d/tmp/wisp/20261020/353v1/_b.bin)" "$(wc -c < /d/tmp/wisp/20261020/353v1/_b.bin)"
  else
    printf '%-56s NOT_AT_HEAD\n' "$(basename "$p")"
  fi
done
echo "== 正控：自造 2 枚分节码位＋2 枚 CR 同尺读回 =="
printf "a${S7}b${S7}c\r\r" > /d/tmp/wisp/20261020/353v1/_ctrl2.bin
printf '  CTRL SECT=%s CR=%s\n' "$(grep -o "$S7" /d/tmp/wisp/20261020/353v1/_ctrl2.bin | wc -l)" "$(tr -cd '\r' < /d/tmp/wisp/20261020/353v1/_ctrl2.bin | wc -c)"
