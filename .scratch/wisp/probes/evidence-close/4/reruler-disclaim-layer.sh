#!/usr/bin/env bash
# evidence-close-4 §1.2b 复算尺：v3 那把尺的"前 8 行免责句"层，在这 146 枚名册上究竟独立决定过几枚？
# 只读；输出只写本目录。⛔ 不改任何仓内件。
set -u
cd "D:/work/workspace/projects plans/Wisp"
D4=.scratch/wisp/probes/evidence-close/4
E=docs/evidence/s1

VERD3='对抗验收|裁决|验收腿|验收程|验收方|验收件|终裁|acceptance|acceptor|非实现者|独立复算|独立验收|独立验证|抽验|回判|reaccept|ratification|accept-r'
IMPL3='实现侧交件|实现者自证|实现程自证|实现件|证据件|交件证据|预做|普查|读数表|凭据索引|只读取证|只读设计核|正向对照|补证|写码位|落地写腿|修复程|实现程|实现方|归因尺|存量盘点'
DISCLAIM='本件不是验收件|不是验收件|本文件不是裁决表|本件不是裁决表|不是对抗验收表|同一程写码又自证|本件不声称结掉'

# 逐枚：标题层结论 vs 免责句命中
while read -r f; do
  [ -d "$E/$f" ] && { echo "$f :: DIR-skip"; continue; }
  t=$(head -1 "$E/$f" | tr -d '\r')
  v=$(printf '%s' "$t" | grep -oE "$VERD3" | head -1)
  i=$(printf '%s' "$t" | grep -oE "$IMPL3" | head -1)
  d=$(head -8 "$E/$f" | tr -d '\r' | grep -m1 -oE "$DISCLAIM")
  [ -z "$d" ] && continue
  if [ -n "$i" ]; then why="标题层已因 IMPL3[$i] 落 OTHER ⇒ 免责句层无活可接"
  elif [ -z "$v" ]; then why="标题层因'无终裁词'(VERD3 不中) 落 OTHER ⇒ 免责句层无活可接"
  else why="★免责句层独立决定（VERD3[$v] 中、IMPL3 不中、免责句[$d]）"; fi
  echo "$f :: dis=[$d] :: $why"
done < <(awk '{print $2}' $D4/pairs-own.txt | sort -u)
