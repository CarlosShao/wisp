#!/bin/sh
# 353-v1 R42 = 入库 blob 件身复量（CR 用 tr -cd 那把尺、分节码位用双字节尺；含正控）
cd '/d/work/workspace/projects plans/Wisp' || exit 9
S7=$(printf '\xc2\xa7')
# ★正控先行＝自造 2 枚 CR 的字节串，用**同一把尺**读回 2（⛔ 用 wc -l 数 CR＝钝尺，本腿第一发就这么错过一次并已改正）
printf 'a\rb\rc\n' > /d/tmp/wisp/20261020/353v1/_ctrl2cr.bin
printf '  CTRL_CR(self-made 2 CR, same ruler wc -c)=%s\n' "$(tr -cd '\r' < /d/tmp/wisp/20261020/353v1/_ctrl2cr.bin | wc -c)"
for p in "$@"; do
  git cat-file blob "HEAD:$p" > /d/tmp/wisp/20261020/353v1/_b.bin 2>/dev/null || { echo "NOT_AT_HEAD $p"; continue; }
  printf '%-62s BYTES=%s LINES=%s CR=%s SECT=%s PH=%s LAST=%s\n' \
    "$(basename "$p")" \
    "$(wc -c < /d/tmp/wisp/20261020/353v1/_b.bin)" \
    "$(wc -l < /d/tmp/wisp/20261020/353v1/_b.bin)" \
    "$(tr -cd '\r' < /d/tmp/wisp/20261020/353v1/_b.bin | wc -c)" \
    "$(grep -o "$S7" /d/tmp/wisp/20261020/353v1/_b.bin | wc -l)" \
    "$(grep -cE 'TODO|FIXME|XXX' /d/tmp/wisp/20261020/353v1/_b.bin || true)" \
    "$(tail -c 1 /d/tmp/wisp/20261020/353v1/_b.bin | od -An -tx1 | tr -d ' \n')"
done
