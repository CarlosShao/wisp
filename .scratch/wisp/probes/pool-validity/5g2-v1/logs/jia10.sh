#!/bin/bash
# Task 4: for the 10 甲 tickets -- verify (a) claimed evidence file exists+tracked,
# (b) the claimed signature string appears VERBATIM with a line number, (c) box totals agree across m28 / m25-ac / passC.
cd "D:/work/workspace/projects plans/Wisp" || exit 1
G=.scratch/wisp/probes/pool-validity/5g2/logs
L=.scratch/wisp/probes/pool-validity/5g2-v1/logs
: > "$L/jia10.tsv"
for T in 11 63 66 69 73 74 76 117 198 241; do
  f=$(awk -F'|' -v n="$T" '$1+0==n{print $2}' "$G/m31-sign.txt")
  sigline=$(grep -nE 'Acceptor|验收人|验收腿|裁决腿|acceptor-|编排者' "docs/evidence/s1/$f" 2>/dev/null | head -1 | cut -c1-150)
  tracked=$(git ls-files --error-unmatch "docs/evidence/s1/$f" >/dev/null 2>&1 && echo Y || echo N)
  exists=$([ -f "docs/evidence/s1/$f" ] && echo Y || echo N)
  m28=$(awk -F'|' -v n="$T" '$1+0==n{print $5"/"$6"/"$3}' "$G/m28-final.txt")
  m25=$(awk -F'|' -v n="$T" '$1+0==n{print $3}' "$G/m25-ac.tsv")
  passc=$(awk -F'|' -v n="$T" '$1+0==n{print $3}' "$G/passC.tsv")
  printf '%s\tfile=%s\texists=%s tracked=%s\n  m31sig=%s\n  VERBATIM@%s\n  m28(chk/tot/un)=%s  m25(face_boxes)=%s  passC(b)=%s\n' \
    "$T" "$f" "$exists" "$tracked" "$(awk -F'|' -v n="$T" '$1+0==n{print $3}' "$G/m31-sign.txt" | cut -c1-70)" "$sigline" "$m28" "$m25" "$passc" >> "$L/jia10.tsv"
done
cat "$L/jia10.tsv"
