#!/bin/bash
# 353-v1 R12 = AC#3 整包成对色⇄逐名新增红 独立复量（量的；尺＝grep/comm 打在 r1 入库底片物理字节）
L=/d/tmp/wisp/20261020/353r1/logs
echo "###### 1) 四数复量（整包两侧各 2 发）"
for f in final-full-B-1 final-full-B-2 final-full-L-1 final-full-L-2; do
  printf '%-18s' "$f"
  printf 'RUN=%s PASS=%s FAIL=%s SKIP=%s  RC未知(见fourcounts)\n' \
    "$(grep -c '^=== RUN' $L/$f.raw.txt||true)" "$(grep -c '^--- PASS' $L/$f.raw.txt||true)" \
    "$(grep -c '^--- FAIL' $L/$f.raw.txt||true)" "$(grep -c '^--- SKIP' $L/$f.raw.txt||true)"
done
echo "###### 2) 腿自己的 RC 台账（final-fourcounts.tsv 逐字）"
cat "$L/final-fourcounts.tsv"
echo "###### 3) 我自己重做名集（⛔ 采信腿的 rednames 文件）"
for f in final-full-B-1 final-full-B-2 final-full-L-1 final-full-L-2; do
  grep -E '^--- FAIL' "$L/$f.raw.txt" | sed 's/^--- FAIL: //; s/ *([0-9.]*s)$//' | sort > "/d/tmp/wisp/20261020/353v1/$f.myred.txt"
  printf '%-18s myred=%s  legred=%s\n' "$f" "$(grep -c . /d/tmp/wisp/20261020/353v1/$f.myred.txt||true)" "$(grep -c . $L/$f.rednames.txt||true)"
done
sort -u /d/tmp/wisp/20261020/353v1/final-full-B-1.myred.txt /d/tmp/wisp/20261020/353v1/final-full-B-2.myred.txt > /d/tmp/wisp/20261020/353v1/Bunion.my.txt
sort -u /d/tmp/wisp/20261020/353v1/final-full-L-1.myred.txt /d/tmp/wisp/20261020/353v1/final-full-L-2.myred.txt > /d/tmp/wisp/20261020/353v1/Lunion.my.txt
echo "-- 新增红（只在落地侧）尺=comm -13 名集 --"
comm -13 /d/tmp/wisp/20261020/353v1/Bunion.my.txt /d/tmp/wisp/20261020/353v1/Lunion.my.txt | sed 's/^/   NEWRED: /'
echo "   NEWRED_COUNT=$(comm -13 /d/tmp/wisp/20261020/353v1/Bunion.my.txt /d/tmp/wisp/20261020/353v1/Lunion.my.txt | grep -c . || true)"
echo "-- 反向（只在基线侧＝落地后变绿的枚）--"
comm -23 /d/tmp/wisp/20261020/353v1/Bunion.my.txt /d/tmp/wisp/20261020/353v1/Lunion.my.txt | sed 's/^/   onlyB: /'
echo "###### 4) 两侧名集并集全表（逐名色）"
cat /d/tmp/wisp/20261020/353v1/Bunion.my.txt /d/tmp/wisp/20261020/353v1/Lunion.my.txt | sort -u | while read -r n; do
  [ -n "$n" ] || continue
  b1=$(grep -qxF "$n" /d/tmp/wisp/20261020/353v1/final-full-B-1.myred.txt && echo RED || echo green)
  b2=$(grep -qxF "$n" /d/tmp/wisp/20261020/353v1/final-full-B-2.myred.txt && echo RED || echo green)
  l1=$(grep -qxF "$n" /d/tmp/wisp/20261020/353v1/final-full-L-1.myred.txt && echo RED || echo green)
  l2=$(grep -qxF "$n" /d/tmp/wisp/20261020/353v1/final-full-L-2.myred.txt && echo RED || echo green)
  echo "   $n |B:$b1 $b2 |L:$l1 $l2"
done
echo "###### 5) container 侧（第二平台）"
ls -1 "$L" | grep -i col | sed 's/^/   /'
for f in final-col-s1-finalB.raw.txt final-col-s2-finalL.raw.txt; do
  printf '%-30s' "$f"
  printf 'RUN=%s PASS=%s FAIL=%s SKIP=%s\n' \
    "$(grep -c '^=== RUN' $L/$f||true)" "$(grep -c '^--- PASS' $L/$f||true)" \
    "$(grep -c '^--- FAIL' $L/$f||true)" "$(grep -c '^--- SKIP' $L/$f||true)"
done
echo "###### 6) 腿报的 container RC（grep fourcounts 里的 col 行）"
grep -i 'col' "$L/final-fourcounts.tsv" || echo "(fourcounts 无 col 行)"
