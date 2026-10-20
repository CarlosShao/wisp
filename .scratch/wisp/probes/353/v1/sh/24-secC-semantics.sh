#!/bin/bash
# 353-v1 R24 = §C 那味的仪器语义现读（flipErrs/provedReparse 到底数的是什么）＋票 118/120/342/348 结案态＋台账号验在
cd '/d/work/workspace/projects plans/Wisp' || exit 9
V=/d/tmp/wisp/20261020/353v1
git cat-file blob "HEAD:internal/winsec/ancestorunlink_348_windows_test.go" > $V/AU348-HEAD.txt
echo "###### 1) flipErrs / provedReparse / cycles 的声明⇄自增处（逐字）"
grep -nE 'flipErrs|provedReparse|cycles|restock|plant' $V/AU348-HEAD.txt | cut -c1-170 | sed -n '1,60p'
echo "###### 2) 该件 :40-90（分子行 ⇄ 红句）"
sed -n '40,90p' $V/AU348-HEAD.txt | nl -ba -v40 | cut -c1-190
echo "###### 3) 票 118 / 120 / 342 / 348 结案态（尺＝ls-tree 数已入库）"
git ls-tree --name-only HEAD .scratch/wisp/issues/ | grep -E '/(118|120|342|348|353)-' | sed 's|.scratch/wisp/issues/||;s/^/   /'
echo "###### 4) 台账号验在（⛔ 引空号）"
for n in A1062 A1066 A1067 A1069 A1070 A1071 A1072; do
  printf '   %s=%s\n' "$n" "$(grep -c "^## $n" docs/reports/pending-and-issues.md || true)"
done
echo "   正控（不存在的号 A9999）=$(grep -c '^## A9999' docs/reports/pending-and-issues.md || true)"
echo "###### 5) 票 118⇄120 里那句『link 只 unlink、⛔ 跟随』的已批文字（现读、供我判⛔ 属本票射程）"
for f in $(git ls-tree --name-only HEAD .scratch/wisp/issues/ | grep -E '/(118|120)-'); do
  echo "  -- $f"
  git cat-file blob "HEAD:$f" | grep -nE 'unlink|跟随|FOLLOW|reparse|link' | cut -c1-160 | sed -n '1,14p'
done
