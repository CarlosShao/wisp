#!/bin/bash
# 353-v1 R30 = 件身尺（CR 枚数／分节码位／行数／字节／占位符／末字节；对入库 blob 或工作树同一形量）
# 用法: sh 30-body-rulers.sh <ref|WT> <path...>
cd '/d/work/workspace/projects plans/Wisp' || exit 9
MODE="$1"; shift
SCRATCH=/d/tmp/wisp/20261020/353v1/_r30.bin
S7=$(printf '\xc2\xa7')
for p in "$@"; do
  if [ "$MODE" = "blob" ]; then git cat-file blob "$p" > "$SCRATCH" 2>/dev/null; SRC="$p";
  else cp -f "$p" "$SCRATCH" 2>/dev/null; SRC="WT:$p"; fi
  [ -f "$SCRATCH" ] || { echo "MISSING $p"; continue; }
  bytes=$(wc -c < "$SCRATCH")
  lines=$(wc -l < "$SCRATCH")
  cr=$(tr -cd '\r' < "$SCRATCH" | wc -c)
  s7=$(grep -o "$S7" "$SCRATCH" | wc -l)
  ph=$(grep -cE 'TODO|FIXME|XXX' "$SCRATCH" || true)
  tailbyte=$(tail -c 1 "$SCRATCH" | od -An -tx1 | tr -d ' \n')
  printf '%s\n  BYTES=%s LINES=%s CR=%s SECT=%s PLACEHOLDERS=%s LASTBYTE=%s\n' "$SRC" "$bytes" "$lines" "$cr" "$s7" "$ph" "$tailbyte"
  if [ "$MODE" = "WT" ]; then
    printf '  ICONV_CHECK='; iconv -f UTF-8 -t UTF-8 "$p" > /dev/null && echo OK || echo BAD
  fi
done
