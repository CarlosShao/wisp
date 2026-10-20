#!/bin/sh
# 353-r1 码位尺：对**入库 blob**（⛔ 工作树）量非 ASCII 码位分布 ⇄ CR 字节数。
# ★CR＝`git cat-file blob <锚>:<path> | tr -cd '\r' | wc -c`（公共块：⛔ `grep -c $'\r'`＝钝尺）。
# ★d22scan 实际扫的带＝U+1F000–1FAFF / U+2200–22FF / U+2600–27BF / U+2B00–2BFF / U+FE0F / U+1F1E6–1F1FF。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 9
REF="${1:-HEAD}"
for p in "$@"; do :; done
shift 2>/dev/null || true
echo "REF=$REF FULL=$(git rev-parse "$REF")"
for p in "$@"; do
  echo "== $p =="
  git cat-file blob "$REF:$p" > "/d/tmp/wisp/20261020/353r1/blob.txt"
  echo "  CR_BYTES=$(tr -cd '\r' < /d/tmp/wisp/20261020/353r1/blob.txt | wc -c)"
  echo "  BYTES=$(wc -c < /d/tmp/wisp/20261020/353r1/blob.txt)"
  echo "  LINES=$(wc -l < /d/tmp/wisp/20261020/353r1/blob.txt)"
  python -c '
import sys, collections, unicodedata
t = open(sys.argv[1], encoding="utf-8").read()
c = collections.Counter(ch for ch in t if ord(ch) > 127)
bands = [(0x1F000,0x1FAFF,"1F000-1FAFF"),(0x2200,0x22FF,"2200-22FF"),(0x2600,0x27BF,"2600-27BF"),
         (0x2B00,0x2BFF,"2B00-2BFF"),(0xFE00,0xFE0F,"FE00-FE0F"),(0x1F1E6,0x1F1FF,"1F1E6-1F1FF"),
         (0x2190,0x21FF,"2190-21FF (arrow: 规格禁/仪器不扫)"),(0x2460,0x24FF,"2460-24FF (带圈数字)"),
         (0x00A7,0x00A7,"00A7 section sign")]
hits = {}
for lo, hi, name in bands:
    n = sum(v for k, v in c.items() if lo <= k <= hi)
    if n:
        hits[name] = n
print("  NONASCII_TOTAL=%d DISTINCT=%d" % (sum(c.values()), len(c)))
print("  SCANNED_BAND_HITS=%s" % (hits if hits else "none"))
' /d/tmp/wisp/20261020/353r1/blob.txt
done
