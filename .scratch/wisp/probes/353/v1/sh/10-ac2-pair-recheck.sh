#!/bin/bash
# 353-v1 R10 = AC#2 成对读数独立复量（尺＝grep/awk 打在 r1 入库的仓外底片物理字节上；量的）
# 射程＝D:\tmp\wisp\20261020\353r1\logs\final-dir-{B,L}.raw.txt（笔②侧另加 base353-s1/s2）
L=/d/tmp/wisp/20261020/353r1/logs
echo "###### 1) 四数独立复量（^=== RUN / ^--- PASS / ^--- FAIL / ^--- SKIP）"
for f in final-dir-B.raw.txt final-dir-L.raw.txt base353-s1-treeB.raw.txt base353-s2-treeL.raw.txt; do
  printf '%-30s' "$f"
  r=$(grep -c '^=== RUN' "$L/$f" || true)
  p=$(grep -c '^--- PASS' "$L/$f" || true)
  fl=$(grep -c '^--- FAIL' "$L/$f" || true)
  sk=$(grep -c '^--- SKIP' "$L/$f" || true)
  printf 'RUN=%s PASS=%s FAIL=%s SKIP=%s\n' "$r" "$p" "$fl" "$sk"
done
echo "###### 2) 逐字分子行（over 3 shots）两侧"
for f in final-dir-B.raw.txt final-dir-L.raw.txt; do
  echo "-- $f"
  grep -n 'over 3 shots' "$L/$f" || echo "(NO MATCH)"
done
echo "###### 3) 逐发 liveness 行（realSealRanFirst / windowEntered / plantErr）"
for f in final-dir-B.raw.txt final-dir-L.raw.txt; do
  echo "-- $f"
  grep -c 'realSealRanFirst=true' "$L/$f" || true
  grep -n 'LIVENESS' "$L/$f" | sed -n '1,6p'
done
echo "###### 4) 红句 / GREEN SCOPE / CONTROL 行"
for f in final-dir-B.raw.txt final-dir-L.raw.txt; do
  echo "-- $f"
  grep -n 'AC#2 RED' "$L/$f" | sed -n '1,4p'
  grep -n 'GREEN SCOPE' "$L/$f" | sed -n '1,4p'
  grep -n 'CONTROL' "$L/$f" | sed -n '1,4p'
done
echo "###### 5) 既有钉（342r2 ancestorwindow）同发读数"
for f in final-dir-B.raw.txt final-dir-L.raw.txt; do
  echo "-- $f"
  grep -n 'AC#1-one' "$L/$f" | sed -n '1,8p'
done
