#!/bin/sh
# 349-r1 命门结构尺：对一枚 .go 文件算三把「剥注释后」的可执行字节 sha256。
# GS1 = 票 343 交付那一形（裁腿＋编排者各自独立复现过的形）：
#       sed 's://.*::'  |  tr -d ' \t\r'  |  sed '/^$/d'
# GS2 = 票 343 第二把独立形：先整行删「行首非空白即 //」的注释行，再 GS1 余下两步。
# GB  = 本腿自加的最严一把：只整行删除「行首非空白即 //」的注释行，其余字节一字不动
#       （不 tr、不删空行）⇒ GB 两侧同串＝除整行注释外没有任何一个字节变过。
# 只读数、不写任何文件。用法：sh strip-sha.sh <file> <label>
set -u
F="$1"
LABEL="$2"
[ -f "$F" ] || { echo "MISSING $F"; exit 91; }
GS1=$(sed 's://.*::' "$F" | tr -d ' \t\r' | sed '/^$/d' | sha256sum | cut -d' ' -f1)
GS1L=$(sed 's://.*::' "$F" | tr -d ' \t\r' | sed '/^$/d' | wc -l)
GS2=$(sed -e '/^[[:space:]]*\/\//d' -e 's://.*::' "$F" | tr -d ' \t\r' | sed '/^$/d' | sha256sum | cut -d' ' -f1)
GS2L=$(sed -e '/^[[:space:]]*\/\//d' -e 's://.*::' "$F" | tr -d ' \t\r' | sed '/^$/d' | wc -l)
GB=$(grep -v '^[[:space:]]*//' "$F" | sha256sum | cut -d' ' -f1)
GBL=$(grep -c -v '^[[:space:]]*//' "$F" || true)
RAW=$(sha256sum "$F" | cut -d' ' -f1)
printf 'label=%s file=%s\n' "$LABEL" "$F"
printf '  RAW_sha256=%s bytes=%s lines=%s\n' "$RAW" "$(wc -c < "$F")" "$(wc -l < "$F")"
printf '  GS1_sha256=%s lines=%s\n' "$GS1" "$GS1L"
printf '  GS2_sha256=%s lines=%s\n' "$GS2" "$GS2L"
printf '  GB_sha256=%s lines=%s\n' "$GB" "$GBL"
