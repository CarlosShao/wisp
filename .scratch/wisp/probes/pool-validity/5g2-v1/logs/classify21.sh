#!/bin/bash
# Classify the 21 待定 tickets that m32-who.tsv reads authorship for:
#  ORCH = matched an orchestrator phrase
#  LEGID = matched 验收腿/裁决腿/本程代号/验收程  (leg-role words)
#  BARENONIMPL = matched ONLY the bare word 非实现者 (no name, no orchestrator phrase)
cd "D:/work/workspace/projects plans/Wisp" || exit 1
G=.scratch/wisp/probes/pool-validity/5g2/logs
L=.scratch/wisp/probes/pool-validity/5g2-v1/logs
ORCH='Acceptor:\*\* ?orchestrator|验收人\*\*：编排者|验收人：编排者|编排者亲自|编排者自裁|编排者执行|由编排者'
LEGW='本程代号|验收腿|裁决腿|验收程'
: > "$L/classify21.tsv"
for T in 67 79 83 119 141 142 143 144 147 149 151 152 153 154 155 157 162 183 212 221 257; do
  files=$(awk -F'|' -v n="$T" '$1+0==n{print substr($0,index($0,"|")+1)}' "$G/k-accfiles.txt" | tr ' ' '\n' | grep -E '^(docs|\.scratch)')
  orch=0; legw=0; bare=0; hitfile=""; hitline=""
  for F in $files; do
    [ -f "$F" ] || continue
    h=$(head -30 "$F")
    o=$(printf '%s\n' "$h" | grep -coE "$ORCH")
    l=$(printf '%s\n' "$h" | grep -coE "$LEGW")
    b=$(printf '%s\n' "$h" | grep -coE '非实现者')
    [ "$o" -gt 0 ] && orch=1
    [ "$l" -gt 0 ] && legw=1
    if [ "$o" -eq 0 ] && [ "$l" -eq 0 ] && [ "$b" -gt 0 ]; then bare=1; fi
    if [ -z "$hitfile" ] && [ $((o+l+b)) -gt 0 ]; then
      hitfile="$F"
      hitline=$(printf '%s\n' "$h" | grep -nE "$ORCH|$LEGW|非实现者" | head -1 | cut -c1-90)
    fi
  done
  cat1="?"
  if [ "$orch" -eq 1 ] && [ "$legw" -eq 1 ]; then cat1="ORCH+LEGWORD"
  elif [ "$orch" -eq 1 ]; then cat1="ORCH"
  elif [ "$legw" -eq 1 ]; then cat1="LEGWORD"
  elif [ "$bare" -eq 1 ]; then cat1="BARE_ONLY"
  fi
  printf '%s\t%s\t%s\t%s\n' "$T" "$cat1" "$hitfile" "$hitline" >> "$L/classify21.tsv"
done
awk -F'\t' '{c[$2]++} END{for(k in c) print k, c[k]}' "$L/classify21.tsv" | sort
echo "--- rows: $(wc -l < "$L/classify21.tsv") ---"
