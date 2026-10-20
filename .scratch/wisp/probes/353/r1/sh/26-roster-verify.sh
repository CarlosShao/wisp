#!/bin/sh
# 353-r1 名册照绿尺＝派单逐字要求：⛔ 改两半门的"错误文字"以外可让 leafrace_342r2:84-85 ⇄
#   ancestorworld_342_348:86 的归因名册照绿（它们按 why 子串作差）＝落地后逐名验。
# 尺＝对**入库 blob** 逐串验：基线锚 a59f5630 ⇄ HEAD 那两枚 `why` 子串命中数必须同值且≥1。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 9
L=/d/tmp/wisp/20261020/353r1/logs
{
  echo "== 两半门 why 子串：基线 blob ⇄ HEAD blob 命中数 =="
  for s in 'the name carries a link, so this handle' 'this handle is not the object the name carries now'; do
    a=$(git cat-file blob a59f5630:internal/winsec/pinnedleaf_342_windows.go | grep -cF "$s" || true)
    b=$(git cat-file blob HEAD:internal/winsec/pinnedleaf_342_windows.go | grep -cF "$s" || true)
    echo "  [$s] BASE=$a HEAD=$b SAME=$([ "$a" = "$b" ] && echo YES || echo NO)"
  done
  echo "== 名册用例逐名色（finalL 整包那两发里取，⛔ 子串＝尺＝^--- (PASS|FAIL|SKIP): 行首名）=="
  for n in TestAC1WindowsLeafRecheckAnswersOnTheBareSiblingForm342r2 \
           TestAC2WindowsAncestorWorldAttributesEveryRefusalToALeg342r1; do
    for f in "$L"/final-full-L-*.raw.txt; do
      line=$(grep -E "^--- (PASS|FAIL|SKIP): $n" "$f" || true)
      echo "  $(basename "$f") $n -> ${line:-ABSENT}"
    done
  done
  echo "== ancestorworld 的归因名册逐 leg 计数（按 why 作差那把尺在场＝证明两半门仍被认得出）=="
  grep -nE 'leaf identity re-check|openPinnedLeaf re-check' "$L"/final-full-L-1.raw.txt | head -6
} 2>&1
