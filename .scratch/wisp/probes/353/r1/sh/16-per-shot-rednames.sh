#!/bin/sh
# 353-r1 逐名定因尺＝对每发 raw 只数名、⛔ 拿四数当定论（公共块 颜色尺："枚数⛔ 是单枚字节的函数"）。
set -u
cd /d/tmp/wisp/20261020/353r1/logs || exit 9
for f in recheck-treeB2-*.raw.txt recheck-treeL-*.raw.txt recheck2-treeB2-*.raw.txt recheck2-treeL-*.raw.txt; do
  [ -f "$f" ] || continue
  echo "== $f =="
  grep -E '^--- FAIL' "$f" | sed 's/^--- FAIL: /   RED /; s/ *([0-9.]*s)$//'
  grep -E 'truncated=[0-9]+' "$f" | sed 's/^/   NUM /' | head -2
done
