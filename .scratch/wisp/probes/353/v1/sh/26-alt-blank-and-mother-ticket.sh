#!/bin/sh
# 353-v1 R26 = 空底片⇄母票签名节⇄腿的交替尺本体（量的；只读）
cd '/d/work/workspace/projects plans/Wisp' || exit 9
echo "===== 1) 腿的 sh/35 交替定因尺本体（HEAD blob 前 30 行＝交替形定义）====="
git cat-file blob 'HEAD:.scratch/wisp/probes/353/r1/sh/35-alternating-attribution.sh' | sed -n '1,30p'
echo "===== 2) 交替批里那枚没打印分子行的底片（逐字全文＝⛔ 藏读数的凭据）====="
cat '/d/tmp/wisp/20261020/353r1/logs/alt-1-finalL.raw.txt'
echo "----- alt-2-finalL 的行数与内容 -----"
wc -l -c < '/d/tmp/wisp/20261020/353r1/logs/alt-2-finalL.raw.txt'
grep -oE 'form=(leaf-symlink|symlink-dir) attempts=400[^|]*' '/d/tmp/wisp/20261020/353r1/logs/alt-2-finalL.raw.txt' | head -2
echo "===== 3) 母票凭据＝342-v3 判语里签'新立一票'那一节（HEAD blob :36-44）====="
git cat-file blob 'HEAD:docs/evidence/s1/342-v3-verdict.md' | sed -n '36,44p' | cut -c1-200
