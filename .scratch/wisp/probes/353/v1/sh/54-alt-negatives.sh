#!/bin/sh
# 353-v1 R54 = 交替批底片在位⇄是否含 leaf-symlink 分子行（量的）
cd /d/tmp/wisp/20261020/353r1/logs || exit 9
echo "raw 枚数=$(ls -1 alt-*.raw.txt | wc -l)  summary 行数=$(wc -l < alt-summary.txt)"
for f in alt-*.raw.txt; do
  printf '%-26s lines=%-4s leafRow=%s symlinkDirRow=%s\n' "$f" "$(wc -l < "$f")" \
    "$(grep -cE 'form=leaf-symlink attempts=400' "$f" || true)" \
    "$(grep -cE 'form=symlink-dir attempts=400' "$f" || true)"
done
echo "== summary 逐字 =="
cat alt-summary.txt
