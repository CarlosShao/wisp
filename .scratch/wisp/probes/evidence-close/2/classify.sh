#!/usr/bin/env bash
# 尺名册用的分类器（本腿自写，可复跑）。
# 只对 docs/evidence/s1/ 里"以票号命名的件"分类：
#   VERDICT = 第一行标题含验收/裁决词，且不含实现侧自称
#   OTHER   = 其余（实现侧交件／普查／读数台件／索引／目录）
set -u
cd "D:/work/workspace/projects plans/Wisp"
PAIRS=.scratch/wisp/probes/evidence-close/2/named-own-pairs.txt
VERD='对抗验收|裁决|验收腿|验收件|终裁|acceptance|独立复算|独立验收|独立验证|verdict|reaccept|ratification|accept-r|抽验|重判|复验'
IMPL='实现侧交件|实现者自证|实现程自证|写码代理|实现程（写代码）|只读预做|普查|读数台|证据件 —|交件证据'
while read -r n f; do
  if [ -d "docs/evidence/s1/$f" ]; then echo "$n|$f|DIR"; continue; fi
  t=$(head -1 "docs/evidence/s1/$f" | tr -d '\r')
  if printf '%s' "$t" | grep -qE "$IMPL"; then echo "$n|$f|OTHER"; continue; fi
  if printf '%s' "$t" | grep -qE "$VERD"; then echo "$n|$f|VERDICT"; else echo "$n|$f|OTHER"; fi
done < "$PAIRS" > .scratch/wisp/probes/evidence-close/2/file-identity2.txt

# 逐票汇总：有 VERDICT 件 / 只有 OTHER 件 / 零枚同名件
while read -r n; do
  v=$(grep "^$n|" .scratch/wisp/probes/evidence-close/2/file-identity2.txt | grep '|VERDICT$' | cut -d'|' -f2 | tr '\n' ' ')
  o=$(grep "^$n|" .scratch/wisp/probes/evidence-close/2/file-identity2.txt | grep -v '|VERDICT$' | cut -d'|' -f2 | tr '\n' ' ')
  if [ -n "$v" ]; then echo "A $n | VERDICT: $v | OTHER: $o";
  elif [ -n "$o" ]; then echo "B $n | OTHERONLY: $o";
  else echo "C $n | NONE"; fi
done < .scratch/wisp/probes/evidence-close/2/done-nums.txt > .scratch/wisp/probes/evidence-close/2/ticket-abc.txt

echo "--- bucket counts (A=有终裁表 / B=只有同名旁证件 / C=零枚同名件) ---"
for b in A B C; do printf '%s %s\n' "$b" "$(grep -c "^$b " .scratch/wisp/probes/evidence-close/2/ticket-abc.txt)"; done
echo "--- B ---"; grep '^B ' .scratch/wisp/probes/evidence-close/2/ticket-abc.txt
echo "--- C ---"; grep '^C ' .scratch/wisp/probes/evidence-close/2/ticket-abc.txt
