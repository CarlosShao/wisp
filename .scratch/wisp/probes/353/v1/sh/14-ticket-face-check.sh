#!/bin/bash
# 353-v1 R14 = 票面 AC#2 字面判据里两味残余在不在（量的；尺＝对入库 blob grep -cF/-n）
cd '/d/work/workspace/projects plans/Wisp' || exit 9
T='.scratch/wisp/issues/353-pinned-name-swap-after-seal-deletes-an-object-the-caller-never-named.md'
git cat-file blob "HEAD:$T" > /d/tmp/wisp/20261020/353v1/blob/TICKET-353.md.txt
echo "###### 0) 票面 blob 行数/字节/CR（tr 尺，对入库 blob）"
wc -l -c < /d/tmp/wisp/20261020/353v1/blob/TICKET-353.md.txt
git cat-file blob "HEAD:$T" | tr -cd '\r' | wc -c
echo "###### 1) 框终量（未勾/已勾）"
echo "UNCHECKED=$(grep -cE '^- \[ \]' /d/tmp/wisp/20261020/353v1/blob/TICKET-353.md.txt || true)"
echo "CHECKED=$(grep -cE '^- \[x\]' /d/tmp/wisp/20261020/353v1/blob/TICKET-353.md.txt || true)"
echo "-- 逐名列出未勾/已勾 --"
grep -nE '^- \[( |x)\]' /d/tmp/wisp/20261020/353v1/blob/TICKET-353.md.txt | sed -e 's/^\(.\{120\}\).*/\1/'
echo "###### 2) AC#2 那一格逐字全文"
grep -n 'AC#2' /d/tmp/wisp/20261020/353v1/blob/TICKET-353.md.txt | sed -n '1,12p'
echo "-- 第 15 行逐字 --"
sed -n '15p' /d/tmp/wisp/20261020/353v1/blob/TICKET-353.md.txt
echo "###### 3) 两味关键词在票面命中几枚（⛔ 任何一处=0 就是没写）"
for w in '残余' '绑定之前' 'NtOpenFile' '测不到' '仍留着调用者那枚名字' '可观测行为变化' '322'; do
  printf '   %-24s ticketTOTAL=%s  AC#2line=%s\n' "$w" \
    "$(grep -cF "$w" /d/tmp/wisp/20261020/353v1/blob/TICKET-353.md.txt || true)" \
    "$(sed -n '15p' /d/tmp/wisp/20261020/353v1/blob/TICKET-353.md.txt | grep -cF "$w" || true)"
done
echo "###### 4) 正控：同一把尺在 AC#1 那行(第14行)命中（证明尺灵）"
for w in '残余' '绑定之前' 'NtOpenFile' '测不到' '仍留着调用者那枚名字' '可观测行为变化'; do
  printf '   %-24s AC#1line=%s\n' "$w" "$(sed -n '14p' /d/tmp/wisp/20261020/353v1/blob/TICKET-353.md.txt | grep -cF "$w" || true)"
done
echo "###### 5) Status 行逐字（第25行）"
sed -n '25p' /d/tmp/wisp/20261020/353v1/blob/TICKET-353.md.txt | sed -e 's/^\(.\{200\}\).*/\1/'
