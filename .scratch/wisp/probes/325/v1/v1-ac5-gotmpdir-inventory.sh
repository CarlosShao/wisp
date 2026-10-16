#!/bin/sh
# 325-v1 AC#5：全仓 GOTMPDIR / WISP_TEST_GOTMP_ROOT / gotmp_* 赋值逐枚清单
# 只跑在仓外导出树上（⛔ 把在飞写腿的字节算成第五处同形）。
# 每枚要答：这处是"当值用"（真流向 export GOTMPDIR）还是"只当展示用"（echo/探针显示值）。
# 用法：sh v1-ac5-gotmpdir-inventory.sh <tree> <out>
set -u
tree=${1:?tree}; out=${2:?out}
cd "$tree" || exit 9
{
    echo "V1 AC#5 GOTMPDIR assignment inventory"
    echo "tree=$tree  start=$(date '+%F %T %z')"
    echo "sha(tools/d22scan/runtests.sh)=$(sha256sum tools/d22scan/runtests.sh | cut -d' ' -f1)"
    echo
    echo '===== A. PRODUCTION SURFACE (scripts/ tools/ cmd/ internal/) : every line mentioning GOTMPDIR or the root vars'
    grep -rn 'GOTMPDIR\|WISP_TEST_GOTMP_ROOT\|gotmp_run=\|gotmp_base=' scripts tools cmd internal \
        2>/dev/null | sed 's/^/  /'
    echo
    echo '===== B. CLASS-LEVEL COUNTS in the production surface (fixed strings)'
    echo "  class_total_lines=$(grep -rn 'GOTMPDIR\|WISP_TEST_GOTMP_ROOT\|gotmp_run=\|gotmp_base=' scripts tools cmd internal 2>/dev/null | wc -l)"
    echo "  export_GOTMPDIR_lines=$(grep -rn 'export GOTMPDIR' scripts tools cmd internal 2>/dev/null | wc -l)   # 当值用：只有这些真的把根交给子进程"
    echo "  gotmp_run_assign=$(grep -rn 'gotmp_run=' scripts tools cmd internal 2>/dev/null | wc -l)"
    echo "  gotmp_base_assign=$(grep -rn 'gotmp_base=' scripts tools cmd internal 2>/dev/null | wc -l)"
    echo
    echo '===== C. THE FIVE ENTRIES THAT REALLY START go test (each export site + the line that builds the root)'
    for f in scripts/portable-tests.sh scripts/winsec-tests.sh scripts/wisp-cli-tests.sh tools/d22scan/runtests.sh scripts/d22scan.sh; do
        echo "--- $f"
        grep -n 'export GOTMPDIR\|gotmp_run=\|gotmp_base=' "$f" 2>/dev/null | sed 's/^/    /' || echo '    (no export site)'
    done
    echo
    echo '===== D. DISPLAY-ONLY SITE (ticket 311 guard reference path, folded to / ON PURPOSE)'
    grep -rn 'wisp_311_ref_path' scripts tools cmd internal 2>/dev/null | sed 's/^/  /'
    echo
    echo '===== E. MIXED-TAIL HUNT OVER THE WHOLE PRODUCTION SURFACE (fixed strings)'
    echo '  CONTROL FIRST (a / -fold string that really is in these files, 4 occurrences expected):'
    printf '      control  ${LOCALAPPDATA//\\//}  hits=%s\n' \
        "$(grep -rnF '${LOCALAPPDATA//\\//}' scripts tools 2>/dev/null | wc -l)"
    echo '  (the 0s below are only readable because that control is 4, i.e. the same grep -rnF ruler sees these files)'
    for pat in '}/d22scan-' '}/winsec-' '}/cli-' '}/portable-' '/wisp-test-gotmp'; do
        hits=$(grep -rnF "$pat" scripts tools cmd internal .github 2>/dev/null)
        printf '  pat=%-16s hits=%s\n' "$pat" "$(printf '%s' "$hits" | grep -c . || true)"
        [ -n "$hits" ] && printf '%s\n' "$hits" | sed 's/^/      /'
    done
    echo
    echo '===== F. THE SAME HUNT OVER .github/ ci.yml (the CI leg writes WISP_TEST_GOTMP_ROOT with / by design)'
    grep -rn 'WISP_TEST_GOTMP_ROOT\|GOTMPDIR' .github 2>/dev/null | sed 's/^/  /' || echo '  (.github names no GOTMPDIR/WISP_TEST_GOTMP_ROOT line)'
    echo
    echo '===== F2. FIFTH MENTIONS THAT ARE NOT ASSIGNMENT SITES (comments/registrations, listed so nobody counts them as a new wound)'
    grep -rn 'Invalid switch\|mklink' internal/winsec/*_test.go 2>/dev/null | sed 's/^/  /' | head -10
    echo
    echo '===== G. .scratch PROBE CODE (not a production entry; listed only so the 全仓 claim is auditable)'
    grep -rln 'export GOTMPDIR' .scratch 2>/dev/null | head -20 | sed 's/^/  /'
    echo "  files_in_.scratch_exporting_GOTMPDIR=$(grep -rln 'export GOTMPDIR' .scratch 2>/dev/null | wc -l)"
    echo "end=$(date '+%F %T %z')"
} > "$out" 2>&1
echo "rc=$? out=$out"
