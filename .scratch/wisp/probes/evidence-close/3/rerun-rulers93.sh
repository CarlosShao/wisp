#!/usr/bin/env bash
# evidence-close-3 复核尺：把死腿 evidence-close/2 的三版分类器逐版复跑，
# 全部输出改写到本腿探针目录（⛔ 不写 evidence-close/2 的任何件、不写仓内件）。
# 三版输入（named-own-pairs.txt / done-nums.txt）只读引用死腿原件。
set -u
cd "D:/work/workspace/projects plans/Wisp"
D2=.scratch/wisp/probes/evidence-close/2
D3=.scratch/wisp/probes/evidence-close/3
PAIRS=$D2/named-own-pairs.txt
NUMS=.scratch/wisp/probes/evidence-close/3/nums93.txt

# ---------- v1 = classify.sh 的尺（标题行含终裁词族 且 不含实现侧自称） ----------
VERD='对抗验收|裁决|验收腿|验收件|终裁|acceptance|独立复算|独立验收|独立验证|verdict|reaccept|ratification|accept-r|抽验|重判|复验'
IMPL='实现侧交件|实现者自证|实现程自证|写码代理|实现程（写代码）|只读预做|普查|读数台|证据件 —|交件证据'
: > $D3/v1-i93.txt
while read -r n f; do
  if [ -d "docs/evidence/s1/$f" ]; then echo "$n|$f|DIR"; continue; fi
  t=$(head -1 "docs/evidence/s1/$f" | tr -d '\r')
  if printf '%s' "$t" | grep -qE "$IMPL"; then echo "$n|$f|OTHER"; continue; fi
  if printf '%s' "$t" | grep -qE "$VERD"; then echo "$n|$f|VERDICT"; else echo "$n|$f|OTHER"; fi
done < "$PAIRS" > $D3/v1-i93.txt

# ---------- v2 = classify2.sh 的尺（前 20 行；实现侧词一票否决，词表更宽） ----------
IMPL2='实现程自证|实现侧交件|实现件 r|修复程|写码代理|实现程（写代码）|实现方第|落地写腿|证据件|交件证据|预做|普查|读数表|凭据索引|只读取证|只读预做|只读设计核|正向对照|写码位|装成规矩|补证'
VERD2='非实现者|验收腿|验收程|验收方|acceptor-|终裁|对抗验收|独立复算|抽验|回判|验收表|裁决表|独立验收|独立验证|reaccept|ratification|accept-r'
: > $D3/v2-i93.txt
while read -r n f; do
  if [ -d "docs/evidence/s1/$f" ]; then echo "$n|$f|DIR"; continue; fi
  head20=$(head -20 "docs/evidence/s1/$f" | tr -d '\r')
  if printf '%s' "$head20" | grep -qE "$IMPL2"; then echo "$n|$f|OTHER"; continue; fi
  if printf '%s' "$head20" | grep -qE "$VERD2"; then echo "$n|$f|VERDICT"; else echo "$n|$f|OTHER"; fi
done < "$PAIRS" > $D3/v2-i93.txt

# ---------- v3 = classify3.sh 的尺（标题行 + 前 8 行免责句否决）—— 死腿最终版 ----------
VERD3='对抗验收|裁决|验收腿|验收程|验收方|终裁|acceptance|acceptor|非实现者|独立复算|独立验收|独立验证|抽验|回判|reaccept|ratification|accept-r'
IMPL3='实现侧交件|实现者自证|实现程自证|实现件|证据件|交件证据|预做|普查|读数表|凭据索引|只读取证|只读设计核|正向对照|补证|写码位|落地写腿|修复程|实现程|实现方|归因尺|存量盘点'
DISCLAIM='本件不是验收件|不是验收件|本文件不是裁决表|本件不是裁决表|不是对抗验收表|同一程写码又自证|本件不声称结掉|本文件是\*\*实现者自证'
: > $D3/v3-i93.txt
while read -r n f; do
  if [ -d "docs/evidence/s1/$f" ]; then echo "$n|$f|DIR|目录"; continue; fi
  t=$(head -1 "docs/evidence/s1/$f" | tr -d '\r')
  d=$(head -8 "docs/evidence/s1/$f" | tr -d '\r' | grep -m1 -E "$DISCLAIM")
  verdict=N
  printf '%s' "$t" | grep -qE "$VERD3" && ! printf '%s' "$t" | grep -qE "$IMPL3" && verdict=Y
  [ -n "$d" ] && verdict=N
  if [ "$verdict" = Y ]; then echo "$n|$f|VERDICT"; else echo "$n|$f|OTHER"; fi
done < "$PAIRS" > $D3/v3-i93.txt

# ---------- 三版逐票三档 ----------
for v in v1 v2 v3; do
  : > $D3/$v-a93.txt
  while read -r n; do
    idf=$D3/$v-i93.txt
    vhit=$(grep "^$n|" $idf | awk -F'|' '$3=="VERDICT"{print $2}' | tr '\n' ' ')
    oth=$(grep "^$n|" $idf | awk -F'|' '$3!="VERDICT"{print $2}' | tr '\n' ' ')
    if [ -n "$vhit" ]; then echo -e "A\t$n\t$vhit\t$oth" >> $D3/$v-a93.txt
    elif [ -n "$oth" ]; then echo -e "B\t$n\t\t$oth" >> $D3/$v-a93.txt
    else echo -e "C\t$n\t\t" >> $D3/$v-a93.txt; fi
  done < "$NUMS"
  printf '=== %s buckets: ' "$v"
  for b in A B C; do printf '%s=%s ' "$b" "$(grep -c "^$b	" $D3/$v-a93.txt)"; done
  echo
done

# ---------- 与死腿在盘原件逐行对照 ----------
echo "--- v3 vs 死腿 file-identity.txt（VERDICT 名册差集） ---"
diff <(grep '|VERDICT' $D3/v3-i93.txt | cut -d'|' -f1,2 | sort) \
     <(grep '|VERDICT' $D2/file-identity.txt | cut -d'|' -f1,2 | sort) && echo "IDENTICAL"
echo "--- v3 vs 死腿 ticket-abc.txt（三档） ---"
diff <(sort $D3/v3-a93.txt) <(sort $D2/ticket-abc.txt) && echo "IDENTICAL"
echo "--- 死腿 ticket-abc3.txt（＝v2 产物）buckets ---"
for b in A B C; do printf '%s=%s ' "$b" "$(grep -c "^$b	" $D2/ticket-abc3.txt 2>/dev/null || grep -c "^$b " $D2/ticket-abc3.txt)"; done; echo
echo "--- v2 vs 死腿 ticket-abc3 逐票档差 ---"
join -t'|' <(awk -F'\t' '{print $2"|"$1}' $D3/v2-a93.txt | sort) \
     <(awk '{print $2"|"$1}' $D2/ticket-abc3.txt | sort) | awk -F'|' '$2!=$3{print "DIFF",$1,$2,$3}' | head -30
