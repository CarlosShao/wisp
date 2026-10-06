#!/usr/bin/env bash
# evidence-close-4 收口尺复跑（只读；⛔ 不写 evidence-close/2、/3 的任何原件，只写本腿目录 /4）。
#
# 本腿选定的尺 = 死腿 evidence-close/3 的 v3 版（`rerun-v3.sh` 里那把），逐字沿用其口径：
#   同名件 := `ls docs/evidence/s1/ | grep -E "^0*<票号>([-/]|$)"`
#             ——票号必须在文件名开头，且后面紧跟 `-`、`/` 或行尾。
#             ⇒ `136-ac14-*` 归票 136（不是票 14）；`245-esc-*` 归票 245（不是票 45／票 5）。
#   VERDICT := 第 1 行标题命中 VERDICT 词族 **且** 标题行不命中实现侧自称词族
#              **且** 前 8 行没有明写的免责句（"本件不是验收件／不是裁决表／同一程写码又自证"）。
#   OTHER   := 其余同名件（实现侧交件／普查／读数台件／取证件／索引／目录）。
#   逐票三档：A 名下有 VERDICT ／ B 名下只有 OTHER ／ C 名下零枚同名件（NOFILE）。
#
# 输入只读：docs/evidence/s1/（名册与标题）＋ .scratch/wisp/issues/*-done.md（分母）。
set -u
cd "D:/work/workspace/projects plans/Wisp"
D4=.scratch/wisp/probes/evidence-close/4
E=docs/evidence/s1

# ---------- 分母：本腿现量（⛔ 不含 README.md、不含任何非 -done 票） ----------
ls .scratch/wisp/issues/*-done.md | sed 's|.*/||' > $D4/done-files-4.txt
awk -F'-' '{print $1}' $D4/done-files-4.txt | sort -u > $D4/nums4.txt
wc -l < $D4/nums4.txt | tr -d ' ' > $D4/denominator.txt

VERD='对抗验收|裁决|验收腿|验收程|验收方|验收件|终裁|acceptance|acceptor|非实现者|独立复算|独立验收|独立验证|抽验|回判|reaccept|ratification|accept-r'
IMPLTITLE='实现侧交件|实现者自证|实现程自证|实现件|证据件|交件证据|预做|普查|读数表|凭据索引|只读取证|只读设计核|正向对照|补证|写码位|落地写腿|修复程|实现程|实现方|归因尺|存量盘点'
DISCLAIM='本件不是验收件|不是验收件|本文件不是裁决表|本件不是裁决表|不是对抗验收表|同一程写码又自证|本件不声称结掉|本文件是\*\*实现者自证'

# ---------- 同名件名册：本腿自己现量（不引用死腿的 pairs） ----------
: > $D4/pairs-own.txt
while read -r n; do
  ls $E/ 2>/dev/null | grep -E "^0*${n}([-/]|$)" | while read -r f; do echo "$n $f"; done
done < $D4/nums4.txt >> $D4/pairs-own.txt

# ---------- 逐件分类 ----------
: > $D4/file-identity-4.txt
while read -r n f; do
  [ -d "$E/$f" ] && { echo "$n|$f|DIR|目录（台件，非裁决表）" >> $D4/file-identity-4.txt; continue; }
  [ -f "$E/$f" ] || { echo "$n|$f|MISSING|原件不在盘" >> $D4/file-identity-4.txt; continue; }
  t=$(head -1 "$E/$f" | tr -d '\r')
  d=$(head -8 "$E/$f" | tr -d '\r' | grep -m1 -E "$DISCLAIM")
  verdict=N
  printf '%s' "$t" | grep -qE "$VERD" && ! printf '%s' "$t" | grep -qE "$IMPLTITLE" && verdict=Y
  [ -n "$d" ] && verdict=N
  if [ "$verdict" = Y ]; then echo "$n|$f|VERDICT|${t:0:60}" >> $D4/file-identity-4.txt
  else why=$(printf '%s' "$t" | grep -oE "$IMPLTITLE" | head -1)
       echo "$n|$f|OTHER|${why:-标题无终裁词}" >> $D4/file-identity-4.txt; fi
done < $D4/pairs-own.txt

# ---------- 逐票三档 ----------
: > $D4/ticket-abc-4.tsv
while read -r n; do
  v=$(awk -F'|' -v n="$n" '$1==n && $3=="VERDICT"{print $2}' $D4/file-identity-4.txt | tr '\n' ' ')
  o=$(awk -F'|' -v n="$n" '$1==n && $3!="VERDICT"{print $2}' $D4/file-identity-4.txt | tr '\n' ' ')
  if [ -n "$v" ]; then printf 'A\t%s\t%s\t%s\n' "$n" "$v" "$o" >> $D4/ticket-abc-4.tsv
  elif [ -n "$o" ]; then printf 'B\t%s\t\t%s\n' "$n" "$o" >> $D4/ticket-abc-4.tsv
  else printf 'C\t%s\t\t\n' "$n" >> $D4/ticket-abc-4.tsv; fi
done < $D4/nums4.txt

echo "--- 本腿复跑三档计数 ---"
for b in A B C; do printf '%s=%s ' "$b" "$(grep -c "^$b	" $D4/ticket-abc-4.tsv)"; done; echo
echo "--- 与死腿 v3-a93（93 枚分母）逐票差集 ---"
join -t'|' <(awk -F'\t' '{print $2"|"$1}' $D4/ticket-abc-4.tsv | sort) \
     <(awk -F'\t' '$2!="268"{print $2"|"$1}' .scratch/wisp/probes/evidence-close/3/v3-a93.txt | sort) \
  | awk -F'|' '$2!=$3{print "DIFF ticket",$1,"mine="$2,"leg3="$3}'
echo "--- 死腿 rows-draft.tsv 缺的那一枚 ---"
comm -23 <(sort $D4/nums4.txt) <(awk -F'\t' '{print $1}' .scratch/wisp/probes/evidence-close/3/rows-draft.tsv | sort)
