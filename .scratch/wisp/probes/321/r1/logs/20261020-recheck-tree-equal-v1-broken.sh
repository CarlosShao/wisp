#!/bin/bash
# 321-r1 终局尺：证明成对台面的字节 == 交件后最终 HEAD 的字节（本腿跑完 7 笔后 HEAD 漂过 3 笔，
# 那 3 笔逐枚 name-only 只含 docs/票面/派单存档＝零 Go 字节 ⇒ 台面等价；这把尺把它量出来，⛔ 靠说）。
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 9
R=/d/tmp/wisp/20261020/321r1
echo "HEAD now = $(git rev-parse --short HEAD)"
echo "--cmd/wisp 整包 tracked blob 串比对（HEAD ⇄ treeB 物理字节）--"
bad=0; n=0
while IFS= read -r line; do
  h=$(printf '%s' "$line" | awk '{print $1}')
  p=$(printf '%s' "$line" | awk '{print $4}')
  n=$((n+1))
  disk=$(sha256sum "$R/treeB/$p" 2>/dev/null | awk '{print $1}')
  headh=$(git hash-object "$R/treeB/$p" 2>/dev/null)
  if [ "$headh" != "$h" ]; then echo "MISMATCH $p"; bad=$((bad+1)); fi
done < <(git ls-tree -r HEAD --format='%(objecttype) %(objectname) %(objectmode) %(path)' | awk '$1=="blob" && $4 ~ /^cmd\/wisp\//')
echo "checked=$n mismatches=$bad"
echo "--internal/winsec blob 是否也等价（352 那枚一行已落＝对两边同串）--"
git hash-object "$R/treeB/internal/winsec/winsec.go"; git rev-parse HEAD:internal/winsec/winsec.go
echo "--判据件 blob 等价--"
git hash-object "$R/treeB/cmd/wisp/exec_holder_toggle_321_test.go"; git rev-parse HEAD:cmd/wisp/exec_holder_toggle_321_test.go
