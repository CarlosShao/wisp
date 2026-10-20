#!/bin/sh
# 342-r2: UTF-8 integrity + codepoint sweep over this leg's authored evidence, captured to
# logs/utf8-and-codepoints.txt. Two rulers on purpose:
#   * iconv -f UTF-8 -t UTF-8  = "does the whole file decode" (catches a byte-cut CJK character, which
#     a grep that tolerates invalid bytes and a byte-level codepoint scan both MISS - measured: my own
#     evidence file carried exactly that defect and only iconv reported it);
#   * the codepoint ruler (sh/codepoint-ruler.pl) = the banned codepoint U+00A7 and the neighbouring
#     U+00B7 per file, so the two are never confused.
cd "/d/work/workspace/projects plans/Wisp" || exit 1
OUT=.scratch/wisp/probes/342/r2/logs/utf8-and-codepoints.txt
{
  echo "== iconv decode ruler, authored files =="
  cd .scratch/wisp/probes/342/r2 || exit 1
  bad=0
  for f in report.md logs/bench-probes-verbatim.txt logs/gates-and-roots.txt logs/appended-lines.txt \
           sh/codepoint-ruler.pl sh/commit-evidence.sh sh/gates-and-roots.sh sh/roots-and-faults.sh \
           sh/self-proof-rulers.sh sh/final-sweep.sh sh/utf8-and-codepoints.sh; do
    if iconv -f UTF-8 -t UTF-8 "$f" > /dev/null 2>&1; then echo "OK  $f"; else echo "BAD $f"; bad=$((bad+1)); fi
  done
  echo "UTF8_BAD=$bad"
  echo
  echo "== raw go-test logs (27 files, verbatim captures) =="
  rawbad=0
  for f in logs/*.log; do
    iconv -f UTF-8 -t UTF-8 "$f" > /dev/null 2>&1 || { echo "RAWBAD $f"; rawbad=$((rawbad+1)); }
  done
  echo "RAW_UTF8_BAD=$rawbad of $(ls logs/*.log | wc -l)"
  echo
  echo "== codepoint ruler (banned U+00A7 must be 0; U+00B7 printed so the two are not conflated) =="
  perl sh/codepoint-ruler.pl report.md logs/bench-probes-verbatim.txt logs/gates-and-roots.txt \
    logs/appended-lines.txt
  echo
  echo "== positive control on the same codepoint ruler (a file that DOES carry U+00A7) =="
  cd "/d/work/workspace/projects plans/Wisp" || exit 1
  perl .scratch/wisp/probes/342/r2/sh/codepoint-ruler.pl docs/reports/HANDOVER.md
  echo
  echo "== inventory / caps =="
  cd .scratch/wisp/probes/342/r2 || exit 1
  echo "PROBE_ROOT_FILES=$(find . -type f | wc -l)"
  echo "REPORT_MD_LINES=$(wc -l < report.md)"
  echo "ZERO_BYTE=$(find . -type f -size 0 | wc -l)"
  echo "OUT_EXTENSION=$(find . -name '*.out' | wc -l)"
} > "$OUT" 2>&1
echo "WROTE $OUT"
cat "$OUT"
