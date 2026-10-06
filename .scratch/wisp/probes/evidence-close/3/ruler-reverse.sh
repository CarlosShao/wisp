#!/usr/bin/env bash
# 反向尺（本腿自加，死腿没跑过这一把）：v3 判 OTHER 的同名件里，
# 标题带终裁词族、且前 40 行没有实现侧自称/免责句 ⇒ 候选"错杀"（其实是终裁表被念成旁证）。
# 与正向尺（VERDICT 件里查实现侧自称）配成一对，防的是两个方向的漏。
set -u
cd "D:/work/workspace/projects plans/Wisp"
D3=.scratch/wisp/probes/evidence-close/3
VERD='对抗验收|裁决|验收腿|验收程|验收方|终裁|acceptance|acceptor|非实现者|独立复算|独立验收|独立验证|抽验|回判|reaccept|ratification|accept-r'
IMPL='实现侧交件|实现者自证|实现程自证|实现件|证据件|交件证据|预做|普查|读数表|凭据索引|只读取证|只读设计核|正向对照|补证|写码位|落地写腿|修复程|实现程|实现方|归因尺|存量盘点|写码代理'
DISCLAIM='本件不是验收件|不是验收件|本文件不是裁决表|本件不是裁决表|不是对抗验收表|同一程写码又自证|本件不声称结掉'
: > $D3/false-veto-candidates.txt
grep '|OTHER' $D3/v3-file-identity.txt | cut -d'|' -f1,2 | while IFS='|' read -r n f; do
  [ -z "$f" ] && continue
  t=$(head -1 "docs/evidence/s1/$f" | tr -d '\r')
  printf '%s' "$t" | grep -qE "$VERD" || continue
  body=$(head -40 "docs/evidence/s1/$f" | tr -d '\r')
  printf '%s' "$body" | grep -qE "$DISCLAIM" && continue
  # 标题自身含实现侧词 ⇒ 不算错杀（那是实现件/证据件的自我命名）
  printf '%s' "$t" | grep -qE "$IMPL" && continue
  echo "$n|$f|$(printf '%s' "$t" | cut -c1-70)" >> $D3/false-veto-candidates.txt
done
echo "--- 候选错杀件数 ---"; wc -l < $D3/false-veto-candidates.txt
echo "--- 其中该票名下另有真 VERDICT 件（桶不变，只是名册少列一枚） ---"
while IFS='|' read -r n f rest; do
  if grep -q "^$n|$f|VERDICT" $D3/v3-file-identity.txt; then echo "?? $n $f"; fi
  if grep "^$n|" $D3/v3-file-identity.txt | grep -q '|VERDICT'; then echo "ALSO-A $n $f"; else echo "ONLY-OTHER $n $f  <== 会改桶"; fi
done < $D3/false-veto-candidates.txt
