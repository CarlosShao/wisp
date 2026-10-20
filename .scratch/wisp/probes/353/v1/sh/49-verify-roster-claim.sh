#!/bin/sh
# 353-v1 R49 = 校验草案第 7 节那句"真正在打 flipErrs/provedReparse 的只有哪几枚件"（量的；对入库 blob）
cd '/d/work/workspace/projects plans/Wisp' || exit 9
echo "== 每枚 winsec 件里 flipErrs / provedReparse 命中数（>0 者才是打这两枚号的手）=="
for p in $(git ls-tree --name-only HEAD internal/winsec/ | grep '_test.go$'); do
  f=$(git cat-file blob "HEAD:$p" | grep -c 'flipErrs' || true)
  r=$(git cat-file blob "HEAD:$p" | grep -c 'provedReparse' || true)
  if [ "$f" != "0" ] || [ "$r" != "0" ]; then printf '  %-64s flipErrs=%s provedReparse=%s\n' "$(basename $p)" "$f" "$r"; fi
done
echo "== 声明处（:= / struct 字段）逐枚 =="
for p in $(git ls-tree --name-only HEAD internal/winsec/ | grep '_test.go$'); do
  git cat-file blob "HEAD:$p" | grep -nE 'flipErrs[[:space:]]*(=|int)|provedReparse[[:space:]]*(=|int)' | sed "s|^|  $(basename $p):|" | cut -c1-150
done
echo "== ancestorworld 那枚 :386 收尾句在不在 =="
git cat-file blob 'HEAD:internal/winsec/ancestorworld_342_348_windows_test.go' | sed -n '386p' | cut -c1-160
echo "== leafrace 自己那行 AC#1-leaf DAMAGE 在不在 =="
git cat-file blob 'HEAD:internal/winsec/leafrace_342r2_windows_test.go' | grep -n 'AC#1-leaf DAMAGE' | cut -c1-140
