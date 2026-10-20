#!/bin/sh
# 353-v1 R46 = 纯追加硬尺（numstat 删除列 + 前缀逐字节 cmp）
cd '/d/work/workspace/projects plans/Wisp' || exit 9
P=.scratch/wisp/probes/353/v1/00-ruler-register.md
V=/d/tmp/wisp/20261020/353v1
echo "== numstat（期望 删除列=0）=="
git diff --numstat HEAD -- "$P"
git cat-file blob "HEAD:$P" > "$V/reg-old.md"
N=$(wc -c < "$V/reg-old.md")
echo "OLD_BYTES=$N"
head -c "$N" "$P" > "$V/reg-prefix.bin"
if cmp "$V/reg-old.md" "$V/reg-prefix.bin" >/dev/null 2>&1; then echo PREFIX_IDENTICAL=YES; else echo PREFIX_IDENTICAL=NO; cmp "$V/reg-old.md" "$V/reg-prefix.bin" | head -3; fi
echo "== 新工作树字节/行 =="
wc -c -l < "$P"
