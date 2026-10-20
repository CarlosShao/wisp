#!/bin/sh
# 353-v1 R47 = 纯追加的"逐行"尺（numstat 删除列 + diff 行形计数；★前缀 cmp 只适用末尾追加，本次含表内插入⇒具名作废那把）
cd '/d/work/workspace/projects plans/Wisp' || exit 9
P="$1"
echo "== numstat =="; git diff --numstat HEAD -- "$P"
A=$(git diff -U0 HEAD -- "$P" | grep -c '^+' || true)
D=$(git diff -U0 HEAD -- "$P" | grep -c '^-' || true)
echo "== diff 行计数（含 +++/--- 头各 1 枚）：加=$A 删=$D ⇒ 净增=$((A-1)) 净删=$((D-1))"
echo "== 被删的行（期望 0 行）=="
git diff -U0 HEAD -- "$P" | grep '^-' | grep -v '^---' | sed 's/^/  /'
echo "== 上下文完整性另一把尺：旧 blob 的每一行是否都仍在新工作树里（逐行 grep -Fx）=="
git cat-file blob "HEAD:$P" > /d/tmp/wisp/20261020/353v1/_old.txt
missing=0
while IFS= read -r line; do
  grep -Fqx -- "$line" "$P" || { missing=$((missing+1)); echo "  LOST: ${line:0:80}"; }
done < /d/tmp/wisp/20261020/353v1/_old.txt
echo "  LOST_LINES=$missing"
