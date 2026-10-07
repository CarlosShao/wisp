#!/bin/bash
# 乙 test: a ticket is 乙 ("只有实现方自述") iff EVERY 名下 file self-declares implementer authorship
#   patterns: 本程角色：**实现程** | 实现腿：`NN-r1`（本文件作者） | 本件作者＝实现
# Also count mixed tickets. Positive control: must hit a known real name (149).
cd "D:/work/workspace/projects plans/Wisp" || exit 1
G=.scratch/wisp/probes/pool-validity/5g2/logs
L=.scratch/wisp/probes/pool-validity/5g2-v1/logs
IMPL='本程角色：\*\*实现程|实现腿：`[0-9]{2,3}-[a-z]?[0-9]*`（本文件作者|实现腿：`[^`]*`（本文件作者|本件作者＝\*\*实现|作者＝实现程|本文件作者） · |（本腿＝实现'
: > "$L/implself.tsv"
while read -r N; do
  files=$(awk -F'|' -v n="$N" '$1+0==n{print substr($0,index($0,"|")+1)}' "$G/k-accfiles.txt" | tr ' ' '\n' | grep -E '^(docs|\.scratch)')
  tot=0; imp=0; names=""
  for F in $files; do
    [ -f "$F" ] || continue
    tot=$((tot+1))
    h=$(head -30 "$F")
    if printf '%s\n' "$h" | grep -qE '本程角色：\*\*实现程|实现腿：`[^`]*`（本文件作者|本件作者＝\*\*实现|作者＝实现程|本程角色：实现'; then
      imp=$((imp+1)); names="$names $(basename "$F")"
    fi
  done
  if [ "$tot" -gt 0 ]; then
    if [ "$imp" -eq "$tot" ]; then v="ALL_IMPL"; elif [ "$imp" -gt 0 ]; then v="MIXED"; else v="NONE"; fi
    printf '%s\t%s\t%s\t%s\t%s\n' "$N" "$tot" "$imp" "$v" "$names" >> "$L/implself.tsv"
  fi
done < <(sort -n -u "$G/roster52.txt")
echo "== breakdown =="; cut -f4 "$L/implself.tsv" | sort | uniq -c
echo "== tickets with 0 files in k-accfiles =="; awk -F'\t' '$2==0{printf "%s ",$1}' "$L/implself.tsv"
echo "== ALL_IMPL =="; awk -F'\t' '$4=="ALL_IMPL"{printf "%s(%s/%s) ",$1,$3,$2}' "$L/implself.tsv"
echo; echo "== MIXED =="; awk -F'\t' '$4=="MIXED"{printf "%s(%s/%s)%s ",$1,$3,$2,$5}' "$L/implself.tsv"
echo; echo "== positive control: 149 row =="; grep -P '^149\t' "$L/implself.tsv"
echo "== nI column from m19 for the flagged =="; for t in $(awk -F'\t' '$4!="NONE"{print $1}' "$L/implself.tsv" | tr '\n' ' '); do echo -n "$t:m19nI=$(awk -F'|' -v n="$t" '$1+0==n{print $6}' "$G/m19-table.tsv") "; done
