#!/bin/bash
# 尺：cmd/wisp 源码里定义的顶层用例 ⇄ 两发实际跑到的名集，逐名差集＋差集归属文件。
set -u
REPO='/d/work/workspace/projects plans/Wisp'
cd "$REPO"
grep -h '^func Test' cmd/wisp/*_test.go | sed 's/^func \(Test[A-Za-z0-9_]*\).*/\1/' | sort -u > /tmp/src.txt
for pair in "orch:/d/tmp/wisp/20261019/orch227-b/full.log" "mine:/d/tmp/wisp/20261019/227v2/full-B.log"; do
  tag=${pair%%:*}; L=${pair#*:}
  grep -E '^=== RUN   ' "$L" | sed 's/^=== RUN   //; s# .*##' | sort -u > /tmp/r-$tag.txt
  echo "== $tag : 定义但未跑到的 cmd/wisp 名（前 30 枚）=="
  comm -23 /tmp/src.txt /tmp/r-$tag.txt | tee /tmp/miss-$tag.txt | head -30
  echo "   count=$(wc -l < /tmp/miss-$tag.txt)"
done
echo "== 差集逐名归属文件（mine 侧）=="
while read -r n; do
  [ -z "$n" ] && continue
  f=$(grep -rl "^func $n(" --include='*_test.go' cmd/wisp 2>/dev/null | head -1)
  printf '%s  %s\n' "$n" "$f"
done < /tmp/miss-mine.txt
