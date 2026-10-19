#!/bin/bash
# 名册归属尺：只在某发出现的名字，归到哪个包的源码（⛔ 假设分隔符，锚行首）。
set -u
REPO='/d/work/workspace/projects plans/Wisp'
cd "$REPO"
D=/d/tmp/wisp/20261019/227v2/namesets
for k in fail-top; do
  echo "== $k only-in-mine 归属 =="
  comm -13 "$D/orch-$k.txt" "$D/mine-$k.txt" | while read -r n; do
    hits=$(grep -rl "^func $n(" --include='*_test.go' internal cmd 2>/dev/null | head -3 | tr '\n' ' ')
    printf '%s  <<  %s\n' "$n" "$hits"
  done
done
echo "== only-in-orch 的 15 枚：本腿 paired-backslash/full-B 里的颜色（逐名）=="
while read -r n; do
  [ -z "$n" ] && continue
  v=$(grep -m1 -E "^(    )*--- (PASS|FAIL|SKIP): ${n}( |\$|/)" /d/tmp/wisp/20261019/227v2/full-B.log | sed 's/^ *//')
  printf '%-72s %s\n' "$n" "${v:-NOT-FOUND-IN-MY-SHOT}"
done < <(comm -23 "$D/orch-fail-top.txt" "$D/mine-fail-top.txt")
echo "== cmd/wisp 源码用例总数 ⇄ 两发块内实际跑到的枚数 =="
grep -h '^func Test' cmd/wisp/*_test.go | sed 's/^func \(Test[A-Za-z0-9_]*\).*/\1/' | sort -u > /tmp/cmdwisp-src.txt
wc -l < /tmp/cmdwisp-src.txt
for pair in "orch:/d/tmp/wisp/20261019/orch227-b/full.log" "mine:/d/tmp/wisp/20261019/227v2/full-B.log"; do
  tag=${pair%%:*}; L=${pair#*:}
  grep -E '^=== RUN   [A-Z]' "$L" | sed 's/^=== RUN   //' | sed 's# .*##' | sort -u > /tmp/run-$tag.txt
  n=$(comm -12 /tmp/cmdwisp-src.txt /tmp/run-$tag.txt | wc -l)
  echo "$tag: cmd/wisp 定义且真跑到的顶层用例 = $n / $(wc -l < /tmp/cmdwisp-src.txt)"
done
