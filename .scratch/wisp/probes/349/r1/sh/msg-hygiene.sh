#!/bin/sh
# 提交正文卫生门：一枚文件里 ① U+00A7 码位计数（want 0）② U+00B7（合法，摊出来证明没混）
# ③ 字节数（want 非 0）。用法：sh msg-hygiene.sh <file...>
set -u
for f in "$@"; do
  [ -f "$f" ] || { echo "MISSING $f"; continue; }
  a7=$(perl -CSD -0ne 'print scalar(() = /\x{00A7}/g)' "$f")
  b7=$(perl -CSD -0ne 'print scalar(() = /\x{00B7}/g)' "$f")
  printf '%-70s U00A7=%s U00B7=%s bytes=%s\n' "$f" "${a7:-0}" "${b7:-0}" "$(wc -c < "$f")"
done
