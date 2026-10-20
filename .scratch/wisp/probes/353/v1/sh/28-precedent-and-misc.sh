#!/bin/bash
# 353-v1 R28 = 先例形（票 322 AC#1 票面）⇄§C 件名册实数 ⇄ 0b75e40b 之后的笔 ⇄ 读数件⛔合并⛔相减尺
cd '/d/work/workspace/projects plans/Wisp' || exit 9
V=/d/tmp/wisp/20261020/353v1
echo "###### 1) 票 322 票面 AC#1 那一格逐字（先例形＝写进票面还是写进仪器？）"
T322=$(git ls-tree --name-only HEAD .scratch/wisp/issues/ | grep -E '/322-' | head -1)
echo "  票面=$T322"
git cat-file blob "HEAD:$T322" > $V/ticket322.md.txt
grep -nE 'AC#1' $V/ticket322.md.txt | cut -c1-260 | sed -n '1,10p'
echo "-- 票 322 面里 '窗口/window/报' 命中枚数 --"
grep -c '报窗口' $V/ticket322.md.txt || true
grep -n '报窗口' $V/ticket322.md.txt | cut -c1-200 | sed -n '1,6p'
echo "###### 2) §C 那 11 枚件名册实数（逐 pattern）"
tot=0
for pat in 'ancestoropen_342' 'ancestorunlink_348' 'ancestorworld_342_348' 'leafrace_342r2' 'openleaf_.*322' 'posix_seal_window_120' 'sealhandle_privateroute_120'; do
  n=$(git ls-tree --name-only HEAD internal/winsec/ | grep -Ec "$pat" || true)
  echo "  pattern[$pat]=$n"
  git ls-tree --name-only HEAD internal/winsec/ | grep -E "$pat" | sed 's|internal/winsec/||;s/^/     /'
  tot=$((tot+n))
done
echo "  合计命中=$tot（派单写 11 枚＝具名对质）"
echo "###### 3) 0b75e40b 之后的每一笔（逐枚列名＋主题起头，⛔ 区间归因）"
git log --format='%h %ad %s' --date=format:'%H:%M' 0b75e40b..HEAD | cut -c1-120 | sed 's/^/  /'
echo "  枚数=$(git log --format='%h' 0b75e40b..HEAD | grep -c . || true)"
echo "###### 4) 读数件里有没有'合并/相减'那形（尺＝grep 符号＋关键词）"
for f in .scratch/wisp/probes/353/r1/readings-ac2-baseline.md .scratch/wisp/probes/353/r1/readings-ac2-landed.md; do
  printf '  %s  minus_sign=%s  合并=%s  相减=%s\n' "$(basename $f)" \
    "$(grep -c ' - ' $f || true)" "$(grep -c '合并' $f || true)" "$(grep -c '相减' $f || true)"
done
echo "###### 5) docs/evidence/s1 下 353* 现有件（我的写面是否空闲）"
git ls-tree --name-only HEAD docs/evidence/s1/ | grep -E '353' | sed 's/^/  /'
ls -1 docs/evidence/s1/ | grep -E '353' | sed 's/^/  盘上: /'
echo "###### 6) 我这封派单自身入库了吗（⛔ 指不存在的件）"
git log --format='%h %s' --all -1 -- .scratch/wisp/probes/orch/20261020/txt/353v1-dispatch.txt | cut -c1-100 | sed 's/^/  /'
git cat-file -e 'HEAD:.scratch/wisp/probes/orch/20261020/txt/353v1-dispatch.txt' && echo "  DISPATCH_AT_HEAD=YES" || echo "  DISPATCH_AT_HEAD=NO"
ls -1 .scratch/wisp/probes/orch/20261020/txt/ | grep 353v1 | sed 's/^/  /'
