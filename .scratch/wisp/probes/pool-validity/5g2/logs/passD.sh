#!/bin/bash
# passD: cheap per-ticket census -> filename, unchecked box counts (two rulers), box totals,
# progress-log verdict-role markers, and evidence-file references written on the ticket itself.
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
Out=$L/passD.tsv; : > "$Out"
while read -r N; do
  [ -z "$N" ] && continue
  FILE=$(ls .scratch/wisp/issues/"$N"-*.md 2>/dev/null | head -1)
  [ -z "$FILE" ] && { echo "$N|NO_FILE_ON_DISK" >> "$Out"; continue; }
  BASE=$(basename "$FILE")
  A=$(grep -c '^- \[ \]' "$FILE")                       # ruler 1: no indentation
  B=$(grep -c '^[[:space:]]*- \[ \]' "$FILE")           # ruler 2: any indentation
  DX=$(grep -c '^[[:space:]]*- \[[xX]\]' "$FILE")       # checked
  TOT=$(grep -c '^[[:space:]]*- \[' "$FILE")            # all boxes
  PL=$(grep -nE '^#{2,3} *(Progress log|进展|进度)' "$FILE" | head -1 | cut -d: -f1)
  [ -n "$PL" ] || PL=1
  PLV=$(tail -n +"$PL" "$FILE" | grep -cE '非实现者|对抗验收|验收|裁决|PASS|FAIL|通过')
  PLI=$(tail -n +"$PL" "$FILE" | grep -cE '实现腿|实现程|落地腿|写码位|写码腿|自述')
  EVREF=$(grep -ohE '(docs|\.scratch)/[A-Za-z0-9_./-]*\.md' "$FILE" | sort -u | tr '\n' ',' | cut -c1-70)
  EVEXIST=0; for C in $(grep -ohE '(docs|\.scratch)/[A-Za-z0-9_./-]*\.md' "$FILE" | sort -u); do test -f "$C" && EVEXIST=$((EVEXIST+1)); done
  printf '%s|%s|unind=%s|anyind=%s|done=%s|boxes=%s|pl_verdict=%s|pl_impl=%s|refmd=%s|refexist=%s\n' \
    "$N" "$BASE" "$A" "$B" "$DX" "$TOT" "$PLV" "$PLI" "$EVEXIST" "$EVREF" >> "$Out"
done < "$L/roster52.txt"
echo "PASSD=$(wc -l < "$Out")"
awk -F'|' '$5=="boxes=0"' "$Out" | wc -l
