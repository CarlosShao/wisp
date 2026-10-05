#!/usr/bin/env bash
# 尺名册用的分类器 v2（本腿自写，可复跑；只读，不改任何仓内件）。
# 步骤：
#   1) 同名件＝第一列票号 + `-` 或 `.` 或目录开头的件（由 done-nums.txt × ls docs/evidence/s1 生成）
#   2) 逐件读**前 20 行**的身份声明：
#      VERDICT（非实现者终裁候选）正词：非实现者|验收腿|验收程|验收方|acceptor-|终裁|对抗验收|独立复算|抽验|回判|验收表|裁决表
#      但**先**看实现侧自称（一票否决）：实现程自证|实现侧交件|实现件 r|修复程|写码代理|实现程（写代码）|实现方第|落地写腿|证据件|交件证据|预做|普查|读数表|凭据索引|只读取证|只读预做|只读设计核|正向对照
#      → 命中实现侧词 ⇒ OTHER
#   3) 逐票汇总三档：A 有 VERDICT 同名件 ／ B 只有 OTHER 同名件 ／ C 零枚同名件
set -u
cd "D:/work/workspace/projects plans/Wisp"
PAIRS=.scratch/wisp/probes/evidence-close/2/named-own-pairs.txt
IMPL='实现程自证|实现侧交件|实现件 r|修复程|写码代理|实现程（写代码）|实现方第|落地写腿|证据件|交件证据|预做|普查|读数表|凭据索引|只读取证|只读预做|只读设计核|正向对照|写码位|装成规矩|补证'
VERD='非实现者|验收腿|验收程|验收方|acceptor-|终裁|对抗验收|独立复算|抽验|回判|验收表|裁决表|独立验收|独立验证|reaccept|ratification|accept-r'

while read -r n f; do
  if [ -d "docs/evidence/s1/$f" ]; then echo "$n|$f|DIR"; continue; fi
  head20=$(head -20 "docs/evidence/s1/$f" | tr -d '\r')
  if printf '%s' "$head20" | grep -qE "$IMPL"; then echo "$n|$f|OTHER"; continue; fi
  if printf '%s' "$head20" | grep -qE "$VERD"; then echo "$n|$f|VERDICT"; else echo "$n|$f|OTHER"; fi
done < "$PAIRS" > .scratch/wisp/probes/evidence-close/2/file-identity3.txt

while read -r n; do
  v=$(grep "^$n|" .scratch/wisp/probes/evidence-close/2/file-identity3.txt | grep '|VERDICT$' | cut -d'|' -f2 | tr '\n' ' ')
  o=$(grep "^$n|" .scratch/wisp/probes/evidence-close/2/file-identity3.txt | grep -v '|VERDICT$' | cut -d'|' -f2 | tr '\n' ' ')
  if [ -n "$v" ]; then echo "A $n | VERDICT: $v | OTHER: $o";
  elif [ -n "$o" ]; then echo "B $n | OTHERONLY: $o";
  else echo "C $n | NONE"; fi
done < .scratch/wisp/probes/evidence-close/2/done-nums.txt > .scratch/wisp/probes/evidence-close/2/ticket-abc3.txt

echo "--- bucket counts (A=有非实现者终裁同名件 / B=只有同名旁证件 / C=零枚同名件) ---"
for b in A B C; do printf '%s %s\n' "$b" "$(grep -c "^$b " .scratch/wisp/probes/evidence-close/2/ticket-abc3.txt)"; done
echo "--- B ---"; grep '^B ' .scratch/wisp/probes/evidence-close/2/ticket-abc3.txt
echo "--- C ---"; grep '^C ' .scratch/wisp/probes/evidence-close/2/ticket-abc3.txt
