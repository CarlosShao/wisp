#!/bin/bash
# passB: per-ticket evidence-location census for the 52-ticket un-censused roster.
# Every finding lands in a file; only summaries go to stdout.
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
Roster=$L/roster52.txt
Out=$L/passB.tsv
Cand=$L/passB-cands.tsv
: > "$Out"; : > "$Cand"
SCOPE=(docs/evidence .scratch/wisp/probes docs/reports)
while read -r N; do
  [ -z "$N" ] && continue
  FILE=$(ls .scratch/wisp/issues/"$N"-*.md 2>/dev/null | head -1)
  if [ -z "$FILE" ]; then
    echo "$N|NO_TICKET_FILE_ON_DISK" >> "$Out"; continue
  fi
  BASE=$(basename "$FILE")
  TLINES=$(wc -l < "$FILE")
  BXI=$(grep -c '^[[:space:]]*- \[ \]' "$FILE"); BXA=$(grep -c '^- \[ \]' "$FILE")
  BXDONEI=$(grep -c '^[[:space:]]*- \[x\]' "$FILE"); BXDONEA=$(grep -c '^- \[x\]' "$FILE")
  BXTOT=$(grep -c '^[[:space:]]*- \[' "$FILE")
  HIT=$(grep -cE "(票|ticket|Ticket|Ticket) ?$N([^0-9]|$)|/$N-" "$FILE")
  # progress-log section: from the Progress log header to EOF
  PLSTART=$(grep -nE '^#{2,3} *(Progress log|进度|进展)' "$FILE" | head -1 | cut -d: -f1)
  if [ -n "$PLSTART" ]; then
    PL=$(tail -n +"$PLSTART" "$FILE")
    PLV=$(printf '%s\n' "$PL" | grep -cE 'PASS|FAIL|裁决|验收方|acceptor|不通过|通过')
    PLN=$(printf '%s\n' "$PL" | wc -l)
  else
    PLV=-1; PLN=-1
  fi
  EN=$(grep -cE '^#{2,3} .*(Evidence|凭据|裁决|验收|归口)' "$FILE")
  echo "$N|$BASE|lines=$TLINES|box_unind=$BXA/box_ind=$BXI/boxdone_ind=$BXDONEI/boxdone_unind=$BXDONEA/boxtot=$BXTOT|tikhits=$HIT|proglog_lines=$PLN/proglog_verdictlines=$PLV|evidsec=$EN" >> "$Out"
  CANDS=$(grep -rlE "票 ?$N([^0-9]|$)|[Tt]icket ?$N([^0-9]|$)|/$N-|${N}-adversarial|${N}-acceptance" "${SCOPE[@]}" 2>/dev/null | grep -v "^$FILE$")
  NC=$(printf '%s\n' "$CANDS" | grep -c .)
  echo "$N|$BASE|cands=$NC" >> "$Out"
  printf '%s\n' "$CANDS" | while read -r C; do
    [ -z "$C" ] && continue
    CL=$(wc -l < "$C")
    CT=$(grep -m1 '^#' "$C" | cut -c1-70)
    [ -z "$CT" ] && CT=$(head -1 "$C" | cut -c1-70)
    RI=$(grep -cE '实现者|落地腿|实现方|自述|self-report' "$C")
    RA=$(grep -cE '验收方|裁决|acceptor|非实现者|独立验收|对抗验收' "$C")
    RV=$(grep -cE 'PASS|FAIL|不通过|通过|拒绝|reject' "$C")
    printf '%s|%s|lines=%s|impl=%s|acc=%s|verdict=%s|%s\n' "$N" "$C" "$CL" "$RI" "$RA" "$RV" "$CT" >> "$Cand"
  done
done < "$Roster"
echo "OUT=$(wc -l < $Out) CAND=$(wc -l < $Cand)"
