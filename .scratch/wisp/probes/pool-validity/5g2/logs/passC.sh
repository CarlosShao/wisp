#!/bin/bash
# passC: classify evidence candidates -> one compact line per ticket + side lists.
cd "D:/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/pool-validity/5g2/logs
Tab=$L/passC.tsv; Zero=$L/passC-zeroverdict.txt; Nonzero=$L/passC-verdictfiles.txt
: > "$Tab"; : > "$Zero"; : > "$Nonzero"
VS='adversarial|acceptance|accept|verdict|ruling|裁决|验收|判定|结案|归口|backfill|终裁|缺口审计|gap-audit|gap_audit|audit'
while read -r N; do
  [ -z "$N" ] && continue
  FILE=$(ls .scratch/wisp/issues/"$N"-*.md 2>/dev/null | head -1)
  BASE=$(basename "$FILE")
  BXI=$(grep -c '^[[:space:]]*- \[ \]' "$FILE"); BXA=$(grep -c '^- \[ \]' "$FILE")
  BXTOT=$(grep -c '^[[:space:]]*- \[' "$FILE")
  PL=$(grep -nE '^#{2,3} *(Progress log|进度|进展)' "$FILE" | head -1 | cut -d: -f1)
  if [ -n "$PL" ]; then PLV=$(tail -n +"$PL" "$FILE" | grep -cE 'PASS|FAIL|裁决|验收|通过|拒绝|不通过')
  else PLV=nosec; fi
  EVF=""; EVN=0
  for C in $(grep -rlE "票 ?$N([^0-9]|$)|[Tt]icket ?$N([^0-9]|$)|/$N-|(^|/)${N}-" docs/evidence .scratch/wisp/probes docs/reports 2>/dev/null | grep -vF "$BASE"); do
    if echo "$C" | grep -qiE "$VS" || grep -m3 -qiE "$VS" "$C"; then
      EVN=$((EVN+1)); EVF="$EVF $C"
      echo "$N|$C|lines=$(wc -l < "$C")|$(grep -m1 '^#' "$C" | cut -c1-60)|$(grep -cE '实现者|落地腿|实现方' "$C")|$(grep -cE '验收方|acceptor|非实现者|独立' "$C")" >> "$Nonzero"
    fi
  done
  OWN=0; OWNSEC=$(grep -cE '^#{2,3} .*(Evidence|凭据|裁决|验收|归口)' "$FILE")
  for C in $(grep -oE '(docs|\.scratch)/[A-Za-z0-9_./-]*\.md' "$FILE" 2>/dev/null | sort -u); do
    test -f "$C" && OWN=$((OWN+1))
  done
  printf '%s|%s|b=%s/%s/%s|plv=%s|EV=%s|ownsec=%s|ownref_exist=%s|%s\n' "$N" "$BASE" "$BXA" "$BXI" "$BXTOT" "$PLV" "$EVN" "$OWNSEC" "$OWN" "$(echo $EVF | cut -c1-120)" >> "$Tab"
  [ "$EVN" = 0 ] && { echo "=== $N $BASE" >> "$Zero"; grep -rlE "票 ?$N([^0-9]|$)|[Tt]icket ?$N([^0-9]|$)|/$N-" docs/evidence .scratch/wisp/probes docs/reports 2>/dev/null | head -8 >> "$Zero"; }
done < "$L/roster52.txt"
echo "TAB=$(wc -l < $Tab) ZEROEV=$(wc -l < $Zero) VF=$(wc -l < $Nonzero)"
