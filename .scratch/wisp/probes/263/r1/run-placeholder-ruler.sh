#!/bin/sh
# 263-r1 占位标记尺（交件判据＝0）
#
# 为什么 pattern 是拼出来的、不是写死的：
# 这三个词只要逐字落在任何被扫的交付物里（包括本脚本自己），尺就会自我命中——
# 这与证据件 §⑤ 里那枚"词面钉的模式串不能匹配到自己"是同一个形状。
# 所以这里只放每个词的 UTF-8 字节八进制转义，运行时再拼。
set -u

repo_root=$(git rev-parse --show-toplevel)
cd "$repo_root" || exit 2

w_dai=$(printf '\345\276\205')          # U+5F85
w_tian=$(printf '\345\241\253')         # U+586B
w_xie=$(printf '\345\206\231')          # U+5199
w_zhong=$(printf '\344\270\255')        # U+4E2D
w_t=$(printf '\124')                    # T
w_b=$(printf '\102')                    # B
w_d=$(printf '\104')                    # D

pat="${w_dai}${w_tian}|${w_tian}${w_xie}${w_zhong}|${w_t}${w_b}${w_d}"

scope=$(git ls-files .scratch/wisp/probes/263/r1 scripts/slo-check.ps1)
# 本次交件"将要 add"的件在 add 之前不在 git ls-files 里（尺不能只扫已入库的那一半），
# 所以允许用参数把它们并进来；实参逐字记在 placeholder-check-2.txt 的 scope 行。
for extra in "$@"; do
    if [ -f "$extra" ]; then
        scope="$scope
$extra"
    fi
done
scope=$(printf '%s\n' "$scope" | sort -u)
n_scope=$(printf '%s\n' "$scope" | grep -c .)

echo "ruler: 对每枚交证件跑 grep -cE <三枚占位词 pattern>（pattern 由 UTF-8 字节转义在运行时拼出，逐字不落在任何交付物里）"
echo "scope: git ls-files 现跑＝本目录全部 tracked 件＋scripts/slo-check.ps1，再加参数并入的待 add 件（实参：$*），共 $n_scope 枚"
echo "--- 命中（只列非零）---"
total=0
evidence_hits=0
ps1_hits=0
for f in $scope; do
    n=$(grep -cE "$pat" "$f" 2>/dev/null)
    n=${n:-0}
    if [ "$n" -gt 0 ]; then
        echo "$n  $f"
    fi
    total=$((total + n))
    case "$f" in
        *fix-and-readings.md) evidence_hits=$n ;;
        scripts/slo-check.ps1) ps1_hits=$n ;;
    esac
done
echo "--- 汇总 ---"
echo "tracked 交付件 = $n_scope；占位词总命中 = $total（其中 scripts/slo-check.ps1 = $ps1_hits，证据件 fix-and-readings.md = $evidence_hits）"
