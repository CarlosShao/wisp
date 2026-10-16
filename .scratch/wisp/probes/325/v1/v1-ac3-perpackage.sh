#!/bin/sh
# 325-v1 AC#3 逐包四数：从整档 stdout 切块归因 ＋ 用三枚小包独立复跑验这把切块尺
# 口径（与脚本自打的总计行同尺）：RUN=^=== RUN  PASS=^--- PASS  FAIL=^--- FAIL  SKIP=^--- SKIP（只数顶层）
# 归因法：go test 把一个包的 -v 输出连续打完，再以 ^ok/FAIL 结果行收尾 ⇒ 以上一枚结果行为界切块。
# 自证：11 块四数之和必须等于脚本自己打印的总计行（不等＝切块被并行打乱，读数作废、改逐包独立跑）。
set -u
tree=${1:?tree}; log=${2:?log}
OUT=${3:?out}
{
echo "V1 AC#3 per-package four numbers"
echo "tree=$tree"
echo "log=$log  lines=$(wc -l < "$log")"
echo "totals line: $(grep 'four numbers' "$log")"
echo
echo '===== A. SEGMENTED ATTRIBUTION (awk, block per package result line)'
awk '
/^=== RUN/   { blk["RUN"]++ }
/^--- PASS/  { blk["PASS"]++ }
/^--- FAIL/  { blk["FAIL"]++ }
/^--- SKIP/  { blk["SKIP"]++ }
/^(ok|FAIL)[ \t]+[^ \t]+[ \t]+/ {
    pkg=$2
    printf "%-46s RUN=%-5s PASS=%-5s FAIL=%-5s SKIP=%-5s\n", pkg, blk["RUN"], blk["PASS"], blk["FAIL"], blk["SKIP"]
    s["RUN"]+=blk["RUN"]; s["PASS"]+=blk["PASS"]; s["FAIL"]+=blk["FAIL"]; s["SKIP"]+=blk["SKIP"]
    n++
    delete blk
}
END {
    printf "BLOCKS=%d  SUM RUN=%d PASS=%d FAIL=%d SKIP=%d\n", n, s["RUN"], s["PASS"], s["FAIL"], s["SKIP"]
}' "$log"
echo
echo '===== B. RED ROSTER (every top-level --- FAIL line, verbatim, with its message block)'
if grep -q '^--- FAIL' "$log"; then
    grep -n '^--- FAIL' "$log" | sed 's/^/  /'
else
    echo '  (empty: zero top-level --- FAIL lines in this run)'
fi
echo "  top-level FAIL lines=$(grep -c '^--- FAIL' "$log" || true)   package FAIL lines=$(grep -c '^FAIL' "$log" || true)"
echo '  positive control for section B: --- PASS lines='$(grep -c '^--- PASS' "$log")' (a ruler that can see lines must not read 0 here)'
echo
echo '===== C. INDEPENDENT CHECK OF THE SEGMENTATION RULER: three small packages run alone through the strict runner'
for p in ./internal/perm/ ./cmd/llmrecord/ ./internal/plugin/; do
    ( cd "$tree" && env -u GOTMPDIR -u WISP_TEST_GOTMP_ROOT \
        PATH="$tree/third_party/sherpa-onnx:$PATH" \
        sh tools/d22scan/runtests.sh "$p" -count=1 2>&1 | grep '^runtests.sh: ' | sed "s|^|  $p |" )
done
echo
echo '===== D. THE FOUR RISK CASES AND THE TWO WORDINGS ARE ABSENT HERE (no Invalid switch / Parameter format in a green run)'
echo "  Invalid_switch_hits=$(grep -cF 'Invalid switch' "$log" || true)  Parameter_format_hits=$(grep -cF 'Parameter format not correct' "$log" || true)   (positive control: PASS lines above are non-zero)"
} > "$OUT" 2>&1
echo "rc=$? out=$OUT"
