#!/usr/bin/env bash
# 尺名册用的分类器 v3（本腿自写，可复跑；只读，不改任何仓内件）。
# 尺的形状（一把尺一个口径，全部可复跑）：
#   同名件 := `ls docs/evidence/s1/ | grep -E "^0*<票号>([-/]|$)"`（票号开头＝这枚票名下的件；
#            不用"票号出现在文件名里"那种含法——`136-ac14-*` 是票 136 的 AC#14、不是票 14）
#   判 VERDICT（该票名下的非实现者终裁件）:=
#       第一行标题命中 VERDICT 词族 **且** 标题行不命中实现侧自称词族
#       **且** 前 8 行没有"本件不是验收件／不是裁决表／同一程写码又自证"这类明写的免责句
#   其余同名件一律记 OTHER（实现侧交件／普查／读数表／取证件／索引／目录）
#   逐票三档：A 名下有 VERDICT 件 ／ B 名下只有 OTHER 件 ／ C 名下零枚同名件
set -u
cd "D:/work/workspace/projects plans/Wisp"
VERD='对抗验收|裁决|验收腿|验收程|验收方|验收件|终裁|acceptance|acceptor|非实现者|独立复算|独立验收|独立验证|抽验|回判|reaccept|ratification|accept-r'
IMPLTITLE='实现侧交件|实现者自证|实现程自证|实现件|证据件|交件证据|预做|普查|读数表|凭据索引|只读取证|只读设计核|正向对照|补证|写码位|落地写腿|修复程|实现程|实现方|归因尺|存量盘点'
DISCLAIM='本件不是验收件|不是验收件|本文件不是裁决表|本件不是裁决表|不是对抗验收表|同一程写码又自证|本件不声称结掉|本文件是\*\*实现者自证'

: > .scratch/wisp/probes/evidence-close/3/file-identity.txt
while read -r n f; do
  if [ -d "docs/evidence/s1/$f" ]; then echo "$n|$f|DIR|目录（台件，非裁决表）"; continue; fi
  t=$(head -1 "docs/evidence/s1/$f" | tr -d '\r')
  d=$(head -8 "docs/evidence/s1/$f" | tr -d '\r' | grep -m1 -E "$DISCLAIM")
  verdict=N
  printf '%s' "$t" | grep -qE "$VERD" && ! printf '%s' "$t" | grep -qE "$IMPLTITLE" && verdict=Y
  [ -n "$d" ] && verdict=N
  if [ "$verdict" = Y ]; then echo "$n|$f|VERDICT|${t:0:60}"; else
    why=$(printf '%s' "$t" | grep -oE "$IMPLTITLE" | head -1); echo "$n|$f|OTHER|${why:-标题无终裁词}｜${t:0:52}"; fi
done < .scratch/wisp/probes/evidence-close/3/named-own-pairs.txt >> .scratch/wisp/probes/evidence-close/3/file-identity.txt

: > .scratch/wisp/probes/evidence-close/3/ticket-abc.txt
while read -r n; do
  v=$(grep "^$n|" .scratch/wisp/probes/evidence-close/3/file-identity.txt | awk -F'|' '$3=="VERDICT"{print $2}' | tr '\n' ' ')
  o=$(grep "^$n|" .scratch/wisp/probes/evidence-close/3/file-identity.txt | awk -F'|' '$3!="VERDICT"{print $2}' | tr '\n' ' ')
  if [ -n "$v" ]; then echo -e "A\t$n\t$v\t$o" >> .scratch/wisp/probes/evidence-close/3/ticket-abc.txt
  elif [ -n "$o" ]; then echo -e "B\t$n\t\t$o" >> .scratch/wisp/probes/evidence-close/3/ticket-abc.txt
  else echo -e "C\t$n\t\t" >> .scratch/wisp/probes/evidence-close/3/ticket-abc.txt; fi
done < .scratch/wisp/probes/evidence-close/3/done-nums.txt

echo "--- 三档计数 ---"
for b in A B C; do printf '%s %s\n' "$b" "$(grep -c "^$b	" .scratch/wisp/probes/evidence-close/3/ticket-abc.txt)"; done
echo "--- VERDICT 件总数 / OTHER 件总数 / DIR 条目 ---"
for k in VERDICT OTHER DIR; do printf '%s %s\n' "$k" "$(grep -c "|$k|" .scratch/wisp/probes/evidence-close/3/file-identity.txt)"; done
echo "--- B 档（只有旁证件） ---"; grep "^B	" .scratch/wisp/probes/evidence-close/3/ticket-abc.txt
echo "--- C 档（零枚同名件） ---"; grep "^C	" .scratch/wisp/probes/evidence-close/3/ticket-abc.txt
