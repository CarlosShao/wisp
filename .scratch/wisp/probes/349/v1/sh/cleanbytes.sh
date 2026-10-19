#!/bin/sh
# 字节级清洗：剥 ANSI 转义序列＋裸控制字符（保留 \t \n），⛔ 做 UTF-8 解码⇒不会因非法字节序列中止。
# 用法：sh cleanbytes.sh <src> <dst>
set -u
SRC="$1"; DST="$2"
perl -0777 -pe 'BEGIN{binmode STDIN; binmode STDOUT} s/\x1b\[[0-9;?]*[A-Za-z]//g; s/[\x00-\x08\x0b\x0c\x0e-\x1f\x7f]/?/g;' "$SRC" > "$DST"
printf 'src=%s bytes=%s -> dst=%s bytes=%s ctrl=%s ESC=%s\n' "$SRC" "$(wc -c < "$SRC")" "$DST" "$(wc -c < "$DST")" \
  "$(tr -cd '\000-\010\013\014\016-\037' < "$DST" | wc -c)" "$(tr -cd '\033' < "$DST" | wc -c)"
