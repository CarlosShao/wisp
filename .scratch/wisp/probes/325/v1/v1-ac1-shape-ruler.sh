#!/bin/sh
# 325-v1 AC#1 形状尺（固定串 grep -F，配已知存在的正控）
# ★纪律：⛔ 用 grep -rn 的反斜杠正则形数命中（反斜杠枚数与文件字节一旦对不上就假 0）；
#   本尺一律 grep -F，且每一段自带一枚"已知真身存在"的正控，报 0 之前先看正控响没响。
# 用法：sh v1-ac1-shape-ruler.sh <树根>
set -eu
root=${1:?need a tree root}
cd "$root"

echo "TREE=$root"
echo '--- FILE BYTE GROUND TRUTH (od -c of the gotmp_run line, one per entry)'
FILES='tools/d22scan/runtests.sh scripts/portable-tests.sh scripts/winsec-tests.sh scripts/wisp-cli-tests.sh'
for f in $FILES; do
    n=$(grep -n 'gotmp_run=' "$f" | head -1 | cut -d: -f1)
    printf '%s:%s  ' "$f" "$n"
    sed -n "${n}p" "$f" | od -c | tr -s ' ' | head -8
    echo
done

echo '--- POSITIVE CONTROL: grep -cF of a two-backslash string that MUST be present in the root builder'
for f in $FILES; do
    printf '%-32s pos_two_bs_hits=%s   pos_one_bs_hits=%s\n' "$f" \
        "$(grep -cF '\\'wisp-test-gotmp "$f" || true)" \
        "$(grep -cF '\'wisp-test-gotmp "$f" || true)"
done
echo

echo '--- TARGET: tail-join separator on the gotmp_run line (fixed strings)'
for f in $FILES; do
    n=$(grep -n 'gotmp_run=' "$f" | head -1 | cut -d: -f1)
    line=$(sed -n "${n}p" "$f")
    bs=$(printf '%s' "$line" | grep -cF '}\\' || true)
    fs=$(printf '%s' "$line" | grep -cF '}/' || true)
    printf '%-32s line=%s  native_join_hits=%s  slash_join_hits=%s\n' "$f" "$n" "$bs" "$fs"
done
echo

echo '--- VERBATIM CROSS-CHECK: "}\\<name>-$(date"'
for pair in 'tools/d22scan/runtests.sh:\\d22scan-' \
            'scripts/winsec-tests.sh:\\winsec-' \
            'scripts/wisp-cli-tests.sh:\\cli-' \
            'scripts/portable-tests.sh:\\portable-'; do
    f=${pair%%:*}
    pat=${pair#*:}
    printf '%-32s pat=%-14s hits=%s\n' "$f" "$pat" "$(grep -cF "$pat"'$(date' "$f" || true)"
done
echo

echo '--- FALSE-RULER DEMONSTRATION (why the regex form is banned for this shape)'
echo '    NB: grep -c counts LINES, and the AC#1 comment block itself quotes \\d22scan-, so a'
echo '    line count above 1 is expected here. The point is the zero, not the number.'
echo "F_two_backslashes       hits=$(grep -cF '}\\d22scan-' tools/d22scan/runtests.sh || true)   <- trustworthy fixed string"
echo "regex_eight_backslashes hits=$(grep -c '\\\\\\\\d22scan-' tools/d22scan/runtests.sh || true)   <- BRE reads 4 literal backslashes = FALSE 0 on a file that really has the shape"
echo "regex_four_backslashes  hits=$(grep -c '\\\\d22scan-' tools/d22scan/runtests.sh || true)   <- BRE reads 2 literal (shape line + comment line)"
echo "regex_two_backslashes   hits=$(grep -c '\\d22scan-' tools/d22scan/runtests.sh || true)   <- BRE reads 1 literal, matches by substring"
echo

echo '--- REPO-WIDE MIXED-TAIL HUNT (fixed strings; 0 must always be read next to the control above)'
for pat in '}/d22scan-' '}/winsec-' '}/cli-' '}/portable-'; do
    printf 'pat=%-14s files=%s\n' "$pat" "$(grep -rlF "$pat" tools scripts 2>/dev/null | tr '\n' ' ')"
done
