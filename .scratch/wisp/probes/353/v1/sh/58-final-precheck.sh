#!/bin/sh
# 353-v1 R58 = 主件落笔前预检（件身＋sh/ 枚数＋名册三把⇄正控＋tracked 脏逐枚 mtime＋框终量）
# ★尺形更正（本腿自抓）：本把初稿在那枚"未完成标记词"的 grep 形里写出了它们本身＝写了就自占一枚命中
#   ＝现改拼装形（尺 R60 与之互证）；＝公共块"凡禁令要禁一枚字面字符、⛔ 在正文里写出那枚字符本身"那条形在尺脚本上的同味。
S7=$(printf '\xc2\xa7')
PHRE=$(printf 'T%soDO|FIX%sME|XXX' 'O' 'X')
cd '/d/work/workspace/projects plans/Wisp' || exit 9
P=docs/evidence/s1/353-v1-verdict.md
echo "== 主件工作树形 =="
printf 'SECT=%s CR=%s LINES=%s BYTES=%s PH=%s LAST=%s ICONV=' "$(grep -o "$S7" "$P" | wc -l)" "$(tr -cd '\r' < "$P" | wc -c)" "$(wc -l < "$P")" "$(wc -c < "$P")" "$(grep -cE "$PHRE" "$P" || true)" "$(tail -c 1 "$P" | od -An -tx1 | tr -d ' \n')"
iconv -f UTF-8 -t UTF-8 "$P" >/dev/null && echo OK || echo BAD
echo "== 正控（自造 2 枚码位＋2 枚 CR 同尺读回＝应 2/2）=="
printf "a${S7}b${S7}c\r\r" > /d/tmp/wisp/20261020/353v1/_ctrl3.bin
printf '  SECT=%s CR=%s\n' "$(grep -o "$S7" /d/tmp/wisp/20261020/353v1/_ctrl3.bin | wc -l)" "$(tr -cd '\r' < /d/tmp/wisp/20261020/353v1/_ctrl3.bin | wc -c)"
echo "== sh/ 枚数与交付件清单（枚数随笔漂＝只当现量、⛔ 当名册）=="
echo "  sh 枚数=$(ls -1 .scratch/wisp/probes/353/v1/sh | wc -l)  v1 根下件数=$(ls -1 .scratch/wisp/probes/353/v1 | grep -vc '^sh$')"
find .scratch/wisp/probes/353/v1 -type f -size 0 | sed 's/^/  ZERO: /'
echo "== 名册三把（本腿⇄被裁腿⇄普查腿）=="
git log --all --format='%s' > /d/tmp/wisp/20261020/353v1/subj-final.txt
printf '  v1 anchored=%s  r1 anchored=%s  a1 anchored=%s\n' "$(grep -c '^票 353 终裁腿 353-v1' /d/tmp/wisp/20261020/353v1/subj-final.txt || true)" "$(grep -c '^票 353 落地腿 353-r1' /d/tmp/wisp/20261020/353v1/subj-final.txt || true)" "$(grep -c '^票 353 普查腿 353-a1' /d/tmp/wisp/20261020/353v1/subj-final.txt || true)"
echo "== tracked 脏逐枚列名＋mtime =="
git status --porcelain -uno | while IFS= read -r line; do
  st=$(printf '%s' "$line" | cut -c1-2); f=$(printf '%s' "$line" | cut -c4-)
  printf '  [%s] %s mtime=%s\n' "$st" "$f" "$(stat -c '%y' "$f" 2>/dev/null | cut -c1-19)"
done
echo "== winsec 漂移⇄本票框终量 =="
echo "  winsec commits after 0b75e40b=$(git log --format='%h' 0b75e40b..HEAD -- internal/winsec | grep -c . || true)"
T=.scratch/wisp/issues/353-pinned-name-swap-after-seal-deletes-an-object-the-caller-never-named.md
printf '  UNCHECKED=%s CHECKED=%s\n' "$(git cat-file blob "HEAD:$T" | grep -cE '^- \[ \]' || true)" "$(git cat-file blob "HEAD:$T" | grep -cE '^- \[x\]' || true)"
